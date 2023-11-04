package service

import (
	"runtime"
	"testing"

	"github.com/1Asi1/metric-track.git/internal/config"
	"github.com/stretchr/testify/assert"
)

func Test_service_GetMetric(t *testing.T) {
	var cfg config.Config
	tests := []struct {
		name string
		want Metric
	}{
		{
			name: "positive",
			want: Metric{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				cfg: cfg,
			}

			m := &(tt.want).MemStats
			runtime.ReadMemStats(m)

			data := s.GetMetric()

			assert.NotEqual(t, 0, data.RandomValue)
		})
	}
}
