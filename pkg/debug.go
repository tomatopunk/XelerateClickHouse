package pkg

import (
	"fmt"
	"sort"
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
		dam.distributeInfo[metric.Timestamp] = v + 1
	} else {
		dam.distributeInfo[metric.Timestamp] = 1
	}
}

func (dam *DebugAppendMetrics) Printf() {
	// Create a temporary slice for sorting
	temp := make([]struct {
		timestamp time.Time
		count     int
	}, 0, len(dam.distributeInfo))

	// Store timestamps and counts in the temporary slice
	for k, v := range dam.distributeInfo {
		temp = append(temp, struct {
			timestamp time.Time
			count     int
		}{
			timestamp: k,
			count:     v,
		})
	}

	// Sort by timestamps
	sort.Slice(temp, func(i, j int) bool {
		return temp[i].timestamp.Before(temp[j].timestamp)
	})

	// Print timestamps and counts in ascending order
	for _, entry := range temp {
		fmt.Printf("Timestamp: %v, Count: %d\n", entry.timestamp, entry.count)
	}
}
