package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"toll-calculator/types"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (c *HTTPClient) Aggregate(ctx context.Context, distance *types.AggregateRequest) error {
	b, err := json.Marshal(distance)

	if err != nil {
		return nil
	}

	_ = ctx

	request, err := http.NewRequest("POST", c.Endpoint, bytes.NewReader(b))

	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the service response with non 200 status code %d", resp.StatusCode)
	}

	return nil
}

func (c *HTTPClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	return &types.Invoice{
		OBUID:         id,
		TotalDistance: 42.12,
		TotalAmount:   125.56,
	}, nil
}
