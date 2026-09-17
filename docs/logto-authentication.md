# Logto 认证与组织配置

平台前端使用 Logto Hosted Sign-in，Backend 只接受 Logto Access Token。用户标识来自 Token 的 `sub`，业务组织来自 `organization_id`，业务表中的 `organization_id` 使用 Logto Organization ID 字符串。

## 本地启动

1. 复制根目录 `.env.example` 为 `.env`，设置 `LOGTO_POSTGRES_PASSWORD` 等部署参数。
2. 在根目录执行以下命令；Compose 会读取前端的 `.env.local` 并将 Logto 配置传入镜像构建，Logto 镜像由 `logto/Dockerfile` 构建，专用 PostgreSQL 和 Logto 与平台服务一起启动：

```bash
docker compose --env-file .env --env-file frontend/.env.local up -d --build
```
3. 在 <http://localhost:3002> 创建 SPA 应用。
4. 将 `http://localhost:5173/callback` 配置为 Redirect URI，将 `http://localhost:5173/` 配置为 Post Logout Redirect URI。
5. 将应用 ID 写入前端未提交的 `.env.local`：

```dotenv
VITE_LOGTO_ENDPOINT=http://localhost:3001
VITE_LOGTO_APP_ID=<logto-spa-app-id>
VITE_LOGTO_RESOURCE=http://localhost:8000
```

Backend 的 Logto 配置必须提供 `issuer`、`audience`、`jwks_url` 和 `organization_claim`；缺失时服务启动失败并报告缺失字段。

## Provider 与组织

Google、GitHub 的 Client ID 和 Client Secret 只在 Logto Admin 或部署环境配置，不提交到仓库。组织创建、成员管理和组织切换都在 Logto 完成；业务平台不维护本地组织副本。

普通 User Token 不包含组织上下文时，组织隔离 API 返回 403，不会回退到默认组织。切换组织后必须使用新的 Organization Token 请求业务 API。

## 回滚边界

本次迁移删除本地用户、组织、成员关系和本地 JWT。清空旧数据并使用新模型初始化后，不支持在同一数据集上回滚到本地认证；如需回滚，必须同时恢复迁移前的数据库备份和应用版本。
