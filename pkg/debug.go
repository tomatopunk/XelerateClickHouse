//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package pkg

import (
	"sort"
	"sync"
	"time"

	"clickhouse-benchmark/pkg/show"
)

type DebugAppendMetrics struct {
	sync.Mutex
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
		show.Debug("Timestamp: %v, Count: %d\n", entry.timestamp, entry.count)
	}
}
