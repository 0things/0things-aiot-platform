## Why

当前平台各微服务之间已彻底解耦去除了强绑定的 Kafka 依赖，但亟需一套通用、高性能、轻量可插拔的跨微服务事件总线。
采用 NATS JetStream 作为默认消息流底座，并基于 Watermill 封装跨进程的 `Producer` 和强类型 `Consumer` 抽象，既能满足 OTA 升级指令下发与进度上报等异步闭环流转，又能为未来遥测、属性变更、状态监控等全平台事件流转提供标准基础设施。

## What Changes

- **统一事件基础设施（`0things/pkg/event`）**：
  - 新增集中式事件主题枚举（`enum.go`），统一定义精准业务语义的主题（如 `TopicOTAUpgradeCommand` 为 `ota.upgrade.command.v1`、`TopicOTAProgressReport` 为 `ota.progress.report.v1` 等）。
  - 新增通用 `Producer` 接口（`producer.go`），提供结构体自动 JSON 序列化与 Metadata/Context 追踪透传。
  - 新增通用 `Consumer` 接口与泛型辅助函数 `Subscribe[T any]`（`consumer.go`），提供自动反序列化、Ack/Nack 确认与 Panic 恢复。
  - 新增总线工厂与连接管理（`bus.go`），默认接入 NATS JetStream 驱动（`watermill-nats/v2`）。
- **OTA 业务解耦与流转对接**：
  - `backend` 在创建/执行 OTA 批量升级计划时，通过 `event.Producer` 向 `TopicOTAUpgradeCommand` 发布指令。
  - `data-engine` 订阅 `TopicOTAUpgradeCommand`，负责根据设备协议类型调度下发至 MQTT 接入层。
  - `mqtt-transport` 接收设备端 OTA 进度上报报文后，通过 `event.Producer` 向 `TopicOTAProgressReport` 发布进度事件。
  - `data-engine` 订阅 `TopicOTAProgressReport`，原子更新批次升级状态、统计指标与进度。
- **配置与依赖更新**：
  - 在各微服务（`backend`、`mqtt-transport`、`data-engine`）配置文件中新增 `event`（默认 NATS URL `nats://127.0.0.1:4222`）配置项。

## Capabilities

### New Capabilities
- `event-bus/nats-watermill`: 统一事件总线抽象与 NATS JetStream 驱动接入，提供强类型发布与订阅机制。
- `ota/event-driven-upgrade`: 基于事件驱动模型的 OTA 升级指令调度与设备进度上报闭环。

### Modified Capabilities
<!-- None: Spec-level behaviors are additive or replacing transport pipelines without changing external HTTP API contract -->

## Impact

- **代码影响**：
  - 新增共享包 `0things/pkg/event`；
  - `backend/internal/service/ota.go` 对接 `event.Producer`；
  - `mqtt-transport/internal/mqtt/client.go` 对接 `event.Producer`；
  - `data-engine/cmd/server/main.go` 与 `data-engine/internal/service/` 对接 `event.Consumer`。
- **依赖变更**：
  - 引入 `github.com/ThreeDotsLabs/watermill` 与 `github.com/ThreeDotsLabs/watermill-nats/v2`。
- **部署环境**：
  - 运行时依赖单二进制轻量级 `nats-server -js`（内存常驻约 30MB）。
