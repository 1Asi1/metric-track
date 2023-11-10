package memory

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("name metric not found error")
)

type Store struct {
	metric map[string]Type
}

type Type struct {
	Gauge   float64
	Counter int64
}

func New() Store {
	return Store{
		make(map[string]Type),
	}
}

func (m Store) Get(ctx context.Context) (map[string]Type, error) {
	return m.metric, nil
}

func (m Store) GetOne(ctx context.Context, name string) (Type, error) {
	if _, ok := m.metric[name]; !ok {
		return Type{}, ErrNotFound
	}

	return m.metric[name], nil
}

func (m Store) Update(ctx context.Context, data map[string]Type) error {
	for k, v := range data {
		m.metric[k] = v
	}
	return nil
}
