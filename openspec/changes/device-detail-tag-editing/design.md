## Context

设备详情 `device-info-tab.tsx` 的「标签信息」卡片当前是占位（静态文案 + 无 `onClick` 的「编辑」按钮）。已存在可复用的 `TagsEditor` 组件（`components/device-detail/tags-editor.tsx`）与 `api/tags.ts`（`useDeviceTags` / `useAddTags` / `useRemoveTags`），后端 `device_tags` 增改删接口已完备。本变更只做前端接线，不改动后端。

## Goals / Non-Goals

**Goals:**
- 在「标签信息」卡片渲染 `TagsEditor`，让用户能给设备打/删自定义 tag。
- 移除当前无作用的「编辑」按钮，改为内联常显编辑态。

**Non-Goals:**
- 运维/契约层面：不改动后端对外接口路径、请求/响应结构、数据模型或 `source` 语义。
- 不处理设备上报标签的展示区分（本变更仅面向自定义标签的增删）。
- 不引入新依赖。

**注意（范围内）**：后端标签接口此前存在服务层 bug（见下方决策），修复该 bug 属于本变更范围。

## Decisions

- **内联常显而非「查看→编辑」切换**：卡片直接展示标签列表与增删输入，去掉原死按钮。理由：交互更轻，且 `TagsEditor` 本身已内置列表+输入，无需再包一层切换状态。
- **复用既有 `TagsEditor`**：不重写组件，仅把它接进卡片并传入 `device.deviceKey`（`TagsEditor` 内部通过 `getDeviceId` 解析为设备 id 调用接口）。
- **后端标签接口按数字 id 解析（bug 修复）**：`Tags`/`SetTags`/`RemoveTags` 原先把路由里的 `:id`（数字设备 id）传给 `DeviceByKey`（按 `device_key` 字符串查），导致接口返回 404。前端 `api/tags.ts` 实际传的是数字 id（先经 `getDeviceId` 解析），swagger 也声明 param 为 `id int`。故改为 `strconv.ParseInt` 后调用 `s.Device(ctx, id)`，并抽出一个 `deviceByIDParam` helper 供三处复用。接口契约（路径/入参/出参）不变。
- **标签键禁止纯数字**：在 `SetTags` 服务层增加校验，key 若仅由 `0-9` 组成则拒绝（错误 `tag key cannot be purely numeric`）。POST 与 PUT 写入路径均经 `SetTags`，故统一覆盖。前端 `tags-editor.tsx` 在提交前用正则 `/^\d+$/` 即时拦截并提示（i18n `deviceDetail.tags.invalidKey`），避免无谓请求。

## Risks / Trade-offs

- [输入法/重复 key] `TagsEditor` 对重复 key 走后端 upsert（更新 value），前端未做 key 格式校验 → 可接受，沿用既有行为，不额外加校验。
- [额外一次解析请求] `TagsEditor` 经 `getDeviceId(deviceKey)` 多一次请求解析设备 id；设备详情页已有 `device.id`，后续若需优化可直接传 id，本变更不做。

## Open Questions

- `tags-editor.tsx` 内硬编码英文文案是否并入本次迁到 i18n？当前不影响功能，可在本变更内一并处理或单独跟进。
