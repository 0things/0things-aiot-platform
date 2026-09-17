## Purpose

为 AIoT 平台提供统一的外部身份认证，使用户注册、登录、密码管理和第三方身份源均由 Logto 负责，业务系统只消费并验证 Logto Access Token。

## ADDED Requirements

### Requirement: Logto Hosted Sign-in

系统 SHALL 将未认证用户引导至 Logto Hosted Sign-in，不得在平台前端收集或提交用户密码。

#### Scenario: User starts sign-in

- **WHEN** 未认证用户访问需要登录的页面
- **THEN** 前端将用户重定向到配置的 Logto Hosted Sign-in 地址

#### Scenario: User selects an external identity provider

- **WHEN** 用户在 Logto 登录页选择 Google 或 GitHub
- **THEN** Logto 完成 Provider 认证，并将授权结果回调到前端 `/callback`

### Requirement: Access token callback

前端 SHALL 在 `/callback` 处理 Logto 授权回调，保存 Access Token，并在后续 Backend 请求中以 Bearer Token 发送。

#### Scenario: Successful callback

- **WHEN** Logto 返回有效授权结果
- **THEN** 前端保存 Access Token、恢复认证状态并跳转到原始受保护页面

#### Scenario: Failed callback

- **WHEN** Logto 返回错误或授权结果无效
- **THEN** 前端不得保存无效 Token，并向用户显示登录失败状态

### Requirement: Backend token validation

Backend SHALL 只接受由配置的 Logto issuer 签发、面向配置 audience、签名有效且未过期的 Access Token。Backend MUST 从 Token 的 `sub` 读取用户标识。

#### Scenario: Valid Logto token

- **WHEN** 请求携带有效的 Logto Bearer Token
- **THEN** Backend 放行认证流程，并将 Token 用户标识放入请求上下文

#### Scenario: Invalid token

- **WHEN** 请求缺少 Token、Token 签名无效、issuer/audience 不匹配或 Token 已过期
- **THEN** Backend 返回 HTTP 401，并拒绝访问受保护资源

### Requirement: Local credential authentication is removed

系统 SHALL 删除本地注册、密码登录、密码校验和本地 JWT 签发能力；`users` 表不再作为认证来源。

#### Scenario: Local login endpoint is requested

- **WHEN** 客户端请求旧的本地登录或注册接口
- **THEN** 系统返回未实现或不存在的接口响应，且不得校验或保存本地密码

#### Scenario: Existing local user data is absent

- **WHEN** 新环境完成初始化
- **THEN** 系统不创建本地用户记录，用户身份以 Logto 的 `sub` 为唯一来源

### Requirement: External provider configuration boundary

Google 和 GitHub SHALL 作为 Logto Provider 配置，不得要求应用代码保存 Client Secret；未配置 Provider 凭据时，平台仍可启动，但对应登录方式不可用。

#### Scenario: Provider credentials are configured later

- **WHEN** 管理员在 Logto Admin 中补充 Google 或 GitHub 的 Client ID、Client Secret 和回调配置
- **THEN** 对应 Provider 可用于 Hosted Sign-in，而无需修改应用代码或提交凭据

