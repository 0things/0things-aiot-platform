## Context

身份数据（users）与设备数据（products/devices/...）分属两个 SQLite 库：user DB（aiot-test.db）与 device DB（aiot-device.db）。组织隔离键 `organization_id`（int64）已贯穿 device 表并由 `scope.Tenant` 在查询时按 `tenant.GetOrganizationID(ctx)` 过滤；但登录签发 token 时该值来自 `tenant` 包的 mock 默认值 `1`，用户真实归属从未被解析。本次补齐 `organizations` / `organization_users`，并把"用户属于哪个组织"在登录与切换时真正算进 token。

## Goals / Non-Goals

**Goals:**
- 在 user DB 中建立 `organizations`、`organization_users` 两张表（真多对多）。
- 登录按用户归属解析 `OrganizationID`（默认取最小 org_id）。
- 提供 `switch-org` 重签 token 实现多组织切换（方案 A）。
- 提供 `GET /organizations` 供前端渲染左上角切换器。
- seed 用真 bcrypt 密码并写入组织关联，使 demo 端到端可登录。

**Non-Goals:**
- 不引入 `role` / 权限模型（后续再议）。
- 不实现组织邀请、组织 CRUD 管理界面（本次仅"我的组织"只读 + 注册自建）。
- 不采用 header 携带活跃组织的方案（已选方案 A：切换即换 token）。
- 不改动 `scope.Tenant` 与 middleware 的既有机制。

## Decisions

**D1. 组织实体放置于 user DB，与 users 同库。**
身份域（user/organization/membership）统一在 user DB；device DB 只保留 `organization_id` 作为共享数值键，不做跨库外键（SQLite 也不支持）。两库的"组织"靠相同的 int64 id 空间对齐。

**D2. `organizations.id` 用 int64 自增，与现有 device 数据回填值 1 及 seed 的 1/2/3 对齐。**
迁移曾把 device 表 `tenant_id` 改名 `organization_id` 并回填 `1`；seed 设备用 `rand.Intn(3)+1`（即 1/2/3）。因此 seed 的 orgs 必须包含 id 1、2、3，否则切换后看不到对应设备。

**D3. `users` 表不冗余组织字段，成员关系全在 `organization_users`。**
真多对多，避免 users 上出现误导性单 org 字段。登录/切换时一律查关联表解析。

**D4. 登录默认 org = 最近一次登录的组织（`organization_users.last_login_at` 最大者）。**
关联表新增 `last_login_at *time.Time` 字段：登录成功与 `switch-org` 切换组织时均将其更新为当前时间。`last_login_at` 为 NULL（首次登录、尚无任何登录记录）时，回退到 `min(org_id)`（即个人组织），保证确定性。该默认值仅用于"首请求/未切换"上下文，`GetOrganizationID → 1` 的 mock 回退保留给非登录公开路径。

**D5. 切换采用方案 A（重签 token），而非 header 覆盖。**
- 理由：完全复用既有 `middleware.StrictAuth` / `NoStrictAuth` 与 `scope.Tenant`（它们只读 `claims.OrganizationID`），无需改 `MyCustomClaims`、无需前端请求拦截器、无需每次校验 header 成员关系。
- 取舍：每次切换重新签发 JWT，前端需替换 token 并刷新数据；token 有效期仍为 90 天，切换频率低，开销可接受。

**D6. `switch-org` 在服务端强制校验成员关系。**
即便前端只展示"我的组织"，后端仍必须以 `organization_users` 校验 `org_id` 归属，防止越权签发他组织 token（防 403 边界）。

**D7. Registration 在事务内完成 user → org → membership。**
复用 `Service.tm.Transaction`（现有 `Register` 已用事务包裹 user 创建），把建组织与写关联并入同一事务，保证归属不丢。

**D8. seed 密码用真实 bcrypt，并向 `organization_users` 写入错落的 `last_login_at`。**
当前 `cmd/seed` 写的是 `hashed_password_xxx` 占位串，bcrypt 校验必失败；需对固定明文（如 `123456`）计算 bcrypt 哈希后写入，使 demo 能真实登录。同时给 `organization_users` 设置不同时间戳的 `last_login_at`（例如 user_001 最近用的是 org 3），以便 demo 中"默认进入最近使用的组织"行为可见。

## Risks / Trade-offs

- **[键空间错位]** 若 `organizations.id` 与 device 数据组织 id 不对齐，切换后设备列表为空。→ 缓解：seed 固定 orgs 为 1/2/3；迁移仅靠 AutoMigrate 建表，不回填 org（device 表已有数据假定属于 1/2/3）。
- **[切换后前端未刷新]** 方案 A 换了 token，若前端不刷新数据缓存，会显示旧组织内容。→ 缓解：OrgSwitcher 切换后调用 `queryClient.invalidateQueries()`（或等价刷新）再替换 token。
- **[事务回滚不一致]** 注册建组织/关联失败应整体回滚。→ 缓解：全部包进 `tm.Transaction`。
- **[无组织用户登录]** 理论上存在历史/异常用户无 membership。→ 缓解：登录时若无归属，按 D4 兜底或新建个人组织（实现时二选一并记录）；当前 fresh DB 不会出现。

## Migration Plan

1. `migration.go` 的 user DB `AutoMigrate` 增加 `&model.Organization{}`、`&model.OrganizationUser{}`。
2. 运行 migrate 建立两张表（空表，无历史数据需迁移）。
3. 运行更新后的 `cmd/seed` 写入 orgs(1/2/3) + 真 bcrypt 用户 + `organization_users` 关联。
4. 回滚：删除两张表即可，不影响 device 表；`organization_id` 仍由 mock 默认值 `1` 兜底（回归现状）。

## Open Questions

- 登录时若用户无任何组织归属，是返回错误还是自动建个人组织？（倾向后者，与注册一致；待实现确认。）
- 前端 OrgSwitcher 是否需要展示组织 `name`/`slug` 之外的元信息（如成员数）？本期不需要。
