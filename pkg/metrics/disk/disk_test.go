package disk

import (
	"fmt"
	"testing"

	gopsutilDisk "github.com/shirou/gopsutil/v4/disk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// func TestMeasureDiskMetrics(t *testing.T) {

// }

func TestRetrieveDeviceMounts(t *testing.T) {
	mockPartitions := func(_ bool) ([]gopsutilDisk.PartitionStat, error) {
		return []gopsutilDisk.PartitionStat{
			{Device: "/dev/nvme01", Mountpoint: "/"},
			{Device: "/dev/nvme02", Mountpoint: "/mnt"},
		}, nil
	}
	t.Parallel()
	got, err := retrieveDeviceMounts(mockPartitions)
	require.Nil(t, err)
	expected := map[string]string{
		"/dev/nvme01": "/",
		"/dev/nvme02": "/mnt",
	}
	assert.Equal(t, fmt.Sprintf("%v", got), fmt.Sprintf("%v", expected))
}

// func TestGetDiskUsage(t *testing.T) {
// 	t.Parallel()
// 	got, err := measureDiskUsage("/")
// 	require.Nil(t, err)
// 	assert.Greater(t, got.Total, 0.0)
// 	assert.Greater(t, got.Free, 0.0)
// 	assert.Greater(t, got.Used, 0.0)
// 	assert.Greater(t, got.Usage, 0.0)
// }

// Takes a long time to get write data....
// func TestGetDiskThroughput(t *testing.T) {
// 	t.Parallel()
// 	var time float64 = 5

// 	got, err := measureDiskThroughput("nvme0n1", time)
// 	require.Nil(t, err)
// 	assert.Greater(t, got.WriteThroughput, 0.0)
// 	assert.Greater(t, got.ReadThroughput, 0.0)
// 	assert.Greater(t, got.WriteOps, 0.0)
// 	assert.Greater(t, got.ReadOps, 0.0)
// 	assert.Greater(t, got.TotalIOPS, 0.0)
// 	assert.Equal(t, got.Interval, time)
// 	t.Log(DiskMetric{DiskThroughput: got}.String())
// }
