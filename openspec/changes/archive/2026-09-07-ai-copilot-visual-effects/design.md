## Context

基于 `assistant-ui` 官方组件体系和 `0things` 前端规范，为了在物联网控制台提供兼具工业严谨与现代 AI 美学的智能助手交互，引入 [libraries.dev](https://libraries.dev) 的 3 个无外部依赖原生 React 视觉特效库：
- `metal-fx`: WebGL 液态金属反光外圈，支持 `chromatic`、`silver`、`gold` 及暗黑/高亮自动映射。
- `border-beam`: 围绕圆角边界流动的彩色光束。
- `thinking-orbs`: 9 种手调点阵动画思维球，支持 64px 与 20px。

## Goals / Non-Goals

**Goals:**
- 将右下角悬浮按钮升级为 `MetalFx`（`preset="chromatic"`）液态金属风格，搭配状态感知。
- 将 `AssistantModal` 弹出窗体外围包裹 `BorderBeam`（`colorVariant="colorful"`, `size="md"`）。
- 在 `ThreadWelcome` 欢迎区、`Reasoning` 推理折叠条、`ToolGroup` 工具调用块及流式指示器中无缝集成 `ThinkingOrb`。
- 确保所有特效在系统暗色/亮色模式（`useTheme`）下自适应，并在后台/折叠时自动降级或休眠，保障低能耗与高帧率。

**Non-Goals:**
- 不改变已有 MCP 工具通信协议或后端网关逻辑。
- 不影响 assistant-ui 官方核心消息状态流转（Runtime）。

## Decisions

1. **浮动按钮使用 MetalFx `chromatic` 预设**:
   - *Rationale*: 炫彩液态金属反光具有极强的未来科技感与辨识度，且内置 WebGL 离屏共享画布（Shared WebGL Context），多实例消耗极低。
   - *Alternative*: 纯 CSS 渐变或 `silver` 银色预设。`chromatic` 相比单调灰白更具 AI 特色。

2. **弹窗外框采用 BorderBeam `colorful` 全光谱光束**:
   - *Rationale*: 光谱流转为深色暗黑背景提供了清晰而富有活力的窗口边缘，且 `strength: 0.65` 可保证克制不刺眼。
   - *Alternative*: 静态边框或单色脉冲。

3. **思维球状态映射与尺寸规范**:
   - Welcome 欢迎区: `size={64}`, `state="breathing"`（呼吸慢变拟态核心）。
   - Reasoning 推理中: `size={20}`, 运行态 `state="solving"`，折叠完成态 `state="breathing"`。
   - IoT 工具调用中: `size={20}`, `state="searching"`（经纬扫描网格象征 IoT 设备与遥测搜索）。
   - 流式光标指示器: `size={20}`, `state="working"`。

## Risks / Trade-offs

- **[WebGL Context 性能与移动端兼容]** → `metal-fx` 采用单个共享 WebGL Context 并在未激活/不可见时自动 pause，内存占用低；在不支持 WebGL 环境下平稳降级为标准 CSS 背景。
- **[Tailwind CSS 与 Canvas 尺寸匹配]** → 在 `BorderBeam` 和 `MetalFx` 外部包裹元素上确保 `position: relative` 与正确的 `border-radius`（如 `rounded-[2.5rem]`），自动保持尺寸对齐。

