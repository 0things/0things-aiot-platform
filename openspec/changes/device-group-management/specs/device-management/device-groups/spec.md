## Purpose

为设备管理提供可按业务组织设备的分组能力，支持手动维护成员和基于设备条件实时计算成员的动态分组，并保证分组资源在多组织环境中的隔离与安全。

## ADDED Requirements

### Requirement: 分组数据结构固定

系统 SHALL 使用 `device_groups` 保存分组基本信息，使用 `device_group_members` 保存手动分组成员。分组表 SHALL 包含 `id`、`group_uuid`、`organization_id`、`name`、`type`、`description`、`rule`、`deleted_at`、`created_at` 和 `updated_at` 字段；成员表 SHALL 包含 `id`、`group_id`、`device_id`、`created_at` 和 `updated_at` 字段。`id`、`group_id` 和 `device_id` 仅用于内部关联，不得出现在公开 API 中。

#### Scenario: 返回分组公开字段

- **WHEN** 用户查询分组详情或分组列表
- **THEN** 系统返回 `groupUuid`、`name`、`type`、`description`、`rule`、`deletedAt`、`createdAt` 和 `updatedAt`，不返回 `organization_id` 或内部自增 ID

#### Scenario: 手动分组成员使用设备 Key

- **WHEN** 用户添加或移除手动分组成员
- **THEN** 请求使用 `deviceKeys`，系统根据设备 Key 建立内部成员关联，不要求调用方提供 `device_id`

### Requirement: 提供分组和设备分组路由

系统 SHALL 提供以下 HTTP 路由：`POST/GET /device-groups`、`GET/PUT/DELETE /device-groups/:groupUuid`、`GET/POST/DELETE /device-groups/:groupUuid/devices`、`POST /device-groups/preview`、`POST /device-groups/:groupUuid/preview`、`GET /devices?groupUuid=:groupUuid` 以及 `GET /devices/:deviceKey/groups`。所有分组路径参数 SHALL 使用 UUID，不得使用数据库自增 ID。

#### Scenario: 创建和查询分组

- **WHEN** 客户端调用 `POST /device-groups` 或 `GET /device-groups`
- **THEN** 系统创建或返回当前组织的分组，并在响应中使用 `groupUuid`

#### Scenario: 管理分组成员

- **WHEN** 客户端调用 `/device-groups/:groupUuid/devices` 的查询、添加或移除接口
- **THEN** 系统仅允许对手动分组操作，并使用 `deviceKey` 定位设备

#### Scenario: 预览动态规则

- **WHEN** 客户端调用未保存规则预览或已保存分组预览接口
- **THEN** 系统返回当前匹配的设备数量和列表，不修改分组或成员数据

#### Scenario: 按分组查询设备

- **WHEN** 客户端调用 `GET /devices?groupUuid=:groupUuid`
- **THEN** 系统按分组筛选设备并遵守原有分页和组织隔离规则

### Requirement: 分组资源使用公开 UUID

系统 SHALL 在 Go 后端为每个分组生成全局唯一的 UUID v4 `groupUuid`，所有公开 API、前端路由和响应 SHALL 使用 `groupUuid` 标识分组，不得返回或要求调用方传递数据库自增 `id` 或 `group_uuid`。

#### Scenario: 创建分组返回 UUID

- **WHEN** 用户创建一个分组
- **THEN** 系统生成 `groupUuid` 并在响应中返回，且该值不会暴露内部数据库主键

#### Scenario: 客户端不能指定分组 UUID

- **WHEN** 创建请求携带 `groupUuid` 或 `group_uuid`
- **THEN** 系统使用 Go 后端新生成的 UUID v4，不采用客户端提供的值

#### Scenario: 使用不存在的 UUID 查询分组

- **WHEN** 用户请求不存在的 `groupUuid`
- **THEN** 系统返回资源不存在错误，不返回其他分组数据

#### Scenario: 使用已删除分组 UUID 查询

- **WHEN** 用户请求已软删除分组的 `groupUuid`
- **THEN** 系统返回资源不存在错误，且该分组不出现在列表、成员和设备分组结果中

### Requirement: 分组删除使用软删除

系统 SHALL 通过设置 `deleted_at` 删除分组，不得删除设备数据。已删除分组的成员关系可以保留，但 SHALL 被所有公开查询过滤。

#### Scenario: 删除分组

- **WHEN** 用户删除一个分组
- **THEN** 系统设置该分组的 `deleted_at`，设备和其他分组不受影响

### Requirement: 支持手动和动态分组

系统 SHALL 支持 `manual` 和 `dynamic` 两种分组类型。手动分组 SHALL 通过成员关系维护设备；动态分组 SHALL 保存筛选规则并根据当前设备数据计算成员，动态分组不得通过手动成员接口修改成员。

#### Scenario: 管理手动分组成员

- **WHEN** 用户向手动分组添加或移除设备
- **THEN** 系统更新该分组的成员关系，并在查询分组设备时反映变更

#### Scenario: 拒绝修改动态分组成员

- **WHEN** 用户调用成员添加或移除接口操作动态分组
- **THEN** 系统拒绝请求并提示动态分组成员由规则计算

### Requirement: 动态规则必须安全可校验

系统 SHALL 对动态规则进行语法和字段校验，只允许平台定义的设备字段、标签字段和操作符。系统 MUST 将规则解释为查询条件，不得直接执行用户提交的 SQL。

#### Scenario: 预览有效动态规则

- **WHEN** 用户提交包含 `product_key`、`device_key`、`name`、`enabled`、`state` 或 `tag.<key>` 的有效规则
- **THEN** 系统返回匹配设备数量和设备列表，不修改已保存的分组规则

#### Scenario: 拒绝非法动态规则

- **WHEN** 规则包含未支持字段、未支持操作符或无法解析的表达式
- **THEN** 系统返回明确的校验错误，不执行查询且不保存规则

### Requirement: 分组数据按组织隔离

系统 SHALL 按当前组织隔离分组、分组成员和动态查询结果。用户不得通过修改 `groupUuid` 访问其他组织的分组或设备。

#### Scenario: 跨组织访问分组

- **WHEN** 用户使用属于其他组织的 `groupUuid` 请求分组或成员
- **THEN** 系统返回资源不存在或无权限错误，且不泄露分组信息

### Requirement: 设备管理集成分组

系统 SHALL 支持在设备列表按分组筛选，并在设备详情展示设备所属的手动分组和命中的动态分组。设备可以属于多个手动分组。

#### Scenario: 按分组筛选设备

- **WHEN** 用户在设备列表选择一个分组
- **THEN** 系统只返回该分组当前可见的设备，并保留原有分页和权限约束

#### Scenario: 查看设备所属分组

- **WHEN** 用户打开设备详情的分组页面
- **THEN** 系统分别展示设备所属的手动分组和根据规则命中的动态分组
