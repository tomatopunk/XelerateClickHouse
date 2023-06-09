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
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/spf13/cobra"
)

type readOption struct {
	bucketCount      int // bucket count like 30
	size             int // bucket size like 100
	concurrencyLevel int //concurrency level like 1
}

var readOpt readOption

var readCommand = &cobra.Command{
	Use:  "read",
	Long: ` benchmarking read `,
	Run: func(cmd *cobra.Command, args []string) {
		if err := benchmarkReadQueries(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	root.AddCommand(readCommand)

	readCommand.Flags().IntVar(&readOpt.bucketCount, "b", 1, "bucket count like 30")
	readCommand.Flags().IntVar(&readOpt.size, "n", 10, "bucket size like 100")
	readCommand.Flags().IntVar(&readOpt.concurrencyLevel, "c", 1, "concurrency level like 1")

}

func benchmarkReadQueries() error {
	db, err := sql.Open("clickhouse", os.Getenv("CLICKHOUSE_URL"))
	if err != nil {
		return err
	}
	defer db.Close()

	// Execute read queries concurrently and measure their execution time
	results := make([][]float64, readOpt.bucketCount)
	var wg sync.WaitGroup
	wg.Add(readOpt.bucketCount)

	failedRequests := 0

	totalStart := time.Now()

	for i := 0; i < readOpt.bucketCount; i++ {
		go func(bucketIndex int) {
			defer wg.Done()

			bucketResults := make([]float64, readOpt.size)

			for j := 0; j < readOpt.size; j++ {
				start := time.Now().Add(time.Duration(bucketIndex) * time.Duration(readOpt.size) * time.Minute)
				end := start.Add(time.Duration(readOpt.size) * time.Minute)

				query := fmt.Sprintf("SELECT * FROM %s.%s WHERE timestamp >= '%s' AND timestamp < '%s' LIMIT 100", databaseName, tableName, start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
				startTime := time.Now()
				_, err := db.Exec(query)
				elapsedTime := time.Since(startTime)

				if err != nil {
					log.Printf("Error executing query: %v", err)
					failedRequests++
				}

				bucketResults[j] = elapsedTime.Seconds()
			}

			results[bucketIndex] = bucketResults
		}(i)
	}

	// Wait for all queries to complete
	wg.Wait()

	totalTime := time.Since(totalStart)

	for i, bucketResults := range results {
		p50, _ := stats.Percentile(bucketResults, 50)
		p80, _ := stats.Percentile(bucketResults, 80)
		p99, _ := stats.Percentile(bucketResults, 99)
		p999, _ := stats.Percentile(bucketResults, 99.9)

		start := time.Now().Add(time.Duration(i) * time.Duration(readOpt.size) * time.Minute)
		end := start.Add(time.Duration(readOpt.size) * time.Minute)

		fmt.Printf("Start %d: %s, End %d: %s, p50: %v,p80: %v, p99: %v, p999: %v\n", i+1, start, i+1, end, p50, p80, p99, p999)
	}

	// Print benchmarking results
	fmt.Printf("\n\n")
	fmt.Printf("ClickHouse URL: %s\n", os.Getenv("CLICKHOUSE_URL"))
	fmt.Printf("Concurrency Level: %d\n", readOpt.bucketCount)
	fmt.Printf("Total queries executed: %d\n", readOpt.bucketCount*readOpt.size)
	fmt.Printf("Failed requests: %d\n", failedRequests)
	fmt.Printf("Time taken for tests: %v\n", totalTime)
	fmt.Printf("Total transferred: N/A\n") // Calculate and print if needed

	return nil
}
