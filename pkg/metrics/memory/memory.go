package memory

import (
	"time"

	gopsutilMem "github.com/shirou/gopsutil/v4/mem"
)

type MemoryMetric struct {
	UsedMemory      uint64
	AvailableMemory uint64
	TimeStamp       time.Time
}

func MeasureMemoryMetrics() (MemoryMetric, error) {
	return measureMemoryMetrics(gopsutilMem.VirtualMemory)
}

type virtualMemoryFunc func() (*gopsutilMem.VirtualMemoryStat, error)

func measureMemoryMetrics(getVirtualMemory virtualMemoryFunc) (MemoryMetric, error) {
	memStats, err := getVirtualMemory()
	if err != nil {
		return MemoryMetric{}, err
	}
	return MemoryMetric{
		UsedMemory:      memStats.Used,
		AvailableMemory: memStats.Available,
		TimeStamp:       time.Now(),
	}, nil
}
