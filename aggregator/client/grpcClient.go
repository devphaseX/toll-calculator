package client

import (
	"context"
	"toll-calculator/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint string
	client   types.DistanceAggregatorClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(endpoint, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}...)

	if err != nil {
		return nil, err
	}

	c := types.NewDistanceAggregatorClient(conn)

	return &GRPCClient{
		Endpoint: endpoint,
		client:   c,
	}, nil
}

func (g *GRPCClient) Aggregate(ctx context.Context, aggReq *types.AggregateRequest) error {
	_, err := g.client.AggregateDistance(ctx, aggReq)
	return err
}

func (g *GRPCClient) GetInvoice(ctx context.Context, aggReq *types.GetInvoiceRequest) (*types.Invoice, error) {
	invoice, err := g.client.GetInvoice(ctx, aggReq)

	if err != nil {
		return nil, err
	}

	return &types.Invoice{
		OBUID:         int(invoice.OBUID),
		TotalDistance: invoice.TotalDistance,
		TotalAmount:   invoice.TotalAmount,
	}, nil
}
