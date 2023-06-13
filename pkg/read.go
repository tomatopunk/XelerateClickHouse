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
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"clickhouse-benchmark/pkg/show"

	"github.com/montanaflynn/stats"
	"github.com/spf13/cobra"
)

type readOption struct {
	startTime string
	endTime   string
	timeStep  string
	sql       string
}

var readOpt readOption

var readCommand = &cobra.Command{
	Use:  "read",
	Long: ` benchmarking read `,
	Run: func(cmd *cobra.Command, args []string) {
		if err := benchmarkReadQueries(); err != nil {
			show.Error("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	root.AddCommand(readCommand)

	readCommand.Flags().StringVar(&readOpt.startTime, "start", "2023-06-09 18:00:00", "start time")
	readCommand.Flags().StringVar(&readOpt.endTime, "end", "2023-06-09 19:00:00", "end time")
	readCommand.Flags().StringVar(&readOpt.timeStep, "step", "minute", "time step")
	readCommand.Flags().StringVar(&readOpt.sql, "sql", "select * from test.metrics", "sql")

}

func benchmarkReadQueries() error {
	// Parse start and end times
	startTime, err := time.Parse("2006-01-02 15:04:05", readOpt.startTime)
	if err != nil {
		return err
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", readOpt.endTime)
	if err != nil {
		return err
	}

	// Calculate the number of iterations based on the time step
	var duration time.Duration
	switch readOpt.timeStep {
	case "day":
		duration = 24 * time.Hour
	case "hour":
		duration = time.Hour
	case "minute":
		duration = time.Minute
	case "second":
		duration = time.Second
	default:
		return fmt.Errorf("invalid time step: %s", readOpt.timeStep)
	}

	iterations := int(endTime.Sub(startTime) / duration)
	failedQuery := 0
	taskStart := time.Now()

	// Construct time condition
	timeCondition := fmt.Sprintf("timestamp > '%s' AND timestamp < '%s'", startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"))

	conn, err := getConn(os.Getenv("CLICKHOUSE_URL"))
	if err != nil {
		return err
	}
	defer conn.Close()

	results := make(map[int]float64)

	for i := 1; i <= iterations; i++ {
		t := startTime.Add(duration * time.Duration(i))
		query := fmt.Sprintf("%s AND timestamp = '%s'", timeCondition, t.Format("2006-01-02 15:04:05"))
		start := time.Now()
		rows, err := conn.Query(context.Background(), readOpt.sql+" WHERE "+query)
		if err != nil {
			failedQuery++
			return err
		}
		defer rows.Close()

		// Calculate query elapsed time
		elapsed := time.Since(start).Seconds()
		results[i] = elapsed
	}

	totalTime := time.Since(taskStart)

	printResults(results)

	bucketSlice := mapToSlice(results)

	p50, _ := stats.Percentile(bucketSlice, 50)
	p80, _ := stats.Percentile(bucketSlice, 80)
	p99, _ := stats.Percentile(bucketSlice, 99)
	p999, _ := stats.Percentile(bucketSlice, 99.9)

	show.Info("p50: %v, p80: %v, p99: %v, p999: %v\n", p50, p80, p99, p999)

	// Print benchmarking results
	show.Info("\n\n")
	show.Info("ClickHouse URL: %s\n", os.Getenv("CLICKHOUSE_URL"))
	show.Info("Total queries executed: %d\n", iterations)
	show.Info("Failed requests: %d\n", failedQuery)
	show.Info("Time taken for tests: %v\n", totalTime)

	return nil
}

func mapToSlice(results map[int]float64) []float64 {
	slice := make([]float64, 0, len(results))

	for _, value := range results {
		slice = append(slice, value)
	}

	return slice
}

func printResults(results map[int]float64) {
	keys := make([]int, 0, len(results))
	for key := range results {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	for _, bucket := range keys {
		bucketResults := results[bucket]
		show.Info("bucket: %v, elapsed: %v\n", bucket, bucketResults)
	}
}
