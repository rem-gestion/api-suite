package clients

import (
	"context"
	"fmt"
	"time"

	rcgrpc "github.com/rem-gestion/rem-common/grpc"
	authpb "github.com/rem-gestion/rem-common/protos/auth-identity/v1"
	"google.golang.org/grpc"
)

// AuthIdentityClient maneja la comunicación gRPC con auth-identity-svc
type AuthIdentityClient struct {
	conn   *grpc.ClientConn
	client authpb.AuthIdentityServiceClient
	addr   string
}

// NewAuthIdentityClient crea una nueva instancia del cliente
func NewAuthIdentityClient(addr string) (*AuthIdentityClient, error) {
	conn, err := rcgrpc.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth-identity service at %s: %w", addr, err)
	}

	client := authpb.NewAuthIdentityServiceClient(conn)

	return &AuthIdentityClient{
		conn:   conn,
		client: client,
		addr:   addr,
	}, nil
}

// NewAuthIdentityClientWithConn crea una nueva instancia del cliente usando una conexión existente
func NewAuthIdentityClientWithConn(conn *grpc.ClientConn) *AuthIdentityClient {
	client := authpb.NewAuthIdentityServiceClient(conn)

	return &AuthIdentityClient{
		conn:   conn,
		client: client,
		addr:   "existing_connection",
	}
}

// Close cierra la conexión gRPC
func (c *AuthIdentityClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ValidateToken valida un token JWT
func (c *AuthIdentityClient) ValidateToken(token string) (*authpb.ValidateTokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &authpb.ValidateTokenRequest{
		Token: token,
	}

	return c.client.ValidateToken(ctx, req)
}

// GetUserById obtiene un usuario por ID
func (c *AuthIdentityClient) GetUserById(userID string) (*authpb.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &authpb.GetUserByIdRequest{
		UserId: userID,
	}

	return c.client.GetUserById(ctx, req)
}

// GetUserByEmail obtiene un usuario por email
func (c *AuthIdentityClient) GetUserByEmail(email string) (*authpb.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &authpb.GetUserByEmailRequest{
		Email: email,
	}

	return c.client.GetUserByEmail(ctx, req)
}

// CreateUser crea un nuevo usuario
func (c *AuthIdentityClient) CreateUser(req *authpb.CreateUserRequest) (*authpb.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return c.client.CreateUser(ctx, req)
}

// UpdateUser actualiza un usuario existente
func (c *AuthIdentityClient) UpdateUser(req *authpb.UpdateUserRequest) (*authpb.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return c.client.UpdateUser(ctx, req)
}

// DeactivateUser desactiva un usuario
func (c *AuthIdentityClient) DeactivateUser(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &authpb.DeactivateUserRequest{
		UserId: userID,
	}

	_, err := c.client.DeactivateUser(ctx, req)
	return err
}

// GetAccount obtiene la información de cuenta de un usuario
func (c *AuthIdentityClient) GetAccount(userID string) (*authpb.AccountResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &authpb.GetAccountRequest{
		AccountId: userID,
	}

	return c.client.GetAccount(ctx, req)
}

// UpdateAccount actualiza la información de cuenta
func (c *AuthIdentityClient) UpdateAccount(req *authpb.UpdateAccountRequest) (*authpb.AccountResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return c.client.UpdateAccount(ctx, req)
}
