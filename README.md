# 0things AIoT 平台

一套面向物联网产品、设备、物模型（TSL）、OTA、规则与运行遥测的统一管理平台。

[![前端](https://img.shields.io/badge/前端-React%2019%20%2B%20Vite-646CFF?logo=react&logoColor=white)](./frontend)
[![后端](https://img.shields.io/badge/后端-Go%20%2B%20Gin-00ADD8?logo=go&logoColor=white)](./backend)
[![遥测服务](https://img.shields.io/badge/遥测服务-Go%20%2B%20Kratos-00ADD8?logo=go&logoColor=white)](./telemetry-service)

> [!NOTE]
> 这是一个多服务仓库。仅做界面开发时可以单独运行前端；需要访问 API 的功能还需要启动后端及其依赖服务。

## 功能概览

- 产品、设备、设备分组与 TSL 管理
- 设备开发流程、Topic 配置与消息解析
- OTA 包管理、规则引擎、告警、用户、任务与 IoT 总览
- REST API、Swagger 文档，以及独立的遥测/事件接入服务

## 架构

```text
frontend（React/Vite，:5173）
        │
        ├── backend（Gin REST API，:8000）
        │      └── PostgreSQL / Redis / 可选 Kafka 与 MongoDB
        │
        └── telemetry-service（Kratos HTTP/gRPC，:8013/:9013）
               └── PostgreSQL / Redis / Kafka
```

## 快速开始

### 前置条件

- Node.js 与 pnpm
- `backend` 需要 Go 1.24.10+，`telemetry-service` 需要 Go 1.25+
- Docker Compose（用于本地 MySQL/Redis 辅助环境）
- 完整运行 API 与遥测功能时，需要可访问的 PostgreSQL、Redis 与 Kafka

### 1. 配置本地服务

```bash
cd backend
cp config/config.example.yml config/local.yml
cd deploy/docker-compose && docker compose up -d
```

在 `backend/config/local.yml` 中填写真实的本地连接信息和密钥。根据已获准的环境配置创建 `telemetry-service/configs/config.yaml`，并提供其中的 PostgreSQL、Redis 与 Kafka 设置。这些本地配置文件均被 Git 忽略。

### 2. 启动平台

在三个独立终端中运行：

```bash
cd backend && go run ./cmd/server -conf ./config/local.yml
cd telemetry-service && go run ./cmd/telemetry-service -conf ./configs
cd frontend && pnpm install && pnpm dev
```

打开 [http://localhost:5173](http://localhost:5173)。后端 Swagger 文档位于 [http://localhost:8000/swagger/index.html](http://localhost:8000/swagger/index.html)。

## 前端配置

复制 `frontend/.env.example` 为 `frontend/.env.local`，即可覆盖服务地址或设置 Clerk 公钥。设备 API 默认使用 `http://localhost:8000`；认证和通知服务可分别使用 `8003`、`8004` 端口。

> [!IMPORTANT]
> 不要提交凭据、DSN、JWT/API 签名密钥或本地运行配置。可共享的模板仅限 `config.example.yml` 与 `.env.example`。

## 常用开发命令

| 范围     | 命令                              | 用途                              |
| -------- | --------------------------------- | --------------------------------- |
| 前端     | `pnpm dev`                        | 启动 Vite 开发服务器              |
| 前端     | `pnpm build`                      | 类型检查并构建生产产物            |
| 前端     | `pnpm lint` / `pnpm format:check` | 运行 ESLint / 检查 Prettier 格式  |
| 后端     | `make test`                       | 运行服务端测试并生成覆盖率报告    |
| 后端     | `make build`                      | 构建 `bin/server`                 |
| 遥测服务 | `make build`                      | 构建服务二进制文件                |
| 遥测服务 | `make api`                        | 修改 protobuf 后重新生成 API 代码 |

## 目录说明

| 路径                                               | 说明                                                                                         |
| -------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| [`frontend/`](./frontend)                | React 前端；功能模块在 `src/features/`，路由在 `src/routes/`，国际化资源在 `public/locales/` |
| [`backend/`](./backend)                  | Gin 后端；应用代码位于 `cmd/` 和 `internal/`，API 文档位于 `docs/`                           |
| [`telemetry-service/`](./telemetry-service)        | Kratos 遥测/事件接入服务；protobuf API 位于 `api/` 和 `internal/`                            |
| [`start-all-services.sh`](./start-all-services.sh) | 适用于已配置环境的本地多服务启动与日志脚本                                                   |

`frontend/src/api/generated/` 下的前端 API 客户端，以及 protobuf 生成产物，都应通过其契约重新生成，不要手工修改。
