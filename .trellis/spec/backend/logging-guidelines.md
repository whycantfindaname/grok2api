# Backend Logging Guidelines

## Logger and Shape

The service uses Go `log/slog`. `backend/internal/infra/observability/logger.go`
creates a JSON handler on stdout at info level. Pass `*slog.Logger` through
application assembly; use `slog.Default()` only as the existing defensive
fallback in `httpserver.New`.

Use a short stable event name followed by structured key/value fields. The
request reference is `middleware.AccessLog` in
`backend/internal/transport/http/middleware/request.go`:

```go
logger.Info("http_request", "request_id", requestID, "method", method,
    "path", routePattern, "status", status, "duration_ms", elapsed)
```

Log route patterns (`c.FullPath()`), not arbitrary user-provided URLs. Carry
the ingress request ID rather than inventing a second correlation field.

## Levels

- `Info`: lifecycle milestones, bounded background summaries, and one access
  event per request.
- `Warn`: recoverable degradation, rejected optional work, or fallback that
  changes behavior but lets the service continue.
- `Error`: an operation failed and needs operator action or makes a required
  component unavailable.
- Debug-level logging is not part of the default production logger; do not add
  high-volume per-token, per-chunk, or per-candidate info logs.

## Sensitive Data Boundary

Never log request/response bodies, authorization headers, API keys, OAuth/SSO
tokens, cookies, credential ciphertext, database DSNs, proxy credentials, or
quality-guard prompts/tokens. Access logging intentionally records only method,
matched path, status, duration, and request ID. Provider diagnostics must be
redacted before entering persistence or logs.

For logging changes, run focused tests for the owning package and `go test
./...`. If adding a potentially sensitive failure path, add an assertion like
`backend/internal/infra/persistence/relational/database_redaction_test.go` that
the secret is absent from rendered errors.

