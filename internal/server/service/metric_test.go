package service

import (
	"context"
	"testing"

	"github.com/1Asi1/metric-track.git/internal/server/models"
	"github.com/stretchr/testify/assert"
)

type storeTest interface {
	Get(context.Context) (models.MemStorage, error)
	Update(context.Context, models.MemStorage) error
}

type store struct {
}

func new() storeTest {
	return store{}
}

func (s store) Get(context.Context) (models.MemStorage, error) {
	var data models.MemStorage
	data.Metrics = map[string]models.Type{}
	return data, nil
}

func (s store) Update(context.Context, models.MemStorage) error {
	return nil
}

func Test_service_UpdateMetric(t *testing.T) {
	st := new()
	srv := service{Store: st}

	type args struct {
		ctx context.Context
		req Request
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "positive",
			args: args{
				ctx: context.Background(),
				req: Request{},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := srv.UpdateMetric(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err, nil)
		})
	}
}
