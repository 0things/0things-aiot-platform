## 1. 后端 MCP 通用封装层 (`pkg/server/mcp`)

- [x] 1.1 引入 `github.com/mark3labs/mcp-go` 依赖。
- [x] 1.2 创建 `backend/pkg/server/mcp/mcp.go`，封装 `Server` 结构体及 `WithStdioSrv`、`WithSSESrv`、`WithStreamableHTTPSrv` 选项，实现优雅停机。

## 2. IoT 工具 DTO 与业务 Handlers

- [x] 2.1 创建 `backend/api/v1/mcp_tool.go`，定义 `QueryDevicesRequest`、`QueryTelemetryHistoryRequest`、`GetActiveAlertsRequest`、`SendDeviceCommandRequest` 等 DTO 结构。
- [x] 2.2 创建 `backend/internal/handler/mcp_tool.go`，实现 `MCPHandler` 接口与 3 个核心只读 IoT Tool 的具体执行逻辑。

## 3. MCP 服务装配中心 (`internal/server/mcp.go`)

- [x] 3.1 创建 `backend/internal/server/mcp.go`，实现 `setupSrv` 与 `newHooks`（拦截并记录 Tool Call 耗时、入参出参与错误日志）。
- [x] 3.2 在 `NewMCPServer` 中注册 3 个核心只读 IoT 工具（`AddTool`）。
- [x] 3.3 在 `backend/config/config.example.yml` 补充 `mcp` 相关配置项。

## 4. 独立 `cmd/mcp` 命令入口与 Wire 注入

- [x] 4.1 创建 `backend/cmd/mcp/main.go`，解析 `-conf` 启动参数与监听退出信号。
- [x] 4.2 创建 `backend/cmd/mcp/wire/wire.go`，配置 Wire 注入 ProviderSet 并执行 `wire` 生成 `wire_gen.go`。
- [x] 4.3 在 `backend/Makefile` 增加 `build-mcp` 构建目标。

## 5. 联调与验证

- [x] 5.1 编写 `backend/internal/server/mcp_test.go` 单元测试，验证参数绑定与工具返回结果。
- [x] 5.2 编译 `backend/bin/mcp` 二进制，使用 Stdio 模式进行本地调用测试。
- [x] 5.3 在 Cursor 或 Claude Desktop 的 `mcpServers` 配置 `0things`，验证 AI 客户端与平台的端到端工具调用。
