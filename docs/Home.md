# Storm Data Shared

Shared Go library for the storm data pipeline. Extracts common patterns from the [ETL](https://github.com/couchcryptid/storm-data-etl) and [API](https://github.com/couchcryptid/storm-data-api) services into a single importable module.

## What It Provides

| Package | Purpose | Used By |
|---------|---------|---------|
| `config` | Environment variable helpers and shared parsers | [ETL](https://github.com/couchcryptid/storm-data-etl), [API](https://github.com/couchcryptid/storm-data-api) |
| `observability` | Structured logging, health endpoints, `ReadinessChecker` interface | [ETL](https://github.com/couchcryptid/storm-data-etl), [API](https://github.com/couchcryptid/storm-data-api) |
| `retry` | Exponential backoff and context-aware sleep | [ETL](https://github.com/couchcryptid/storm-data-etl) |

## Pages

- [[Architecture]] -- Package design, interface contracts, and design decisions

## Related Repositories

| Repository | Relationship |
|------------|-------------|
| [storm-data-etl](https://github.com/couchcryptid/storm-data-etl) | Imports `config`, `observability`, `retry` |
| [storm-data-api](https://github.com/couchcryptid/storm-data-api) | Imports `config`, `observability` |
| [storm-data-collector](https://github.com/couchcryptid/storm-data-collector) | TypeScript -- does not use this library |
| [storm-data-system](https://github.com/couchcryptid/storm-data-system) | System orchestration and E2E tests |
