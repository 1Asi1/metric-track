package service

import (
	"os"
	"testing"

	"github.com/1Asi1/metric-track.git/internal/config"
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
			l := newLogger()
			s := Service{
				cfg: cfg,
				log: l,
			}

			data := s.GetMetric()

			assert.NotEqual(t, 0, data.Type["RandomValue"])
		})
	}
}
