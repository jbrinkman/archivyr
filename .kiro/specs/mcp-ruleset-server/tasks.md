# Implementation Plan

- [x] 1. Initialize project structure and dependencies
  - Create Go module with appropriate name and version
  - Add required dependencies: valkey-glide v2, mcp-go, zerolog, testify
  - Create directory structure following the design (cmd/, internal/, test/, docker/, examples/, docs/)
  - Initialize git repository with .gitignore for Go projects
  - _Requirements: 13.3_

- [x] 2. Implement configuration management
  - Create Config struct in internal/config package with fields for ValkeyHost, ValkeyPort, and LogLevel
  - Implement LoadConfig() function to read from environment variables with defaults (localhost:6379)
  - Implement Validate() function to ensure configuration values are valid
  - _Requirements: 9.1, 9.2, 9.3_

- [x] 2.1 Write unit tests for configuration
  - Test LoadConfig with various environment variable combinations
  - Test default value application
  - Test Validate function with valid and invalid inputs
  - _Requirements: 9.1, 9.2, 9.3_

- [x] 3. Create GitHub Actions CI workflow
  - Create .github/workflows/ci.yml
  - Configure triggers for pull requests and pushes to main
  - Set up Go environment (version 1.21+)
  - Add step to install dependencies
  - Add step to run golangci-lint
  - Add step to run unit tests with race detector and coverage
  - Add step to run integration tests
  - Add step to build Docker image
  - Add step to run container smoke tests
  - Add step to upload coverage reports
  - _Requirements: 11.1, 11.2, 11.3, 11.4_

- [x] 4. Implement Valkey client wrapper
  - Create Client struct in internal/valkey package wrapping valkey-glide GlideClient
  - Implement NewClient(host, port string) to create and connect to Valkey
  - Implement Close() for graceful connection shutdown
  - Implement Ping() for health checks
  - Add error handling for connection failures
  - _Requirements: 1.5, 9.4, 9.5_

- [x] 4.1 Write unit tests for Valkey client
  - Test NewClient with mock connections
  - Test connection error scenarios
  - Test Ping functionality
  - _Requirements: 1.5, 9.4, 9.5_

- [x] 5. Implement core ruleset data types and validation
  - Create Ruleset struct in internal/ruleset package with all metadata fields
  - Create RulesetUpdate struct for partial updates
  - Implement ValidateRulesetName in internal/util to enforce snake_case naming
  - Implement timestamp formatting and parsing utilities
  - _Requirements: 1.2, 1.3_

- [x] 5.1 Write unit tests for validation utilities
  - Test snake_case validation with valid and invalid names
  - Test timestamp formatting and parsing
  - _Requirements: 1.2_

- [x] 6. Implement ruleset service - Create operation
  - Create Service struct in internal/ruleset with Valkey client dependency
  - Implement NewService constructor
  - Implement Exists(name string) to check if ruleset exists using Valkey EXISTS command
  - Implement ListNames() to retrieve all ruleset names using KEYS command
  - Implement Create(ruleset *Ruleset) with duplicate name checking
  - Set created_at and last_modified timestamps on creation
  - Store ruleset in Valkey hash with key pattern "ruleset:{name}"
  - Return error with existing names list if duplicate detected
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 1.1, 1.2, 1.3_

- [x] 6.1 Write unit tests for Create operation
  - Test successful creation
  - Test duplicate name detection and error response
  - Test timestamp setting
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_

- [x] 7. Implement ruleset service - Read operations
  - Implement Get(name string) to retrieve ruleset by exact name using HGETALL
  - Parse hash fields back into Ruleset struct
  - Return error if ruleset not found
  - Implement List() to retrieve all rulesets with metadata
  - Implement Search(pattern string) for glob pattern matching using KEYS with pattern
  - _Requirements: 3.1, 3.2, 3.3, 4.1, 4.2, 4.3, 4.4, 5.1, 5.2, 5.3_

- [x] 7.1 Write unit tests for Read operations
  - Test Get with existing and non-existent rulesets
  - Test List with empty and populated store
  - Test Search with various patterns
  - _Requirements: 3.1, 3.2, 3.3, 4.1, 4.2, 4.3, 4.4, 5.1, 5.2, 5.3_

- [x] 8. Implement ruleset service - Update operation
  - Implement Update(name string, updates *RulesetUpdate) to modify existing rulesets
  - Check if ruleset exists before updating
  - Update only provided fields (description, tags, markdown)
  - Update last_modified timestamp, preserve created_at
  - Return error if ruleset not found
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

- [x] 8.1 Write unit tests for Update operation
  - Test successful update of each field
  - Test partial updates
  - Test timestamp handling
  - Test error on non-existent ruleset
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

- [x] 9. Implement ruleset service - Delete operation
  - Implement Delete(name string) to remove ruleset from Valkey using DEL command
  - Check if ruleset exists before deleting
  - Return confirmation message on success
  - Return error with list of existing names if ruleset not found
  - _Requirements: 7.1, 7.2, 7.3, 7.4_

- [x] 9.1 Write unit tests for Delete operation
  - Test successful deletion
  - Test error response with existing names list
  - Test deletion of non-existent ruleset
  - _Requirements: 7.1, 7.2, 7.3, 7.4_

- [x] 10. Implement MCP resource handler
  - Create Handler struct in internal/mcp package with ruleset service dependency
  - Implement RegisterResources() to register ruleset resources with URI scheme "ruleset://{name}"
  - Implement resource handler function to call service.Get() for exact name matches
  - Format response as text/markdown with metadata
  - Handle resource not found errors
  - _Requirements: 3.1, 3.2, 3.3, 3.4_

- [x] 11. Implement MCP tool handlers
  - Implement RegisterTools() to register all CRUD tools
  - Implement create_ruleset tool handler with parameters: name, description, tags, markdown
  - Implement get_ruleset tool handler with parameter: name
  - Implement update_ruleset tool handler with parameters: name, description, tags, markdown (all optional except name)
  - Implement delete_ruleset tool handler with parameter: name
  - Implement list_rulesets tool handler with no parameters
  - Implement search_rulesets tool handler with parameter: pattern
  - Map tool calls to corresponding service methods
  - Format tool responses according to MCP protocol
  - _Requirements: 2.1, 2.2, 2.3, 3.1, 4.1, 4.2, 4.3, 4.4, 5.1, 5.2, 5.3, 6.1, 6.2, 6.3, 6.4, 6.5, 7.1, 7.2, 7.3, 7.4_

- [x] 11.1 Write unit tests for MCP tool handlers
  - Test each tool handler with valid inputs
  - Test error scenarios for each tool
  - Test parameter validation
  - _Requirements: 2.1, 3.1, 4.1, 5.1, 6.1, 7.1_

- [x] 12. Implement MCP server initialization and stdio transport
  - Implement NewHandler(service *ruleset.Service) constructor
  - Implement Start() to initialize MCP server with stdio transport using mcp-go
  - Configure server with registered resources and tools
  - Set up stdio reader/writer for MCP protocol communication
  - Implement MCP protocol initialization and capability negotiation
  - Add structured logging for server events using zerolog
  - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5_

- [ ] 13. Implement main application entry point
  - Create main.go in cmd/mcp-ruleset-server
  - Load configuration using config package
  - Initialize zerolog logger with configured log level
  - Create Valkey client and test connection with Ping
  - Exit with error if Valkey connection fails
  - Create ruleset service with Valkey client
  - Create and start MCP handler
  - Implement graceful shutdown on SIGINT/SIGTERM
  - Close Valkey connection on shutdown
  - _Requirements: 9.1, 9.2, 9.3, 9.4, 9.5, 8.1, 8.2, 8.3, 8.4, 8.5_

- [ ] 14. Write integration tests for Valkey operations
  - Use testcontainers to spin up Valkey instance
  - Test full CRUD workflow against real Valkey
  - Test concurrent operations
  - Test connection handling and retries
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_

- [ ] 15. Write integration tests for MCP protocol
  - Test tool invocations with sample JSON-RPC payloads
  - Test resource retrieval via URI
  - Test error responses format
  - Test stdio transport communication
  - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5_

- [ ] 16. Create Dockerfile with multi-stage build
  - Create Dockerfile in docker/ directory
  - Stage 1: Use golang:1.21-alpine to build Go binary
  - Copy go.mod, go.sum and download dependencies
  - Copy source code and build static binary
  - Stage 2: Use alpine:latest as runtime base
  - Install Valkey server in runtime image
  - Copy compiled binary from builder stage
  - Set up working directory and permissions
  - _Requirements: 10.1, 10.2, 10.3_

- [ ] 17. Create Docker entrypoint script
  - Create docker-entrypoint.sh in docker/ directory
  - Start Valkey server in daemon mode
  - Implement health check loop using valkey-cli ping
  - Start MCP server after Valkey is ready
  - Handle shutdown signals (SIGTERM, SIGINT)
  - Make script executable
  - _Requirements: 10.4, 10.5_

- [ ] 18. Write end-to-end Docker tests
  - Build Docker image in test
  - Start container with testcontainers
  - Test MCP server availability via stdio
  - Test full CRUD workflow through container
  - Test graceful shutdown
  - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5_

- [ ] 19. Create golangci-lint configuration
  - Create .golangci.yml with linter configuration
  - Enable recommended linters (govet, errcheck, staticcheck, etc.)
  - Configure linter rules for code quality
  - Set timeout and concurrency settings
  - _Requirements: 11.4_

- [ ] 20. Create GitHub Actions CD workflow
  - Create .github/workflows/cd.yml
  - Configure manual workflow dispatch trigger
  - Set up semantic-release with conventional commits analyzer
  - Configure semantic-release to generate changelog and determine version
  - Add step to create Git tag based on semantic version
  - Add step to build multi-arch Docker images (amd64, arm64)
  - Add step to push images to registry with version tags (latest, vX.Y.Z, vX.Y, vX)
  - Add step to create GitHub release with generated changelog
  - _Requirements: 11.5, 12.1, 12.2, 12.3, 12.4, 12.5_

- [ ] 21. Create semantic-release configuration
  - Create .releaserc.json with conventional commits configuration
  - Configure commit analyzer for feat, fix, docs, chore, etc.
  - Configure release notes generator
  - Configure changelog generation to CHANGELOG.md
  - Configure Git assets and tagging
  - _Requirements: 12.1, 12.2, 12.5_

- [ ] 22. Create project documentation
  - Create README.md with project overview, quick start, Docker usage, and MCP client configuration examples
  - Create CONTRIBUTING.md with conventional commit guidelines and development workflow
  - Create LICENSE file (choose appropriate license)
  - Create examples/mcp-config.json with sample MCP client configuration
  - Create examples/sample-rulesets/ with example ruleset markdown files
  - _Requirements: 13.1, 13.2, 13.4, 13.5_

- [ ] 23. Create API documentation
  - Create docs/API.md with complete MCP tool reference
  - Document all tool parameters and return values
  - Document resource URI schemes
  - Include request/response examples for each operation
  - Document error codes and handling
  - _Requirements: 13.1, 13.5_

- [ ] 24. Create architecture documentation
  - Create docs/ARCHITECTURE.md with detailed component descriptions
  - Include data flow diagrams
  - Document design decisions and rationale
  - Document extension points for future features
  - _Requirements: 13.1, 13.5_

- [ ] 25. Final integration and testing
  - Build complete Docker image
  - Test full workflow: start container, create/read/update/delete rulesets, stop container
  - Verify data persistence across container restarts
  - Test with actual MCP client (e.g., Claude Desktop, Cursor)
  - Verify conventional commit enforcement in CI
  - Test manual CD workflow trigger
  - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5, 11.1, 11.2, 11.3, 11.4, 11.5, 12.1, 12.2, 12.3, 12.4, 12.5_
