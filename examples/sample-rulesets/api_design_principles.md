# API Design Principles

## Overview

This ruleset defines principles for designing RESTful APIs that are consistent, intuitive, and maintainable.

## Resource Naming

### Use Nouns, Not Verbs

- Resources should be nouns representing entities
- Use HTTP methods to indicate actions

```
# Good
GET    /users
POST   /users
GET    /users/123
PUT    /users/123
DELETE /users/123

# Bad
GET    /getUsers
POST   /createUser
GET    /getUserById/123
```

### Use Plural Nouns

- Use plural nouns for collections
- Keep naming consistent

```
# Good
/users
/products
/orders

# Bad
/user
/product
/order
```

### Use Hierarchical Structure

- Represent relationships with nested paths
- Keep nesting to 2-3 levels maximum

```
# Good
/users/123/orders
/users/123/orders/456

# Bad
/users/123/orders/456/items/789/details
```

## HTTP Methods

### Use Standard Methods

- `GET`: Retrieve resources (safe, idempotent)
- `POST`: Create new resources
- `PUT`: Update entire resources (idempotent)
- `PATCH`: Partial updates
- `DELETE`: Remove resources (idempotent)

### Method Semantics

```
GET    /users          # List all users
GET    /users/123      # Get specific user
POST   /users          # Create new user
PUT    /users/123      # Replace user
PATCH  /users/123      # Update user fields
DELETE /users/123      # Delete user
```

## Status Codes

### Use Appropriate Status Codes

- `200 OK`: Successful GET, PUT, PATCH, or DELETE
- `201 Created`: Successful POST
- `204 No Content`: Successful DELETE with no response body
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Authenticated but not authorized
- `404 Not Found`: Resource doesn't exist
- `409 Conflict`: Conflict with current state (e.g., duplicate)
- `422 Unprocessable Entity`: Validation errors
- `500 Internal Server Error`: Server error

### Consistent Error Responses

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      {
        "field": "email",
        "message": "Invalid email format"
      }
    ]
  }
}
```

## Request and Response Format

### Use JSON

- Use JSON as the default format
- Set `Content-Type: application/json`
- Accept `Accept: application/json`

### Consistent Field Naming

- Use `camelCase` for JSON fields
- Be consistent across all endpoints

```json
{
  "userId": 123,
  "firstName": "John",
  "lastName": "Doe",
  "emailAddress": "john@example.com",
  "createdAt": "2025-10-28T10:30:00Z"
}
```

### Use ISO 8601 for Dates

```json
{
  "createdAt": "2025-10-28T10:30:00Z",
  "updatedAt": "2025-10-28T15:45:00Z"
}
```

## Pagination

### Use Consistent Pagination

- Support pagination for list endpoints
- Use query parameters for pagination

```
GET /users?page=1&limit=20
GET /users?offset=0&limit=20
```

### Include Pagination Metadata

```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "totalPages": 8
  }
}
```

## Filtering and Sorting

### Use Query Parameters

```
GET /users?status=active
GET /users?role=admin&status=active
GET /users?sort=createdAt&order=desc
```

### Support Common Filters

- Equality: `?status=active`
- Comparison: `?age_gt=18` (greater than)
- Multiple values: `?status=active,pending`
- Search: `?q=john`

## Versioning

### Use URL Versioning

- Include version in the URL path
- Use major version numbers only

```
GET /v1/users
GET /v2/users
```

### Maintain Backward Compatibility

- Don't break existing clients
- Deprecate old versions gradually
- Provide migration guides

## Security

### Use HTTPS

- Always use HTTPS in production
- Redirect HTTP to HTTPS

### Authentication

- Use standard authentication (OAuth 2.0, JWT)
- Include authentication in headers

```
Authorization: Bearer <token>
```

### Rate Limiting

- Implement rate limiting
- Return rate limit headers

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1635724800
```

## Documentation

### Provide API Documentation

- Document all endpoints
- Include request/response examples
- Document error responses
- Use OpenAPI/Swagger specification

### Example Documentation

```yaml
/users:
  get:
    summary: List all users
    parameters:
      - name: page
        in: query
        schema:
          type: integer
    responses:
      200:
        description: Successful response
        content:
          application/json:
            schema:
              type: array
              items:
                $ref: '#/components/schemas/User'
```

## Best Practices

### Keep It Simple

- Design for the common case
- Don't over-engineer
- Start simple, evolve as needed

### Be Consistent

- Use consistent naming conventions
- Use consistent response formats
- Use consistent error handling

### Think About the Client

- Design from the client's perspective
- Minimize round trips
- Provide useful error messages

### Performance

- Use caching headers (ETag, Cache-Control)
- Support compression (gzip)
- Optimize database queries
- Consider pagination for large datasets

### Idempotency

- Make PUT and DELETE idempotent
- Consider idempotency keys for POST

```
POST /orders
Idempotency-Key: <unique-key>
```
