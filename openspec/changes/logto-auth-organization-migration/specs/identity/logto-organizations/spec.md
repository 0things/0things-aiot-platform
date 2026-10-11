## Purpose

以 Logto Organization 作为平台唯一组织来源，让业务请求直接依据 Organization Token 建立租户上下文，移除本地组织和成员关系的第二套状态。

## ADDED Requirements

### Requirement: Logto organization is the source of truth

系统 SHALL 使用 Logto Organization 作为组织的唯一来源，不得创建或维护本地 `organizations`、`organization_users` 组织主数据和成员关系。

#### Scenario: Organization is created or managed

- **WHEN** 管理员创建、修改或管理组织成员
- **THEN** 操作在 Logto 中完成，业务平台不创建本地组织副本

#### Scenario: New environment initialization

- **WHEN** 新环境执行初始化
- **THEN** 系统不创建本地组织或组织成员记录

### Requirement: Organization Token is required for tenant APIs

所有需要组织隔离的业务 API SHALL 要求 Logto Access Token 包含有效的 Organization 标识，并将该标识作为当前租户上下文。

#### Scenario: Request with organization token

- **WHEN** 请求携带有效且包含组织声明的 Logto Organization Token
- **THEN** Backend 使用 Token 中的组织 ID 执行产品、设备、规则链等业务查询和写入

#### Scenario: Request without organization context

- **WHEN** 请求携带有效但不包含组织声明的普通 User Token
- **THEN** 需要组织隔离的 API 返回 HTTP 401 或 HTTP 403，并且不得回退到固定组织或默认组织

### Requirement: Organization ID uses Logto string identifiers

所有业务租户字段 SHALL 使用 Logto Organization ID 字符串保存和查询，不得依赖本地自增数字组织 ID。

#### Scenario: Business data is written

- **WHEN** 用户创建或更新带组织归属的业务数据
- **THEN** 数据保存 Token 中的 Logto Organization ID 字符串

#### Scenario: Cross-organization access is attempted

- **WHEN** 用户使用组织 A 的 Token 请求组织 B 的业务数据
- **THEN** 查询返回空或拒绝访问，且不得因本地 ID 转换绕过租户隔离

### Requirement: Organization APIs are removed from the platform auth surface

系统 SHALL 删除依赖本地组织成员关系的组织列表和组织切换接口；组织选择和切换由 Logto 完成。

#### Scenario: Client requests local organization switching

- **WHEN** 客户端请求旧的 `/organizations` 或 `/auth/switch-org` 接口
- **THEN** 系统不再执行本地成员校验或重新签发平台 JWT

#### Scenario: User changes organization in Logto

- **WHEN** 用户在 Logto 流程中切换组织并获得新的 Organization Token
- **THEN** 后续业务请求使用新 Token 的组织上下文，而无需调用平台组织切换接口

### Requirement: Organization creation is explicit

系统 SHALL 仅在用户显式提交组织名称时创建 Logto Organization；登录或接受邀请不得自动创建组织。

#### Scenario: User creates an organization

- **WHEN** 已认证用户在应用的组织创建页面提交有效名称
- **THEN** Backend 在 Logto 创建该组织并将当前用户加入该组织

#### Scenario: User accepts an invitation without an organization

- **WHEN** 尚未加入组织的用户接受一个有效邀请
- **THEN** 用户加入邀请指定的组织，系统不得为该用户创建个人组织

### Requirement: Organization members and invitations use Logto as their source

组织成员和邀请管理 SHALL 通过 Backend 调用 Logto Management API；Backend 不得为这些数据创建本地副本。成员与邀请列表 SHALL 支持分页，成员列表 SHALL 支持关键词查询。

#### Scenario: List members or invitations

- **WHEN** 用户请求当前组织的成员或邀请列表并指定分页参数
- **THEN** Backend 返回当前组织对应的数据及分页信息

#### Scenario: Manage current organization membership

- **WHEN** 用户移除成员、撤销邀请或查询邀请详情
- **THEN** Backend 仅对当前组织内的相应 Logto 记录执行操作

### Requirement: Invitation acceptance uses a one-time token

系统 SHALL 通过 Logto 创建组织邀请和有效期为一小时的 one-time token，并由 Logto 邮件连接器发送包含应用接受页链接的邀请邮件。接受时必须校验登录用户主邮箱与邀请邮箱一致，并确认邀请处于未过期的待处理状态；成功后由 Backend 调用 Logto 接受该邀请。

#### Scenario: Invitee opens a valid invitation

- **WHEN** 被邀请者打开邮件链接并使用 one-time token 完成 Logto 认证
- **THEN** 应用恢复邀请接受流程，不要求用户创建新组织

#### Scenario: Invitee accepts the invitation

- **WHEN** 当前登录用户的主邮箱匹配有效邀请邮箱并确认接受
- **THEN** Backend 将 Logto 邀请状态更新为已接受，用户加入邀请所指定的组织

#### Scenario: Invitation is expired or belongs to another email

- **WHEN** 邀请已过期、已处理，或登录用户邮箱与邀请邮箱不匹配
- **THEN** Backend 拒绝接受邀请并保留邀请及组织状态不变

