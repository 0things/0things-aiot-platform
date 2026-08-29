## Context

当前设备模型使用 bigint `id` 作为数据库主键、`device_key` 作为业务标识，设备标签已存在，但没有分组实体或成员关系。前端设备详情仅有未接入导航的分组占位组件。实现需遵循现有 Gin → Service → Repository → GORM Model 分层，并重新生成 DAL、Swagger 和 Orval 客户端。

## Goals / Non-Goals

**Goals:**

- 建立组织隔离的扁平分组资源。
- 用服务端生成的 UUID 作为所有分组 API 的公开标识。
- 支持手动成员关系和安全的动态规则查询。
- 将分组接入设备列表筛选和设备详情展示。

**Non-Goals:**

- 本次不改变现有设备标签 API，也不把标签替换为分组。
- 本次不支持任意 SQL、脚本执行或复杂规则编排。
- 本次不实现分组级权限模型、设备批量导入或定时物化动态成员。

## Decisions

### 数据表和字段

新增以下两张表，表名和字段固定如下：

```text
device_groups
├── id              BIGINT      内部自增主键，不对外返回
├── group_uuid      UUID        Go 服务端生成的 UUID v4，对外唯一，NOT NULL、UNIQUE
├── organization_id BIGINT      所属组织，NOT NULL
├── name            VARCHAR(128) 分组名称，NOT NULL
├── type            VARCHAR(16) 分组类型：manual 或 dynamic，NOT NULL
├── description     TEXT        分组描述，可为空
├── rule            TEXT        动态规则表达式，manual 时为空
├── deleted_at      TIMESTAMP   软删除时间，可为空
├── created_at      TIMESTAMP   创建时间，NOT NULL
└── updated_at      TIMESTAMP   更新时间，NOT NULL

device_group_members
├── id              BIGINT      内部自增主键，不对外返回
├── group_id        BIGINT      device_groups.id，NOT NULL
├── device_id       BIGINT      devices.id，NOT NULL
├── created_at      TIMESTAMP   加入时间，NOT NULL
└── updated_at      TIMESTAMP   更新时间，NOT NULL
```

约束固定如下：

- `device_groups.group_uuid` 建立唯一索引。
- `device_groups.deleted_at` 建立索引；已软删除分组不参与查询和名称校验。
- 活跃分组的 `(organization_id, name)` 必须组织内唯一，由 Service 校验并配合查询索引实现。
- `device_group_members (group_id, device_id)` 建立唯一约束，禁止重复成员。
- `device_group_members.group_id` 和 `device_id` 分别关联分组和设备的内部主键；API 不接收这两个字段。
- 动态分组的 `rule` 必须有值，手动分组的 `rule` 必须为空。

对外分组对象只返回以下字段（组织由当前认证上下文确定，不返回 `organization_id`）：

```text
groupUuid
name
type
description
rule
deletedAt
createdAt
updatedAt
```

成员接口使用 `deviceKey`，不暴露 `device_id`：

```json
{
  "deviceKeys": ["dk_device_001", "dk_device_002"]
}
```

### 1. 内部主键与公开标识分离

新增 `device_groups` 使用 bigint `id` 作为内部主键，同时使用 Go 服务端通过 `github.com/google/uuid` 生成的 UUID v4 `group_uuid` 作为公开标识；`device_group_members` 使用内部 `group_id` 和 `device_id` 建立高效关联。DTO、路由参数和前端缓存只出现 `groupUuid`。UUID 必须在创建分组时由后端生成，客户端传入的 `groupUuid`/`group_uuid` 被忽略或拒绝，避免客户端伪造和主键枚举。

相比完全使用 UUID 外键，该方案保留现有数据库关联习惯和索引效率；相比直接暴露 bigint，则避免公开接口泄露内部数据规模。

### 2. 手动成员与动态规则分开处理

`type=manual` 的成员写入成员表，并通过 `(group_id, device_id)` 唯一约束防止重复；`type=dynamic` 只保存规则文本或规范化规则 JSON，不保存可被手动编辑的成员快照。查询动态分组设备时实时计算，保证产品、状态和标签变化立即生效。

动态规则采用受限表达式解析器：字段白名单为 `device_key`、`name`、`product_key`、`enabled`、`state`、`created_at`、`updated_at` 和 `tag.<key>`，操作符首期为 `=`、`!=`、`LIKE`、`IN`，逻辑连接为 `AND`、`OR`。解析结果转换为 Repository 查询条件，禁止拼接或执行原始 SQL。规则预览复用同一解析和查询路径，但不持久化。

相比直接保存 SQL，白名单 DSL 可校验、可审计且避免注入；相比只提供固定表单，表达式能覆盖截图中的组合条件，同时保留后续扩展空间。

### 3. 组织约束在 Service 层统一校验

创建和更新分组时，Service 按当前组织校验名称唯一。Repository 的所有分组和成员查询都带组织条件，避免仅依赖 Handler 参数校验。

删除分组只设置 `device_groups.deleted_at`，不删除设备；已删除分组的成员关系保留在成员表中但不再对外可见，也不参与动态分组查询。

### 4. API 资源边界

本次新增和修改的路由固定如下：

```text
POST   /device-groups                         创建分组
GET    /device-groups                         查询分组列表
GET    /device-groups/:groupUuid              查询分组详情
PUT    /device-groups/:groupUuid              更新分组
DELETE /device-groups/:groupUuid              软删除分组

GET    /device-groups/:groupUuid/devices      查询分组设备
POST   /device-groups/:groupUuid/devices      添加手动分组设备
DELETE /device-groups/:groupUuid/devices      移除手动分组设备

POST   /device-groups/preview                 预览未保存的动态规则
POST   /device-groups/:groupUuid/preview      预览已保存分组的动态规则

GET    /devices?groupUuid=:groupUuid          按分组筛选设备
GET    /devices/:deviceKey/groups             查询设备所属分组
```

所有 `:groupUuid` 均为 UUID 字符串；成员请求体使用设备 `deviceKey` 列表。`/device-groups/preview` 为静态路径，必须在 `/:groupUuid` 路由之前注册，避免被通配参数匹配。预览接口只返回匹配数量和设备列表，不修改分组数据。

这样分组是独立资源，设备仍以 `deviceKey` 为业务标识，不新增暴露设备内部 ID 的接口。

### 5. 前端交互

在设备管理下新增分组列表和详情页；创建/编辑使用抽屉，分组类型切换为“默认（手动）/动态”。动态类型显示规则编辑器、字段/操作符提示、搜索预览和重置操作；手动类型显示成员选择器。设备列表提供分组筛选，设备详情的分组页展示手动所属分组和动态命中分组。所有文案同时维护中英文资源，React Query key 使用 `groupUuid`。

## Risks / Trade-offs

- [动态规则查询成本不可控] → 限制字段和操作符，复用分页查询并限制预览结果；后续再评估索引和缓存。
- [规则语法升级导致旧规则无法解析] → 保存规则版本或规范化 JSON，并在读取时返回校验错误而不是执行不确定查询。
- [UUID 不能替代权限控制] → 所有 Service 查询继续强制注入当前组织条件，并统一返回 404/无权限语义。
- [动态成员不是快照] → 页面明确标识“实时匹配”，批量操作前提供当前匹配数量和列表预览。

## Migration Plan

1. 迁移创建 `device_groups` 和 `device_group_members`，不修改现有设备和标签表。
2. 部署后端分组 API 和设备查询的可选分组条件；旧设备数据无需迁移。
3. 部署前端分组页面、设备列表筛选和设备详情展示。
4. 回滚时先停用前端入口，再回滚 API；分组表可保留，不影响既有设备业务。
