package service

import (
	"context"
	"os"
	"testing"

	"github.com/1Asi1/metric-track.git/internal/server/config"
	"github.com/1Asi1/metric-track.git/internal/server/repository/memory"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func Test_service_UpdateMetric(t *testing.T) {
	l := newLogger()
	st := memory.New(l, config.Config{})
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

func TestService_GetMetric(t *testing.T) {
	l := newLogger()
	st := memory.New(l, config.Config{})
	srv := Service{Store: st}

	tests := []struct {
		name string
		want string
	}{
		{
			name: "positive",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := srv.GetMetric(context.Background())
			require.NoErrorf(t, err, "srv.GetMetric")
			assert.NotNil(t, got)
		})
	}
}

func TestService_GetOneMetric(t *testing.T) {
	l := newLogger()
	st := memory.New(l, config.Config{})
	srv := Service{Store: st}

	tests := []struct {
		name string
		req  MetricsRequest
		want error
	}{
		{
			name: "negative",
			req:  MetricsRequest{},
			want: memory.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := srv.GetOneMetric(context.Background(), tt.req)
			assert.ErrorIs(t, err, tt.want, "srv.GetOneMetric")
		})
	}
}

func TestService_Updates(t *testing.T) {
	l := newLogger()
	st := memory.New(l, config.Config{})
	srv := Service{Store: st}

	tests := []struct {
		name string
		req  []MetricsRequest
	}{
		{
			name: "positive",
			req:  []MetricsRequest{{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := srv.Updates(context.Background(), tt.req)
			require.NoErrorf(t, err, "srv.Updates")
		})
	}
}
