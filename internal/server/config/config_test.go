package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func TestNew(t *testing.T) {
	type args struct {
		log zerolog.Logger
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name: "positive",
			want: Config{
				MetricServerAddr: "localhost:8080",
				MetricPPROFAddr:  "localhost:8081",
				StoreInterval:    5 * time.Minute,
				StorePath:        "./tmp/metrics-db.json",
				StoreRestore:     true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}
