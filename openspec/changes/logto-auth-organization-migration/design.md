## Context

当前 Backend 使用本地 `users`、`organizations`、`organization_users` 和自签发 JWT；前端也保留本地登录状态。目标环境允许清空现有数据，因此不需要兼容旧用户、旧组织或旧 Token。新部署通过独立的 `logto/` Compose 管理 Logto 与专用 PostgreSQL。

## Goals / Non-Goals

**Goals:**

- 让 Logto 成为用户身份、登录会话和组织成员关系的唯一来源。
- 让 Backend 只验证 Logto Organization Token，并从声明建立租户上下文。
- 让前端使用 Hosted Sign-in 和 `/callback`，不再处理密码。
- 让所有业务租户字段使用 Logto 的字符串 Organization ID。

**Non-Goals:**

- 不迁移现有用户、组织或业务数据。
- 不在应用中实现 Google/GitHub 的 OAuth 协议细节。
- 不在 Backend 运行时调用 Logto Management API 查询每次组织成员关系。
- 不在本次变更中实现微信或其他身份 Provider。

## Decisions

### 1. Logto Hosted Sign-in + OIDC callback

前端使用 Logto 的浏览器 OIDC 客户端完成登录跳转、回调处理、Token 缓存和登出。应用只配置 Logto endpoint、app ID、回调地址和登出地址；Provider 的 Client Secret 留在 Logto Admin 或部署环境。

选择 Hosted Sign-in 而不是自建登录表单，是为了避免应用重新接管密码、Provider 回调和会话安全。Google/GitHub 凭据延后配置不影响代码和容器启动。

### 2. Backend 使用 Logto verifier，不保留双 Token 兼容期

认证 middleware 通过 Logto issuer 的 JWKS 验证 Bearer Token，校验签名、issuer、audience、过期时间和 Organization 声明。认证成功后，将 `sub` 和组织 ID 写入请求上下文；组织 API 不再回退到固定组织。

选择直接移除本地 JWT 兼容逻辑，是因为用户已确认清空旧数据，不需要双认证来源或迁移窗口。这样可以避免同一请求在两套 claims 之间产生歧义。

### 3. Logto Organization ID 直接作为业务租户键

删除本地组织主数据和成员关系表；所有业务模型中的 `organization_id` 改为字符串列，并让 Repository 使用请求上下文中的 Logto Organization ID 过滤。

选择直接改键而不是保留本地映射表，是因为本次允许清空数据，映射表只会重新引入第二个组织来源。数据库迁移按全新初始化处理，不提供旧数字 ID 转换。

### 4. 独立 Logto PostgreSQL

`logto/docker-compose.yml` 单独定义 `logto` 与 `logto-postgres`，Logto 数据不与 AIoT 业务数据库共享。Logto Compose 使用与平台相同的外部 `0things-net`，开发环境由根 Compose 先创建网络。

选择独立数据库是为了隔离 Logto schema、升级周期和备份边界；Logto 官方示例中的初始化命令保留在 Logto 容器 entrypoint 中。

### 5. Remove local auth surface

删除本地 `/login`、`/register`、`/organizations`、`/auth/switch-org` 及其 service、repository、model 和 JWT 签发依赖。保留业务 API 的认证上下文读取，但读取 Logto claims。

## Risks / Trade-offs

- [Logto 不可用会导致所有登录失败] → 将 Logto 与 PostgreSQL 纳入 Compose health/dependency，并在前端显示明确的登录服务错误。
- [Organization Token 误用普通 User Token] → Backend 对组织隔离 API 强制要求组织声明，禁止默认组织回退。
- [所有 organization_id 改为字符串影响面较大] → 先统一模型、Repository、租户上下文和 API DTO，再执行全新数据库初始化；用全仓库搜索确认没有数字比较和强制转换。
- [Provider 凭据尚未准备] → 先完成 Logto Hosted Sign-in、回调和 Token 验证；Google/GitHub 在 Logto Admin 配置后再进行端到端登录验证。
- [删除本地认证接口是 breaking change] → 前后端同一版本发布，旧前端不再作为兼容客户端支持。

## Migration Plan

1. 部署 Logto 和独立 PostgreSQL，完成 Logto 应用、回调地址和 Organization 配置。
2. 配置开发环境 endpoint：Logto `3001`、Admin `3002`、前端回调 `http://localhost:5173/callback`。
3. 清空并重新初始化 AIoT 业务数据库；不创建本地用户、组织和成员关系表。
4. 发布 Backend 的 Logto verifier、字符串租户键和旧认证接口删除。
5. 发布前端 Hosted Sign-in、callback、Token 管理和路由保护。
6. 在 Logto Admin 中配置 Google、GitHub Provider 后验证登录和组织 Token。

回滚策略：不支持在同一数据集上回滚到本地认证，因为旧用户、组织和数字租户键已被删除。若需要回滚，只能恢复变更前的数据库备份和应用版本。
