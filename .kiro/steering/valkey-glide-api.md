---
inclusion: always
---

# Valkey GLIDE Go Client API Reference

## Context7 Library ID

When you need to look up Valkey GLIDE Go client documentation, use this Context7-compatible library ID:

**Library ID:** `/valkey-io/valkey-glide`

## Go Module

```go
import glide "github.com/valkey-io/valkey-glide/go/v2"
```

## Quick API Reference

### Common Operations

**Hash Operations:**

```go
// HSet - Set hash fields
fields := map[string]string{
    "field1": "value1",
    "field2": "value2",
}
client.HSet(ctx, "key", fields)

// HGet - Get single field
result, err := client.HGet(ctx, "key", "field")
if !result.IsNil() {
    value := result.Value()
}

// HGetAll - Get all fields
allFields, err := client.HGetAll(ctx, "key")
```

**Key Operations:**

```go
// Exists - Check if keys exist
count, err := client.Exists(ctx, []string{"key1", "key2"})

// Scan - Iterate through keys
cursor := models.NewCursor()
for {
    result, err := client.Scan(ctx, cursor)
    keys := result.Data
    cursor = result.Cursor
    if cursor.IsFinished() {
        break
    }
}
```

**String Operations:**

```go
// Set
client.Set(ctx, "key", "value")

// Get
result, err := client.Get(ctx, "key")
if !result.IsNil() {
    value := result.Value()
}
```

## Important Notes

- Use `models.NewCursor()` without arguments to create a new cursor
- `Scan()` takes only context and cursor, no options parameter
- Hash operations use `map[string]string` for field-value pairs
- Results have `.IsNil()` method to check for nil values
- Use `.Value()` to extract the actual value from results
