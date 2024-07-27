package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"tolling/types"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3000", "the listen address of HTTP server")
	flag.Parse()

	store := NewMemoryStore()
	var (
		svc = NewInvoicerAggregator(store)
	)
	makeHTTPTransport(*listenAddr, svc)
}

func makeHTTPTransport(addr string, svc Aggregator) {
	fmt.Printf("HTTP transport running on port %s", addr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.ListenAndServe(addr, nil)
}

func handleAggregate(Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
