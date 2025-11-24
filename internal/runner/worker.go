package runner

import (
	"context"

	"github.com/calummacc/g0/internal/httpclient"
)

// Worker sends HTTP requests in a loop until the context is cancelled
type Worker struct {
	client  *httpclient.Client
	request httpclient.Request
	results chan<- Result
}

// NewWorker creates a new worker
func NewWorker(client *httpclient.Client, request httpclient.Request, results chan<- Result) *Worker {
	return &Worker{
		client:  client,
		request: request,
		results: results,
	}
}

// Start begins the worker loop, sending requests until ctx is cancelled
func (w *Worker) Start(ctx context.Context) {
	defer func() {
		// Recover from any panic (e.g., sending on closed channel)
		// This should not happen with proper synchronization, but provides safety
		recover()
	}()

	for {
		// Check if context is done before starting a new request
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Send request
		resp := w.client.Do(w.request)

		// Check context again before sending result (request might have taken time)
		select {
		case <-ctx.Done():
			// Context cancelled, don't send result
			return
		case w.results <- Result{
			Latency:    resp.Latency,
			StatusCode: resp.StatusCode,
			Error:      resp.Error,
		}:
			// Successfully sent result, continue loop
		}
	}
}

