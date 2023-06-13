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
	"time"

	"clickhouse-benchmark/pkg/clickhouse"
	"clickhouse-benchmark/pkg/show"

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

	// Create a channel to synchronize the completion of batches
	done := make(chan struct{})

	// Create multiple batches based on concurrency limit
	batches := make([]*clickhouse.Batch, writeOpt.concurrencyLimit)
	for i := 0; i < writeOpt.concurrencyLimit; i++ {
		batch, err := clickhouse.Prepare(conn, databaseName, tableName)
		if err != nil {
			return err
		}
		batches[i] = batch

		// Start a goroutine to process each batch
		go func(batch *clickhouse.Batch) {
			defer func() {
				done <- struct{}{} // Signal the completion of the batch
			}()

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
							show.Debug("append is failed: %v", err)
						}
					}
				}
			}

			// Send the batch for execution
			err := sendBatch(batch, debugFlag, debugInfo, totalRecords)
			if err != nil {
				show.Error("Failed to send batch: %v\n", err)
			}
		}(batches[i])
	}

	// Wait for all batches to complete
	for i := 0; i < writeOpt.concurrencyLimit; i++ {
		<-done
	}

	// Perform benchmarking calculations
	elapsedTime := time.Since(startTime)
	completeRequests := totalRecords / writeOpt.size

	// Print benchmarking results
	show.Info("ClickHouse URL: %s", os.Getenv("CLICKHOUSE_URL"))
	show.Info("Benchmarking Bucket Count: %d", writeOpt.bucketCount)
	show.Info("Benchmarking Size: %d", writeOpt.size)
	show.Info("Benchmarking Concurrency", writeOpt.concurrencyLimit)
	show.Info("Benchmarking Bucket Unit: %s", "Seconds")
	show.EmptyLine()

	show.Info("Time taken for tests: %v", elapsedTime)
	show.Info("Complete requests: %d", completeRequests)
	show.Info("Total transferred: %d", totalRecords) // Update this based on the actual transferred data size

	return nil
}

// Send the batch for execution and handle progress output
func sendBatch(batch *clickhouse.Batch, debugFlag bool, debugInfo *DebugAppendMetrics, totalRecords int) error {
	if debugFlag {
		debugInfo.Printf()
		return nil
	}

	err := batch.Send()
	if err != nil {
		return err
	}

	progress := (batch.TotalRows() / writeOpt.size) * 100 / (totalRecords / writeOpt.size)
	show.Info("Progress: %d%%", progress)

	return nil
}
