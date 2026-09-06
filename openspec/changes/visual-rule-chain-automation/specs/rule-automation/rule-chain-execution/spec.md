## Purpose

以 ThingsBoard 的统一消息和具名关系模型异步、可追踪地执行已发布规则链。

## ADDED Requirements

### Requirement: 规范化并投递规则消息
系统 SHALL 将已持久化的遥测、属性、RPC、设备状态、告警、实体、分组、评论和受控 REST 事件转换为含 ID、类型、发生时间、组织、来源实体、产品、数据、元数据和关联 ID 的 `RuleMessage`，经持久化 Outbox 异步分发给范围匹配且已启用的已发布规则链。消息类型 SHALL 是受长度限制的字符串，以支持受控的 `custom.<integration>.<event>` 类型而无需数据库迁移。

#### Scenario: 匹配产品范围
- **WHEN** 产品内设备产生已持久化遥测
- **THEN** 系统 SHALL 仅为该产品范围内已启用规则链建立执行

### Requirement: 使用节点定义和关系执行图
数据库节点定义目录 SHALL 定义前端可见的配置 schema、端口、关系和分类；Go 执行器注册表 SHALL 白名单校验其 `executor_key`。执行器 SHALL 沿所有匹配关系扇出；条件至少声明 `True`、`False`、`Failure`，动作至少声明 `Success`、`Failure`。

#### Scenario: 条件 True 扇出
- **WHEN** 条件节点返回 `True` 且存在两条 `True` 边
- **THEN** 系统 SHALL 向两个下游节点投递同一规则消息上下文

#### Scenario: 节点失败
- **WHEN** 节点配置或动作执行失败
- **THEN** 系统 SHALL 写入错误事件并仅沿 `Failure` 边继续

### Requirement: 幂等、延迟与调试事件
系统 SHALL 以规则版本和消息 ID 保证执行幂等；每次节点尝试 SHALL 保留 IN、OUT、ERROR、消息类型、关系、耗时和脱敏摘要。延迟和重试 SHALL 存入可恢复计划任务；限流、去重、计数和聚合状态 SHALL 存入按节点和实体范围隔离的节点状态。

#### Scenario: 重复消息
- **WHEN** 同一消息再次到达同一规则版本
- **THEN** 系统 SHALL 复用既有执行或拒绝重复副作用，不得再次发送已成功非幂等动作

### Requirement: 支持受控节点与子链
系统 SHALL 为每条链提供不可编辑的系统入口，并提供 Filter（消息类型、来源实体、属性、阈值）、Enrichment（来源属性、关联实体、最新遥测）、Transformation（字段映射、变更来源）、Action（保存属性/遥测、创建/恢复告警、设备命令、RPC 回复）、External（Webhook）、Flow（延迟、限流、去重、消息计数、子链、结束）和 Analytics（窗口聚合、聚合值计算）节点。短信、邮件和 IM 通知动作不属于首期范围。子链 SHALL 仅调用同组织已发布链、返回终止关系并参与跨链循环校验。

#### Scenario: 子链返回父链
- **WHEN** 子链通过 `Success` 结束
- **THEN** 父链 SHALL 从子链节点的 `Success` 关系继续

### Requirement: 安全测试运行
系统 SHALL 用示例或最近脱敏消息模拟执行。模拟运行 MUST 不调用 Webhook、不下发命令、不产生外部子链副作用，但 SHALL 生成完整节点事件。

#### Scenario: 测试 Webhook
- **WHEN** 用户测试含 Webhook 节点的规则链
- **THEN** 系统 SHALL 记录模拟投递并明确标记为未发送
