package cpu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCpuUsage(t *testing.T) {
	got, err := getCpuUsage()
	require.Nil(t, err)
	t.Log(got)
}
