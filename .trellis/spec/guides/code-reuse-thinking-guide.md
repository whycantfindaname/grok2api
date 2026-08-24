# Code Reuse and Ownership

## Search by Contract, Not Only Name

Before creating code, search the owning layer and its callers:

```bash
rg -n "symbol|route|jsonField|errorCode" backend frontend
rg --files backend/internal/<area> frontend/src/<area>
```

Inspect at least one implementation and its tests. Search for the serialized
field as well as the Go/TypeScript name because API and database names differ.

## Preferred Reuse Points

- Admin response envelopes: `backend/internal/shared/response/response.go`.
- Repository error translation: `infra/persistence/relational/errors.go`.
- Provider capability declaration: `infra/provider/definition.go`.
- Frontend authenticated JSON/SSE transport: `frontend/src/shared/api/client.ts`.
- Frontend runtime shape validation: `frontend/src/shared/api/decoder.ts`.
- Low-level UI primitives: `frontend/src/components/ui/`.
- Product-level loading/empty/error and table composition:
  `frontend/src/shared/components/`.
- Pure formatting/period/sort helpers: `frontend/src/shared/lib/`.

Extend these owners instead of reimplementing their contracts in a feature.
For example, a new Admin endpoint should not add a second fetch wrapper or
error-envelope parser, and a new handler should not locally translate GORM
errors.

## When Similar Code Must Stay Separate

Provider adapters under `infra/provider/{cli,web,console}` have different
authentication, model catalogs, quota sources, capabilities, and wire
protocols. Share a helper only when the invariant is genuinely identical and
provider-neutral. Do not merge behavior merely because two upstream JSON
objects currently look alike.

Likewise, Admin API responses and public OpenAI/Anthropic-compatible responses
have different contracts. The Admin `data`/`error` envelope must not leak into
public inference responses.

## Extraction Threshold

Extract when multiple real callers own the same stable rule and one owner can
name it clearly. Keep a small helper local when it has one feature owner. Do not
create `common`, `misc`, or cross-layer utility packages that mix HTTP, domain,
database, and provider knowledge.

Tests must follow the abstraction: when parsing/normalization becomes shared,
move its behavioral cases to the shared owner and retain caller tests for the
integration boundary.

