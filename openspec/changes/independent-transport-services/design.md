## Context

参考 ThingsBoard 架构，将原来耦合在单体或通用网关中的协议接入逻辑，拆分为 4 个独立的 Go 模块工程：
1. `mqtt-transport`：专职 MQTT 协议接入与下行命令分发
2. `http-transport`：专职 HTTP 设备数据上报 REST API
3. `coap-transport`：专职 CoAP UDP 数据接收
4. `rule-engine`：专职 Kafka 数据消费、物模型处理与规则触发

## Goals / Non-Goals

**Goals:**
- 将 4 个模块分别打造成轻量独立编译、可独立部署运行的 Go 微服务。
- 统一 Kafka 上行主题 `device.message.v1` 和下行主题 `device.command.v1` 契约。
- 移除各传输服务中冗余的 CRUD、DAL、页面路由与无关迁移代码，保持网络层极简高吞吐。
- 在 `rule-engine` 中实现独立的 Kafka 消费者流式处理引擎。

**Non-Goals:**
- 不重写核心后台 `backend/` 的现有业务接口与数据模型。
- 不引入重型的第三方 Flink/Spark 大数据引擎，保持纯 Go 实现。

## Decisions

### 1. 传输服务与业务服务解耦
- **决策**: 各 Transport 服务不直连 MySQL / 业务表，只通过 Kafka `device.message.v1` 投递上行消息。
- **替代方案**: 直连数据库或调用后端 HTTP 内部接口（缺点：高并发下容易打满连接池或增加调用链延迟）。

### 2. 下行指令基于 Channel 路由
- **决策**: `mqtt-transport` 消费 `device.command.v1` 时，根据 `transport == "mqtt"` 或广播机制，由持有该设备连接的实例通过 MQTT 下发。

### 3. `rule-engine` 作为独立 Worker
- **决策**: `rule-engine` 仅作为后台计算型 Worker 启动，无对外管理 HTTP 端口，专职消费 Kafka 并执行物模型解析与规则动作。

## Risks / Trade-offs

- **[Kafka 单点依赖]** → 各 Transport 和 Rule Engine 均支持断网重连与重试缓冲机制。
- **[本地多进程启动复杂度]** → 提供统一的 shell 脚本或 Makefile 支持一键启动/构建。

