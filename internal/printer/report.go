package printer

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/calummacc/g0/internal/runner"
)

// PrintLogo prints the g0 logo
func PrintLogo() {
	logo := `
	\033[36m    ______      \033[33m__ 
	\033[36m   / ____/___ _ \033[33m/ /____  ____ 
	\033[36m  / / __/ __  /\033[33m / ___/ / __ \
	\033[36m / /_/ / /_/ /\033[33m (__  ) / /_/ /
	\033[36m \____/\__,_/\033[33m /____(_)\____/ 
	\033[0m
	\033[32m            g0 — High-Performance Load Tester\033[0m
	`
	fmt.Print(logo)
	fmt.Println()
}

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

// PrintProgress displays a progress bar with current test statistics
func PrintProgress(elapsed time.Duration, totalDuration time.Duration, stats *runner.ProgressStats) {
	// Calculate progress percentage
	progress := float64(elapsed) / float64(totalDuration)
	if progress > 1.0 {
		progress = 1.0
	}

	// Create progress bar (50 characters wide)
	barWidth := 50
	filled := int(progress * float64(barWidth))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	// Calculate current RPS
	var rps float64
	if elapsed > 0 {
		rps = float64(stats.TotalRequests) / elapsed.Seconds()
	}

	// Format elapsed and remaining time
	elapsedStr := formatDurationShort(elapsed)
	remaining := totalDuration - elapsed
	if remaining < 0 {
		remaining = 0
	}
	remainingStr := formatDurationShort(remaining)

	// Clear previous line and print progress
	fmt.Fprintf(os.Stderr, "\r[%s] %.1f%% | Elapsed: %s | Remaining: %s | Requests: %d | Success: %d | Failed: %d | RPS: %.1f",
		bar, progress*100, elapsedStr, remainingStr, stats.TotalRequests, stats.SuccessRequests, stats.FailedRequests, rps)
}

// ClearProgress clears the progress line
func ClearProgress() {
	fmt.Fprintf(os.Stderr, "\r%s\r", strings.Repeat(" ", 150))
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

// formatDurationShort formats a duration in a short, readable way for progress display
func formatDurationShort(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.0fms", float64(d.Nanoseconds())/1000000.0)
	} else if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%dm%ds", minutes, seconds)
}
