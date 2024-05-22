package metric

import (
	"errors"
	"syscall"
	"testing"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config/agent"
	"github.com/stretchr/testify/require"
)

func TestWriteGopsutil(t *testing.T) {
	metric := NewMetric(nil)
	require.Empty(t, metric.String())
	metric.WriteGopsutil()
	require.NotEmpty(t, metric.String())
}

func TestWriteMain(t *testing.T) {
	metric := NewMetric(nil)
	require.Empty(t, metric.String())
	metric.WriteMain()
	require.NotEmpty(t, metric.String())
}

func TestAddClient(t *testing.T) {
	metric := NewMetric(nil)
	metric.AddMetricClient(nil)
}

func TestSend(t *testing.T) {
	conf := agent.NewAgentConfig()
	conf.ServerAddress = ":8080"
	metric := NewMetric(conf)
	err := metric.Send()
	if !errors.Is(err, syscall.ECONNREFUSED) {
		t.Error(err)
	}
}

func TestSendByGRPC(t *testing.T) {
	conf := agent.NewAgentConfig()
	metric := NewMetric(conf)
	err := metric.SendByGRPC()
	if err.Error() != "metricClient has a nil pointer" {
		t.Error(err)
	}
}
