package clients

import (
	"fmt"

	"google.golang.org/grpc"
)

// ClientManager gestiona todos los clientes gRPC del servicio organization
type ClientManager struct {
	AuthIdentity *AuthIdentityClient
	Person       *PersonClient
	Address      *AddressClient
}

// NewClientManager crea una nueva instancia del gestor de clientes
func NewClientManager(authConn, personConn, addressConn *grpc.ClientConn) *ClientManager {
	return &ClientManager{
		AuthIdentity: NewAuthIdentityClientWithConn(authConn),
		Person:       NewPersonClientWithConn(personConn),
		Address:      NewAddressClientWithConn(addressConn),
	}
}

// NewClientManagerWithAddresses crea una nueva instancia del gestor de clientes usando direcciones
func NewClientManagerWithAddresses(authAddr, personAddr, addressAddr string) (*ClientManager, error) {
	// Crear cliente auth-identity
	authClient, err := NewAuthIdentityClient(authAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth-identity client: %w", err)
	}

	// Crear cliente person
	personClient, err := NewPersonClient(personAddr)
	if err != nil {
		authClient.Close() // Limpiar conexión anterior
		return nil, fmt.Errorf("failed to create person client: %w", err)
	}

	// Crear cliente address
	addressClient, err := NewAddressClient(addressAddr)
	if err != nil {
		authClient.Close() // Limpiar conexiones anteriores
		personClient.Close()
		return nil, fmt.Errorf("failed to create address client: %w", err)
	}

	return &ClientManager{
		AuthIdentity: authClient,
		Person:       personClient,
		Address:      addressClient,
	}, nil
}

// Close cierra todas las conexiones gRPC
func (cm *ClientManager) Close() error {
	var errs []error

	if cm.AuthIdentity != nil {
		if err := cm.AuthIdentity.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close auth-identity client: %w", err))
		}
	}

	if cm.Person != nil {
		if err := cm.Person.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close person client: %w", err))
		}
	}

	if cm.Address != nil {
		if err := cm.Address.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close address client: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing clients: %v", errs)
	}

	return nil
}
