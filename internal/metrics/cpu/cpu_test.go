package cpu

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCpuMetrics(t *testing.T) {
	got, err := getCpuMetrics(WithSeconds(0.2))
	require.Nil(t, err)
	assert.Equal(t, 0.2, got.TimeInterval)
	assert.False(t, got.TimeStamp.IsZero())
	assert.NotZero(t, len(got.Usage))
	assert.Greater(t, got.LoadAvg1, 0.0)
	assert.Greater(t, got.LoadAvg5, 0.0)
	assert.Greater(t, got.LoadAvg15, 0.0)
}

func TestString(t *testing.T) {
	input := CpuMetric{
		Usage:        []float64{1, 2, 3},
		TimeInterval: 0.3,
		LoadAvg1:     0.1,
		LoadAvg5:     0.2,
		LoadAvg15:    0.3,
	}
	expected := `Usage: 1.00 2.00 3.00 
        TimeInterval: 0.30
        LoadAvg1: 0.10
        LoadAvg5: 0.20
        LoadAvg15: 0.30
        TimeStamp: 0001-01-01 00:00:00 +0000 UTC`

	assert.Equal(t, strings.Fields(expected), strings.Fields(input.String()))
}
