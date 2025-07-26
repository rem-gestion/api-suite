package clients

import (
	"context"
	"fmt"

	"github.com/rem-gestion/rem-common/config"
	"go.uber.org/zap"
)

type ClientManager struct {
	AddressClient      *AddressClient
	PersonClient       *PersonClient
	OrganizationClient *OrganizationClient
	logger             *zap.Logger
}

type ClientConfig struct {
	AddressServiceHost      string
	AddressServicePort      int
	PersonServiceHost       string
	PersonServicePort       int
	OrganizationServiceHost string
	OrganizationServicePort int
}

func NewClientManager(cfg *ClientConfig, logger *zap.Logger) (*ClientManager, error) {
	lg := logger.Named("grpc-clients")

	// Address Client
	addressClient, err := NewAddressClient(cfg.AddressServiceHost, cfg.AddressServicePort)
	if err != nil {
		lg.Warn("failed to connect to address service", zap.Error(err))
		// No retornamos error aquí, permitimos que el servicio funcione sin validaciones externas
	}

	// Person Client
	personClient, err := NewPersonClient(cfg.PersonServiceHost, cfg.PersonServicePort)
	if err != nil {
		lg.Warn("failed to connect to person service", zap.Error(err))
		// No retornamos error aquí, permitimos que el servicio funcione sin validaciones externas
	}

	// Organization Client
	organizationClient, err := NewOrganizationClient(cfg.OrganizationServiceHost, cfg.OrganizationServicePort)
	if err != nil {
		lg.Warn("failed to connect to organization service", zap.Error(err))
		// No retornamos error aquí, permitimos que el servicio funcione sin validaciones externas
	}

	return &ClientManager{
		AddressClient:      addressClient,
		PersonClient:       personClient,
		OrganizationClient: organizationClient,
		logger:             lg,
	}, nil
}

// ValidateAddressExists valida que existe una dirección
func (cm *ClientManager) ValidateAddressExists(ctx context.Context, addressID string) error {
	if cm.AddressClient == nil {
		cm.logger.Warn("address client not available, skipping validation", zap.String("address_id", addressID))
		return nil // No fallar si el cliente no está disponible
	}

	return cm.AddressClient.ValidateAddressExists(ctx, addressID)
}

// ValidatePersonExists valida que existe una persona
func (cm *ClientManager) ValidatePersonExists(ctx context.Context, personID string) error {
	if cm.PersonClient == nil {
		cm.logger.Warn("person client not available, skipping validation", zap.String("person_id", personID))
		return nil // No fallar si el cliente no está disponible
	}

	return cm.PersonClient.ValidatePersonExists(ctx, personID)
}

// ValidateOrganizationExists valida que existe una organización
func (cm *ClientManager) ValidateOrganizationExists(ctx context.Context, organizationID string) error {
	if cm.OrganizationClient == nil {
		cm.logger.Warn("organization client not available, skipping validation", zap.String("organization_id", organizationID))
		return nil // No fallar si el cliente no está disponible
	}

	return cm.OrganizationClient.ValidateOrganizationExists(ctx, organizationID)
}

// ValidateManagerExists valida que existe un manager (puede ser persona u organización)
// según el manager_type_id (1=person, 2=organization, etc.)
func (cm *ClientManager) ValidateManagerExists(ctx context.Context, managerID string, managerTypeID int32) error {
	// Según el ERD, manager_type indica si es person u organization
	// Asumimos: 1 = PERSON, 2 = ORGANIZATION (deberías confirmar con los datos reales)
	switch managerTypeID {
	case 1: // PERSON
		return cm.ValidatePersonExists(ctx, managerID)
	case 2: // ORGANIZATION
		return cm.ValidateOrganizationExists(ctx, managerID)
	default:
		return fmt.Errorf("unknown manager_type_id: %d", managerTypeID)
	}
}

func (cm *ClientManager) Close() error {
	var errors []error

	if cm.AddressClient != nil {
		if err := cm.AddressClient.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close address client: %w", err))
		}
	}

	if cm.PersonClient != nil {
		if err := cm.PersonClient.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close person client: %w", err))
		}
	}

	if cm.OrganizationClient != nil {
		if err := cm.OrganizationClient.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close organization client: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors closing clients: %v", errors)
	}

	return nil
}

// LoadFromConfig carga la configuración de los servicios gRPC desde rem-common/config
func LoadClientConfigFromEnv(cfg *config.Config) *ClientConfig {
	return &ClientConfig{
		// Address Service
		AddressServiceHost: getEnvOrDefault("REM_ADDRESS_HOST", "localhost"),
		AddressServicePort: getEnvPortOrDefault("REM_ADDRESS_PORT", 50052),
		// Person Service
		PersonServiceHost: getEnvOrDefault("REM_PERSON_HOST", "localhost"),
		PersonServicePort: getEnvPortOrDefault("REM_PERSON_PORT", 50053),
		// Organization Service (usando los mismos defaults por ahora)
		OrganizationServiceHost: getEnvOrDefault("REM_ORGANIZATION_HOST", "localhost"),
		OrganizationServicePort: getEnvPortOrDefault("REM_ORGANIZATION_PORT", 50055),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	// TODO: Usar os.Getenv o el sistema de config de rem-common
	// Por ahora devolvemos el default hasta que implementemos la lectura de env
	return defaultValue
}

func getEnvPortOrDefault(key string, defaultValue int) int {
	// TODO: Usar os.Getenv + strconv.Atoi o el sistema de config de rem-common
	// Por ahora devolvemos el default hasta que implementemos la lectura de env
	return defaultValue
}
