package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"toll-calculator/types"

	"google.golang.org/grpc"
)

func main() {
	httpListenAddr := flag.String("listenaddr", ":5000", "the listen address of the server")
	grpcListendAddr := flag.String("grpc_listen_add", ":5001", "the listen address of the server")

	flag.Parse()
	var (
		store            = NewMemoryStore()
		srv   Aggregator = NewInvoiceAggregator(store)
	)

	srv = NewLogMiddleware(srv)
	go func() {
		_, err := makeGRPCTransport(*grpcListendAddr, srv)

		if err != nil {
			log.Fatal(err)
		}
	}()
	log.Fatal(makeHTTPTransport(*httpListenAddr, srv))
}

func makeHTTPTransport(listenaddr string, srv Aggregator) error {
	fmt.Println("HTTP transport running on port", listenaddr)
	http.HandleFunc("/agg", handleAggregate(srv))

	return http.ListenAndServe(listenaddr, nil)
}

func makeGRPCTransport(listenaddr string, srv Aggregator) (*types.None, error) {
	fmt.Println("GRPC transport running on port", listenaddr)

	ln, err := net.Listen("tcp", listenaddr)

	if err != nil {
		return nil, err
	}

	defer ln.Close()

	server := grpc.NewServer([]grpc.ServerOption{}...)

	types.RegisterDistanceAggregatorServer(server, NewGRPCServer(srv))
	return nil, server.Serve(ln)
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}
