## Why

设备详情中的「物模型数据 → 服务调用」目前仅为占位内容，用户无法回溯设备服务调用的输入和返回结果。平台需要一份可筛选、可分页查看的调用记录，以支持设备服务的排障与审计。

## What Changes

- 在设备详情新增物模型服务调用记录列表，按当前设备展示调用时间、服务标识符、服务名称、输入参数和输出参数。
- 提供按服务标识符和时间范围筛选、服务端分页、JSON 参数详情查看，以及加载、错误和空数据状态。
- 新增设备服务调用记录的只读分页 API 与 `device_service_invocations` 持久化模型。

## Capabilities

### New Capabilities
- `device-management/service-invocation-records`: 查看当前设备物模型服务调用记录，并按服务标识符和时间范围分页筛选。

### Modified Capabilities

- 无。

## Impact

- 后端：调用记录模型、仓储、服务、处理器、路由和 Swagger/OpenAPI 契约。
- 前端：设备详情物模型服务调用 Tab、查询 Hook、表格与中英文文案；通过生成的 API 客户端访问接口。
- 数据库：新增 `device_service_invocations` 表及设备/服务/时间查询索引。
