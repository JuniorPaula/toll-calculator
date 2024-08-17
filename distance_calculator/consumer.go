package main

import (
	"encoding/json"
	"fmt"
	"time"
	"tolling/aggregator/client"
	"tolling/types"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer   *kafka.Consumer
	isRunning  bool
	service    CalculatorServicer
	aggrClient *client.Client
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, aggrClient *client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)
	return &KafkaConsumer{
		consumer:   c,
		service:    svc,
		aggrClient: aggrClient,
	}, nil
}

func (kc *KafkaConsumer) Start() {
	logrus.Info("Kafka consumer started")
	kc.isRunning = true
	kc.readMessagesLoop()
}

func (kc *KafkaConsumer) readMessagesLoop() {
	for kc.isRunning {
		msg, err := kc.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("Kafka consumer error: %v", err)
			continue
		}

		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("Error unmarshalling OBUData: %v", err)
			continue
		}

		distance, err := kc.service.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("Error calculating distance: %v", err)
			continue
		}
		logrus.WithFields(logrus.Fields{
			"obu_id": data.OBUID,
			"lat":    data.Lat,
			"long":   data.Long,
			"dist":   fmt.Sprintf("%.2f", distance),
		}).Info("Distance calculated")

		req := types.Distance{
			Value: distance,
			Unix:  time.Now().Unix(),
			OBUID: data.OBUID,
		}

		if err := kc.aggrClient.AggregateInvoice(req); err != nil {
			logrus.Errorf("aggregate error: [%v]", err)
			continue
		}
	}
}
