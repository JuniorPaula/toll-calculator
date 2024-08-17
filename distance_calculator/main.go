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

	kc, err := NewKafkaConsumer(kafkaTopic, svc, client.NewClient(aggregatorEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	kc.Start()
}
