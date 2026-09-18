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

| 字段 | 本地 Docker Compose 值 |
| --- | --- |
| 服务器地址 | `postgres:5432` |
| 数据库 | `aiot` |
| 用户名 | `postgres` |
| 密码 | `password` |
| 密码加密方式 | `sha256` |
| 加盐方式 | `suffix` |
| TLS | 关闭 |

> EMQX 运行在 Docker 容器内，数据库地址必须使用 Compose 服务名 `postgres:5432`。`127.0.0.1:5432` 指向 EMQX 容器自身，无法连接项目数据库。

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
