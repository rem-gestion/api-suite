package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rem-gestion/api-suite/auth-identity/src/dto"
	"github.com/rem-gestion/api-suite/auth-identity/src/models"
	"github.com/rem-gestion/api-suite/auth-identity/src/repository"
	"github.com/rem-gestion/rem-common/circuit"
	rerrors "github.com/rem-gestion/rem-common/errors"
	"github.com/rem-gestion/rem-common/saga"
	"github.com/rem-gestion/rem-common/utils"
	"github.com/rem-gestion/rem-common/validation"
)

// AuthService maneja la lógica de negocio de autenticación
type AuthService struct {
	accountRepo     repository.AccountRepository
	userRepo        repository.UserRepository
	personClient    *PersonServiceClient
	jwtSecret       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	validator       *validation.PersonValidator
	healthChecker   *circuit.ServiceHealthChecker
}

// NewAuthService crea una nueva instancia del servicio
func NewAuthService(
	accountRepo repository.AccountRepository,
	userRepo repository.UserRepository,
	personClient *PersonServiceClient,
	jwtSecret string,
) *AuthService {
	healthChecker := circuit.NewServiceHealthChecker()

	// Registrar health checker para person-svc si está disponible
	if personClient != nil {
		healthChecker.RegisterService("person-svc", &PersonServiceHealthChecker{client: personClient})
	}

	return &AuthService{
		accountRepo:     accountRepo,
		userRepo:        userRepo,
		personClient:    personClient,
		jwtSecret:       jwtSecret,
		accessTokenTTL:  time.Hour * 24,      // 24 horas
		refreshTokenTTL: time.Hour * 24 * 30, // 30 días
		validator:       validation.NewPersonValidator(),
		healthChecker:   healthChecker,
	}
}

// PersonServiceHealthChecker implementa health checking para person-svc
type PersonServiceHealthChecker struct {
	client *PersonServiceClient
}

// Check verifica la salud del servicio person-svc
func (hc *PersonServiceHealthChecker) Check(ctx context.Context) error {
	if hc.client == nil {
		return fmt.Errorf("person service client not available")
	}

	// Intentar un health check básico - podríamos implementar un método específico en PersonServiceClient
	// Por ahora, asumimos que si el cliente existe, está saludable
	return nil
}

// Register registra un nuevo usuario
func (s *AuthService) Register(req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	ctx := context.Background()

	// Validar datos básicos
	if req.Email == "" {
		return nil, &rerrors.BadRequestError{Msg: "email is required"}
	}
	if req.Password == "" {
		return nil, &rerrors.BadRequestError{Msg: "password is required"}
	}
	if len(req.Password) < 6 {
		return nil, &rerrors.BadRequestError{Msg: "password must be at least 6 characters"}
	}

	// Validar persona si se proporcionan datos
	if req.PersonData != nil {
		log.Printf("DEBUG: PersonData received: %+v", req.PersonData)
		if err := ValidatePersonDataHelper(s.validator, req.PersonData); err != nil {
			return nil, &rerrors.BadRequestError{Msg: fmt.Sprintf("person data validation failed: %v", err)}
		}
	} else {
		log.Printf("DEBUG: No PersonData provided")
	}

	// Verificar si el cliente de persona está disponible
	if req.PersonData != nil && s.personClient == nil {
		log.Printf("WARNING: PersonData provided but person client is nil")
	}

	// Verificar si el email ya existe
	exists, err := s.accountRepo.Exists(req.Email)
	if err != nil {
		log.Printf("Error checking email existence: %v", err)
		return nil, &rerrors.InternalServerError{Msg: "error checking email availability"}
	}
	if exists {
		return nil, &rerrors.ConflictError{Msg: "email already registered"}
	}

	// Health check antes de operaciones críticas
	if s.personClient != nil {
		if err := s.healthChecker.CheckService(ctx, "person-svc"); err != nil {
			log.Printf("Person service health check failed: %v", err)
			// Continuar sin crear persona si el servicio no está disponible
		}
	}

	// Crear saga para transacción distribuida
	sagaInstance := saga.NewSaga()
	var accountID, userID, personID string

	// Paso 1: Crear cuenta
	sagaInstance.AddStep(saga.Step{
		Name: "create_account",
		Execute: func(ctx context.Context) (interface{}, error) {
			// Hash de la contraseña
			passwordHash, err := utils.HashPassword(req.Password)
			if err != nil {
				return nil, fmt.Errorf("error hashing password: %w", err)
			}

			// Determinar proveedor
			provider := models.ProviderEmail
			if req.Provider == "google" {
				provider = models.ProviderGoogle
			}

			// Crear Account
			account := &models.Account{
				Provider:     provider,
				Email:        req.Email,
				PasswordHash: passwordHash,
				Status:       models.StatusPending,
			}

			if err := s.accountRepo.Create(account); err != nil {
				return nil, fmt.Errorf("error creating account: %w", err)
			}

			accountID = account.ID
			return account.ID, nil
		},
		Rollback: func(ctx context.Context, result interface{}) error {
			if aid, ok := result.(string); ok {
				return s.accountRepo.Delete(aid)
			}
			return nil
		},
	})

	// Paso 2: Crear persona si se proporcionan datos
	if req.PersonData != nil {
		log.Printf("DEBUG: Adding create_person step to saga")
		sagaInstance.AddStep(saga.Step{
			Name: "create_person",
			Execute: func(ctx context.Context) (interface{}, error) {
				log.Printf("DEBUG: Executing create_person step")
				if s.personClient == nil {
					log.Printf("ERROR: person service client not available")
					return nil, errors.New("person service client not available")
				}

				log.Printf("DEBUG: Calling personClient.CreatePerson with data: %+v", req.PersonData)
				personResp, err := s.personClient.CreatePerson(req.PersonData)
				if err != nil {
					log.Printf("ERROR: Failed to create person via gRPC: %v", err)
					return nil, err // Ahora retornamos el error mapeado directamente
				}

				log.Printf("DEBUG: Person created successfully with ID: %s", personResp.Id)
				personID = personResp.Id
				return personResp.Id, nil
			},
			Rollback: func(ctx context.Context, result interface{}) error {
				if pid, ok := result.(string); ok && s.personClient != nil {
					// Aquí necesitaríamos implementar DeletePerson en PersonServiceClient
					log.Printf("TODO: Rollback person creation for ID: %s", pid)
				}
				return nil
			},
		})
	} else {
		log.Printf("DEBUG: Skipping create_person step - no PersonData provided")
	}

	// Paso 3: Crear usuario
	sagaInstance.AddStep(saga.Step{
		Name: "create_user",
		Execute: func(ctx context.Context) (interface{}, error) {
			log.Printf("DEBUG: Creating user with accountID=%s, personID=%s", accountID, personID)
			user := &models.User{
				AccountID:     accountID,
				OnboardStatus: models.OnboardNew,
			}

			if personID != "" {
				user.PersonID = &personID
				user.OnboardStatus = models.OnboardInProgress
				log.Printf("DEBUG: User will be linked to person ID: %s", personID)
			} else {
				log.Printf("DEBUG: User will be created without person link")
			}

			if err := s.userRepo.Create(user); err != nil {
				return nil, fmt.Errorf("error creating user: %w", err)
			}

			userID = user.ID
			log.Printf("DEBUG: User created successfully with ID: %s", userID)
			return user.ID, nil
		},
		Rollback: func(ctx context.Context, result interface{}) error {
			if uid, ok := result.(string); ok {
				return s.userRepo.Delete(uid)
			}
			return nil
		},
	})

	// Paso 4: Activar cuenta
	sagaInstance.AddStep(saga.Step{
		Name: "activate_account",
		Execute: func(ctx context.Context) (interface{}, error) {
			if err := s.accountRepo.UpdateStatus(accountID, models.StatusActive); err != nil {
				return nil, fmt.Errorf("error activating account: %w", err)
			}
			return "activated", nil
		},
		Rollback: func(ctx context.Context, result interface{}) error {
			return s.accountRepo.UpdateStatus(accountID, models.StatusSuspended)
		},
	})

	// Ejecutar saga
	if err := sagaInstance.Execute(ctx); err != nil {
		log.Printf("Saga execution failed: %v", err)

		// Extraer el error original de la saga para preservar el tipo
		errorMsg := err.Error()

		// Buscar patrones específicos que indican errores de conflicto
		if strings.Contains(errorMsg, "duplicate key value violates unique constraint") ||
			strings.Contains(errorMsg, "unique constraint") ||
			strings.Contains(errorMsg, "already exists") ||
			strings.Contains(errorMsg, "email already registered") {
			return nil, &rerrors.ConflictError{Msg: "resource already exists"}
		}

		// Buscar patrones que indican errores de validación
		if strings.Contains(errorMsg, "validation failed") ||
			strings.Contains(errorMsg, "is required") ||
			strings.Contains(errorMsg, "must be at least") ||
			strings.Contains(errorMsg, "invalid") {
			return nil, &rerrors.BadRequestError{Msg: "validation failed"}
		}

		// Para otros errores, mantener como error interno
		return nil, &rerrors.InternalServerError{Msg: "registration failed"}
	}

	// Generar tokens
	accessToken, err := utils.GenerateToken(userID, "user", s.jwtSecret, s.accessTokenTTL)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		return nil, &rerrors.InternalServerError{Msg: "error generating tokens"}
	}

	refreshToken, err := utils.GenerateToken(userID, "refresh", s.jwtSecret, s.refreshTokenTTL)
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		return nil, &rerrors.InternalServerError{Msg: "error generating tokens"}
	}

	response := &dto.RegisterResponse{
		UserID:       userID,
		Email:        req.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.accessTokenTTL.Seconds()),
	}

	if personID != "" {
		response.PersonID = personID
	}

	return response, nil
}

// Login autentica un usuario
func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Validar datos
	if req.Email == "" {
		return nil, &rerrors.BadRequestError{Msg: "email is required"}
	}
	if req.Password == "" {
		return nil, &rerrors.BadRequestError{Msg: "password is required"}
	}

	// Buscar cuenta por email
	account, err := s.accountRepo.GetByEmail(req.Email)
	if err != nil {
		log.Printf("Error getting account by email: %v", err)
		return nil, &rerrors.InternalServerError{Msg: "error during authentication"}
	}
	if account == nil {
		return nil, &rerrors.UnauthorizedError{Msg: "invalid credentials"}
	}

	// Verificar estado de la cuenta
	if account.Status == models.StatusSuspended {
		return nil, &rerrors.ForbiddenError{Msg: "account suspended"}
	}
	if account.Status == models.StatusPending {
		return nil, &rerrors.ForbiddenError{Msg: "account pending activation"}
	}

	// Verificar contraseña
	if err := utils.CheckPassword(account.PasswordHash, req.Password); err != nil {
		return nil, &rerrors.UnauthorizedError{Msg: "invalid credentials"}
	}

	// Obtener usuario
	user, err := s.userRepo.GetByAccountID(account.ID)
	if err != nil {
		log.Printf("Error getting user by account ID: %v", err)
		return nil, &rerrors.InternalServerError{Msg: "error during authentication"}
	}
	if user == nil {
		return nil, &rerrors.InternalServerError{Msg: "user not found"}
	}

	// Actualizar último login
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		log.Printf("Error updating last login: %v", err)
		// No fallar por esto
	}

	// Generar tokens
	accessToken, err := utils.GenerateToken(user.ID, "user", s.jwtSecret, s.accessTokenTTL)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		return nil, &rerrors.InternalServerError{Msg: "error generating tokens"}
	}

	refreshToken, err := utils.GenerateToken(user.ID, "refresh", s.jwtSecret, s.refreshTokenTTL)
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		return nil, &rerrors.InternalServerError{Msg: "error generating tokens"}
	}

	response := &dto.LoginResponse{
		UserID:       user.ID,
		Email:        account.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.accessTokenTTL.Seconds()),
	}

	if user.PersonID != nil {
		response.PersonID = *user.PersonID
	}

	return response, nil
}

// RefreshToken renueva un token de acceso
func (s *AuthService) RefreshToken(refreshToken string) (*dto.RefreshTokenResponse, error) {
	// Validar y parsear refresh token
	claims, err := utils.ParseToken(refreshToken, s.jwtSecret)
	if err != nil {
		return nil, &rerrors.UnauthorizedError{Msg: "invalid refresh token"}
	}

	// Verificar que es un refresh token
	if claims.Role != "refresh" {
		return nil, &rerrors.UnauthorizedError{Msg: "invalid token type"}
	}

	// Verificar que el usuario existe y está activo
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
		return nil, &rerrors.InternalServerError{Msg: "error validating token"}
	}
	if user == nil {
		return nil, &rerrors.UnauthorizedError{Msg: "user not found"}
	}

	// Verificar estado de la cuenta
	if user.Account != nil {
		if user.Account.Status == models.StatusSuspended {
			return nil, &rerrors.ForbiddenError{Msg: "account suspended"}
		}
		if user.Account.Status != models.StatusActive {
			return nil, &rerrors.ForbiddenError{Msg: "account not active"}
		}
	}

	// Generar nuevos tokens
	newAccessToken, err := utils.GenerateToken(user.ID, "user", s.jwtSecret, s.accessTokenTTL)
	if err != nil {
		log.Printf("Error generating new access token: %v", err)
		return nil, &rerrors.InternalServerError{Msg: "error generating new token"}
	}

	newRefreshToken, err := utils.GenerateToken(user.ID, "refresh", s.jwtSecret, s.refreshTokenTTL)
	if err != nil {
		log.Printf("Error generating new refresh token: %v", err)
		return nil, &rerrors.InternalServerError{Msg: "error generating new token"}
	}

	return &dto.RefreshTokenResponse{
		UserID:       user.ID,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.accessTokenTTL.Seconds()),
	}, nil
}

// GetUserProfile obtiene el perfil del usuario
func (s *AuthService) GetUserProfile(userID string) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
		return nil, &rerrors.InternalServerError{Msg: "error getting user profile"}
	}
	if user == nil {
		return nil, &rerrors.NotFoundError{Msg: "user not found"}
	}

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

// ForgotPassword inicia el proceso de recuperación de contraseña
func (s *AuthService) ForgotPassword(email string) (*dto.ForgotPasswordResponse, error) {
	// Verificar que el email existe
	account, err := s.accountRepo.GetByEmail(email)
	if err != nil {
		log.Printf("Error getting account by email: %v", err)
		return nil, &rerrors.InternalServerError{Msg: "error processing request"}
	}

	// Siempre retornar éxito para no revelar si el email existe
	response := &dto.ForgotPasswordResponse{
		Message: "If the email exists, you will receive a password reset link",
		Success: true,
	}

	if account != nil && account.Status == models.StatusActive {
		// TODO: Aquí se integraría con el servicio de comunicación
		// para enviar el email de reset de contraseña
		log.Printf("Password reset requested for email: %s", email)

		// Por ahora solo logueamos, pero aquí iría:
		// 1. Generar token de reset temporal
		// 2. Almacenar en cache/BD con TTL
		// 3. Enviar email via servicio de comunicación
	}

	return response, nil
}

// ResetPassword resetea la contraseña con un token válido
func (s *AuthService) ResetPassword(req dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, error) {
	// TODO: Implementar cuando esté listo el servicio de comunicación
	// Por ahora retornar no implementado

	return &dto.ResetPasswordResponse{
		Message: "Password reset functionality will be available soon",
		Success: false,
	}, &rerrors.InternalServerError{Msg: "password reset not yet implemented"}
}

// createPersonIfNeeded crea una persona si se proporcionan datos
func (s *AuthService) createPersonIfNeeded(personData *dto.PersonData) (string, error) {
	if s.personClient == nil {
		return "", errors.New("person service client not available")
	}

	// Validar datos mínimos
	if personData.Type == "individual" && personData.Individual == nil {
		return "", errors.New("individual data required for individual type")
	}
	if personData.Type == "company" && personData.Company == nil {
		return "", errors.New("company data required for company type")
	}

	// Crear persona via gRPC
	personResp, err := s.personClient.CreatePerson(personData)
	if err != nil {
		return "", fmt.Errorf("failed to create person: %w", err)
	}

	return personResp.Id, nil
}

// ValidateToken valida un token JWT
func (s *AuthService) ValidateToken(token string) (*utils.Claims, error) {
	claims, err := utils.ParseToken(token, s.jwtSecret)
	if err != nil {
		return nil, &rerrors.UnauthorizedError{Msg: "invalid token"}
	}

	// Verificar que el usuario existe y está activo
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		log.Printf("Error getting user by ID during token validation: %v", err)
		return nil, &rerrors.UnauthorizedError{Msg: "invalid token"}
	}
	if user == nil {
		return nil, &rerrors.UnauthorizedError{Msg: "user not found"}
	}

	// Verificar estado de la cuenta
	if user.Account != nil && user.Account.Status != models.StatusActive {
		return nil, &rerrors.UnauthorizedError{Msg: "account not active"}
	}

	return claims, nil
}

// validatePersonData valida datos de persona usando rem-common validator
func (s *AuthService) validatePersonData(personData *dto.PersonData) error {
	if personData == nil {
		return nil
	}

	// Validar tipo de persona
	if err := s.validator.ValidatePersonType(personData.Type, "person_data.type"); err != nil {
		return fmt.Errorf("invalid person type: %s", err.Message)
	}

	// Validar según tipo
	switch personData.Type {
	case "individual":
		if personData.Individual == nil {
			return fmt.Errorf("individual data is required for individual type")
		}
		if errors := s.validator.ValidateIndividualData(&IndividualDataAdapter{data: personData.Individual}, "person_data.individual"); errors.HasErrors() {
			return fmt.Errorf("individual validation failed: %v", errors)
		}
	case "company":
		if personData.Company == nil {
			return fmt.Errorf("company data is required for company type")
		}
		if errors := s.validator.ValidateCompanyData(&CompanyDataAdapter{data: personData.Company}, "person_data.company"); errors.HasErrors() {
			return fmt.Errorf("company validation failed: %v", errors)
		}
	}

	// Validar contactos si existen
	if len(personData.Contacts) > 0 {
		var contactAdapters []validation.ContactData
		for _, contact := range personData.Contacts {
			contactAdapters = append(contactAdapters, &ContactDataAdapter{data: &contact})
		}
		if errors := s.validator.ValidateContactData(contactAdapters, "person_data.contacts"); errors.HasErrors() {
			return fmt.Errorf("contacts validation failed: %v", errors)
		}
	}

	// Validar dirección si existe
	if personData.Address != nil {
		if errors := s.validator.ValidateAddressData(&AddressDataAdapter{data: personData.Address}, "person_data.address"); errors.HasErrors() {
			return fmt.Errorf("address validation failed: %v", errors)
		}
	}

	return nil
}
