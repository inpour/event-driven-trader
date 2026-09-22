package backtest

type Type string

type Marker struct {
	Time     int64   `json:"time"` // Unix time
	Price    float64 `json:"price"`
	Position string  `json:"position"`
	Shape    string  `json:"shape"`
	Color    string  `json:"color"`
	Size     float32 `json:"size"`
	Text     string  `json:"text"`
}
