package client

import (
	"tolling/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPClient struct {
	Endpoint string
	types.AggregatorClient
}

func NewGRPClient(endpoint string) (*GRPClient, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)

	return &GRPClient{
		Endpoint:         endpoint,
		AggregatorClient: c,
	}, nil
}
