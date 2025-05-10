package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "eco"}

	var cmdRun = &cobra.Command{
		Use:   "run",
		Short: "Run a system",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("System running...")
		},
	}

	rootCmd.AddCommand(cmdRun)
	rootCmd.Execute()
}
