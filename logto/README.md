# Logto

Logto 使用独立的 PostgreSQL 数据库，并由平台 Helm Chart 与其他服务一起部署。镜像通过 `logto/Dockerfile` 构建；工作负载、Service、持久化存储和数据库连接均定义在 `deploy/helm/0things`。

在所选环境 values 中配置 Logto：

- `logto.image`：Logto 镜像与版本
- `logtoPostgres.auth`：独立数据库账号密码
- `logto.ingress`：公开端点与管理端点
- `logto.secret`：直接创建 Secret 或引用已有 Secret

测试环境示例：

```bash
helm upgrade --install 0things-test deploy/helm/0things \
  --namespace 0things-test \
  --create-namespace \
  -f deploy/helm/0things/values-test.yaml
```

Provider 凭据在 Logto Admin 中配置，不应提交到仓库。完整部署说明见 [`deploy/helm/0things/README.md`](../deploy/helm/0things/README.md)。
