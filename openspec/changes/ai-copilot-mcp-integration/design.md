## Context

Model Context Protocol (MCP) 是当前连接大模型与本地/远程工具的通用标准。0things AIoT 平台已有完整的设备管理、TSDB 时序遥测、告警中心与指令下发能力。

参考 `tmp/nunu-layout-mcp-main` 的工程规范，在 `backend` 内部构建一个原生的、纯粹的 **0things IoT MCP Server**，通过 STDIO 供本地 Cursor / Claude Desktop 拉起，或通过 SSE 供远程 AI 平台连接。

## Goals / Non-Goals

**Goals:**
- **标准 MCP 协议支持**：基于 `github.com/mark3labs/mcp-go`，提供 STDIO 与 SSE 两种标准传输协议。
- **3 个核心只读 IoT Tools 契约**：覆盖设备列表检索、设备详情、时序历史点与统计。
- **强类型参数与 DTO 校验**：基于 `api/v1/mcp_tool.go` 与 `req.BindArguments(&params)` 进行入参解析与校验。
- **全链路 Hooks 与可观测性**：记录每次 Tool Call 的入参、出参、执行耗时与错误日志。
- **Wire 自动化依赖注入**：独立二进制入口 `cmd/mcp`，通过 Wire 完成数据库、TSDB 客户端、Service 与 Handler 的装配。

**Non-Goals:**
- 不在服务端实现 LLM 客户端与 Agent 协调循环（由外部 Cursor / Claude 客户端负责大模型调用与决策）。
- 不在首期实现动态加载外部三方 MCP Server（聚焦于作为提供方暴露 0things 自身工具）。
- 不在首期提供活动告警查询或设备控制；两者依赖的领域服务尚未落地，后续单独变更实现。

## Decisions

### 1. 纯 MCP 服务端分层架构（对齐 nunu-layout-mcp-main）

```text
       外部 AI 客户端 (Cursor / Claude Desktop / Dify)
                            │
               ┌────────────┴────────────┐
               ▼ (STDIO)                 ▼ (SSE / HTTP)
     ┌───────────────────────────────────────────────┐
     │              backend/cmd/mcp/                 │
     │  (独立二进制: 0things-mcp -conf config.yaml)  │
     └──────────────────────┬────────────────────────┘
                            │
                            ▼
     ┌───────────────────────────────────────────────┐
     │           backend/pkg/server/mcp/             │
     │  (Server 封装: Stdio/SSE 监听, 优雅停机)      │
     └──────────────────────┬────────────────────────┘
                            │
                            ▼
     ┌───────────────────────────────────────────────┐
     │          backend/internal/server/mcp.go       │
     │  (MCP 装配中心: setupSrv, newHooks, AddTool)  │
     └──────────────┬─────────────────┬──────────────┘
                    │                 │
                    ▼                 ▼
   ┌───────────────────────────┐  ┌───────────────────┐
   │ backend/internal/handler/ │  │ backend/api/v1/   │
   │ (mcp_tool.go: 业务逻辑)   │  │ (mcp_tool.go: DTO)│
   └────────────┬──────────────┘  └───────────────────┘
                │
                ▼ (复用现有底层)
   ┌───────────────────────────┐  ┌───────────────────┐
   │ GORM DAL (PostgreSQL)     │  │ pkg/tsdb (时序库) │
   └───────────────────────────┘  └───────────────────┘
```

### 2. 核心 IoT MCP Tools 详细定义

| Tool 名称 | 参数 (DTO) | 描述 | 返回值 |
| :--- | :--- | :--- | :--- |
| `iot_query_devices` | `product_id?`, `status?`, `keyword?`, `limit?` | 多维度查询设备列表及在线状态 | 设备简要列表 JSON |
| `iot_get_device_detail` | `device_key!` | 查询单台设备的物模型、标签与影子数据 | 设备完整详情 JSON |
| `iot_query_telemetry_history` | `device_key!`, `property!`, `start_time?`, `end_time?`, `limit?` | 查询指定属性的时序历史采样点与极值统计 | 采样点列表 + Min/Max/Avg 统计摘要 |

### 3. 配置定义 (`config.yaml`)

```yaml
mcp:
  name: "0things-iot-mcp"
  version: "1.0.0"
  stdio_enabled: true
  sse_enabled: false
  sse_addr: ":8008"
  streamable_http_enabled: true
  streamable_http_addr: ":8009"
```

## Risks / Trade-offs

- [时序数据点过多打爆 LLM 上下文] → `iot_query_telemetry_history` 限制最大点数为 100，并在 Handler 内计算好 `stats: { min, max, avg, latest }` 返回，帮助大模型直接获取结构化结论。
- [弱网超时] → MCP Server 封装中设置 10 秒请求超时上下文与优雅停止（Graceful Shutdown）。

## Migration Plan

1. 引入 `github.com/mark3labs/mcp-go`。
2. 实现 `backend/pkg/server/mcp/mcp.go`。
3. 实现 `backend/api/v1/mcp_tool.go` 与 `backend/internal/handler/mcp_tool.go`。
4. 实现 `backend/internal/server/mcp.go`。
5. 实现 `backend/cmd/mcp/main.go` 与 Wire 注入。
6. 本地配置 Cursor `mcpServers` 联调验证三个只读工具。
