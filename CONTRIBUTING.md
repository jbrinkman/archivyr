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

## Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes with proper signoffs
4. Run tests and linter
5. Submit a pull request

## Code Standards

- Follow Go best practices
- Maintain test coverage above 80%
- Run `golangci-lint run` before submitting
- Add tests for new features

## Questions?

Open an issue for discussion before starting major changes.
