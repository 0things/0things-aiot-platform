## Why

设备详情「标签信息」卡片目前只是占位（静态文案 + 一个无 `onClick` 的「编辑」按钮），而已写好的 `TagsEditor` 组件并未被引用。后端 `device_tags` 的增改删接口（GET/POST/PUT/DELETE `/devices/{id}/tags`）已齐备，但用户在前端无法给设备打 tag。需要在设备详情页把标签编辑能力真正暴露出来。

## What Changes

- 在 `device-info-tab.tsx` 的「标签信息」卡片中渲染已有的 `TagsEditor` 组件（传入 `device.deviceKey`），替换占位文案。
- 处理卡片「编辑」按钮交互：本提案采用「内联常显」——直接展示标签列表与增删输入，去掉当前无作用的「编辑」按钮。
- 修复后端标签接口的 id 解析：`internal/service/device.go` 的 `Tags`/`SetTags`/`RemoveTags` 原先把 URL 里的数字设备 id 误传给 `DeviceByKey`（按 `device_key` 字符串查），导致接口 404；改为按数字 id 查（`s.Device(ctx, id)`）。
- （可选，是否并入本次待定）将 `tags-editor.tsx` 内硬编码的英文文案迁到 i18n（`deviceManagement.json` 的 zh/en 两语言），以符合仓库约定。

## Capabilities

### New Capabilities
- `device-management/device-tags`: 设备详情页的标签查看与编辑能力（新增、删除自定义标签），基于已存在的 `device_tags` 后端接口。

### Modified Capabilities
（无——后端接口与行为不变，仅前端暴露既有能力）

## Impact

- 前端：`src/features/devices/components/device-detail/device-info-tab.tsx`、`tags-editor.tsx`（已存在，基本无需改）、`api/tags.ts`（已存在，无需改）、i18n `public/locales/{zh,en}/deviceManagement.json`（仅当并入 i18n 清理时改动）。
- 后端：`internal/service/device.go`（标签接口服务层按 id 解析的设备查询）。
- 依赖：沿用既有 TanStack Query、shadcn/ui、Lucide，不引入新依赖。
