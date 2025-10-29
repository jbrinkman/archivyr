# Command Safety Guidelines

## Avoiding Dangerous Command Patterns

When executing terminal commands, avoid patterns that trigger manual approval requirements. These patterns are designed to catch potentially dangerous operations, but can interfere with trusted workflows.

## Dangerous Patterns to Avoid

The following patterns will always require manual approval:

1. **Shell operators**: `|`, `>`, `>>`, `<`, `&&`, `||`, `;`
2. **Command substitution**: `$()`, `` ` ` ``
3. **Wildcards in sensitive contexts**: `rm *`, `chmod *`, etc.
4. **Sudo/privilege escalation**: `sudo`, `su`, `doas`
5. **System modification**: `rm -rf`, `mkfs`, `dd`, `fdisk`
6. **Network operations**: `curl | sh`, `wget | bash`

## Best Practices for Test Execution

### Running Tests

**Avoid:**

```bash
# Pipes trigger approval
dotnet test --logger "console;verbosity=normal" 2>&1 | head -200

# Semicolons in arguments can trigger approval
dotnet test --logger "console;verbosity=detailed"
```

**Prefer:**

```bash
# Simple, clean commands
dotnet test tests/Valkey.Glide.IntegrationTests/ --framework net8.0

# Filter specific tests
dotnet test tests/Valkey.Glide.IntegrationTests/ --framework net8.0 --filter "FullyQualifiedName~TestClassName"

# Use minimal verbosity options
dotnet test --framework net8.0 --verbosity minimal
```

### Viewing Output

**Avoid:**

```bash
# Pipes require approval
cat file.txt | grep pattern
dotnet build 2>&1 | tail -50
```

**Prefer:**

```bash
# Use tool-specific options
dotnet build --verbosity quiet

# Read files directly with tools
# Use readFile or grepSearch tools instead of cat/grep
```

### Building Projects

**Avoid:**

```bash
# Chaining commands
dotnet clean && dotnet build

# Output redirection
dotnet build > output.log 2>&1
```

**Prefer:**

```bash
# Single commands
dotnet build sources/Valkey.Glide/ --framework net8.0

# Use --verbosity to control output
dotnet build --verbosity minimal
```

## Command Alternatives

| Instead of | Use |
|------------|-----|
| `command \| head` | Run command without pipe, output is automatically limited |
| `command \| tail` | Run command without pipe, check end of output in result |
| `command \| grep` | Use `grepSearch` tool instead |
| `cat file \| grep` | Use `grepSearch` tool instead |
| `command1 && command2` | Make separate tool calls for each command |
| `command > file` | Capture output in tool result, then write with `fsWrite` |
| `--logger "console;verbosity=X"` | Omit logger flag (console is default) or use `--verbosity X` |

## Test Execution Patterns

### Unit Tests

```bash
# Simple execution
dotnet test tests/Valkey.Glide.UnitTests/ --framework net8.0

# With filtering
dotnet test tests/Valkey.Glide.UnitTests/ --framework net8.0 --filter "ClassName"
```

### Integration Tests

```bash
# All integration tests
dotnet test tests/Valkey.Glide.IntegrationTests/ --framework net8.0

# Specific test class
dotnet test tests/Valkey.Glide.IntegrationTests/ --framework net8.0 --filter "FullyQualifiedName~PubSubCommandTests"

# Specific test method
dotnet test tests/Valkey.Glide.IntegrationTests/ --framework net8.0 --filter "FullyQualifiedName~PubSubCommandTests.PublishAsync_WithNoSubscribers_ReturnsZero"
```

### Coverage

```bash
# Use Task commands when available
task coverage:unit

# Or direct dotnet commands
dotnet test --framework net8.0 --collect:"XPlat Code Coverage"
```

## General Principles

1. **Keep commands simple**: Single-purpose commands are less likely to trigger approval
2. **Avoid shell features**: Don't use pipes, redirects, or command chaining
3. **Use tool-specific options**: Prefer built-in flags over shell operators
4. **Leverage Kiro tools**: Use `readFile`, `grepSearch`, `fsWrite` instead of shell commands
5. **Trust the output**: Kiro automatically captures and presents command output

## When Manual Approval is Necessary

Some operations legitimately require manual approval for safety:

- Installing system packages
- Modifying system files
- Running scripts from the internet
- Deleting large numbers of files
- Operations requiring elevated privileges

For these cases, manual approval is appropriate and should be requested.
