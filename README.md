<div align="center">

<img src="./frontend/public/images/favicon.svg" alt="0things logo" width="72" />

# 0things AIoT 平台

面向物联网产品、设备、物模型、遥测、OTA 与自动化场景的一体化开发平台。

[![Frontend](https://img.shields.io/badge/frontend-React%2019%20%2B%20Vite-646CFF?logo=react&logoColor=white)](./frontend)
[![Backend](https://img.shields.io/badge/backend-Go%20%2B%20Gin-00ADD8?logo=go&logoColor=white)](./backend)
[![License](https://img.shields.io/badge/license-MIT-green)](./LICENSE)

[功能概览](#功能概览) · [快速开始](#快速开始) · [本地开发](#本地开发) · [项目结构](#项目结构)

</div>

## 功能概览

- 产品、设备、设备分组与物模型（TSL）管理
- HTTP / MQTT 设备接入与设备凭据管理
- 设备遥测、属性、事件、RPC 与历史数据查询
- OTA 固件、批量升级与升级状态跟踪
- 规则链、自动化、告警、任务和运行监控
- Logto 组织与用户认证
- MCP 服务与 AI Copilot，用自然语言访问 IoT 能力
- 中英文界面、深色模式和响应式管理控制台

## 架构

```text
                    ┌──────────────────────────────┐
                    │ frontend  React / Vite        │ :5173 / :80
                    └──────────────┬───────────────┘
                                   │ REST / MCP
             ┌─────────────────────┴─────────────────────┐
             │                                           │
┌────────────▼────────────┐               ┌──────────────▼─────────────┐
│ backend  Gin REST API    │ :8000         │ ai-copilot  AI Gateway      │ :8005
│ 管理、认证、设备与业务    │               │ MCP 客户端与对话运行时       │
└───────┬─────────┬────────┘               └──────────────┬─────────────┘
        │         │                                       │
        │         └──────────────┐              ┌────────▼────────┐
        │                        │              │ mcp-server :8009 │
        │                        │              └─────────────────┘
        │                        │
┌───────▼────────┐     ┌─────────▼─────────┐
│ data-engine     │     │ transport-http   │ :8081
│ 计算与任务中心   │     │ HTTP 设备网关     │
└───────┬────────┘     └──────────────────┘
        │
┌───────▼────────┐     ┌──────────────────┐
│ transport-mqtt │     │ EMQX :1883       │
│ MQTT 设备网关   │────▶│ MQTT Broker      │
└────────────────┘     └──────────────────┘

       NATS :4222 · Redis :6379 · Logto :3001/:3002 · SQLite / MySQL / PostgreSQL
```

服务通过 NATS 传递事件；业务数据库由 backend 管理，时序数据通过 `pkg/tsdb` 的可插拔驱动写入。开发环境可以使用 SQLite，完整 Docker 环境也可以接入 PostgreSQL、Redis 和 EMQX。

## 快速开始

### 前置条件

- Docker Desktop 与 Docker Compose
- Node.js 25+ 与 pnpm
- Go 1.25+
- Git

### 使用 Docker Compose 启动完整环境

首次启动前复制本地配置模板，并填写必要的密码、应用 ID 和模型 API Key：

```bash
cp logto/.env.example logto/.env
cp ai-copilot/.env.example ai-copilot/.env
cp backend/config/config.example.yml backend/config/docker.yml
cp data-engine/config/config.example.yml data-engine/config/docker.yml
cp transport-http/config/config.example.yml transport-http/config/docker.yml
cp transport-mqtt/config/config.example.yml transport-mqtt/config/docker.yml
```

然后启动所有容器：

```bash
docker compose up -d --build
```

常用入口：

| 服务 | 地址 |
| --- | --- |
| 管理控制台 | [http://localhost:5173](http://localhost:5173) |
| 管理 API | [http://localhost:8000](http://localhost:8000) |
| Swagger | [http://localhost:8000/swagger/index.html](http://localhost:8000/swagger/index.html) |
| AI Copilot | [http://localhost:8005](http://localhost:8005) |
| MCP Streamable HTTP | [http://localhost:8009/mcp](http://localhost:8009/mcp) |
| HTTP 设备网关 | [http://localhost:8081](http://localhost:8081) |
| Logto 管理控制台 | [http://localhost:3002](http://localhost:3002) |
| Logto OIDC 服务 | [http://localhost:3001](http://localhost:3001) |
| EMQX Dashboard | [http://localhost:18083](http://localhost:18083) |

停止环境：

```bash
docker compose down
```

> [!IMPORTANT]
> 所有 `.env`、`config/local.yml` 和 `config/docker.yml` 都可能包含密钥或数据库凭据，不要提交到 Git。生产环境请替换示例中的默认密码和加密密钥。

## 本地开发

### 启动依赖服务

只启动基础设施：

```bash
docker compose up -d nats emqx redis logto-postgres logto
```

### 配置并执行数据库迁移

backend 默认可以使用 SQLite：

```bash
cp backend/config/config.example.yml backend/config/local.yml
cd backend && go run ./cmd/migration -conf ./config/local.yml
```

如果使用其他数据库，请先修改 `backend/config/local.yml` 中的 `data.db.aiot` 配置。迁移程序使用 GORM `AutoMigrate` 创建和更新业务表。

### 启动服务

每个服务使用自己的 `config/local.yml`。可以分别打开终端启动：

```bash
cd backend && go run ./cmd/server -conf ./config/local.yml
cd data-engine && go run ./cmd/server -conf ./config/local.yml
cd transport-http && go run ./cmd/server -conf ./config/local.yml
cd transport-mqtt && go run ./cmd/server -conf ./config/local.yml
cd backend && go run ./cmd/mcp -conf ./config/local.yml
cd ai-copilot && pnpm install && pnpm dev
cd frontend && pnpm install && pnpm dev
```

也可以在已经准备好各服务配置后，从仓库根目录执行：

```bash
make start
```

停止由脚本启动的本地进程：

```bash
make stop
```

前端默认访问 `http://localhost:5173`，后端 API 默认访问 `http://localhost:8000`。前端环境变量模板位于 [`frontend/.env.example`](./frontend/.env.example)，AI Copilot 模板位于 [`ai-copilot/.env.example`](./ai-copilot/.env.example)。

## 设备模拟器

仓库提供一个简单的 MQTT 虚拟设备，可用于发送遥测数据或触发高温告警：

```bash
make mock-device
make mock-alarm
```

设备网关和 Broker 的地址分别通过 `transport-mqtt/config/` 与 backend 的 `device_gateway.mqtt` 配置。默认设备端点 MQTT Host 为 `mqtt.0things.com`，本地 Broker 示例为 `tcp://127.0.0.1:1883`。

## 常用命令

| 命令 | 用途 |
| --- | --- |
| `make test` | 运行所有 Go 服务和共享包测试 |
| `make build` | 编译各 Go 服务 |
| `pnpm -C frontend format` | 格式化前端代码 |
| `pnpm -C frontend build` | 类型检查并构建前端 |
| `pnpm -C frontend lint` | 检查前端代码 |
| `pnpm -C ai-copilot test` | 运行 AI Copilot 测试 |
| `make -C backend gen` | 生成 GORM Gen 查询代码 |
| `make -C backend swag` | 生成 Swagger 文档 |
| `pnpm -C frontend generate:api` | 根据 OpenAPI 重新生成前端 API 客户端 |

> [!NOTE]
> `frontend/src/api/generated/`、backend 的 `wire_gen.go`、GORM Gen 文件和 Swagger 文档属于生成产物。修改契约或依赖后，请运行对应的生成命令，不要手工编辑这些文件。

## 项目结构

| 路径 | 说明 |
| --- | --- |
| [`frontend/`](./frontend) | React 管理控制台，包含设备、产品、规则、告警、组织和 AI Copilot 页面 |
| [`backend/`](./backend) | Gin REST API、认证、业务服务、数据库迁移和 MCP Server |
| [`data-engine/`](./data-engine) | 遥测时序数据、设备影子、任务和数据计算服务 |
| [`transport-http/`](./transport-http) | HTTP 设备接入网关 |
| [`transport-mqtt/`](./transport-mqtt) | MQTT 设备接入、事件和 OTA 消息处理 |
| [`ai-copilot/`](./ai-copilot) | 基于 OpenAI-compatible 模型和 MCP 的 AI Gateway |
| [`pkg/event/`](./pkg/event) | 事件总线抽象与 NATS/Watermill 集成 |
| [`pkg/protocol/`](./pkg/protocol) | JSON、JavaScript、Modbus、GB28181 等协议编解码器 |
| [`pkg/tsdb/`](./pkg/tsdb) | SQLite、PostgreSQL、TDengine、InfluxDB 等时序存储适配器 |
| [`logto/`](./logto) | Logto 身份服务及其 Docker 配置 |
| [`openspec/`](./openspec) | 功能变更提案、设计和任务规格 |
| [`docs/`](./docs) | 架构、认证、数据模型和运行说明 |

## 文档

- [Logto 认证与组织说明](./docs/logto-authentication.md)
- [Logto 组织租户盘点](./docs/logto-organization-tenancy-inventory.md)
- [OTA 架构对比](./docs/ota-architecture-comparison.md)
- [规则链数据模型](./docs/rule-chain-data-model.md)
- [OpenSpec 变更目录](./openspec/changes)
