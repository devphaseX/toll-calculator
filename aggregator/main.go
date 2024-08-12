package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"toll-calculator/types"
)

func main() {
	listenAddr := flag.String("listenaddr", ":5000", "the listen address of the server")
	flag.Parse()
	var (
		store            = NewMemoryStore()
		srv   Aggregator = NewInvoiceAggregator(store)
	)

	srv = NewLogMiddleware(srv)
	makeHTTPTransport(*listenAddr, srv)
}

func makeHTTPTransport(listenaddr string, srv Aggregator) {
	fmt.Println("HTTP transport running on port", listenaddr)
	http.HandleFunc("/agg", handleAggregate(srv))

	http.ListenAndServe(listenaddr, nil)
}

func handleAggregate(srv Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance

		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		if err := srv.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func handleGetInvoice(srv Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		obuIDS, ok := r.URL.Query()["obu_id"]

		if !ok || len(obuIDS) == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing obu_id in request query"})
			return
		}

		obuID, err := strconv.Atoi(obuIDS[0])

		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid obu_id passed in request query"})
			return
		}

		invoice, err := srv.CalculateInvoice(obuID)

		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"data": invoice})
	}
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	w.Header().Add("content-type", "application/json")

	return json.NewEncoder(w).Encode(v)
}
