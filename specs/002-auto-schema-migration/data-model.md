# Data Model: Auto Schema Migration with Graceful Fallback

**Feature**: 002-auto-schema-migration
**Date**: 2026-04-06

## Entities

### IntegrationStatus

Tracks the operational state of each database integration.

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| Name | string | Integration identifier ("postgres", "redis", "mongo") | Required, one of three values |
| Enabled | bool | Whether integration is operational | Default false, set true after successful migration |
| Error | string | Failure reason if disabled | Empty if enabled, descriptive message if failed |
| LastCheckTime | timestamp | When status was last updated | Set on every status change |
| MigrationVersion | string | Schema version migrated (optional) | Empty for Redis (no schema) |

**State Transitions**:

```
Initial (Enabled=false, Error="")
    → [Migration Success] → Operational (Enabled=true, Error="")
    → [Migration Failure] → Disabled (Enabled=false, Error="<reason>")
    → [Connection Failure] → Disabled (Enabled=false, Error="connection failed: <reason>")
```

**Relationships**: None - each integration status is independent.

---

### SchemaMigrationRecord

Tracks which schemas have been successfully migrated (internal tracking).

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| IntegrationName | string | Which integration this applies to | Required |
| Version | string | Schema version identifier | Default "v1" |
| MigratedAt | timestamp | When migration completed | Set on success |
| TablesCreated | []string | List of tables/collections created | PostgreSQL/MongoDB only |

**Notes**:
- Not exposed via API (internal tracking only)
- Redis has empty TablesCreated (no schema)
- Used for debugging startup sequence

---

### HealthCheckResponse

Response format for health check endpoints.

| Field | Type | Description |
|-------|------|-------------|
| Status | string | "healthy" or "unhealthy" |
| Database | string | Integration name (from existing endpoints) |
| MigrationStatus | string | NEW: "migrated", "not_configured", "failed" |
| Error | string | NEW: failure reason if failed |

**Existing Endpoints to Extend**:
- `/postgres/health` → add MigrationStatus, Error fields
- `/redis/health` → add MigrationStatus, Error fields
- `/mongo/health` → add MigrationStatus, Error fields
- `/health` → aggregate all integration statuses

---

### HelloResponse

Response for the guaranteed `/hello` endpoint.

| Field | Type | Description |
|-------|------|-------------|
| Message | string | Simple greeting message |
| Timestamp | timestamp | Server time |
| Integrations | map | Summary of integration statuses (optional) |

**Validation**: No validation - always returns successful response.

---

## Storage

**In-Memory Storage**: IntegrationStatus and SchemaMigrationRecord stored in application memory (no persistence required).

- StatusTracker holds current integration states
- Not persisted to disk (application restart resets status)
- Acceptable per spec assumptions: "Schema migrations are additive only"

## Validation Rules

| Entity | Rule | Enforcement |
|--------|------|-------------|
| IntegrationStatus.Name | Must be one of: postgres, redis, mongo | Enum validation in code |
| IntegrationStatus.Error | Max 500 characters | Truncate if exceeded |
| HealthCheckResponse.Status | Must be "healthy" or "unhealthy" | Const values |
| HelloResponse | No validation required | Always succeeds |