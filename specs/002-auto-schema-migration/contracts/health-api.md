# Health API Contract (Extended)

**Endpoints**: `/health`, `/postgres/health`, `/redis/health`, `/mongo/health`
**Purpose**: Health check with migration status information
**Feature**: 002-auto-schema-migration

## Global Health Endpoint

### Request

```
GET /health
```

### Response (Extended)

**Success (200 OK)**:

```json
{
  "status": "healthy",
  "integrations": {
    "postgres": {
      "status": "healthy",
      "migration_status": "migrated"
    },
    "redis": {
      "status": "healthy",
      "migration_status": "migrated"
    },
    "mongo": {
      "status": "unhealthy",
      "migration_status": "failed",
      "error": "connection timeout"
    }
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| status | string | Overall status ("healthy" if at least one integration works) |
| integrations | object | Per-integration status |
| integrations[].status | string | "healthy" or "unhealthy" |
| integrations[].migration_status | string | "migrated", "failed", or "not_configured" |
| integrations[].error | string | Failure reason if unhealthy (optional) |

## Integration Health Endpoints (Extended)

### PostgreSQL Health

```
GET /postgres/health
```

**Success (200 OK)** - integration enabled and healthy:

```json
{
  "status": "healthy",
  "database": "postgresql",
  "migration_status": "migrated"
}
```

**Partial (200 OK)** - integration disabled:

```json
{
  "status": "disabled",
  "database": "postgresql",
  "migration_status": "failed",
  "error": "connection refused"
}
```

**Not Configured (200 OK)** - no environment variables:

```json
{
  "status": "not_configured",
  "database": "postgresql",
  "migration_status": "not_configured"
}
```

### Redis Health

```
GET /redis/health
```

**Success (200 OK)** - integration enabled:

```json
{
  "status": "healthy",
  "database": "redis",
  "migration_status": "migrated"
}
```

Same partial and not_configured responses as PostgreSQL.

### MongoDB Health

```
GET /mongo/health
```

**Success (200 OK)** - integration enabled:

```json
{
  "status": "healthy",
  "database": "mongodb",
  "migration_status": "migrated"
}
```

Same partial and not_configured responses as PostgreSQL.

## Status Values

| Status | Meaning | HTTP Code |
|--------|---------|-----------|
| healthy | Integration operational, ping succeeded | 200 |
| unhealthy | Integration enabled but ping failed | 503 |
| disabled | Integration migration/connection failed at startup | 200 |
| not_configured | Environment variables not set | 200 |

## Behavior Changes

- Existing endpoints return 200 for disabled/not_configured (previously would not exist)
- Health check still performs ping if integration is enabled
- Migration status added to all responses

## Contract Tests

1. Integration healthy → status="healthy", migration_status="migrated"
2. Integration disabled at startup → status="disabled", migration_status="failed", error populated
3. Integration not configured → status="not_configured", migration_status="not_configured"
4. Integration enabled but runtime ping fails → status="unhealthy", HTTP 503