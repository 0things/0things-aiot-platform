## 上下文

规则链借鉴 ThingsBoard 的运行模型，而不复制其旧前端：所有业务事件先成为统一 `RuleMessage`；画布上的节点按具名关系继续消息；已发布图不可变；调试可以回放消息经过每个节点的结果。FastBee 的“设备变量触发、条件、延迟、告警、变量执行”保留为面向用户的业务语言。

ThingsBoard 的真实实现中，规则图主要是 `rule_chain`、`rule_node`、`relation` 三类数据；节点类型是后端 Java 类名，并非数据库动作目录。0things 采用相同的图和消息模型，但额外持久化一个节点定义目录，使左侧节点库、配置表单和端口规则可以从数据库获取。定义目录只能描述和启停已安装能力，不能让数据库直接执行任意代码。

## 用户可见效果

```text
左侧节点库（后端目录）                 React Flow 画布                       右侧配置
Filter                  拖入      [消息入口]                              选中：调用 Webhook
Enrichment               ───▶        │                                     URL：告警服务
Transformation                         ▼ True                                请求体：高温告警 JSON
Action                             [温度 > 30℃] ─ False ─▶ [结束]
External                                  │
Flow                                      ▼
Analytics                           [创建高温告警] ─ Created ─▶ [调用 Webhook]
```

首期用户拖动的是“调用 Webhook”动作。它可调用用户自己的告警、短信、飞书或任何 HTTP 服务；URL、方法、请求体、重试和幂等策略就是节点配置，安全凭据加密保存在该节点。平台短信是后续独立扩展，不属于本期范围。

## 技术栈与边界

| 层级 | 技术 | 职责 |
|---|---|---|
| 管理端 | React 19、TypeScript、Vite、Tailwind | 路由、数据状态、编辑器和运行回放 |
| 画布 | `@xyflow/react`（React Flow） | 节点、端口、边、拖拽、选择、缩放、小地图；浏览器不执行规则 |
| 组件 | shadcn/ui、Radix、Lucide、React Hook Form、Zod | 节点卡片、抽屉、动态表单、菜单与无障碍 |
| 服务端 | Go、Gin、GORM、gormgen、Wire、Swagger | API、组织隔离、图校验、节点执行和持久化 |
| 异步 | 既有 Kafka/`franz-go`，本地 Worker 回退 | Transactional Outbox、规则执行、延迟恢复、Webhook 投递 |
| Webhook | 平台 `WebhookSender` | 安全的 HTTP 调用、签名、重试与投递审计 |

不引入 `nikoksr/notify`、Novu 或独立通知 SaaS。本期只提供平台自己的 `WebhookSender`；外部系统自行决定将收到的 HTTP 事件转为短信、IM 或其他通知。后续新增平台短信、邮件或 OAuth 型企业集成时，只新增动作定义和其领域表，不改变规则图表结构。

## 核心模型

### 统一消息信封

```text
RuleMessage {
  id, type, occurredAt,
  organizationId, originatorType, originatorUuid, productUuid,
  payload, metadata, correlationId, source
}
```

`type` 是受长度限制的字符串，不使用数据库枚举。平台定义标准消息类型，集成也可使用 `custom.<integration>.<event>`；消息类型筛选节点决定如何路由。这与 ThingsBoard 的 Message Type switch/filter 心智模型相同，新增事件不需要新增表或迁移。

`payload` 是事件数据，`metadata` 保存设备名、设备类型、请求 ID、RPC 方法、告警 ID、时间戳等上下文。敏感字段进入 Outbox、执行快照和调试日志前均按节点定义的脱敏规则处理。

### 节点目录、图和执行器

`rule_node_definitions` 是左侧节点库的唯一数据来源。它保存中文名称、说明、图标、分类、动态表单 schema、默认配置、输入输出端口和排序；前端只读取已启用且非系统定义。

每条定义的 `executor_key` 必须被 Go 内置执行器注册表白名单识别。数据库管理员可以隐藏节点、调整顺序和表单元数据，不能写入脚本、SQL、URL 模板以外的可执行逻辑，亦不能把未注册的 `executor_key` 投入生产。发布版本会冻结定义 key、定义版本和配置快照，避免后来修改节点目录影响已发布规则。

左侧目录的一级分类固定为 ThingsBoard 风格的 `Filter`、`Enrichment`、`Transformation`、`Action`、`External`、`Flow`、`Analytics`。分类 key 为稳定英文值，显示名称首期直接使用英文，不考虑国际化。每条链都有一个 `system.entry` 系统入口节点：由系统创建、不可删除、不出现在左侧库；它接收匹配范围的 `RuleMessage` 并从 `Success` 输出。用户要筛选遥测、属性、RPC 等消息类型时，使用 `Filter / 消息类型` 节点。

默认定义按以下 key 管理，不为每个动作建业务表。Seeder 同时写入 ThingsBoard 当前开源版的 77 个 `@RuleNode` 对应项；其配置先使用受限 JSON 编辑器，待相应 Go 执行器实现后再收紧为专用表单 schema。ThingsBoard 核心没有 `Analytics` 类，因此该分类保留 0things 的两个扩展节点。

| 分类 | 首期定义 key | 具名输出关系 |
|---|---|---|
| Filter | `filter.message-type`（消息类型）、`filter.originator-type`（来源实体）、`filter.property`（属性条件）、`filter.value-condition`（数值条件） | `True`、`False`、`Failure` |
| Enrichment | `enrichment.originator-attributes`（来源属性）、`enrichment.related-entity`（关联实体）、`enrichment.latest-telemetry`（最新遥测） | `Success`、`Failure` |
| Transformation | `transformation.map`（字段映射）、`transformation.change-originator`（变更来源实体） | `Success`、`Failure` |
| Action | `action.save-attributes`、`action.save-telemetry`、`action.create-alert`、`action.clear-alert`、`action.device-command`、`action.rpc-reply` | `Success`、`Failure`、`Created`、`Updated`、`Cleared` |
| External | `external.webhook`（调用 Webhook） | `Success`、`Failure` |
| Flow | `flow.delay`、`flow.rate-limit`、`flow.deduplicate`、`flow.message-count`、`flow.sub-chain`、`flow.end` | `Success`、`Failure`、子链终止关系 |
| Analytics | `analytics.aggregate`（时间窗口聚合）、`analytics.calculate`（聚合值计算） | `Success`、`Failure` |

### 图、版本和发布

草稿由节点、边和布局组成，可频繁保存；发布通过完整校验后生成不可变版本，生产只加载当前已发布版本。恢复历史版本总是建立新版本，绝不改写旧版本。

发布校验：恰有一个系统入口节点；用户节点 key 唯一；来源关系已声明且目标端口兼容；节点从入口可达；普通边和子链调用图均无环；可达非终止节点配置完整；范围、物模型、Webhook URL 和子链属于当前组织且可用；限制节点、边和一次执行的最大步数。

### 执行、状态和调试

```text
设备/HTTP/实体事件
        │  业务事务内写入
        ▼
rule_outbox_events（RuleMessage） ── Worker ──▶ 已发布版本
                                               │
                                               ├─ rule_executions
                                               ├─ rule_node_events
                                               ├─ rule_node_states / rule_scheduled_jobs
                                               └─ 告警、Webhook、设备命令等业务服务
```

- `(version_id, message_id)` 保证同一已发布版本对同一消息只建立一次执行。
- `(execution_id, node_key, attempt_no)` 记录每次节点 IN、OUT、ERROR、具名关系和耗时；非幂等动作用节点级幂等键保护。
- `rule_node_states` 保存限流、去重、计数、聚合等需要跨消息记忆的状态；`rule_scheduled_jobs` 保存延迟、重试和恢复任务。两者都不是画布配置。
- 测试运行标记 `simulation=true`；会产生完整节点事件，但不会真正调用 Webhook 或下发设备命令。

## 完整数据表设计

所有本模块表均有 `id BIGINT AUTO_INCREMENT PRIMARY KEY`、`created_at`、`updated_at`。对外资源另有不可猜测的 `uuid VARCHAR(36)`；HTTP 路径、DTO、前端状态、规则节点配置、导入导出一律只用 UUID，内部关联才用自增 ID。租户根资源都有 `organization_id` 并以组织范围查询；无权访问和资源不存在返回等价 not-found。JSON 均以 `TEXT` 保存并在写入前按 schema 校验。

### 1. 节点目录

| 表 | 核心字段 | 约束与职责 |
|---|---|---|
| `rule_node_definitions` | `uuid`, `key`, `version`, `category`, `name`, `description`, `icon`, `config_schema`, `default_config`, `input_ports`, `output_ports`, `executor_key`, `enabled`, `is_system`, `sort_order` | `key` 唯一；左侧节点库和配置抽屉的数据源；`executor_key` 必须由 Go 白名单注册；禁止删除已被发布版本引用的定义，只能停用或升版本 |

`config_schema`、端口和默认值由系统种子数据维护；不需要国际化列。本表是 0things 相比 ThingsBoard 的额外层，解决“用户从哪里看到所有动作”的问题。

### 2. 规则图与审计

| 表 | 核心字段 | 约束与职责 |
|---|---|---|
| `rule_chains` | `uuid`, `organization_id`, `name`, `description`, `enabled`, `current_draft_version_id`, `current_published_version_id`, `deleted_at` | 规则链根资源；组织内未删除名称唯一 |
| `rule_chain_scopes` | `rule_chain_id`, `scope_type`, `scope_uuid` | 一个链可绑定组织、产品、设备组或具体设备；`(rule_chain_id, scope_type, scope_uuid)` 唯一，用于分发匹配 |
| `rule_chain_versions` | `uuid`, `rule_chain_id`, `version_no`, `status`, `graph_snapshot`, `graph_checksum`, `published_by`, `published_at` | `(rule_chain_id, version_no)` 唯一；`draft` 可更新，`published`/`superseded` 不可更新；快照含节点定义版本 |
| `rule_nodes` | `version_id`, `node_key`, `definition_id`, `definition_key`, `definition_version`, `name`, `config`, `secret_config_ciphertext`, `credential_key_version`, `position_x`, `position_y` | `(version_id, node_key)` 唯一；Webhook 的 URL、方法、请求体和重试在 `config`，认证头/HMAC 密钥在密文列；`definition_key/version` 让历史图脱离目录后仍可解释 |
| `rule_edges` | `version_id`, `source_node_key`, `source_port`, `relation_type`, `target_node_key`, `target_port`, `label` | 保存 ThingsBoard 风格具名关系；禁止同一来源端口的重复目标边 |
| `rule_chain_audit_logs` | `rule_chain_id`, `version_id`, `actor_user_id`, `action`, `before_summary`, `after_summary` | 追加保存创建、保存、发布、停用、恢复、导入、测试操作；摘要不含密钥 |

### 3. 消息、运行和恢复

| 表 | 核心字段 | 约束与职责 |
|---|---|---|
| `rule_outbox_events` | `event_uuid`, `organization_id`, `message_id`, `message_type`, `originator_type`, `originator_uuid`, `product_uuid`, `source_type`, `source_uuid`, `occurred_at`, `payload`, `metadata`, `correlation_id`, `status`, `available_at`, `locked_at` | 事务内写入的统一消息信封和可靠 Outbox；`message_id` 全局唯一；Worker 原子领取 |
| `rule_executions` | `uuid`, `rule_chain_id`, `version_id`, `outbox_event_id`, `message_id`, `status`, `simulation`, `input_snapshot`, `correlation_id`, `started_at`, `finished_at` | `(version_id, message_id)` 唯一；一次消息命中多个链时每个版本各有一条执行 |
| `rule_node_events` | `execution_id`, `node_key`, `attempt_no`, `direction`, `message_type`, `relation_type`, `status`, `duration_ms`, `summary`, `error_message` | 调试回放的追加事件；`direction` 为 `IN`/`OUT`/`ERROR`；日志均脱敏 |
| `rule_node_states` | `version_id`, `node_key`, `scope_type`, `scope_uuid`, `state_key`, `state_data`, `expires_at` | `(version_id, node_key, scope_type, scope_uuid, state_key)` 唯一；服务于限流、聚合、去重、计数等有状态节点 |
| `rule_scheduled_jobs` | `execution_id`, `node_key`, `job_type`, `run_at`, `status`, `payload`, `idempotency_key`, `attempt_count`, `locked_at`, `last_error` | 延迟、重试、定时恢复任务；`idempotency_key` 唯一，Worker 重启可恢复 |

### 4. 规则动作产生的业务数据

| 表 | 核心字段 | 约束与职责 |
|---|---|---|
| `alerts` | `uuid`, `organization_id`, `device_uuid`, `rule_chain_id`, `rule_node_key`, `alert_type`, `severity`, `status`, `active_dedupe_key`, `count`, `latest_context` | 规则可创建、聚合、确认、恢复和关闭的告警；活动去重键唯一 |
| `alert_transitions` | `alert_id`, `execution_id`, `transition_type`, `actor_type`, `reason`, `context` | 追加式告警历史和规则执行关联 |
Webhook 调用不增加独立业务表：每次尝试复用 `rule_node_events` 的 `attempt_no`、`request_summary`、`response_status`、`response_summary` 和 `error_message`。设备命令、遥测、属性、实体关系、评论等动作也应调用其所属业务服务并复用那些领域的既有表；规则引擎不复制一套 `device_commands`、`telemetry` 或 `comments` 表。执行、节点事件和 `source_type/source_uuid` 保留可追溯关联。

### 5. 一条真实规则的完整数据样例

以下 14 条样例围绕同一个场景：组织 `42` 的设备 `dk_device_047` 上报温度 `36℃`，命中“温度高于 30℃”规则，创建高温告警，然后调用用户配置的 Webhook。省略所有表共同拥有的 `id`、`created_at`、`updated_at` 时会明确说明；真实 UUID 均用简化值表示以提高可读性。

#### `rule_node_definitions`：左侧“阈值判断”节点目录项

```json
{
  "id": 3,
  "uuid": "3a8d6b91-0c62-4e67-a9fc-5ef129b006a3",
  "key": "filter.value-condition",
  "version": 1,
  "category": "Filter",
  "name": "数值条件",
  "description": "按物模型数值决定 True 或 False 分支",
  "icon": "GitBranch",
  "config_schema": { "property": "string", "operator": [">", ">=", "<", "<="], "value": "number" },
  "input_ports": ["input"],
  "output_ports": ["true", "false", "failure"],
  "executor_key": "value-condition-v1",
  "enabled": true,
  "sort_order": 20
}
```

#### `rule_chains`：用户看到的规则链

```json
{
  "id": 101,
  "uuid": "49eafe9d-56f4-4d78-950a-4c3ea9bf9101",
  "organization_id": 42,
  "name": "设备高温告警",
  "description": "温度高于 30℃ 时通知值班系统",
  "enabled": true,
  "current_draft_version_id": 202,
  "current_published_version_id": 201,
  "deleted_at": null
}
```

#### `rule_chain_scopes`：只监听某个产品的设备

```json
{
  "id": 301,
  "rule_chain_id": 101,
  "scope_type": "product",
  "scope_uuid": "782f69c4-b788-4b2b-ae91-a9d6b4ef2001"
}
```

#### `rule_chain_versions`：生产正在使用的不可变版本

```json
{
  "id": 201,
  "uuid": "cb9f0546-09a5-41ac-bbc6-00c5608a0201",
  "rule_chain_id": 101,
  "version_no": 1,
  "status": "published",
  "graph_snapshot": {
    "nodes": ["telemetry-input", "temperature-over-30", "create-high-temp-alert", "call-oncall-webhook"],
    "edges": ["telemetry-input:Success", "temperature-over-30:True", "create-high-temp-alert:Created"]
  },
  "graph_checksum": "sha256:5c5f...d43e",
  "published_by": 7,
  "published_at": "2026-09-06T10:00:00+08:00"
}
```

#### `rule_nodes`：画布上的“调用 Webhook”节点

```json
{
  "id": 405,
  "version_id": 201,
  "node_key": "call-oncall-webhook",
  "definition_id": 12,
  "definition_key": "external.webhook",
  "definition_version": 1,
  "name": "通知值班系统",
  "config": {
    "url": "https://ops.example.com/hooks/device-alerts",
    "method": "POST",
    "body_template": {
      "event": "device.high_temperature",
      "deviceKey": "{{originator.key}}",
      "temperature": "{{message.temperature}}",
      "alertId": "{{alert.uuid}}"
    },
    "timeout_ms": 5000,
    "retry_count": 3,
    "idempotency_key_template": "{{message.id}}:high-temperature"
  },
  "secret_config_ciphertext": "enc:v1:7v7dJ...9KQ=",
  "credential_key_version": 1,
  "position_x": 900,
  "position_y": 120
}
```

读取这个节点时，API 不返回 `secret_config_ciphertext`，只返回 `credentials_configured: true`。

#### `rule_edges`：画布上的具名连线

```json
{
  "id": 502,
  "version_id": 201,
  "source_node_key": "temperature-over-30",
  "source_port": "true",
  "relation_type": "True",
  "target_node_key": "create-high-temp-alert",
  "target_port": "input",
  "label": "温度过高"
}
```

#### `rule_chain_audit_logs`：谁发布了版本

```json
{
  "id": 601,
  "rule_chain_id": 101,
  "version_id": 201,
  "actor_user_id": 7,
  "action": "published",
  "before_summary": { "draft_version": 1 },
  "after_summary": { "published_version": 1, "node_count": 4, "edge_count": 3 }
}
```

#### `rule_outbox_events`：设备上报后写入的统一消息

```json
{
  "id": 701,
  "event_uuid": "90f0ce43-c2f2-4c50-a7a7-dc54c3940701",
  "organization_id": 42,
  "message_id": "7b7df5e3-6636-456e-8071-909ef0200701",
  "message_type": "telemetry.posted",
  "originator_type": "device",
  "originator_uuid": "6a095d62-1f9e-4a85-a250-e99b96230470",
  "product_uuid": "782f69c4-b788-4b2b-ae91-a9d6b4ef2001",
  "source_type": "telemetry",
  "source_uuid": "d659d1f3-adf8-40ea-b0cd-8b4fc2fe0701",
  "occurred_at": "2026-09-06T10:30:05+08:00",
  "payload": { "temperature": 36.0 },
  "metadata": { "device_key": "dk_device_047", "device_name": "A 区温湿度传感器", "ts": 1788661805000 },
  "correlation_id": "mqtt:dk_device_047:1788661805000",
  "status": "processed"
}
```

#### `rule_executions`：该消息命中生产版本的一次执行

```json
{
  "id": 801,
  "uuid": "6f8d89d2-2527-4a72-8da8-5a025ccb0801",
  "rule_chain_id": 101,
  "version_id": 201,
  "outbox_event_id": 701,
  "message_id": "7b7df5e3-6636-456e-8071-909ef0200701",
  "status": "succeeded",
  "simulation": false,
  "input_snapshot": { "temperature": 36.0, "device_key": "dk_device_047" },
  "correlation_id": "mqtt:dk_device_047:1788661805000",
  "started_at": "2026-09-06T10:30:05.120+08:00",
  "finished_at": "2026-09-06T10:30:05.486+08:00"
}
```

#### `rule_node_events`：Webhook 的第二次尝试成功

```json
{
  "id": 903,
  "execution_id": 801,
  "node_key": "call-oncall-webhook",
  "attempt_no": 2,
  "direction": "OUT",
  "message_type": "telemetry.posted",
  "relation_type": "Success",
  "status": "succeeded",
  "duration_ms": 186,
  "request_summary": { "method": "POST", "host": "ops.example.com", "path": "/hooks/device-alerts" },
  "response_status": 202,
  "response_summary": { "accepted": true },
  "error_message": null
}
```

第一次请求若超时，也是在这张表中增加一条 `attempt_no: 1`、`direction: "ERROR"`、`error_message: "request timeout"` 的记录。

#### `rule_node_states`：一分钟内的去重状态

```json
{
  "id": 1001,
  "version_id": 201,
  "node_key": "deduplicate-high-temperature",
  "scope_type": "device",
  "scope_uuid": "6a095d62-1f9e-4a85-a250-e99b96230470",
  "state_key": "high-temperature",
  "state_data": { "last_message_id": "7b7df5e3-6636-456e-8071-909ef0200701", "count": 1 },
  "expires_at": "2026-09-06T10:31:05+08:00"
}
```

#### `rule_scheduled_jobs`：延迟节点暂存的任务

```json
{
  "id": 1101,
  "execution_id": 801,
  "node_key": "delay-5-minutes",
  "job_type": "resume-node",
  "run_at": "2026-09-06T10:35:05+08:00",
  "status": "pending",
  "payload": { "resume_relation": "Success" },
  "idempotency_key": "execution:801:node:delay-5-minutes",
  "attempt_count": 0,
  "locked_at": null
}
```

#### `alerts`：由告警节点创建或聚合的告警

```json
{
  "id": 1201,
  "uuid": "9d32ce1f-23a7-4e04-a769-03c0de7e1201",
  "organization_id": 42,
  "device_uuid": "6a095d62-1f9e-4a85-a250-e99b96230470",
  "rule_chain_id": 101,
  "rule_node_key": "create-high-temp-alert",
  "alert_type": "high_temperature",
  "severity": "critical",
  "status": "open",
  "active_dedupe_key": "device:dk_device_047:rule:101:type:high_temperature",
  "count": 1,
  "latest_context": { "temperature": 36.0, "threshold": 30.0 }
}
```

#### `alert_transitions`：告警创建历史

```json
{
  "id": 1301,
  "alert_id": 1201,
  "execution_id": 801,
  "transition_type": "created",
  "actor_type": "rule_engine",
  "reason": "temperature is greater than 30",
  "context": { "temperature": 36.0, "threshold": 30.0 }
}
```

### 6. 表之间的关系

```text
rule_chains ──< rule_chain_versions ──< rule_nodes >── rule_node_definitions
      │                    │                │
      └──< rule_chain_scopes│                └── rule_edges
                           │
rule_outbox_events ──< rule_executions ──< rule_node_events
                              │    │
                              │    └──< rule_scheduled_jobs / rule_node_states
                              ├──< alerts ──< alert_transitions
                              └── external.webhook 的每次尝试记录在 rule_node_events
```

## 与 ThingsBoard 消息类型的真实案例对比

ThingsBoard 将以下事件包装为 Message：属性/遥测、RPC、设备活动与连接、实体与分组变动、属性变动、告警、评论、REST API 和自定义事件。下表不是“只存得下”的假设，而是说明消息从产生到画布执行的实际落点。参考：[ThingsBoard Rule Engine Message Types](https://thingsboard.io/docs/pe/reference/rule-engine/message-types/)。

| ThingsBoard 消息类型 | 0things 等价消息与实际规则 | 数据落点 | 是否满足 |
|---|---|---|---|
| `POST_TELEMETRY_REQUEST` | 设备上报 `{temperature: 36}` 写入遥测后产生 `telemetry.posted`；`system.entry → filter.message-type → filter.value-condition → action.create-alert → external.webhook` | `rule_outbox_events`、执行/节点事件、`alerts`、Webhook 投递 | 满足；这是首期验收案例 |
| `POST_ATTRIBUTES_REQUEST`、`ATTRIBUTES_UPDATED`、`ATTRIBUTES_DELETED` | 设备或平台修改属性后产生 `attributes.posted`/`attributes.updated`/`attributes.deleted`；筛选 `mode=manual` 后写属性、告警或调用 Webhook | Outbox 保存变化前后值和来源；属性仍写原有物模型/属性表 | 满足；属性服务须在提交事务内写 Outbox |
| `TO_SERVER_RPC_REQUEST` | 设备调用 RPC 后产生 `rpc.request`，节点可条件分流、调用 Webhook 或 `action.device-command`/RPC 回复 | Outbox `metadata` 保存 request ID、method、params；命令/RPC 回应写既有设备服务 | 满足；需由 RPC 服务发布事件 |
| `RPC_CALL_FROM_SERVER_TO_DEVICE` | 规则命中后由 `action.device-command` 发起 `rpc.server-requested`，回执可再次进入规则链 | 节点事件关联既有命令/RPC 记录，`correlation_id` 串起请求和回执 | 满足；不复制设备命令表 |
| `CONNECT_EVENT`、`DISCONNECT_EVENT`、`ACTIVITY_EVENT`、`INACTIVITY_EVENT` | 连接状态变化产生 `device.connected`、`device.disconnected`、`device.active`、`device.inactive`；可做离线告警并调用 Webhook | Outbox + 执行；告警/Webhook 按需产生 | 满足；设备会话服务是事件源 |
| `ALARM`、`ALARM_ASSIGNED`、`ALARM_UNASSIGNED` | 告警创建/状态变化产生 `alert.created`、`alert.updated`、`alert.assigned`、`alert.unassigned`；可继续调用 Webhook | `alerts`、`alert_transitions`、Outbox | 满足；覆盖规则创建和人工处置两条路径 |
| `ENTITY_CREATED`、`ENTITY_UPDATED`、`ENTITY_DELETED`、`ENTITY_ASSIGNED`、`ENTITY_UNASSIGNED` | 产品、设备、资产等实体变化产生 `entity.*`；可自动初始化属性或调用 Webhook | Outbox 仅保存实体引用和变更摘要；实体由所属表维护 | 架构满足；每个实体服务接入 Outbox 后生效 |
| `ADDED_TO_ENTITY_GROUP`、`REMOVED_FROM_ENTITY_GROUP` | 设备加入/移出分组产生 `entity-group.added`/`entity-group.removed`；可重算权限、调用 Webhook 或触发子链 | Outbox + `rule_chain_scopes` 用于匹配该组规则 | 架构满足；分组服务发布事件后生效 |
| `REST_API_REQUEST` | 受控入口将 HTTP 请求标准化为 `rest.request`，通过图完成校验、转换、Webhook/设备命令，并按关联 ID 返回响应 | Outbox metadata 保存请求关联 ID；响应由 API 服务完成 | 满足；不是允许画布暴露任意未鉴权 HTTP 接口 |
| `COMMENT_CREATED`、`COMMENT_UPDATED` | 评论模块存在后产生 `comment.created`、`comment.updated`，可触发协同或调用 Webhook | Outbox 引用评论记录；规则引擎无需新增评论表 | 架构满足；当前评论模块若未实现，就没有事件源 |
| 自定义 Integration/Rule Node 消息 | 集成发布 `custom.vendor.event`，用 `filter.message-type` 路由到任意节点 | `message_type` 字符串、payload/metadata、节点定义目录 | 满足；无需新增消息表或动作表 |

因此，本设计在**消息种类、关系扇出、有状态节点、告警、外部动作、调试、版本回滚**这些 ThingsBoard 规则引擎能力上有对应承载模型。尚未实现的具体节点（例如脚本、Kafka、MQTT、邮件、复杂聚合）是 `rule_node_definitions + Go executor` 的实现增量，而不是数据表结构缺口；每增加一个节点必须新增定义、执行器、配置校验和测试。

## Webhook 动作

```text
external.webhook
        │  `rule_nodes.config`：URL、方法、请求体、重试、幂等键
        │  `secret_config_ciphertext`：认证头、Bearer Token、HMAC 密钥
        ▼
Worker ──▶ WebhookSender ──▶ rule_node_events（每次尝试）
```

- Webhook 仅允许 HTTPS 公网目标；解析后拒绝 loopback、私网、链路本地和保留地址，禁止跨主机重定向，并限制响应体、超时和退避。
- 节点配置保存 URL、方法、受限请求体映射、超时、重试和幂等策略；认证头和签名密钥只保存在 `rule_nodes.secret_config_ciphertext`。版本快照、导出 JSON 和读取 API 均不得返回该密文。
- 飞书、钉钉、企业微信机器人可以先以各自 Webhook URL 配置；格式辅助只是前端体验增强，不需要额外表或动作。

## API 设计

```text
GET                   /rule-node-definitions
GET/POST              /rule-chains
GET/PATCH/DELETE      /rule-chains/:uuid
POST                  /rule-chains/:uuid/copy
POST                  /rule-chains/:uuid/validate
POST                  /rule-chains/:uuid/test-runs
POST                  /rule-chains/:uuid/publish
POST                  /rule-chains/:uuid/disable
GET                   /rule-chains/:uuid/versions
POST                  /rule-chains/:uuid/versions/:versionUuid/restore
GET                   /rule-chains/:uuid/export
POST                  /rule-chains/import

GET                   /rule-executions
GET                   /rule-executions/:uuid
GET/POST              /alerts
POST                  /alerts/:uuid/ack | /resolve | /close
```

规则读取 API 仅返回“Webhook 凭据是否存在”，绝不返回密文或明文。真实 Webhook 测试使用同一安全校验并把每次尝试写入节点事件；模拟测试不出站。

## 前端工作区

- 路由：自动化 / 规则链、告警中心、运行记录；规则编辑页为独立路由。
- 编辑器先请求 `GET /rule-node-definitions`；按目录的分类、图标、端口和 schema 渲染左侧库、React Flow 节点和 shadcn 配置抽屉。
- 支持撤销/重做、复制/粘贴、框选、自动布局、Fit View、MiniMap、键盘访问和未保存离开确认。
- 顶部显式提供保存、校验、测试、发布；测试和生产执行都能以节点/边状态覆盖层回放。

## 非目标与迁移

- 不提供任意 JS/Groovy/SQL 脚本节点；后续需要时必须单独设计隔离沙箱、资源限制和审计。
- 普通边不得成环；延迟和重试使用计划任务；子链可复用但调用图无环。
- 不复制 ThingsBoard 前端，也不嵌入 Node-RED/FastBee。
- 先部署节点目录、图、Outbox、Worker 和模拟 Webhook；用“遥测阈值 → 告警 → 模拟/真实 Webhook”验证后，再为各业务服务逐项接入消息源。短信作为后续独立变更。规则链与既有场景联动并行，历史保持可读。
