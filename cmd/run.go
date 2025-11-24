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

	// Print logo
	printer.PrintLogo()

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

	// Channel to receive test result
	resultChan := make(chan *runner.RunResult, 1)
	errChan := make(chan error, 1)
	statsChan := make(chan *runner.Stats, 1)

	// Start progress monitoring in a goroutine
	progressDone := make(chan struct{})
	startTime := time.Now()
	var stats *runner.Stats

	// Start the test in a goroutine
	go func() {
		result, err := runner.RunWithStatsAndChannel(config, statsChan)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- result
	}()

	// Progress monitoring goroutine
	go func() {
		// Wait for stats to be available
		select {
		case s := <-statsChan:
			stats = s
		case <-time.After(2 * time.Second):
			// Stats not available yet, continue anyway (shouldn't happen normally)
		}

		ticker := time.NewTicker(100 * time.Millisecond) // Update every 100ms
		defer ticker.Stop()

		for {
			select {
			case s := <-statsChan:
				// Stats instance is now available (if not received earlier)
				stats = s
			case <-ticker.C:
				elapsed := time.Since(startTime)
				if elapsed >= testDuration {
					// Test duration reached, stop updating
					return
				}
				if stats != nil {
					progressStats := stats.GetProgressStats()
					printer.PrintProgress(elapsed, testDuration, &progressStats)
				} else {
					// Stats not available yet, show basic progress with zero stats
					zeroStats := runner.ProgressStats{}
					printer.PrintProgress(elapsed, testDuration, &zeroStats)
				}
			case <-progressDone:
				return
			}
		}
	}()

	// Wait for test to complete
	var result *runner.RunResult
	select {
	case err := <-errChan:
		close(progressDone)
		time.Sleep(50 * time.Millisecond)
		printer.ClearProgress()
		return fmt.Errorf("load test failed: %w", err)
	case result = <-resultChan:
		// Test completed
	}

	// Stop progress updates
	close(progressDone)
	time.Sleep(150 * time.Millisecond) // Give progress goroutine time to finish and print final update

	// Clear progress line
	printer.ClearProgress()
	fmt.Println() // Add a newline after clearing progress

	// Print results
	printer.PrintResults(result.Summary)

	return nil
}

