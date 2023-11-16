package service

import (
	"context"
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

type Type struct {
	Gauge   float64
	Counter int64
}

type Request struct {
	Metric string
	Name   string
	Type   Type
}

type Service struct {
	Store memory.Store
	log   zerolog.Logger
}

func New(store memory.Store, log zerolog.Logger) Service {
	return Service{
		Store: store,
		log:   log,
	}
}

func (s Service) GetMetric(ctx context.Context) (string, error) {
	l := s.log.With().Str("apiserver", "GetMetric").Logger()

	data, err := s.Store.Get(ctx)
	if err != nil {
		l.Error().Err(err).Msg("s.Store.Get")
		return "", err
	}

	l.Info().Msgf("data value: %+v", data)

	res := s.parseToHTML(data)

	return res, nil
}

func (s Service) GetOneMetric(ctx context.Context, metric, name string) (string, error) {
	l := s.log.With().Str("apiserver", "GetOneMetric").Logger()

	data, err := s.Store.GetOne(ctx, name)
	if err != nil {
		l.Error().Err(err).Msgf("s.Store.GetOne metric name: %s", name)
		return "", fmt.Errorf("metric name: %s, error:%w", name, err)
	}

	l.Info().Msgf("data value: %+v", data)

	if metric == Gauge {
		frmt := strconv.FormatFloat(data.Gauge, 'f', -1, 64)
		return frmt, nil
	}

	return strconv.FormatInt(data.Counter, 10), nil
}

func (s Service) UpdateMetric(ctx context.Context, req Request) error {
	l := s.log.With().Str("apiserver", "UpdateMetric").Logger()

	data, err := s.Store.Get(ctx)
	if err != nil {
		l.Error().Err(err).Msg("s.Store.Get")
		return err
	}

	l.Info().Msgf("data value: %+v", data)

	value := data[req.Name]
	if req.Metric == Gauge {
		value.Gauge = req.Type.Gauge
	} else {
		value.Counter += req.Type.Counter
	}
	data[req.Name] = value

	err = s.Store.Update(ctx, data)
	if err != nil {
		l.Error().Err(err).Msgf("s.Store.Update, data: %+v", data)
		return err
	}

	return nil
}

func (s Service) parseToHTML(data map[string]memory.Type) string {
	var insert string

	for k, v := range data {
		insert += fmt.Sprintf(`
	<p><b>Имя: %s</p>
	<p><b>Guage: %f</p>
	<p><b>Counter: %d</p>`,
			k,
			v.Gauge,
			v.Counter) +
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
