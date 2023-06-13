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
	"time"

	"clickhouse-benchmark/pkg/show"

	"github.com/spf13/cobra"
)

type WriteOption struct {
	bucketCount int // bucket count like 30
	size        int // bucket size like 100
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
	writeCommand.Flags().IntVar(&writeOpt.bucketCount, "b", 3, "bucket count like 30")
	writeCommand.Flags().IntVar(&writeOpt.size, "n", 1, "bucket size like 100")
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
				failedRequests++
			}
		}
	}

	// Send the batch for execution
	if debugFlag {
		debugInfo.Printf()
		return nil
	}
	err = batch.Send()
	if err != nil {
		show.Error("Failed to send batch: %v\n", err)
	}

	// Perform benchmarking calculations
	elapsedTime := time.Since(startTime)
	completeRequests := totalRecords / writeOpt.size

	// Print benchmarking results
	show.Info("ClickHouse URL: %s\n", os.Getenv("CLICKHOUSE_URL"))
	show.Info("Benchmarking Bucket Count: %d\n", writeOpt.bucketCount)
	show.Info("Benchmarking Size: %d\n", writeOpt.size)
	show.Info("Benchmarking Bucket Unit: %s\n", "Seconds")
	show.Info("\n")

	show.Info("Time taken for tests: %v\n", elapsedTime)
	show.Info("Complete requests: %d\n", completeRequests)
	show.Info("Failed requests: %d\n", failedRequests)
	show.Info("Total transferred: %d\n", totalRecords) // Update this based on the actual transferred data size

	return nil
}
