package service

import (
	"context"

	"github.com/1Asi1/metric-track.git/internal/server/repository/memory"
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

var TypeMetric = map[string]struct{}{
	Gauge:   {},
	Counter: {},
}

type MemStorage struct {
	Metrics map[string]Type
}

type Type struct {
	Gauge   float64
	Counter int64
}

type Request struct {
	Metric string
	Name   string
	Type   Type
}

type Service interface {
	UpdateMetric(context.Context, Request) error
}

type service struct {
	Store memory.Store
}

func New(store memory.Store) Service {
	return service{
		Store: store,
	}
}

func (s service) UpdateMetric(ctx context.Context, req Request) error {
	data, err := s.Store.Get(ctx)
	if err != nil {
		return err
	}

	value := data.Metrics[req.Name]
	if req.Metric == Gauge {
		value.Gauge = req.Type.Gauge
	} else {
		value.Counter += req.Type.Counter
	}
	data.Metrics[req.Name] = value

	err = s.Store.Update(ctx, data)
	if err != nil {
		return err
	}

	return nil
}
