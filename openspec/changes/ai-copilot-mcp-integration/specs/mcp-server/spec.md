## Purpose

定义 0things 原生 IoT MCP (Model Context Protocol) Server 的实现规范，包含独立命令入口 `cmd/mcp`、核心 IoT 工具集契约、多传输协议支持，以及面向多租户数据的安全访问边界。

## ADDED Requirements

### Requirement: 独立命令行入口与多传输协议支持
系统 SHALL 在 `backend/cmd/mcp` 提供独立二进制入口，支持 STDIO、SSE 和 Streamable HTTP 三种 MCP 传输协议。Streamable HTTP SHALL 作为远程服务端集成的标准传输协议。

#### Scenario: 外部 IDE 本地连接
- **WHEN** Cursor 或 Claude Desktop 通过命令行参数拉起 `0things-mcp -conf config.yaml`
- **THEN** 系统 SHALL 启动 STDIO 监听，正确响应 MCP 初始化、工具列表查询（tools/list）与工具调用请求

#### Scenario: 服务端通过 Streamable HTTP 连接

- **WHEN** 已配置的 MCP 客户端向 Streamable HTTP 端点初始化会话
- **THEN** 系统 SHALL 按 MCP Streamable HTTP 协议响应初始化、工具列表查询与工具调用请求

### Requirement: Streamable HTTP 调用身份认证

系统 MUST 要求每个 Streamable HTTP MCP 会话携带由 0things 签发的 Bearer JWT，并 SHALL 从经验证的 token 建立调用用户和组织身份上下文。缺少、无效或过期 token 的请求 MUST NOT 获得工具列表或执行工具，且系统 MUST NOT 回退至默认组织。STDIO 的本地开发行为不受本要求影响。

#### Scenario: 未认证客户端请求工具列表

- **WHEN** 客户端未携带有效 Bearer JWT 调用 Streamable HTTP 的初始化、`tools/list` 或 `tools/call`
- **THEN** 系统 SHALL 拒绝该请求并返回认证失败响应
- **AND** 系统 MUST NOT 返回任何工具定义、组织信息或业务数据

#### Scenario: 已认证客户端调用工具

- **WHEN** 客户端携带有效 Bearer JWT 通过 Streamable HTTP 调用任一已注册工具
- **THEN** 系统 SHALL 使用 token 对应的调用用户和组织身份处理该请求

### Requirement: 核心 IoT MCP 工具集定义
系统 SHALL 提供以下核心只读工具：
- `iot_query_devices`：多条件检索设备列表与在线状态；
- `iot_get_device_detail`：获取单台设备的物模型、标签与影子数据；
- `iot_query_telemetry_history`：查询设备时序指标点与极值/平均值统计；

#### Scenario: 查询时序遥测指标
- **WHEN** 调用 `iot_query_telemetry_history` 工具，传入合法的 `device_key` 与 `property`
- **THEN** 系统 SHALL 调用 TSDB 客户端获取时序采样点，并计算 Min/Max/Avg/Latest 统计摘要后以 JSON 字符串返回

### Requirement: IoT 工具按组织隔离数据

系统 SHALL 对 `iot_query_devices`、`iot_get_device_detail` 和 `iot_query_telemetry_history` 的每一次调用，使用经验证身份的组织范围执行查询。工具入参 MUST NOT 指定、覆盖或扩大组织范围；任何调用均 MUST NOT 返回当前组织无权访问的设备、设备详情或遥测数据。

#### Scenario: 查询当前组织设备

- **WHEN** 组织 A 的已认证客户端调用任一 IoT 查询工具
- **THEN** 系统 SHALL 仅返回组织 A 授权范围内的资源

#### Scenario: 猜测其他组织设备标识

- **WHEN** 组织 A 的客户端使用组织 B 的 `deviceKey` 调用设备详情或遥测工具
- **THEN** 系统 MUST NOT 返回组织 B 的资源、遥测值、物模型、影子状态或标签

### Requirement: 统一 Hooks 与可观测性
系统 SHALL 为每次 Tool Call 记录可审计的调用摘要：请求标识、已验证的调用用户与组织标识、工具名称、脱敏后的参数摘要、耗时和成功/失败结果类别。日志 MUST NOT 记录 Authorization header、Bearer JWT、模型密钥、完整工具返回 JSON、遥测采样点、设备影子、设备标签或设备元数据。

#### Scenario: 记录一次成功的工具调用
- **WHEN** 客户端成功调用任一 IoT MCP 工具
- **THEN** 系统 SHALL 记录调用用户与组织标识、工具名称、脱敏参数摘要、耗时和成功结果类别

#### Scenario: 工具返回敏感设备数据

- **WHEN** 任一 IoT 工具返回遥测采样点、设备影子、标签或元数据
- **THEN** 系统 MUST NOT 将这些原始内容写入应用日志
