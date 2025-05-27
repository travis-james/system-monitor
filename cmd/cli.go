package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/travis-james/system-monitor/pkg/metrics/cpu"
)

func RunCLI() {
	metricsToCollect := flag.String("metric", "", "metrics to retrieve (cpu, disk)")
	seconds := flag.Float64("seconds", 5, "Duration to measure metric(s) where applicable")
	flag.Parse()

	if *metricsToCollect == "" {
		fmt.Println("no metric was chosen (ex: -metric=cpu,disk)")
		os.Exit(1)
	}

	metricsType := strings.Split(*metricsToCollect, ",")
	for _, metric := range metricsType {
		switch metric {
		case "cpu":
			cpuMetrics, err := cpu.MeasureCpuMetrics(*seconds)
			if err != nil {
				fmt.Println("Error measuring CPU:", err)
			} else {
				fmt.Printf("CPU Metrics: %s\n", cpuMetrics.String())
			}
		default:
			fmt.Printf("Invalid metric type: %s\n", metric)
		}
	}
}
