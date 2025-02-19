package client

import (
	"context"
	"tolling/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPClient struct {
	Endpoint string
	client   types.AggregatorClient
}

func NewGRPClient(endpoint string) (*GRPClient, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)

	return &GRPClient{
		Endpoint: endpoint,
		client:   c,
	}, nil
}

func (c *GRPClient) Aggregate(ctx context.Context, req *types.AggregateRequest) error {
	_, err := c.client.Aggregate(ctx, req)
	return err
}
