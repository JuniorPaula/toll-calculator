package types

type Invoice struct {
	OBUID         int     `json:"obu_id"`
	TotalDistance float64 `json:"total_distance"`
	TotalAmount   float64 `json:"total_amount"`
}

type Distance struct {
	Value float64 `json:"value"`
	OBUID int     `json:"ObuID"`
	Unix  int64   `json:"unix"`
}

type OBUData struct {
	OBUID int     `json:"obu_id"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}
