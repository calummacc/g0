package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/calummacc/g0/internal/printer"
	"github.com/calummacc/g0/internal/runner"
	"github.com/spf13/cobra"
)

var (
	url         string
	concurrency int
	duration    string
	method      string
	body        string
	headers     []string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a load test",
	Long: `Run a load test against a target URL with specified concurrency and duration.

Example:
  g0 run --url https://api.example.com --c 100 --d 10s
  g0 run --url https://api.example.com --c 50 --d 30s --method POST --body '{"key":"value"}' --headers "Content-Type: application/json"`,
	RunE: runLoadTest,
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&url, "url", "u", "", "Target URL (required)")
	runCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 10, "Number of concurrent workers")
	runCmd.Flags().StringVarP(&duration, "duration", "d", "10s", "Test duration (e.g., 10s, 1m, 30s)")
	runCmd.Flags().StringVarP(&method, "method", "m", "GET", "HTTP method")
	runCmd.Flags().StringVarP(&body, "body", "b", "", "Request body")
	runCmd.Flags().StringArrayVarP(&headers, "headers", "H", []string{}, "HTTP headers (can be specified multiple times)")

	runCmd.MarkFlagRequired("url")
}

func runLoadTest(cmd *cobra.Command, args []string) error {
	// Parse duration
	testDuration, err := time.ParseDuration(duration)
	if err != nil {
		return fmt.Errorf("invalid duration format: %w", err)
	}

	// Validate concurrency
	if concurrency <= 0 {
		return fmt.Errorf("concurrency must be greater than 0")
	}

	// Parse headers
	headerMap := make(map[string]string)
	for _, h := range headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid header format: %s (expected 'Key: Value')", h)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		headerMap[key] = value
	}

	// Print test configuration
	printer.PrintTestStart(url, concurrency, testDuration)

	// Create and run the load test
	config := runner.Config{
		URL:         url,
		Concurrency: concurrency,
		Duration:    testDuration,
		Method:      method,
		Body:        body,
		Headers:     headerMap,
	}

	stats, err := runner.Run(config)
	if err != nil {
		return fmt.Errorf("load test failed: %w", err)
	}

	// Print results
	printer.PrintResults(stats)

	return nil
}

