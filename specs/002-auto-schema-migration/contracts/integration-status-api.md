# Integration Status API Contract

**Endpoint**: `/integrations/status` (NEW)
**Purpose**: Expose operational status of all database integrations
**Feature**: 002-auto-schema-migration

## Request

```
GET /integrations/status
```

No parameters required.

## Response

**Success (200 OK)**:

```json
{
  "integrations": [
    {
      "name": "postgres",
      "enabled": true,
      "error": "",
      "last_check": "2026-04-06T12:00:00Z",
      "migration_status": "migrated"
    },
    {
      "name": "redis",
      "enabled": false,
      "error": "connection refused: localhost:6379",
      "last_check": "2026-04-06T12:00:00Z",
      "migration_status": "failed"
    },
    {
      "name": "mongo",
      "enabled": false,
      "error": "",
      "last_check": "2026-04-06T12:00:00Z",
      "migration_status": "not_configured"
    }
  ],
  "summary": {
    "total": 3,
    "enabled": 1,
    "disabled": 2
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| integrations | array | List of all integration statuses |
| integrations[].name | string | Integration identifier |
| integrations[].enabled | boolean | Whether operational |
| integrations[].error | string | Failure reason (empty if enabled or not configured) |
| integrations[].last_check | string | ISO 8601 timestamp of last status update |
| integrations[].migration_status | string | "migrated", "failed", or "not_configured" |
| summary | object | Aggregate counts |
| summary.total | number | Total integrations (always 3) |
| summary.enabled | number | Count of operational integrations |
| summary.disabled | number | Count of non-operational integrations |

## Migration Status Values

| Value | Meaning |
|-------|---------|
| migrated | Schema successfully migrated, integration operational |
| failed | Migration or connection failed, integration disabled |
| not_configured | Environment variables not set, integration skipped |

## Behavior

- Always returns 200 OK with current status
- Reflects status at startup time (not runtime re-check)
- No dependencies on specific integrations

## Error Handling

No error responses - endpoint always returns current status.

## Contract Tests

1. All integrations configured and working → all enabled=true, migration_status="migrated"
2. One integration misconfigured → that integration enabled=false, migration_status="failed", error populated
3. No integrations configured → all migration_status="not_configured", enabled=false