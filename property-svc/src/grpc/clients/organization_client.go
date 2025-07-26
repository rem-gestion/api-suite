package clients

import (
	"context"
	"fmt"
	"time"

	organizationpb "github.com/rem-gestion/rem-common/protos/organization/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrganizationClient struct {
	client organizationpb.OrganizationServiceClient
	conn   *grpc.ClientConn
}

func NewOrganizationClient(host string, port int) (*OrganizationClient, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to organization service at %s: %w", addr, err)
	}

	client := organizationpb.NewOrganizationServiceClient(conn)

	return &OrganizationClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *OrganizationClient) GetOrganization(ctx context.Context, organizationID string) (*organizationpb.Organization, error) {
	req := &organizationpb.GetOrganizationRequest{
		Id: organizationID,
	}

	resp, err := c.client.GetOrganization(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization %s: %w", organizationID, err)
	}

	return resp, nil
}

func (c *OrganizationClient) ValidateOrganizationExists(ctx context.Context, organizationID string) error {
	_, err := c.GetOrganization(ctx, organizationID)
	return err
}

func (c *OrganizationClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
