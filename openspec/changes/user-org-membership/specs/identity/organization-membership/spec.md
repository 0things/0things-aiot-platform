## Purpose

定义用户与组织的多对多归属关系，并打通登录时的组织解析与运行时的组织切换（方案 A：切换即重签 token）。

## ADDED Requirements

### Requirement: 注册自动创建个人组织并关联
系统 SHALL 在用户注册时，先创建 user，再创建一个属于该用户的个人 organization，并向 `organization_users` 写入一条关联记录，使新用户至少属于一个组织。

#### Scenario: 注册成功建立归属
- **WHEN** 用户以新邮箱完成注册
- **THEN** 系统创建 user、创建一个个人 organization，并在 `organization_users` 写入 (user, 该 organization) 关联

#### Scenario: 邮箱已存在时拒绝
- **WHEN** 注册邮箱已存在
- **THEN** 系统返回邮箱已占用错误，不创建任何记录

### Requirement: 登录解析用户所属组织
系统 SHALL 在登录校验密码成功后，从 `organization_users` 解析该用户所属组织；默认 `OrganizationID` 取 `last_login_at` 最近的组织，当所有关联记录的 `last_login_at` 均为 NULL 时回退到最小的 `organization_id`（即个人组织）。登录成功时系统 SHALL 将该组织在 `organization_users` 中的 `last_login_at` 更新为当前时间。不再使用 `tenant.GetOrganizationID` 的 mock 默认值 `1`。

#### Scenario: 单组织用户登录
- **WHEN** 用户仅属于组织 ID=5
- **THEN** 签发的 token 中 `OrganizationID` 为 5，且该关联记录的 `last_login_at` 被更新

#### Scenario: 多组织用户登录取最近登录组织
- **WHEN** 用户属于组织集合 {2, 3, 7} 且 `organization_users` 中 org=3 的 `last_login_at` 最新
- **THEN** 签发的 token 中 `OrganizationID` 为 3

#### Scenario: 无任何登录记录时回退最小组织
- **WHEN** 用户所属组织的 `last_login_at` 均为 NULL
- **THEN** 签发的 token 中 `OrganizationID` 为这些组织中最小的 org_id

#### Scenario: 密码错误
- **WHEN** 邮箱存在但密码 bcrypt 校验失败
- **THEN** 系统返回 401 未授权，不签发 token，且不更新任何 `last_login_at`

### Requirement: 在所属组织间切换并重签 token
系统 SHALL 提供 `POST /auth/switch-org`，接收 `org_id`；仅当该用户通过 `organization_users` 确属该组织时，重新签发 token（新 token 的 `OrganizationID` 为该 `org_id`），并将该关联记录的 `last_login_at` 更新为当前时间，否则返回 403。

#### Scenario: 切换到所属组织
- **WHEN** 已登录用户对其所属组织 ID=3 调用 `POST /auth/switch-org`
- **THEN** 系统返回新的 accessToken，其 `OrganizationID` 为 3，且该关联的 `last_login_at` 被更新

#### Scenario: 切换到非所属组织被拒绝
- **WHEN** 已登录用户对其不归属的组织 ID=99 调用 `POST /auth/switch-org`
- **THEN** 系统返回 403，且不签发 token，且不更新任何 `last_login_at`

### Requirement: 关联表不含角色字段
系统 SHALL 在 `organization_users` 中仅持久化组织与用户的归属关系（含 `organization_id`、`user_id`、`last_login_at`、创建时间戳与唯一约束 `(organization_id, user_id)`），当前版本不存储 `role` 等角色字段。

#### Scenario: 关联唯一
- **WHEN** 同一用户对同一组织重复写入关联
- **THEN** 唯一约束阻止重复记录
