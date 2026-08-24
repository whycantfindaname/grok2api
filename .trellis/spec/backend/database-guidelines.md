# Database Guidelines

## Storage Boundaries

Relational state uses GORM with SQLite for a typical single instance and
PostgreSQL for multi-instance deployments. Coordination state such as rate
limits, concurrency leases, sticky routing, locks, and invalidation events
belongs to `infra/runtime/memory` or `infra/runtime/redis`, not relational
tables. Process-local registries such as
`backend/internal/application/requeststatus/` are not durable ledgers.

## Repository Pattern

- Define the capability in `backend/internal/repository/` with `context.Context`
  as the first argument.
- Implement SQL/GORM behavior in
  `backend/internal/infra/persistence/relational/`.
- Convert rows to domain values through local mapping functions; do not expose
  GORM models above `infra`.
- Use `db.WithContext(ctx)` on queries and transactions. Preserve stable
  repository errors through `mapError` in
  `backend/internal/infra/persistence/relational/errors.go`.
- Keep bounded list operations explicit with page/limit, ID cursor, or batch
  size. The account repository contains concrete cursor-based examples.

## Schema and Upgrades

`backend/internal/infra/persistence/relational/schema.go` is the schema owner.
New persisted models must be added to `schemaModels`; required indexes belong
in `schemaIndexes` or a dedicated idempotent upgrade helper.

Upgrades must work for both dialects:

- PostgreSQL schema initialization takes an advisory transaction lock.
- SQLite migration may rebuild referenced tables, so the existing controlled
  foreign-key handling must be preserved.
- Data backfills must be bounded, restart-safe, and exclude already completed
  rows. `migrateStandaloneEgressProxyProfileBatch` is the reference pattern.
- Constraints or compatibility migrations need focused upgrade tests such as
  `schema_client_key_limits_upgrade_test.go` and
  `schema_media_job_upgrade_test.go`.

Do not add an unbounded startup scan or assume SQLite SQL works unchanged on
PostgreSQL. Do not silently discard a failed migration: wrap the operation with
context and return it so startup cannot claim readiness.

## Security and Verification

Never include DSNs, credentials, encrypted payloads, or secrets in errors or
logs. `database.go` redacts PostgreSQL URLs while preserving the error chain;
`database_redaction_test.go` is the regression example.

For persistence changes, run from `backend/`:

```bash
go test ./internal/infra/persistence/relational
go test ./...
go vet ./...
```

Run `go test -race ./...` when transactions interact with concurrent workers,
caches, or runtime invalidation.

