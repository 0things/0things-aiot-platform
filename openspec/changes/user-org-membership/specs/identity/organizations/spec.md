## Purpose

让已登录用户查询自己所属的组织，并标识当前活跃组织，支撑前端左上角组织切换器。

## ADDED Requirements

### Requirement: 查询我的组织列表
系统 SHALL 提供 `GET /organizations`，返回当前用户通过 `organization_users` 关联的所有组织，每个组织包含 `id`、`name` 以及 `is_current` 标记（`is_current` 为真表示与当前 token 的 `OrganizationID` 一致）。

#### Scenario: 用户属于多个组织
- **WHEN** 已登录用户调用 `GET /organizations` 且其通过 `organization_users` 关联了 3 个组织
- **THEN** 响应返回这 3 个组织，且仅其中一个的 `is_current` 为真

#### Scenario: 未携带有效 token
- **WHEN** 请求未携带有效 token
- **THEN** 系统返回 401 未授权

### Requirement: 组织标识与设备数据键空间对齐
系统 SHALL 保证 `organizations.id` 为 int64 且其在 `organization_users` 中关联出的 `organization_id` 与 device 数据表（products/devices/...）的 `organization_id` 处于同一数值空间；切换组织后，该组织的设备数据按既有租户隔离规则可见。

#### Scenario: 切换到组织后看到该组织设备
- **WHEN** 用户切换到组织 ID=2 的 token 后查询设备列表
- **THEN** 仅返回 `organization_id = 2` 的设备（由 `scope.Tenant` 过滤）
