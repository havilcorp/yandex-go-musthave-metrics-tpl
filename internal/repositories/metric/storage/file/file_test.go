package file

import (
	"context"
	"testing"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/domain"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config"
	"github.com/stretchr/testify/require"
)

func TestFileStorage_AddGaugeBulk(t *testing.T) {
	conf := config.Config{
		StoreInterval: 999,
		IsRestore:     false,
	}
	store, err := NewFileStorage(&conf)
	if err != nil {
		t.Errorf("NewFileStorage %v", err)
	}
	list := []domain.Gauge{{
		Key:   "GAUGE1",
		Value: float64(1.1),
	}}
	t.Run("AddGaugeBulk", func(t *testing.T) {
		err := store.AddGaugeBulk(context.Background(), list)
		if err != nil {
			t.Errorf("AddGaugeBulk %v", err)
		}
		listGauges, err := store.GetAllGauge(context.Background())
		if err != nil {
			t.Errorf("GetAllGauge %v", err)
		}
		require.Equal(t, listGauges["GAUGE1"], float64(1.1))
		val, err := store.GetGauge(context.Background(), "GAUGE1")
		if err != nil {
			t.Errorf("GetGauge %v", err)
		}
		require.Equal(t, val, float64(1.1))
	})
}

func TestFileStorage_AddCounterBulk(t *testing.T) {
	conf := config.Config{
		StoreInterval: 999,
		IsRestore:     false,
	}
	store, err := NewFileStorage(&conf)
	if err != nil {
		t.Errorf("NewFileStorage %v", err)
	}
	list := []domain.Counter{{
		Key:   "COUNTER1",
		Value: int64(1),
	}}
	t.Run("AddCounterBulk", func(t *testing.T) {
		err := store.AddCounterBulk(context.Background(), list)
		if err != nil {
			t.Errorf("AddCounterBulk %v", err)
		}
		listCounters, err := store.GetAllCounters(context.Background())
		if err != nil {
			t.Errorf("GetAllGauge %v", err)
		}
		require.Equal(t, listCounters["COUNTER1"], int64(1))
		val, err := store.GetCounter(context.Background(), "COUNTER1")
		if err != nil {
			t.Errorf("GetCounter %v", err)
		}
		require.Equal(t, val, int64(1))
	})
}
