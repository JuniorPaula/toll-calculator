package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"tolling/types"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	httpAddr := flag.String("httpAddr", ":4000", "the listen address of HTTP server")
	rpcAddr := flag.String("rpcAddr", ":4001", "the listen address of gRPC server")
	flag.Parse()

	var (
		store = NewMemoryStore()
		svc   = NewInvoicerAggregator(store)
	)

	svc = NewMetrictsMeddleware(svc)
	svc = NewLogMiddleware(svc)

	// start GRPC transport for an goroutine
	go func() { makeGRPCTransport(*rpcAddr, svc) }()

	// start HTTP transport an main goroutine
	log.Fatal(makeHTTPTransport(*httpAddr, svc))
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Printf("[RPC] transport running on port (:%s)\n", listenAddr)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	server := grpc.NewServer([]grpc.ServerOption{}...)

	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))
	return server.Serve(ln)
}

func makeHTTPTransport(addr string, svc Aggregator) error {
	fmt.Printf("[HTTP] transport running on port (:%s)\n", addr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))

	// prometheus metrics route
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(addr, nil)
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, ok := r.URL.Query()["obu"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing OBU ID"})
			return
		}

		OBUID, err := strconv.Atoi(q[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid OBU ID"})
			return
		}

		invoice, err := svc.CalculateInvoice(OBUID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, invoice)
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
