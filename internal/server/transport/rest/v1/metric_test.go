package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/1Asi1/metric-track.git/internal/server/transport/rest"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
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

type storeTest struct {
	metric map[string]service.Type
}

func newStore() storeTest {
	res := make(map[string]service.Type)
	gauge := 3.14
	res["Test"] = service.Type{Gauge: &gauge, Counter: nil}
	return storeTest{metric: res}
}

func (s storeTest) Get(ctx context.Context) (map[string]service.Type, error) {
	return s.metric, nil
}

func (s storeTest) GetOne(ctx context.Context, name string) (service.Type, error) {
	return s.metric["Test"], nil
}

func (s storeTest) Update(ctx context.Context, data map[string]service.Type) error {
	return nil
}

func TestV1_UpdateMetric(t *testing.T) {
	l := newLogger()
	st := newStore()
	se := service.New(st, l)

	router := chi.NewRouter()
	h := rest.Handler{
		Mux:     router,
		Service: se}
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
				metricType:  "gauge",
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

func TestV1_UpdateMetric2(t *testing.T) {
	l := newLogger()
	st := newStore()
	se := service.New(st, l)

	router := chi.NewRouter()
	h := rest.Handler{
		Mux:     router,
		Service: se}
	New(h)

	s := httptest.NewServer(router)
	defer s.Close()

	value := 1.1

	tests := []struct {
		name string
		args service.MetricsRequest
	}{
		{
			name: "positive",
			args: service.MetricsRequest{
				ID:    "test",
				MType: "gauge",
				Value: &value,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/update/", s.URL)

			res, err := resty.New().R().SetHeader("Content-Type", "application/json; charset=utf-8").SetBody(tt.args).Post(url)

			require.NoError(t, err)

			assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

			assert.Equal(t, http.StatusOK, res.StatusCode())

			err = res.RawBody().Close()
			require.NoError(t, err)
		})
	}
}
