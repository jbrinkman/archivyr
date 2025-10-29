# Testing Guidelines

## Test Execution Requirements

Before marking ANY task as complete, ALL tests in the project MUST pass. This is a strict requirement with no exceptions.

### Pre-Task Assumptions

- When starting any task, assume there are no failing tests in the codebase
- The codebase should be in a clean, passing state before beginning work

### During Task Execution

- If tests fail during task execution, assume YOUR changes caused the failures
- Regardless of the actual cause, YOU are responsible for fixing all failing tests
- Tests from "unrelated" features may fail due to your changes - fix them

### Before Task Completion

Run the full test suite before marking a task complete:

```bash
go test ./... -v -race -cover
```

All tests must pass with:

- No test failures
- No race conditions detected
- Acceptable code coverage maintained

### Test Failure Protocol

If tests fail:

1. **Identify the failure**: Read the test output carefully
2. **Understand the cause**: Determine if your changes broke existing functionality
3. **Fix the issue**: Update your code or the tests as needed
4. **Verify the fix**: Re-run all tests to confirm they pass
5. **Only then**: Mark the task as complete

### No Exceptions

- Do not mark tasks complete with failing tests
- Do not skip test failures assuming they're "unrelated"
- Do not defer test fixes to future tasks
- Do not commit code with failing tests

## Test Quality Standards

### Focus on Core Functionality

- Write tests that validate core business logic
- Avoid over-testing edge cases unless critical
- Keep tests minimal but comprehensive

### Test Organization

- Unit tests alongside source files (*_test.go)
- Integration tests in test/integration/
- E2E tests in test/e2e/
- Minimum 80% code coverage target

### Test Naming

- Use descriptive test names that explain what is being tested
- Follow the pattern: TestFunctionName or TestFunctionName_Scenario
- Use table-driven tests for multiple scenarios

### Test Independence

- Each test should be independent and not rely on other tests
- Tests should clean up after themselves
- Use testcontainers for integration tests requiring external services

## Continuous Integration

All tests run automatically in CI/CD pipelines. Local test failures will result in CI failures, blocking merges.
