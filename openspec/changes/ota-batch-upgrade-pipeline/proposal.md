## Why

当前 `BatchUpgrade` 主要完成数据库批次和设备记录创建，定时任务只推进数据库状态，尚未形成 Kafka、MQTT 与设备端之间可追踪、可重试的 OTA 闭环。同一设备在不同批次中的任务也缺少稳定的任务关联，导致进度、失败和最终版本无法可靠归属到具体批次。

参考阿里云 OTA 的静态批次流程，将升级包推送、设备下载、进度上报和最终版本确认拆分为明确阶段，建立以设备任务为核心的可靠发布链路，才能支持离线补发、重试、超时、取消以及准确的批次统计。

## What Changes

- 将批量升级建模为批次任务和设备任务，使用 `batch_id + device_id` 唯一定位批次内设备任务。
- 在创建批次事务提交后，直接发布版本化 Kafka OTA 命令，不再通过定时任务模拟下发状态。
- 增加 Kafka 到 MQTT 的 OTA 传输适配能力，向设备发布升级包信息，并支持设备离线后的上线补发和消息幂等。
- 接收设备版本上报和升级进度，按 `batch_id + device_id` 精确更新设备任务；只有设备重启后上报目标版本才判定升级成功。
- 支持设备任务的发送、执行、失败、超时、重试和取消状态，以及批次级统计和查询。
- 支持静态批次的升级速率、重试间隔、重试次数、设备升级超时和任务冲突策略，为后续灰度/动态升级保留扩展边界。
- 保留现有 OTA MQTT 主题约定，默认使用 MQTT 下发升级元数据、HTTPS 下载固件，并记录摘要/签名等校验信息。
- **BREAKING**：OTA 设备上报接口需要携带 `batch_id`；仅使用升级包和设备标识的旧上报方式不能用于多批次精确更新。

## Capabilities

### New Capabilities

- `ota-batch-dispatch`: 管理 OTA 批次、设备任务、状态机、调度策略和批次统计。
- `ota-device-transport`: 通过 Kafka、MQTT 完成 OTA 命令下发，以及设备进度和版本事件回传。

### Modified Capabilities

- （无现有 OTA capability spec；现有 `device-management/device-tags` 与本变更无关。）

## Impact

- 后端 OTA handler、service、repository、任务调度和 Kafka 集成。
- 数据库中的 OTA 批次、设备任务、发布结果、重试和索引结构。
- Kafka topic、消息 schema、消费组、死信和重放策略。
- MQTT 适配器及设备 OTA 主题、进度/版本上报协议。
- OTA 批次详情、设备部署列表和统计接口，以及前端批次进度展示。
- 设备端需要实现升级信息校验、HTTPS 下载、固件校验、安装重启、进度上报和最终版本上报。
- `telemetry-service` 不作为本变更的实现依赖；传输能力应归属 backend 内的 OTA worker 或独立 OTA transport 服务。
