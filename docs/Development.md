# Development

## Prerequisites

- Go 1.25+
- [golangci-lint](https://golangci-lint.run/) (for linting)
- [pre-commit](https://pre-commit.com/) (optional, for git hooks)

No Docker required -- this is a pure Go library with no infrastructure dependencies.

## Setup

```sh
git clone <repo-url>
cd storm-data-shared
```

Install pre-commit hooks (optional):

```sh
pre-commit install
```

## Testing

### Unit Tests

```sh
make test
```

Runs all tests with the race detector enabled (`-race -count=1`).

### Coverage

```sh
make test-cover
```

Generates `coverage.out` and opens an HTML coverage report in the browser.

## Linting

```sh
make lint
```

Uses `golangci-lint` with the configuration in `.golangci.yml`. Enabled linters include:

- `errcheck`, `govet`, `staticcheck` -- correctness
- `gosec` -- security
- `gocyclo` -- complexity (threshold: 15)
- `revive`, `gocritic` -- style
- `gofmt`, `goimports` -- formatting
- `misspell`, `unparam` -- hygiene

This matches the linter configuration used by the [ETL](https://github.com/couchcryptid/storm-data-etl) and [API](https://github.com/couchcryptid/storm-data-api) services.

## Formatting

```sh
make fmt
```

Runs `gofmt` and `goimports` across the project.

## Vulnerability Scanning

```sh
make vuln
```

Runs `govulncheck` against all packages.

## Pre-commit Hooks

The `.pre-commit-config.yaml` configures hooks that run on every commit:

- Trailing whitespace removal
- End-of-file newline
- YAML and JSON validation
- Merge conflict markers
- Secret detection (`gitleaks`, `detect-private-key`)
- Large file detection
- `gofmt` and `goimports`
- `golangci-lint`

## CI Pipeline

The `.github/workflows/ci.yml` workflow runs on pushes and pull requests to `main`:

| Job | What It Does |
|-----|-------------|
| `test` | `make test` (unit tests with race detector) |
| `lint` | `make lint` (golangci-lint) |
| `sonarcloud` | Unit tests with coverage + SonarCloud scan |

There is no release workflow or Docker image -- this is a library consumed via Go modules. Consuming services reference it by commit hash (pseudo-version).

## Versioning

This module uses Go pseudo-versions based on commit hashes rather than semantic version tags:

```
github.com/couchcryptid/storm-data-shared v0.0.0-20260211182606-5c0ac15abbdf
```

To update a consuming service to the latest commit:

```sh
go get github.com/couchcryptid/storm-data-shared@main
go mod tidy
```

## Project Conventions

- **Primitive arguments**: Shared functions take primitive types (`string`, `time.Duration`), not service-specific config structs. This keeps the library independent of consuming services.
- **Structural typing for interfaces**: The `ReadinessChecker` interface is defined here but services implement it without importing this module in their checker code -- Go's structural typing handles the match.
- **Co-located tests**: Tests live alongside the code they test (`env_test.go` next to `env.go`).
- **No generated code**: Pure library code with no code generation steps.

## Consuming Services

| Service | Import Path |
|---------|------------|
| [storm-data-etl](https://github.com/couchcryptid/storm-data-etl) | `config`, `observability`, `retry` |
| [storm-data-api](https://github.com/couchcryptid/storm-data-api) | `config`, `observability` |

## Related

- [System Development](https://github.com/couchcryptid/storm-data-system/wiki/Development) -- multi-repo workflow and cross-service conventions
- [ETL Development](https://github.com/couchcryptid/storm-data-etl/wiki/Development) -- consuming service development and testing
- [API Development](https://github.com/couchcryptid/storm-data-api/wiki/Development) -- consuming service development and testing
- [[Architecture]] -- package layout and design decisions
- [[Configuration]] -- shared environment variable parsers
- [[Code Quality]] -- linting, static analysis, and quality gates
