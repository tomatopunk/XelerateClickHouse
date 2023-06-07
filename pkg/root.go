package pkg

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "cb",
	Short: "A CLI application to create database and tables in ClickHouse",
	Long:  ``,
}

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
