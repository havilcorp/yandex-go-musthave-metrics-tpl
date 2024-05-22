package file

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/domain"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config/server"
	"github.com/stretchr/testify/require"
)

func TestFileStorage_AddGaugeBulk(t *testing.T) {
	conf := server.Config{
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
	err = store.AddGaugeBulk(context.Background(), list)
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
}

func TestFileStorage_AddCounterBulk(t *testing.T) {
	conf := server.Config{
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
	err = store.AddCounterBulk(context.Background(), list)
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
}

func TestFileStorage_SaveToFile_LoadFromFile(t *testing.T) {
	conf := server.Config{
		StoreInterval:   999,
		IsRestore:       false,
		FileStoragePath: "/tmp/test-metrics-db.json",
	}
	store, err := NewFileStorage(&conf)
	if err != nil {
		t.Errorf("NewFileStorage %v", err)
	}
	err = store.AddGauge(context.Background(), "GAUGE1", 1.1)
	if err != nil {
		t.Errorf("AddGauge %v", err)
	}
	err = store.AddCounter(context.Background(), "COUNTER1", 1)
	if err != nil {
		t.Errorf("AddCounter %v", err)
	}
	err = store.SaveToFile(context.Background())
	if err != nil {
		t.Errorf("SaveToFile %v", err)
	}
	err = store.LoadFromFile(context.Background())
	if err != nil {
		t.Errorf("LoadFromFile %v", err)
	}
	confForTestRestore := server.Config{
		StoreInterval:   999,
		IsRestore:       true,
		FileStoragePath: "/tmp/test-metrics-db.json",
	}
	_, err = NewFileStorage(&confForTestRestore)
	if err != nil {
		t.Errorf("NewFileStorage %v", err)
	}
}

func TestFileStorage_AddGauge(t *testing.T) {
	conf := server.Config{
		StoreInterval:   0,
		IsRestore:       false,
		FileStoragePath: "/tmp/test-metrics-db.json",
	}
	store, err := NewFileStorage(&conf)
	if err != nil {
		t.Errorf("NewFileStorage %v", err)
	}
	err = store.AddGauge(context.Background(), "GAUGE", 1.1)
	if err != nil {
		t.Errorf("AddGauge %v", err)
	}
}

func TestFileStorage_AddCounter(t *testing.T) {
	conf := server.Config{
		StoreInterval:   0,
		IsRestore:       false,
		FileStoragePath: "/tmp/test-metrics-db.json",
	}
	store, err := NewFileStorage(&conf)
	if err != nil {
		t.Errorf("NewFileStorage %v", err)
	}
	err = store.AddCounter(context.Background(), "COUNTER", 1)
	if err != nil {
		t.Errorf("AddCounter %v", err)
	}
	err = store.AddCounter(context.Background(), "COUNTER", 1)
	if err != nil {
		t.Errorf("AddCounter %v", err)
	}
}

func TestFileStorage_LoadFromFile(t *testing.T) {
	conf := server.Config{
		StoreInterval:   999,
		IsRestore:       false,
		FileStoragePath: "/tmp/not-found.json",
	}
	store, err := NewFileStorage(&conf)
	if err != nil {
		t.Errorf("NewFileStorage %v", err)
	}
	err = store.LoadFromFile(context.Background())
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("LoadFromFile %v", err)
	}
}
