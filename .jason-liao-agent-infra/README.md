# Grok2API 受管同步

本目录是 Grok2API 工作流的唯一项目内入口。初版绑定 `lwj_dev`，仓库收敛由
Agent Infra registry 执行，项目工作流运行 `backend` 下的完整 Go 测试。

```bash
cd backend && go test ./...
```

服务配置、数据库恢复、启动和真实账户验收仍由平台 profile 管理，初版合同不会修改
配置或重启服务。失败时按 [errors.md](errors.md) 处理。
