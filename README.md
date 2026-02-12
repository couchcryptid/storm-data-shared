# Storm Data Shared

Shared Go library for the storm data pipeline. Provides common configuration helpers, observability primitives (logging, health endpoints), and retry utilities used by both the [ETL](https://github.com/couchcryptid/storm-data-etl) and [API](https://github.com/couchcryptid/storm-data-api) services.

## Packages

### `config`

Environment variable helpers and parsers for settings shared across services.

| Function | Description |
|----------|-------------|
| `EnvOrDefault(key, fallback)` | Read env var with fallback |
| `ParseBrokers(value)` | Split comma-separated broker list |
| `ParseBatchSize()` | Parse `BATCH_SIZE` (default 50, range 1--1000) |
| `ParseBatchFlushInterval()` | Parse `BATCH_FLUSH_INTERVAL` (default 500ms) |
| `ParseShutdownTimeout()` | Parse `SHUTDOWN_TIMEOUT` (default 10s) |

### `observability`

Structured logging and standard health endpoints.

| Function | Description |
|----------|-------------|
| `NewLogger(level, format)` | Create `slog.Logger` (JSON or text) and set as default |
| `LivenessHandler()` | `GET /healthz` -- always 200 `{"status": "healthy"}` |
| `ReadinessHandler(checker)` | `GET /readyz` -- 200 or 503 based on `ReadinessChecker` |
| `WriteJSON(w, status, v)` | JSON response helper |

The `ReadinessChecker` interface:

```go
type ReadinessChecker interface {
    CheckReadiness(ctx context.Context) error
}
```

### `retry`

Backoff and context-aware sleep utilities.

| Function | Description |
|----------|-------------|
| `NextBackoff(current, max)` | Exponential backoff with cap |
| `SleepWithContext(ctx, d)` | Context-cancellable sleep |

## Usage

```go
import (
    "github.com/couchcryptid/storm-data-shared/config"
    "github.com/couchcryptid/storm-data-shared/observability"
    "github.com/couchcryptid/storm-data-shared/retry"
)

// Configuration
brokers := config.ParseBrokers(config.EnvOrDefault("KAFKA_BROKERS", "localhost:9092"))

// Logging
logger := observability.NewLogger("info", "json")

// Health endpoints
http.Handle("/healthz", observability.LivenessHandler())
http.Handle("/readyz", observability.ReadinessHandler(myChecker))

// Retry
backoff = retry.NextBackoff(backoff, 5*time.Second)
retry.SleepWithContext(ctx, backoff)
```

## Development

```
make test         # Run unit tests with race detector
make test-cover   # Generate and open HTML coverage report
make lint         # Run golangci-lint
make fmt          # Format with gofmt + goimports
make vuln         # Run govulncheck
make clean        # Remove coverage artifacts
```

## Project Structure

```
config/
  env.go              Environment variable helpers and parsers
  env_test.go         Tests for all parsers and edge cases
observability/
  logging.go          Structured slog logger factory
  logging_test.go     Tests for level parsing and logger creation
  health.go           Liveness, readiness handlers and ReadinessChecker interface
  health_test.go      Tests for health endpoint responses
retry/
  backoff.go          Exponential backoff and context-aware sleep
  backoff_test.go     Tests for backoff calculation and cancellation
```

## Documentation

See the [project wiki](../../wiki) for detailed documentation:

- [Architecture](../../wiki/Architecture) -- Package design, interface contracts, and design decisions
- [Configuration](../../wiki/Configuration) -- Shared environment variables and parsing rules
- [Development](../../wiki/Development) -- Build, test, lint, CI, and conventions
