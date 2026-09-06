## Purpose

定义 0things 原生 IoT MCP (Model Context Protocol) Server 的实现规范，包含独立命令入口 `cmd/mcp`、通用包装器 `pkg/server/mcp`、核心 IoT 工具集契约以及多租户数据隔离安全机制。

## Requirements

### Requirement: 独立命令行入口与多传输协议支持
系统 SHALL 在 `backend/cmd/mcp` 提供独立二进制入口，基于 `mark3labs/mcp-go` 封装通用的 `pkg/server/mcp` 包装器，支持 STDIO 与 SSE 两种标准传输协议。

#### Scenario: 外部 IDE 本地连接
- **WHEN** Cursor 或 Claude Desktop 通过命令行参数拉起 `0things-mcp -conf config.yaml`
- **THEN** 系统 SHALL 启动 STDIO 监听，正确响应 MCP 初始化、工具列表查询（tools/list）与工具调用请求

### Requirement: 核心 IoT MCP 工具集定义
系统 SHALL 提供以下核心只读与控制工具：
- `iot_query_devices`：多条件检索设备列表与在线状态；
- `iot_query_telemetry_history`：查询设备时序指标点与极值/平均值统计；
- `iot_get_active_alerts`：获取当前组织内未解除的活动告警；
- `iot_send_device_command`：向设备下发属性设置或服务调用指令；
- `iot_generate_rule_chain`：根据自然语言生成规则链拓扑与节点配置 JSON。

#### Scenario: 查询时序遥测指标
- **WHEN** 调用 `iot_query_telemetry_history` 工具，传入合法的 `device_key` 与 `property`
- **THEN** 系统 SHALL 调用 TSDB 客户端获取时序采样点，并计算 Min/Max/Avg 统计摘要后以 JSON 字符串返回

### Requirement: 严格的多租户组织隔离与审计
所有 MCP 工具执行时 SHALL 从请求上下文中提取当前用户的 `organization_id`，并验证目标设备/告警是否属于当前组织。系统 SHALL 记录所有 Tool Call 的调用方、参数、耗时与执行结果。

#### Scenario: 跨组织越权拦截
- **WHEN** 用户尝试查询属于其他组织的设备数据
- **THEN** 系统 SHALL 拒绝执行并返回“设备不存在或无权限访问”错误

