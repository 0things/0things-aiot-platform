## 背景

当前平台能接收遥测、事件和设备状态，但自动化仍是零散表单。用户无法直观看到消息为何触发、经过哪些判断、向哪里通知，也无法安全测试、发布或回退逻辑。

本变更以 ThingsBoard Rule Chain 为基准：统一消息进入有向图，节点按具名关系（`True`、`False`、`Success`、`Failure` 等）路由，节点拥有独立配置和调试轨迹。它不复制 ThingsBoard 的旧式 UI 或脚本泛滥；前端以 React Flow 和 shadcn/ui 提供现代体验。

## 变更内容

- 新增规则链草稿、校验、发布版本、停用、复制、回滚、导入导出与模板。
- 新增 React Flow 工作区：拖拽节点库、具名端口和连线、缩放、小地图、快捷操作、自动布局、配置抽屉、校验提示和运行覆盖层。
- 新增 ThingsBoard 风格 `RuleMessage`、数据库节点定义目录、Go 执行器白名单和有向关系执行器；生产仅执行不可变已发布版本。
- 新增输入/触发、筛选、转换、流控、告警、通知、设备命令、子规则链和结束节点。
- 新增通用“调用 Webhook”动作：URL、方法、请求体、重试和幂等策略保存在节点配置，认证信息加密保存于节点；用户可将规则事件交给自己的短信、IM 或业务服务。
- 新增告警、执行和节点事件记录，支持测试回放和生产调试；Webhook 尝试复用节点事件记录。
- 保留现有场景联动；未来仅提供显式复制为规则链草稿的迁移入口。

## 能力范围

### 新增能力

- `rule-automation/rule-chain-management`：草稿、发布版本、模板和审计。
- `rule-automation/rule-chain-execution`：消息、关系、节点注册、队列执行、子链和调试。
- `rule-automation/visual-editor`：画布、节点配置、测试、运行回放。
- `alerting/alert-lifecycle`：创建、聚合、确认、恢复和关闭告警。
- `notifications/notification-channels`：首期 Webhook 节点凭据、请求渲染和重试；短信是后续扩展。

### 修改能力

- 无。现有场景联动 API 和页面继续可用。

## 影响范围

- 后端新增模型、迁移、DTO、Handler、Service、Repository、节点注册表、消费者和 Worker，遵循 `Handler → Service → Repository`。
- 前端新增自动化功能域和 `@xyflow/react`；shadcn/ui 承担节点外观、抽屉、表单和反馈。
- Swagger/OpenAPI 与 Orval 客户端重新生成；新增 UI 同步中英文 i18n。
- 首期唯一的外部动作是安全 HTTP Webhook。短信、邮件和 IM 平台能力均不在本变更实现范围；用户可通过 Webhook 调用自建服务完成这些能力。
