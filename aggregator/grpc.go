package main

import (
	"context"
	"errors"
	"strconv"
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

func (g *GRPCServer) GetInvoice(ctx context.Context, aggReq *types.GetInvoiceRequest) (*types.InvoiceData, error) {
	if aggReq == nil {
		return nil, errors.New("agg request payload missing")
	}

	obuID, err := strconv.Atoi(aggReq.ObuID)

	if err != nil {
		return nil, errors.New("invalid OBUID")
	}

	data, err := g.srv.CalculateInvoice(obuID)

	if err != nil {
		return nil, err
	}

	return &types.InvoiceData{
		OBUID:         int32(data.OBUID),
		TotalDistance: data.TotalDistance,
		TotalAmount:   data.TotalAmount,
	}, nil
}
