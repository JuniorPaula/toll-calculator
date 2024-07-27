package main

import "log"

// type DistanceCalculator struct {
// 	consumer DataConsumer
// }

var kafkaTopic = "obudata"

func main() {
	var (
		err error
		svc CalculatorServicer
	)

	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	kc, err := NewKafkaConsumer(kafkaTopic, svc)
	if err != nil {
		log.Fatal(err)
	}
	kc.Start()
}
