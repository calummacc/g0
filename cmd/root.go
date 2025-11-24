package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "g0",
	Short: "g0 - A minimal high-performance HTTP load tester",
	Long: `g0 is a fast, lightweight CLI tool that sends concurrent HTTP requests
and measures load-testing metrics. It's designed to be simple yet powerful.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Here you can define flags and configuration settings
}

