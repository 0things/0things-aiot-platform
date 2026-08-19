# 仓库协作指南

## 项目结构与模块组织

本仓库包含三个可独立部署的部分：

- `frontend/`：基于 Vite + React + TypeScript 的管理后台 UI。业务特性代码放在 `src/features/`，通用组件放在 `src/components/`，API 客户端放在 `src/api/`，翻译文件位于 `public/locales/{zh,en}/`。
- `backend/`：Go 编写的 HTTP 服务。业务代码统一放在 `cmd/`、`internal/` 下，可复用包放在 `pkg/`；服务端的测试位于 `test/server/`。
- `telemetry-service/`：基于 Go/Kratos 的遥测服务。Protobuf 定义在 `api/` 和 `internal/` 中；生成代码与源码放在同一目录。

请勿手工编辑 `frontend/src/api/generated/` 中的生成代码；契约变更时应通过 OpenAPI 源重新生成。

## 构建、测试与开发命令

请在对应包目录下执行命令：

```bash
cd frontend && pnpm dev          # 启动本地 Vite 开发服务器
cd frontend && pnpm build        # TypeScript 类型检查并产出生产构建
cd frontend && pnpm lint         # 运行 ESLint
cd frontend && pnpm format:check # Prettier 格式校验
cd backend && make test          # 运行 Go 服务端测试并生成覆盖率报告
cd backend && make build         # 构建 bin/server
cd telemetry-service && make build    # 构建遥测服务可执行文件
```

在 `telemetry-service/` 下修改 protobuf 或生成的 Go 代码后，请运行 `make api`、`make config` 或 `make generate`。

## 代码风格与命名约定

沿用既有格式：TypeScript 使用 2 空格缩进、单引号、不加分号；Go 使用 `gofmt`。React 组件使用 `PascalCase`，自定义 Hook 以 `use-*.ts` 命名，特性目录使用 kebab-case，Go 文件使用小写下划线。UI 文案需要同时维护在两个语言文件中，不要硬编码特定语言。优先复用 shadcn 组件和 Lucide 图标，再考虑新增依赖。

## 测试指南

提交前先运行最小范围的检查，再跑整个包的构建。Go 测试以 `*_test.go` 文件名追加在源码旁边，整体组织沿用 `test/server` 的现有约定。前端目前没有专门的单元测试脚本，请在本地浏览器中验证改动流程，并在条件允许时执行 `pnpm build`、`pnpm lint` 和 `pnpm format:check`。

## 提交与 Pull Request 规范

提交说明使用简洁的 Conventional Commit 主题，例如 `fix: unify list page scrolling` 或 `chore: update project configuration`。每次提交只关注一个主题。PR 应说明受影响的服务或路由、描述行为变更、关联相关 Issue，并附带改动前后的截图（针对可见的前端变更）。切勿提交凭据、本地工具目录、构建产物或生成文件，除非本次提交的目的就是重新生成这些文件。


Behavioral guidelines to reduce common LLM coding mistakes. Merge with project-specific instructions as needed.

Tradeoff: These guidelines bias toward caution over speed. For trivial tasks, use judgment.

1. Think Before Coding
Don't assume. Don't hide confusion. Surface tradeoffs.

Before implementing:

State your assumptions explicitly. If uncertain, ask.
If multiple interpretations exist, present them - don't pick silently.
If a simpler approach exists, say so. Push back when warranted.
If something is unclear, stop. Name what's confusing. Ask.
2. Simplicity First
Minimum code that solves the problem. Nothing speculative.

No features beyond what was asked.
No abstractions for single-use code.
No "flexibility" or "configurability" that wasn't requested.
No error handling for impossible scenarios.
If you write 200 lines and it could be 50, rewrite it.
Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

3. Surgical Changes
Touch only what you must. Clean up only your own mess.

When editing existing code:

Don't "improve" adjacent code, comments, or formatting.
Don't refactor things that aren't broken.
Match existing style, even if you'd do it differently.
If you notice unrelated dead code, mention it - don't delete it.
When your changes create orphans:

Remove imports/variables/functions that YOUR changes made unused.
Don't remove pre-existing dead code unless asked.
The test: Every changed line should trace directly to the user's request.

4. Goal-Driven Execution
Define success criteria. Loop until verified.

Transform tasks into verifiable goals:

"Add validation" → "Write tests for invalid inputs, then make them pass"
"Fix the bug" → "Write a test that reproduces it, then make it pass"
"Refactor X" → "Ensure tests pass before and after"
For multi-step tasks, state a brief plan:

1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.