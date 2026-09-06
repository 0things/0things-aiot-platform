## Why

随着大语言模型和 AI Agent（如 Cursor、Claude Desktop、Dify、Coze 等）的普及，Model Context Protocol (MCP) 已成为 AI 连接外部系统与工具的标准协议。

0things 作为 AIoT 平台，拥有设备元数据、TSL 物模型、时序遥测、活动告警与控制命令等丰富资产。通过实现一个标准、纯粹的 **0things IoT MCP Server**，可以让任意外部 AI 客户端（Cursor/Claude）通过 **STDIO** 本地直连或 **SSE** 远程网络直连，赋予 AI 实时感知与操作物理物联网世界的能力。

## What Changes

- **标准 MCP 通信协议封装（`pkg/server/mcp`）**：
  - 参考 `nunu-layout-mcp-main`，封装 `Server` 结构体，支持 `WithStdioSrv`（本地 IDE 进程管道）与 `WithSSESrv`（远程 HTTP SSE）以及优雅停机机制。
- **IoT 核心只读工具集（`internal/handler/mcp_tool.go` & `api/v1/mcp_tool.go`）**：
  - `iot_query_devices`：查询设备列表、产品归属与在线状态。
  - `iot_query_telemetry_history`：查询设备时序历史采样点与统计摘要（Min/Max/Avg/Latest）。
  - `iot_get_device_detail`：获取单台设备的完整详细信息（物模型、标签、影子状态）。
  - 活动告警与真实设备控制待对应领域服务具备后单独实现。
- **MCP 服务装配中心（`internal/server/mcp.go`）**：
  - 装配 `setupSrv`，配置全局 `newHooks`（拦截并记录 Tool Call 输入输出日志与错误）。
  - 注册工具列表（`AddTool`）。
- **独立命令行二进制入口（`cmd/mcp`）**：
  - 提供 `cmd/mcp/main.go`，支持 `-conf` 参数指定配置文件。
  - 基于 Google Wire 实现数据库、TSDB 客户端、Service、Handler 与 MCP Server 的自动化依赖注入。

## Capabilities

### New Capabilities
- `mcp-server`: 0things 原生 IoT MCP 服务端规范，包含 STDIO/SSE 双传输协议支持、5 大核心物联工具集契约与多租户组织数据隔离。

### Modified Capabilities
<!-- 无现有 Spec 需求变更 -->

## Impact

- **后端**：
  - 新增 `backend/pkg/server/mcp/mcp.go`。
  - 新增 `backend/api/v1/mcp_tool.go`。
  - 新增 `backend/internal/handler/mcp_tool.go`。
  - 新增 `backend/internal/server/mcp.go`。
  - 新增 `backend/cmd/mcp/main.go` 及 `backend/cmd/mcp/wire/`。
  - `config.yaml` 补充 `mcp` 配置块（`name`, `version`, `stdio_enabled`, `sse_enabled`, `sse_addr`）。
  - 依赖引入 `github.com/mark3labs/mcp-go`。
