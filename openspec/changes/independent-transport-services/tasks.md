## 1. MQTT Transport 传输服务实现

- [ ] 1.1 精简 `mqtt-transport/`，移除多余的 web 路由和数据库迁移，保留轻量服务结构并通过 `go build ./cmd/server` 验证构建
- [ ] 1.2 实现 `mqtt-transport` 的 MQTT 订阅与 Kafka 上行生产（`device.message.v1`）逻辑
- [ ] 1.3 实现 `mqtt-transport` 的 Kafka 下行消费（`device.command.v1`）与 MQTT 命令下发

## 2. HTTP Transport 传输服务实现

- [ ] 2.1 精简 `http-transport/` 并实现独立的设备 HTTP 遥测/属性接入路由（`/api/v1/:deviceKey/telemetry`）
- [ ] 2.2 实现 HTTP 设备报文向 Kafka `device.message.v1` 的异步投递并通过单元测试或构建验证

## 3. CoAP Transport 传输服务实现

- [ ] 3.1 精简 `coap-transport/` 并配置 UDP 监听与 `/v1/device-ingress/{deviceKey}` 路由
- [ ] 3.2 实现 CoAP 上行转 Kafka `device.message.v1` 生产逻辑并验证构建

## 4. Rule Engine 规则计算服务实现

- [ ] 4.1 改造 `rule-engine/`，构建专门的 Kafka 消费者主循环，订阅 `device.message.v1`
- [ ] 4.2 实现物模型解析、设备影子状态更新与告警规则处理逻辑，并通过 `go build ./cmd/server` 验证构建

## 5. 编译构建与综合验证

- [ ] 5.1 验证 4 个独立项目（`mqtt-transport`、`http-transport`、`coap-transport`、`rule-engine`）全部独立通过 `go test` 与 `go build`
- [ ] 5.2 确保 `backend/` 业务核心服务正常编译与测试通过
