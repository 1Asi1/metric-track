package service

import (
	"runtime"

	"github.com/1Asi1/metric-track.git/internal/agent/config"
	"github.com/rs/zerolog"
)

type Gauge float64
type Counter int64

type Metric struct {
	Type map[string]any
}

type Service struct {
	cfg config.Config
	log zerolog.Logger
}

func New(cfg config.Config, log zerolog.Logger) Service {
	return Service{
		cfg: cfg,
		log: log,
	}
}

func (s Service) GetMetric() Metric {
	l := s.log.With().Str("service", "GetMetric").Logger()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	var res = map[string]any{
		"Alloc":         m.Alloc,
		"BuckHashSys":   m.BuckHashSys,
		"Frees":         m.Frees,
		"GCCPUFraction": m.GCCPUFraction,
		"GCSys":         m.GCSys,
		"HeapAlloc":     m.HeapAlloc,
		"HeapIdle":      m.HeapIdle,
		"HeapInuse":     m.HeapInuse,
		"HeapObjects":   m.HeapObjects,
		"HeapReleased":  m.HeapReleased,
		"HeapSys":       m.HeapSys,
		"LastGC":        m.LastGC,
		"Lookups":       m.Lookups,
		"MCacheInuse":   m.MCacheInuse,
		"MCacheSys":     m.MCacheSys,
		"MSpanInuse":    m.MSpanInuse,
		"MSpanSys":      m.MSpanSys,
		"Mallocs":       m.Mallocs,
		"NextGC":        m.NextGC,
		"NumForcedGC":   m.NumForcedGC,
		"NumGC":         m.NumGC,
		"OtherSys":      m.OtherSys,
		"PauseTotalNs":  m.PauseTotalNs,
		"StackInuse":    m.StackInuse,
		"StackSys":      m.StackSys,
		"Sys":           m.Sys,
		"TotalAlloc":    m.TotalAlloc,
		"RandomValue":   Gauge(0),
		"PollCount":     Counter(0),
	}

	l.Debug().Msgf("data value: %+v", res)

	return Metric{Type: res}
}
