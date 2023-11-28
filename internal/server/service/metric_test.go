package service

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func newLogger() zerolog.Logger {
	out := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02 15:04:05 -0700",
		NoColor:    true,
	}

	l := zerolog.New(out)

	return l.Level(zerolog.InfoLevel).With().Timestamp().Logger()
}

type storeTest struct {
	metric map[string]Type
}

func newStore() storeTest {
	res := make(map[string]Type)
	gauge := 3.14
	res["Test"] = Type{Gauge: &gauge, Counter: nil}
	return storeTest{metric: res}
}

func (s storeTest) Get(ctx context.Context) (map[string]Type, error) {
	return s.metric, nil
}

func (s storeTest) GetOne(ctx context.Context, name string) (Type, error) {
	return s.metric["Test"], nil
}

func (s storeTest) Update(ctx context.Context, data map[string]Type) error {
	return nil
}

func Test_service_UpdateMetric(t *testing.T) {
	st := newStore()
	srv := Service{Store: st}

	type args struct {
		ctx context.Context
		req MetricsRequest
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
				req: MetricsRequest{},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := srv.UpdateMetric(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err, nil)
		})
	}
}
