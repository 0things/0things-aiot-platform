## 1. MCP HTTP 身份边界

- [x] 1.1 为 MCP Streamable HTTP transport 增加 Bearer JWT 验证、用户/组织 `context.Context` 注入与未认证拒绝逻辑；验证 MCP 的 `tools/list` 和三个 `tools/call` 在无 token 时均返回认证失败。
- [x] 1.2 移除 Streamable HTTP 请求到默认组织的回退路径，并让 `iot_query_devices`、`iot_get_device_detail`、`iot_query_telemetry_history` 使用 token 中的组织上下文；验证组织 A token 无法读取组织 B 设备、详情和遥测。
- [x] 1.3 为 Web Copilot 的 MCP HTTP 认证补充配置说明和集成测试，同时确认 STDIO transport 行为不变；验证 `GOTOOLCHAIN=local go test ./internal/server ./internal/handler ./pkg/server/mcp` 通过。

## 2. AI Gateway 基础设施

- [x] 2.1 新建独立的 `ai-copilot/` TypeScript/Node.js 服务结构、Node 25 引擎约束、启动脚本、health/readiness 端点和最小化结构化日志；验证独立安装依赖后可启动并返回健康状态。
- [x] 2.2 添加 Vercel AI SDK、`@ai-sdk/mcp`、OpenAI-compatible provider 和 JWT 验证依赖，锁定兼容版本并提供不含秘密的 `.env.example`；验证 `pnpm install --frozen-lockfile`（或该服务定义的等效安装命令）成功。
- [x] 2.3 实现 Gateway 配置校验、JWT 校验、中间件请求 ID 与敏感配置脱敏；验证缺少模型/JWT/MCP 必需配置时拒绝启动，合法 token 能通过且无 token 返回 401。

## 3. 流式聊天与 MCP 编排

- [x] 3.1 实现 `POST /v1/ai/chat`：校验 UI message 请求、调用 OpenAI-compatible 模型，并以 Vercel AI SDK `toUIMessageStreamResponse()` 返回流；验证模拟模型响应的契约测试可解析文本增量。
- [x] 3.2 使用 Streamable HTTP MCP client 建立每请求受控连接，将已验证 Bearer token 仅转发给 MCP，并在请求结束时关闭连接；验证 Gateway 到 MCP 的请求带有身份且模型 provider 请求不含该 token。
- [x] 3.3 仅将 `iot_query_devices`、`iot_get_device_detail`、`iot_query_telemetry_history` 注册给模型，限制工具循环和工具结果大小；验证请求“修改/删除/下发命令”不会产生写操作调用。
- [x] 3.4 建立统一的模型、MCP、超时和取消错误映射及安全日志字段；验证 UI stream 不包含 token、API Key、MCP URL、数据库错误或堆栈，且 MCP 不可用时返回可重试错误码。

## 4. 管理端 Copilot 体验

- [x] 4.1 在 `frontend` 安装 `@assistant-ui/react`、`@assistant-ui/ai-sdk`、`ai` 和 `@ai-sdk/react`，配置 `VITE_AI_GATEWAY_URL` 与开发代理，不在前端配置模型或 MCP 凭据；验证 `corepack pnpm@11.8.0 install` 与类型检查成功。
- [x] 4.2 在 `src/features/ai-copilot/` 实现基于 `useChatRuntime` 的 Copilot runtime、Bearer transport、内存新建对话和取消/重试操作；验证刷新页面或点击新建对话后消息状态清空，且现有页面不发生导航。
- [x] 4.3 使用现有 shadcn 组件在认证布局挂载可收起 Copilot 入口和抽屉，呈现流式消息、输入状态和工具“查询中/完成/失败”状态；验证在设备管理和规则链页面均可打开并关闭。
- [x] 4.4 为 Copilot 的按钮、空状态、只读能力提示和服务失败提示补齐中英文 i18n；验证切换 `zh`/`en` 时无新增硬编码可见文案，MCP 不可用时显示安全的中文重试提示。
- [x] 4.5 以 `VITE_AI_COPILOT_ENABLED` 实现入口开关，默认按部署环境配置；验证关闭时不渲染入口且不向 Gateway 发请求，开启时正常建立聊天流。

## 5. 端到端验证与交付

- [x] 5.1 编写 AI Gateway 的 JWT、工具白名单、MCP token 透传、流式成功与安全错误测试；验证服务测试命令全部通过。
- [x] 5.2 使用组织 A/B 测试数据完成端到端测试：查询设备列表、设备详情、遥测历史，以及猜测跨组织设备标识；验证只有本组织数据返回。
- [x] 5.3 执行全量回归 `make test && make build`、`corepack pnpm@11.8.0 -C frontend format:check`、`corepack pnpm@11.8.0 -C frontend build` 以及 AI Gateway 测试；验证所有命令通过并记录部署所需环境变量。
