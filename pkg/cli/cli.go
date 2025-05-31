package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/travis-james/system-monitor/pkg/metrics/cpu"
	"github.com/travis-james/system-monitor/pkg/metrics/disk"
)

func RunCLI() {
	metricsToCollect := flag.String("metric", "", "metrics to retrieve (cpu, disk)")
	seconds := flag.Float64("seconds", 5, "Duration to measure metric(s) where applicable")
	md := flag.String("md", "", "mounted directory")
	flag.Parse()

	if *metricsToCollect == "" {
		fmt.Println("no metric was chosen (ex: -metric=cpu,disk)")
		os.Exit(1)
	}

	metricsType := strings.SplitSeq(*metricsToCollect, ",")
	for metric := range metricsType {
		switch metric {
		case "cpu":
			cpuMetrics, err := cpu.MeasureCpuMetrics(*seconds)
			if err != nil {
				fmt.Println("Error measuring CPU:", err)
				os.Exit(1)
			} else {
				fmt.Printf("CPU Metrics: %s\n", cpuMetrics.String())
			}
		case "deviceMounts":
			deviceMounts, err := disk.RetrieveDeviceMountsToString()
			if err != nil {
				fmt.Println("Error getting device mounts:", err)
			} else {
				fmt.Printf("device mounts:\n%s\n", deviceMounts)
			}
		case "disk":
			diskMetrics, err := disk.MeasureDiskMetrics(*md, *seconds)
			if err != nil {
				fmt.Println("Error getting disk metrics:", err)
				os.Exit(1)
			} else {
				fmt.Printf("Disk Metrics: %s\n", diskMetrics)
			}
		default:
			fmt.Printf("Invalid metric type: %s\n", metric)
		}
	}
}
