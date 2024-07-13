package main

import "log"

// type DistanceCalculator struct {
// 	consumer DataConsumer
// }

var kafkaTopic = "obudata"

func main() {
	kc, err := NewKafkaConsumer(kafkaTopic)
	if err != nil {
		log.Fatal(err)
	}
	kc.Start()
}
