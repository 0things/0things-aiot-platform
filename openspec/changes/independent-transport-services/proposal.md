## Why

参考 ThingsBoard 传输层架构，解决原单一 `device-gateway` 中不同协议混杂运行、长短连接相互影响、无法针对特定协议独立伸缩与部署的问题。将传输层按协议边界（MQTT、HTTP、CoAP）彻底拆解为独立的进程与微服务，通过 Kafka 消息总线与核心服务解耦。

## What Changes

- **移除** 集中式的 `cmd/device-gateway/`。
- **新增** 独立的 MQTT 传输服务入口：`backend/cmd/mqtt-transport/main.go`。
- **新增** 独立的 HTTP 传输服务入口：`backend/cmd/http-transport/main.go`。
- **新增** 独立的 CoAP 传输服务入口：`backend/cmd/coap-transport/main.go`。
- **重构** `backend/internal/transport/` 内部公共模块，统一上行投递契约（`device.message.v1`）与下行指令消费契约（`device.command.v1`）。
- **更新** `Makefile`，支持各个独立 Transport 服务的独立编译构建与打包（`bin/mqtt-transport`、`bin/http-transport`、`bin/coap-transport`）。

## Capabilities

### New Capabilities
- `device-management/transport-services`: 定义独立协议传输服务（MQTT/HTTP/CoAP）与核心消息总线上下行标准契约。

### Modified Capabilities
- 无

## Impact

- **Affected Code**: `backend/cmd/`、`backend/internal/transport/`、`backend/Makefile`、`backend/config/local.yml`
- **APIs**: 设备 HTTP 接入端口由各 Transport 配置独立提供，核心后台 HTTP REST API 不受影响
- **Dependencies**: 依赖 Kafka / 内部消息队列作为传输层与业务服务之间的上下行总线
