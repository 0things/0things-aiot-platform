## Context

0things 已经完成了对 Kafka 代码的全面剥离。为了支持跨微服务（`backend`、`mqtt-transport`、`http-transport`、`data-engine`）的解耦异步通信，需要在公共包 `0things/pkg/event` 中构建一套统一的事件总线。

## Goals / Non-Goals

**Goals:**
- 提供精准业务语义的集中事件主题枚举 `event.Topic`（如 `TopicOTAUpgradeCommand`、`TopicOTAProgressReport` 等）。
- 基于 Watermill 和 NATS JetStream 驱动构建通用的 `Producer` 和 `Consumer` 接口。
- 提供强类型、泛型辅助方法 `event.Subscribe[T any]`，自动处理 JSON 序列化/反序列化、Ack/Nack 与 Panic 保护。
- 完成协议感知（Protocol-Aware）的 OTA 升级指令下发（`TopicOTAUpgradeCommand`）与设备进度上报（`TopicOTAProgressReport`）闭环流转。
- 单机和本地单元测试环境下支持内存测试与快速启动。

**Non-Goals:**
- 本变更不重写规则引擎复杂的流式 DAG 计算（属于后续规则链演进）。
- 本变更不修改前端 API 契约与 UI 交互界面。

## Decisions

### 1. 消息中间件与驱动选型：NATS JetStream + Watermill
- **决定**：使用 NATS JetStream 作为默认消息底座，并在 Go 代码层通过 `github.com/ThreeDotsLabs/watermill-nats/v2` 提供封装。
- **理由**：NATS 是纯 Go 编写的极轻量中间件（单二进制 25MB、内存占用 30MB、支持持久化流与微秒级延迟），避免了 JVM 资源的沉重负担。Watermill 则提供了开箱即用的生命周期与中间件能力。
- **替代方案对比**：
  - *纯 GoChannel*：仅支持单进程内内存共享，无法直接满足多微服务跨进程场景；
  - *Apache Kafka*：依赖过重、JVM 内存开销大、本地开发不便；
  - *Dapr*：Sidecar 额外增加 4 次跨进程 IPC 与性能损耗。

### 2. 接口抽象层设计：精准业务语义枚举 + 泛型订阅
- **决定**：
  - 在 `pkg/event/enum.go` 中按 `领域.动作.事件.版本` 规范统一声明 `Topic` 常量（如 `TopicOTAUpgradeCommand = "ota.upgrade.command.v1"`, `TopicOTAProgressReport = "ota.progress.report.v1"`）；
  - `Producer.Publish(ctx, topic, payload)` 接收任意 Go 结构体自动序列化；
  - `event.Subscribe[T any](ctx, consumer, topic, handler)` 提供类型安全的泛型回调函数。
- **理由**：彻底杜绝各微服务内字符串硬编码与重复编写 `json.Unmarshal` / Panic 恢复模板代码。

### 3. 设备传输协议感知（Protocol-Aware OTA Routing）
- **决定**：
  - **MQTT 设备**：`backend` 触发升级时，向 `TopicOTAUpgradeCommand` 发送事件（附带 `transport: mqtt`），`mqtt-transport` 消费并主动 Push 至设备长连接；
  - **HTTP 设备**：无长连接，系统将固件元数据写入设备期望状态（Desired State），设备通过定期轮询 `GET /api/v1/ota/info` 拉取，不触发下推事件；
  - **进度上报**：所有 Transport（MQTT/HTTP/CoAP）统一解析设备上报进度并向 `TopicOTAProgressReport` 发布标准事件，由 `data-engine` 原子更新批次统计与状态。

## Risks / Trade-offs

- **[NATS 依赖启动]**：微服务启动需要 NATS JetStream 服务可用。
  - *缓解措施*：在 `pkg/event` 中提供优雅的重连机制；在单元测试中支持 Watermill `GoChannel` 零外部依赖测试。
- **[消息顺序性]**：设备上报 OTA 进度需保证基本时序。
  - *缓解措施*：NATS JetStream 消费组支持按 Subject / Consumer Group 顺序交付，上报数据中包含客户端 `timestamp` 进行时间戳保序校验。

## Migration Plan

1. 在 `pkg/event` 中实现总线模块与单元测试。
2. 在 `backend`、`mqtt-transport`、`data-engine` 中引入 `pkg/event` 并注入 `Producer`/`Consumer`。
3. 更新配置文件 `config.example.yml` 和 `local.yml`。
4. 运行 `make test && make build` 全局验证。
