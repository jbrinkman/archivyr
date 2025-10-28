# Product Overview

MCP Ruleset Server is a centralized storage and management system for AI editor rulesets. It solves the problem of ruleset duplication and drift across multiple projects by providing a single source of truth for common rulesets.

## Core Purpose

Store and manage AI editor guidelines, rules, and steering documents in a centralized Valkey-backed system accessible via the Model Context Protocol (MCP).

## Key Features

- CRUD operations for rulesets (create, read, update, delete)
- Pattern-based search and listing capabilities
- Exact-match retrieval via MCP resources
- Metadata tracking (timestamps, tags, descriptions)
- Snake_case naming convention for rulesets
- Distributed as a self-contained Docker image with bundled Valkey instance

## Target Users

Developers using AI editors (Claude Desktop, Cursor, etc.) across multiple projects who need consistent, centralized ruleset management.
