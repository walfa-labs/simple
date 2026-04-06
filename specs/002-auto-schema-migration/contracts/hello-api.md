# Hello API Contract

**Endpoint**: `/hello`
**Purpose**: Guaranteed liveness probe - always available regardless of integration status
**Feature**: 002-auto-schema-migration

## Request

```
GET /hello
```

No parameters required.

## Response

**Success (200 OK)**:

```json
{
  "message": "Hello from Simple API",
  "timestamp": "2026-04-06T12:00:00Z",
  "integrations": {
    "postgres": "enabled",
    "redis": "disabled",
    "mongo": "enabled"
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| message | string | Static greeting message |
| timestamp | string | ISO 8601 timestamp of server time |
| integrations | object | Summary of integration statuses (optional enhancement) |

## Behavior

- **Always succeeds**: Endpoint must respond even when all integrations failed
- **No dependencies**: Does not require any database connection
- **Fast response**: Should respond within milliseconds

## Error Handling

No error responses - endpoint always returns 200 OK.

## Contract Tests

1. Call `/hello` with all integrations working → expect 200
2. Call `/hello` with no integrations configured → expect 200
3. Call `/hello` after all integrations failed → expect 200