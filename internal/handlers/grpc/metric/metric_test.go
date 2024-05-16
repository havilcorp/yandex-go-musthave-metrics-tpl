package metric

import (
	"context"
	"testing"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/domain"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/mocks"
	pb "github.com/havilcorp/yandex-go-musthave-metrics-tpl/pkg/proto/metric"
	"github.com/stretchr/testify/mock"
)

func TestMetricServer_UpdateMetricBulk(t *testing.T) {
	metricHandler := mocks.NewIMetric(t)

	metricHandler.On("AddGaugeBulk", mock.Anything, []domain.Gauge{{Key: "GAUGE", Value: float64(1.1)}}).Return(nil)
	metricHandler.On("AddCounterBulk", mock.Anything, []domain.Counter{{Key: "COUNTER", Value: int64(1)}}).Return(nil)

	gauge := make([]*pb.Gauge, 0)
	counter := make([]*pb.Counter, 0)

	gauge = append(gauge, &pb.Gauge{
		Key:   "GAUGE",
		Value: float64(1.1),
	})

	counter = append(counter, &pb.Counter{
		Key:   "COUNTER",
		Value: int64(1),
	})

	h := NewMetric(metricHandler)
	_, err := h.UpdateMetricBulk(context.Background(), &pb.UpdateMetricBulkRequest{
		Gauge:   gauge,
		Counter: counter,
	})
	if err != nil {
		t.Error(err)
	}
}
