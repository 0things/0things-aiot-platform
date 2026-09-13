# 规则链数据表与作用

本文说明当前代码中已经落地的规则链数据表、字段作用和后续扩展方向。表名和字段以 `backend/internal/model` 中的 GORM 模型为准；规划中的表不会被当作当前数据库结构使用。

## 一、当前已落地的数据表

### 1. `rule_node_definitions`：节点定义目录

该表是规则链编辑器左侧节点面板和节点配置表单的数据来源。它描述“节点能做什么”，不直接执行脚本或 SQL。

| 字段 | 作用 |
| --- | --- |
| `id` | 内部自增主键，仅供数据库内部关联使用 |
| `uuid` | 节点定义的公开标识，唯一 |
| `key` | 稳定的节点类型标识，例如 `filter.value-condition`，唯一 |
| `version` | 节点定义版本；规则图保存时应记录使用的版本 |
| `category` | 节点分类，例如 `Filter`、`Action`、`External` |
| `name` / `description` | 节点显示名称和说明 |
| `icon` | 前端显示图标名称 |
| `config_schema` | 节点配置的 JSON Schema |
| `default_config` | 拖入画布时使用的默认配置 |
| `input_ports` / `output_ports` | 输入、输出端口及具名关系定义 |
| `executor_key` | 后端执行器注册表中的白名单 key |
| `enabled` | 是否允许新编辑使用该节点 |
| `is_system` | 是否为系统定义；系统入口等特殊节点使用此标记 |
| `sort_order` | 同分类下的显示顺序 |
| `created_at` / `updated_at` | 创建和更新时间 |

约束：`uuid` 和 `key` 唯一。停用节点定义不会影响已保存图中的历史定义引用。

### 2. `rule_chains`：规则链根资源

该表对应列表中的一条规则链，负责保存规则链的基本信息和当前草稿图。当前实现尚未拆分版本、范围和执行记录表。

| 字段 | 作用 |
| --- | --- |
| `id` | 内部自增主键，不通过 API 暴露 |
| `uuid` | 规则链公开标识；详情、编辑和删除路由使用该值 |
| `organization_id` | 所属组织，用于租户隔离 |
| `name` | 规则链名称；同一组织内未软删除记录不能重复 |
| `description` | 规则链说明 |
| `status` | 当前状态，现阶段主要使用 `draft` / `published` |
| `version` | 当前图版本号；每次保存更新时递增 |
| `graph` | React Flow 图数据，包含节点、边和布局；以 JSON 存储 |
| `created_at` / `updated_at` | 创建和更新时间 |
| `deleted_at` | GORM 软删除时间；删除后列表和查询默认不返回 |

#### 名称唯一规则

Service 在写入前通过 Repository 查询保证：

```text
UNIQUE (organization_id, name) WHERE deleted_at IS NULL
```

因此：

- 同一组织内，两个未删除规则链不能使用相同名称；
- 不同组织可以使用相同名称；
- 软删除后，原名称可以重新创建；
- Service 会先去除名称首尾空格，空名称直接拒绝；
- 名称校验属于规则链业务校验，不把数据库驱动错误文本泄漏到业务层。

## 二、当前保存链路

```text
HTTP Handler
    ↓ 解析请求、组织 UUID 和图数据
RuleChainService
    ↓ 校验名称、图结构、状态和版本
RuleChainRepository
    ↓ 组织范围查询和写入数据库
rule_chains
```

节点目录的读取链路为：

```text
rule_node_definitions → Repository → Service → Handler → 前端节点面板/配置表单
```

前端只展示后端返回的节点定义，不在前端补造缺失的节点或名称。

## 三、规则图 JSON 的职责

`rule_chains.graph` 保存画布状态，典型结构如下：

```json
{
  "nodes": [
    {
      "id": "entry",
      "type": "ruleNode",
      "position": { "x": 80, "y": 120 },
      "data": { "definitionKey": "system.entry", "config": {} }
    }
  ],
  "edges": []
}
```

- `nodes` 保存节点 ID、类型、位置以及节点配置；
- `edges` 保存来源节点、目标节点和关系/端口信息；
- 图 JSON 是规则链草稿的一部分，不能把凭据明文写入其中；
- 浏览器只负责编辑和展示，不能在前端执行规则；
- 当前后端已校验 JSON 合法且至少存在一个节点，完整的端口兼容、可达性和循环校验属于后续执行能力建设。

## 四、后续拆分建议

当规则链从“草稿 CRUD”进入发布和运行阶段，建议按职责新增以下表。它们不是当前迁移已创建的表，落地时应新增 GORM 模型、迁移、DAL 和测试。

| 表 | 主要作用 |
| --- | --- |
| `rule_chain_scopes` | 保存规则链作用范围，例如组织、产品、设备组或设备 |
| `rule_chain_versions` | 保存不可变发布快照、版本号、发布人和发布时间 |
| `rule_nodes` | 从版本快照中拆出节点，保存定义版本、配置和位置 |
| `rule_edges` | 从版本快照中拆出边，保存具名关系和端口连接 |
| `rule_chain_audit_logs` | 记录创建、保存、发布、停用、恢复、导入和测试操作 |
| `rule_outbox_events` | 在业务事务内保存统一规则消息，供 Worker 异步分发 |
| `rule_executions` | 记录某个发布版本处理某条消息的执行实例 |
| `rule_node_events` | 记录节点输入、输出、错误、关系和耗时，用于调试回放 |
| `rule_node_states` | 保存限流、去重、计数和聚合等跨消息状态 |
| `rule_scheduled_jobs` | 保存延迟、重试和恢复任务 |

后续拆分的关键约束：

1. 规则链和版本对外只使用 UUID，内部自增 ID 只用于数据库关联。
2. 发布版本不可变，恢复历史版本应创建新的版本记录。
3. `rule_executions` 应以 `(version_id, message_id)` 做幂等约束。
4. 节点密钥、Webhook 凭据等敏感数据必须单独加密保存，不能进入图 DTO、日志或导出文件。
5. 规则动作调用既有设备、告警和消息服务，不重复创建同类业务表。

## 五、迁移入口

当前规则链相关模型通过 `MigrateServer.Start` 的 GORM `AutoMigrate` 接入服务启动迁移：

- `RuleNodeDefinition` → `rule_node_definitions`
- `RuleChain` → `rule_chains`

测试环境若单独使用 `AutoMigrate`，也必须包含对应模型，才能验证唯一索引和软删除行为。
