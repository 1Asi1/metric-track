package models

type Metric struct {
	ID      string   `db:"id"`
	Gauge   *float64 `db:"gauge"`
	Counter *int64   `db:"counter"`
}
