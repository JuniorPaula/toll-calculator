package main

import (
	"log"
	"tolling/aggregator/client"
)

const (
	kafkaTopic         = "obudata"
	aggregatorEndpoint = "http://127.0.0.1:4000/aggregate"
)

func main() {
	var (
		err error
		svc CalculatorServicer
	)

	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	// httpClient := client.NewHTTPClient(aggregatorEndpoint)
	grpcClient, err := client.NewGRPClient(aggregatorEndpoint)

	kc, err := NewKafkaConsumer(kafkaTopic, svc, grpcClient)
	if err != nil {
		log.Fatal(err)
	}
	kc.Start()
}
