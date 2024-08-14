package main

import (
	"context"
	"toll-calculator/types"
)

type GRPCServer struct {
	types.UnimplementedDistanceAggregatorServer
	srv Aggregator
}

func NewGRPCServer(srv Aggregator) *GRPCServer {
	return &GRPCServer{
		srv: srv,
	}
}

func (g *GRPCServer) AggregateDistance(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
	data := types.Distance{
		OBUID: int(req.ObuID),
		Value: req.Value,
		Unix:  int64(req.Unix),
	}

	_ = ctx

	return nil, g.srv.AggregateDistance(data)
}
