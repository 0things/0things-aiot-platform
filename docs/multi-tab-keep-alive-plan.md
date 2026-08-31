# 多标签页与 KeepAlive 计划

## 目标

在 `frontend` 的认证区新增路由标签栏。第一版覆盖设备列表和设备详情：标签可在刷新后恢复，页面组件仅在当前 SPA 会话内保活。缓存容量为 10 个标签，采用 LRU 淘汰策略。

在 TanStack Router 与 React 19 `StrictMode` 下先完成真实组件保活兼容性验证；只有验证通过才合入生产功能，不移除 `StrictMode`。

## 实施内容

- 新增路由标签状态模块，维护 `id`、URL、标题、可关闭状态、最近访问时间及当前激活标签。
  - 首页固定且不可关闭。
  - 设备列表按 pathname 仅保留一个标签，忽略筛选 query。
  - 每个设备详情按完整 pathname 独立建标签，例如不同 `deviceKey` 为不同实例。
  - 切换标签按保存的 URL 导航；关闭当前标签时优先切到右侧标签，否则左侧，最后回首页。
  - 达到 10 个标签时淘汰最久未访问、非首页、非当前的标签，并同步销毁其保活实例。
  - 标签清单写入 `sessionStorage`，刷新恢复；退出登录清空。刷新不会恢复 React 组件实例。

- 在 `frontend/src/components/layout/authenticated-layout.tsx` 中加入全局 `RouteTabs` 和认证区内容缓存容器，保证试点路由处于同一标签与保活边界内。

- 新增设备路由的标签元数据/映射：
  - `/device-management/devices`：标题取现有中英文导航资源，`keepAlive: true`，单标签。
  - `/device-management/devices/$deviceKey`：标题为“设备详情”，`keepAlive: true`，动态路径多实例。
  - 登录、错误、回调及非试点页面不进入标签栏、不参与缓存。

- 以 React 19 兼容的保活核心实现 TanStack Router 适配 POC，缓存键严格使用标签 ID。
  - 保留入口的 `StrictMode`。
  - 不使用仅支持 React Router 的路由适配包，也不手改生成的 `routeTree.gen.ts`。
  - POC 通过后，将其封装为 `KeepAliveRouteContent`，对不缓存路由保持普通 `<Outlet />` 行为。
  - POC 若无法同时满足 StrictMode、路由参数正确性和实例状态保留，则停止接入严格 KeepAlive，并报告兼容性结论；不以移除 StrictMode 作为替代。

- 标签栏提供点击切换、关闭、关闭其他、关闭右侧及横向滚动；复用现有 shadcn UI、Lucide 图标及中英文 locale，不硬编码用户可见文案。

## 验证

- 在设备列表输入筛选、滚动页面或打开局部交互后，进入设备详情并返回；确认列表组件未卸载且状态保持。
- 打开两个不同 `deviceKey` 的详情，分别修改可见局部状态，往返切换后状态互不串扰。
- 关闭标签、LRU 淘汰、退出登录时，对应页面实例被销毁；当前标签与首页不会被错误淘汰。
- 重复进入设备列表不产生重复标签；详情页产生独立标签；浏览器前进/后退与标签激活态一致。
- 刷新后恢复标签列表和当前 URL，但页面状态按预期重建。
- 执行 `pnpm build`、`pnpm lint`、`pnpm format:check`，并完成浏览器验收。

## 边界与默认值

- 第一版仅覆盖设备列表和设备详情；OTA 等其他业务页在 POC 成功后按同一元数据机制逐步接入。
- “KeepAlive”指当前浏览器会话中的真实 React 组件实例保活，不包含跨刷新恢复未提交表单内容。
- 缓存总量最多 10 个，首页固定，未引入用户可配置的上限或缓存开关。
