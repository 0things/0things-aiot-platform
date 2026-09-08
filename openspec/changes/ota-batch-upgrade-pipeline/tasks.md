## 1. 协议与受控基础设施

- [x] 1.1 固化 `ota.upgrade.command.v1` 与 `ota.upgrade.report.v1` JSON schema 和示例，覆盖 `batch_id`、设备标识、默认模块、目标版本、HTTPS 下载信息、SHA-256、进度和错误；用 schema 校验示例。
- [ ] 1.2 部署 EMQX，配置 TLS、设备认证和按产品/设备隔离的 ACL；使用测试客户端验证未授权 topic 不能发布或订阅。
- [ ] 1.3 为 data-engine 增加 Kafka、EMQX、数据库和 TLS 配置校验；缺失必需凭据时安全拒绝启动。
- [ ] 1.4 引入并锁定 Eclipse Paho Go MQTT 客户端，验证 QoS 1 发布、订阅、重连与优雅关闭。

## 2. 批次与设备任务

- [ ] 2.1 扩展 OTA 设备任务，增加 `batch_id`、默认模块、目标版本、发送结果、进度、当前版本、最后上报时间和错误，并建立 `(batch_id, device_id)` 唯一约束及查询索引；验证旧数据可读。
- [ ] 2.2 改造 `BatchUpgrade`，在一个事务中创建批次和每台设备的独立任务；重复设备输入不得生成重复任务，同一设备可保留跨批次历史。
- [ ] 2.3 实现 `created`、`dispatching`、`in_progress`、`succeeded`、`partial_success`、`failed` 的任务与批次状态聚合；覆盖成功、发送失败、设备失败和混合结果测试。

## 3. 命令下发与设备回报

- [ ] 3.1 在批次事务提交后发布 Kafka OTA 命令，按发布结果更新任务发送状态和错误；使用 mock 或测试容器覆盖成功与失败。
- [ ] 3.2 在 data-engine 独立消费 `ota.upgrade.command.v1`，以 QoS 1 向 `/ota/device/upgrade/{productKey}/{deviceName}` 下发完整升级元数据，并直接写入下发结果；用 EMQX 测试客户端验证。
- [ ] 3.3 在 data-engine 独立消费 `ota.upgrade.report.v1`，要求 `batch_id` 并直接写入共享数据库的设备进度、错误和最后上报时间；进度不得标记任务成功。
- [ ] 3.4 在 data-engine 的 `ota.upgrade.report.v1` consumer 中处理 inform；只有默认模块版本等于目标版本时标记成功，验证 100% progress 但版本不匹配仍不成功。
- [ ] 3.5 移除 mqtt-transport 的 Kafka 消费和 MQTT 下行能力；验证数据库状态只由 data-engine 的派发结果或设备上报改变。

## 4. 批次可见性与验收

- [ ] 4.1 扩展批次详情、设备任务列表与统计 API，按 `batch_id` 返回状态、进度、错误、当前版本和目标版本；生成 Swagger 并覆盖 handler 测试。
- [ ] 4.2 更新 OTA 批次详情和设备部署列表，展示批次维度的任务状态、进度、错误、当前和目标版本；补齐中英文文案并执行前端构建。
- [ ] 4.3 完成端到端验收：创建批次 → Kafka → EMQX → 测试设备收到命令 → progress → 目标版本 inform → 批次成功；保存验收记录并执行 `make test && make build` 与 `pnpm -C frontend format && pnpm -C frontend build`。
