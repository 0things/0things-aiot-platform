## Purpose

为生产环境提供自动续期的 HTTPS 证书和统一入口路由，同时保证测试环境不会意外暴露到公网，也不接管现有域名解析。

## ADDED Requirements

### Requirement: Production domain routing
生产部署 SHALL 通过 HTTPS 提供 `console.0thing.com`、`gateway.0thing.com`、`auth.0thing.com` 和 `auth-admin.0thing.com` 四个入口。gateway 的 `/api` 路径 SHALL 去掉 `/api` 前缀后到达 Backend，`/mcp` 路径 SHALL 原样到达 MCP 服务。

Backend 的 Logto audience 与前端请求令牌使用的 API Resource Identifier SHALL 均为 `https://gateway.0thing.com/api`，不能使用 gateway 根地址代替 Backend API 资源。

#### Scenario: Gateway API request
- **WHEN** 客户端请求 `https://gateway.0thing.com/api/v1/devices`
- **THEN** Backend 收到 `/v1/devices`

#### Scenario: Gateway MCP request
- **WHEN** 客户端请求 `https://gateway.0thing.com/mcp`
- **THEN** MCP 服务收到 `/mcp`

#### Scenario: Application and identity entrypoints
- **WHEN** 客户端访问 console、auth 或 auth-admin 域名
- **THEN** 分别到达 Frontend、Logto 用户端或 Logto 管理端

### Requirement: One automatically renewed production certificate
生产部署 SHALL 使用一张由免费 ACME CA 签发、自动续期且覆盖 `0thing.com` 和 `*.0thing.com` 的证书；全部生产 Ingress SHALL 引用同一 TLS Secret。生产 Chart SHALL 根据私有 prod values 创建 Cloudflare Token Secret 和专用的 ClusterIssuer。证书签发 SHALL 使用 Cloudflare DNS-01，而不能把 Cloudflare 凭据写入仓库。

#### Scenario: Initial issuance and renewal
- **WHEN** 有效的 Cloudflare API Token 已配置且证书尚不存在或临近到期
- **THEN** 证书控制器创建或续期共享的 TLS Secret，无须人工上传证书

### Requirement: Private test environment
test 环境 SHALL 不创建公网 Ingress、Certificate 或 ClusterIssuer。

#### Scenario: Render test manifests
- **WHEN** 使用 test values 渲染 Helm Chart
- **THEN** 输出不包含上述公网入口或证书
