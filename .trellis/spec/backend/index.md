# Backend Development Guidelines

The backend is the Go module under `backend/`. Its dependency direction is
`transport/http -> application -> domain`; repository interfaces sit between
application logic and `infra` implementations. Start with `backend/README.md`
and `STRUCTURE.md`, then use the topic guide that matches the change.

## Pre-Development Checklist

- Identify the owning layer and inspect a neighboring implementation plus its
  `_test.go` file before editing.
- Keep provider-specific protocol details under
  `backend/internal/infra/provider/{cli,web,console}` rather than generic HTTP
  handlers or domain types.
- Decide whether state is relational persistence, runtime coordination
  (Memory/Redis), or process-local state; these are deliberately separate.
- For public API changes, check route registration, response shape, auth
  middleware, Swagger annotations, and the frontend decoder together.

## Guides

| Guide | Use it for |
| --- | --- |
| [Directory Structure](./directory-structure.md) | Package ownership and dependency direction |
| [Database Guidelines](./database-guidelines.md) | GORM repositories, schema upgrades, SQLite/PostgreSQL behavior |
| [Error Handling](./error-handling.md) | Error ownership, translation, and HTTP responses |
| [Logging Guidelines](./logging-guidelines.md) | `slog`, request correlation, and secret boundaries |
| [Quality Guidelines](./quality-guidelines.md) | Tests, CI-equivalent commands, and generated Swagger |

## Quality Check

Run from `backend/` for ordinary backend changes:

```bash
go test ./...
go vet ./...
go build ./cmd/grok2api
```

Use `go test -race ./...` for changes involving goroutines, caches, registries,
leases, or in-memory coordination. If public API annotations change, run
`make swagger` at the repository root and verify that only
`backend/docs/{docs.go,swagger.json,swagger.yaml}` changed as expected.

