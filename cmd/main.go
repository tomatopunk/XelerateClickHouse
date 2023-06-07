package main

import (
	"log"

	"clickhouse-benchmark/pkg"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	pkg.Execute()
}
