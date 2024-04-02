package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/1Asi1/metric-track.git/internal/server/config"
	"github.com/1Asi1/metric-track.git/internal/server/repository/memory"
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

func TestV1_UpdateMetric(t *testing.T) {
	l := newLogger()
	st := memory.New(l, config.Config{})
	se := service.New(st, l)

	router := chi.NewRouter()
	h := rest.Handler{
		Mux:     router,
		Service: se}
	New(h, "", "")

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
	st := memory.New(l, config.Config{})
	se := service.New(st, l)

	router := chi.NewRouter()
	h := rest.Handler{
		Mux:     router,
		Service: se}
	New(h, "", "")

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

func TestV1_GetMetric(t *testing.T) {
	l := newLogger()
	st := memory.New(l, config.Config{})
	se := service.New(st, l)

	router := chi.NewRouter()
	h := rest.Handler{
		Mux:     router,
		Service: se}
	New(h, "", "")

	s := httptest.NewServer(router)
	defer s.Close()

	tests := []struct {
		name string
	}{
		{
			name: "positive",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/", s.URL)

			res, err := resty.New().R().SetHeader("Content-Type", "application/json; charset=utf-8").Get(url)

			require.NoError(t, err)

			assert.Equal(t, "text/html; charset=utf-8", res.Header().Get("Content-Type"))

			assert.Equal(t, http.StatusOK, res.StatusCode())

			err = res.RawBody().Close()
			require.NoError(t, err)
		})
	}
}

func TestV1_GetOneMetric(t *testing.T) {
	l := newLogger()
	st := memory.New(l, config.Config{})
	se := service.New(st, l)

	router := chi.NewRouter()
	h := rest.Handler{
		Mux:     router,
		Service: se}
	New(h, "", "")

	s := httptest.NewServer(router)
	defer s.Close()

	data := make(map[string]memory.Type)
	gauge := 1.0
	counter := int64(1)
	data["test"] = memory.Type{
		Gauge:   &gauge,
		Counter: &counter,
	}
	st.Update(context.Background(), "test", data)

	tests := []struct {
		name string
		want string
	}{
		{
			name: "positive",
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/value/gauge/test", s.URL)

			res, _ := resty.New().R().SetHeader("Content-Type", "application/json; charset=utf-8").Get(url)

			assert.Equal(t, "text/plain; charset=utf-8", res.Header().Get("Content-Type"))

			act := string(res.Body())

			assert.Equal(t, tt.want, act, "resty.New().R()")
		})
	}
}

func TestV1_GetOneMetric2(t *testing.T) {
	l := newLogger()
	st := memory.New(l, config.Config{})
	se := service.New(st, l)

	router := chi.NewRouter()
	h := rest.Handler{
		Mux:     router,
		Service: se}
	New(h, "", "")

	s := httptest.NewServer(router)
	defer s.Close()

	tests := []struct {
		name string
		arg  service.MetricsRequest
	}{
		{
			name: "positive",
			arg: service.MetricsRequest{
				ID:    "test",
				MType: "gauge",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/value/", s.URL)

			res, _ := resty.New().R().SetHeader("Content-Type", "application/json; charset=utf-8").SetBody(tt.arg).Post(url)

			assert.Equal(t, "text/plain; charset=utf-8", res.Header().Get("Content-Type"))

			assert.Equal(t, http.StatusNotFound, res.StatusCode())

			err := res.RawBody().Close()
			require.NoError(t, err)
		})
	}
}

func TestV1_Ping(t *testing.T) {
	l := newLogger()
	st := memory.New(l, config.Config{})
	se := service.New(st, l)

	router := chi.NewRouter()
	h := rest.Handler{
		Mux:     router,
		Service: se}
	New(h, "", "")

	s := httptest.NewServer(router)
	defer s.Close()

	tests := []struct {
		name string
	}{
		{
			name: "positive",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/ping", s.URL)

			res, _ := resty.New().R().SetHeader("Content-Type", "application/json; charset=utf-8").Get(url)

			assert.Equal(t, http.StatusOK, res.StatusCode())
		})
	}
}

func TestV1_Updates(t *testing.T) {
	l := newLogger()
	st := memory.New(l, config.Config{})
	se := service.New(st, l)

	router := chi.NewRouter()
	h := rest.Handler{
		Mux:     router,
		Service: se}
	New(h, "", "")

	s := httptest.NewServer(router)
	defer s.Close()

	tests := []struct {
		name string
		req  []service.MetricsRequest
	}{
		{name: "positive",
			req: []service.MetricsRequest{
				{
					ID:    "",
					MType: "counter",
					Delta: nil,
					Value: nil,
				},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/updates/", s.URL)

			res, _ := resty.New().R().SetHeader("Content-Type", "application/json; charset=utf-8").SetBody(tt.req).Post(url)

			assert.Equal(t, http.StatusOK, res.StatusCode())
		})
	}
}
