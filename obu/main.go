package main

import (
	"log"
	"math"
	"math/rand"
	"time"
	"tolling/types"

	"github.com/gorilla/websocket"
)

const wsEndpoint = "ws://127.0.0.1:3000/ws"

var sendInterval = time.Second * 5

func genLatLong() (float64, float64) {
	return genCoord(), genCoord()
}

func genCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func generateOBUIDs(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}

func main() {
	obuIDs := generateOBUIDs(20)
	wsConn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatalf("Could not connect with Socket: %v\n", err)
	}

	for {
		for i := 0; i < len(obuIDs); i++ {
			lat, long := genLatLong()
			data := types.OBUData{
				OBUID: obuIDs[i],
				Lat:   lat,
				Long:  long,
			}
			if err := wsConn.WriteJSON(data); err != nil {
				log.Fatalf("Could not write json: %v\n", err)
			}
		}
		time.Sleep(sendInterval)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
