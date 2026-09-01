# Grok2API 受管同步错误

## GROK2API_SOURCE_VERIFY_FAILED

- 阶段：项目工作流。
- 含义：`backend` 的 `go test ./...` 未通过或超时。
- 检查：定位首个失败 package 和测试，确认当前 Go 版本与依赖解析状态。
- 处理：修复源码或测试后重跑完整命令。
- 停止条件：测试通过前，不进入服务配置、数据库迁移或重启。
