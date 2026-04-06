# API Contract: Simple Go API

**Version**: 1.0  
**Base URL**: `http://localhost:8080`  
**Documentation**: `/swagger/index.html`

## Overview

This document describes the API contract for the Simple Go API with optional PostgreSQL, Redis, and MongoDB integrations.

## Global Endpoints

### Health Check

**GET** `/health`

Returns the overall health status of the API.

**Response** (200 OK):
```json
{
  "status": "healthy"
}
```

---

## PostgreSQL Endpoints

*Enabled only when POSTGRES_HOST, POSTGRES_USER, and POSTGRES_DB are configured*

### Create Record

**POST** `/postgres/records`

**Request Body**:
```json
{
  "title": "string (required)"
}
```

**Response** (201 Created):
```json
{
  "id": 1,
  "title": "string",
  "created_at": "2026-04-06T12:00:00Z"
}
```

**Errors**:
- 400 Bad Request: Invalid input (missing title)
- 500 Internal Server Error: Database error

### Get All Records

**GET** `/postgres/records`

**Response** (200 OK):
```json
[
  {
    "id": 1,
    "title": "string",
    "created_at": "2026-04-06T12:00:00Z"
  }
]
```

### Get Record by ID

**GET** `/postgres/records/:id`

**Response** (200 OK):
```json
{
  "id": 1,
  "title": "string",
  "created_at": "2026-04-06T12:00:00Z"
}
```

**Errors**:
- 404 Not Found: Record does not exist

### Delete Record

**DELETE** `/postgres/records/:id`

**Response** (200 OK):
```json
{
  "message": "record deleted"
}
```

**Errors**:
- 404 Not Found: Record does not exist

### PostgreSQL Health

**GET** `/postgres/health`

**Response** (200 OK):
```json
{
  "status": "healthy",
  "database": "postgresql"
}
```

**Response** (503 Service Unavailable):
```json
{
  "status": "unhealthy",
  "error": "connection error details"
}
```

---

## Redis Endpoints

*Enabled only when REDIS_HOST is configured*

### Set Cache

**POST** `/redis/cache`

**Request Body**:
```json
{
  "key": "string (required)",
  "value": "string (required)",
  "expires_in": 3600
}
```

- `expires_in`: Optional TTL in seconds (0 or omitted = no expiration)

**Response** (200 OK):
```json
{
  "message": "cache set",
  "key": "string"
}
```

**Errors**:
- 400 Bad Request: Invalid input
- 500 Internal Server Error: Redis error

### Get Cache

**GET** `/redis/cache/:key`

**Response** (200 OK):
```json
{
  "key": "string",
  "value": "string",
  "expires_in": 3600
}
```

- `expires_in`: Remaining TTL in seconds (omitted if no expiration)

**Errors**:
- 404 Not Found: Key does not exist

### Delete Cache

**DELETE** `/redis/cache/:key`

**Response** (200 OK):
```json
{
  "message": "cache deleted",
  "key": "string"
}
```

**Errors**:
- 500 Internal Server Error: Redis error

### Redis Health

**GET** `/redis/health`

**Response** (200 OK):
```json
{
  "status": "healthy",
  "database": "redis",
  "info": "1234 bytes"
}
```

**Response** (503 Service Unavailable):
```json
{
  "status": "unhealthy",
  "error": "connection error details"
}
```

---

## MongoDB Endpoints

*Enabled only when MONGO_HOST and MONGO_DB are configured*

### Create Document

**POST** `/mongo/documents`

**Request Body**:
```json
{
  "title": "string (required)"
}
```

**Response** (201 Created):
```json
{
  "id": "507f1f77bcf86cd799439011",
  "title": "string",
  "created_at": "2026-04-06T12:00:00Z"
}
```

**Errors**:
- 400 Bad Request: Invalid input (missing title)
- 500 Internal Server Error: Database error

### Get All Documents

**GET** `/mongo/documents`

**Response** (200 OK):
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "title": "string",
    "created_at": "2026-04-06T12:00:00Z"
  }
]
```

### Get Document by ID

**GET** `/mongo/documents/:id`

**Response** (200 OK):
```json
{
  "id": "507f1f77bcf86cd799439011",
  "title": "string",
  "created_at": "2026-04-06T12:00:00Z"
}
```

**Errors**:
- 400 Bad Request: Invalid ID format
- 404 Not Found: Document does not exist

### Delete Document

**DELETE** `/mongo/documents/:id`

**Response** (200 OK):
```json
{
  "message": "document deleted"
}
```

**Errors**:
- 400 Bad Request: Invalid ID format
- 404 Not Found: Document does not exist

### MongoDB Health

**GET** `/mongo/health`

**Response** (200 OK):
```json
{
  "status": "healthy",
  "database": "mongodb"
}
```

**Response** (503 Service Unavailable):
```json
{
  "status": "unhealthy",
  "error": "connection error details"
}
```

---

## Common Error Responses

### 400 Bad Request
```json
{
  "error": "error message describing the validation failure"
}
```

### 404 Not Found
```json
{
  "error": "record not found" | "document not found" | "key not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "error message (may contain internal details)"
}
```

### 503 Service Unavailable
```json
{
  "status": "unhealthy",
  "error": "connection error details"
}
```

---

## Configuration

All database integrations are optional and enabled via environment variables:

### PostgreSQL
```env
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=user
POSTGRES_PASSWORD=pass
POSTGRES_DB=mydb
```

### Redis
```env
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```

### MongoDB
```env
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_USER=user
MONGO_PASSWORD=pass
MONGO_DB=mydb
```

---

## Observations

1. **Missing UPDATE Operations**: PostgreSQL and MongoDB endpoints lack PUT/PATCH for updates
2. **Error Exposure**: 500 errors may expose internal database error details
3. **No Pagination**: List endpoints return all records without pagination
4. **No Filtering**: No query parameters for filtering or sorting
5. **No Authentication**: All endpoints are publicly accessible
