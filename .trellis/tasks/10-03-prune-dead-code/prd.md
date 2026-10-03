# Remove approved grok2api wrappers and helpers

## Goal

删除已批准的 S2（grok2api 部分）与 S8 旧内部实现，保留当前生产行为和有用回归断言。

## Authority and background

[已批准删除计划](/Users/jasonliao/Desktop/code/Artifacts/infra-dead-code-audit-20261003/DELETION_PLAN.md) 的 S2、S8 是对象和语义边界的完整依据。用户的“可以请继续”批准计划，后续“确认”批准任务创建与启动；本任务只物化原决定。交付仅 source、未提交，不进行发布或激活。当前基线见 [baseline.json](/Users/jasonliao/Desktop/code/Artifacts/infra-dead-code-cleanup-20261003/task-preparation/baseline.json)。

## Requirements

- R1 删除批准对象及仅服务旧实现的导入；测试转向计划中的当前接口，保留默认/空 metadata、reasoning、video 参数与 idle timeout 语义。
- R2 产品文件范围仅以下路径及同包直接覆盖它们的现有 `_test.go`；测试文件在编辑前按符号引用确认，不扩展产品功能：
  - `backend/internal/pkg/neterror/classify.go`: `ErrBuildStreamIdleTimeout`, `IsBuildStreamIdleTimeout`。
  - `backend/internal/application/gateway/selector.go`: `claimAccountSlot`；`video.go`: `encodeVideoInput`。
  - `backend/internal/infra/persistence/relational/client_key_repository.go`: `expiredBillingReservationCount`。
  - `backend/internal/infra/provider/cli/normalize.go`: `normalizeBuildRequestPayload`；`responses_tool_declarations.go`: `dedupeSlice`。
  - `backend/internal/infra/provider/conversation/chat_request.go`: `convertChatMessages`；`messages_request.go`: `convertMessagesRequest`, `convertAnthropicMessages`。
  - `backend/internal/infra/provider/console/normalize.go`: `normalizeRequest`, `toolIdentity`。
  - `backend/internal/infra/provider/web/account_identity.go`: `parseAccountIdentity`；`quota.go`: `isImagineQuotaMode`。
- R3 保留 `validateRemoteImageURLWithResolver`、`autoCleanConfig`、泛型函数、当前生产接口及配置兼容字段。回读当前内容；若正常生产消费者出现，报告证据并交给主代理裁决。
- R4 不修改前端、部署、配置、凭据、API、生成 Swagger 或其他候选项；不 stage/commit/push/fetch、启动服务或进行真实 provider 请求。

## Acceptance Criteria

- [ ] R1/R2 全部批准旧符号移除，活跃接口与实际行为断言保留；定向引用复查无旧符号残留。
- [ ] R3 活跃测试接缝、取消/跟踪/metadata/reasoning/video/domain 行为保留。
- [ ] 后端相关包测试及全量 `go test ./...`、`go vet ./...`、`go build ./cmd/grok2api`、`go test -race ./...` 通过；Go 文件符合 gofmt，`git diff --check` 通过。
- [ ] 仅授权源码和测试变更，未提交、未应用；并发和已有修改保留。检查失败/缺失如实报告。

## Settled decisions

任务属于复杂、多包清理，配套 `design.md` 与 `implement.md`。本准备阶段维持 planning，由主代理审阅产物并执行已授权的启动。
