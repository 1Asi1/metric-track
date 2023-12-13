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

	"github.com/1Asi1/metric-track.git/internal/server/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	ErrNotFound = errors.New("name metric not found")
)

type Store interface {
	Get(ctx context.Context) (map[string]Type, error)
	GetOne(ctx context.Context, name string) (Type, error)
	Update(ctx context.Context, data map[string]Type)
	Ping() error
}

type StoreMemory struct {
	metric map[string]Type
	log    zerolog.Logger
}

type FileStore struct {
	memoryStore   StoreMemory
	storeRestore  bool
	storeInterval time.Duration
	storePath     string
}

type Metric struct {
	Name  string   `json:"name"`
	Value *float64 `json:"value"`
	Delta *int64   `json:"delta"`
}

type Type struct {
	Gauge   *float64
	Counter *int64
}

func New(log zerolog.Logger, cfg config.Config) Store {
	l := log.With().Str("memory", "New").Logger()

	store := StoreMemory{
		metric: make(map[string]Type),
		log:    log,
	}

	if len(cfg.StorePath) != 0 {
		// блок чтения данных из файла.
		if cfg.StoreRestore {
			metric, err := getData(cfg.StorePath, log)
			if err != nil {
				l.Err(err).Msg("s.getData")
			}

			if metric != nil {
				store.metric = metric
			}
		}

		// блок инициализации хранилища с записью в файл.
		fileStore := FileStore{
			memoryStore:   store,
			storeRestore:  cfg.StoreRestore,
			storeInterval: cfg.StoreInterval,
			storePath:     cfg.StorePath,
		}

		go fileStore.dataRetentionPeriodic()
		return fileStore
	}

	return store
}

func (m StoreMemory) Get(ctx context.Context) (map[string]Type, error) {
	return m.metric, nil
}

func (m StoreMemory) GetOne(ctx context.Context, name string) (Type, error) {
	if _, ok := m.metric[name]; !ok {
		return Type{}, fmt.Errorf("problem with m.metric[%s]: %w", name, ErrNotFound)
	}

	return m.metric[name], nil
}

func (m StoreMemory) Update(ctx context.Context, data map[string]Type) {
	for k, v := range data {
		m.metric[k] = v
	}
}

func (m StoreMemory) Ping() error {
	return errors.New("db not included")
}

func (f FileStore) Get(ctx context.Context) (map[string]Type, error) {
	return f.memoryStore.metric, nil
}

func (f FileStore) GetOne(ctx context.Context, name string) (Type, error) {
	if _, ok := f.memoryStore.metric[name]; !ok {
		return Type{}, fmt.Errorf("problem with m.metric[%s]: %w", name, ErrNotFound)
	}

	return f.memoryStore.metric[name], nil
}

func (f FileStore) Update(ctx context.Context, data map[string]Type) {
	l := f.memoryStore.log.With().Str("memory", "Update").Logger()
	for k, v := range data {
		f.memoryStore.metric[k] = v
	}

	err := f.dataRetention()
	if err != nil {
		l.Err(err).Msg("m.DataRetention")
	}
}

func (f FileStore) Ping() error {
	return errors.New("db not included")
}

func (f FileStore) dataRetentionPeriodic() {
	l := f.memoryStore.log.With().Str("memory", "dataRetentionPeriodic").Logger()

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

			_, err = file.Write(data)
			l.Err(err).Msg("file.Write")

			err = file.Close()
			l.Err(err).Msg("file.Close")
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
		defer func() { err = file.Close() }()

		_, err = file.Write(data)
		log.Err(err).Msg("file.Write")
	}

	return nil
}

func (f FileStore) toData() ([]byte, error) {
	var metrics []Metric
	if len(f.memoryStore.metric) != 0 {
		for n, v := range f.memoryStore.metric {
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

func getData(pathFile string, log zerolog.Logger) (map[string]Type, error) {
	file, err := os.OpenFile(pathFile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	var data []byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Bytes()...)
	}

	metric := make(map[string]Type)
	if len(data) != 0 {
		var metrics []Metric
		err = json.Unmarshal(data, &metrics)
		if err != nil {
			return nil, err
		}

		for _, v := range metrics {
			metric[v.Name] = Type{
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
