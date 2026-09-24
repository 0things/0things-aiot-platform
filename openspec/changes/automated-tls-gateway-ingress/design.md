## Context

参见 proposal.md。现有 Chart 仅支持 frontend、Logto 的可选 Ingress，默认 class 为 nginx；生产 values 仍使用公网 IP 和 LoadBalancer。k3s 自带 Traefik，test 与 prod 位于不同 namespace。

## Goals / Non-Goals

**Goals:** 部署配置可重复渲染，生产通过指定域名访问并自动维护 TLS；test 仍仅供内部访问。

**Non-Goals:** 不维护 Cloudflare A 记录，不直接操作线上 k3s/Cloudflare，不改变应用 API，不迁移已有数据库。

## Decisions

1. cert-manager 控制器使用 `deploy/helm/cluster` 下的 k3s `HelmChart` 清单引用官方 Chart，由 k3s Helm Controller 作为独立的集群级 release 管理，避免 test/prod 重复安装控制器和 CRD。专用于 0things 的 `ClusterIssuer` 和 Cloudflare Token Secret 由应用 Chart 根据仅在 prod 提供的配置创建；卸载 prod release 会删除它们。应用 Chart 还声明单个 `Certificate`、Ingress 和 Traefik `Middleware`。
2. `Certificate` 显式管理 `0thing.com` 与 `*.0thing.com`，四个 Ingress 引用同一个 namespace 内的 TLS Secret；不在每个 Ingress 上配置 cert-manager ingress-shim 注解，以免多个资源争夺同一证书。
3. gateway 使用一个 Ingress，分别将 `/api` 转发到 Backend、`/mcp` 转发到 MCP。Ingress 附 Traefik StripPrefix Middleware，其唯一匹配前缀为 `/api`，因此 `/mcp` 路径保持原样。
4. prod values 使用 Traefik、ClusterIP Service 和 HTTPS public URLs；Ingress 不设置 enabled 开关，配置了域名且服务启用时生成。test 不配置公网域名，也不启用 Certificate。`gateway.0thing.com` 是 Ingress 主机，Backend API 的 Logto Resource Identifier 是 `https://gateway.0thing.com/api`；Backend audience 与前端构建资源地址必须一致。同源 `/api` 和 `/copilot` 仍由前端 Nginx 代理。

## Risks / Trade-offs

- [Cloudflare Token 权限过大或泄露] → 仅授予该 zone 的 DNS Edit/Zone Read；真实值仅写入本地忽略的 prod values。Helm release 记录会保存 Token，限制读取该 release 的权限。
- [现有 Cloudflare A 记录可能未覆盖全部入口] → 部署前由用户自行确认四个域名都指向服务器；本变更不会创建或修改 A 记录。
- [首次签发受 DNS 传播/ACME 限额影响] → 先安装控制器，再安装包含 Token Secret、Issuer 与 Certificate 的应用 Chart；检查 Certificate/Challenge 就绪后切流。
- [Traefik CRD 版本随 k3s 变化] → 部署前确认集群有 `traefik.io/v1alpha1` Middleware CRD。
- [不同 namespace 不能直接共享 TLS Secret] → Certificate 与所有生产 Ingress 均位于 `0things-prod`。

## Migration Plan

1. 准备限定权限的 Cloudflare API Token，并在集群独立安装 cert-manager。
2. 先渲染检查生产和测试清单，再升级生产 Chart，由其创建 Token Secret、ClusterIssuer 和 Certificate。
3. 检查证书 Ready、各域名路由；确认现有 A 记录及前端镜像中的生产 Logto 地址。
4. 回滚时恢复此前 values/Helm revision；独立安装的 cert-manager 不随应用 release 自动删除，但 prod release 中的 ClusterIssuer 会随其卸载，现有 A 记录不受影响。
