# Code Quality

Quality philosophy, tooling, and enforcement for the shared library. For pipeline-wide quality standards, see the [system Code Quality page](https://github.com/couchcryptid/storm-data-system/wiki/Code-Quality).

## Coding Philosophy

### Library, not framework

This module provides utility functions that services call -- it never calls back into service code. No init functions, no package-level side effects, no required initialization order. Services choose which packages to import and how to use them.

### Primitive arguments over config structs

Shared functions take primitive types (`string`, `time.Duration`, `int`) rather than service-specific config structs. `NewLogger("info", "json")` instead of `NewLogger(cfg)`. This keeps the library independent of consuming services and makes each function's requirements explicit.

### Structural typing for interfaces

The `ReadinessChecker` interface is defined here, but services implement it without importing this module in their checker code. Go's structural typing handles the match -- any type with `CheckReadiness(ctx context.Context) error` satisfies the interface.

### Thin wrappers in consumers

Services wrap shared functions in thin adapters that accept their own config types. For example, the ETL's `observability.NewLogger(cfg)` calls `sharedobs.NewLogger(cfg.LogLevel, cfg.LogFormat)`. This keeps the shared library decoupled while giving services a convenient, typed API.

### No generated code

Pure library code. No code generation, no build steps, no Makefiles beyond convenience targets. `go test ./...` is the entire test suite.

## Static Analysis

### golangci-lint

12 enabled linters -- a focused subset appropriate for a library with no HTTP clients or database connections:

| Category | Linters |
|----------|---------|
| Correctness | `errcheck`, `govet`, `staticcheck`, `errorlint`, `exhaustive` |
| Security | `gosec` |
| Style | `gocritic` (diagnostic/style/performance), `revive` (exported) |
| Complexity | `gocyclo` (threshold: 15) |
| Hygiene | `misspell`, `unparam`, `errname`, `unconvert`, `prealloc` |
| Test quality | `testifylint` |

`bodyclose`, `noctx`, and `sqlclosecheck` are omitted because this library makes no HTTP requests and has no database operations.

### govulncheck

Vulnerability scanning via `make vuln` runs `govulncheck` against all packages. This checks dependencies for known vulnerabilities in the Go vulnerability database.

### SonarQube Cloud

Analyzed via CI on every push and pull request: [SonarCloud dashboard](https://sonarcloud.io/summary/overall?id=couchcryptid_storm-data-shared)

SonarCloud configuration (`sonar-project.properties`):
- Reports Go coverage via `coverage.out`
- Allows idiomatic Go test naming (`TestX_Y_Z`) on test files

## Security

| Layer | What It Catches |
|-------|----------------|
| `gosec` | Security-sensitive patterns, weak crypto |
| `govulncheck` | Known vulnerabilities in dependencies |
| `gitleaks` | Secrets in source code |
| `detect-private-key` | Private key files accidentally committed |

## Quality Gates

### Pre-commit Hooks

`.pre-commit-config.yaml` runs on every commit:

- File hygiene: trailing whitespace, end-of-file newline, YAML/JSON validation
- Formatting: `gofmt`, `goimports`
- Linting: `golangci-lint` (5-minute timeout)
- Security: `gitleaks`, `detect-private-key`, `check-added-large-files`

### CI Pipeline

Every push and pull request to `main` runs:

| Job | Command | What It Enforces |
|-----|---------|-----------------|
| `test` | `make test` | Unit tests with race detector (`-race -count=1`) |
| `lint` | `make lint` | golangci-lint with 12 linters |
| `sonarcloud` | SonarCloud scan | Coverage, bugs, vulnerabilities, code smells, security hotspots |

There is no `build` job -- this is a library with no main package.

### SonarCloud Quality Gate

Uses the default "Sonar way" gate on new code: >= 80% coverage, <= 3% duplication, A ratings for reliability/security/maintainability, 100% security hotspots reviewed.

## Testing

All tests are unit tests with no infrastructure dependencies. See [[Development]] for commands.

| Scope | What It Covers |
|-------|---------------|
| `config` | Environment variable parsing, defaults, validation, edge cases |
| `observability` | Logger creation, health handler responses, JSON writer |
| `retry` | Backoff calculation, context-aware sleep |

All tests run with `-race -count=1`. Tests are co-located with source files (e.g., `env_test.go` next to `env.go`).

## Related

- [System Code Quality](https://github.com/couchcryptid/storm-data-system/wiki/Code-Quality) -- pipeline-wide quality standards
- [[Development]] -- commands, CI pipeline, versioning, conventions
- [[Architecture]] -- package design, interface contracts
- Consuming services: [ETL](https://github.com/couchcryptid/storm-data-etl/wiki/Code-Quality), [API](https://github.com/couchcryptid/storm-data-api/wiki/Code-Quality)
