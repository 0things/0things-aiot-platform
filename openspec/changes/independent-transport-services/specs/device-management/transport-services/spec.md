## Purpose

提供协议独立的传输微服务体系（MQTT Transport、HTTP Transport、CoAP Transport）以及独立的规则计算微服务（Rule Engine），通过统一的 Kafka 上下行总线完成设备数据的高吞吐接入与命令下发。

## ADDED Requirements

### Requirement: MQTT Transport Service
`mqtt-transport` 服务必须独立监听或连接 MQTT Broker，订阅设备上行 Topic，并将设备数据打包发送至 Kafka `device.message.v1` 主题；同时订阅 Kafka `device.command.v1` 主题，将云端指令转化为 MQTT Publish 推送给设备。

#### Scenario: MQTT telemetry uplink
- **WHEN** MQTT 设备向属性上报主题发送 JSON 遥测数据
- **THEN** `mqtt-transport` 接收并将其封装为标准 DeviceMessage，生产至 Kafka `device.message.v1`

#### Scenario: MQTT command downlink
- **WHEN** Kafka `device.command.v1` 产生针对该设备的控制指令
- **THEN** `mqtt-transport` 消费并向设备的命令 Topic 发送 MQTT 报文

### Requirement: HTTP Transport Service
`http-transport` 服务必须提供独立的 HTTP 接入端口，支持设备通过 POST 请求上报遥测与属性数据，并异步写入 Kafka `device.message.v1`。

#### Scenario: HTTP telemetry uplink
- **WHEN** 设备通过 HTTP POST 发送遥测数据至 `/api/v1/:deviceKey/telemetry`
- **THEN** `http-transport` 验证设备并投递至 Kafka `device.message.v1`，返回 HTTP 202 Accepted

### Requirement: CoAP Transport Service
`coap-transport` 服务必须监听 UDP 端口，接收 CoAP POST 报文并投递至 Kafka `device.message.v1`。

#### Scenario: CoAP telemetry uplink
- **WHEN** 设备向 `/v1/device-ingress/{deviceKey}` 发送 CoAP POST 请求
- **THEN** `coap-transport` 接收并投递至 Kafka `device.message.v1`，返回 CoAP Changed 响应

### Requirement: Rule Engine Consumer Service
`rule-engine` 服务必须作为独立的 Kafka 消费者，从 `device.message.v1` 持续拉取设备数据，完成物模型解析、时序数据写入、设备影子更新与告警规则触发。

#### Scenario: Rule engine processes device message
- **WHEN** Kafka `device.message.v1` 接收到上行遥测数据
- **THEN** `rule-engine` 消费该消息并解析物模型，更新设备状态与历史时序记录
