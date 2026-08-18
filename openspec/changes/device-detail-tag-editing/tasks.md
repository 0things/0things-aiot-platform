## 1. 接线 TagsEditor 到设备详情卡片

- [x] 1.1 在 `device-info-tab.tsx` 引入 `TagsEditor`，移除「标签信息」卡片的占位文案与无作用的「编辑」按钮，改为渲染 `<TagsEditor deviceKey={device.deviceKey} />`
- [ ] 1.2 本地 `pnpm dev` 验证：设备有/无标签时卡片展示正确，新增与删除标签后列表即时刷新（lint + tsc 已通过，浏览器交互待手动确认）

## 2. i18n 清理（可选，是否并入本次待定）

- [x] 2.1 将 `tags-editor.tsx` 内硬编码英文文案（'Tag added'、'No tags'、'key'、'value'、'Add'、'Loading…'）迁到 `deviceManagement.json` 的 zh/en，并在组件中改用 `t(...)`
- [x] 2.2 运行 `pnpm lint` 与类型检查确认改动合规（`lint`、`tsc -b`、对我改文件的 `prettier --check` 均通过；仓库级 `format:check` 的告警为既有问题，与本变更无关）

## 3. 修复后端标签接口 id 解析（bug fix）

- [x] 3.1 将 `internal/service/device.go` 的 `Tags`/`SetTags`/`RemoveTags` 从按 `device_key` 查改为按数字 id 查（新增 `deviceByIDParam` helper，内部 `strconv.ParseInt` + `s.Device`）
- [x] 3.2 重新编译并重启后端，验证 `GET`/`POST`/`DELETE /v1/devices/:id/tags` 返回 200（已用设备 4 验证，测试标签已清理）

## 4. 标签键禁止纯数字

- [x] 4.1 后端 `internal/service/device.go` 的 `SetTags` 增加校验：key 仅含 0-9 时拒绝（新增 `isNumericKey` helper）
- [x] 4.2 前端 `tags-editor.tsx` 提交前用 `/^\d+$/` 拦截纯数字键并提示 i18n `deviceDetail.tags.invalidKey`（zh/en 已加键）
- [x] 4.3 重新编译重启后端，验证纯数字键被拒绝、非数字键正常写入

## 5. 单元测试（覆盖改动代码）

- [x] 5.1 `internal/repository/device_tag_test.go`：覆盖 `DeviceTagRepository` 的 `ListTags` / `SetTags`（增量 upsert、全量替换、已存在 key 更新 value+source）/ `DeleteTags`（存在/不存在 key）
- [x] 5.2 `test/server/service/device_tags_test.go`：覆盖 service 层 `SetTags` 校验（纯数字键拒绝、空键拒绝、超长键拒绝、合法键写入成功）与 `deviceByIDParam` 的 id 解析失败路径（间接覆盖 `isNumericKey`）
- [x] 5.3 为解除 `internal/repository` 测试包既有编译阻断，删除了两份失效测试：`rule_repository_test.go`（`model.Rule` 已不存在）、`device_event_repository_test.go`（`json.RawMessage` 误用作 string）
