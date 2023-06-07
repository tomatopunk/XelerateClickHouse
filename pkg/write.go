package pkg

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

type WriteOption struct {
	bucketUnit       string //bucket unit like s
	bucketCount      int    // bucket count like 30
	size             int    // bucket size like 100
	concurrencyLevel int    //concurrency level like 1
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
	writeCommand.Flags().StringVar(&writeOpt.bucketUnit, "u", "day", "bucket unit like s")
	writeCommand.Flags().IntVar(&writeOpt.bucketCount, "b", 1, "bucket count like 30")
	writeCommand.Flags().IntVar(&writeOpt.size, "n", 10, "bucket size like 100")
	writeCommand.Flags().IntVar(&writeOpt.concurrencyLevel, "c", 1, "concurrency level like 1")
}

func writeToClickhouse() error {
	db, err := sql.Open("clickhouse", clickhouseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	// Prepare the insert statement
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s.%s (timestamp, metric_group, number_field_keys, number_field_values, string_field_keys, string_field_values, tag_keys, tag_values) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", databaseName, tableName))
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Calculate the total number of data records
	totalRecords := writeOpt.size * writeOpt.bucketCount
	failedRequests := 0
	startTime := time.Now()

	// Generate and insert data concurrently
	var wg sync.WaitGroup
	wg.Add(writeOpt.concurrencyLevel)
	for i := 0; i < writeOpt.concurrencyLevel; i++ {
		go func() {
			defer wg.Done()

			// Generate and insert data for each time bucket
			for j := 0; j < writeOpt.bucketCount; j++ {
				// Generate data for the current time bucket
				bucketStartTime := time.Now().Add(-time.Duration(j) * getBucketDuration(writeOpt.bucketUnit))
				data := generateData(writeOpt.size, bucketStartTime)

				// Batch insert the generated data into ClickHouse
				if err := batchInsertData(db, stmt, data); err != nil {
					failedRequests++
					fmt.Printf("Error inserting data: %v\n", err)
				}
			}
		}()
	}

	wg.Wait()

	// Perform benchmarking calculations
	elapsedTime := time.Since(startTime)
	completeRequests := totalRecords / writeOpt.size

	// Print benchmarking results
	fmt.Printf("ClickHouse URL: %s\n", clickhouseURL)
	fmt.Printf("Benckmarking Bucket Count : %d\n", writeOpt.bucketCount)
	fmt.Printf("Benckmarking Size : %d\n", writeOpt.size)
	fmt.Printf("Benckmarking Bucket Unit : %s\n", writeOpt.bucketUnit)
	fmt.Printf("Concurrency Level: %d\n", writeOpt.concurrencyLevel)
	fmt.Printf("\n\n")

	fmt.Printf("Time taken for tests: %v\n", elapsedTime)
	fmt.Printf("Complete requests: %d\n", completeRequests)
	fmt.Printf("Failed requests: %d\n", failedRequests)
	fmt.Printf("Total transferred: %d\n", totalRecords) // Update this based on the actual transferred data size

	return nil
}

// Perform batch insert of data into ClickHouse
func batchInsertData(db *sql.DB, stmt *sql.Stmt, data [][]interface{}) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, row := range data {
		_, err := stmt.Exec(row...)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Generate data for a time bucket
func generateData(size int, bucketStartTime time.Time) [][]interface{} {
	data := make([][]interface{}, 0, size)

	for i := 0; i < size; i++ {
		// Generate data for each record in the bucket
		timestamp := bucketStartTime.Add(time.Duration(i) * time.Minute) // Update this based on the desired time granularity
		metricGroup := "group" + strconv.Itoa(rand.Intn(10))             // Example metric group generation

		// Generate number field keys and values
		numberFieldKeys := make([]string, 0)
		numberFieldValues := make([]float64, 0)
		for j := 0; j < 3; j++ {
			key := fmt.Sprintf("field%d", j+1)
			numberFieldKeys = append(numberFieldKeys, key)
			numberFieldValues = append(numberFieldValues, rand.Float64()) // Example number field value generation
		}

		// Generate string field keys and values
		stringFieldKeys := make([]string, 0)
		stringFieldValues := make([]string, 0)
		for j := 0; j < 3; j++ {
			key := fmt.Sprintf("field%d", j+4)
			stringFieldKeys = append(stringFieldKeys, key)
			stringFieldValues = append(stringFieldValues, "value"+strconv.Itoa(rand.Intn(3)+1)) // Example string field value generation
		}

		// Generate tag keys and values
		tagKeys := make([]string, 0)
		tagValues := make([]string, 0)
		for j := 0; j < 3; j++ {
			key := fmt.Sprintf("tag%d", j+1)
			tagKeys = append(tagKeys, key)
			tagValues = append(tagValues, "value"+strconv.Itoa(rand.Intn(3)+1)) // Example tag value generation
		}

		data = append(data, []interface{}{
			timestamp,
			metricGroup,
			numberFieldKeys,
			numberFieldValues,
			stringFieldKeys,
			stringFieldValues,
			tagKeys,
			tagValues,
		})
	}

	return data
}

// Get the duration for a single time bucket based on the bucket unit
func getBucketDuration(bucketUnit string) time.Duration {
	switch bucketUnit {
	case "day":
		return 24 * time.Hour
	case "hour":
		return time.Hour
	case "minute":
		return time.Minute
	default:
		return 24 * time.Hour // Default to day
	}
}
