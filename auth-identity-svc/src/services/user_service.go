package services

import (
	"context"
	"fmt"
	"log"

	"github.com/rem-gestion/api-suite/auth-identity/src/dto"
	"github.com/rem-gestion/api-suite/auth-identity/src/models"
	"github.com/rem-gestion/api-suite/auth-identity/src/repository"
	"github.com/rem-gestion/rem-common/circuit"
	rerrors "github.com/rem-gestion/rem-common/errors"
	"github.com/rem-gestion/rem-common/saga"
	"github.com/rem-gestion/rem-common/validation"
)

// UserService maneja operaciones relacionadas con usuarios
type UserService struct {
	userRepo      repository.UserRepository
	accountRepo   repository.AccountRepository
	personClient  *PersonServiceClient
	validator     *validation.PersonValidator
	healthChecker *circuit.ServiceHealthChecker
}

// NewUserService crea una nueva instancia del servicio
func NewUserService(
	userRepo repository.UserRepository,
	accountRepo repository.AccountRepository,
	personClient *PersonServiceClient,
) *UserService {
	healthChecker := circuit.NewServiceHealthChecker()

	// Registrar health checker para person-svc si está disponible
	if personClient != nil {
		healthChecker.RegisterService("person-svc", &PersonServiceHealthChecker{client: personClient})
	}

	return &UserService{
		userRepo:      userRepo,
		accountRepo:   accountRepo,
		personClient:  personClient,
		validator:     validation.NewPersonValidator(),
		healthChecker: healthChecker,
	}
}

// UpdateProfile actualiza el perfil del usuario con validación exhaustiva y Saga
func (s *UserService) UpdateProfile(userID string, req dto.UpdateProfileRequest) (*dto.UserProfileResponse, error) {
	ctx := context.Background()

	// Obtener usuario actual
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
		return nil, &rerrors.InternalServerError{Msg: "error getting user"}
	}
	if user == nil {
		return nil, &rerrors.NotFoundError{Msg: "user not found"}
	}

	// Validar datos de persona si se proporcionan
	if req.PersonData != nil {
		if err := ValidatePersonDataHelper(s.validator, req.PersonData); err != nil {
			return nil, &rerrors.BadRequestError{Msg: fmt.Sprintf("person data validation failed: %v", err)}
		}
	}

	// Health check antes de operaciones críticas
	if s.personClient != nil && req.PersonData != nil {
		if err := s.healthChecker.CheckService(ctx, "person-svc"); err != nil {
			log.Printf("Person service health check failed: %v", err)
			return nil, &rerrors.InternalServerError{Msg: "person service temporarily unavailable"}
		}
	}

	// Si hay datos de persona, usar Saga para consistencia transaccional
	if req.PersonData != nil {
		return s.updateProfileWithSaga(ctx, user, req)
	}

	// Si no hay datos de persona, actualización simple
	if req.PersonData == nil && user.OnboardStatus == models.OnboardNew {
		if err := s.userRepo.UpdateOnboardStatus(userID, models.OnboardInProgress); err != nil {
			log.Printf("Error updating onboard status: %v", err)
			// No fallar por esto
		}
		user.OnboardStatus = models.OnboardInProgress
	}

	// Retornar perfil actualizado
	return s.buildUserProfile(user)
}

// updateProfileWithSaga actualiza el perfil usando patrón Saga para consistencia
func (s *UserService) updateProfileWithSaga(ctx context.Context, user *models.User, req dto.UpdateProfileRequest) (*dto.UserProfileResponse, error) {
	sagaInstance := saga.NewSaga()
	var personID string
	var originalPersonID *string

	if user.PersonID == nil {
		// Caso 1: Crear nueva persona
		sagaInstance.AddStep(saga.Step{
			Name: "create_person",
			Execute: func(ctx context.Context) (interface{}, error) {
				if s.personClient == nil {
					return nil, fmt.Errorf("person service client not available")
				}

				personResp, err := s.personClient.CreatePerson(req.PersonData)
				if err != nil {
					return nil, fmt.Errorf("failed to create person: %w", err)
				}

				personID = personResp.Id
				return personResp.Id, nil
			},
			Rollback: func(ctx context.Context, result interface{}) error {
				if pid, ok := result.(string); ok && s.personClient != nil {
					// TODO: Implementar DeletePerson en PersonServiceClient para rollback completo
					log.Printf("TODO: Rollback person creation for ID: %s", pid)
				}
				return nil
			},
		})

		sagaInstance.AddStep(saga.Step{
			Name: "update_user_person_id",
			Execute: func(ctx context.Context) (interface{}, error) {
				if err := s.userRepo.UpdatePersonID(user.ID, personID); err != nil {
					return nil, fmt.Errorf("error updating user person ID: %w", err)
				}
				return personID, nil
			},
			Rollback: func(ctx context.Context, result interface{}) error {
				// Rollback: quitar la referencia a la persona
				return s.userRepo.UpdatePersonID(user.ID, "")
			},
		})
	} else {
		// Caso 2: Actualizar persona existente (placeholder - necesita implementación)
		originalPersonID = user.PersonID
		log.Printf("TODO: Implementar actualización de persona existente para ID: %s", *originalPersonID)
		// Por ahora, solo actualizamos el estado de onboarding
	}

	// Paso final: Actualizar estado de onboarding
	sagaInstance.AddStep(saga.Step{
		Name: "update_onboard_status",
		Execute: func(ctx context.Context) (interface{}, error) {
			if user.OnboardStatus == models.OnboardNew {
				if err := s.userRepo.UpdateOnboardStatus(user.ID, models.OnboardInProgress); err != nil {
					return nil, fmt.Errorf("error updating onboard status: %w", err)
				}
				return models.OnboardInProgress, nil
			}
			return user.OnboardStatus, nil
		},
		Rollback: func(ctx context.Context, result interface{}) error {
			// Rollback: restaurar estado anterior
			return s.userRepo.UpdateOnboardStatus(user.ID, user.OnboardStatus)
		},
	})

	// Ejecutar saga
	if err := sagaInstance.Execute(ctx); err != nil {
		log.Printf("Saga execution failed: %v", err)
		return nil, &rerrors.InternalServerError{Msg: "profile update failed"}
	}

	// Actualizar objeto user local para respuesta
	if personID != "" {
		user.PersonID = &personID
		user.OnboardStatus = models.OnboardInProgress
	}

	// Retornar perfil actualizado
	return s.buildUserProfile(user)
}

// GetUserByEmail obtiene un usuario por email
func (s *UserService) GetUserByEmail(email string) (*dto.UserProfileResponse, error) {
	account, err := s.accountRepo.GetByEmail(email)
	if err != nil {
		log.Printf("Error getting account by email: %v", err)
		return nil, &rerrors.InternalServerError{Msg: "error getting user"}
	}
	if account == nil {
		return nil, &rerrors.NotFoundError{Msg: "user not found"}
	}

	user, err := s.userRepo.GetByAccountID(account.ID)
	if err != nil {
		log.Printf("Error getting user by account ID: %v", err)
		return nil, &rerrors.InternalServerError{Msg: "error getting user"}
	}
	if user == nil {
		return nil, &rerrors.NotFoundError{Msg: "user not found"}
	}

	return s.buildUserProfile(user)
}

// ListUsers obtiene una lista de usuarios con filtros
func (s *UserService) ListUsers(page, perPage int, filter *repository.UserFilter) (*dto.ListUsersResponse, error) {
	users, total, err := s.userRepo.List(page, perPage, filter)
	if err != nil {
		log.Printf("Error listing users: %v", err)
		return nil, &rerrors.InternalServerError{Msg: "error listing users"}
	}

	// Convertir a DTOs
	userProfiles := make([]*dto.UserProfileResponse, len(users))
	for i, user := range users {
		profile, err := s.buildUserProfile(user)
		if err != nil {
			log.Printf("Error building user profile: %v", err)
			continue // Skip this user but continue with others
		}
		userProfiles[i] = profile
	}

	return &dto.ListUsersResponse{
		Users:   userProfiles,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	}, nil
}

// DeactivateUser desactiva un usuario
func (s *UserService) DeactivateUser(userID, updatedBy string) error {
	// Obtener usuario
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
		return &rerrors.InternalServerError{Msg: "error getting user"}
	}
	if user == nil {
		return &rerrors.NotFoundError{Msg: "user not found"}
	}

	// Desactivar cuenta
	if err := s.accountRepo.UpdateStatus(user.AccountID, models.StatusSuspended); err != nil {
		log.Printf("Error updating account status: %v", err)
		return &rerrors.InternalServerError{Msg: "error deactivating user"}
	}

	return nil
}

// buildUserProfile construye un perfil de usuario completo
func (s *UserService) buildUserProfile(user *models.User) (*dto.UserProfileResponse, error) {
	profile := &dto.UserProfileResponse{
		UserID:        user.ID,
		OnboardStatus: string(user.OnboardStatus),
		LastLogin:     user.LastLogin,
		CreatedAt:     user.CreatedAt,
	}

	// Agregar datos de account si están disponibles
	if user.Account != nil {
		profile.Email = user.Account.Email
		profile.Status = string(user.Account.Status)
	}

	// Agregar datos de persona si están disponibles
	if user.PersonID != nil {
		profile.PersonID = *user.PersonID

		// Obtener datos de persona desde person-svc
		if s.personClient != nil {
			personData, err := s.personClient.GetPerson(*user.PersonID)
			if err != nil {
				log.Printf("Error getting person data: %v", err)
				// No fallar por esto, solo no incluir los datos
			} else {
				profile.PersonData = personData
			}
		}
	}

	return profile, nil
}
