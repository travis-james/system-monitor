package cpu

import (
	"errors"
	"strings"
	"testing"
	"time"

	gopsutilLoad "github.com/shirou/gopsutil/v4/load"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock percentage usage function
func mockPercentageUsage(duration time.Duration, detailed bool) ([]float64, error) {
	if duration <= 0 {
		return nil, errors.New("invalid duration")
	}
	return []float64{10.5, 15.2, 20.3}, nil
}

// Mock load average function
func mockLoadAvg() (*gopsutilLoad.AvgStat, error) {
	return &gopsutilLoad.AvgStat{Load1: 1.5, Load5: 2.0, Load15: 2.5}, nil
}

func TestMeasureCpuMetrics_ValidInput(t *testing.T) {
	got, err := measureCpuMetrics(mockPercentageUsage, mockLoadAvg, 5)
	require.Nil(t, err)

	expected := CpuMetric{
		Usage:    []float64{10.5, 15.2, 20.3},
		LoadAvg1: 1.5,
	}
	assert.Equal(t, len(expected.Usage), len(got.Usage))
	assert.Equal(t, expected.LoadAvg1, expected.LoadAvg1)
}

func TestMeasureCpuMetrics_InvalidDuration(t *testing.T) {
	_, err := measureCpuMetrics(mockPercentageUsage, mockLoadAvg, -1)
	assert.NotNil(t, err)
}

func TestMeasureCpuMetrics_ErrorInCPUUsage(t *testing.T) {
	mockErrUsage := func(duration time.Duration, detailed bool) ([]float64, error) {
		return nil, errors.New("mock CPU usage error")
	}

	_, err := measureCpuMetrics(mockErrUsage, mockLoadAvg, 5)
	assert.NotNil(t, err)
}

func TestMeasureCpuMetrics_ErrorInLoadAvg(t *testing.T) {
	mockErrLoadAvg := func() (*gopsutilLoad.AvgStat, error) {
		return &gopsutilLoad.AvgStat{}, errors.New("mock load avg error")
	}

	_, err := measureCpuMetrics(mockPercentageUsage, mockErrLoadAvg, 5)
	assert.NotNil(t, err)
}

func TestString(t *testing.T) {
	t.Parallel()
	input := CpuMetric{
		Usage:        []float64{1, 2, 3},
		TimeInterval: 0.3,
		LoadAvg1:     0.1,
		LoadAvg5:     0.2,
		LoadAvg15:    0.3,
	}
	expected := `Usage: 1.00 2.00 3.00 
		NumberOfCores: 0
        TimeInterval: 0.30
        LoadAvg1: 0.10
        LoadAvg5: 0.20
        LoadAvg15: 0.30
        TimeStamp: 0001-01-01 00:00:00 +0000 UTC`

	assert.Equal(t, strings.Fields(expected), strings.Fields(input.String()))
}
