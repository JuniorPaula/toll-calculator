package main

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer  *kafka.Consumer
	isRunning bool
}

func NewKafkaConsumer(topic string) (*KafkaConsumer, error) {
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
		consumer: c,
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
		fmt.Println(string(msg.Value))
	}
}
