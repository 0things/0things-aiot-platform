# EMQX MQTT 数据库认证

EMQX 通过 `aiot` PostgreSQL 数据库中的 `device_credentials` 表校验 MQTT 设备账号。

## 数据来源

创建设备时，后端会生成 MQTT 凭据：

- `username`：设备 UUID
- `password`：`SHA256(设备明文密码 + salt)` 的十六进制哈希值
- `salt`：随机生成的 16 字节盐值（十六进制）
- `credential_type`：固定为 `mqtt`
- `enabled`：凭据是否有效

明文密码只保存在 `password_ciphertext` 中，供设备管理接口解密后展示；EMQX 不使用该字段。

## EMQX Dashboard 配置

在 EMQX Dashboard 进入 **访问控制 → 认证 → 创建认证器**，选择 **PostgreSQL**，填写：

| 字段 | Helm 配置来源 |
| --- | --- |
| 服务器地址 | Helm 渲染的 PostgreSQL Service（`<release>-postgres:5432`） |
| 数据库 | `postgres.auth.database` |
| 用户名 | `postgres.auth.username` |
| 密码 | `postgres.auth.password` 或对应的已有 Secret |
| 密码加密方式 | `sha256` |
| 加盐方式 | `suffix` |
| TLS | 关闭 |

> EMQX 与 PostgreSQL 部署在同一 Kubernetes namespace 时，应使用 Helm 渲染的 Service DNS。运行时配置由 Chart 自动生成，无需在 Pod 内使用 `127.0.0.1`。

SQL 查询填写：

```sql
SELECT password, salt
FROM device_credentials
WHERE username = ${username}
  AND credential_type = 'mqtt'
  AND enabled = true
LIMIT 1
```

EMQX 会将客户端提交的密码与查询结果按 `sha256 + suffix` 规则校验，即：

```text
SHA256(客户端密码 + salt) == password
```

## 验证

创建设备后，在设备详情中获取 MQTT 用户名和密码，再执行：

```bash
./scripts/mqtt-publish.sh 127.0.0.1 1883 <username> <password> <productKey> <deviceKey>
```

发布成功说明 EMQX 数据库认证和 MQTT 上报均已生效。
