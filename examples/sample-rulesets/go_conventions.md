# Go Coding Conventions

## Overview

This ruleset defines Go coding standards following idiomatic Go practices and team conventions.

## Project Structure

Follow the standard Go project layout:

```
project/
├── cmd/              # Main applications
├── internal/         # Private application code
├── pkg/              # Public library code
├── test/             # Additional test files
├── go.mod
└── go.sum
```

## Naming Conventions

### Packages

- Use short, lowercase, single-word names
- Avoid underscores or mixed caps
- Choose names that describe the package's purpose

```go
// Good
package user
package http
package config

// Bad
package user_management
package httpUtils
```

### Variables and Functions

- Use `camelCase` for unexported names
- Use `PascalCase` for exported names
- Use short names for local variables with limited scope
- Use descriptive names for package-level variables

```go
// Good
var userCount int
func calculateTotal() int { }
func ParseRequest(r *http.Request) error { }

// Bad
var UserCount int  // Should be unexported
func CalculateTotal() int { }  // Should be unexported
```

### Interfaces

- Use `-er` suffix for single-method interfaces
- Keep interfaces small and focused

```go
// Good
type Reader interface {
    Read(p []byte) (n int, err error)
}

type UserRepository interface {
    FindByID(id int) (*User, error)
    Save(user *User) error
}

// Bad
type IUser interface { }  // Don't use "I" prefix
```

## Error Handling

### Return Errors

- Return errors as the last return value
- Don't panic in library code
- Use `errors.New()` or `fmt.Errorf()` for error messages

```go
// Good
func processUser(id int) (*User, error) {
    user, err := findUser(id)
    if err != nil {
        return nil, fmt.Errorf("failed to find user %d: %w", id, err)
    }
    return user, nil
}

// Bad
func processUser(id int) *User {
    user, err := findUser(id)
    if err != nil {
        panic(err)  // Don't panic
    }
    return user
}
```

### Error Wrapping

- Use `%w` verb to wrap errors
- Provide context when wrapping errors
- Check for specific errors with `errors.Is()` and `errors.As()`

```go
if err != nil {
    return fmt.Errorf("failed to process request: %w", err)
}

if errors.Is(err, ErrNotFound) {
    // Handle not found case
}
```

## Code Organization

### Function Length

- Keep functions short and focused
- Extract complex logic into helper functions
- Aim for functions under 50 lines

### Package Organization

- Group related functionality in packages
- Keep packages focused on a single responsibility
- Avoid circular dependencies

## Documentation

### Comments

- Write godoc comments for all exported types, functions, and constants
- Start comments with the name of the thing being described
- Use complete sentences

```go
// User represents a user account in the system.
type User struct {
    ID   int
    Name string
}

// FindByID retrieves a user by their unique identifier.
// It returns ErrNotFound if the user does not exist.
func FindByID(id int) (*User, error) {
    // Implementation
}
```

### Package Documentation

- Add package documentation in a `doc.go` file or at the top of a main file
- Describe the package's purpose and main concepts

```go
// Package user provides user account management functionality.
//
// It includes operations for creating, retrieving, updating, and
// deleting user accounts, as well as authentication and authorization.
package user
```

## Testing

### Test Files

- Place tests in `*_test.go` files
- Use the same package name with `_test` suffix for black-box tests
- Use the same package name without suffix for white-box tests

```go
// user_test.go (black-box)
package user_test

import "myapp/user"

// user_internal_test.go (white-box)
package user
```

### Test Functions

- Name tests with `Test` prefix followed by the function name
- Use table-driven tests for multiple scenarios
- Use subtests with `t.Run()` for better organization

```go
func TestCalculateDiscount(t *testing.T) {
    tests := []struct {
        name     string
        price    float64
        rate     float64
        expected float64
    }{
        {"10% discount", 100.0, 0.1, 90.0},
        {"50% discount", 100.0, 0.5, 50.0},
        {"no discount", 100.0, 0.0, 100.0},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := CalculateDiscount(tt.price, tt.rate)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## Concurrency

### Goroutines

- Don't leak goroutines
- Use context for cancellation
- Use channels for communication

```go
// Good
func processItems(ctx context.Context, items []Item) error {
    errCh := make(chan error, 1)
    
    go func() {
        errCh <- doWork(ctx, items)
    }()
    
    select {
    case err := <-errCh:
        return err
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

### Mutexes

- Use mutexes to protect shared state
- Keep critical sections small
- Consider using channels instead of mutexes when appropriate

```go
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    c.value++
    c.mu.Unlock()
}
```

## Code Formatting

- Use `gofmt` or `goimports` for formatting
- Run formatters before committing
- Configure your editor to format on save

## Linting

- Use `golangci-lint` for comprehensive linting
- Address all linting warnings
- Configure linters in `.golangci.yml`

## Dependencies

- Use Go modules for dependency management
- Keep dependencies minimal
- Regularly update dependencies
- Use `go mod tidy` to clean up unused dependencies
