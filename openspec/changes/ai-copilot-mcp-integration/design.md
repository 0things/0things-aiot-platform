## Context

0things 是一个典型的现代 AIoT 平台，支持多租户组织（Organization）、设备生命周期、TSL 物模型、TSDB 时序遥测及规则链。随着平台功能丰富，跨页面多步骤排障（例如：先去设备列表查在线状态 -> 去时序图表查异常温度 -> 去告警中心核对级别 -> 去控制台下发指令 -> 去规则链配置阈值）存在割裂与操作门槛。

通过集成大语言模型（LLM）与开放的 **Model Context Protocol (MCP)** 标准，配合前沿的 **`assistant-ui` + Vercel AI SDK**，打造端到端的智能运维 Copilot，使平台兼具“AI 对话驱动运维”与“对外开放 MCP 工具供外部 IDE/Agent 接入”的双重能力。

## Goals / Non-Goals

**Goals:**
- **双模 MCP 架构**：实现统一的 `internal/mcp` IoT 工具集，既可通过独立的 `cmd/mcp` 编译为支持 STDIO / SSE 的独立二进制（供 Cursor/Claude Desktop 等外部接入），也可在 `cmd/server` 进程内（In-Process）无损高效运行。
- **现代化前端对话界面**：基于 `@assistant-ui/react` 与 `@assistant-ui/react-ai-sdk`，提供支持折叠/全屏的全局 Copilot 抽屉、打字机流式输出与思考过程展示。
- **安全的 Human-in-the-Loop 控制卡片**：对查询类只读工具自动执行，对设备控制（如开关机、下发命令）、规则发布等高危写操作，在对话流中渲染交互式二次确认卡片。
- **多租户数据隔离**：所有 MCP Tool 调用强制继承当前认证用户的 `organization_id`，杜绝越权访问。
- **开箱即用的静态模型配置**：通过后端 `config.yaml` 统一配置兼容 OpenAI 协议的主流大模型（DeepSeek, Qwen, OpenAI, Ollama 等）。

**Non-Goals:**
- 首期不在前端提供动态添加/编辑大模型 API Key 的管理页面，仅支持服务端 `config.yaml` 静态配置。
- 首期不引入重型的 Python/Node.js 外部 AI 框架，全面使用 Go 原生组件实现。

## Decisions

### 1. 双模 MCP 架构与模块分层（Dual-Mode MCP Architecture）

```text
                               ┌────────────────────────────────────────┐
                               │     外部 Client (Cursor / Claude)      │
                               └──────────────────┬─────────────────────┘
                                                  │ (STDIO / SSE)
                                                  ▼
                               ┌────────────────────────────────────────┐
                               │          cmd/mcp/main.go (独立入口)     │
                               │          pkg/server/mcp (通用包装层)    │
                               └──────────────────┬─────────────────────┘
                                                  │
                                                  ▼
┌─────────────────────────────────┐   ┌─────────────────────────────────┐
│        cmd/server/main.go       │──►│        internal/mcp/            │
│  (Web 业务服务: /v1/ai/chat)    │   │  (IoT Tools: 查设备/遥测/告警/控│
└─────────────────────────────────┘   │   制/规则，统一租户校验与审计)  │
                                      └────────────────┬────────────────┘
                                                       │
                                      ┌────────────────┴────────────────┐
                                      ▼                                 ▼
                             GORM DAL (PostgreSQL)             pkg/tsdb (时序库)
```

- **`pkg/server/mcp/server.go`**：通用封装层，基于 `github.com/mark3labs/mcp-go`，提供 Stdio、SSE 和 HTTP 模式的优雅启停。
- **`cmd/mcp/`**：独立命令行入口，供本地 IDE（Cursor / Claude Desktop）通过 `stdio` 方式拉起，或作为独立远程 MCP 微服务部署。
- **`cmd/server/`**：主 Web API 服务，直接在内存中（In-Process）加载 `internal/mcp`，供 Web 界面 AI Copilot 调用，0 网络开销。

*备选方案比较*：若仅在 `cmd/server` 进程内硬编码 Tool，将无法对外输出标准 MCP能力；若仅做外部 HTTP MCP，则内部 Copilot 每次 Tool Call 都要多走一次网络序列化与鉴权握手。双模共享 `internal/mcp` 为最优解。

### 2. 前端技术栈：`assistant-ui` + Vercel AI SDK

- 前端引入 `@assistant-ui/react`、`@assistant-ui/react-ai-sdk` 和 `ai`。
- 全局布局（`AuthenticatedLayout`）中嵌入 `<CopilotDrawer>` 浮动抽屉。
- 使用 `useVercelUseChatRuntime` 或标准 SSE 协议与后端 `POST /v1/ai/chat` 建立连接。
- 利用 `assistant-ui` 的 `makeAssistantToolUI` / `ToolUI` 注册自定义组件：
  - 工具调用执行进度条（如 `正在查询时序历史... [耗时 45ms]`）；
  - 命令下发二次确认卡片（`<DeviceCommandConfirmCard />`）。

### 3. 安全与 Human-in-the-Loop 二次确认流程

```text
用户: "将设备 dk_device_047 关机"
       │
       ▼ (LLM 判定需要调用 iot_send_device_command)
后端识别到该工具为 [高危/写操作]
       │
       ▼ (后端暂停执行，向 SSE 推送 Tool Call 待确认帧)
前端 assistant-ui 渲染确认卡片:
┌────────────────────────────────────────────────────────┐
│ ⚠️ 设备控制二次确认 (Human-in-the-loop)                 │
│ 目标设备: dk_device_047 (空调机组 #047)                │
│ 参数: {"power_switch": false}                          │
│                                                        │
│       [ 确认执行 ]               [ 取消操作 ]          │
└────────────────────────────────────────────────────────┘
       │
       ▼ (用户在卡片中点击【确认执行】)
前端将确认动作发送至后端 -> 后端校验 token -> 调用下发逻辑 -> 返回执行成功结果 -> LLM 生成最终总结
```

### 4. 0things 原生 MCP Tool 定义规范

| Tool 名称 | 参数 (JSON Schema) | 风险级别 | 底层依赖 |
| :--- | :--- | :--- | :--- |
| `iot_query_devices` | `product_id?`, `group_id?`, `keyword?`, `status?`, `limit?` | 只读 (自动执行) | `DeviceService.List` |
| `iot_query_telemetry_history` | `device_key!`, `property!`, `start_time?`, `end_time?`, `limit?` | 只读 (自动执行) | `tsdb.Client.QueryPoints` |
| `iot_get_active_alerts` | `device_key?`, `severity?`, `limit?` | 只读 (自动执行) | `AlertService.ListActive` |
| `iot_send_device_command` | `device_key!`, `service_name!`, `params!` | **高危写操作 (二次确认)** | `DeviceService.SendCommand` / Kafka |
| `iot_generate_rule_chain` | `description!`, `scope_type?`, `scope_id?` | **写操作 (生成草稿待确认)** | `RuleService.SaveDraft` |

### 5. 配置文件与提供商设计（`config.yaml`）

```yaml
ai:
  enabled: true
  provider: "openai" # 支持 openai, deepseek, qwen, ollama 等兼容接口
  base_url: "https://api.deepseek.com/v1"
  api_key: "sk-your-api-key"
  model: "deepseek-chat"
  temperature: 0.3
  max_tokens: 4096
  system_prompt: |
    你是一个专业的 0things 工业物联网平台智能助手。
    你可以使用提供的 MCP 工具查询设备状态、时序历史、告警排查并协助运维。
    在涉及修改状态或下发物理设备控制命令前，必须明确提示用户。

mcp:
  name: "0things-iot-mcp"
  version: "1.0.0"
  stdio_enabled: true
  sse_enabled: false
  sse_addr: ":8008"
```

## Risks / Trade-offs

- [大模型 Token 超长与上下文截断] → 单次遥测时序点查询限制 `limit=100`，并在 Tool 内部计算好统计聚合指标（Min/Max/Avg/Latest），将结构化摘要提供给 LLM，而非倾倒全量原始采样点。
- [弱网与多租户越权] → 每个 Tool 执行严格从 Gin/JWT context 提取 `organization_id`，若目标设备不属于当前组织直接拒绝执行。
- [二次确认卡片状态时效性] → 确认动作携带唯一 `action_token`，设置 5 分钟有效期，过期后卡片置灰失效，防止重复下发。

## Migration Plan

1. 后端引入 `github.com/mark3labs/mcp-go` 依赖，新建 `pkg/server/mcp` 与 `internal/mcp`。
2. 建立 `cmd/mcp` 入口并编写单元测试验证 Stdio / In-Process 工具调用。
3. 后端 `cmd/server` 接入 `POST /v1/ai/chat` SSE 端点与 OpenAI 兼容客户端。
4. 前端引入 `@assistant-ui/react` 与 `@assistant-ui/react-ai-sdk`，实现 Copilot 抽屉与 Tool UI 确认卡片。
5. 在 `config.yaml` 补充 `ai` 与 `mcp` 默认配置项。

