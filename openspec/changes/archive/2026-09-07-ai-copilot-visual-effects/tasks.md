## 1. 依赖与基础特效接入

- [x] 1.1 安装并确认 `border-beam`, `thinking-orbs`, `metal-fx` 依赖已就绪并完成前端类型适配
- [x] 1.2 确认特效组件利用内置 DOM 变动监听（MutationObserver）无缝同步系统 Dark/Light 切换

## 2. 悬浮触发入口与主窗口特效增强

- [x] 2.1 将 `AssistantModalButton` 接入 `MetalFx`（`preset="chromatic"`, `variant="circle"`），并在按钮状态中集成 `ThinkingOrb`
- [x] 2.2 在 `AssistantModal` 弹窗外框集成 `BorderBeam`（`colorVariant="colorful"`, `size="md"`, `strength={0.65}`），实现全光谱流转光束

## 3. 会话与工具状态动效优化

- [x] 3.1 在 `ThreadWelcome` 欢迎区居中引入 64px 呼吸思维球（`ThinkingOrb state="breathing"`）
- [x] 3.2 在 `ReasoningTrigger` / `group-reasoning` 中集成 20px 思维球（运行态 `state="solving"`, 完成态 `state="breathing"`）
- [x] 3.3 在 `ToolGroupTrigger` / `ToolFallback` 中集成 20px 搜索思维球（运行态 `state="searching"`）
- [x] 3.4 替换 `AssistantMessage` 流式字符指示器为微型点阵脉冲

## 4. 全链路构建与视觉验证

- [x] 4.1 运行 `pnpm -C frontend format && pnpm -C frontend build` 验证 TypeScript 编译与构建无误
- [x] 4.2 验证各页面下暗色/亮色切换、浮动按钮悬停/展开、思考生成与工具调用的动效流畅度

