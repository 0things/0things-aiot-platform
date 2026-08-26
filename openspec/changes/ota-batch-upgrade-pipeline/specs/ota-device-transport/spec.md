## Purpose

为 OTA 服务与设备之间建立基于 Kafka 和 MQTT 的可靠传输协议，支持在线推送、离线补发、设备进度回传以及最终版本确认。

## ADDED Requirements

### Requirement: OTA 命令必须通过版本化消息传输

系统 SHALL 发布版本化的 OTA 命令消息，消息 MUST 包含 `batch_id`、设备标识、目标版本、下载地址或文件标识、摘要/签名和有效期。消息幂等键为 `batch_id + device_id + target_version`。

#### Scenario: 发布升级命令

- **WHEN** 设备任务进入可发送状态
- **THEN** 系统向 OTA Kafka 命令主题发布一条可幂等消费的升级消息

#### Scenario: 消息过期

- **WHEN** 适配器消费到已超过有效期的升级命令
- **THEN** 适配器不得向设备发布该命令，并将任务标记为过期或失败

### Requirement: MQTT 适配器必须支持设备在线和离线场景

MQTT 适配器 SHALL 将有效 OTA 命令发布到 `/ota/device/upgrade/{productKey}/{deviceName}`。在线设备应立即接收，离线设备上线后系统 MUST 能补发仍有效且未完成的任务。重复命令 MUST 以 `batch_id + device_id + target_version` 幂等处理。

#### Scenario: 在线设备推送

- **WHEN** 批次任务发送时设备在线
- **THEN** 适配器向对应 OTA upgrade topic 发布升级包信息并记录发布结果

#### Scenario: 离线设备上线

- **WHEN** 设备创建批次时离线，之后重新上线
- **THEN** 系统向该设备补发未过期且未完成的升级任务

### Requirement: 设备进度必须可关联到具体任务

系统 SHALL 接收 `/ota/device/progress/{productKey}/{deviceName}` 上报，并要求消息携带 `batch_id` 关联信息、进度、阶段和错误描述。进度上报不得单独决定升级成功。

#### Scenario: 上报下载进度

- **WHEN** 设备上报下载阶段和百分比
- **THEN** 系统更新对应任务的进度和最近上报时间，并保持任务为执行中

#### Scenario: 上报失败

- **WHEN** 设备上报下载、校验或烧写失败
- **THEN** 系统记录错误码和描述，并根据批次重试策略安排重试或标记失败

### Requirement: 最终版本上报是成功判定依据

系统 SHALL 接收 `/ota/device/inform/{productKey}/{deviceName}` 的版本上报，并将其与任务目标版本和模块进行比较。只有版本匹配时，设备任务才能标记为成功；进度为 100% 不能替代最终版本确认。

#### Scenario: 版本匹配

- **WHEN** 设备重启后上报的模块版本等于任务目标版本
- **THEN** 对应任务标记为成功，并更新批次成功统计

#### Scenario: 版本不匹配

- **WHEN** 设备上报的版本与目标版本不一致
- **THEN** 任务不得标记为成功，并按失败或重试策略继续处理

### Requirement: 下载信息必须支持安全校验

系统 SHALL 在升级信息中提供固件大小、目标版本和摘要或签名。默认下载协议应支持 HTTPS；使用 MQTT 下载时 MUST 限制为设备和固件协议支持的单文件场景。设备必须在安装前完成摘要或签名校验。

#### Scenario: 固件校验失败

- **WHEN** 设备下载完成后摘要或签名校验不通过
- **THEN** 设备上报校验失败，系统不得将任务标记为成功
