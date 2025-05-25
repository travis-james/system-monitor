package cpu

import (
	"fmt"
	"time"

	syscpu "github.com/shirou/gopsutil/v4/cpu"
)

type CpuMetric struct {
	CpuUsages []CpuUsage
	LoadAvg   float64
	TimeStamp time.Time
}

type CpuUsage struct {
	CpuName string
	Usage   float64
}

func getCpuMetrics() (CpuMetric, error) {
	time := time.Now()

	percentages, err := syscpu.Percent(5, true)
	if err != nil {
		return CpuMetric{}, err
	}
	var cpuUsages []CpuUsage
	var name string
	for index, percent := range percentages {
		if index != 0 {
			name = fmt.Sprintf("cpu%d", index)
		} else {
			name = "cpuOverall"
		}
		cpuUsages = append(cpuUsages, CpuUsage{
			CpuName: name,
			Usage:   percent,
		})
	}
	loadAvg, err := syscpu.LoadAvg()
	return CpuMetric{
		CpuUsages: cpuUsages,

		TimeStamp: time,
	}, nil
}

func getCpuUsage() ([]float64, error) {
	percentages, err := syscpu.Percent(5, true)
	if err != nil {
		return []float64{}, err
	}
	return percentages, nil
}
