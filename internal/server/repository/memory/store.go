package memory

import (
	"context"
)

type Store interface {
	Get(context.Context) (map[string]Type, error)
	Update(context.Context, map[string]Type) error
}

type memoryStore struct {
	metric map[string]Type
}

type Type struct {
	Gauge   float64
	Counter int64
}

func New(path string) Store {
	return memoryStore{
		make(map[string]Type),
	}
}

func (m memoryStore) Get(ctx context.Context) (map[string]Type, error) {
	return m.metric, nil
}

func (m memoryStore) Update(ctx context.Context, data map[string]Type) error {
	for k, v := range data {
		m.metric[k] = v
	}
	return nil
}
