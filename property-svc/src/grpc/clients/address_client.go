package clients

import (
	"context"
	"fmt"
	"time"

	addresspb "github.com/rem-gestion/rem-common/protos/address/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AddressClient struct {
	client addresspb.AddressServiceClient
	conn   *grpc.ClientConn
}

func NewAddressClient(host string, port int) (*AddressClient, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to address service at %s: %w", addr, err)
	}

	client := addresspb.NewAddressServiceClient(conn)

	return &AddressClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *AddressClient) GetAddress(ctx context.Context, addressID string) (*addresspb.Address, error) {
	req := &addresspb.GetAddressRequest{
		Id: addressID,
	}

	resp, err := c.client.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get address %s: %w", addressID, err)
	}

	return resp.Address, nil
}

func (c *AddressClient) ValidateAddressExists(ctx context.Context, addressID string) error {
	_, err := c.GetAddress(ctx, addressID)
	return err
}

func (c *AddressClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
