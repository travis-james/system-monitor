package disk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDiskName(t *testing.T) {
	t.Parallel()
	got, err := getDiskNames()
	require.Nil(t, err)
	assert.NotZero(t, got)
	assert.Greater(t, len(got), 0)
}

func TestGetDiskUsage(t *testing.T) {
	t.Parallel()
	got, err := getDiskUsage("/")
	require.Nil(t, err)
	assert.Greater(t, got.Total, 0.0)
	assert.Greater(t, got.Free, 0.0)
	assert.Greater(t, got.Used, 0.0)
	assert.Greater(t, got.Usage, 0.0)
}

// Takes a long time to get write data....
func TestGetDiskThroughput(t *testing.T) {
	t.Parallel()
	var time float64 = 20
	got, err := measureDiskThroughput("nvme0n1", time)
	require.Nil(t, err)
	assert.Greater(t, got.WriteThroughput, 0.0)
	assert.Greater(t, got.ReadThroughput, 0.0)
	assert.Greater(t, got.WriteOps, 0.0)
	assert.Greater(t, got.ReadOps, 0.0)
	assert.Greater(t, got.TotalIOPS, 0.0)
	assert.Equal(t, got.Interval, time)
	t.Log(got)
}
