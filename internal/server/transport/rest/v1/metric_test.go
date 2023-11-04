package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/1Asi1/metric-track.git/internal/server/transport/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type storeTest interface {
	UpdateMetric(context.Context, service.Request) error
}

type store struct {
}

func new() storeTest {
	return store{}
}

func (s store) UpdateMetric(context.Context, service.Request) error {
	return nil
}

func TestV1_UpdateMetric(t *testing.T) {
	st := new()
	h := V1{
		rest.Handler{},
		st,
	}

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
			target := fmt.Sprintf("/update/%s/%s/%s", tt.args.metricType, tt.args.metricName, tt.args.metricValue)
			r := httptest.NewRequest(http.MethodPost, target, nil)

			w := httptest.NewRecorder()

			h.UpdateMetric(w, r)

			res := w.Result()

			assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))

			assert.Equal(t, res.StatusCode, http.StatusOK)

			err := res.Body.Close()

			require.NoError(t, err)
		})
	}
}
