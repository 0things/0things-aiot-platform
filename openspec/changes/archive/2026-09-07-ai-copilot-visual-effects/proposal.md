## Why

随着 0things AI Copilot 与 IoT MCP 生态全链路贯通，AI 问答交互已具备多步工具调用和数据汇总能力。为了让 AI Copilot 在视觉品质、科技感和动态反馈层面达到行业顶级水准，引入 [libraries.dev](https://libraries.dev) 视觉特效组件体系（`border-beam`, `thinking-orbs`, `metal-fx`），为浮动入口、主会话窗口、思考推理过程、工具调用及欢迎屏提供富有生命力的视觉呈现与状态反馈。

## What Changes

- **浮动入口按钮质感增强**：使用 `metal-fx`（`preset="chromatic"`，炫彩流体液态金属反光与灵动光晕）包裹右下角浮动 Copilot 触发入口，结合 `thinking-orbs` 实现状态响应。
- **对话弹窗光束外框**：在 `AssistantModal` 主容器外周集成 `border-beam`（`colorVariant="colorful"`，全光谱流转光束），赋予助手窗口现代灵动的流光科技边界。
- **动态思维球与思考状态**：
  - 欢迎屏（Welcome Hero）居中展示 64px 呼吸思维球（`state="breathing"`）。
  - 推理状态（Reasoning）采用 20px 解算思维球（`state="solving"`）。
  - IoT 工具调用（Tool Calls）运行态展示 20px 扫描思维球（`state="searching"`）。
  - 流式指示器（Indicator）与待机态融入点阵动效。
- **主题与系统性能适配**：全量效果自适应 Dark/Light 主题切换，利用 GPU 硬件加速并自带离屏自动休眠机制。

## Capabilities

### New Capabilities
- `ai-copilot-visual-effects`: 覆盖 AI Copilot 的 libraries.dev 视觉特效层，包含 MetalFx 触发按钮、BorderBeam 会话外框、ThinkingOrb 多状态思考动画规范与主题响应。

### Modified Capabilities
<!-- 无现有已固化主规范需求变更 -->

## Impact

- **前端依赖**：新增 `border-beam`, `thinking-orbs`, `metal-fx` 运行时依赖（无第三方二次依赖，纯轻量级现代 React 组件）。
- **前端组件**：优化 `frontend/src/components/assistant-ui/elements/`（`assistant-modal.aui.tsx`, `thread.aui.tsx`, `reasoning.aui.tsx`, `tool-group.aui.tsx`）及 `frontend/src/features/ai-copilot/`。
- **后端/API**：无影响，纯前端视觉与状态呈现优化。

