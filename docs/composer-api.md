# Composer API Design Document

## Overview

The Composer API provides a structured interface for managing and interacting with software packages in the Open Registry system. This document outlines the design principles, endpoints, authentication mechanisms, and data models for the Composer API.

## API Design Principles

1. **RESTful Architecture**: The API follows REST principles with standard HTTP methods and status codes.
2. **JSON Format**: All request and response payloads use JSON format.
3. **Versioning**: API endpoints are versioned (e.g., `/api/v1/`) to ensure backward compatibility.
4. **Authentication**: JWT-based authentication with proper scope control.
5. **Rate Limiting**: Prevents abuse through appropriate rate limiting.
6. **Idempotency**: Operations that modify state support idempotency keys.

## Base URL

```
https://api.open-registry.example/api/v1
```

## Authentication

The API uses JSON Web Tokens (JWT) for authentication:

```
Authorization: Bearer <jwt_token>
```

## Endpoints

### Package Management

#### List Packages

```
GET /packages
```

Query parameters:
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20)
- `filter`: Filter by package properties
- `sort`: Sort by package properties

Response:
```json
{
  "data": [
    {
      "id": "package-id",
      "name": "package-name",
      "version": "1.0.0",
      "description": "Package description",
      "author": "Author Name",
      "created_at": "2025-06-01T12:00:00Z",
      "updated_at": "2025-06-01T12:00:00Z"
    }
  ],
  "meta": {
    "total": 100,
    "page": 1,
    "limit": 20
  }
}
```

#### Get Package Details

```
GET /packages/{package_id}
```

Response:
```json
{
  "id": "package-id",
  "name": "package-name",
  "version": "1.0.0",
  "description": "Package description",
  "dependencies": [
    {
      "name": "dependency-1",
      "version": "^2.0.0"
    }
  ],
  "author": "Author Name",
  "license": "MIT",
  "created_at": "2025-06-01T12:00:00Z",
  "updated_at": "2025-06-01T12:00:00Z"
}
```

#### Create Package

```
POST /packages
```

Request:
```json
{
  "name": "package-name",
  "version": "1.0.0",
  "description": "Package description",
  "dependencies": [
    {
      "name": "dependency-1",
      "version": "^2.0.0"
    }
  ],
  "license": "MIT"
}
```

Response:
```json
{
  "id": "package-id",
  "name": "package-name",
  "version": "1.0.0",
  "created_at": "2025-06-06T12:00:00Z"
}
```

#### Update Package

```
PUT /packages/{package_id}
```

Request:
```json
{
  "description": "Updated package description",
  "version": "1.0.1"
}
```

#### Delete Package

```
DELETE /packages/{package_id}
```

### Version Management

#### List Package Versions

```
GET /packages/{package_id}/versions
```

#### Get Specific Version

```
GET /packages/{package_id}/versions/{version}
```

#### Create Version

```
POST /packages/{package_id}/versions
```

### Dependency Management

```
GET /packages/{package_id}/dependencies
```

```
POST /packages/{package_id}/dependencies
```

### Publishing Workflow

```
POST /packages/{package_id}/publish
```

```
POST /packages/{package_id}/unpublish
```

## Error Handling

All errors follow a consistent format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {
      "field": "Additional error details"
    }
  }
}
```

Common HTTP status codes:
- `200 OK`: Request succeeded
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict (e.g., version already exists)
- `422 Unprocessable Entity`: Validation failed
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

## Rate Limiting

Rate limits are applied per API key and are included in response headers:

```
X-RateLimit-Limit: 5000
X-RateLimit-Remaining: 4999
X-RateLimit-Reset: 1654532461
```

## Pagination

All list endpoints support pagination with consistent response format:

```json
{
  "data": [...],
  "meta": {
    "total": 100,
    "page": 1,
    "limit": 20,
    "next_page": 2,
    "prev_page": null
  }
}
```

## Webhooks

The API supports webhooks for event notifications:

```
POST /webhooks
```

Request:
```json
{
  "url": "https://example.com/webhook",
  "events": ["package.published", "package.updated"],
  "secret": "webhook_secret"
}
```

## SDK Support

Official SDKs will be provided for:
- JavaScript/TypeScript
- Python
- Go
- Ruby

## Future Considerations

1. **GraphQL API**: Consider implementing a GraphQL API for more flexible querying.
2. **Package Metrics**: API endpoints for package download statistics and usage metrics.
3. **Team Collaboration**: API support for team-based package management.
4. **Advanced Security Features**: Vulnerability scanning and package signing APIs.

## Implementation Timeline

1. **Phase 1 (Q3 2025)**: Core package management endpoints
2. **Phase 2 (Q4 2025)**: Versioning and dependency management
3. **Phase 3 (Q1 2026)**: Webhooks and advanced features
