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
	"io/ioutil"
	"os"
	"strings"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/spf13/cobra"
)

var initCommand = &cobra.Command{
	Use:  "init",
	Long: `create database, create tables `,
	Run: func(cmd *cobra.Command, args []string) {
		if err := initClickhouse(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	root.AddCommand(initCommand)
}

func initClickhouse() error {
	conn, err := getConn(os.Getenv("CLICKHOUSE_URL"))
	if err != nil {
		return fmt.Errorf("failed to connect to ClickHouse: %v", err)
	}
	defer conn.Close()

	// Create the database
	if err := executeSQLFile(conn, "scripts/test_database.sql"); err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}

	// Create the tables
	if err := executeSQLFile(conn, "scripts/test_table.sql"); err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	fmt.Println("Database and tables created successfully")

	return nil
}

func executeSQLFile(conn driver.Conn, filePath string) error {
	// Read the content of the SQL file
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %v", err)
	}

	// Split the file content into individual SQL statements
	statements := strings.Split(string(fileContent), ";")

	// Execute each SQL statement
	for _, statement := range statements {
		trimmedStatement := strings.TrimSpace(statement)
		if trimmedStatement == "" {
			continue // Skip empty statements
		}

		err := conn.Exec(context.Background(), trimmedStatement)
		if err != nil {
			return fmt.Errorf("failed to execute SQL statement: %v ,sql: %v", err, trimmedStatement)
		}
	}

	return nil
}
