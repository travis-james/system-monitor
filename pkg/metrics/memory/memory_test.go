package memory

import (
	"errors"
	"testing"

	gopsutilMem "github.com/shirou/gopsutil/v4/mem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock function for testing
func mockVirtualMemory(stats *gopsutilMem.VirtualMemoryStat, err error) virtualMemoryFunc {
	return func() (*gopsutilMem.VirtualMemoryStat, error) {
		return stats, err
	}
}

func TestMeasureMemoryMetrics_ValidStats(t *testing.T) {
	var (
		used      uint64 = 2048
		available uint64 = 4096
		mockStats        = &gopsutilMem.VirtualMemoryStat{
			Used:      used,
			Available: available,
		}
	)
	got, err := measureMemoryMetrics(mockVirtualMemory(mockStats, nil))
	require.Nil(t, err)
	assert.Equal(t, used, got.UsedMemory)
	assert.Equal(t, available, got.AvailableMemory)
	assert.NotZero(t, got.TimeStamp)
}

func TestMeasureMemoryMetrics_ErrorCase(t *testing.T) {
	_, err := measureMemoryMetrics(mockVirtualMemory(nil, errors.New("failed to get memory stats")))
	assert.NotNil(t, err)
}
