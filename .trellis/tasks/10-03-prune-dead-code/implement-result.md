# grok2api implementation evidence

Assignment: `grok-dead-code`. Authority: [DELETION_PLAN S2](/Users/jasonliao/Desktop/code/Artifacts/infra-dead-code-audit-20261003/DELETION_PLAN.md:30), [S8](/Users/jasonliao/Desktop/code/Artifacts/infra-dead-code-audit-20261003/DELETION_PLAN.md:108), [PRD](./prd.md), [design](./design.md), [implementation plan](./implement.md). Implementation DONE and ready for independent checking. Source-only, uncommitted; a separate checker follows.

## File map

All paths below are relative to `backend/`.

| File | Change |
| --- | --- |
| `internal/pkg/neterror/classify.go` | Remove Build idle sentinel alias and predicate; keep shared upstream errors/classifiers. |
| `internal/pkg/neterror/classify_test.go` | Preserve wrapped idle error, generic deadline rejection, cancellation and progress assertions on the upstream predicate. |
| `internal/application/gateway/selector.go` | Remove the unused `claimAccountSlot` wrapper; tracked coordination unchanged. |
| `internal/application/gateway/video.go` | Remove the old two-argument encoder wrapper. |
| `internal/application/gateway/video_test.go` | Use `encodeVideoInputFull` with explicit Generate, nil optional audio and empty optional video; preserve size and split/legacy input assertions. |
| `internal/infra/persistence/relational/client_key_repository.go` | Remove unused expired reservation count helper; cleanup/accounting untouched. |
| `internal/infra/provider/cli/normalize.go` | Remove the unused payload wrapper; metadata implementation unchanged. |
| `internal/infra/provider/cli/responses_tool_declarations.go` | Remove unused `dedupeSlice`; no replacement introduced. |
| `internal/infra/provider/cli/streamidle_test.go` | Retain context cause, read error and HTTP/2 cancellation assertions using the upstream sentinel. |
| `internal/infra/provider/conversation/chat_request.go` | Remove unused messages wrapper; reasoning replay implementation unchanged. |
| `internal/infra/provider/conversation/messages_request.go` | Remove unused request/message wrappers; reasoning replay and request dispatch unchanged. |
| `internal/infra/provider/console/normalize.go` | Remove the old normalizer wrapper and test-only `toolIdentity`. |
| `internal/infra/provider/console/console_test.go` | Call metadata normalizer with explicit nil defaults; retain metadata/reasoning tests; assert normalized tool types and function names directly. |
| `internal/infra/provider/web/account_identity.go` | Remove unused identity parser wrapper. |
| `internal/infra/provider/web/account_identity_test.go` | Test `sessionidentity.Parse` with unchanged missing/enveloped identity assertions. |
| `internal/infra/provider/web/quota.go` | Remove unused quota predicate wrapper. |
| `internal/infra/provider/web/protocol_test.go` | Call `account.IsWebImagineQuotaMode` in the existing catalog assertion. |

No retained production function body was edited. No imports became unused. DNS injection at `infra/provider/web/attachments.go` and maintenance `application/account/auto_clean.go` remain intact. No spec/documentation change is needed: ownership, APIs, current errors and runtime behavior remain the same.

## Evidence and verification status

Evidence directory: [/Users/jasonliao/Desktop/code/Artifacts/infra-dead-code-cleanup-20261003/grok/](/Users/jasonliao/Desktop/code/Artifacts/infra-dead-code-cleanup-20261003/grok/).

- `final.patch`, `patch-summary.json`: exact source patch and 17-file map. `gofmt-write`, `gofmt-check`, `diff-check`, `removed-symbols`, `preserved-seams`, `staged-files`, `final-status` each have exact argv/exit/stdout/stderr JSON and raw logs.
- Formatting, diff check and scope pass. Exact old symbol search in `backend/` returns no matches. The two required seams remain.
- `toolchain.{json,log}`: Go 1.27.1, darwin/arm64, CGO enabled. HOME, XDG config/cache, temp, GOPATH, build/module caches are isolated under the evidence directory; GOENV is disabled, local toolchain and readonly module flags explicit.
- `baseline-test.{json,log}`: pre-edit `go test ./...` failed at dependency setup with public official module route EOF; dependency-free packages passed. This failed baseline is retained.
- `final-test.{json,log}`: slow mirror download attempt explicitly stopped (Go child exit -15) when a faster verified official route became available; it is not a passing check.
- `dependency-prefetch.{json,log}`: public Go modules only via official GOPROXY and child-only HTTP/HTTPS proxy, exit 0 after 444.673s. No parent/system configuration change. `dependency-readiness` confirms the proxy-free cache is usable (exit 0).
- `related-test.{json,log}`: `go test -v` for neterror, gateway, relational, CLI, conversation, Console and Web passes (exit 0; 105.624s). Existing optional PostgreSQL/Redis integration cases skip because no external test credentials/addresses are injected; verbose logs retain the skip names and reasons.
- `full-test.{json,log}`: `go test ./...` passes (exit 0; 106.742s).
- `full-vet.{json,log}`: `go vet ./...` passes (exit 0; 22.91s).
- `full-build.{json,log}`: `go build -o <evidence>/grok2api-validation ./cmd/grok2api` passes (exit 0; 17.773s). The binary is retained only in the evidence directory and was not executed.
- `full-race.{json,log}`: `go test -race ./...` passes (exit 0; 466.862s), including tracked gateway coordination and relational tests. Original timeouts and assertions remain unchanged.

All validation after prefetch runs with `GOPROXY=off`, with HTTP/HTTPS proxy and external service/provider variables absent. This verifies source behavior and compilation; runtime and live consumers are outside the approved scope.

The final tracked patch contains exactly the 17 source/test paths above: 61 additions and 134 deletions. All 14 approved legacy symbols are absent from backend source/tests. The final source patch is frozen as `final.patch`; exact commands, exits and raw output are retained alongside it. No implementation blocker remains. External PostgreSQL/Redis integration and runtime/live provider behavior were not exercised, consistent with the source-only scope.

No staging, commits, push/fetch, runtime activation, service start, provider request, frontend/config/Swagger changes were performed. Existing task artifacts and source edits are preserved.
