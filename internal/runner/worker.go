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
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Send request
			resp := w.client.Do(w.request)

			// Send result to channel
			result := Result{
				Latency:    resp.Latency,
				StatusCode: resp.StatusCode,
				Error:      resp.Error,
			}

			select {
			case w.results <- result:
			case <-ctx.Done():
				return
			}
		}
	}
}

