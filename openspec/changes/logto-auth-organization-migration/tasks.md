## 1. Logto deployment

- [x] 1.1 将 Logto Dockerfile、专用 PostgreSQL、持久化卷和共享 Docker 网络纳入根目录 Compose，并用 `docker-compose config` 验证配置
- [ ] 1.2 配置开发环境 Logto 应用、回调地址 `http://localhost:5173/callback` 和登出地址，并验证 Logto 与 Admin 端口可访问
- [ ] 1.3 在 Logto Admin 中配置 Google、GitHub Provider；在凭据未配置前确认应用仍可启动且不包含任何仓库密钥

## 2. Backend authentication

- [x] 2.1 增加 Logto issuer、audience、JWKS 和组织声明配置，并验证配置缺失时启动错误清晰
- [ ] 2.2 替换认证 middleware，使其验证 Logto Access Token 的签名、issuer、audience、过期时间和组织声明；用受保护 API 手工验证 401/403 行为
- [x] 2.3 将请求上下文中的用户标识改为 Token 的 `sub`，将组织标识改为 Token 的 Organization claim，并搜索确认没有默认组织回退
- [x] 2.4 更新 Wire、Router、Handler 和服务依赖，移除本地 JWT 签发、解析和登录依赖；用 `go test ./...` 验证 Backend 编译和现有测试

## 3. Remove local identity domain

- [x] 3.1 删除本地 `/login`、`/register`、`/organizations` 和 `/auth/switch-org` 接口及对应 service、repository、model 依赖，并用路由检查确认接口不再注册
- [x] 3.2 删除本地 `users`、`organizations`、`organization_users` 的 migration 和初始化逻辑，重新初始化数据库后确认这些表不会创建
- [x] 3.3 清理本地密码、bcrypt 和用户组织成员查询代码，并用全仓库搜索确认业务认证不再读取本地密码或成员表

## 4. String organization tenancy

- [x] 4.1 盘点所有使用 `organization_id` 的模型、DTO、Repository、查询和索引，形成迁移清单并确认没有遗漏
- [x] 4.2 将业务租户字段和相关查询统一改为 Logto Organization ID 字符串，验证产品、设备、规则链等 API 使用 Token 组织进行读写
- [ ] 4.3 删除跨组织访问和默认组织回退路径，使用两个 Logto Organization Token 手工验证组织 A 无法访问组织 B 数据

## 5. Frontend Logto flow

- [ ] 5.1 接入 Logto Hosted Sign-in，移除本地密码登录/注册表单，并验证未认证页面会跳转 Logto
- [ ] 5.2 增加 `/callback` 回调处理、Token 保存、登出和过期清理，验证成功回调、失败回调和过期 Token 行为
- [ ] 5.3 更新 API 请求拦截器和路由守卫以使用 Logto Access Token，并同步中英文 UI 文案
- [x] 5.4 通过 `pnpm -C frontend format && pnpm -C frontend build` 验证前端类型、格式和生产构建

## 6. Clean initialization and acceptance

- [ ] 6.1 在全新数据库上执行 migration，确认只创建业务数据表和 Logto 外部身份配置所需数据，不创建本地身份组织表
- [ ] 6.2 完成 Google/GitHub 登录、Organization Token 获取、Backend 访问产品/设备/规则链和跨组织拒绝的端到端验证
- [x] 6.3 更新部署和认证文档，说明 Provider 凭据只在 Logto Admin/部署环境配置，并记录旧本地认证不支持回滚
