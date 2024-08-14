package client

import (
	"toll-calculator/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint string
	types.DistanceAggregatorClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endpoint, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}...)

	if err != nil {
		return nil, err
	}

	c := types.NewDistanceAggregatorClient(conn)

	return &GRPCClient{
		Endpoint:                 endpoint,
		DistanceAggregatorClient: c,
	}, nil
}
