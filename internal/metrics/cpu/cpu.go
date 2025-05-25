package cpu

import (
	"fmt"
	"time"

	syscpu "github.com/shirou/gopsutil/v4/cpu"
	sysload "github.com/shirou/gopsutil/v4/load"
)

type CpuMetric struct {
	Usage        []float64 // CPU usage as a percentage over a given time interval, each entry represents a core.
	TimeInterval float64   // The time interval for which usage percentage of the cpu is taken from.
	LoadAvg1     float64   // Average system load (number of processes running/waiting) over the past 1 minute.
	LoadAvg5     float64   // Average system load (number of processes running/waiting) over the past 5 minutes.
	LoadAvg15    float64   // Average system load (number of processes running/waiting) over the past 15 minutes.
	TimeStamp    time.Time // Time the measurement was taken.
}

type Option func(*float64)

func WithSeconds(seconds float64) Option {
	return func(input *float64) {
		if seconds < 0 || seconds > 180 {
			// TODO: log
			return
		}
		*input = seconds
		/// This seems unneccessarily comples just for an int...
	}
}

func getCpuMetrics(options ...Option) (CpuMetric, error) {
	seconds := 1.0
	for _, opt := range options {
		opt(&seconds)
	}

	percentages, err := syscpu.Percent(time.Duration(seconds)*time.Second, true)
	if err != nil {
		return CpuMetric{}, fmt.Errorf("Error getting CPU usage: %v", err)
	}

	loadAvg, err := sysload.Avg()
	if err != nil {
		return CpuMetric{}, fmt.Errorf("Error in getting load average: %v", err)
	}
	return CpuMetric{
		Usage:        percentages,
		TimeInterval: seconds,
		LoadAvg1:     loadAvg.Load1,
		LoadAvg5:     loadAvg.Load5,
		LoadAvg15:    loadAvg.Load15,
		TimeStamp:    time.Now(),
	}, nil
}

func (cm CpuMetric) String() string {
	retval := "Usage: "
	for _, percentage := range cm.Usage {
		retval += fmt.Sprintf("%.2f ", percentage)
	}
	retval += fmt.Sprintf("\nTimeInterval: %.2f\nLoadAvg1: %.2f\nLoadAvg5: %.2f\nLoadAvg15: %.2f\nTimeStamp: %v", cm.TimeInterval, cm.LoadAvg1, cm.LoadAvg5, cm.LoadAvg15, cm.TimeStamp)
	return retval
}
