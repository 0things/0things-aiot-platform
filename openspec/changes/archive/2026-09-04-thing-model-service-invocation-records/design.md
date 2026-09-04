## Context

设备详情的「物模型数据 → 服务调用」目前为占位卡片。现有设备事件列表已采用 TanStack Table 与服务端分页，并提供时间范围筛选和 JSON 详情弹窗；本变更沿用这一交互模式。动机与范围见 [proposal.md](proposal.md)，行为约束见 [服务调用记录规格](specs/device-management/service-invocation-records/spec.md)。

## Goals / Non-Goals

**Goals:**

- 为单一设备提供稳定、可筛选、可分页的服务调用记录查询契约。
- 保存调用当时的服务名称与 JSON 入/出参，避免 TSL 后续修改影响历史可读性。
- 复用项目既有分页组件和生成式 API 客户端流程。

**Non-Goals:**

- 不提供“立即调用服务”的表单、重试、取消或批量操作。
- 不在主表增加调用状态、失败原因、设备列或可隐藏列。
- 不改变现有物模型事件接口及其数据。

## Decisions

### 设备资源下的只读列表接口

定义 `GET /devices/{deviceKey}/thing-model-service-invocations`，支持 `serviceIdentifier`、`startAt`、`endAt`、`page` 和 `pageSize` 查询参数。路径绑定设备上下文，避免客户端重复传入设备标识，并与既有 `/devices/:deviceKey/*` 子资源保持一致。

响应数据为 `invocations`、`total`、`page`、`pageSize`；每项包含 `uuid`、`invokedAt`、`serviceIdentifier`、`serviceName`、`inputParams`、`outputParams`。时间使用既有 `yyyy-MM-dd HH:mm:ss` 查询格式，页码从 1 开始，默认 `page=1`、`pageSize=20`，页大小范围 1–100。

备选方案是复用全局 `/device-events` 风格的筛选接口；未采用，因为页面只服务当前设备且设备子资源路径更清晰。

### 调用记录数据模型

新增 `device_service_invocations` 表：

| 字段 | 类型/约束 | 用途 |
|---|---|---|
| `id` | bigint 主键 | 内部主键 |
| `uuid` | varchar(36) 唯一且非空 | 对外记录标识 |
| `device_id` | bigint 非空 | 关联设备 |
| `service_identifier` | varchar(128) 非空 | TSL 服务标识符 |
| `service_name` | varchar(255) 非空 | 调用时服务名称快照 |
| `input_params` | text/json 非空 | 输入参数 JSON |
| `output_params` | text/json 可空 | 设备输出参数 JSON |
| `invoked_at` | datetime 非空 | 调用发起时间 |
| `created_at` | datetime | 创建时间 |
| `updated_at` | datetime | 更新时间 |

创建唯一索引 `uk_device_service_invocations_uuid(uuid)`，以及 `device_id, invoked_at DESC` 和 `device_id, service_identifier, invoked_at DESC` 复合索引，分别服务默认分页与标识符筛选。输出尚未产生时保留 `NULL`，而不是引入状态字段，以遵守页面固定五列的产品决定。

### 后端与前端边界

后端按既有 Handler → Service → Repository 分层实现查询，并使用 GORM 模型与生成的 DAL。Swagger 注释定义 API，随后运行现有生成命令更新前端客户端；不得直接修改生成产物。

前端在 `ThingModelDataTab` 中以调用记录组件替换占位卡片。组件维护草稿筛选条件和已应用查询条件；首次加载以最近 7 天请求。参数列使用截断摘要与弹窗详情，分页复用 `DataTablePagination`，不进行客户端全量过滤或分页。

### 调用记录写入边界

本变更包含记录表和读取契约。实际服务下发入口在后续接入时必须在调用发起后写入输入参数，并在收到设备回复后补写输出参数；本变更不新增新的用户发起调用入口。

## Risks / Trade-offs

- [未回复记录无法区分处理中、失败或超时] → 输出参数为 `null`，首期不增加状态语义；后续如需状态，作为独立需求扩展字段与主表。
- [JSON 参数体过大影响列表响应] → 主列表只返回原始字符串并在前端截断预览；调用方应限制单次参数体积，详情按当前页记录查看。
- [TSL 服务被重命名] → 保存 `service_name` 快照，历史记录不回查当前 TSL。
- [页码在数据删除后越界] → 前端在页大小变更时回到第一页；接口返回空数组与实际 total，不将越界视为错误。

## Migration Plan

1. 增加 GORM 模型并执行生成 DAL 命令，随应用迁移创建表和索引。
2. 发布读取 API、Swagger 文档和生成前端 API 客户端。
3. 发布设备详情列表页面与双语文案。
4. 回滚时先撤回前端入口和 API 路由；保留调用记录表，避免删除审计数据。
