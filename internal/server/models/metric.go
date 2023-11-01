package models

type MemStorage struct {
	Metrics map[string]Type `json:"metrics"`
}

type Type struct {
	Gauge   float64 `json:"gauge"`
	Counter int64   `json:"counter"`
}
