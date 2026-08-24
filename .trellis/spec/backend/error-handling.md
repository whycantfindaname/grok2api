# Backend Error Handling

## Ownership and Translation

Errors are translated at layer boundaries:

1. `infra/persistence/relational/errors.go` maps GORM not-found and duplicate
   errors to `repository.ErrNotFound` and `repository.ErrConflict`.
2. Application packages expose feature errors and wrap operational failures
   with `%w`, so callers can use `errors.Is` without losing context.
3. HTTP handlers map known application errors to stable status/code/message
   responses and hide unexpected internal details.

Follow the switches in `backend/internal/transport/http/account/handler.go` and
auth mappings in `backend/internal/transport/http/middleware/auth.go`. Validate
request shape at the handler boundary; validate business invariants in the
application service.

## HTTP Contracts

Admin APIs use `response.Success` and `response.Error` from
`backend/internal/shared/response/response.go`. Errors have a stable `code`, a
safe user-facing `message`, and the request ID when available. Public `/v1`
inference routes follow their OpenAI/Anthropic-compatible protocol shapes, so
do not force the Admin envelope onto them.

Use specific statuses:

- malformed input or unsupported filter: `400`;
- missing/expired authentication: `401`;
- resource conflict: `409` where the existing handler contract does so;
- runtime dependency unavailable or startup reconciliation: `503`;
- upstream/provider failure: the protocol-specific gateway mapping.

Do not return raw database, provider body, stack trace, DSN, token, cookie, or
encryption error to clients. Preserve the detailed cause for safe internal
logging only after redaction.

## Cancellation and Partial Work

Propagate request contexts into repositories and providers. Treat
`context.Canceled` and `context.DeadlineExceeded` deliberately in background
workers rather than logging them as unexplained failures. Batch operations
must report actual succeeded/failed/skipped counts and must not claim atomicity
unless the repository transaction provides it.

Tests should assert status and stable code, not only that “an error occurred.”
Use `httptest.NewRecorder` patterns from
`backend/internal/transport/http/server_test.go` and handler tests. Add a
redaction assertion whenever an error could contain credentials.

