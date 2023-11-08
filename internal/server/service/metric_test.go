package service

import (
	"context"
	"testing"

	"github.com/1Asi1/metric-track.git/internal/server/repository/memory"
	"github.com/stretchr/testify/assert"
)

//type storeTest interface {
//	Get(context.Context) (map[string]Type, error)
//	Update(context.Context, map[string]Type) error
//}

type store struct {
	metric map[string]memory.Type
}

func (s store) Get(ctx context.Context) (map[string]memory.Type, error) {
	return s.metric, nil
}

func (s store) Update(ctx context.Context, m map[string]memory.Type) error {
	return nil
}

func newStoreTest() memory.Store {
	return store{metric: make(map[string]memory.Type)}
}

func Test_service_UpdateMetric(t *testing.T) {
	st := newStoreTest()
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
