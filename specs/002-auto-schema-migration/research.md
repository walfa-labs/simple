# Research: Auto Schema Migration with Graceful Fallback

**Feature**: 002-auto-schema-migration
**Date**: 2026-04-06

## Research Topics

### 1. PostgreSQL Schema Migration Pattern

**Decision**: Extend existing `CREATE TABLE IF NOT EXISTS` pattern in NewHandler

**Rationale**:
- Existing code already creates tables during initialization (postgres/handler.go:47-56)
- Pattern is simple, reliable, and idempotent
- No need for external migration tools (golang-migrate, sqlc migrations)
- Meets RMM Level 2 automated operations requirement

**Alternatives Considered**:
| Alternative | Rejected Because |
|-------------|------------------|
| golang-migrate library | Adds external dependency, complexity not needed for simple schema |
| Separate migration binary | Violates embedded schema assumption from spec |
| Versioned migrations | Spec states migrations are additive only, no rollback needed |

**Implementation Notes**:
- Extract schema creation from NewHandler into separate `MigrateSchema()` method
- Return specific errors for connection failure vs schema failure
- Use `CREATE TABLE IF NOT EXISTS` for idempotency
- Add index creation: `CREATE INDEX IF NOT EXISTS idx_records_created_at ON records(created_at DESC)`

---

### 2. MongoDB Schema Migration Pattern

**Decision**: Use collection creation with index creation in Connect phase

**Rationale**:
- MongoDB is schemaless, but indexes improve performance
- Collection created implicitly on first insert, but explicit creation preferred
- Index creation is idempotent with `createIndexes` command

**Alternatives Considered**:
| Alternative | Rejected Because |
|-------------|------------------|
| Explicit collection creation only | No benefit over implicit creation |
| Schema validation rules | Over-engineering for simple document model |

**Implementation Notes**:
- Create index on `created_at` field for sorting documents: `collection.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{"created_at", -1}}})`
- Handle "index already exists" as success (not error)
- Extract into `MigrateSchema()` method on Handler

---

### 3. Redis Schema Pattern

**Decision**: Initialize placeholder keys or skip (Redis has no schema)

**Rationale**:
- Redis is key-value store, no traditional schema
- Optional: initialize default keys for health check visibility
- Primary concern is connection validation, not schema migration

**Alternatives Considered**:
| Alternative | Rejected Because |
|-------------|------------------|
| Pre-populate sample data | Not required by spec, adds unnecessary complexity |
| No initialization at all | Connection test sufficient for Redis |

**Implementation Notes**:
- Current implementation (redis/handler.go:29-45) already validates connection with Ping
- No schema migration needed for Redis - just connection validation
- Add `MigrateSchema()` method that returns nil (no-op) for consistency across integrations

---

### 4. Integration Status Tracking

**Decision**: Create shared `IntegrationStatus` type in models package, tracker in integrations package

**Rationale**:
- Centralized status tracking enables health endpoint aggregation
- Enables runtime queries for integration state
- Matches existing Go patterns for service health

**Alternatives Considered**:
| Alternative | Rejected Because |
|-------------|------------------|
| Individual handler tracking | Status scattered across handlers, harder to aggregate |
| Global map variable | No encapsulation, thread safety concerns |

**Implementation Notes**:
- Define `IntegrationStatus` struct with: Name, Enabled, Error, LastCheckTime
- Create `StatusTracker` with methods: SetStatus(name, status), GetStatus(name), GetAllStatuses()
- Use sync.RWMutex for thread-safe access
- Register tracker in main.go, pass to handlers or expose globally

---

### 5. Graceful Degradation Pattern

**Decision**: Continue initializing other integrations after failure, disable failed integration, always start server

**Rationale**:
- Meets FR-004, FR-005: independent initialization, disable only failing integration
- Meets FR-007: `/hello` endpoint always available
- Pattern: attempt each integration → log failure → disable → continue
- Server starts regardless of integration status

**Alternatives Considered**:
| Alternative | Rejected Because |
|-------------|------------------|
| Fail fast on any error | Violates graceful degradation requirement |
| Retry loop on failure | Adds startup latency, spec says 30-second max |
| Circuit breaker pattern | Over-engineering for startup phase |

**Implementation Notes**:
- Each integration initialization wrapped in separate function
- Error logged with integration name and failure reason (FR-006)
- Failed integration not registered with router (routes unavailable)
- Status tracker updated for each integration
- Server starts after all initialization attempts complete

---

### 6. RMM Level 2 Compliance

**Decision**: Implement automated operations, graceful degradation, comprehensive logging, and basic monitoring

**Rationale**:
- RMM Level 2 requirements defined in FR-010 clarification
- Four key capabilities required:
  1. Automated operations: schema migration is automatic (no manual steps)
  2. Graceful degradation: failed integrations disabled, others continue
  3. Comprehensive logging: all failures logged with actionable messages
  4. Basic monitoring: health endpoints expose integration status

**Alternatives Considered**:
| Alternative | Rejected Because |
|-------------|------------------|
| Skip RMM requirements | Violates explicit requirement in spec |
| Add metrics/telemetry library | Not in scope, basic monitoring via health endpoints sufficient |

**Implementation Notes**:
- Structured logging format: `[INTEGRATION] migration failed: <error>`
- Health endpoint returns detailed status for each integration
- `/hello` endpoint as guaranteed liveness probe (always responds)