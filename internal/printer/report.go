package printer

import (
	"fmt"
	"time"

	"github.com/calummacc/g0/internal/runner"
)

// PrintTestStart prints the test configuration
func PrintTestStart(url string, concurrency int, duration time.Duration) {
	fmt.Println("Load Test Started")
	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Concurrency: %d\n", concurrency)
	fmt.Printf("Duration: %s\n", duration)
	fmt.Println()
}

// PrintResults prints the test results in a formatted way
func PrintResults(summary *runner.Summary) {
	fmt.Println("Results:")
	fmt.Printf("Total Requests: %d\n", summary.TotalRequests)
	fmt.Printf("Success: %d\n", summary.SuccessRequests)
	fmt.Printf("Failed: %d\n", summary.FailedRequests)
	fmt.Printf("RPS: %.1f\n", summary.RPS)
	fmt.Println()

	fmt.Println("Latency:")
	fmt.Printf("  Min: %s\n", formatDuration(summary.MinLatency))
	fmt.Printf("  Avg: %s\n", formatDuration(summary.AvgLatency))
	fmt.Printf("  Max: %s\n", formatDuration(summary.MaxLatency))
	fmt.Printf("  p90: %s\n", formatDuration(summary.P90Latency))
	fmt.Printf("  p95: %s\n", formatDuration(summary.P95Latency))
	fmt.Printf("  p99: %s\n", formatDuration(summary.P99Latency))

	// Print status code distribution if there are any
	if len(summary.StatusCodeCounts) > 0 {
		fmt.Println()
		fmt.Println("Status Codes:")
		for code, count := range summary.StatusCodeCounts {
			fmt.Printf("  %d: %d\n", code, count)
		}
	}
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%.2fns", float64(d.Nanoseconds()))
	} else if d < time.Millisecond {
		return fmt.Sprintf("%.2fµs", float64(d.Nanoseconds())/1000.0)
	} else if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1000000.0)
	}
	return d.Round(time.Millisecond).String()
}

