package cpu

import (
	"errors"
	"fmt"
	"time"

	gopsutilCPU "github.com/shirou/gopsutil/v4/cpu"
	gopsutilLoad "github.com/shirou/gopsutil/v4/load"
)

const ERR_INVALID_SECONDS = "seconds must be greater than zero"

// CpuMetric contains data for usage (how busy each core is) and load average (how much demand there is for cpu resources)
type CpuMetric struct {
	Usage         []float64 // CPU usage as a percentage over a given time interval, each entry represents a core.
	NumberOfCores int       // Number of cores the CPU has.
	TimeInterval  float64   // The time interval for which usage percentage of the cpu is taken from.
	LoadAvg1      float64   // Average system load (number of processes running/waiting) over the past 1 minute.
	LoadAvg5      float64   // Average system load (number of processes running/waiting) over the past 5 minutes.
	LoadAvg15     float64   // Average system load (number of processes running/waiting) over the past 15 minutes.
	TimeStamp     time.Time // Time the measurement was taken.
}

// MeasureCpuMetrics is the public wrapper for measureCpuMetrics.
// Will get all related cpu metrics and return CpuMetric.
func MeasureCpuMetrics(seconds float64) (CpuMetric, error) {
	return measureCpuMetrics(gopsutilCPU.Percent, gopsutilLoad.Avg, seconds)
}

// percentFunc is dependency injection for measureCpuMetrics and
// gopsutilCPU.Percent.
type percentFunc func(time.Duration, bool) ([]float64, error)

// loadAvgFunc is dependency injection for measureCpuMetrics and
// gopsutilLoad.Avg.
type loadAvgFunc func() (*gopsutilLoad.AvgStat, error)

// measureCpuMetrics gets all related cpu metrics to put them
// in a CpuMetric struct.
func measureCpuMetrics(getPercentageUsage percentFunc, getLoadAvg loadAvgFunc, seconds float64) (CpuMetric, error) {
	if seconds <= 0 {
		return CpuMetric{}, errors.New(ERR_INVALID_SECONDS)
	}
	percentages, err := getPercentageUsage(time.Duration(seconds)*time.Second, true)
	if err != nil {
		return CpuMetric{}, fmt.Errorf("error getting CPU usage: %v", err)
	}

	loadAvg, err := getLoadAvg()
	if err != nil {
		return CpuMetric{}, fmt.Errorf("error in getting load average: %v", err)
	}
	return CpuMetric{
		Usage:         percentages,
		NumberOfCores: len(percentages),
		TimeInterval:  seconds,
		LoadAvg1:      loadAvg.Load1,
		LoadAvg5:      loadAvg.Load5,
		LoadAvg15:     loadAvg.Load15,
		TimeStamp:     time.Now(),
	}, nil
}

// String returns a string representation of CpuMetric.
func (cm CpuMetric) String() string {
	retval := "Usage: "
	for _, percentage := range cm.Usage {
		retval += fmt.Sprintf("%.2f ", percentage)
	}
	retval += fmt.Sprintf("\nNumberOfCores: %d\nTimeInterval: %.2f\nLoadAvg1: %.2f\nLoadAvg5: %.2f\nLoadAvg15: %.2f\nTimeStamp: %v", cm.NumberOfCores, cm.TimeInterval, cm.LoadAvg1, cm.LoadAvg5, cm.LoadAvg15, cm.TimeStamp)
	return retval
}
