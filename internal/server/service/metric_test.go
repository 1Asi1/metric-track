package service

import (
	"context"
	"testing"

	"github.com/1Asi1/metric-track.git/internal/server/repository/memory"
	"github.com/stretchr/testify/assert"
)

func Test_service_UpdateMetric(t *testing.T) {
	st := memory.New()
	srv := Service{Store: st}

	type args struct {
		ctx context.Context
		req Request
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
				req: Request{},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := srv.UpdateMetric(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err, nil)
		})
	}
}
