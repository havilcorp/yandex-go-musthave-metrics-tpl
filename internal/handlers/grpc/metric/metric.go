// Package metric grpc соединение
package metric

import (
	"context"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/domain"
	pb "github.com/havilcorp/yandex-go-musthave-metrics-tpl/pkg/proto/metric"
)

type IMetric interface {
	AddGaugeBulk(ctx context.Context, list []domain.Gauge) error
	AddCounterBulk(ctx context.Context, list []domain.Counter) error
}

type MetricServer struct {
	pb.UnimplementedMetricServer
	metric IMetric
}

func NewMetric(metric IMetric) *MetricServer {
	return &MetricServer{
		metric: metric,
	}
}

func (s *MetricServer) UpdateMetricBulk(ctx context.Context, in *pb.UpdateMetricBulkRequest) (*pb.UpdateMetricBulkResponse, error) {
	var response pb.UpdateMetricBulkResponse
	gauge := make([]domain.Gauge, 0)
	for _, g := range in.Gauge {
		gauge = append(gauge, domain.Gauge{
			Key:   g.Key,
			Value: g.Value,
		})
	}
	counter := make([]domain.Counter, 0)
	for _, c := range in.Counter {
		counter = append(counter, domain.Counter{
			Key:   c.Key,
			Value: c.Value,
		})
	}
	if err := s.metric.AddGaugeBulk(ctx, gauge); err != nil {
		return &response, err
	}
	if err := s.metric.AddCounterBulk(ctx, counter); err != nil {
		return &response, err
	}
	return &response, nil
}
