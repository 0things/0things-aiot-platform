## Why

当前 OTA 批次只创建数据库记录，定时任务会模拟推进状态；控制台无法确认设备是否实际收到升级命令、完成安装或重启到目标版本。先交付一条可在测试设备验证的静态批次闭环，才能让 OTA 成为可信的设备运维能力。

## What Changes

- 为每个静态批次中的设备创建可按 `batch_id + device_id` 关联的独立任务，并按任务聚合批次状态。
- 批次事务提交后由 backend 发布版本化 Kafka OTA 命令；data-engine 以 MQTT QoS 1 下发到设备，并独立消费下发结果与设备回报写入共享 OTA 数据库。
- 接收带 `batch_id` 的设备进度和版本上报；只有默认模块版本等于目标版本才判定任务成功。
- 提供按批次查询任务、进度、错误和目标版本的 API 与控制台展示，并使用一个测试设备完成端到端验收。
- **BREAKING**：新的 OTA progress 与 inform 上报必须携带 `batch_id`；旧上报不得用于更新新批次任务。

## Capabilities

### New Capabilities

- `ota-batch-dispatch`: 静态 OTA 批次、设备任务、状态聚合与批次查询。
- `ota-device-transport`: Kafka 到 MQTT 的命令下发、设备进度上报与最终版本确认。

### Modified Capabilities

- 无；仓库没有已归档的 OTA capability spec。

## Impact

- backend 的 OTA 模型、服务、API 和 Kafka producer；data-engine 的两个 OTA topic consumer、MQTT dispatcher 与 OTA 状态存储。
- EMQX 的 TLS、设备认证及按产品/设备隔离的 ACL；设备端的 OTA MQTT 与上报契约。
- OTA 批次详情与设备部署列表，以及中英文文案。
- 本期不含离线补发、DLQ、自动重试、取消、灰度/限速、多模块或生产指标告警；这些在 MVP 验收后另建变更。
