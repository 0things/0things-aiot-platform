# 0things Helm Chart

该 Chart 用于在 Kubernetes 中部署完整的 0things AIoT 平台，并通过独立 values 文件管理 test 和 prod 环境。

## 部署内容

Chart 默认部署以下组件：

- 应用：backend、MCP server、data-engine、transport-mqtt、AI copilot、frontend、Logto。
- 基础设施：业务 PostgreSQL、Logto PostgreSQL、Redis、NATS JetStream、EMQX、TDengine。
- 配置：服务运行配置、数据库凭据、应用凭据、镜像仓库凭据和前端 Nginx 配置。
- 网络：内部 Service，以及可选的 frontend、Logto 和 Logto Admin Ingress。

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

> [!IMPORTANT]
> `values-prod.yaml` 中的 `change-me-*` 是占位值，生产部署前必须替换。真实生产 values 可以保存在受控的私有仓库或服务器文件中。

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

先检查并替换生产占位凭据：

```bash
rg 'change-me-' deploy/helm/0things/values-prod.yaml
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

`ingress.enabled=true` 时，frontend 使用 `ingress.frontend`，Logto 使用 `ingress.logto`，管理控制台使用 `ingress.logtoAdmin`。每个入口都可以独立设置 host、path 和 TLS Secret。

没有域名或 Ingress Controller 时，可以将对应 Service 的 `type` 设置为 `LoadBalancer`。生产 values 默认通过 LoadBalancer 暴露 frontend、backend、Logto 和 EMQX。

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
