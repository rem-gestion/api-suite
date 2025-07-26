package clients

import (
	"context"
	"fmt"
	"time"

	personpb "github.com/rem-gestion/rem-common/protos/person/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PersonClient struct {
	client personpb.PersonServiceClient
	conn   *grpc.ClientConn
}

func NewPersonClient(host string, port int) (*PersonClient, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to person service at %s: %w", addr, err)
	}

	client := personpb.NewPersonServiceClient(conn)

	return &PersonClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *PersonClient) GetPerson(ctx context.Context, personID string) (*personpb.PersonResponse, error) {
	req := &personpb.GetPersonRequest{
		Id: personID,
	}

	resp, err := c.client.GetPerson(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get person %s: %w", personID, err)
	}

	return resp, nil
}

func (c *PersonClient) ValidatePersonExists(ctx context.Context, personID string) error {
	_, err := c.GetPerson(ctx, personID)
	return err
}

func (c *PersonClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
