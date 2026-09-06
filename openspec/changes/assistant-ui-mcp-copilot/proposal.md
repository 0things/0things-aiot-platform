## Why

0things 已提供三个只读 IoT MCP 工具，但平台前端没有面向运维人员的 AI 对话入口。需要将自然语言提问、模型推理和 MCP 工具调用组成受认证保护的完整闭环，而不是让浏览器直接暴露 MCP 或模型凭据。

## What Changes

- 在 React 管理端引入 `assistant-ui`，使用 `useChatRuntime` 提供全局可展开的 AI Copilot 抽屉。
- 新增独立的 Node.js AI Gateway：受 JWT 保护地提供流式聊天 API，采用 Vercel AI SDK 调用兼容 OpenAI 的模型，并原生输出 UI stream。
- AI Gateway 在服务端通过 Streamable HTTP 连接现有 0things MCP，动态发现并调用三个只读 IoT 工具。
- 将当前用户的组织身份安全传递到内部 MCP 调用，禁止浏览器直接访问 MCP 地址或模型 API Key。
- 首期仅保留内存会话；不持久化聊天记录、不开放用户自定义 MCP、不包含任何写操作工具。

## Capabilities

### New Capabilities

- `ai-copilot-chat`: 全局 Copilot UI、流式聊天接口、服务端 MCP 工具编排与组织隔离。

### Modified Capabilities

无。

## Impact

- 前端：新增 `@assistant-ui/react`、`@assistant-ui/ai-sdk`、`ai`、`@ai-sdk/react` 和 AI Copilot feature；在认证布局挂载入口。
- AI Gateway：新增 TypeScript/Node.js 服务，使用 Vercel AI SDK、`@ai-sdk/mcp` 与 OpenAI-compatible provider，提供 `POST /v1/ai/chat`。
- MCP：为 Streamable HTTP 内部服务间调用增加可信 JWT 身份上下文；继续使用 `http://127.0.0.1:8009/mcp`，不改变对外工具名称。
