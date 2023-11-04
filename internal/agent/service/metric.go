package service

import (
	"math/rand"
	"runtime"

	"github.com/1Asi1/metric-track.git/internal/config"
)

type Gauge float64
type Counter int64

type Metric struct {
	runtime.MemStats
	PollCount   Counter
	RandomValue Gauge
}

type Service interface {
	GetMetric() Metric
}

type service struct {
	cfg config.Config
}

func New(cfg config.Config) Service {
	return service{
		cfg: cfg,
	}
}

func (s service) GetMetric() Metric {
	var m Metric
	memStat := &(m).MemStats
	runtime.ReadMemStats(memStat)

	m.RandomValue = Gauge(rand.ExpFloat64())

	return m
}
