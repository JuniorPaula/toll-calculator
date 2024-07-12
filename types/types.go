package types

type OBUData struct {
  OBUID   int       `json:"obu_id"`
  Lat     float64   `json:"lat"`
  Long    float64   `json:"long"`
}
