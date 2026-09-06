## Purpose

定义 0things 前端与后端之间的 AI Copilot 对话交互规范，包含基于 assistant-ui 的 UI 组件体系、Vercel AI SDK 流式通信、会话管理以及敏感操作的 Human-in-the-loop 二次确认卡片机制。

## Requirements

### Requirement: 全局 AI Copilot 侧边栏与悬浮抽屉
系统 SHALL 在前端全局布局中提供可折叠、可全屏的 AI Copilot 抽屉容器，基于 `@assistant-ui/react` 实现打字机流式响应、智能触底滚动与 Markdown / 代码高亮渲染。

#### Scenario: 唤出对话抽屉
- **WHEN** 用户在任意页面点击右下角 AI 悬浮气泡或顶部导航栏的 AI 助手入口
- **THEN** 系统 SHALL 弹出侧边抽屉，加载当前组织上下文与历史对话记录，并展示快捷 Prompt 引导标签

### Requirement: 流式对话与 Tool Call 过程可视化
系统 SHALL 提供 `POST /v1/ai/chat` 接口，兼容 Vercel AI SDK Data Stream 协议。当大模型调用底层 MCP 工具时，前端 SHALL 实时展示工具执行状态与耗时。

#### Scenario: 只读数据查询
- **WHEN** 用户输入“查询 dk_device_047 最近 1 小时的温度变化”
- **THEN** 系统 SHALL 流式通知前端正在调用 `iot_query_telemetry_history`，并在工具返回后由大模型输出时序分析与统计总结

### Requirement: 高危动作与物理设备控制的 Human-in-the-Loop 二次确认
系统 SHALL 对所有涉及设备控制、状态修改、规则发布的写操作工具实施二次确认。后端不得直接执行指令，而是向对话流中推送确认请求，由前端渲染交互式确认卡片。

#### Scenario: 用户确认下发关机指令
- **WHEN** 用户指示“将设备 dk_device_047 关机”
- **THEN** 系统 SHALL 暂停直接下发，在聊天流中渲染包含目标设备、参数预览的确认卡片；当用户点击【确认执行】后，系统方可向设备下发指令并返回执行结果

#### Scenario: 用户取消操作
- **WHEN** 用户在确认卡片上点击【取消操作】
- **THEN** 系统 SHALL 取消待执行命令，并将卡片标记为已取消，不向底层发送任何控制指令

