# grok2api 结构与维护边界

本文是 grok2api 在 Jason Liao Infra 工作区中的结构说明。内容以当前
lwj_dev checkout 的源码、配置模板、构建文件、测试和 Git 状态为准；它说明
代码边界，不是部署清单，也不证明任何消费者或服务当前已经上线。

本文核验时的工作分支为 `lwj_dev`；具体提交和远端同步状态始终以 Git 回读为准。

## 1. Infra 角色、消费者与仓库身份

在 Infra companion manifest 中，本仓库的 component id 是 grok2api，工作区
路径是 Tools/grok2api。它承担 Grok 多账号 API 网关和管理端源码的维护，属于
Infra 的 companion source checkout；它不是 smart-search-cli 或
grok2api-service 的源码目录。

当前登记的消费者是：

| 消费者 | 关系边界 |
| --- | --- |
| smart-search-cli | 可把本仓库提供的网关作为搜索/推理后端；是否已经配置、导入账号并完成真实 consumer probe，由 Infra 的激活门禁单独证明。 |
| grok2api-service | 服务编排层登记的下游消费者；本仓库不拥有该服务的启动器、部署状态或运行时凭据。 |

仓库远端和分支关系如下：

| 身份 | 当前 Git 配置 |
| --- | --- |
| 个人 fork | fork → https://github.com/whycantfindaname/grok2api.git |
| 上游 | origin → https://github.com/chenyme/grok2api.git |
| 工作分支 | lwj_dev，当前跟踪个人 fork 的 lwj_dev |
| 个人分支基线 | 本地 main/fork/main，当前为 9fadb5b6… |

fork/main 是当前 checkout 可直接比较的个人 fork 基线；本仓库没有在本次
任务中 fetch，因此不能把未读取的远端 origin/main 当作当前上游快照。服务
是否运行、账号是否导入、消费者是否能成功请求，必须分别由对应的运行态检查和
真实 consumer probe 验证；README、Dockerfile、端口配置或 checkout 存在本身都
不能替代这些证据。

## 2. 根目录与 tracked 顶层结构

当前 tracked 的顶层目录只有以下五个；.git/ 是 Git 元数据，不属于项目源码。

    .
    ├── .github/                    Issue 模板与 CI 工作流
    ├── backend/                    Go 后端、Provider、持久化与 Swagger
    ├── docker/                     容器启动入口
    ├── frontend/                   React/Vite 管理端源码和静态资源
    └── tools/                      可选运维工具；当前为出口质量守护程序

当前 tracked 的重要根文件如下：

| 文件 | 维护边界 |
| --- | --- |
| .dockerignore | Docker build context 的排除规则。 |
| .gitignore | 本地配置、凭据、运行态、缓存、构建物和临时文件的 Git 边界。 |
| Dockerfile | 前端构建、后端编译和最终 Alpine 运行镜像的多阶段构建定义。 |
| LICENSE | 上游项目许可证文本；除非许可证要求，不在此文档中改写。 |
| Makefile | 根级 run 和 swagger 入口。 |
| README.md、README.zh-CN.md | 面向使用者的能力、部署、API 和 Provider 说明；用户可见行为变化时同步更新。 |
| VERSION | 镜像和管理端显示使用的版本文件。 |
| config.example.yaml | 脱敏的启动配置模板和字段边界，不放真实配置。 |
| docker-compose.yml | 主服务、可选 FlareSolverr/WARP 注释示例、quality-guard profile 和持久化卷。默认镜像名仍是 ghcr.io/chenyme/grok2api:latest，不等于个人 fork 已发布或已被使用。 |
| STRUCTURE.md | 本文；只记录源码结构、所有权和验证边界，不保存账号、密钥或运行日志。 |

.DS_Store 被 .gitignore 忽略。本次不删除、不纳入结构文档所描述的源码资产，
也不以它作为项目状态证据。

.github/ 下的 tracked 内容包括：

- workflows/ghcr-image.yml：后端测试、go vet、Swagger 生成物校验、前端 lint/build，
  以及按架构构建/发布 GHCR 镜像。
- workflows/codeql.yml：Actions、Go 和 JavaScript/TypeScript 的 CodeQL 分析。
- workflows/stale.yml：定时处理长期无活动 Issue/PR。
- ISSUE_TEMPLATE/：Bug、文档和功能请求模板；模板明确要求脱敏凭据和环境信息。

## 3. Backend 结构

后端是独立 Go module（见 [backend/go.mod](backend/go.mod)，当前 Go 版本为 1.26）。进程入口是
[backend/cmd/grok2api/main.go](backend/cmd/grok2api/main.go)，启动参数由 internal/cli 解析，再由 internal/app
完成配置、数据库、运行态、Provider、应用服务和 HTTP 路由装配。

    backend/
    ├── cmd/grok2api/main.go       进程入口；只调用 CLI runner
    ├── docs/                       已生成并纳入 Git 的 Swagger Go/JSON/YAML
    ├── internal/
    │   ├── app/                    应用装配、启动恢复、拓扑和生命周期
    │   ├── domain/                 账号、模型、推理、媒体、审计、密钥等领域规则
    │   ├── application/            账号同步、认证、网关、审计、媒体、模型、出口等用例
    │   ├── infra/                  配置、数据库、运行态、Provider、出口、安全和观测实现
    │   ├── repository/              持久化、运行态、事件和限流等接口
    │   ├── transport/http/          Gin 路由、鉴权、中间件和协议适配
    │   ├── pkg/                    批处理、网络错误、缓存、性能、推理回放等通用机制
    │   └── shared/response/         共享响应结构
    └── README.md                   后端运行、代码层次和验证入口

层次依赖的维护方向是 Transport → Application → Domain；数据库、运行态和
Provider 等具体实现通过接口从 infra 接入，不应把 HTTP、数据库或某个上游
协议细节下沉到领域层。

### 3.1 HTTP 和消费者接口

路由集中在 [backend/internal/transport/http/server.go](backend/internal/transport/http/server.go)：

- /healthz：存活响应。
- /readyz：分层启动恢复和运行态就绪响应；不等于某个外部消费者验证成功。
- /v1/*：客户端 API。推理 handler 覆盖 Responses、Chat Completions、Anthropic
  Messages、Images 和 Videos；还包括 stored response、compact、媒体读取和当前
  fork 增加的 /v1/request-status/:requestId 查询。
- /api/admin/v1/*：管理员登录、账号/模型/Client Key、审计、Dashboard、媒体、设置、
  出口和系统更新接口，受管理员鉴权保护。
- /api/internal/v1/quality-guard/*：仅在启用质量守护并带内部凭据时注册，供
  sidecar 使用，不是公共 API。
- /swagger/*：只有 server.swaggerEnabled: true 时注册；生产模板默认关闭。
- 前端静态文件：由 frontend.staticPath 指向的构建目录由同一 Go 服务托管。

公共推理请求需要客户端 Key；管理员接口使用管理员会话。当前分支还在进程内按
client key + request id 短暂保留同步文本请求状态，代码位于
internal/application/requeststatus 和 internal/transport/http/middleware/request_status.go。
该状态不是跨进程数据库账本，不能当作持久化审计或跨实例运行态。

### 3.2 Provider 边界

Provider 通过 [internal/infra/provider/definition.go](backend/internal/infra/provider/definition.go) 声明能力，生产适配器由
internal/infra/provider/{cli,web,console} 实现；internal/application/gateway
依据能力、模型路由、账号健康/额度/并发和出口策略调度，不在通用 handler 中拼接
Provider 私有请求。

| Provider | 当前代码中声明的边界 |
| --- | --- |
| Grok Build（cli） | OAuth/Device OAuth，远程模型目录，Billing 额度；Responses、Chat、Messages、compact、stored responses 和视频。 |
| Grok Web（web） | SSO 导入，静态模型目录按等级过滤，远程额度窗口；Responses、Chat、Messages、stored responses、图片、图片编辑和视频。Web 的 Gateway/chat/responses stream 还负责上游搜索工具、引用和 SSE/协议转换。 |
| Grok Console（console） | SSO 导入，静态模型目录；无状态 Responses/Chat/Messages，以及图片、图片编辑和视频；与 Web 共用浏览器/出口挑战处理面。 |

三个 Provider 各自维护凭据、额度、健康、冷却、并发和模型能力。故障切换只在
选定 Provider 的候选账号内进行，不能因为公开模型同名就混用另一个 Provider 的
账号状态。

### 3.3 internal 关键实现分区

- application/account、accountsync：账号导入/导出、OAuth/SSO、额度和模型同步、
  凭据刷新、账号转换、清理与 Provider 关联。
- application/gateway：模型路由、选号、粘滞会话、并发/额度门禁、重试、媒体任务、
  请求结束结算和质量探测。
- application/clientkey、adminauth：客户端 Key 生成/鉴权/限制和管理员会话。
- application/audit、dashboard、media、model、settings、quotarecovery、
  invalidation、updatecheck：对应管理与后台用例。
- infra/persistence/relational：SQLite/PostgreSQL、迁移/初始化和 Repository 实现。
- infra/runtime/memory、infra/runtime/redis：限流、并发租约、粘滞路由、分布式锁、
  事件通知、额度恢复和推理回放的运行态实现。
- infra/egress：按 Provider/资源作用域管理直连、HTTP/SOCKS/Resin、订阅、探测、
  固定代理、回退和 FlareSolverr clearance。
- infra/media、infra/security、infra/observability、infra/qualityguard：本地
  媒体对象存储、凭据加密/密码/令牌、日志以及质量守护 bootstrap。
- transport/http/{account,adminauth,audit,clientkey,dashboard,egress,inference,media,model,settings,system}：
  管理和公共协议的 handler；middleware 统一处理请求 ID、大小/超时、鉴权、并发、
  请求状态和安全响应头。

## 4. Frontend 结构

frontend/ 是 React 19 + TypeScript + Vite 项目，锁定包管理器为 pnpm 11.5.2。
启动入口是 [frontend/src/main.tsx](frontend/src/main.tsx)，路由和壳层分别在
[src/app/router.tsx](frontend/src/app/router.tsx) 与 [src/app/app-shell.tsx](frontend/src/app/app-shell.tsx)。

    frontend/
    ├── public/                     favicon、grok2api 标识、runtime-config.js、赞助图片
    ├── src/app/                    providers、鉴权边界、路由、应用壳层和懒加载页面
    ├── src/entities/               模型、系统等跨页面实体 API/类型
    ├── src/features/               accounts、models、client-keys、audits、dashboard、
    │                               media、creative-console、quality-guard、settings、docs 等
    ├── src/components/ui/          Radix/shadcn 风格的可复用 UI 组件
    ├── src/shared/                 API client/decoder、认证、运行配置、i18n、表格和工具
    ├── src/types/                  运行配置等全局类型声明
    ├── package.json                dev/build/lint/preview 和依赖声明
    ├── pnpm-lock.yaml              依赖锁定文件
    └── vite.config.ts              @ 别名、5173 开发端口、API 代理和 dist 输出目录

管理端通过 /api/admin/v1/* 调用后端；public/runtime-config.js 和
src/shared/config/runtime-config.ts 允许部署时注入 API 地址。开发服务器将
/api、/v1、/healthz、/readyz 代理到 VITE_DEV_API_TARGET 或本机 8000。

frontend/dist/、frontend/node_modules/、.pnpm-store/ 和 TypeScript/Vite 缓存
都是生成物或运行态，均被忽略，不应手工提交。Docker build 会先在前端 stage 运行
pnpm build，再把 dist 复制进最终镜像；后端源码仓库不依赖已提交的 dist。

## 5. Docker 与可选运维工具

### 5.1 镜像和 Compose

[`Dockerfile`](Dockerfile) 是三阶段构建：Node Alpine 构建前端，Go Alpine 编译无 CGO 后端，最终
使用 Alpine 非 root 用户 grok2api 运行。运行镜像包含二进制、前端 dist、VERSION
和 [docker/entrypoint.sh](docker/entrypoint.sh)，监听默认 8000，并以 /healthz 做 Docker healthcheck。

entrypoint 的边界是：从只读挂载的 /run/grok2api/config.yaml 复制配置到容器内
/app/config.yaml，收紧权限后再以非 root 用户执行主进程。它不生成账号、不导入凭据，
也不替代配置校验。

[docker-compose.yml](docker-compose.yml) 的主服务使用 grok2api-data 保存数据库和媒体；
quality_guard_state 保存质量守护 bootstrap、状态和锁。FlareSolverr、WARP 只是
注释示例，egress-quality-guard 只有在 quality-guard profile 显式启用时才构建
和运行。Compose 默认镜像仍指向上游 GHCR；要运行当前源码构建的镜像，必须由部署者
显式设置 GROK2API_IMAGE 或执行自己的 build，不能从源码 checkout 推断容器已经切换。

### 5.2 tools/egress-quality-guard

这是当前仓库唯一的运维工具目录，使用 Python 标准库实现（详见
[中文运维说明](tools/egress-quality-guard/README.zh-CN.md)）：

- quality_guard.py：读取主程序生成的受限 bootstrap，支持 passive/active/hybrid
  检测，按固定模型探测出口质量，隔离/恢复节点但不删除节点或账号绑定。
- session_rotator.py：可选的 1024Proxy 粘性会话轮换器，通过回环 HTTP 接口更新
  Mihomo 配置并验证出口 IP；凭据文件和 Mihomo 配置是外部 owner-controlled 输入。
- 两个 Dockerfile：分别打包质量守护 sidecar 和 session rotator。
- README(.zh-CN).md、SECURITY.md：操作、安全和限制说明；两个 *_test.py 是
  Python 单元测试入口。

质量守护的内部 API、bootstrap 和状态卷是服务运行态边界，不是公共 API，也不能把
sidecar 的启动状态当作主服务或 Smart Search consumer 的 live 证据。

## 6. 源码、生成物、运行态与凭据边界

| 类别 | 归属与规则 |
| --- | --- |
| 源码/契约 | Go、TypeScript/React、Python、YAML/Compose、Dockerfile、测试、Issue 模板、锁文件和许可证进入 Git；变更应落在其拥有的目录。 |
| 生成物 | backend/docs/docs.go、backend/docs/swagger.json、backend/docs/swagger.yaml 由 make swagger 从后端注释生成并纳入 Git；frontend/dist、镜像层、Go/Node 缓存不纳入 Git。生成 Swagger 后必须检查只出现预期 diff。 |
| 运行态 | 根 config.yaml、data/、backend/data/、SQLite/WAL/SHM、媒体目录、日志、Redis 数据、Quality Guard state/bootstrap/lock、Docker volumes、frontend/node_modules、.pnpm-store、.gocache 等均在部署或本机运行时产生。 |
| 凭据 | secrets.jwtSecret、credentialEncryptionKey、bootstrap 管理员密码、Build OAuth、Web/Console SSO/Cookie、Client Key 明文、数据库/Redis DSN 密码、代理/FlareSolverr/Statsig/rotation token 和质量守护内部 token 都必须留在 owner-controlled secret store 或权限受限运行态文件。 |

数据库会保存账号、凭据密文、模型、额度、Client Key、审计和媒体任务；
credentialEncryptionKey 负责凭据加密，写入账号后必须长期保留，但它本身不能
提交到 Git。多实例模式要求 PostgreSQL、Redis、稳定 instance/cluster 标识和共享
媒体目录；单实例默认是 SQLite + Memory。上述是代码和配置契约，不代表当前工作区
已经完成这些部署条件。

禁止把账号导出、Access/Refresh Token、Cookie、代理密码、管理员密码、生产日志、
数据库快照或 sidecar 状态卷写入 STRUCTURE.md、Issue、commit message 或测试输出。

## 7. 当前 fork 的个人改动边界

当前分支相对 `fork/main` 的可见差异由 `git diff fork/main...lwj_dev` 核对，
可归纳为两类：

1. **Web 搜索/引用和协议适配**：集中在
   backend/internal/infra/provider/web/{chat.go,gateway.go,responses_stream.go,image.go}、
   backend/internal/infra/provider/conversation/response.go 及对应测试，处理
   Gateway 的 render_citation、tool_usage_card/tool_result、Hosted Search
   call、URL citation、Chat/Responses SSE 形状和 xAI 工具使用信息。
2. **客户端请求状态跟踪**：新增/修改
   backend/internal/application/requeststatus、clientkey/service.go、
   backend/internal/transport/http/{inference/handler.go,middleware/auth.go,middleware/request_status.go,server.go}
   及测试，以客户端身份隔离短期 request status 查询，并保留系统质量守护 Client
   Key 不可人工操作的边界。

此外 VERSION 在该分支提升到 v3.1.2。上述是当前分支与个人 fork 基线的差异
摘要，不声称 fork/main 之前的所有提交都由本地个人编写；合并提交带入的上游
内容仍按上游代码维护。新增 fork-specific 逻辑应优先放在对应 Provider/transport
边界，配套测试，并在本节和相关运行文档中留下可核对的文件范围。

## 8. 开发、验证与上游同步入口

### 8.1 本地开发入口

前置工具版本以构建文件和 CI 为准：Go 1.26、Node 22、pnpm 11.5.2、Python 3
和可选的 Docker/Compose。源码运行需要从 config.example.yaml 复制出本地
config.yaml，填入真实密钥和 Provider 配置；该文件被忽略，不应由本仓库文档代填。

    # 根目录
    make run

    # 后端（backend/README.md 也提供了同一入口）
    cd backend
    go test ./...
    go test -race ./...
    go vet ./...
    go build ./cmd/grok2api

    # 前端
    cd frontend
    pnpm install --frozen-lockfile
    pnpm lint
    pnpm build

### 8.2 契约、工具与镜像验证

    # 从根目录重新生成已提交的 Swagger 文件，并确认没有未预期修改
    make swagger
    git diff --exit-code -- backend/docs/docs.go backend/docs/swagger.json backend/docs/swagger.yaml

    # 出口质量守护测试
    python3 -m unittest -v \
      tools/egress-quality-guard/quality_guard_test.py \
      tools/egress-quality-guard/session_rotator_test.py

    # 仅验证 Compose 展开，不启动服务
    docker compose --profile quality-guard config --quiet

[.github/workflows/ghcr-image.yml](.github/workflows/ghcr-image.yml) 是 CI 验证的权威组合：后端 test/vet、Swagger
一致性、前端 frozen install/lint/build，以及 amd64/arm64 Docker build。/healthz
或 /readyz 只有在明确的本地/部署进程上执行时才证明该进程的健康/就绪，不能证明
smart-search-cli 或 grok2api-service 的真实调用链。

上游同步的安全顺序：

1. 先检查 git status --short --branch、当前 HEAD、远端和待保留的 dirty files。
2. 在允许联网的独立步骤中 git fetch origin main，再用
   git log --left-right --cherry-pick origin/main...lwj_dev 和
   git diff origin/main...lwj_dev 审阅上游与个人差异；不要拿未 fetch 的 ref
   猜测上游状态。
3. 保留 lwj_dev 的 fork-specific 提交和未提交改动，确认冲突归属后再由 owner
   明确选择 merge/rebase；同步本身不应覆盖 STRUCTURE.md、运行配置或 dirty files。
4. 同步后重新运行后端、前端、Swagger、工具和必要的 Docker 构建检查，更新本文
   的快照 SHA、分支差异和消费者/运行态边界。
5. fetch、merge/rebase、commit、push、发布镜像或切换消费者都是独立的状态变更，
   不由“更新文档”授权；push 只在 owner 明确要求时执行。

## 9. 文档更新规则

- 顶层目录、backend layer、Provider Definition、HTTP 路由、frontend feature、
  Compose profile、生成物/运行态/凭据边界或 Infra consumer manifest 改变时，更新
  本文对应章节。
- 面向使用者的安装、配置、Provider 能力和 API 示例改动同步 README.md、
  README.zh-CN.md；后端代码层次和验证入口同步 backend/README.md；质量守护行为同步
  tools/egress-quality-guard/README.zh-CN.md 和安全文档。
- 修改 Swagger 注释或公开接口后，从根目录运行 make swagger，提交生成的
  backend/docs/* 变化；不要直接手改生成文件来隐藏契约差异。
- 只记录已在源码、配置、Git 差异或验证命令中可复核的事实；“已启动”“已激活”、
  “consumer 可用”必须引用对应阶段的运行态/真实请求证据，不能由 README 或结构
  文档推断。
- 文档不得包含明文凭据、账号导出、生产日志、临时本地路径中的秘密或不可复现的
  live 结论。每次结构性同步都要重新核验当前分支和远端关系。
