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

