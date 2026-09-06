## Why

0things AIoT 平台目前拥有设备管理、TSL 物模型、时序遥测与规则链等核心能力，但运维人员在面对海量设备、复杂告警排查、时序指标分析和跨模块操作时，仍依赖繁琐的多页面手动跳转与表单配置。

引入 AI Copilot 智能运维对话助手，借助大语言模型（LLM）的自然语言理解与推理能力，结合 Model Context Protocol (MCP) 开放标准，使运维人员能够通过自然语言对话一键完成“设备状态查询、时序指标分析、活动告警诊断、物理设备控制指令下发与规则链辅助编排”，极大降低平台使用门槛并提升运维排障效率。

## What Changes

- **前端 AI 交互中心（assistant-ui + Vercel AI SDK）**：
  - 引入 `@assistant-ui/react` 与 `@assistant-ui/react-ai-sdk`，在平台右下角/侧边栏提供现代化的 AI Copilot 对话抽屉（支持折叠、全屏、Thread 多会话管理与智能触底）。
  - 对接 Vercel AI SDK Data Stream 协议，支持流式打字机输出、思考链展示与实时 MCP Tool Call 状态渲染。
  - 基于 `assistant-ui` 的 Tool UI / Generative UI 机制，实现设备控制与高危动作的 **Human-in-the-loop 交互式二次确认卡片**。
- **独立 MCP 独立命令入口（`cmd/mcp`）**：
  - 新增 `cmd/mcp/main.go` 独立入口，支持标准 **STDIO**（供 Cursor / Claude Desktop 等本地 IDE 直接调用）与 **独立 SSE/HTTP 远程服务**。
  - 封装通用 `pkg/server/mcp` 服务包装器，支持优雅启停与中间件 Hooks。
- **后端 AI 编排与流式服务（`cmd/server`）**：
  - 提供 `POST /v1/ai/chat` 流式端点，兼容 Vercel AI SDK 数据流协议（SSE）。
  - 进程内（In-Process）直接复用 `internal/mcp` 注册的 IoT Tools 执行函数，实现零网络开销的 Function Calling 闭环。
  - 集成大模型客户端，支持通过 `config.yaml` 静态配置主流大模型提供商（OpenAI, DeepSeek, Qwen 等兼容接口）。
  - 强制注入多租户 `organization_id` 与 JWT 身份上下文，保证数据隔离与权限校验。
- **0things 原生 MCP Server（内置工具集）**：
  - 基于 Go MCP SDK (`mark3labs/mcp-go`) 实现 0things IoT 核心工具集：
    - `iot_query_devices`：多维度设备列表与在线状态查询。
    - `iot_query_telemetry_history`：时序指标历史查询与统计分析（调用 `pkg/tsdb`）。
    - `iot_get_active_alerts`：活动告警检索与严重等级排查。
    - `iot_send_device_command`：物理设备控制指令下发（需二次确认）。
    - `iot_generate_rule_chain`：自然语言辅助生成规则链拓扑与节点配置 JSON。
  - 支持在 `config.yaml` 中配置外部标准 MCP Server（HTTP/SSE）以实现能力横向扩展。

## Capabilities

### New Capabilities
- `ai-copilot`: 平台级智能运维助手交互规范，包含 assistant-ui 对话界面、Vercel AI SDK 流式通信协议、会话管理、Tool UI 渲染与敏感写操作的 Human-in-the-loop 二次确认流。
- `mcp-server`: 0things 原生 IoT MCP 工具集规范、独立 `cmd/mcp` 服务与外部 MCP 协议网关，定义只读/写入工具契约、租户数据边界与参数 Schema。

### Modified Capabilities
<!-- 无现有 Spec 需求变更 -->

## Impact

- **后端**：
  - 新增 `backend/cmd/mcp/`（独立 MCP 二进制入口与 Wire 依赖注入）。
  - 新增 `backend/pkg/server/mcp/`（通用 MCP 服务包装层）。
  - 新增 `backend/internal/mcp/`（IoT Tools 实现与审计 Hooks）。
  - 新增 `backend/internal/ai/`（LLM 客户端、Agent 协调器、Vercel AI SDK 兼容 SSE 适配器）。
  - `config.yaml` 新增 `ai` 与 `mcp` 配置块（模型提供商、BaseURL、APIKey、Model、Temperature、System Prompt、STDIO/SSE 端口等）。
  - 依赖新增 `github.com/mark3labs/mcp-go`。
- **前端**：
  - 新增 `frontend/src/features/ai-copilot/` 模块。
  - `package.json` 引入 `@assistant-ui/react`、`@assistant-ui/react-ai-sdk`、`ai` 及相关 Markdown 依赖。
  - 在全局布局组件（`AuthenticatedLayout`）中嵌入 Copilot 悬浮气泡与抽屉面板。
