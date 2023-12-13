package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/1Asi1/metric-track.git/internal/server/repository/memory"
	"github.com/rs/zerolog"
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

var TypeMetric = map[string]struct{}{
	Gauge:   {},
	Counter: {},
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type MetricsRequest struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type Store interface {
	Get(ctx context.Context) (map[string]memory.Type, error)
	GetOne(ctx context.Context, name string) (memory.Type, error)
	Update(ctx context.Context, data map[string]memory.Type)
	Ping() error
}

type Service struct {
	Store Store
	log   zerolog.Logger
}

func New(store Store, log zerolog.Logger) Service {
	return Service{
		Store: store,
		log:   log,
	}
}

func (s Service) GetMetric(ctx context.Context) (string, error) {
	l := s.log.With().Str("service", "GetMetric").Logger()

	data, err := s.Store.Get(ctx)
	if err != nil {
		l.Error().Err(err).Msg("s.Store.Get")
		return "", err
	}

	l.Debug().Msgf("data value: %+v", data)

	res := s.parseToHTML(data)

	return res, nil
}

func (s Service) GetOneMetric(ctx context.Context, req MetricsRequest) (Metrics, error) {
	l := s.log.With().Str("service", "GetOneMetric").Logger()

	data, err := s.Store.GetOne(ctx, req.ID)
	if err != nil {
		l.Error().Err(err).Msgf("s.Store.GetOne metric id: %s", req.ID)
		return Metrics{}, fmt.Errorf("metric name: %s, error:%w", req.ID, err)
	}

	l.Debug().Msgf("data value: %+v", data)

	if req.MType == Gauge {
		return Metrics{
			ID:    req.ID,
			MType: Gauge,
			Value: data.Gauge,
			Delta: nil,
		}, nil
	}

	return Metrics{
		ID:    req.ID,
		MType: Counter,
		Value: nil,
		Delta: data.Counter,
	}, nil
}

func (s Service) UpdateMetric(ctx context.Context, req MetricsRequest) (Metrics, error) {
	l := s.log.With().Str("service", "UpdateMetric").Logger()

	data, err := s.Store.Get(ctx)
	if err != nil {
		l.Error().Err(err).Msg("s.Store.Get")
		return Metrics{}, err
	}

	l.Debug().Msgf("data value: %+v", data)
	value := data[req.ID]
	if req.MType == Gauge {
		value.Gauge = req.Value
	} else {
		if value.Counter != nil {
			*value.Counter += *req.Delta
		} else {
			value.Counter = req.Delta
		}
	}
	data[req.ID] = value

	s.Store.Update(ctx, data)

	return Metrics{
		ID:    req.ID,
		MType: req.MType,
		Value: value.Gauge,
		Delta: value.Counter,
	}, nil
}

func (s Service) Ping(ctx context.Context) error {

	if err := s.Store.Ping(); err != nil {
		return err
	}

	return nil
}

func (s Service) parseToHTML(data map[string]memory.Type) string {
	var insert string

	for k, v := range data {
		var gauge float64
		var counter int64

		if v.Gauge != nil {
			gauge = *v.Gauge
		}

		if v.Counter != nil {
			counter = *v.Counter
		}

		insert += fmt.Sprintf(`
	<p><b>Имя: %s</p>
	<p><b>Gauge: %s</p>
	<p><b>Counter: %d</p>`,
			k,
			strconv.FormatFloat(gauge, 'f', -1, 64),
			counter) +
			"\n_______________\n"
	}

	res := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
	    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
	    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
	    <meta charset="UTF-8">
	    <title>Metrics</title>
	</head>
	<body>
	%s
	</body>
	</html>`,
		insert)

	return res
}

func (m Metrics) MarshalJSON() ([]byte, error) {
	type MetricsAlias Metrics
	aliasValue := struct {
		MetricsAlias
	}{
		MetricsAlias: MetricsAlias(m),
	}
	aliasValue.ID = m.ID
	aliasValue.MType = m.MType
	aliasValue.Delta = m.Delta
	if m.Value != nil {
		if float64(*m.Value) == float64(int(*m.Value)) {
			jsonData, err := json.Marshal(aliasValue)
			if err != nil {
				return nil, err
			}

			jsonString := string(jsonData)
			formattedJSON := strconv.FormatFloat(float64(*m.Value), 'f', 6, 64)
			jsonString = jsonString[:len(jsonString)-1] + formattedJSON[1:] + "}"

			return []byte(jsonString), nil
		}
	}

	return json.Marshal(aliasValue)
}
