## 1. 后端 MCP 核心工具集与通用包装层

- [ ] 1.1 引入 `github.com/mark3labs/mcp-go` 依赖，创建 `backend/pkg/server/mcp/server.go` 通用服务包装器（支持 Stdio 与 SSE）。
- [ ] 1.2 创建 `backend/internal/mcp/` 模块，初始化 `IoTMCPServer` 并实现组织隔离拦截 Hooks。
- [ ] 1.3 实现只读 MCP 工具：`iot_query_devices`、`iot_query_telemetry_history` 与 `iot_get_active_alerts`。
- [ ] 1.4 实现写操作 MCP 工具：`iot_send_device_command`（命令下发）与 `iot_generate_rule_chain`（规则链辅助生成）。
- [ ] 1.5 编写 `backend/internal/mcp/` 单元测试，验证工具参数校验、租户隔离与返回值结构。

## 2. 独立 `cmd/mcp` 命令行服务

- [ ] 2.1 创建 `backend/cmd/mcp/main.go` 独立入口，支持通过 `-conf` 指定配置文件。
- [ ] 2.2 配置 `cmd/mcp/wire/wire.go`，完成数据库、TSDB Client 与 MCP Server 的依赖注入装配。
- [ ] 2.3 在 `Makefile` 中增加 `make build-mcp` 构建指令，验证 STDIO 本地拉起与外部 Cursor 接入。

## 3. 后端 AI 编排与 SSE 流式接口 (`cmd/server`)

- [ ] 3.1 创建 `backend/internal/ai/llm/client.go`，封装支持 OpenAI / DeepSeek / Qwen 规范的流式大模型客户端。
- [ ] 3.2 创建 `backend/internal/ai/agent/orchestrator.go`，实现 Tool Call 轮转循环与敏感写操作待确认状态拦截。
- [ ] 3.3 实现 `POST /v1/ai/chat` 接口，输出兼容 Vercel AI SDK Data Stream 协议的 SSE 事件流。
- [ ] 3.4 在 `backend/config/config.yaml` 补充 `ai` 与 `mcp` 默认配置项。

## 4. 前端 `assistant-ui` 与 Vercel AI SDK 对接

- [ ] 4.1 在 `frontend/package.json` 中引入 `@assistant-ui/react`、`@assistant-ui/react-ai-sdk` 与 `ai` 依赖。
- [ ] 4.2 创建 `frontend/src/features/ai-copilot/` 目录，封装基于 `useVercelUseChatRuntime` 的流式客户端。
- [ ] 4.3 构建 `<CopilotDrawer>` 抽屉组件与聊天容器，嵌入全局布局 `AuthenticatedLayout`。
- [ ] 4.4 完善 Markdown、表格、代码高亮与快捷 Prompt 引导标签。

## 5. 前端 Tool UI 与二次确认卡片

- [ ] 5.1 实现 MCP 工具调用进度组件（展示只读查询工具的调用状态与耗时）。
- [ ] 5.2 实现 `<DeviceCommandConfirmCard>` 交互式二次确认卡片（支持【确认执行】与【取消操作】）。
- [ ] 5.3 联调命令下发确认流程，验证确认后向设备发送指令并由 AI 总结结果的完整闭环。

## 6. 综合验证与规范检查

- [ ] 6.1 运行 `make test && make build` 确保 Go 后端与 `cmd/mcp` 编译测试通过。
- [ ] 6.2 运行 `pnpm -C frontend format && pnpm -C frontend build` 确保前端类型与生产构建通过。
- [ ] 6.3 验证 Cursor / Claude Desktop 本地通过 STDIO 连接 `0things-mcp` 工具集。

