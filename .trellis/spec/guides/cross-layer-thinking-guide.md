# Cross-Layer Contract Changes

## Map the Actual Flow

For a management feature, trace:

```text
React component/form
  -> feature *-api.ts + runtime decoder
  -> /api/admin/v1 route and middleware
  -> HTTP request/response DTO
  -> application service
  -> repository/runtime/provider interface
  -> concrete infra implementation
```

For public inference, additionally trace provider capability selection,
provider-specific request/stream conversion, protocol-compatible errors, and
audit/media side effects. `backend/internal/transport/http/server.go` is the
route/auth map; `STRUCTURE.md` is the architecture and runtime-boundary map.

## Contract Checklist

When adding or changing a field, route, or behavior:

- Find every backend writer/reader and every frontend decoder/consumer.
- Keep request validation at HTTP/form boundaries and business validation in
  the application/domain owner.
- Update the runtime decoder; a TypeScript DTO alone does not validate JSON.
- Preserve stable Admin error codes and public protocol response shapes.
- Check auth group, request size, timeout, readiness, and concurrency
  middleware; route placement is a security and availability decision.
- If persisted, update schema models/indexes/migrations and test upgrade from an
  existing shape for both supported dialect assumptions.
- If it is coordination state, implement both Memory and Redis paths where the
  repository interface requires distributed behavior.
- If public API annotations changed, regenerate and review Swagger artifacts.
- Update README/STRUCTURE only when the public behavior or ownership map truly
  changed.

## Provider and State Boundaries

A public model name does not erase provider boundaries. Build, Web, and Console
accounts retain separate credentials, health, quota, concurrency, egress, and
capabilities. Generic gateway/handler code selects declared capabilities; it
must not infer private provider behavior from a model string.

Relational database state, Memory/Redis runtime coordination, and process-local
registries have different durability. Do not present process-local request
status as persistent audit data, or store high-frequency coordination in SQL
because a table is convenient.

## Verification and Evidence

Run the affected backend/frontend checks from their package directories. For a
cross-layer API change, the minimum is backend tests/vet/build plus frontend
test/lint/build. Swagger is an additional contract check when annotations
change.

Separate evidence explicitly:

- source: code, config template, tests, and docs agree;
- build: binaries/assets compile and generated contracts are current;
- runtime: a specific configured process passes `/healthz` and `/readyz`;
- live: a real authenticated provider or downstream consumer request succeeds.

Never claim the later stage from evidence that proves only an earlier one.
