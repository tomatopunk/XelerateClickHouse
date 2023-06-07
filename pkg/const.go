package pkg

import (
	"os"
)

const (
	tableName      = "metrics"
	databaseName   = "test"
	partitionLimit = 10
)

var clickhouseURL = os.Getenv("CLICKHOUSE_URL")
