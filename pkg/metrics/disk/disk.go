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

// DiskUsage has all values in bytes, except the field Usage which is a
// percentage.
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

// MeasureDiskMetrics is a wrapper for measureDiskUsage and measureDiskThroughput.
func MeasureDiskMetrics(diskName string, interval float64) (DiskMetric, error) {
	diskUsage, err := measureDiskUsage(diskName)
	if err != nil {
		return DiskMetric{}, err
	}
	diskThroughput, err := measureDiskThroughput(diskName, interval)
	if err != nil {
		return DiskMetric{}, err
	}
	return DiskMetric{
		DiskUsage:      diskUsage,
		DiskThroughput: diskThroughput,
		TimeStamp:      time.Now(),
	}, nil
}

// RetrieveDeviceMounts returns a mapping of storage devices and their corresponding
// mount points in the system. The keys represent phsyical paritions or storage
// devices. The values are the mountpoints of these physical paritions.
// This function can be used to see what available devices there are, then
// a user can pass the corresponding value to MeasureDiskMetrics.
func RetrieveDeviceMounts() (map[string]string, error) {
	partitions, err := sysdisk.Partitions(false) // False returns all physical devices.
	if err != nil {
		return map[string]string{}, fmt.Errorf("error when getting paritions: %v", err)
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

func measureDiskThroughput(blockDeviceName string, interval float64) (DiskThroughput, error) {
	ioStatsStart, err := sysdisk.IOCounters(blockDeviceName)
	if err != nil {
		return DiskThroughput{}, fmt.Errorf("error when getting start stats: %v", err)
	}

	startStat, exists := ioStatsStart[blockDeviceName]
	if !exists {
		return DiskThroughput{}, fmt.Errorf("disk name %q not found in start stat", blockDeviceName)
	}

	time.Sleep(time.Duration(interval) * time.Second)

	ioStatsEnd, err := sysdisk.IOCounters(blockDeviceName)
	if err != nil {
		return DiskThroughput{}, fmt.Errorf("error when getting end stats: %v", err)
	}
	endStat, exists := ioStatsEnd[blockDeviceName]
	if !exists {
		return DiskThroughput{}, fmt.Errorf("disk name %q not found in end stat", blockDeviceName)
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

func (dm DiskMetric) String() string {
	return fmt.Sprintf(
		"DiskUsage: {\nTotal: %.2f\nUsed: %.2f\nFree: %.2f\nUsage: %.2f\n}\n"+
			"DiskThroughput: {\nReadThroughput: %.2f\nWriteThroughput: %.2f\n"+
			"ReadOps: %.2f\nWriteOps: %.2f\nTotalIOPS: %.2f\nInterval: %.2f\n}\n"+
			"%v",
		dm.DiskUsage.Total, dm.DiskUsage.Used, dm.DiskUsage.Free, dm.DiskUsage.Usage,
		dm.DiskThroughput.ReadThroughput, dm.DiskThroughput.WriteThroughput,
		dm.DiskThroughput.ReadOps, dm.DiskThroughput.WriteOps, dm.DiskThroughput.TotalIOPS,
		dm.DiskThroughput.Interval, dm.TimeStamp,
	)
}
