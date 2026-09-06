## Purpose

为 0things 的已认证用户提供一个可在任何管理页面使用的 AI Copilot。用户可用自然语言查询自己组织内的设备、设备详情和遥测历史；模型只能经由受保护的 Streamable HTTP MCP 调用既有只读工具，不能在浏览器暴露模型凭据或跨组织读取数据。

## ADDED Requirements

### Requirement: 已认证页面提供全局 Copilot 入口

系统 SHALL 在所有已认证的管理端页面提供一致的 AI Copilot 入口。入口打开后显示可收起的对话抽屉、消息列表、输入框和“新建对话”操作，并使用 assistant-ui 的 `useChatRuntime` 驱动流式消息状态。

#### Scenario: 用户在设备管理页面打开 Copilot

- **WHEN** 已认证用户点击全局 Copilot 入口
- **THEN** 系统 SHALL 打开对话抽屉且不导航离开当前页面
- **AND** 用户 SHALL 能输入自然语言问题并提交

#### Scenario: 用户新建对话

- **WHEN** 用户在 Copilot 中点击“新建对话”
- **THEN** 系统 SHALL 清空当前浏览器内存中的消息和工具执行状态
- **AND** 系统 MUST NOT 删除或修改任何 IoT 业务数据

### Requirement: 聊天服务以 AI SDK 流式协议编排 MCP 工具

系统 SHALL 由独立的 Node.js AI Gateway 提供 `POST /v1/ai/chat`，并使用 Vercel AI SDK 的 UI stream 与 `@ai-sdk/mcp` 的 Streamable HTTP 客户端编排模型和 MCP 工具。浏览器 SHALL 仅连接聊天 API，MUST NOT 直接连接 MCP 服务或模型提供方。

#### Scenario: 模型需要查询设备列表

- **WHEN** 用户提问需要设备列表且模型选择 `iot_query_devices`
- **THEN** AI Gateway SHALL 经内部 MCP 地址调用该工具
- **AND** UI SHALL 在最终回答生成期间显示工具正在查询的状态
- **AND** 最终回答 SHALL 以流式文本显示在同一轮对话中

#### Scenario: 浏览器尝试直接调用 MCP

- **WHEN** 浏览器向 MCP Streamable HTTP 地址发起未受信任的请求
- **THEN** MCP 服务 MUST NOT 将该请求作为当前用户的工具调用执行
- **AND** 模型提供方 API Key MUST NOT 被发送至浏览器

### Requirement: MCP 查询必须按调用用户的组织隔离

系统 SHALL 将 AI Gateway 请求中的已验证用户身份安全传递至 MCP Streamable HTTP 请求。MCP 服务 SHALL 基于该可信身份构建组织上下文，并对三个 IoT 查询工具执行现有的组织范围限制。

#### Scenario: 用户查询本组织设备

- **WHEN** 属于组织 A 的用户通过 Copilot 查询设备或遥测数据
- **THEN** 工具调用 SHALL 仅返回组织 A 可访问的资源

#### Scenario: 用户猜测其他组织设备标识

- **WHEN** 属于组织 A 的用户要求查询组织 B 的设备标识或遥测数据
- **THEN** 系统 MUST NOT 返回组织 B 的任何设备详情、遥测值或存在性信息

### Requirement: 首期工具范围仅限三个只读 IoT 查询

系统 SHALL 只向模型注册 `iot_query_devices`、`iot_get_device_detail` 和 `iot_query_telemetry_history` 三个 MCP 工具。每个工具 MUST 保持只读；聊天服务不得根据模型输出执行写数据库、下发设备命令、修改规则链或发送外部通知的操作。

#### Scenario: 用户要求修改设备

- **WHEN** 用户要求 Copilot 修改设备属性、下发命令或删除资源
- **THEN** 系统 SHALL 明确说明当前 Copilot 只能查询
- **AND** 系统 MUST NOT 发起任何写操作工具调用

### Requirement: 失败信息对用户可恢复且不泄露内部细节

系统 SHALL 将模型服务、MCP 连接和工具调用失败转换为可理解、可重试的中文提示。系统 MUST NOT 在 UI 流中暴露访问令牌、模型密钥、内部 MCP 地址、数据库错误或未脱敏的堆栈信息。

#### Scenario: MCP 暂时不可用

- **WHEN** AI Gateway 无法连接 MCP 或工具调用超时
- **THEN** UI SHALL 显示“暂时无法查询设备数据，请稍后重试”的失败状态和重试入口
- **AND** 系统 MUST NOT 展示底层错误文本或凭据
