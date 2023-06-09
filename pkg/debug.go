package pkg

import (
	"fmt"
	"time"
)

type DebugAppendMetrics struct {
	metrics        []Metric
	distributeInfo map[time.Time]int
}

func NewDebugAppendMetrics() *DebugAppendMetrics {
	return &DebugAppendMetrics{distributeInfo: make(map[time.Time]int)}
}

func (dam *DebugAppendMetrics) Add(metric Metric) {
	dam.metrics = append(dam.metrics, metric)
	if v, ok := dam.distributeInfo[metric.Timestamp]; ok {
		v++
	} else {
		dam.distributeInfo[metric.Timestamp] = 1
	}
}

func (dam *DebugAppendMetrics) Printf() {
	for k, count := range dam.distributeInfo {
		fmt.Printf("Timestamp: %v, Count: %d\n", k, count)
	}
}
