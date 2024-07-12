package main

import (
  "fmt"
  "log"
  "encoding/json"
  "net/http"

  "tolling/types"
  "github.com/gorilla/websocket"

  "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var kafkaTopic = "obudata"

type DataReceive struct {
  msgChan chan types.OBUData
  conn *websocket.Conn
  kProducer *kafka.Producer
}

func NewDataReceive() (*DataReceive, error) {
  p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

  return &DataReceive{
    msgChan: make(chan types.OBUData, 128),
    kProducer: p,
  }, nil
}

func main() {
  recv, err := NewDataReceive()
  if err != nil {
    log.Fatal(err)
  }

  http.HandleFunc("/ws", recv.handlerWS)

  http.ListenAndServe(":3000", nil)
}

func (dr *DataReceive) producerData(data types.OBUData) error {
  b, err := json.Marshal(data)
  if err != nil {
    return nil
  }

	err = dr.kProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
		Value:          b,
	}, nil)

  return err
}

func (dr *DataReceive) handlerWS(w http.ResponseWriter, r *http.Request) {
  u := websocket.Upgrader{
    ReadBufferSize: 1028,
    WriteBufferSize: 1028,
  }

  conn, err := u.Upgrade(w, r, nil)
  if err != nil {
    log.Fatal(err)
  }
  dr.conn = conn

  go dr.wsReceiveLoop()
}

func (dr *DataReceive) wsReceiveLoop() {
  fmt.Println("New OBU Connected!!!")
  for {
    var data types.OBUData
    if err := dr.conn.ReadJSON(&data); err != nil {
      log.Println("Erro READ JSON", err)
      continue
    }
    if err := dr.producerData(data); err != nil {
      fmt.Println("Kafka producer error:", err)
    }
    //fmt.Printf("receive OBU data from [%d] :: <lat %.2f, long %.2f>\n", data.OBUID, data.Lat, data.Long)
    // dr.msgChan <- data
  }
}
