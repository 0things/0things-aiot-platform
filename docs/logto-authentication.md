# Logto 认证与组织配置

平台前端使用 Logto Hosted Sign-in，Backend 只接受 Logto Access Token。用户标识来自 Token 的 `sub`，业务组织来自 `organization_id`，业务表中的 `organization_id` 使用 Logto Organization ID 字符串。

## 环境配置

1. 在 `deploy/helm/0things/values-test.yaml` 或 `values-prod.yaml` 中设置 Logto 镜像、独立 PostgreSQL 凭据、公开地址与管理地址。
2. 按 Chart README 的命令安装对应环境，Helm 会创建 Logto、专用 PostgreSQL、Service、持久化存储和运行时 Secret。
3. 在配置的 Logto Admin 地址创建 SPA 应用。
4. 将前端回调地址配置为 Redirect URI，将前端根地址配置为 Post Logout Redirect URI。
5. 将应用 ID 和端点写入所选环境的前端配置：

```dotenv
VITE_LOGTO_ENDPOINT=<logto-public-url>
VITE_LOGTO_APP_ID=<logto-spa-app-id>
VITE_LOGTO_RESOURCE=<backend-public-url>
```

Backend 的 Logto 配置必须提供 `issuer`、`audience`、`jwks_url` 和 `organization_claim`；缺失时服务启动失败并报告缺失字段。

## Provider 与组织

Google、GitHub 的 Client ID 和 Client Secret 只在 Logto Admin 或部署环境配置，不提交到仓库。组织创建、成员管理和组织切换都在 Logto 完成；业务平台不维护本地组织副本。

普通 User Token 不包含组织上下文时，组织隔离 API 返回 403，不会回退到默认组织。切换组织后必须使用新的 Organization Token 请求业务 API。

## 回滚边界

本次迁移删除本地用户、组织、成员关系和本地 JWT。清空旧数据并使用新模型初始化后，不支持在同一数据集上回滚到本地认证；如需回滚，必须同时恢复迁移前的数据库备份和应用版本。
