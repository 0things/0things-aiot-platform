## Why

当前平台同时维护本地密码、用户、组织、组织成员关系和自签发 JWT，认证与租户上下文存在两套来源，增加了登录、Token 校验和组织隔离的一致性成本。新系统不需要保留现有用户与组织数据，因此将身份与组织统一交给自托管 Logto，可以删除本地认证状态并让业务 API 直接依赖 Logto Organization Token。

## What Changes

- **BREAKING** 使用 Logto Hosted Sign-in 替代本地注册、登录和密码校验。
- **BREAKING** Backend 只接受并校验 Logto Access Token，不再签发本地 JWT。
- **BREAKING** 删除本地 `users`、`organizations`、`organization_users` 及其相关认证接口和业务逻辑。
- **BREAKING** 业务表的租户标识改为 Logto Organization ID 字符串，禁止继续依赖本地数字组织 ID。
- 增加自托管 Logto 与独立 PostgreSQL 服务，开发环境默认使用 `localhost:3001`、`localhost:3002`。
- 前端通过 `/callback` 接收 Logto 登录回调并保存 Access Token，Google 和 GitHub Provider 后续在 Logto Admin 中配置。
- 业务请求从 Logto Token 读取用户 `sub` 与组织声明，并以组织声明建立租户上下文。
- 清理现有数据后按新模型初始化，不设计旧用户、旧组织或旧 Token 的兼容迁移。

## Capabilities

### New Capabilities

- `identity/logto-auth`: 通过 Logto Hosted Sign-in 完成登录回调、Token 保存、Token 校验和 Google/GitHub Provider 配置边界。
- `identity/logto-organizations`: 以 Logto Organization 作为唯一组织来源，并使用 Organization Token 建立业务租户上下文。

### Modified Capabilities

- （无既有认证或组织 spec；本变更创建新的身份与组织能力。）

## Impact

- **Backend**：认证 middleware、handler、router、依赖注入、用户 service、JWT 包和租户上下文；删除本地认证与组织 API。
- **Backend model/storage**：删除本地身份组织模型；所有引用 `organization_id` 的业务模型、查询和迁移改为字符串。
- **Frontend**：删除本地登录/注册交互，增加 Logto callback、Token 管理和登录跳转；更新路由守卫及 API 请求拦截器。
- **Deployment**：新增 `logto/` 独立 Compose 管理 Logto 和专用 PostgreSQL；Google/GitHub 密钥只在 Logto Admin 或部署环境配置，不提交到仓库。
- **Data**：仅支持清空现有身份、组织及关联业务数据后重新初始化，不提供存量数据转换。
