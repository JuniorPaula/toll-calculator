package main

import (
  "fmt"
  "time"
  "math/rand"
  "math"
)

var sendInterval = time.Second

type OBUData struct {
  OBUID   int       `json:"obu_id"`
  Lat     float64   `json:"lat"`
  Long    float64   `json:"long"`
}

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
  for {
    for i := 0; i < len(obuIDs); i++ {
      lat, long := genLatLong()
      data := OBUData{
        OBUID: obuIDs[i],
        Lat: lat,
        Long: long,
      }
      fmt.Printf("%+v\n", data)
    }
    time.Sleep(sendInterval)
  }
}

func init() {
  rand.Seed(time.Now().UnixNano())
}
