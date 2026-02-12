# Architecture

## Overview

`storm-data-shared` is a Go library module (`github.com/couchcryptid/storm-data-shared`) that provides common infrastructure code for the storm data pipeline's Go services. It is imported as a regular Go module dependency by the [ETL](https://github.com/couchcryptid/storm-data-etl) and [API](https://github.com/couchcryptid/storm-data-api) services.

This is a library, not a service -- it has no `main` package, no Docker image, and no runtime deployment.

## Package Layout

```
config/
  env.go              Environment variable helpers and shared parsers
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

## Package Details

### `config`

Provides `EnvOrDefault` for reading environment variables with fallbacks, `ParseBrokers` for splitting comma-separated Kafka broker lists, and typed parsers for shared settings (`BATCH_SIZE`, `BATCH_FLUSH_INTERVAL`, `SHUTDOWN_TIMEOUT`).

Each service's `internal/config` package imports these helpers and combines them with service-specific settings. For example, the ETL's `config.Load()` calls `sharedcfg.ParseShutdownTimeout()` for the shared timeout parser, then adds Mapbox-specific configuration on top.

### `observability`

#### Logging

`NewLogger(level, format)` creates a structured `slog.Logger` and sets it as the process default via `slog.SetDefault()`. Supports `debug`, `info`, `warn`, and `error` levels. Format is `json` (default) or `text` for human-readable output.

Each service wraps this in a thin `observability.NewLogger(cfg)` function that extracts `LogLevel` and `LogFormat` from the service's config struct.

#### Health Endpoints

`LivenessHandler()` returns 200 `{"status": "healthy"}` unconditionally. `ReadinessHandler(checker)` calls the `ReadinessChecker` interface and returns 200 or 503 based on the result.

The `ReadinessChecker` interface is the key contract:

```go
type ReadinessChecker interface {
    CheckReadiness(ctx context.Context) error
}
```

Each service provides its own implementation:

| Service | Implementation | Ready When |
|---------|---------------|------------|
| [ETL](https://github.com/couchcryptid/storm-data-etl) | `pipeline.Pipeline` (atomic bool) | At least one message processed |
| [API](https://github.com/couchcryptid/storm-data-api) | `database.PoolReadiness` | PostgreSQL pool responds to ping |

Go's structural typing means services implement this interface without importing the shared module in their readiness checker code -- they only need a method with the matching signature.

#### JSON Response Helper

`WriteJSON(w, status, v)` sets `Content-Type: application/json`, writes the status code, and encodes the value. Used by both health handlers and available for other JSON responses.

### `retry`

`NextBackoff(current, max)` implements capped exponential backoff by doubling the current duration up to a maximum. `SleepWithContext(ctx, d)` sleeps for the specified duration but returns early (with `false`) if the context is cancelled.

Used by the ETL pipeline's extract-transform-load loop to back off on Kafka broker failures.

## Design Decisions

### Library, Not Framework

The shared module provides standalone functions and interfaces. It does not impose structure on consuming services. Each service's `internal/` packages remain in control of how shared code is composed.

**Why**: Services have different architectures (hexagonal ETL vs layered API). A framework would force them into a common shape that doesn't fit both.

### Thin Service Wrappers

Services wrap shared functions in thin adapters (e.g., `observability.NewLogger(cfg)` calls `sharedobs.NewLogger(cfg.LogLevel, cfg.LogFormat)`). The shared library takes primitive arguments, not service-specific config structs.

**Why**: Keeps the shared module free of service-specific types. Services can change their config structure without affecting the library.

### Interface-Based Health Checks

The `ReadinessChecker` interface is defined in the shared module and implemented by each service's infrastructure layer.

**Why**: Provides a consistent readiness pattern (same HTTP response shape, same timeout behavior) while letting each service define its own readiness criteria. Go's structural typing means implementors don't need to import the shared module.

### No Prometheus Metrics in Shared

Each service defines its own Prometheus metrics in `internal/observability/metrics.go`. The shared module does not provide a common metrics struct.

**Why**: Metric definitions are service-specific (HTTP metrics, Kafka metrics, database metrics, geocoding metrics). A shared metrics struct would either be too generic to be useful or would need to know about every service's concerns.

## Related

- [System Architecture](https://github.com/couchcryptid/storm-data-system/wiki/Architecture) -- full pipeline design and improvement roadmap
- [ETL Architecture](https://github.com/couchcryptid/storm-data-etl/wiki/Architecture) -- hexagonal design that imports shared config, observability, and retry
- [API Architecture](https://github.com/couchcryptid/storm-data-api/wiki/Architecture) -- layered design that imports shared config and observability
- [System Observability](https://github.com/couchcryptid/storm-data-system/wiki/Observability) -- health checks, metrics, and logging across all services
- [[Configuration]] -- shared environment variable parsers
- [[Development]] -- testing, linting, and versioning
