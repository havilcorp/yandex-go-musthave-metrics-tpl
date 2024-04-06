package metric

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/domain"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/repositories/metric/storage/file"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/repositories/metric/storage/memory"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/repositories/metric/storage/psql"
)

type Provider interface {
	AddGauge(ctx context.Context, key string, gauge float64) error
	AddCounter(ctx context.Context, key string, counter int64) error
	AddGaugeBulk(ctx context.Context, list []domain.Gauge) error
	AddCounterBulk(ctx context.Context, list []domain.Counter) error
	GetCounter(ctx context.Context, key string) (int64, error)
	GetGauge(ctx context.Context, key string) (float64, error)
	GetAllCounters(ctx context.Context) (map[string]int64, error)
	GetAllGauge(ctx context.Context) (map[string]float64, error)
}

func MetricFactory(provider string, conf *config.Config, db *sql.DB) (Provider, error) {
	switch provider {
	case "memory":
		return memory.NewMemStorage(), nil
	case "file":
		return file.NewFileStorage(conf)
	case "psql":
		return psql.NewPsqlStorage(conf, db)
	default:
		return nil, fmt.Errorf("unknown provider %s", provider)
	}
}
