package memory

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
	"syscall"
	"time"

	"github.com/1Asi1/metric-track.git/internal/server/config"
	"github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/rs/zerolog"
)

var (
	ErrNotFound = errors.New("name metric not found")
)

type Store struct {
	metric map[string]service.Type
	log    zerolog.Logger
	cfg    config.Config
}

type Metric struct {
	Name  string   `json:"name"`
	Value *float64 `json:"value"`
	Delta *int64   `json:"delta"`
}

func New(log zerolog.Logger, cfg config.Config) Store {
	s := Store{
		metric: make(map[string]service.Type),
		log:    log,
		cfg:    cfg,
	}
	l := s.log.With().Str("memory", "New").Logger()

	if s.cfg.StoreRestore {
		metric, err := getData(cfg, log)
		if err != nil {
			l.Err(err).Msg("s.getData")
		}

		if metric != nil {
			s.metric = metric
		}
	}

	go s.dataRetentionPeriodic()

	return s
}

func (m Store) Get(ctx context.Context) (map[string]service.Type, error) {
	return m.metric, nil
}

func (m Store) GetOne(ctx context.Context, name string) (service.Type, error) {
	l := m.log.With().Str("memory", "GetOne").Logger()

	if _, ok := m.metric[name]; !ok {
		l.Error().Err(ErrNotFound).Msgf("m.metric[name], name: %s", name)
		return service.Type{}, ErrNotFound
	}

	return m.metric[name], nil
}

func (m Store) Update(ctx context.Context, data map[string]service.Type) error {
	l := m.log.With().Str("memory", "DataRetentionPeriodic").Logger()

	for k, v := range data {
		m.metric[k] = v
	}

	err := m.dataRetention()
	if err != nil {
		l.Err(err).Msg("m.DataRetention")
	}

	return nil
}

func (m Store) dataRetentionPeriodic() {
	l := m.log.With().Str("memory", "DataRetentionPeriodic").Logger()

	if m.cfg.StoreInterval != 0 {
		ticker := time.NewTicker(m.cfg.StoreInterval)
		for {
			select {
			case <-ticker.C:
				data, err := m.toData()
				if err != nil {
					l.Err(err).Msg("m.getData")
				}

				file, err := os.OpenFile(m.cfg.StorePath, syscall.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					l.Err(err).Msg("os.OpenFile")

				}

				file.Write(data)
				file.Close()

			default:
				continue
			}
		}
	}
}

func (m Store) dataRetention() error {
	l := m.log.With().Str("memory", "DataRetention").Logger()

	if m.cfg.StoreInterval == 0 {
		data, err := m.toData()
		if err != nil {
			l.Err(err).Msg("m.getData")
			return err
		}

		file, err := os.OpenFile(m.cfg.StorePath, syscall.O_TRUNC|os.O_SYNC|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			l.Err(err).Msg("os.OpenFile")
			return err
		}
		defer file.Close()

		file.Write(data)
	}

	return nil
}

func (m Store) toData() ([]byte, error) {
	l := m.log.With().Str("memory", "getData").Logger()

	var metrics []Metric
	if len(m.metric) != 0 {
		for n, v := range m.metric {
			metric := Metric{
				Name:  n,
				Value: v.Gauge,
				Delta: v.Counter,
			}

			metrics = append(metrics, metric)
		}
	} else {
		return nil, nil
	}

	data, err := json.Marshal(metrics)
	if err != nil {
		l.Err(err).Msg("json.Marshal")
		return nil, err
	}

	return data, nil
}

func getData(cfg config.Config, log zerolog.Logger) (map[string]service.Type, error) {
	l := log.With().Str("memory", "getData").Logger()

	file, err := os.OpenFile(cfg.StorePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		l.Err(err).Msg("os.OpenFile")
		return nil, err
	}

	var data []byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Bytes()...)
	}

	metric := make(map[string]service.Type)
	if len(data) != 0 {
		var metrics []Metric
		err = json.Unmarshal(data, &metrics)
		if err != nil {
			l.Err(err).Msg("json.Unmarshal")
			return nil, err
		}

		for _, v := range metrics {
			metric[v.Name] = service.Type{
				Gauge:   v.Value,
				Counter: v.Delta,
			}
		}
	}

	err = file.Close()
	if err != nil {
		l.Err(err).Msg("file.Close")
	}
	return metric, nil
}
