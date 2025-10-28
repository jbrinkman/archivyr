# Requirements Document

## Introduction

This document specifies the requirements for an MCP (Model Context Protocol) server that provides centralized storage and management of AI editor rulesets. The system addresses the problem of ruleset duplication and drift across multiple projects by providing a single source of truth for common rulesets stored in Valkey. The server will be implemented in Go, use valkey-glide for Go V2 for Valkey interactions, leverage the MCP-Go module for MCP protocol implementation, and be distributed as a Docker image with a bundled Valkey instance.

## Glossary

- **MCP Server**: The Model Context Protocol server application that handles client requests for ruleset management
- **Valkey**: An open-source key-value data store used as the backing storage system
- **Ruleset**: A collection of guidelines, rules, or steering documents for AI editors, stored as markdown content with metadata
- **valkey-glide**: The Go client library (V2) used to interact with Valkey
- **MCP-Go**: The Go module/library used to implement the MCP protocol
- **stdio Transport**: Standard input/output communication mechanism for MCP protocol
- **Docker Image**: The containerized distribution package containing both the MCP Server and bundled Valkey instance
- **Resource**: An MCP protocol concept for retrieving exact-match content by URI
- **Tool**: An MCP protocol concept for executing operations like search, create, update, delete
- **Hash**: A Valkey data structure that stores field-value pairs, used to store ruleset metadata and content

## Requirements

### Requirement 1

**User Story:** As a developer using multiple AI editors across different projects, I want to store rulesets in a centralized location, so that I can maintain consistency and avoid duplication.

#### Acceptance Criteria

1. THE MCP Server SHALL store rulesets in Valkey using hash data structures
2. WHEN a ruleset is stored, THE MCP Server SHALL use snake_case naming convention for the ruleset key
3. THE MCP Server SHALL store the following fields in each ruleset hash: created_at, last_modified, description, tags, and markdown
4. THE MCP Server SHALL persist ruleset data in Valkey with durability guarantees
5. THE MCP Server SHALL connect to a Valkey instance using valkey-glide for Go V2 client library

### Requirement 2

**User Story:** As a user of the MCP server, I want to create new rulesets with metadata, so that I can organize and describe my AI editor guidelines.

#### Acceptance Criteria

1. THE MCP Server SHALL provide a tool to create new rulesets with name, description, tags, and markdown content
2. WHEN a create operation is requested with a name that already exists, THE MCP Server SHALL return an error message containing a list of existing ruleset names
3. WHEN a create operation is requested with a name that already exists, THE MCP Server SHALL prompt the user to provide a different name
4. WHEN a ruleset is created, THE MCP Server SHALL set the created_at field to the current timestamp
5. WHEN a ruleset is created, THE MCP Server SHALL set the last_modified field to the current timestamp

### Requirement 3

**User Story:** As a user of the MCP server, I want to retrieve rulesets by exact name, so that I can quickly access specific guidelines for my projects.

#### Acceptance Criteria

1. THE MCP Server SHALL expose rulesets as Resources for exact-match retrieval by name
2. WHEN a Resource is requested with an exact ruleset name, THE MCP Server SHALL return the complete ruleset including all metadata fields and markdown content
3. IF a requested ruleset name does not exist, THEN THE MCP Server SHALL return an appropriate error message
4. THE MCP Server SHALL format Resource URIs using a consistent scheme for ruleset identification

### Requirement 4

**User Story:** As a user of the MCP server, I want to search for rulesets by pattern, so that I can discover available rulesets that match my needs.

#### Acceptance Criteria

1. THE MCP Server SHALL provide a tool to search for rulesets using pattern matching
2. WHEN a search operation is requested, THE MCP Server SHALL return a list of ruleset names matching the provided pattern
3. THE MCP Server SHALL support wildcard pattern matching for ruleset name searches
4. WHEN a search operation returns results, THE MCP Server SHALL include basic metadata (name, description, tags) for each matching ruleset

### Requirement 5

**User Story:** As a user of the MCP server, I want to list all available rulesets, so that I can see what guidelines are stored in the system.

#### Acceptance Criteria

1. THE MCP Server SHALL provide a tool to list all rulesets in the system
2. WHEN a list operation is requested, THE MCP Server SHALL return all ruleset names with their basic metadata
3. THE MCP Server SHALL return list results in a consistent, readable format
4. WHEN no rulesets exist, THE MCP Server SHALL return an empty list with an appropriate message

### Requirement 6

**User Story:** As a user of the MCP server, I want to update existing rulesets, so that I can refine and improve my guidelines over time.

#### Acceptance Criteria

1. THE MCP Server SHALL provide a tool to update existing rulesets
2. WHEN an update operation is requested, THE MCP Server SHALL allow modification of description, tags, and markdown content
3. WHEN a ruleset is updated, THE MCP Server SHALL update the last_modified field to the current timestamp
4. WHEN a ruleset is updated, THE MCP Server SHALL preserve the original created_at field value
5. IF an update operation is requested for a non-existent ruleset, THEN THE MCP Server SHALL return an appropriate error message

### Requirement 7

**User Story:** As a user of the MCP server, I want to delete rulesets I no longer need, so that I can keep my ruleset collection clean and relevant.

#### Acceptance Criteria

1. THE MCP Server SHALL provide a tool to delete rulesets by name
2. WHEN a delete operation is requested, THE MCP Server SHALL remove the ruleset and all associated metadata from Valkey
3. WHEN a delete operation completes successfully, THE MCP Server SHALL return a confirmation message
4. IF a delete operation is requested for a non-existent ruleset, THEN THE MCP Server SHALL return an error message containing a list of all existing ruleset names

### Requirement 8

**User Story:** As a system administrator, I want the MCP server to communicate via stdio transport, so that it can integrate with MCP clients using standard input/output.

#### Acceptance Criteria

1. THE MCP Server SHALL implement the stdio transport mechanism using the MCP-Go module
2. THE MCP Server SHALL read MCP protocol messages from standard input
3. THE MCP Server SHALL write MCP protocol responses to standard output
4. THE MCP Server SHALL handle MCP protocol initialization, capability negotiation, and message routing
5. THE MCP Server SHALL log errors and diagnostic information to standard error

### Requirement 9

**User Story:** As a system administrator, I want to configure the Valkey connection, so that the MCP server can connect to the appropriate Valkey instance.

#### Acceptance Criteria

1. THE MCP Server SHALL accept Valkey host configuration via environment variable
2. THE MCP Server SHALL accept Valkey port configuration via environment variable
3. THE MCP Server SHALL use default values for host (localhost) and port (6379) when environment variables are not provided
4. THE MCP Server SHALL establish connection to Valkey during startup
5. IF Valkey connection fails during startup, THEN THE MCP Server SHALL log an error message and exit with a non-zero status code

### Requirement 10

**User Story:** As a user, I want to run the MCP server using a Docker image, so that I can easily deploy and use the server without complex setup.

#### Acceptance Criteria

1. THE Docker Image SHALL be based on Alpine Linux
2. THE Docker Image SHALL include a bundled Valkey instance
3. THE Docker Image SHALL include the compiled MCP Server binary
4. WHEN the Docker container starts, THE Docker Image SHALL start the Valkey instance before starting the MCP Server
5. THE Docker Image SHALL expose the MCP Server via stdio transport for container interaction

### Requirement 11

**User Story:** As a project maintainer, I want automated CI/CD pipelines, so that code quality is maintained and releases are streamlined.

#### Acceptance Criteria

1. THE project SHALL include GitHub Actions workflow for Continuous Integration
2. WHEN a pull request is created or updated, THE CI workflow SHALL execute automatically
3. WHEN code is pushed to the main branch, THE CI workflow SHALL execute automatically
4. THE CI workflow SHALL compile the Go code, run tests, and perform linting
5. THE project SHALL include GitHub Actions workflow for Continuous Deployment that is manually triggered

### Requirement 12

**User Story:** As a project maintainer, I want semantic versioning and automated release notes, so that releases are consistent and well-documented.

#### Acceptance Criteria

1. THE CD workflow SHALL use semantic-release to determine version numbers
2. THE CD workflow SHALL generate changelogs based on conventional commit messages
3. THE CD workflow SHALL create GitHub releases with generated release notes
4. THE CD workflow SHALL build and publish Docker images with semantic version tags
5. THE project SHALL enforce conventional commit message format for all commits

### Requirement 13

**User Story:** As a developer contributing to the project, I want clear project structure and documentation, so that I can understand and contribute effectively.

#### Acceptance Criteria

1. THE project SHALL include a README.md file with setup instructions, usage examples, and architecture overview
2. THE project SHALL include a CONTRIBUTING.md file with guidelines for conventional commits and development workflow
3. THE project SHALL organize Go code following standard Go project layout conventions
4. THE project SHALL include example MCP client configuration for common AI editors
5. THE project SHALL include inline code documentation following Go documentation standards
