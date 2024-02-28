package v1

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1Asi1/metric-track.git/internal/server/config"
	"github.com/1Asi1/metric-track.git/internal/server/repository/memory"
	"github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/1Asi1/metric-track.git/internal/server/transport/rest"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkUpdateMetric2(b *testing.B) {
	l := newLogger()
	st := memory.New(l, config.Config{})
	se := service.New(st, l)

	router := chi.NewRouter()
	h := rest.Handler{
		Mux:     router,
		Service: se}
	New(h, "")

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
		b.Run(tt.name, func(b *testing.B) {
			url := fmt.Sprintf("%s/update/", s.URL)

			res, err := resty.New().R().SetHeader("Content-Type", "application/json; charset=utf-8").SetBody(tt.args).Post(url)

			require.NoError(b, err)

			assert.Equal(b, "application/json", res.Header().Get("Content-Type"))

			assert.Equal(b, http.StatusOK, res.StatusCode())

			err = res.RawBody().Close()
			require.NoError(b, err)
		})
	}
}
