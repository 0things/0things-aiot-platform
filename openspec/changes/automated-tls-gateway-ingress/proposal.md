## Why

生产环境目前通过 IP 和 LoadBalancer 暴露服务，缺少统一域名入口与自动续期的 HTTPS 证书。需要在现有 k3s/Traefik 集群上为 `0thing.com` 提供可重复部署的公网访问，同时保持 test 环境私有。

## What Changes

- 为生产环境配置 `console.0thing.com`、`gateway.0thing.com`、`auth.0thing.com` 和 `auth-admin.0thing.com` 的 Traefik Ingress；gateway 的 `/api` 转发给 Backend 并去掉前缀，`/mcp` 原样转发给 MCP。
- 使用 cert-manager、Let's Encrypt 和 Cloudflare DNS-01 申请一张覆盖 `0thing.com` 与 `*.0thing.com` 的证书并自动续期；生产 Ingress 共用该 TLS Secret。
- 保留现有 Cloudflare A 记录，由用户自行维护；cert-manager 仅通过 DNS-01 验证申请证书。
- 更新生产环境地址和部署说明；test 不启用公网 Ingress 或证书。

## Capabilities

### New Capabilities

- `deployment/public-ingress-tls`: 定义生产域名路由、自动证书，以及 test 环境的隔离行为。

### Modified Capabilities

None.

## Impact

- 涉及 `deploy/helm/0things` 的 Ingress、证书模板、环境 values 与文档，以及生产前端构建时的 Logto 地址。
- 集群需预装 Traefik，并独立安装 cert-manager；需要受限的 Cloudflare API Token 和 Let's Encrypt 联系邮箱。
- 仅新增/修改部署配置；不改变服务 HTTP API，不直接变更线上集群或 Cloudflare DNS。
