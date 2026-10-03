# grok2api cleanup design

## Boundaries and contracts

对象清单与替代接口遵循 [批准计划 S2/S8](/Users/jasonliao/Desktop/code/Artifacts/infra-dead-code-audit-20261003/DELETION_PLAN.md)。仅删除内部冗余层，保持 transport → application → domain 与 provider/relational 的现有所有权。

测试直接调用当前 `ErrUpstreamStreamIdleTimeout`/predicate、Tracked、WithMetadata、WithReasoningReplay、`sessionidentity.Parse`、`account.IsWebImagineQuotaMode` 和 `encodeVideoInputFull`。具体参数取当前签名并显式表达旧包装器的默认值；不得猜测等价签名。`toolIdentity` 的有用断言落在当前规范化输出，纯旧 helper 自测删除。计数与去重 helper 无调用，不引入替代算法。

## Compatibility and operational effects

当前生产协议、API、数据、配置、错误处理和并发协调行为保持。活跃的 DNS 注入/maintenance 接缝继续存在。若当前源码推翻无消费者前提，停止该对象，提供引用证据，其他独立对象可继续。

source-only 清理无 rollout/服务动作。普通 Git 历史可用于后续恢复；如回滚本次未提交修改，仅手工撤销任务拥有的局部差异，保留并发内容，不使用全文件或 Git reset 恢复。
