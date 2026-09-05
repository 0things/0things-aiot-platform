## Why

设备详情的属性页目前错误地以设备影子或未定义结构的原始遥测 JSON 推导属性和值，导致页面数据源与物模型定义脱节，也不能稳定呈现单位、读写方式和最后上报时间。属性页需要以产品 TSL 为权威定义，并为每个定义属性提供时序库中的最后一个上报值。

## What Changes

- 新增设备物模型属性最后值查询能力，按产品 TSL 的属性定义返回结构化属性列表及各属性最后一个时序点。
- 在设备详情属性页使用该专用接口展示属性名称、标识符、数据类型、单位、读写方式、当前值和最后上报时间。
- 保留现有时序历史接口 `GET /v1/devices/:deviceKey/telemetry/history` 用于单个属性的趋势图查询。
- 移除属性页对设备影子和原始最新遥测 JSON 接口的依赖。

## Capabilities

### New Capabilities

- `device-management/thing-model-properties`: 查询设备产品 TSL 中定义的属性及其最后上报值，并在设备详情中展示。

### Modified Capabilities

- 无。

## Impact

- 后端：新增属性最后值 DTO、Handler、Service 与时序数据访问能力，新增设备路由并更新 Swagger。
- 前端：替换设备详情属性页的数据请求与呈现逻辑，更新中英文文案并重新生成 API 客户端。
- 数据源：产品 TSL 提供属性定义，时序存储提供最后上报值；设备影子不参与该页面。
