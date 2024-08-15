package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"toll-calculator/aggregator/client"
	"toll-calculator/types"

	"github.com/sirupsen/logrus"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	listenaddr := flag.String("listenaddr", ":6000", "the listen address of the server")
	flag.Parse()

	var (
		client         = client.NewHTTPClient("localhost:5000")
		invoiceHandler = NewInvoiceHandler(client)
	)
	http.HandleFunc("/invoice", makeApiFunc(invoiceHandler.handleGetInvoice))
	logrus.Infof("Gateway http server listening on port %s\n", *listenaddr)
	if err := http.ListenAndServe(*listenaddr, nil); err != nil {
		log.Fatal(err)
	}
}

type InvoiceHandler struct {
	client client.Client
}

func NewInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: c,
	}
}

func (c *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	obuQuery, ok := r.URL.Query()["obu_id"]

	if !ok || len(obuQuery) == 0 {
		return fmt.Errorf("missing obu_id in request query")

	}

	obuID := obuQuery[0]

	inv, err := c.client.GetInvoice(context.Background(), &types.GetInvoiceRequest{
		ObuID: obuID,
	})

	if err != nil {
		return writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	return writeJSON(w, http.StatusOK, map[string]any{"data": inv})
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}

func makeApiFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}
