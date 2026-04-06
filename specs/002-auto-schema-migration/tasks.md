# Tasks: Auto Schema Migration with Graceful Fallback

**Input**: Design documents from `/specs/002-auto-schema-migration/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Not explicitly requested in specification. Tests omitted per template guidance.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Project Type**: web-service (single Go application)
- **Repository root**: `/home/walfa/projekt/walfa/simple/`
- **Structure**: Existing Go project - extend current files and add new files

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Verify project state and establish shared models/tracking infrastructure

- [x] T001 Verify existing project structure matches plan.md (config/, handlers/, integrations/, models/, main.go)
- [x] T002 [P] Create IntegrationStatus model in models/integration.go
- [x] T003 [P] Create IntegrationSummary model in models/integration.go (for hello response)
- [x] T004 Create StatusTracker in integrations/status.go with SetStatus, GetStatus, GetAllStatuses methods
- [x] T005 Add sync.RWMutex to StatusTracker for thread-safe access in integrations/status.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T006 Create guaranteed /hello endpoint handler in handlers/hello.go
- [x] T007 Register /hello route in main.go (always registered, no conditions)
- [x] T008 Create /integrations/status endpoint handler in handlers/integration.go
- [x] T009 Create StatusTracker instance in main.go and pass to initialization functions
- [x] T010 Add integration name constants ("postgres", "redis", "mongo") to constants package

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Automatic Schema Migration on Startup (Priority: P1) 🎯 MVP

**Goal**: Automatically detect and create database schemas for all configured integrations on startup

**Independent Test**: Start application with valid database credentials, verify schemas created without manual intervention

### Implementation for User Story 1

- [x] T011 [P] [US1] Create PostgreSQL schema migration function in integrations/postgres/schema.go
- [x] T012 [P] [US1] Create MongoDB schema migration function in integrations/mongo/schema.go
- [x] T013 [P] [US1] Create Redis schema validation function in integrations/redis/schema.go (no-op, connection only)
- [x] T014 [US1] Add MigrateSchema method to postgres Handler in integrations/postgres/handler.go
- [x] T015 [US1] Add MigrateSchema method to mongo Handler in integrations/mongo/handler.go
- [x] T016 [US1] Add MigrateSchema method to redis Handler in integrations/redis/handler.go
- [x] T017 [US1] Extend postgres NewHandler to call MigrateSchema after connection in integrations/postgres/handler.go
- [x] T018 [US1] Extend mongo NewHandler to call MigrateSchema after connection in integrations/mongo/handler.go
- [x] T019 [US1] Extend redis NewHandler to validate connection (MigrateSchema) in integrations/redis/handler.go

**Checkpoint**: User Story 1 complete - schemas auto-migrate on startup when valid credentials provided

---

## Phase 4: User Story 2 - Graceful Integration Degradation (Priority: P2)

**Goal**: Failed integrations are automatically disabled while working integrations remain operational

**Independent Test**: Misconfigure one database while others valid, verify disabled integration logged, others work

### Implementation for User Story 2

- [x] T020 [US2] Modify main.go integration initialization to use defer/recover or error checking pattern
- [x] T021 [US2] Add per-integration error handling in main.go (disable on failure, continue to next)
- [x] T022 [US2] Update StatusTracker with failure reason when integration fails in main.go
- [x] T023 [US2] Add structured logging format "[INTEGRATION] migration failed: <error>" in main.go
- [x] T024 [P] [US2] Extend postgres HealthCheck to include migration_status field in integrations/postgres/handler.go
- [x] T025 [P] [US2] Extend redis HealthCheck to include migration_status field in integrations/redis/handler.go
- [x] T026 [P] [US2] Extend mongo HealthCheck to include migration_status field in integrations/mongo/handler.go
- [x] T027 [US2] Extend global HealthCheck handler to aggregate integration statuses in handlers/health.go

**Checkpoint**: User Story 2 complete - failed integrations gracefully disabled, others operational

---

## Phase 5: User Story 3 - Complete Fallback Mode (Priority: P3)

**Goal**: Application starts successfully even when all integrations fail, /hello always responds

**Independent Test**: Make all databases invalid, verify app starts and /hello responds

### Implementation for User Story 3

- [x] T028 [US3] Verify main.go server start happens after all integration attempts (not dependent on success)
- [x] T029 [US3] Add integration status summary to hello response in handlers/hello.go
- [x] T030 [US3] Register /integrations/status route in main.go
- [x] T031 [US3] Add error response for disabled integration endpoints (503 Service Unavailable) in handlers/integration.go
- [x] T032 [US3] Test scenario: no integrations configured → app starts, /hello works, all show "not_configured"

**Checkpoint**: User Story 3 complete - guaranteed availability regardless of integration status

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, Swagger updates, and final validation

- [x] T033 [P] Add Swagger comments for /hello endpoint in handlers/hello.go
- [x] T034 [P] Add Swagger comments for /integrations/status endpoint in handlers/integration.go
- [x] T035 Regenerate Swagger docs with `swag init -g main.go -o docs`
- [x] T036 Run quickstart.md validation scenarios (all integrations, one failed, all failed, none configured)
- [x] T037 Run linter with `golangci-lint run` (go vet used as fallback - passed)
- [x] T038 Run tests with `go test ./...`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion (T002-T005 models/tracker needed) - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational phase completion
- **User Story 2 (Phase 4)**: Depends on User Story 1 (needs migration methods to test failure handling)
- **User Story 3 (Phase 5)**: Depends on User Story 2 (needs graceful degradation pattern)
- **Polish (Phase 6)**: Depends on all user stories complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - Independent
- **User Story 2 (P2)**: Depends on US1 migration methods being implemented
- **User Story 3 (P3)**: Depends on US2 graceful degradation being implemented

### Within Each User Story

- Models before services
- Services before handlers
- Handler extensions before main.go integration
- Core implementation before cross-cutting concerns

### Parallel Opportunities

- Setup phase: T002, T003 can run in parallel (different files/models)
- US1: T011, T012, T013 can run in parallel (different integration schemas)
- US2: T024, T025, T026 can run in parallel (different health check files)
- Polish: T033, T034 can run in parallel (different swagger comments)

---

## Parallel Example: User Story 1

```bash
# Launch all schema migration files together:
Task: "Create PostgreSQL schema migration function in integrations/postgres/schema.go"
Task: "Create MongoDB schema migration function in integrations/mongo/schema.go"
Task: "Create Redis schema validation function in integrations/redis/schema.go"
```

## Parallel Example: User Story 2

```bash
# Launch all health check extensions together:
Task: "Extend postgres HealthCheck to include migration_status field"
Task: "Extend redis HealthCheck to include migration_status field"
Task: "Extend mongo HealthCheck to include migration_status field"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (models, status tracker)
2. Complete Phase 2: Foundational (hello endpoint, status registration)
3. Complete Phase 3: User Story 1 (schema migration)
4. **STOP and VALIDATE**: Test with valid database credentials, verify schemas auto-created
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test with valid DB → Deploy (MVP!)
3. Add User Story 2 → Test with misconfigured DB → Deploy
4. Add User Story 3 → Test with all DBs invalid → Deploy
5. Polish → Swagger docs, validation → Final release

---

## Summary

| Metric | Value |
|--------|-------|
| Total Tasks | 38 |
| Phase 1 (Setup) | 5 tasks |
| Phase 2 (Foundational) | 5 tasks |
| Phase 3 (US1 - MVP) | 9 tasks |
| Phase 4 (US2) | 8 tasks |
| Phase 5 (US3) | 5 tasks |
| Phase 6 (Polish) | 6 tasks |
| Parallel Opportunities | 8 tasks marked [P] |
| MVP Scope | Phases 1-3 (19 tasks) |

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story is independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Existing files: extend carefully, preserve current functionality
- New files: follow existing patterns from similar files