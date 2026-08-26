# 仓库协作指南

## 项目结构与模块组织

本仓库当前维护前端和后端两个部分：

- `frontend/`：Vite + React 19 + TypeScript 管理后台。路由位于 `src/routes/`，业务模块位于 `src/features/`，通用 UI 位于 `src/components/`，API 客户端位于 `src/api/`，中英文资源位于 `public/locales/{zh,en}/`。
- `backend/`：Go + Gin REST API。启动入口在 `cmd/server/`，业务代码在 `internal/`，公共包在 `pkg/`；单元测试分布在 `internal/**/*_test.go`，服务级测试位于 `test/server/`，Swagger 产物在 `docs/`。

前端 `src/api/generated/` 和后端 Swagger 文档不可手工修改；契约变化后使用对应生成命令更新，并检查生成差异。

## 构建、测试与开发命令

请在对应包目录下执行命令。首次运行后端前，先准备本地配置和依赖服务：

```bash
cd frontend && pnpm dev          # 启动本地 Vite 开发服务器
cd frontend && pnpm build        # TypeScript 类型检查并产出生产构建
cd frontend && pnpm lint         # 运行 ESLint
cd frontend && pnpm format:check # Prettier 格式校验
cd frontend && pnpm generate:api # 根据 orval.config.ts 重新生成 API 客户端
cd backend && make test          # 运行 Go 服务端测试并生成覆盖率报告
cd backend && make build         # 构建 bin/server
cd backend && make gen           # 生成 GORM DAL 代码
cd backend && make swag          # 根据 handler 注释生成 Swagger 文档
```

后端完整测试由 `make test` 调用 `go test ./test/server/...`，修改业务代码时还应先运行受影响包的 `go test ./internal/...`。前端没有专用 `test` script，关键流程需要浏览器验收；`pnpm build` 会执行 `tsc -b` 和 Vite 生产构建。

本地启动通常为：

```bash
cd backend && go run ./cmd/server -conf ./config/local.yml
cd frontend && pnpm dev
```

后端也提供 `make bootstrap`（启动 Docker Compose 依赖、执行迁移并启动服务），使用前确认本地配置可用。

## 代码风格与命名约定

沿用既有格式：TypeScript 使用 2 空格缩进、单引号、不加分号；Go 使用 `gofmt`。React 组件使用 `PascalCase`，Hook 文件沿用现有 `use-*.ts` 命名，特性目录使用 kebab-case，Go 文件使用小写下划线。UI 文案需要同时维护中英文资源，跨 namespace 使用 `namespace:key`（例如 `common:createdAt`），不要硬编码特定语言。优先复用现有 shadcn 组件和 Lucide 图标，再考虑新增依赖。

## 测试指南

提交前先运行最小范围的检查，再跑整个包的构建。Go 测试以 `*_test.go` 文件名放在被测包旁边；服务级测试沿用 `test/server` 的现有约定。前端目前没有专用测试脚本，已有集成测试文件仍需按项目配置执行，并在本地浏览器验证可见流程。前端至少执行受影响文件的 Prettier 检查，条件允许时执行 `pnpm build`、`pnpm lint` 和 `pnpm format:check`；后端条件允许时执行 `make test` 和 `make build`。

## 提交与 Pull Request 规范

提交说明使用简洁的 Conventional Commit 主题，例如 `fix: preserve OTA batch history` 或 `chore: update project configuration`。每次提交只关注一个主题；用户明确要求“全部提交”时才合并无关改动。PR 应说明受影响的服务或路由、描述行为变更、关联相关 Issue，并附带改动前后的截图（针对可见的前端变更）。切勿提交凭据、`config/local.yml`、`.env.local`、本地工具目录、覆盖率/构建产物或生成文件，除非本次提交的目的就是重新生成这些文件。推送后需核对远程分支和提交哈希。


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
