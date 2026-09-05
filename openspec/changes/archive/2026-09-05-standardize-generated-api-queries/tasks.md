## 1. 设备模块重构 (Devices)

- [x] 1.1 重构 `src/features/devices/api/queries.ts`，导出 `GetDevicesParams` 并使用 `useGetDevices`、`useGetDevicesDeviceKey`、`useGetDevicesDeviceKeyTelemetry`、`useGetDeviceStatistics` 与生成的 Query Key 工厂函数。
- [x] 1.2 重构 `src/features/devices/api/` 下的子查询文件（`push-records.ts`、`tags.ts`、`telemetry.ts`、`endpoints.ts`、`shadow.ts`），替换为系统生成的 Hook 与类型。
- [x] 1.3 规范化 `src/features/devices/components/devices-provider.tsx` 中的弹窗 Context Hook 命名（`useDevicesDialog`），并更新相关弹窗组件的引用。
- [x] 1.4 适配设备表格、详情页各个 Tab 及对话框组件，传递符合规范的查询参数并验证设备模块 TypeScript 类型检查通过。

## 2. 产品模块重构 (Products)

- [x] 2.1 重构 `src/features/products/api/queries.ts`，导出 `GetProductsParams` 并使用 `useGetProducts`、`useGetProductsProductKey` 与生成的 Query Key 工厂函数。
- [x] 2.2 重构 `src/features/products/api/categories.ts` 及产品详情页子 Tab，全面使用自动生成的 Hook 与类型。
- [x] 2.3 适配产品列表表格及消费组件，验证产品模块 TypeScript 类型检查通过。

## 3. 其他业务模块标准化 (Other Features)

- [x] 3.1 重构设备分组模块（`src/features/device-groups/`），使用 `GetDeviceGroupsParams` / `useGetDeviceGroups` 及生成的 Hook。
- [x] 3.2 重构场景联动模块（`src/features/scene-linkage/api/queries.ts`），全面接入生成的参数类型与 Hook。
- [x] 3.3 重构告警中心模块（`src/features/alert-center/api/queries.ts`），全面接入生成的参数类型与 Hook。
- [x] 3.4 重构运维监控 OTA 模块（`src/features/operations-monitoring/ota/`），全面接入生成的参数类型与 Hook。

## 4. 清理与全量验证 (Verification & Cleanup)

- [x] 4.1 清理各业务模块中冗余废弃的手写查询参数类型与不安全的类型强转。
- [x] 4.2 执行 `pnpm -C frontend format && pnpm -C frontend build` 验证全局类型安全、代码格式与生产构建完全通过。
