package main

import (
	"fmt"
	"log"
	"net/http"

	"tolling/types"

	"github.com/gorilla/websocket"
)

type DataReceive struct {
	msgChan chan types.OBUData
	conn    *websocket.Conn
	prod    DataProducer
}

func NewDataReceive() (*DataReceive, error) {
	var (
		p          DataProducer
		err        error
		kafkaTopic = "obudata"
	)

	p, err = NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}

	p = NewLoggerMiddleware(p)
	return &DataReceive{
		msgChan: make(chan types.OBUData, 128),
		prod:    p,
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
	return dr.prod.ProducerData(data)
}

func (dr *DataReceive) handlerWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
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
