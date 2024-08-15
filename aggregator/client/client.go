package client

import (
	"context"
	"toll-calculator/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
	GetInvoice(context.Context, *types.GetInvoiceRequest) (*types.Invoice, error)
}
