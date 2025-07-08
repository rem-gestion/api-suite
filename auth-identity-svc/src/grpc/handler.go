package grpc

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/rem-gestion/api-suite/auth-identity/src/dto"
	"github.com/rem-gestion/api-suite/auth-identity/src/services"
	authpb "github.com/rem-gestion/rem-common/protos/auth-identity/v1"
)

// AuthIdentityHandler implementa el servidor gRPC
type AuthIdentityHandler struct {
	authpb.UnimplementedAuthIdentityServiceServer
	authService *services.AuthService
	userService *services.UserService
}

// New crea una nueva instancia del handler gRPC
func New(authService *services.AuthService, userService *services.UserService) *AuthIdentityHandler {
	return &AuthIdentityHandler{
		authService: authService,
		userService: userService,
	}
}

// ValidateToken valida un token JWT
func (h *AuthIdentityHandler) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	claims, err := h.authService.ValidateToken(req.Token)
	if err != nil {
		log.Printf("Error validating token: %v", err)
		return &authpb.ValidateTokenResponse{
			IsValid: false,
		}, nil // No retornar error, solo indicar que no es válido
	}

	return &authpb.ValidateTokenResponse{
		IsValid:   true,
		UserId:    claims.UserID,
		Role:      claims.Role,
		ExpiresAt: claims.ExpiresAt.Unix(),
	}, nil
}

// GetUserById obtiene un usuario por ID
func (h *AuthIdentityHandler) GetUserById(ctx context.Context, req *authpb.GetUserByIdRequest) (*authpb.UserResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	profile, err := h.authService.GetUserProfile(req.UserId)
	if err != nil {
		log.Printf("Error getting user profile: %v", err)
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return h.buildUserResponse(profile), nil
}

// GetUserByEmail obtiene un usuario por email
func (h *AuthIdentityHandler) GetUserByEmail(ctx context.Context, req *authpb.GetUserByEmailRequest) (*authpb.UserResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	profile, err := h.userService.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return h.buildUserResponse(profile), nil
}

// CreateUser crea un nuevo usuario (para uso interno entre servicios)
func (h *AuthIdentityHandler) CreateUser(ctx context.Context, req *authpb.CreateUserRequest) (*authpb.UserResponse, error) {
	// TODO: Implementar cuando sea necesario para comunicación entre servicios
	return nil, status.Error(codes.Unimplemented, "create user via gRPC not implemented yet")
}

// UpdateUser actualiza un usuario
func (h *AuthIdentityHandler) UpdateUser(ctx context.Context, req *authpb.UpdateUserRequest) (*authpb.UserResponse, error) {
	// TODO: Implementar cuando sea necesario para comunicación entre servicios
	return nil, status.Error(codes.Unimplemented, "update user via gRPC not implemented yet")
}

// DeactivateUser desactiva un usuario
func (h *AuthIdentityHandler) DeactivateUser(ctx context.Context, req *authpb.DeactivateUserRequest) (*emptypb.Empty, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	err := h.userService.DeactivateUser(req.UserId, req.UpdatedBy)
	if err != nil {
		log.Printf("Error deactivating user: %v", err)
		return nil, status.Error(codes.Internal, "failed to deactivate user")
	}

	return &emptypb.Empty{}, nil
}

// GetAccount obtiene una cuenta por ID
func (h *AuthIdentityHandler) GetAccount(ctx context.Context, req *authpb.GetAccountRequest) (*authpb.AccountResponse, error) {
	// TODO: Implementar cuando sea necesario
	return nil, status.Error(codes.Unimplemented, "get account via gRPC not implemented yet")
}

// UpdateAccount actualiza una cuenta
func (h *AuthIdentityHandler) UpdateAccount(ctx context.Context, req *authpb.UpdateAccountRequest) (*authpb.AccountResponse, error) {
	// TODO: Implementar cuando sea necesario
	return nil, status.Error(codes.Unimplemented, "update account via gRPC not implemented yet")
}

// buildUserResponse construye la respuesta gRPC desde el perfil de usuario
func (h *AuthIdentityHandler) buildUserResponse(profile *dto.UserProfileResponse) *authpb.UserResponse {
	response := &authpb.UserResponse{
		Id:        profile.UserID,
		PersonId:  profile.PersonID,
		CreatedAt: timestamppb.New(profile.CreatedAt),
	}

	if profile.LastLogin != nil {
		response.LastLogin = timestamppb.New(*profile.LastLogin)
	}

	// Mapear estado de onboarding
	switch profile.OnboardStatus {
	case "new":
		response.OnboardStatus = authpb.OnboardStatus_NEW
	case "in_progress":
		response.OnboardStatus = authpb.OnboardStatus_IN_PROGRESS
	case "done":
		response.OnboardStatus = authpb.OnboardStatus_DONE
	default:
		response.OnboardStatus = authpb.OnboardStatus_ONBOARD_STATUS_UNSPECIFIED
	}

	// Si hay datos de cuenta, agregarlos
	if profile.Email != "" {
		account := &authpb.AccountResponse{
			Email:     profile.Email,
			CreatedAt: timestamppb.New(profile.CreatedAt),
		}

		// Mapear estado de cuenta
		switch profile.Status {
		case "pending":
			account.Status = authpb.AccountStatus_PENDING
		case "active":
			account.Status = authpb.AccountStatus_ACTIVE
		case "suspended":
			account.Status = authpb.AccountStatus_SUSPENDED
		default:
			account.Status = authpb.AccountStatus_ACCOUNT_STATUS_UNSPECIFIED
		}

		response.Account = account
	}

	return response
}
