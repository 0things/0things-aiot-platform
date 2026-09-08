## 1. 统一事件总线基础包 (`0things/pkg/event`)

- [x] 1.1 创建 `pkg/event/enum.go`，精准定义系统事件 Topic 常量（`TopicOTAUpgradeCommand`、`TopicOTAProgressReport`、`TopicDeviceTelemetryReport` 等），并通过单元测试验证常量定义。
- [x] 1.2 创建 `pkg/event/bus.go`、`pkg/event/producer.go`、`pkg/event/consumer.go`，基于 Watermill NATS JetStream 实现通用 `Producer` 和泛型 `Subscribe[T any]` 辅助函数，并编写完整单测 `pkg/event/event_test.go` 验证消息收发、泛型反序列化与 Ack/Nack。

## 2. OTA 业务流转改造与微服务对接

- [x] 2.1 在 `backend` 服务中引入 `pkg/event`，在 `wire` 中配置 `Producer` 依赖注入，将 `OTAService` 的升级指令下发对接至 `TopicOTAUpgradeCommand`，并编写/更新单元测试。
- [x] 2.2 在 `mqtt-transport` 中引入 `pkg/event`，在监听到设备端 OTA 进度上报报文时向 `TopicOTAProgressReport` 发送 `OTAUpgradeReport` 事件，并验证单元测试。
- [x] 2.3 在 `data-engine` 中引入 `pkg/event`，启动针对 `TopicOTAUpgradeCommand`（调度 MQTT 下发）与 `TopicOTAProgressReport`（原子更新数据库状态）的泛型订阅消费协程，并验证单元测试。

## 3. 配置文件与全局集成验证

- [x] 3.1 更新各微服务配置文件（`backend`、`transport-mqtt`、`data-engine` 的 `config.example.yml` 和 `local.yml`），增加 `event` NATS JetStream 配置。
- [x] 3.2 运行根目录 `go mod tidy`、`go work sync`、`make test && make build`，确保所有服务编译与单元测试 100% 通过。
