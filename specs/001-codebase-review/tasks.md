# Tasks: Codebase Quality Review

**Input**: Design documents from `/specs/001-codebase-review/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Test tasks are included as this is a code quality improvement effort.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: Repository root with `config/`, `handlers/`, `integrations/`
- Paths shown are relative to repository root

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure for code improvements

- [x] T001 Create test directory structure: `mkdir -p tests/unit tests/integration tests/contract`
- [x] T002 [P] Create constants package for timeout values: `constants/timeouts.go`
- [x] T003 [P] Initialize test dependencies in `go.mod` (testify, mock packages if needed)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 Extract timeout constants from hardcoded values in `integrations/postgres/handler.go`, `integrations/redis/handler.go`, `integrations/mongo/handler.go`
- [x] T005 Create shared error handling utilities to prevent sensitive data exposure
- [x] T006 Add graceful shutdown signal handling in `main.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Code Quality Assessment (Priority: P1) 🎯 MVP

**Goal**: Fix code quality issues including unused types, inconsistent patterns, and missing tests

**Independent Test**: Run `go test ./...` and all tests pass; run `gofmt -l .` and no files need formatting

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T007 [P] [US1] Create unit test for config loading in `tests/unit/config_test.go`
- [ ] T008 [P] [US1] Create unit test for health handler in `tests/unit/handlers/health_test.go`
- [ ] T009 [P] [US1] Create integration test for PostgreSQL operations in `tests/integration/postgres_test.go`
- [ ] T010 [P] [US1] Create integration test for Redis operations in `tests/integration/redis_test.go`
- [ ] T011 [P] [US1] Create integration test for MongoDB operations in `tests/integration/mongo_test.go`

### Implementation for User Story 1

- [x] T012 [P] [US1] Remove unused `RecordInput` type from `integrations/mongo/handler.go` lines 31-33
- [x] T013 [P] [US1] Replace inline struct with `RecordInput` type in `integrations/postgres/handler.go` CreateRecord function
- [x] T014 [P] [US1] Replace inline struct with `DocumentInput` type in `integrations/mongo/handler.go` CreateDocument function
- [x] T015 [US1] Replace hardcoded timeouts with constants from `constants/timeouts.go` in `integrations/postgres/handler.go`
- [x] T016 [US1] Replace hardcoded timeouts with constants from `constants/timeouts.go` in `integrations/redis/handler.go`
- [x] T017 [US1] Replace hardcoded timeouts with constants from `constants/timeouts.go` in `integrations/mongo/handler.go`
- [x] T018 [US1] Add graceful shutdown with signal handling in `main.go`
- [x] T019 [US1] Run `go fmt ./...` to ensure all files are properly formatted
- [x] T020 [US1] Run `go vet ./...` and fix any issues

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Architecture Pattern Review (Priority: P2)

**Goal**: Complete CRUD operations by adding missing UPDATE endpoints for PostgreSQL and MongoDB

**Independent Test**: Test PUT/PATCH endpoints for PostgreSQL and MongoDB; verify records are updated correctly

### Tests for User Story 2

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T021 [P] [US2] Create contract test for PostgreSQL UPDATE endpoint in `tests/contract/postgres_update_test.go`
- [ ] T022 [P] [US2] Create contract test for MongoDB UPDATE endpoint in `tests/contract/mongo_update_test.go`
- [ ] T023 [P] [US2] Create integration test for PostgreSQL update flow in `tests/integration/postgres_update_test.go`
- [ ] T024 [P] [US2] Create integration test for MongoDB update flow in `tests/integration/mongo_update_test.go`

### Implementation for User Story 2

- [x] T025 [US2] Add `UpdateRecordInput` type in `integrations/postgres/handler.go`
- [x] T026 [US2] Implement UpdateRecord handler in `integrations/postgres/handler.go`
- [x] T027 [US2] Register UpdateRecord route in `integrations/postgres/handler.go` RegisterRoutes function
- [x] T028 [US2] Add Swagger documentation for PostgreSQL UPDATE endpoint
- [x] T029 [US2] Add `UpdateDocumentInput` type in `integrations/mongo/handler.go`
- [x] T030 [US2] Implement UpdateDocument handler in `integrations/mongo/handler.go`
- [x] T031 [US2] Register UpdateDocument route in `integrations/mongo/handler.go` RegisterRoutes function
- [x] T032 [US2] Add Swagger documentation for MongoDB UPDATE endpoint
- [x] T033 [US2] Regenerate Swagger documentation: `swag init -g main.go -o docs`

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Security and Error Handling Review (Priority: P2)

**Goal**: Improve error handling to prevent sensitive information exposure

**Independent Test**: Trigger database errors and verify no sensitive info (connection strings, internal paths) is returned to client

### Tests for User Story 3

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T034 [P] [US3] Create test for error message sanitization in `tests/unit/error_handler_test.go`
- [ ] T035 [P] [US3] Create integration test verifying no sensitive data in error responses in `tests/integration/security_test.go`

### Implementation for User Story 3

- [ ] T036 [US3] Create error sanitization function in `handlers/errors.go`
- [ ] T037 [P] [US3] Update error responses in `integrations/postgres/handler.go` to use sanitized errors
- [ ] T038 [P] [US3] Update error responses in `integrations/redis/handler.go` to use sanitized errors
- [ ] T039 [P] [US3] Update error responses in `integrations/mongo/handler.go` to use sanitized errors
- [ ] T040 [US3] Add error logging to stderr for internal debugging while returning safe errors to clients
- [ ] T041 [US3] Document SSL configuration requirements for production in `.env.example`

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T042 [P] Update `README.md` or create `TESTING.md` with test running instructions
- [ ] T043 [P] Add test coverage reporting configuration
- [ ] T044 Run full test suite: `go test ./...`
- [ ] T045 Run linter: `golangci-lint run` (or `go vet ./...` if golangci-lint not available)
- [ ] T046 Verify quickstart.md instructions still work after changes
- [ ] T047 Commit all changes with clear commit messages

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Builds on code patterns from US1
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - Can run in parallel with US2

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Code cleanup before adding new features
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel
- T012-T014 (removing unused types) can run in parallel
- T015-T017 (replacing timeouts) can run in parallel
- T037-T039 (updating error handling) can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Create unit test for config loading in tests/unit/config_test.go"
Task: "Create unit test for health handler in tests/unit/handlers/health_test.go"
Task: "Create integration test for PostgreSQL operations in tests/integration/postgres_test.go"
Task: "Create integration test for Redis operations in tests/integration/redis_test.go"
Task: "Create integration test for MongoDB operations in tests/integration/mongo_test.go"

# Launch all cleanup tasks for User Story 1 together:
Task: "Remove unused RecordInput type from integrations/mongo/handler.go"
Task: "Replace inline struct with RecordInput type in integrations/postgres/handler.go"
Task: "Replace inline struct with DocumentInput type in integrations/mongo/handler.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (code cleanup, tests, constants)
4. **STOP and VALIDATE**: Run tests, verify code quality improvements
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (code quality + tests)
   - Developer B: User Story 2 (UPDATE endpoints)
   - Developer C: User Story 3 (error handling)
3. Stories complete and integrate independently

---

## Task Summary

| Phase | Tasks | Story | Focus |
|-------|-------|-------|-------|
| Setup | T001-T003 | - | Infrastructure |
| Foundational | T004-T006 | - | Shared utilities |
| US1 | T007-T020 | P1 | Code quality, tests, cleanup |
| US2 | T021-T033 | P2 | UPDATE endpoints |
| US3 | T034-T041 | P2 | Error handling |
| Polish | T042-T047 | - | Documentation, final validation |

**Total Tasks**: 47
**High Priority**: 20 tasks (US1)
**Medium Priority**: 22 tasks (US2, US3)
**Low Priority**: 5 tasks (Setup, Foundational, Polish)

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
