package metric

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteGopsutil(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		metric := NewMetric(nil)
		require.Empty(t, metric.String())
		metric.WriteGopsutil()
		require.NotEmpty(t, metric.String())
	})
}

func TestWriteMain(t *testing.T) {
	t.Run("test2", func(t *testing.T) {
		metric := NewMetric(nil)
		require.Empty(t, metric.String())
		metric.WriteMain()
		require.NotEmpty(t, metric.String())
	})
}
