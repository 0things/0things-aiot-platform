## 1. 数据模型与迁移

- [ ] 1.1 在 `backend/internal/model` 新增 `organization.go`（`Organization`：Id int64 PK、Name、Slug、时间戳、软删）与 `organization_user.go`（`OrganizationUser`：Id、OrganizationID int64、UserID string、`LastLoginAt *time.Time`、时间戳，唯一约束 (organization_id, user_id)），并确认 `TableName` 分别为 `organizations` / `organization_users`。验证：go build 通过。
- [ ] 1.2 在 `backend/internal/server/migration.go` 的 user DB `AutoMigrate` 中追加 `&model.Organization{}`、`&model.OrganizationUser{}`。验证：运行 migrate 后 user DB 存在两张表且无报错。

## 2. Repository 层

- [ ] 2.1 新增 `OrganizationRepository`（`Create`、`GetByID`、`ListByIDs`），置于 user DB repository。验证：单元测试或编译通过。
- [ ] 2.2 新增 `OrganizationUserRepository`：`Create(ctx, orgID, userID)`、`ListOrgIDsByUser(ctx, userID) ([]int64, error)`、`IsMember(ctx, userID, orgID) (bool, error)`。验证：单测覆盖唯一约束与成员判定。

## 3. Service 层（登录/注册/组织/切换）

- [ ] 3.1 改造 `Register`（`internal/service/user.go`）：在既有事务内，建 user 后创建个人 organization（name 取邮箱前缀或固定规则）并写 `organization_users`。验证：事务回滚测试 + 注册后查 `organization_users` 存在关联。
- [ ] 3.2 改造 `Login`：密码校验后调用 `ListOrgIDsByUser`（含 `last_login_at`）；若无归属则新建个人组织并关联（与注册一致）；默认 `OrganizationID` 取 `last_login_at` 最近的组织，全为 NULL 时回退 `min(org_id)`；并将该关联的 `last_login_at` 更新为当前时间，不再依赖 `tenant.GetOrganizationID` 的 mock 默认值。验证：单测覆盖单组织/多组织(取最近)/全 NULL 回退最小/密码错误。
- [ ] 3.3 新增 `ListMyOrganizations(ctx, userID) ([]*v1.OrganizationItem, error)`，标记 `is_current`（对比 `claims.OrganizationID`）。验证：返回列表含 `is_current` 且仅一个为真。
- [ ] 3.4 新增 `SwitchOrganization(ctx, userID, orgID) (string, error)`：经 `IsMember` 校验后调用 `jwt.GenToken` 重签（方案 A），并将该关联的 `last_login_at` 更新为当前时间。验证：单测覆盖所属(成功重签且 last_login_at 更新)/非所属(返回错误且不签发)。

## 4. API / Handler / 路由 / DI

- [ ] 4.1 在 `api/v1` 新增 DTO：`OrganizationItem`（id、name、is_current）、`SwitchOrgRequest`（org_id）、`SwitchOrgResponseData`（accessToken），并补 swag 注释。验证：`go build` 与 swag 生成无错。
- [ ] 4.2 在 `internal/handler/user.go`（或新 handler）新增 `ListMyOrganizations`、`SwitchOrganization` 方法，并使用 `StrictAuth`。验证：路由注册后编译通过。
- [ ] 4.3 在 gin 路由与 wire（`cmd/server/wire`）中注册新 handler/路由（`GET /organizations`、`POST /auth/switch-org`）并注入新 repository。验证：`make build` 通过。

## 5. Seed 脚本

- [ ] 5.1 更新 `backend/cmd/seed/main.go`：对用户 DB 增加 `organizations`(id=1/2/3) 与 `organization_users` 关联；用户密码改为对固定明文（如 `123456`）计算 bcrypt 哈希后写入（替换原假值）；并向 `organization_users` 写入错落的 `last_login_at`（如 user_001 最近用 org 3）以演示默认进入最近组织。验证：运行 seed 后，用 `123456` 调用 login 能拿到 token且默认 org 为最近登录组织。

## 6. 前端（左上角组织切换器）

- [ ] 6.1 在 `src/features/`（如 `org/` 或 `auth/`）新增 api 封装：`getMyOrganizations`、`switchOrganization`，并重新生成/补充 `src/api/generated`（若走 OpenAPI 则需 regenerate，勿手改生成文件）。验证：类型检查通过。
- [ ] 6.2 新增 org store（当前组织 id、组织列表、切换动作），选中态持久化到 localStorage；切换时调用 `switchOrganization` 替换 auth store 中的 token 并 `queryClient.invalidateQueries()` 刷新数据。验证：切换后设备列表随新 org 变化。
- [ ] 6.3 新增 `OrgSwitcher` 组件置于页面左上角，渲染组织下拉与当前标记；i18n 在 `public/locales/{zh,en}/` 同步新增文案。验证：`pnpm build` 与 `pnpm lint` 通过。

## 7. 联调与验证

- [ ] 7.1 端到端验证：migrate → seed → 用 seed 用户登录拿到 token → `GET /organizations` 返回列表 → `POST /auth/switch-org` 重签 → 新 token 查询设备仅见该 org 数据。验证：整链在本地跑通。
- [ ] 7.2 运行 `backend` `make test` 与 `frontend` `pnpm build && pnpm lint && pnpm format:check` 全部通过。
