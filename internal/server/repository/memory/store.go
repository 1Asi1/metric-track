package memory

import (
	"bufio"
	"context"
	"encoding/json"
	"os"

	"github.com/1Asi1/metric-track.git/internal/server/models"
)

type Store interface {
	Get(context.Context) (models.MemStorage, error)
	Update(context.Context, models.MemStorage) error
}

type memoryStore struct {
	Path string
}

func New(path string) Store {
	return memoryStore{
		Path: path,
	}
}

func (m memoryStore) Get(ctx context.Context) (models.MemStorage, error) {
	file, err := os.Open(m.Path)
	if err != nil {
		return models.MemStorage{}, err
	}
	defer file.Close()

	var result models.MemStorage
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err = json.Unmarshal(scanner.Bytes(), &result); err != nil {
			return models.MemStorage{}, err
		}
	}

	if len(result.Metrics) == 0 {
		var data models.MemStorage
		data.Metrics = map[string]models.Type{}
		return data, nil
	}

	return result, nil
}

func (m memoryStore) Update(ctx context.Context, data models.MemStorage) error {
	file, err := os.OpenFile(m.Path, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	result, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = file.Write(result)
	if err != nil {
		return err
	}

	return nil
}
