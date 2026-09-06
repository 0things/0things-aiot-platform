## Purpose

将加密管理、可审计的 Webhook 配置提供给画布“调用 Webhook”动作。短信和其他通知提供方不属于本变更。

## ADDED Requirements

### Requirement: 管理 Webhook 节点配置和凭据
系统 SHALL 将 URL、方法、请求体映射、超时、重试和幂等策略保存于 Webhook 节点的公开配置；认证头、Bearer Token 和签名密钥 SHALL 加密保存于该节点的密文配置。凭据 MUST 不在规则图快照、导出 JSON、日志或 API 读取响应中明文或密文出现。

#### Scenario: 保存带凭据的 Webhook 节点
- **WHEN** 用户保存包含 Authorization 请求头的 Webhook 节点
- **THEN** 系统 SHALL 将该请求头写入节点密文配置，并在后续读取节点时只返回“凭据已配置”标记

### Requirement: 受控 Webhook 动作
系统 SHALL 通过平台 `WebhookSender` 调用节点配置的 URL。节点 SHALL 只允许受控的请求体映射和 HTTP 方法；系统 MUST NOT 将 `nikoksr/notify`、Novu 或第三方提供方类型引入本变更。

#### Scenario: 调用用户自建短信服务
- **WHEN** 用户把 Webhook 节点配置为调用自己的短信服务
- **THEN** 系统 SHALL 按 Webhook 安全策略投递，不管理其短信供应商密钥、模板或计费

### Requirement: 可靠安全投递
系统 SHALL 在执行与节点事件中为逻辑调用和每次尝试留记录，应用去重、超时、有限指数退避和错误脱敏。Webhook MUST 仅请求 HTTPS 公网目标，并拒绝本地、私网、保留地址和跨主机重定向。

#### Scenario: Webhook 失败
- **WHEN** Webhook 耗尽重试
- **THEN** 系统 SHALL 标记投递失败、保存脱敏错误，并让 Webhook 节点走 `Failure`

#### Scenario: 重复调用
- **WHEN** 同一节点在同一幂等键下重复提交
- **THEN** 系统 SHALL 复用已有逻辑投递而不得重复发送

### Requirement: 渲染和测试请求
系统 SHALL 在服务端使用受限变量映射渲染请求体，并支持受控终端测试。模拟规则测试 MUST 不出站；真实终端测试 SHALL 记录投递并遵守生产安全校验。

#### Scenario: 模拟 Webhook
- **WHEN** 测试含 Webhook 节点的规则
- **THEN** 系统 SHALL 显示渲染请求结果且不得真实调用终端
