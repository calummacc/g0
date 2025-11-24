package runner

import (
	"context"
	"time"

	"github.com/calummacc/g0/internal/httpclient"
)

// Config holds the configuration for a load test
type Config struct {
	URL         string
	Concurrency int
	Duration    time.Duration
	Method      string
	Body        string
	Headers     map[string]string
}

// Run executes a load test with the given configuration
func Run(config Config) (*Summary, error) {
	// Create HTTP client
	client := httpclient.New()

	// Create request configuration
	request := httpclient.Request{
		Method:  config.Method,
		URL:     config.URL,
		Body:    config.Body,
		Headers: config.Headers,
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), config.Duration)
	defer cancel()

	// Create results channel
	results := make(chan Result, config.Concurrency*10)

	// Create stats collector
	stats := NewStats()

	// Start stats collector goroutine
	statsDone := make(chan struct{})
	go func() {
		defer close(statsDone)
		for {
			select {
			case result, ok := <-results:
				if !ok {
					return
				}
				stats.AddResult(result)
			case <-ctx.Done():
				// Drain remaining results after context is done
				for {
					select {
					case result := <-results:
						stats.AddResult(result)
					default:
						return
					}
				}
			}
		}
	}()

	// Start workers
	for i := 0; i < config.Concurrency; i++ {
		worker := NewWorker(client, request, results)
		go worker.Start(ctx)
	}

	// Wait for duration to complete
	<-ctx.Done()

	// Give workers a moment to stop sending (they check ctx.Done() before sending)
	time.Sleep(50 * time.Millisecond)

	// Close results channel to signal stats collector to finish
	close(results)

	// Wait for stats collector to finish processing
	<-statsDone

	// Finalize stats
	stats.Finalize()

	// Get summary
	summary := stats.GetSummary()

	return &summary, nil
}

