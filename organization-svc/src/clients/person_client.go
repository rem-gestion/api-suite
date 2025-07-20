package clients

import (
	"context"
	"fmt"
	"time"

	rcgrpc "github.com/rem-gestion/rem-common/grpc"
	personpb "github.com/rem-gestion/rem-common/protos/person/v1"
	"google.golang.org/grpc"
)

// PersonClient maneja la comunicación gRPC con person-svc
type PersonClient struct {
	conn   *grpc.ClientConn
	client personpb.PersonServiceClient
	addr   string
}

// NewPersonClient crea una nueva instancia del cliente
func NewPersonClient(addr string) (*PersonClient, error) {
	conn, err := rcgrpc.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to person service at %s: %w", addr, err)
	}

	client := personpb.NewPersonServiceClient(conn)

	return &PersonClient{
		conn:   conn,
		client: client,
		addr:   addr,
	}, nil
}

// NewPersonClientWithConn crea una nueva instancia del cliente usando una conexión existente
func NewPersonClientWithConn(conn *grpc.ClientConn) *PersonClient {
	client := personpb.NewPersonServiceClient(conn)

	return &PersonClient{
		conn:   conn,
		client: client,
		addr:   "existing_connection",
	}
}

// Close cierra la conexión gRPC
func (c *PersonClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// CreatePerson crea una nueva persona
func (c *PersonClient) CreatePerson(req *personpb.CreatePersonRequest) (*personpb.PersonResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return c.client.CreatePerson(ctx, req)
}

// GetPerson obtiene una persona por ID
func (c *PersonClient) GetPerson(personID string) (*personpb.PersonResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &personpb.GetPersonRequest{
		Id: personID,
	}

	return c.client.GetPerson(ctx, req)
}

// ListPersons lista personas con filtros opcionales
func (c *PersonClient) ListPersons(req *personpb.ListPersonsRequest) (*personpb.ListPersonsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return c.client.ListPersons(ctx, req)
}

// UpdatePerson actualiza una persona existente
func (c *PersonClient) UpdatePerson(req *personpb.UpdatePersonRequest) (*personpb.PersonResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return c.client.UpdatePerson(ctx, req)
}

// DeletePerson elimina una persona
func (c *PersonClient) DeletePerson(personID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &personpb.DeletePersonRequest{
		Id: personID,
	}

	_, err := c.client.DeletePerson(ctx, req)
	return err
}

// AddContact agrega un contacto a una persona
func (c *PersonClient) AddContact(req *personpb.AddContactRequest) (*personpb.ContactResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return c.client.AddContact(ctx, req)
}

// UpdateContact actualiza un contacto existente
func (c *PersonClient) UpdateContact(req *personpb.UpdateContactRequest) (*personpb.ContactResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return c.client.UpdateContact(ctx, req)
}

// DeleteContact elimina un contacto
func (c *PersonClient) DeleteContact(contactID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &personpb.DeleteContactRequest{
		Id: contactID,
	}

	_, err := c.client.DeleteContact(ctx, req)
	return err
}

// ListContacts lista los contactos de una persona
func (c *PersonClient) ListContacts(personID string) (*personpb.ListContactsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &personpb.ListContactsRequest{
		PersonId: personID,
	}

	return c.client.ListContacts(ctx, req)
}
