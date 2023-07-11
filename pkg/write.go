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
	"os"
	"sync"
	"time"

	"clickhouse-benchmark/pkg/clickhouse"
	"clickhouse-benchmark/pkg/show"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

type WriteOption struct {
	bucketCount      int // bucket count like 30
	size             int // bucket size like 100
	concurrencyLimit int
}

var writeOpt WriteOption

var writeCommand = &cobra.Command{
	Use:  "write",
	Long: ` write some data to clickhouse`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := writeToClickhouse(); err != nil {
			show.Error("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	root.AddCommand(writeCommand)
	writeCommand.Flags().IntVar(&writeOpt.bucketCount, "b", 100, "bucket count like 30")
	writeCommand.Flags().IntVar(&writeOpt.size, "n", 1, "bucket size like 100")
	writeCommand.Flags().IntVar(&writeOpt.concurrencyLimit, "c", 1, "concurrency limit like 1")
}

func writeToClickhouse() error {
	conn, err := getConn(os.Getenv("CLICKHOUSE_URL"))
	if err != nil {
		return err
	}
	defer conn.Close()

	startTime := time.Now()
	// Calculate the total number of data records
	totalRecords := writeOpt.size * writeOpt.bucketCount

	debugInfo := NewDebugAppendMetrics()

	wg := sync.WaitGroup{}
	wg.Add(writeOpt.concurrencyLimit)

	bar := pb.StartNew(totalRecords)

	for i := 1; i <= writeOpt.concurrencyLimit; i++ {
		batch, err := clickhouse.Prepare(conn, databaseName, tableName)
		if err != nil {
			return err
		}

		// Start a goroutine to process each batch
		go func(step int) {
			defer func() {
				wg.Done()
			}()

			// Generate data for each bucket
			for bucket := step; bucket <= writeOpt.bucketCount*writeOpt.concurrencyLimit; bucket += writeOpt.concurrencyLimit {
				if bucket > writeOpt.bucketCount {
					break
				}
				//step concurrency
				timestamp := startTime.Add(time.Duration(bucket) * time.Second)

				// Generate metrics data
				for j := 0; j < writeOpt.size; j++ {
					//t := timestamp.Add(time.Duration(j) * time.Second)
					t := timestamp
					metric := generateMetric(t)
					err := batch.AppendStruct(&metric)
					bar.Increment()
					if debugFlag {
						debugInfo.Lock()
						debugInfo.Add(metric)
						debugInfo.Unlock()
					}

					if err != nil {
						show.Error("append is failed: %v", err)
					}
				}
			}

			// Send the batch for execution
			if debugFlag {
				return
			}

			err := batch.Send()
			if err != nil {
				show.Error("Failed to send batch: %v\n", err)
				return
			}
		}(i)
	}

	// Wait for all batches to complete
	wg.Wait()

	bar.Finish()
	if debugFlag {
		debugInfo.Printf()
	}

	// Perform benchmarking calculations
	elapsedTime := time.Since(startTime)
	completeRequests := totalRecords / writeOpt.size

	// Print benchmarking results
	show.Info("ClickHouse URL: %s", os.Getenv("CLICKHOUSE_URL"))
	show.Info("Benchmarking Bucket Count: %d", writeOpt.bucketCount)
	show.Info("Benchmarking Size: %d", writeOpt.size)
	show.Info("Benchmarking Concurrency: %v", writeOpt.concurrencyLimit)
	show.Info("Benchmarking Bucket Unit: %s", "Seconds")
	show.EmptyLine()

	show.Info("Time taken for tests: %v", elapsedTime)
	show.Info("Complete requests: %d", completeRequests)
	show.Info("Total transferred: %d", totalRecords) // Update this based on the actual transferred data size

	return nil
}
