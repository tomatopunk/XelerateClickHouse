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

	"clickhouse-benchmark/pkg/show"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/spf13/cobra"
)

var descCommand = &cobra.Command{
	Use:  "desc",
	Long: ` describe the table `,
	Run: func(cmd *cobra.Command, args []string) {
		if err := descClickhouse(); err != nil {
			show.Error("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	root.AddCommand(descCommand)
}

func descClickhouse() error {
	conn, err := getConn(os.Getenv("CLICKHOUSE_URL"))
	if err != nil {
		return err
	}
	defer conn.Close()

	// Get create table SQL
	createTableQuery, err := getCreateTableSQL(conn)
	if err != nil {
		return err
	}

	// Get partition information
	partitions, err := getPartitionsInfo(conn)
	if err != nil {
		return err
	}

	// Print the description
	show.Info("Description <%s>\n\n", tableName)
	show.Info("Clickhouse URL: %s\n\n", os.Getenv("CLICKHOUSE_URL"))
	show.Info("Create Table SQL:\n\n%s\n\n", createTableQuery)
	show.Info("Partition: \n\n")
	for i, partition := range partitions {
		if i > partitionLimit {
			show.Info("...")
			break
		}
		show.Info("Partition %s, disk: %s, total_row: %d, all_disk: %d\n", partition.Name, partition.DiskName, partition.RowCount, partition.DiskSize)
	}

	printPartitionAggregation(partitions)

	return nil
}

type PartitionInfo struct {
	Name     string
	DiskName string
	RowCount int
	DiskSize int
}

func printPartitionAggregation(partitions []PartitionInfo) {
	// Map to aggregate partitions by disk
	aggregatedPartitions := make(map[string]*PartitionInfo)

	// Aggregate partitions by disk
	for _, partition := range partitions {
		// Check if the disk already exists in the map
		if _, ok := aggregatedPartitions[partition.DiskName]; ok {
			// If the disk exists, update the aggregated values
			aggregatedPartitions[partition.DiskName].RowCount += partition.RowCount
			aggregatedPartitions[partition.DiskName].DiskSize += partition.DiskSize
		} else {
			// If the disk doesn't exist, add it to the map
			aggregatedPartitions[partition.DiskName] = &PartitionInfo{
				DiskName: partition.DiskName,
				RowCount: partition.RowCount,
				DiskSize: partition.DiskSize,
			}
		}
	}

	// Print the aggregated partition information
	for disk, partition := range aggregatedPartitions {
		show.Info("Disk: %s\n", disk)
		show.Info("Partition Count: %d\n", len(partitions))
		show.Info("Sum RowCount: %d\n", partition.RowCount)
		show.Info("Sum DiskSize: %.2f MB\n\n", float64(partition.DiskSize)/1024/1024)
	}
}

func getPartitionsInfo(conn driver.Conn) ([]PartitionInfo, error) {
	query := fmt.Sprintf("SELECT partition, disk_name, sum(rows) AS total_row, sum(bytes_on_disk) AS all_disk FROM system.parts WHERE active AND database = '%s' AND partition != '19700101' AND table = '%s' GROUP BY partition, disk_name ORDER BY partition", databaseName, tableName)
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	partitions := make([]PartitionInfo, 0)
	var partition, diskName string
	var totalRow, allDisk int
	for rows.Next() {
		err := rows.Scan(&partition, &diskName, &totalRow, &allDisk)
		if err != nil {
			return nil, err
		}
		partitions = append(partitions, PartitionInfo{
			Name:     partition,
			DiskName: diskName,
			RowCount: totalRow,
			DiskSize: allDisk,
		})
	}
	return partitions, nil
}

func getCreateTableSQL(conn driver.Conn) (string, error) {
	query := fmt.Sprintf("SELECT name AS table, create_table_query FROM system.tables WHERE database = '%s' AND (engine = 'ReplicatedMergeTree' OR engine = 'ReplicatedReplacingMergeTree')", databaseName)
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var tableName, createTableQuery string
	for rows.Next() {
		err := rows.Scan(&tableName, &createTableQuery)
		if err != nil {
			return "", err
		}
	}

	return createTableQuery, nil
}
