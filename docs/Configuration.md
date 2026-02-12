# Configuration

The shared library provides parsers for environment variables that are common across the [ETL](https://github.com/couchcryptid/storm-data-etl) and [API](https://github.com/couchcryptid/storm-data-api) services. Each service's `internal/config` package uses these parsers alongside service-specific configuration.

## Shared Environment Variables

These variables are parsed by the shared `config` package and used by both Go services:

| Variable | Default | Validation | Parser |
|----------|---------|-----------|--------|
| `BATCH_SIZE` | `50` | Integer, 1--1000 | `ParseBatchSize()` |
| `BATCH_FLUSH_INTERVAL` | `500ms` | Positive Go duration | `ParseBatchFlushInterval()` |
| `SHUTDOWN_TIMEOUT` | `10s` | Positive Go duration | `ParseShutdownTimeout()` |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warn`, `error` | Used by `observability.NewLogger()` |
| `LOG_FORMAT` | `json` | `json` or `text` | Used by `observability.NewLogger()` |

`KAFKA_BROKERS` is parsed by `ParseBrokers()` but the default value is service-specific (`localhost:9092` for ETL, `localhost:29092` for API).

## Helpers

### `EnvOrDefault(key, fallback)`

Reads an environment variable. Returns the fallback if the variable is unset or empty.

```go
port := config.EnvOrDefault("PORT", "8080")
```

### `ParseBrokers(value)`

Splits a comma-separated string into a slice of broker addresses, trimming whitespace:

```go
brokers := config.ParseBrokers("broker1:9092, broker2:9092")
// ["broker1:9092", "broker2:9092"]
```

## Error Handling

All `Parse*` functions return `(value, error)`. On invalid input, the error message includes the variable name:

```
invalid BATCH_SIZE: must be 1-1000
invalid SHUTDOWN_TIMEOUT
invalid BATCH_FLUSH_INTERVAL: must be a positive duration
```

Services call these in their `config.Load()` function and propagate errors to `main`, which logs and exits.

## How Services Use It

Each service has its own `internal/config` package that combines shared parsers with service-specific settings:

```go
import sharedcfg "github.com/couchcryptid/storm-data-shared/config"

func Load() (*Config, error) {
    shutdownTimeout, err := sharedcfg.ParseShutdownTimeout()
    if err != nil {
        return nil, err
    }

    return &Config{
        KafkaBrokers:    sharedcfg.ParseBrokers(sharedcfg.EnvOrDefault("KAFKA_BROKERS", "localhost:9092")),
        LogLevel:        sharedcfg.EnvOrDefault("LOG_LEVEL", "info"),
        ShutdownTimeout: shutdownTimeout,
        // ... service-specific fields
    }, nil
}
```

For the full list of environment variables per service, see:
- [ETL Configuration](https://github.com/couchcryptid/storm-data-etl/wiki/Configuration)
- [API Configuration](https://github.com/couchcryptid/storm-data-api/wiki/Configuration)

## Related

- [System Configuration](https://github.com/couchcryptid/storm-data-system/wiki/Configuration) -- environment variables across all services
- [ETL Configuration](https://github.com/couchcryptid/storm-data-etl/wiki/Configuration) -- ETL-specific settings built on shared parsers
- [API Configuration](https://github.com/couchcryptid/storm-data-api/wiki/Configuration) -- API-specific settings built on shared parsers
- [[Architecture]] -- package design and how services compose shared code
- [[Development]] -- testing and project conventions
