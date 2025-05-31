package disk

import (
	"fmt"
	"testing"
	"time"

	gopsutilDisk "github.com/shirou/gopsutil/v4/disk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// func TestMeasureDiskMetrics(t *testing.T) {

// }

func TestRetrieveDeviceMounts(t *testing.T) {
	t.Parallel()
	mockPartitions := func(_ bool) ([]gopsutilDisk.PartitionStat, error) {
		return []gopsutilDisk.PartitionStat{
			{Device: "/dev/nvme01", Mountpoint: "/"},
			{Device: "/dev/nvme02", Mountpoint: "/mnt"},
		}, nil
	}
	got, err := retrieveDeviceMounts(mockPartitions)
	require.Nil(t, err)
	expected := map[string]string{
		"/dev/nvme01": "/",
		"/dev/nvme02": "/mnt",
	}
	assert.Equal(t, fmt.Sprintf("%v", got), fmt.Sprintf("%v", expected))
}

func TestMeasureDiskUsage(t *testing.T) {
	t.Parallel()
	var (
		total       uint64  = 54
		used        uint64  = 44
		free        uint64  = 34
		usedPercent float64 = 2.4
	)
	mockUsage := func(_ string) (*gopsutilDisk.UsageStat, error) {
		return &gopsutilDisk.UsageStat{
			Total:       total,
			Used:        used,
			Free:        free,
			UsedPercent: usedPercent,
		}, nil
	}
	got, err := measureDiskUsage(mockUsage, "/")
	require.Nil(t, err)
	assert.Equal(t, got.Total, total)
	assert.Equal(t, got.Free, free)
	assert.Equal(t, got.Used, used)
	assert.Equal(t, got.Usage, usedPercent)
}

func TestGetDiskThroughput(t *testing.T) {
	t.Parallel()
	mockUsage := func() ioCountersFunc {
		count := 0
		return func(...string) (map[string]gopsutilDisk.IOCountersStat, error) {
			count += 1000 // Simulating an increase in disk stats
			return map[string]gopsutilDisk.IOCountersStat{
				"mockDisk": {
					ReadBytes:  uint64(count),
					WriteBytes: uint64(count * 2),
					ReadCount:  uint64(count / 100),
					WriteCount: uint64(count / 200),
				},
			}, nil
		}
	}
	var time float64 = 0.1

	got, err := measureDiskThroughput(mockUsage(), "mockDisk", time)
	require.Nil(t, err)
	assert.Greater(t, got.WriteThroughput, 0.0)
	assert.Greater(t, got.ReadThroughput, 0.0)
	assert.Greater(t, got.WriteOps, 0.0)
	assert.Greater(t, got.ReadOps, 0.0)
	assert.Greater(t, got.TotalIOPS, 0.0)
	assert.Equal(t, got.Interval, time)
}

func TestString(t *testing.T) {
	dm := DiskMetric{
		DiskUsage: DiskUsage{
			Total: 5,
			Used:  4,
			Free:  3,
			Usage: 0.9,
		},
		DiskThroughput: DiskThroughput{
			ReadThroughput:  99,
			WriteThroughput: 98,
			ReadOps:         97,
			WriteOps:        96,
			TotalIOPS:       95,
			Interval:        1.5,
		},
		TimeStamp: time.Now(),
	}
	t.Log(dm.String())
}
