# Contributing to Archivyr

Thank you for your interest in contributing to Archivyr!

## Developer Certificate of Origin (DCO)

All commits must include a DCO signoff. This certifies that you have the right to submit the contribution under the project's open source license.

### Adding Signoffs

Always use the `-s` flag when committing:

```bash
git commit -s -m "Your commit message"
```

Or configure Git to automatically add signoffs:

```bash
git config --global format.signOff true
```

### Signoff Format

Each commit must end with:

```
Signed-off-by: Your Name <your.email@example.com>
```

## Commit Message Format

Follow conventional commit format:

```
<type>(<scope>): <description>

[optional body]

Signed-off-by: Your Name <your.email@example.com>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

## Development Workflow

### Prerequisites

- Go 1.21 or later
- Docker (for integration tests and local deployment)
- [Task](https://taskfile.dev) (recommended for build automation)

### Setting Up Your Development Environment

```bash
# Clone the repository
git clone https://github.com/joebrinkman/archivyr.git
cd archivyr

# Install dependencies
task deps

# Set up development environment
task dev:setup
```

### Development Tasks

This project uses [Task](https://taskfile.dev) for build automation. Always use Task commands instead of direct Go commands to ensure consistency with CI.

```bash
# Show all available tasks
task

# Format code
task fmt

# Run linter (same as CI)
task lint

# Run linter with auto-fix
task lint:fix

# Run all tests (same as CI)
task test

# Run unit tests only
task test:unit

# Run integration tests
task test:integration

# Run end-to-end tests
task test:e2e

# Generate coverage report
task test:coverage

# Build binary
task build

# Run full CI pipeline locally
task ci

# Verify code before commit (format, lint, test)
task verify
```

### Making Changes

1. Create a feature branch from `main`:

   ```bash
   git checkout -b feat/your-feature-name
   ```

2. Make your changes following the code standards below

3. Add tests for new functionality

4. Run verification before committing:

   ```bash
   task verify
   ```

5. Commit with proper signoff and conventional commit format:

   ```bash
   git commit -s -m "feat(scope): description"
   ```

6. Push your branch and create a pull request

### Testing Requirements

Before marking any task complete or submitting a PR:

1. **All tests must pass**:

   ```bash
   task test
   ```

2. **All linting checks must pass**:

   ```bash
   task lint
   ```

3. **Code coverage must be maintained** (minimum 80%)

### Pull Request Process

1. Fork the repository
2. Create a feature branch with a descriptive name
3. Make your changes with proper DCO signoffs
4. Ensure all tests and linting pass
5. Update documentation if needed
6. Submit a pull request with a clear description
7. Address any review feedback
8. Wait for CI checks to pass
9. Maintainers will merge when approved

### Code Review Guidelines

- PRs should be focused and address a single concern
- Include tests for new functionality
- Update documentation for user-facing changes
- Respond to review feedback promptly
- Keep commits clean and well-organized

## Code Standards

### Go Best Practices

- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` for formatting (automated via `task fmt`)
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Keep functions small and focused
- Use meaningful variable and function names

### Testing Standards

- Maintain test coverage above 80%
- Write unit tests for all business logic
- Use integration tests for Valkey interactions
- Use end-to-end tests for Docker workflows
- Test error paths and edge cases
- Use table-driven tests for multiple scenarios

### Project Structure

Follow the established project structure:

```
cmd/              # Application entry points
internal/         # Private application code
  config/         # Configuration management
  valkey/         # Valkey client wrapper
  ruleset/        # Core business logic
  mcp/            # MCP protocol handlers
  util/           # Shared utilities
test/             # Test suites
  integration/    # Integration tests
  e2e/            # End-to-end tests
```

### Naming Conventions

- Packages: lowercase, single word
- Files: lowercase with underscores
- Types: PascalCase
- Functions: PascalCase (exported), camelCase (private)
- Ruleset names: snake_case (enforced)

### Documentation

- Add godoc comments for exported types and functions
- Update README.md for user-facing changes
- Update API documentation for new tools or resources
- Include examples for new features

## Continuous Integration

All pull requests run through automated CI checks:

- Code formatting verification
- Linting with golangci-lint
- Unit tests with race detector
- Integration tests with testcontainers
- End-to-end Docker tests
- Code coverage reporting

Ensure your code passes all CI checks before requesting review.

## Release Process

Releases are automated using semantic-release:

- Version numbers determined by conventional commits
- Changelog generated automatically
- Docker images built and published
- GitHub releases created with release notes

This is why following conventional commit format is critical.

## Questions?

- Open an issue for bug reports or feature requests
- Start a discussion for questions or ideas
- Reach out to maintainers for major changes before starting work
