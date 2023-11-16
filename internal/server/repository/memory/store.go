package memory

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
)

var (
	ErrNotFound = errors.New("name metric not found")
)

type Store struct {
	metric map[string]Type
	log    zerolog.Logger
}

type Type struct {
	Gauge   *float64
	Counter *int64
}

func New(log zerolog.Logger) Store {
	return Store{
		metric: make(map[string]Type),
		log:    log,
	}
}

func (m Store) Get(ctx context.Context) (map[string]Type, error) {
	return m.metric, nil
}

func (m Store) GetOne(ctx context.Context, name string) (Type, error) {
	l := m.log.With().Str("memory", "GetOne").Logger()

	if _, ok := m.metric[name]; !ok {
		l.Error().Err(ErrNotFound).Msgf("m.metric[name], name: %s", name)
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
