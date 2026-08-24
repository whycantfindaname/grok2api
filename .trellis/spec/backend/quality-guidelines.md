# Backend Quality Guidelines

## Required Patterns

- Keep dependencies explicit through constructors or assembly structs; avoid
  package-global mutable services.
- Pass `context.Context` through blocking repository/provider work.
- Use UTC for persisted timestamps and externally compared deadlines. GORM's
  `NowFunc` in `relational/database.go` is the baseline.
- Protect concurrent process state with an owned mutex/atomic abstraction and
  test lifecycle behavior. `application/requeststatus/registry.go` is a small
  reference; routing, leases, and caches require the same ownership clarity.
- Keep auth boundaries visible in `transport/http/server.go`: public, admin,
  internal quality guard, and client-key routes are not interchangeable.

## Tests

Place tests beside code and prefer behavior-level assertions:

- table-driven unit tests for parsers, validation, mappings, and domain rules;
- application tests with narrow fakes for orchestration and error mapping;
- `httptest` for middleware, route protection, status, headers, and envelopes;
- temporary SQLite files for repository and migration behavior;
- explicit PostgreSQL integration tests only where dialect behavior matters.

`backend/internal/transport/http/server_test.go` demonstrates route/auth and
SPA fallback contracts. Persistence upgrade tests under
`infra/persistence/relational/schema_*_upgrade_test.go` demonstrate migrations
from an old shape, not just a fresh schema.

## Forbidden Shortcuts

- Do not bypass repository interfaces with GORM from application/domain code.
- Do not expose provider-specific payloads through shared domain contracts.
- Do not weaken auth, readiness, request-size, timeout, or concurrency
  middleware for a local feature.
- Do not treat `/healthz`, `/readyz`, or a successful build as proof that a
  configured provider or downstream consumer is live.
- Do not hand-edit generated Swagger files.
- Do not put runtime `config.yaml`, databases, WAL/SHM files, credentials,
  logs, media, or cache output in Git.

## Verification Matrix

```bash
cd backend
go test ./...                    # all backend behavior
go vet ./...                     # static checks used in CI
go build ./cmd/grok2api          # production entrypoint compiles
go test -race ./...              # concurrency-sensitive changes
```

For API annotation changes, additionally run `make swagger` from the root and
compare the three tracked files under `backend/docs/`. The authoritative CI
sequence is `.github/workflows/ghcr-image.yml`.
