# Code Coverage Guide

## Overview

This project uses Go's built-in coverage tooling with the `-coverpkg` flag to capture coverage across unit tests and integration tests.

## Understanding Coverage

### Coverage Types

1. **Unit Test Coverage** (`./internal/...`)
   - Tests individual functions and methods
   - Fast, isolated tests
   - Located in `*_test.go` files alongside source code

2. **Integration Test Coverage** (`./test/integration/...`)
   - Tests components working together with real dependencies (Valkey via testcontainers)
   - Uses `-coverpkg=./internal/...` to instrument application code
   - Captures coverage for code paths that unit tests can't easily reach

3. **E2E Test Coverage** (`./test/e2e/...`)
   - Tests full Docker container workflow
   - **Cannot capture coverage** (runs in separate Docker process)
   - Provides confidence but no coverage metrics

## Running Coverage

### Combined Coverage (Recommended)

Run all tests and generate a combined coverage report:

```bash
task test:coverage:combined
```

This will:

1. Run unit tests with coverage
2. Run integration tests with `-coverpkg=./internal/...` to capture application coverage
3. Merge both coverage reports into `reports/coverage/coverage.out`
4. Save test results to `reports/test/`
4. Display total coverage percentage

### View HTML Report

```bash
task test:coverage:html
```

Then open `reports/coverage/coverage.html` in your browser.

### Report Locations

All generated reports are stored in the `reports/` directory (gitignored):

- `reports/coverage/coverage.out` - Combined coverage data
- `reports/coverage/coverage.html` - HTML coverage visualization  
- `reports/coverage/data/` - Individual coverage files (unit, integration)
- `reports/test/*.json` - Test results in JSON format

### Individual Test Types

```bash
# Unit tests only
task test:unit

# Integration tests only  
task test:integration

# E2E tests only (no coverage)
task test:e2e
```

## The `-coverpkg` Flag

The key to capturing integration test coverage is the `-coverpkg` flag:

```bash
go test ./test/integration/... \
  -coverprofile=integration.out \
  -coverpkg=./internal/... \
  -covermode=atomic
```

**What it does:**

- Tells Go to instrument ALL packages matching `./internal/...`
- Even though tests are in `./test/integration/...`, coverage is captured for `./internal/...`
- Without this flag, integration tests would show 0% coverage

## Coverage Targets

- **Minimum Target**: 80% combined coverage
- **Current Coverage**: Check with `task test:coverage:combined`

### Package Breakdown

- `internal/config`: 100% (configuration loading and validation)
- `internal/mcp`: 86%+ (MCP protocol handlers)
- `internal/ruleset`: 88%+ (business logic)
- `internal/validation`: 100% (validation utilities)
- `internal/valkey`: 58% unit + integration coverage (infrastructure code)

**Note**: `internal/valkey` has lower unit test coverage because:

- Unit tests focus on validation and error handling
- Integration tests cover successful client creation and operations
- Combined coverage shows the true picture

## Why E2E Tests Don't Contribute to Coverage

E2E tests run the application in Docker containers, which means:

- The code runs in a completely separate process
- Go's coverage instrumentation can't reach across process boundaries
- These tests still provide valuable end-to-end validation
- They just don't contribute to the coverage percentage

## Troubleshooting

### Integration tests show 0% coverage

Make sure you're using `-coverpkg=./internal/...`:

```bash
go test ./test/integration/... -coverpkg=./internal/... -coverprofile=coverage.out
```

### Coverage report is empty

Check that the coverage file has content:

```bash
wc -l coverdata/integration.out
```

If it only has 1 line (`mode: atomic`), the `-coverpkg` flag wasn't applied.

### Want to see what integration tests cover?

Run just integration tests with coverage:

```bash
go test ./test/integration/... -coverpkg=./internal/... -coverprofile=integration.out
go tool cover -html=integration.out -o integration-coverage.html
```

## References

- [Go Coverage Documentation](https://go.dev/blog/cover)
- [Using -coverpkg for integration tests](https://go.dev/testing/coverage/)
