# grok2api execution plan

1. 主代理审阅 PRD/design/manifests 后启动任务；准备代理不得启动。复查 Git 身份、dirty 和当前目标内容，确认每个批准对象的消费者。
2. 分组清理 neterror 与 provider wrapper；迁移测试到当前接口，明确默认 metadata/reasoning/domain 参数。保留原行为断言。
3. 清理 gateway、relational 和无调用 helper；video 测试明确 Generate 与空可选参数，selector 测试保留 tracked 协调行为。
4. 从 `backend/` 运行相关包测试；最后运行 `go test ./...`、`go test -race ./...`、`go vet ./...`、`go build ./cmd/grok2api`。Go 格式仅作用于实际改动文件；以 `gofmt -l <changed Go files>` 无输出验收。
5. 定向 `rg` 查全部旧符号，检查有效测试断言、保留对象与 diff；根目录 `git diff --check`。失败或无法运行的检查报告 blocked/failed，不修改验收条件。
6. 主代理按整个清单验收；评估是否存在必要 spec 更新，无新知识可记录 no-op。保留 source-only/uncommitted 决定，收尾不进行 stage/commit/push/activation。

无需访问真实 provider、运行服务、fetch、更新 Swagger 或验证前端。回滚点按 neterror/provider、gateway/relational 两组；只撤销任务局部修改并保留并发内容。
