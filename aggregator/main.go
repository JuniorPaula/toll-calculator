package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"tolling/types"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file %v", err)
	}

	httpAddr := os.Getenv("AGG_HTTP_ENDPOINT")
	rpcAddr := os.Getenv("AGG_GRPC_ENDPOINT")

	var (
		store = makeStore()
		svc   = NewInvoicerAggregator(store)
	)

	svc = NewMetrictsMeddleware(svc)
	svc = NewLogMiddleware(svc)

	// start GRPC transport for an goroutine
	go func() { makeGRPCTransport(fmt.Sprintf(":%s", rpcAddr), svc) }()

	// start HTTP transport an main goroutine
	log.Fatal(makeHTTPTransport(fmt.Sprintf(":%s", httpAddr), svc))
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Printf("[RPC] transport running on port (%s)\n", listenAddr)
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
	fmt.Printf("[HTTP] transport running on port (%s)\n", addr)

	var (
		aggMetricHandler  = newHTTPMetrictHandler("aggregate")
		invMetricHandler  = newHTTPMetrictHandler("invoice")
		invoiceHandler    = makeHTTPHandlerFunc(invMetricHandler.instrument(handleGetInvoice(svc)))
		aggregatorHandler = makeHTTPHandlerFunc(aggMetricHandler.instrument(handleAggregate(svc)))
	)

	http.HandleFunc("/aggregate", aggregatorHandler)
	http.HandleFunc("/invoice", invoiceHandler)

	// prometheus metrics route
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(addr, nil)
}
