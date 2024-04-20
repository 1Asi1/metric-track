package grpc

import (
	"context"
	"fmt"

	"github.com/1Asi1/metric-track.git/internal/server/service"
	proto "github.com/1Asi1/metric-track.git/rpc/gen"
)

type MetricGrpcService struct {
	proto.UnsafeMetricGrpcServer
	service service.Service
}

func NewMetricGrpcServer(service service.Service) *MetricGrpcService {
	return &MetricGrpcService{service: service}
}

func (s *MetricGrpcService) Updates(ctx context.Context, req *proto.UpdatesRequest) (*proto.UpdatesResponse, error) {
	model := make([]service.MetricsRequest, len(req.Metrics))
	for i, v := range req.Metrics {
		model[i] = service.MetricsRequest{
			ID:    v.ID,
			MType: v.MType,
			Delta: &v.Delta,
			Value: &v.Value,
		}
	}

	if err := s.service.Updates(ctx, model); err != nil {
		return &proto.UpdatesResponse{
			Error: err.Error(),
		}, fmt.Errorf("s.service.Updates: %w", err)
	}

	return nil, nil
}
