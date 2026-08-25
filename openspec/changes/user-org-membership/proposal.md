## Why

当前登录链路里 JWT 的 `OrganizationID` 是空壳：登录时 `tenant.GetOrganizationID(ctx)` 因 context 无 tenant 而永远回退到 mock 默认值 `1`，用户的真实组织归属从未被解析。平台已有基于 `organization_id` 的租户数据隔离（`scope.Tenant`），但缺少 `organizations` / `organization_users` 实体与"用户↔组织"的关联，也没有多组织切换能力。需要补齐身份域的组织模型、把登录与组织归属打通，并支持前端左上角实时切换组织。

## What Changes

- 新增 `organizations` 表（int64 主键，与 device 表的 `organization_id` 共享键空间，id 对齐已有 seed 的 1/2/3）。
- 新增 `organization_users` 关联表（真多对多，一人可属多组织；暂不含 `role`）。
- 改造 `Register`：建 user → 建一个个人 organization → 写一条 `organization_users`。
- 改造 `Login`：校验密码后，从 `organization_users` 解析用户所属组织，默认取最小 `org_id` 写入 token（不再依赖 mock 默认值 `1`）。
- 新增 `POST /auth/switch-org {org_id}`：校验该用户确属该组织后，**重新签发 token**（方案 A——切换即换 token，统一且复用既有中间件与 `scope.Tenant`）。
- 新增 `GET /organizations`：返回当前用户所属组织列表，并标记 `is_current`（与当前 token 的 org 对比）。
- 前端新增左上角 `OrgSwitcher` 组件 + org store + 选中态持久化；切换时调用 `switch-org` 替换 token，并刷新数据。
- seed 脚本新增 organizations(1/2/3) 与 `organization_users` 关联，且用户密码改为**真 bcrypt** 哈希（当前 seed 为假值，登录校验会失败），使 demo 端到端可登录。

## Capabilities

### New Capabilities
- `identity/organizations`: 组织实体与"我的组织"查询（`GET /organizations`，返回所属组织列表及当前组织标记）。
- `identity/organization-membership`: 用户↔组织多对多关联；注册自动建个人组织并写入关联；通过 `switch-org` 在已属组织间重签 token 实现切换。

### Modified Capabilities
- （无既有 spec 需变更；登录行为变化归入上述新 capability 描述。）

## Impact

- **后端 model**：新增 `Organization`、`OrganizationUser`（位于 user DB / aiot-test.db，与 identity 同库）。
- **后端 repository**：新增 `OrganizationRepo` 与 `OrganizationUserRepo`（Create / ListByUser / IsMember）。
- **后端 service**：`user.go` 的 `Register`、`Login` 改造；新增 `ListMyOrganizations`、`SwitchOrganization`。
- **后端 middleware/jwt**：无需改（方案 A 复用 claims.OrganizationID）。
- **后端 handler**：新增 organizations / switch-org 路由（需接入 wire DI）。
- **后端 migration**：`AutoMigrate` 增加两张表。
- **后端 seed**：`cmd/seed/main.go` 增加 orgs 与关联，密码改真 bcrypt。
- **前端**：`src/features/` 下新增 org 相关 feature（api 封装、store、OrgSwitcher 组件），i18n 同步 zh/en；请求层从 auth store 读取 token（切换即替换）。
- **依赖/数据**：device 表 `organization_id` 键空间须与 `organizations.id` 对齐（seed 用 1/2/3），否则切换后看不到设备。
