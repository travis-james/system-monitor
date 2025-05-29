package disk

import (
	"fmt"
	"time"

	sysdisk "github.com/shirou/gopsutil/v4/disk"
)

type DiskMetric struct {
	DiskUsage
	DiskThroughput
	TimeStamp time.Time // Time the measurement was taken.
}

// DiskUsage has all values in bytes, except the field
// Usage which is a percentage.
type DiskUsage struct {
	Total float64
	Used  float64
	Free  float64
	Usage float64
}

type DiskThroughput struct {
	ReadThroughput  float64
	WriteThroughput float64
	ReadOps         float64
	WriteOps        float64
	TotalIOPS       float64
	Interval        float64
}

func MeasureDiskMetrics(diskName string) (DiskMetric, error) {
	return DiskMetric{}, nil
}

// RetrieveDeviceMounts returns a map that has it's keys (device name/path)
// mapped to a mounted file system.
func RetrieveDeviceMounts() (map[string]string, error) {
	partitions, err := sysdisk.Partitions(false) // False returns all physical devices.
	if err != nil {
		return map[string]string{}, nil
	}
	// Map devices to their mount points
	deviceMap := make(map[string]string)
	for _, p := range partitions {
		deviceMap[p.Device] = p.Mountpoint
	}
	return deviceMap, nil
}

func measureDiskUsage(diskName string) (DiskUsage, error) {
	usage, err := sysdisk.Usage(diskName)
	if err != nil {
		return DiskUsage{}, err
	}
	return DiskUsage{
		Total: float64(usage.Total),
		Used:  float64(usage.Used),
		Free:  float64(usage.Free),
		Usage: usage.UsedPercent,
	}, nil
}

func measureDiskThroughput(diskName string, interval float64) (DiskThroughput, error) {
	ioStatsStart, err := sysdisk.IOCounters(diskName)
	if err != nil {
		return DiskThroughput{}, fmt.Errorf("error when getting start stats: %v", err)
	}
	startStat, exists := ioStatsStart[diskName]
	if !exists {
		return DiskThroughput{}, fmt.Errorf("disk name %q not found in start stat", diskName)
	}

	time.Sleep(time.Duration(interval) * time.Second)

	ioStatsEnd, err := sysdisk.IOCounters(diskName)
	if err != nil {
		return DiskThroughput{}, fmt.Errorf("error when getting end stats: %v", err)
	}
	endStat, exists := ioStatsEnd[diskName]
	if !exists {
		return DiskThroughput{}, fmt.Errorf("disk name %q not found in end stat", diskName)
	}

	readOps := float64(endStat.ReadCount-startStat.ReadCount) / interval
	writeOps := float64(endStat.WriteCount-startStat.WriteCount) / interval

	return DiskThroughput{
		ReadThroughput:  float64(endStat.ReadBytes-startStat.ReadBytes) / interval,
		WriteThroughput: float64(endStat.WriteBytes-startStat.WriteBytes) / interval,
		ReadOps:         readOps,
		WriteOps:        writeOps,
		TotalIOPS:       readOps + writeOps,
		Interval:        interval,
	}, nil
}
