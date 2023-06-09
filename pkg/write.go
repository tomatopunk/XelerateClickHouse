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
	"sync"
	"time"

	"github.com/spf13/cobra"
)

type WriteOption struct {
	bucketCount      int // bucket count like 30
	size             int // bucket size like 100
	concurrencyLevel int //concurrency level like 1
}

var writeOpt WriteOption

var writeCommand = &cobra.Command{
	Use:  "write",
	Long: ` write some data to clickhouse`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := writeToClickhouse(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	root.AddCommand(writeCommand)
	writeCommand.Flags().IntVar(&writeOpt.bucketCount, "b", 3, "bucket count like 30")
	writeCommand.Flags().IntVar(&writeOpt.size, "n", 1, "bucket size like 100")
	writeCommand.Flags().IntVar(&writeOpt.concurrencyLevel, "c", 1, "concurrency level like 1")
}

func writeToClickhouse() error {
	conn, err := getConn(os.Getenv("CLICKHOUSE_URL"))
	if err != nil {
		return err
	}
	defer conn.Close()

	// Prepare the insert statement
	batch, err := conn.PrepareBatch(context.Background(), fmt.Sprintf("INSERT INTO %s.%s", databaseName, tableName))
	if err != nil {
		return err
	}

	// Calculate the total number of data records
	totalRecords := writeOpt.size * writeOpt.bucketCount
	failedRequests := 0
	startTime := time.Now()

	debugInfo := NewDebugAppendMetrics()

	// Generate and insert data concurrently
	var wg sync.WaitGroup
	wg.Add(writeOpt.concurrencyLevel)

	for i := 0; i < writeOpt.concurrencyLevel; i++ {
		go func() {
			defer wg.Done()

			// Generate data for each bucket
			for bucket := 0; bucket < writeOpt.bucketCount; bucket++ {
				timestamp := startTime.Add(time.Duration(bucket) * time.Second)

				// Generate metrics data
				for j := 0; j < writeOpt.size; j++ {
					//t := timestamp.Add(time.Duration(j) * time.Second)
					t := timestamp
					metric := generateMetric(t)
					err := batch.AppendStruct(&metric)
					if debugFlag {
						debugInfo.Add(metric)
					}

					if err != nil {
						if debugFlag {
							fmt.Printf("append is failed: %v", err)
						}
						failedRequests++
					}
				}
			}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Send the batch for execution
	if debugFlag {
		debugInfo.Printf()
		return nil
	}
	err = batch.Send()
	if err != nil {
		fmt.Printf("Failed to send batch: %v\n", err)
	}

	// Perform benchmarking calculations
	elapsedTime := time.Since(startTime)
	completeRequests := totalRecords / writeOpt.size

	// Print benchmarking results
	fmt.Printf("ClickHouse URL: %s\n", os.Getenv("CLICKHOUSE_URL"))
	fmt.Printf("Benchmarking Bucket Count: %d\n", writeOpt.bucketCount)
	fmt.Printf("Benchmarking Size: %d\n", writeOpt.size)
	fmt.Printf("Benchmarking Bucket Unit: %s\n", "Seconds")
	fmt.Printf("Concurrency Level: %d\n", writeOpt.concurrencyLevel)
	fmt.Printf("\n")

	fmt.Printf("Time taken for tests: %v\n", elapsedTime)
	fmt.Printf("Complete requests: %d\n", completeRequests)
	fmt.Printf("Failed requests: %d\n", failedRequests)
	fmt.Printf("Total transferred: %d\n", totalRecords) // Update this based on the actual transferred data size

	return nil
}
