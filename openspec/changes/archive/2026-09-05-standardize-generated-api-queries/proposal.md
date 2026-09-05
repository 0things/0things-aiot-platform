## Why

目前，前端多个业务模块（如 `devices`、`products`、`device-groups`、`scene-linkage`、`alert-center` 等）仍然采用手工定义查询参数接口、重复编写 `useQuery` 包装逻辑以及手工进行 `as unknown as ...` 类型强转的方式。而较新的 `events` 模块则直接复用了由 Orval 根据后端 Swagger 契约自动生成的参数类型（如 `GetDeviceEventsParams`）与 React Query Hooks（如 `useGetDeviceEvents`）。

将所有前端业务模块的 API 查询参数与 Hook 全面标准化为使用自动生成的参数与 Hook，能够彻底消除冗余的样板代码，确保前后端数据类型（如数字型 ID 与字符串）的严格一致性，避免因参数格式不匹配引发的潜在缺陷。

## What Changes

- **统一查询参数类型**：将所有业务模块中的手写参数类型（如 `params: { page?: number; ... }`）替换为 Orval 生成的类型（如 `GetDevicesParams`、`GetProductsParams`、`GetDeviceGroupsParams`、`ListSceneLinkagesParams`、`GetAlertsParams` 等）。
- **统一 React Query Hooks**：重构各业务模块的查询 Hook（`useDevices`、`useProducts`、`useDeviceGroups`、`useSceneLinkages`、`useAlerts` 等），基于 Orval 生成的 Hooks 进行二次包装，并统一通过 `select: (res) => res?.data` 解包数据。
- **统一 Query Key 管理**：在各模块的 `keys` 对象中，统一使用 Orval 生成的 Query Key 工厂函数（如 `getGetDevicesQueryKey`、`getGetProductsQueryKey`）。
- **消费组件类型对齐**：更新表格、筛选栏和弹窗等 UI 组件，传递符合生成类型规范的入参。
- **区分 Provider 弹窗 Hook 命名**：将 `devices-provider.tsx` 中用于弹窗状态的 Context Hook 进行区分命名（如 `useDevicesDialog`），避免与数据查询 Hook `useDevices` 产生命名冲突。

## Capabilities

### New Capabilities
<!-- 无新增业务 Spec 变更：属于前端纯架构重构 -->

### Modified Capabilities
<!-- 无已有业务 Spec 变更：行为与业务功能保持不变 -->

## Impact

- **前端受影响模块**：`src/features/devices/`、`src/features/products/`、`src/features/device-groups/`、`src/features/scene-linkage/`、`src/features/alert-center/`、`src/features/operations-monitoring/`
- **后端/API**：无影响（已有 Swagger 契约完备）。
- **破坏性变更**：无（对最终用户无感知，仅为前端代码重构与类型收敛）。
