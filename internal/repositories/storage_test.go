// Package repositories репозиторий для сохранения данных метрик
package repositories

import (
	"testing"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/domain"
)

func TestStorage_SaveToFile(t *testing.T) {
	type args struct {
		metric domain.Metric
		path   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test error path",
			args: args{
				path: "/tmp-notfound/test-metrics-db.json",
				metric: domain.Metric{
					Gauge:   map[string]float64{"GAUGE": 1.1},
					Counter: map[string]int64{"COUNTER": 1},
				},
			},
			wantErr: true,
		},
		{
			name: "Test good",
			args: args{
				path: "/tmp/test-metrics-db.json",
				metric: domain.Metric{
					Gauge:   map[string]float64{"GAUGE": 1.1},
					Counter: map[string]int64{"COUNTER": 1},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewStorage(tt.args.path)
			if err := store.SaveToFile(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("Storage.SaveToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
