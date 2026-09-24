# 0things Helm Chart

该 Chart 用于在 Kubernetes 中部署完整的 0things AIoT 平台，并通过独立 values 文件管理 test 和 prod 环境。

## 部署内容

Chart 默认部署以下组件：

- 应用：backend、MCP server、data-engine、transport-mqtt、AI copilot、frontend、Logto。
- 基础设施：业务 PostgreSQL、Logto PostgreSQL、Redis、NATS JetStream、EMQX、TDengine。
- 配置：服务运行配置、数据库凭据、应用凭据、镜像仓库凭据和前端 Nginx 配置。
- 网络：内部 Service，以及可选的 frontend、gateway API/MCP、Logto 和 Logto Admin Ingress。

模板按部署模块组织：

```text
templates/
├── application/
├── infrastructure/
├── secrets/
└── networking/
```

## 前置条件

- Kubernetes 1.24 或更高版本。
- Helm 3。
- 已安装可用的 StorageClass；如果集群没有默认 StorageClass，需要设置 `global.storageClass`。
- 已构建并推送 values 中配置的业务镜像。
- 已在目标数据库中准备应用所需表结构。该 Chart 不执行数据库 migration。
- 生产公网入口依赖 k3s Traefik 的 `traefik.io/v1alpha1` Middleware CRD 和 cert-manager；先按下文安装集群级组件。

## Values 管理

| 文件 | 用途 |
| --- | --- |
| `values.yaml` | 完整配置结构和公共默认值 |
| `values-test.yaml` | test 的账号密码、域名、资源和存储覆盖 |
| `values-prod.yaml` | prod 的账号密码、公网地址、副本、资源和存储覆盖 |

环境文件包含部署凭据，因此已加入 `.gitignore`，不会随仓库分发。首次部署时，在仓库根目录运行以下命令创建本地文件，再按环境修改账号密码、域名、资源和存储设置：

```bash
cp -n deploy/helm/0things/values.yaml deploy/helm/0things/values-test.yaml
cp -n deploy/helm/0things/values.yaml deploy/helm/0things/values-prod.yaml
```

`-n` 会保留已有的本地环境文件。Helm 会自动加载 `values.yaml`，部署时再追加对应的环境文件。

仓库中的本地 test values 默认不创建镜像仓库 Secret；如果镜像是私有的，部署前须配置 `registry` 凭据或引用已有的 `imagePullSecrets`。

> [!IMPORTANT]
> `values-prod.yaml` 中的 `change-me-*` 或 `REPLACE_WITH_*` 是占位值，生产部署前必须替换。真实生产 values 可以保存在受控的私有仓库或服务器文件中。

数据库和中间件凭据由 Chart 自动创建为 Kubernetes Secret，并与相应组件定义放在同一个模板文件。镜像仓库拉取凭据单独定义在 `secrets/registry.yaml`；如果 Secret 已由平台预先创建，可通过顶层 `imagePullSecrets` 引用。

各项目的 ConfigMap 与工作负载定义放在同一模板文件：Backend/MCP 共用 `application/backend.yaml` 中的 `backend-config`；Data Engine、MQTT Transport 和 Frontend 的 ConfigMap 分别放在各自的 `application/*.yaml` 中。前三份项目配置的键名均为 `local.yml`。这些配置可能包含密码和应用密钥；ConfigMap 内容不是加密存储，请限制集群访问权限，并避免提交真实凭据。

私有镜像仓库默认使用以下配置自动创建 `kubernetes.io/dockerconfigjson` Secret：

```yaml
registry:
  create: true
  server: ccr.ccs.tencentyun.com
  username: your-user
  password: your-password
```

如果 Secret 已由平台统一管理：

```yaml
registry:
  create: false
imagePullSecrets:
  - tencent-registry
```

## 安装

### Test

```bash
helm upgrade --install 0things-test deploy/helm/0things \
  --namespace 0things-test \
  --create-namespace \
  -f deploy/helm/0things/values-test.yaml \
  --atomic \
  --wait \
  --timeout 10m
```

### Production

先准备好镜像、数据库结构及本地生产 values，替换占位凭据（包括 `clusterIssuer.apiToken`）：

```bash
rg 'change-me-|REPLACE_WITH_' deploy/helm/0things/values-prod.yaml
```

确认数据库结构已经准备完毕后执行：

```bash
helm upgrade --install 0things-prod deploy/helm/0things \
  --namespace 0things-prod \
  --create-namespace \
  -f deploy/helm/0things/values-prod.yaml \
  --atomic \
  --wait \
  --timeout 15m
```

资源使用固定组件名；test 和 prod 必须部署到不同 namespace。Helm release 名称可以区分版本记录，但不会改变资源名。

## 外部访问

配置了对应 `ingress.*.host` 且服务启用时，会生成 Ingress；未配置域名的 test 环境不会生成。frontend 使用 `ingress.frontend`，gateway 使用 `ingress.gateway`，Logto 使用 `ingress.logto`，管理控制台使用 `ingress.logtoAdmin`。生产路由为：

| 域名与路径 | 后端 | 路径处理 |
| --- | --- | --- |
| `console.0thing.com/` | frontend | 原样转发 |
| `gateway.0thing.com/api/...` | backend | 去掉 `/api` 前缀 |
| `gateway.0thing.com/mcp` | mcp-server | 原样转发 |
| `auth.0thing.com/` | Logto app | 原样转发 |
| `auth-admin.0thing.com/` | Logto admin | 原样转发 |

生产 Ingress 共用固定名称的 `aiot-platform-tls`，由一张覆盖 `0thing.com` 与 `*.0thing.com` 的 Let's Encrypt 证书提供。证书由 cert-manager 自动续期；Cloudflare DNS-01 仅用于域名所有权验证，不创建或修改现有 A 记录。四个入口的 A 记录由你自行维护。test values 不配置公网域名或证书资源，不申请公网证书。

### 生产集群准备（只执行一次）

按 [独立的集群级部署说明](../cluster/README.md) 确认 Traefik 并应用 `deploy/helm/cluster/cert-manager.yaml`。k3s Helm Controller 维护 cert-manager 控制器；prod 0things release 根据本地 values 创建 Cloudflare Token Secret 和 `ClusterIssuer`，test 不创建。卸载 prod release 会删除该 Secret 和 Issuer，但不会卸载 cert-manager 控制器。

前端镜像须用生产 Logto application ID 构建：`docker build -f frontend/Dockerfile --build-arg VITE_LOGTO_APP_ID=YOUR_LOGTO_APP_ID -t YOUR_FRONTEND_IMAGE frontend`。`frontend/.dockerignore` 会排除 `.env.production`，因此必须传这个构建参数。`VITE_LOGTO_ENDPOINT` 默认 `https://auth.0thing.com`，`VITE_LOGTO_RESOURCE` 默认 `https://gateway.0thing.com/api`；后者必须与 Logto 内配置的 Backend API Resource Identifier、生产 values 中的 `backend.config.logto.audience` 完全一致。`gateway.0thing.com` 是入口主机，只有 `/api` 路径路由到 Backend。修改资源标识或 application ID 后必须重新构建并推送前端镜像；仅更新 Helm values 不会改变已打包的 Vite 环境变量。前端保持同源 `/api` 和 `/copilot` 代理。

生产部署后确认：

```bash
kubectl -n 0things-prod get certificate,ingress
kubectl -n cert-manager get secret cloudflare-api-token
kubectl wait --for=condition=Ready clusterissuer/letsencrypt-cloudflare --timeout=120s
kubectl -n 0things-prod wait --for=condition=Ready certificate/0thing-com --timeout=10m
kubectl -n 0things-prod get secret aiot-platform-tls
kubectl -n 0things-prod get pods,svc
```

再由你确认四个域名的现有 A 记录及 HTTPS 路由。若证书未就绪，查看 `kubectl -n 0things-prod describe certificate 0thing-com` 和同 namespace 的 `order,challenge`。回滚应用使用上文 `helm rollback`；独立安装的 cert-manager 不会随应用 release 回滚或删除。

没有域名或 Ingress Controller 时，可以将对应 Service 的 `type` 设置为 `LoadBalancer`。当前生产 values 通过 Traefik 暴露 HTTP 服务，EMQX 仍保留独立 LoadBalancer（MQTT 不经过这些 HTTP Ingress）。

## 持久化

PostgreSQL、Redis、NATS、EMQX 和 TDengine 使用 StatefulSet 的 `volumeClaimTemplates`。容量在环境 values 中分别设置，例如：

```yaml
tdengine:
  persistence:
    storageClass: fast-ssd
    dataSize: 100Gi
    logSize: 20Gi
```

> [!WARNING]
> 修改 values 中的数据库密码只会更新 Kubernetes Secret，不会修改已经初始化的数据库内部账号。密码轮换必须先在数据库中完成，再更新 values。

## 验证

```bash
helm lint deploy/helm/0things
helm lint deploy/helm/0things -f deploy/helm/0things/values-test.yaml
helm lint deploy/helm/0things -f deploy/helm/0things/values-prod.yaml

helm template 0things-test deploy/helm/0things \
  -n 0things-test \
  -f deploy/helm/0things/values-test.yaml

helm template 0things-prod deploy/helm/0things \
  -n 0things-prod \
  -f deploy/helm/0things/values-prod.yaml
```

安装后检查：

```bash
kubectl -n 0things-prod get pods,svc,ingress,pvc
helm status 0things-prod -n 0things-prod
```

## 升级与回滚

升级继续使用对应环境文件：

```bash
helm upgrade 0things-prod deploy/helm/0things \
  -n 0things-prod \
  -f deploy/helm/0things/values-prod.yaml \
  --atomic --wait --timeout 15m
```

查看历史并回滚：

```bash
helm history 0things-prod -n 0things-prod
helm rollback 0things-prod <REVISION> -n 0things-prod --wait
```

回滚 Chart 不会自动回滚数据库结构或 PVC 数据。
