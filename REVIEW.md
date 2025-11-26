# PR Review: Load testing runner

## Summary
- Reviewed CLI runner and HTTP client for request execution and progress reporting.
- Attempted to run `go test ./...`, but the command hung (likely waiting on external module downloads) and was interrupted.

## Blocking Issues
1. **Load test duration is not enforced for in-flight requests.** Workers rely on `context.Context` to stop looping, but individual HTTP requests are issued without context cancellation. If a request stalls longer than the configured duration (e.g., a slow server or hanging connection), `worker.Start` waits for `client.Do` to return before checking `ctx.Done()`, meaning the test can overrun the requested duration by up to the HTTP client's timeout (30s). Use `http.NewRequestWithContext` (or add context to `httpclient.Request`) so in-flight requests are canceled when the test duration elapses. 【F:internal/runner/worker.go†L25-L56】【F:internal/httpclient/client.go†L47-L88】

## Tests
- `go test ./...` (hangs; interrupted)【eb9fa5†L1】【63a850†L1-L2】
