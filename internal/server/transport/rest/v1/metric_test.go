package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/1Asi1/metric-track.git/internal/server/transport/rest"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type storeTest interface {
	GetMetric(ctx context.Context) (string, error)
	GetOneMetric(ctx context.Context, metric, name string) (string, error)
	UpdateMetric(context.Context, service.Request) error
}

type store struct {
}

func new() storeTest {
	return store{}
}

func (s store) GetMetric(ctx context.Context) (string, error) {
	return "", nil
}

func (s store) GetOneMetric(ctx context.Context, metric, name string) (string, error) {
	return "", nil
}

func (s store) UpdateMetric(context.Context, service.Request) error {
	return nil
}

func TestV1_UpdateMetric(t *testing.T) {
	st := new()

	router := chi.NewRouter()
	h := rest.Handler{
		Mux:     router,
		Service: st}
	New(h)

	s := httptest.NewServer(router)
	defer s.Close()

	type args struct {
		metricType  string
		metricName  string
		metricValue string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive",
			args: args{
				metricType:  service.Gauge,
				metricName:  "Test",
				metricValue: "3.14",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/update/%s/%s/%s", s.URL, tt.args.metricType, tt.args.metricName, tt.args.metricValue)

			req := resty.New().R()
			req.Method = http.MethodPost
			req.URL = url

			res, err := req.Send()
			require.NoError(t, err)

			assert.Equal(t, "text/plain; charset=utf-8", res.Header().Get("Content-Type"))

			assert.Equal(t, http.StatusOK, res.StatusCode())

			err = res.RawBody().Close()
			require.NoError(t, err)
		})
	}
}
