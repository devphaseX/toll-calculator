package main

import (
	"fmt"
	"log"
	"os"
	"toll-calculator/types"
)

const basePrice = 3.15

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) *InvoiceAggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	fmt.Println("processing and inserting distance in the storage", distance)
	return i.store.Insert(distance)
}

func (i *InvoiceAggregator) CalculateInvoice(obuID int) (*types.Invoice, error) {
	dist, err := i.store.Get(obuID)

	if err != nil {
		return nil, fmt.Errorf("obu id (%d) not found", obuID)
	}

	invoice := types.Invoice{
		OBUID:         obuID,
		TotalDistance: dist,
		TotalAmount:   basePrice * dist,
	}

	return &invoice, nil
}

func makeStore() Storer {
	switch storeType := os.Getenv("AGG_STORE_TYPE"); storeType {
	case "memory":
		return NewMemoryStore()
	default:
		log.Fatalf("invalid store type given %s", storeType)
		return nil // this return is unreachable, but it is good practice to include it
	}
}
