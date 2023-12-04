package memory

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/rs/zerolog"
)

var (
	ErrNotFound = errors.New("name metric not found")
)

type MemoryStore struct {
	metric        map[string]service.Type
	log           zerolog.Logger
	storeRestore  bool
	storeInterval time.Duration
	storePath     string
}

type FileStore struct {
	MemoryStore
}

type Metric struct {
	Name  string   `json:"name"`
	Value *float64 `json:"value"`
	Delta *int64   `json:"delta"`
}

func New(log zerolog.Logger, storeRestore bool, storeInterval time.Duration, storePath string) MemoryStore {
	s := MemoryStore{
		metric:        make(map[string]service.Type),
		log:           log,
		storeRestore:  storeRestore,
		storeInterval: storeInterval,
		storePath:     storePath,
	}
	l := s.log.With().Str("memory", "New").Logger()

	if s.storeRestore {
		metric, err := getData(s.storePath, log)
		if err != nil {
			l.Err(err).Msg("s.getData")
		}

		if metric != nil {
			s.metric = metric
		}
	}

	fileStore := FileStore{s}
	go fileStore.dataRetentionPeriodic()

	return s
}

func (m MemoryStore) Get(ctx context.Context) (map[string]service.Type, error) {
	return m.metric, nil
}

func (m MemoryStore) GetOne(ctx context.Context, name string) (service.Type, error) {
	if _, ok := m.metric[name]; !ok {
		return service.Type{}, fmt.Errorf("problem with m.metric[%s]: %w", name, ErrNotFound)
	}

	return m.metric[name], nil
}

func (m MemoryStore) Update(ctx context.Context, data map[string]service.Type) {
	l := m.log.With().Str("memory", "DataRetentionPeriodic").Logger()

	for k, v := range data {
		m.metric[k] = v
	}

	fileStore := FileStore{m}
	err := fileStore.dataRetention()
	if err != nil {
		l.Err(err).Msg("m.DataRetention")
	}
}

func (f FileStore) dataRetentionPeriodic() {
	l := f.log.With().Str("memory", "DataRetentionPeriodic").Logger()

	if f.storeInterval != 0 {
		ticker := time.NewTicker(f.storeInterval)
		for range ticker.C {
			data, err := f.toData()
			if err != nil {
				l.Err(err).Msg("m.getData")
			}

			file, err := os.OpenFile(f.storePath, syscall.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				l.Err(err).Msg("os.OpenFile")

			}

			file.Write(data)
			file.Close()
		}
	}
}

func (f FileStore) dataRetention() error {
	if f.storeInterval == 0 {
		data, err := f.toData()
		if err != nil {
			return err
		}

		file, err := os.OpenFile(f.storePath, syscall.O_TRUNC|os.O_SYNC|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer file.Close()

		file.Write(data)
	}

	return nil
}

func (f FileStore) toData() ([]byte, error) {
	var metrics []Metric
	if len(f.metric) != 0 {
		for n, v := range f.metric {
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
		return nil, err
	}

	return data, nil
}

func getData(pathFile string, log zerolog.Logger) (map[string]service.Type, error) {
	file, err := os.OpenFile(pathFile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
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
		return nil, err
	}

	return metric, nil
}
