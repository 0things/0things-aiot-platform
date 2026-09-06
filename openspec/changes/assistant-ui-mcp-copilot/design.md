## Context

详见 [proposal.md](./proposal.md)。管理端是 Vite + React 单页应用，后端是 Go 服务；两者都不是可运行 Vercel AI SDK 服务端中间件的 Node.js 运行时。现有 MCP 服务已提供 Streamable HTTP（`/mcp`）和三个只读 IoT 工具，但其 HTTP 调用尚未携带调用用户的组织身份。浏览器也不应获得模型密钥或内部 MCP 地址。

## Goals / Non-Goals

**Goals:**

- 使用官方 Vercel AI SDK 与 MCP 客户端实现稳定的 UI Message Stream，避免在 Go 中手写、维护 AI SDK 流协议。
- 让 React UI 只知道一个受认证的聊天入口；模型、MCP 地址和凭据均处于服务端。
- 让工具执行沿用 MCP 服务已有的组织数据隔离，并在 UI 中有可观察的执行状态。
- 将 AI Gateway 保持为无状态服务，使其可水平扩展；首期会话由浏览器内存持有。

**Non-Goals:**

- 不将 AI Gateway 变成通用业务 API 或替代 Go 后端。
- 不保存聊天消息、提示词、工具入参与输出，也不提供会话回放、审计检索或记忆。
- 不实现模型供应商管理 UI、多模型路由、用户自带模型 Key、用户自定义 MCP、知识库/RAG 或写操作审批。
- 不改变既有三个 MCP 工具的名称、入参和只读语义。

## Decisions

### 1. 使用独立的 Node.js AI Copilot 服务，而不是由 Go 后端实现 AI SDK 流协议

新增 `ai-copilot/` 独立 TypeScript 服务，运行时使用 Node.js 25，与前端的 Node 版本约束一致。其 HTTP 层采用轻量标准 Request/Response 适配器；核心依赖为 `ai`、`@ai-sdk/mcp`、OpenAI-compatible provider 和 JWT 校验库。它暴露唯一的业务入口：`POST /v1/ai/chat`。

Gateway 负责：验证 Bearer JWT、建立模型请求、连接 MCP、白名单化工具、把 `streamText` 的 `toUIMessageStreamResponse()` 原样返回，以及将内部失败转为安全的流错误。它不访问数据库，不持久化对话，也不作为前端静态资源服务。

选择原因：`useChatRuntime` 消费 Vercel AI SDK 的 UI stream；Node Gateway 可直接使用 AI SDK 和 `@ai-sdk/mcp`，从而降低协议兼容和 MCP 客户端维护成本。

### 2. 采用受控的三段式调用链与 JWT 身份转发

调用路径固定为：

```text
authenticated React client
  -- Authorization: Bearer <user JWT> --> AI Gateway /v1/ai/chat
  -- Authorization: Bearer <same verified JWT>, private network --> MCP /mcp
  -- organization context --> existing device / telemetry services
```

AI Gateway 和 MCP HTTP transport 都使用与主后端一致的 JWT 验证配置，并从 token 取得用户和组织声明；MCP 不再对 Streamable HTTP 请求回退到默认组织。Gateway 只把验证后的调用所需 Bearer token 透传给 MCP，模型供应商永远不会收到该 token。

MCP 的 Streamable HTTP 处理器需要在工具列举和调用之前验证请求，并把身份写入 `context.Context`，使现有 Handler/Service 保持组织隔离。MCP 地址只部署在私有网络或由网络策略限制为 Gateway；即使网络策略配置错误，无 Bearer token 的请求也必须被拒绝。

STDIO 保持给本地开发者 MCP 客户端使用，不作为 Web Copilot 的身份通道；本变更不改变其既有行为。

备选方案是 Gateway 以服务账户调用 MCP 并在自定义 header 传组织 ID。该 header 更容易被伪造、无法表达当前用户授权边界，且会绕过既有 JWT 认证模型，故不采用。

### 3. 前端只集成 assistant-ui 作为展示与交互层

在 `frontend/src/features/ai-copilot/` 放置 Copilot feature，并从认证布局挂载全局入口。组件采用 shadcn 的 `Sheet` / `Button` / `ScrollArea` 等现有 UI 原语承载抽屉和入口；assistant-ui 的 `useChatRuntime` 与 AI SDK transport 管理消息、提交、流式增量和工具调用状态。

自定义 transport 将当前 Clerk/应用认证获得的 Bearer token 带到 Gateway，API 基址由 `VITE_AI_GATEWAY_URL` 配置。开发环境通过 Vite proxy 或明确环境变量访问 Gateway；生产环境必须使用 HTTPS 的同源反向代理或受信任跨域配置。前端不得出现 MCP URL、模型 URL 或任意 provider Key。

聊天状态仅存于 runtime 内存；刷新页面、关闭页面或点击“新建对话”都会开始新会话。所有可见文案写入 `public/locales/zh` 和 `public/locales/en`。

### 4. MCP 工具在 Gateway 中显式白名单，而非运行时全量暴露

Gateway 可以通过 MCP 发现工具以便连接校验，但传给模型的工具集只能筛选为：

- `iot_query_devices`
- `iot_get_device_detail`
- `iot_query_telemetry_history`

每轮最大工具循环次数设置为有限值（建议 5），并沿用工具端的遥测 `limit` 上限。UI 仅展示本地化的工具显示名和“查询中/已完成/失败”状态；默认不展示原始请求参数、完整返回值或内部错误。

这样即使将来 MCP 新增写操作，Copilot 也不会自动获得执行权限。备选的“发现即全部注册”会把今后新增工具意外暴露给模型，故不采用。

### 5. 模型与服务配置全部在服务端环境变量中管理

AI Gateway 使用 OpenAI-compatible 配置：`AI_MODEL_BASE_URL`、`AI_MODEL_API_KEY`、`AI_MODEL_ID`；还需要 `MCP_STREAMABLE_HTTP_URL`、JWT 验证所需的安全配置、`AI_GATEWAY_PORT` 和可选 `AI_GATEWAY_ALLOWED_ORIGINS`。生产配置由部署平台 Secret 注入，示例配置只保留占位符，日志对以上字段做脱敏。

当未配置模型或 MCP 时，Gateway 启动应失败并给出运维可诊断日志；运行时模型/工具故障则返回面向用户的可重试提示。绝不将配置值、上游响应 body 或堆栈写进 UI 流。

## Risks / Trade-offs

- [新增 Gateway 增加一个部署单元和监控点] → 提供独立 health/readiness 端点、容器启动配置及最小日志字段（请求 ID、组织 ID、耗时、工具名、结果类别），不记录聊天正文或秘密。
- [Gateway 与 MCP 都要验证 JWT，密钥/公钥配置可能漂移] → 使用同一部署来源注入认证配置；启动检查校验必要字段；以跨组织访问集成测试作为发布门槛。
- [模型可能幻觉或用词不精确] → 系统提示明确只可基于工具结果回答；限制三项只读工具；没有数据时要求明确说明无法确认。
- [MCP 或模型慢会拖住浏览器连接] → Gateway 设置模型、MCP 和总请求超时；前端支持取消与重试；超时使用统一中文错误。
- [流式协议或依赖升级导致 assistant-ui 兼容性变化] → 以 AI SDK UI stream 为 Gateway 的唯一对外协议，并为流响应与 tool part 写契约测试；依赖采用锁定的兼容版本组升级。
- [遥测返回量或工具输出过大] → 强制工具 limit 上限、限制模型上下文中的工具结果大小，并向用户提示缩小时间范围。

## Migration Plan

1. 先部署 MCP HTTP 认证改动和三个只读工具，使用服务间 JWT 集成测试验证组织隔离；不启用 Gateway 流量。
2. 部署 AI Gateway，配置模型、JWT 与私网 MCP URL，运行 health/readiness 和三个工具的冒烟检查。
3. 部署前端 Copilot，但通过 `VITE_AI_COPILOT_ENABLED=false` 保持入口关闭；在测试组织启用后验证流式问答、工具状态和失败提示。
4. 验证通过后按组织或环境打开入口。首期不需要数据迁移，因为不写任何聊天数据表。

回滚时先关闭前端 feature flag，再停止 Gateway 流量；MCP 的 Streamable HTTP 保持可用，不影响现有 STDIO 连接。若 MCP 身份校验改动出现问题，可只回滚该 HTTP transport 改动，已存在的设备业务 API 不受影响。

## Open Questions

- 首个生产模型供应商与额度策略尚未确定；这不改变 OpenAI-compatible 接口或本变更的工具边界，部署时由环境变量选择。
- Copilot 入口的最终摆放位置（顶部导航或右下角悬浮按钮）可在实现时按现有认证布局视觉规范决定，不影响通信和安全设计。
