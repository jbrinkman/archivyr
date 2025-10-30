# Archivyr Examples

This directory contains example configurations and sample rulesets to help you get started with Archivyr.

## MCP Client Configuration

### mcp-config.json

Example MCP client configuration for connecting to Archivyr. This configuration can be adapted for various MCP-compatible clients like Claude Desktop, Cursor, and others.

**Usage:**

For Claude Desktop, add the contents to your configuration file:

- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

For other clients, consult their documentation for MCP server configuration.

## Sample Rulesets

The `sample-rulesets/` directory contains example rulesets demonstrating different use cases:

### python_style_guide.md

Python coding standards based on PEP 8 and team conventions. Covers:

- Naming conventions (variables, functions, classes, constants)
- Code organization and imports
- Documentation with docstrings
- Error handling best practices
- Type hints usage
- Testing standards
- Code formatting and linting

**Use this for:** Python projects requiring consistent coding standards.

### go_conventions.md

Go coding conventions following idiomatic Go practices. Covers:

- Project structure and package organization
- Naming conventions for packages, types, and functions
- Error handling and wrapping
- Documentation with godoc comments
- Testing with table-driven tests
- Concurrency patterns (goroutines, channels, mutexes)
- Code formatting and linting

**Use this for:** Go projects following standard Go idioms.

### api_design_principles.md

RESTful API design principles for building consistent, intuitive APIs. Covers:

- Resource naming conventions
- HTTP method semantics
- Status code usage
- Request/response format standards
- Pagination, filtering, and sorting
- API versioning strategies
- Security best practices
- API documentation requirements

**Use this for:** Backend projects exposing RESTful APIs.

### testing_best_practices.md

Comprehensive testing guidelines across all testing levels. Covers:

- Testing pyramid (unit, integration, E2E)
- AAA pattern (Arrange-Act-Assert)
- Test organization and naming
- Mocking and stubbing strategies
- Test data management
- Test coverage goals
- Performance considerations
- CI/CD integration

**Use this for:** Any project requiring robust testing practices.

## Using Sample Rulesets

### Option 1: Import into Archivyr

Once you have Archivyr running, you can create these rulesets using the `create_ruleset` tool:

```
Create a ruleset named "python_style_guide" with the content from examples/sample-rulesets/python_style_guide.md
```

### Option 2: Customize for Your Team

1. Copy a sample ruleset to your local directory
2. Modify it to match your team's conventions
3. Import the customized version into Archivyr

### Option 3: Use as Templates

Use these samples as templates to create your own rulesets:

- Copy the structure and sections you need
- Add team-specific guidelines
- Remove sections that don't apply
- Add examples from your codebase

## Creating Your Own Rulesets

When creating rulesets, follow these guidelines:

### Naming Convention

Use `snake_case` for ruleset names:

- ✅ `python_style_guide`
- ✅ `api_design_principles`
- ✅ `react_component_patterns`
- ❌ `PythonStyleGuide`
- ❌ `api-design-principles`

### Structure

Organize your rulesets with clear sections:

1. **Overview**: Brief description of the ruleset's purpose
2. **Main Sections**: Organized by topic or category
3. **Examples**: Code examples showing good and bad practices
4. **Best Practices**: Summary of key points

### Content

- Use markdown formatting for readability
- Include code examples with syntax highlighting
- Provide both good and bad examples
- Keep it concise and actionable
- Focus on "why" not just "what"

### Metadata

When creating rulesets, include useful metadata:

- **Description**: Brief summary of the ruleset
- **Tags**: Categorization tags (e.g., ["python", "style", "pep8"])

## Listing and Searching Rulesets

The `search_rulesets` tool serves dual purposes:

### Listing All Rulesets

To see all available rulesets, simply omit the pattern parameter:

```
List all available rulesets
```

Or explicitly use the wildcard pattern:

```
Search for rulesets matching "*"
```

Both approaches return the same result: a complete list of all rulesets with their metadata.

### Searching with Patterns

To find specific rulesets, provide a glob pattern:

```
Search for rulesets matching "*python*"
Search for rulesets matching "style_*"
Search for rulesets matching "*_guide"
```

**Pattern Examples**:

- `*python*` - Matches any ruleset with "python" in the name
- `style_*` - Matches rulesets starting with "style_"
- `*_guide` - Matches rulesets ending with "_guide"
- `api_*` - Matches rulesets starting with "api_"

## Example Workflow

Here's a typical workflow for using Archivyr with these samples:

1. **Start Archivyr**:

   ```bash
   docker run -i ghcr.io/joebrinkman/archivyr:latest
   ```

2. **Import sample rulesets** using your AI editor

3. **List all rulesets** to see what's available:

   ```
   List all available rulesets
   ```

4. **Search for specific rulesets** by pattern:

   ```
   Search for rulesets matching "*style*"
   ```

5. **Reference rulesets** in your projects:

   ```
   Use the python_style_guide ruleset for this Python code
   ```

6. **Update rulesets** as your standards evolve:

   ```
   Update the python_style_guide with new section about async/await patterns
   ```

## Contributing

Have a useful ruleset to share? Consider contributing it to this repository:

1. Create a new markdown file in `sample-rulesets/`
2. Follow the naming and structure guidelines above
3. Submit a pull request with your ruleset

## Questions?

For more information about Archivyr, see the main [README.md](../README.md) in the project root.
