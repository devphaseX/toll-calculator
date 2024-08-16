package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

	request, err := http.NewRequest("POST", c.Endpoint+"/agg", bytes.NewReader(b))

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

func (c *HTTPClient) GetInvoice(ctx context.Context, aggReq *types.GetInvoiceRequest) (*types.Invoice, error) {
	dataUrl := url.URL{
		Scheme: "http", // or "https" if your endpoint uses HTTPS
		Host:   c.Endpoint,
		Path:   "/invoice",
	}

	q := dataUrl.Query()
	q.Add("obu_id", aggReq.ObuID)

	// Set the query values back to the URL
	dataUrl.RawQuery = q.Encode()

	fullURL := dataUrl.String()
	fmt.Println(fullURL)

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body struct {
		Error *string        `json:"error"`
		Data  *types.Invoice `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errorMsg := body.Error

		if errorMsg == nil {
			return nil, errors.New("request failed with status code " + string(resp.StatusCode))
		}

		return nil, errors.New(*errorMsg)
	}

	invoice := body.Data

	if invoice == nil {
		return nil, errors.New("no data sent")
	}

	return invoice, nil
}
