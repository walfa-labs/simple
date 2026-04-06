# Quickstart: Auto Schema Migration with Graceful Fallback

**Feature**: 002-auto-schema-migration

## Prerequisites

- Go 1.26+ installed
- PostgreSQL, Redis, or MongoDB available (optional - any combination works)

## Quick Test Scenarios

### Scenario 1: All Integrations Working

1. Configure all databases in `.env`:

```env
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=user
POSTGRES_PASSWORD=pass
POSTGRES_DB=mydb

REDIS_HOST=localhost
REDIS_PORT=6379

MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_USER=user
MONGO_PASSWORD=pass
MONGO_DB=mydb
```

2. Start the application:

```bash
go run main.go
```

3. Verify migration:

```bash
# Check global health
curl http://localhost:8080/health

# Check integration status
curl http://localhost:8080/integrations/status

# Check PostgreSQL
curl http://localhost:8080/postgres/health

# Verify /hello works
curl http://localhost:8080/hello
```

Expected: All integrations show `migration_status: "migrated"`

---

### Scenario 2: One Integration Failed

1. Configure PostgreSQL and MongoDB correctly, but Redis with invalid host:

```env
REDIS_HOST=invalid-host  # Will fail
```

2. Start the application:

```bash
go run main.go
```

3. Check logs - should show:

```
[REDIS] migration failed: connection refused: invalid-host:6379
Redis integration disabled
```

4. Verify:

```bash
curl http://localhost:8080/integrations/status
```

Expected:
- postgres: enabled=true, migration_status="migrated"
- redis: enabled=false, migration_status="failed", error populated
- mongo: enabled=true, migration_status="migrated"

---

### Scenario 3: All Integrations Failed

1. Configure all databases with invalid credentials:

```env
POSTGRES_HOST=invalid
REDIS_HOST=invalid
MONGO_HOST=invalid
```

2. Start the application:

```bash
go run main.go
```

3. Application should start successfully (not crash)

4. Verify `/hello` works:

```bash
curl http://localhost:8080/hello
```

Expected: Returns 200 OK with message

5. Check integration status:

```bash
curl http://localhost:8080/integrations/status
```

Expected: All show migration_status="failed"

---

### Scenario 4: No Integrations Configured

1. Remove or empty all database config in `.env`

2. Start the application:

```bash
go run main.go
```

3. Application starts with only `/hello` endpoint

4. Verify:

```bash
curl http://localhost:8080/hello
curl http://localhost:8080/integrations/status
```

Expected: All show migration_status="not_configured"

---

## Key Endpoints

| Endpoint | Purpose | Always Available |
|----------|---------|------------------|
| `/hello` | Liveness probe | Yes |
| `/health` | Global health check | Yes |
| `/integrations/status` | Integration status details | Yes |
| `/postgres/health` | PostgreSQL health | No (only if enabled) |
| `/redis/health` | Redis health | No (only if enabled) |
| `/mongo/health` | MongoDB health | No (only if enabled) |

---

## Startup Behavior

1. Load configuration from `.env`
2. For each configured integration:
   - Attempt connection
   - Attempt schema migration
   - If failure: log error, disable integration, continue
   - If success: mark enabled, register routes
3. Start HTTP server (always starts)
4. `/hello` endpoint registered regardless of integration status

---

## Troubleshooting

| Symptom | Check |
|---------|-------|
| Integration disabled | Check logs for `[INTEGRATION] migration failed: <error>` |
| No routes for integration | Integration was disabled at startup - check `/integrations/status` |
| Application crashes | Bug - should never crash on integration failure |
| Slow startup | Migration timeout - check database connectivity |