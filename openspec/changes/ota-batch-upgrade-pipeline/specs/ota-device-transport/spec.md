## Purpose

为 OTA MVP 建立 Kafka 到 MQTT 的命令传输和设备回报协议，以最终版本确认升级结果。

## ADDED Requirements

### Requirement: OTA 命令必须以版本化消息下发

系统 SHALL 在设备任务的批次事务提交后发布 `ota.upgrade.command.v1` Kafka 消息。消息 MUST 含 `batch_id`、设备标识、默认模块、目标版本、HTTPS 下载 URL、文件大小、SHA-256 摘要和有效期；data-engine MUST 通过 TLS 连接的 EMQX，以 QoS 1 向 `/ota/device/upgrade/{productKey}/{deviceName}` 发布升级元数据。mqtt-transport MUST 只生产 Kafka 上报事件，不得消费 Kafka 或执行 MQTT 下行。

#### Scenario: 下发升级命令

- **WHEN** 后端成功创建批次设备任务
- **THEN** data-engine 向对应设备 topic 下发完整升级元数据，并直接更新命令发布结果

#### Scenario: 未授权 topic

- **WHEN** 未获该产品或设备授权的 MQTT 客户端尝试发布或订阅 OTA topic
- **THEN** EMQX 拒绝该操作

### Requirement: 进度上报必须关联任务但不得判定成功

系统 SHALL 接收带 `batch_id` 的 progress 上报，并按批次和设备更新阶段、百分比、错误和最后上报时间。progress 上报 MUST NOT 单独将任务标记为成功。

#### Scenario: 上报完成进度

- **WHEN** 设备上报百分比为 100 的 progress
- **THEN** 系统保留任务为未成功，直至收到匹配版本的 inform 上报

### Requirement: OTA topic 必须由 data-engine 分别消费

data-engine SHALL 使用独立消费组分别消费 `ota.upgrade.command.v1` 和 `ota.upgrade.report.v1`，并直接写入 OTA 任务及批次状态。backend MUST NOT 消费这些 topic。

#### Scenario: MQTT 派发完成

- **WHEN** data-engine 向设备 topic 的 QoS 1 发布成功或失败
- **THEN** 命令 consumer 直接更新对应设备任务

### Requirement: 最终版本上报是成功的唯一设备依据

系统 SHALL 接收带 `batch_id` 的默认模块 inform 上报。只有该版本等于任务目标版本时，系统才可将任务标记为成功；版本不匹配时 MUST 保持非成功并记录当前版本。

#### Scenario: 版本匹配

- **WHEN** 设备上报的默认模块版本等于目标版本
- **THEN** 对应任务标记成功，并重新计算批次状态

#### Scenario: 版本不匹配

- **WHEN** 设备上报的默认模块版本不同于目标版本
- **THEN** 系统更新当前版本但不得将任务标记成功
