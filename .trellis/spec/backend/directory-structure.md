# Backend Directory Structure

## Runtime Assembly

- `backend/cmd/grok2api/main.go` is intentionally thin: it delegates argument
  parsing and startup to `internal/cli`.
- `backend/internal/app/` owns composition, startup recovery, background
  lifecycle, and concrete dependency wiring.
- `backend/internal/transport/http/server.go` owns route groups and middleware
  order. Handlers are grouped by API area below `transport/http/<area>/`.

Do not construct databases, providers, or background workers inside handlers.
Add dependencies to the application assembly and pass the narrow service into
the handler, following `httpserver.Dependencies` in
`backend/internal/transport/http/server.go`.

## Layer Ownership

| Path | Ownership |
| --- | --- |
| `backend/internal/domain/` | Provider-neutral entities, enums, and domain rules |
| `backend/internal/application/` | Use cases, validation, orchestration, retries, and scheduling |
| `backend/internal/repository/` | Interfaces and shared query/error contracts |
| `backend/internal/infra/persistence/relational/` | GORM models, mappings, repositories, and schema upgrades |
| `backend/internal/infra/runtime/{memory,redis}/` | Ephemeral/distributed coordination implementations |
| `backend/internal/infra/provider/{cli,web,console}/` | Upstream-specific authentication and protocol adapters |
| `backend/internal/transport/http/` | Gin routes, auth boundaries, request/response DTOs, protocol adaptation |
| `backend/internal/pkg/` | Reusable mechanisms without feature ownership |
| `backend/internal/shared/response/` | Admin API success/error envelope only |

The application layer depends on repository interfaces, not GORM. For example,
`backend/internal/application/clientkey/service.go` accepts repository and
runtime interfaces, while `backend/internal/infra/persistence/relational/`
implements the relational side.

## Feature Placement

- Extend an existing feature package when the behavior belongs to that use
  case; account behavior belongs in `application/account`, not a generic util.
- Put provider capability declarations in
  `backend/internal/infra/provider/definition.go` and private wire formats in
  the provider package.
- Keep HTTP request/response conversion next to its handler. Domain and
  repository packages must not import Gin or public API DTOs.
- Put tests beside the package under test. Handler contract tests use
  `httptest`, while persistence tests use temporary databases.

Avoid a cross-layer “helpers” package. A helper that knows about HTTP status,
GORM rows, and domain accounts has no stable owner and breaks the dependency
direction documented in `backend/README.md`.

