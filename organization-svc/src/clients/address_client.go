package clients

import (
	"context"
	"fmt"
	"time"

	rcgrpc "github.com/rem-gestion/rem-common/grpc"
	addresspb "github.com/rem-gestion/rem-common/protos/address/v1"
	"google.golang.org/grpc"
)

// AddressClient maneja la comunicación gRPC con address-svc
type AddressClient struct {
	conn   *grpc.ClientConn
	client addresspb.AddressServiceClient
	addr   string
}

// NewAddressClient crea una nueva instancia del cliente
func NewAddressClient(addr string) (*AddressClient, error) {
	conn, err := rcgrpc.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to address service at %s: %w", addr, err)
	}

	client := addresspb.NewAddressServiceClient(conn)

	return &AddressClient{
		conn:   conn,
		client: client,
		addr:   addr,
	}, nil
}

// NewAddressClientWithConn crea una nueva instancia del cliente usando una conexión existente
func NewAddressClientWithConn(conn *grpc.ClientConn) *AddressClient {
	client := addresspb.NewAddressServiceClient(conn)

	return &AddressClient{
		conn:   conn,
		client: client,
		addr:   "existing_connection",
	}
}

// Close cierra la conexión gRPC
func (c *AddressClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// CreateAddress crea una nueva dirección
func (c *AddressClient) CreateAddress(req *addresspb.CreateAddressRequest) (*addresspb.CreateAddressResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return c.client.Create(ctx, req)
}

// GetAddress obtiene una dirección por ID
func (c *AddressClient) GetAddress(addressID string) (*addresspb.GetAddressResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &addresspb.GetAddressRequest{
		Id: addressID,
	}

	return c.client.Get(ctx, req)
}
