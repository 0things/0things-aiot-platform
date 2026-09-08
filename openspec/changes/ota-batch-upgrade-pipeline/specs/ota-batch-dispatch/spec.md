## Purpose

为静态 OTA 批次提供按设备可追踪的任务和状态聚合，确保控制台状态来自实际下发结果及设备回报。

## ADDED Requirements

### Requirement: 静态批次必须生成独立设备任务

系统 SHALL 在创建静态 OTA 批次时，为每台目标设备创建一条由 `batch_id + device_id` 唯一定位的任务。任务 MUST 保存默认模块、目标版本、发送结果、进度、当前版本、最后上报时间和错误；同一设备进入不同批次时 MUST 保留不同历史任务。

#### Scenario: 创建批次

- **WHEN** 用户提交升级包和目标设备列表
- **THEN** 系统在同一事务中创建批次及每台设备的独立任务，并返回批次标识

#### Scenario: 重复设备输入

- **WHEN** 一个创建请求中重复包含同一设备
- **THEN** 系统至多为该设备创建一条本批次任务

### Requirement: 批次状态必须由设备任务聚合

设备任务 SHALL 表示已创建、已发送、执行中、成功或失败；批次 SHALL 表示 `created`、`dispatching`、`in_progress`、`succeeded`、`partial_success` 或 `failed`。系统 MUST 仅根据任务的实际发送结果和设备回报聚合批次状态。

#### Scenario: 批次部分成功

- **WHEN** 一个批次至少一台设备成功，且至少一台设备失败
- **THEN** 批次状态为 `partial_success`

### Requirement: 查询必须以批次为边界

系统 SHALL 提供按 `batch_id` 查询批次详情、设备任务和统计的能力。响应 MUST 返回任务状态、进度、错误、当前版本和目标版本，且不得混入其他批次任务。

#### Scenario: 查询批次设备

- **WHEN** 用户查询指定批次的设备任务
- **THEN** 系统只返回该批次的设备任务和聚合统计
