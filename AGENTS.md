## 1. 项目结构与核心命令

```bash
# 核心验证（无需切目录，根目录直接执行）
make test && make build                             # 验证所有微服务与共享包
pnpm -C frontend format && pnpm -C frontend build   # 前端代码格式化、类型检查与生产构建

# 契约与代码生成
make -C backend swag && pnpm -C frontend generate:api # 同步 Swagger 文档与前端 API 客户端
make -C backend gen                                   # 生成 GORM DAL 代码
```

## 2. 核心架构与分层铁律

### 依赖方向与生成代码禁令
- **单向向内依赖**：`Router/Handler` $\rightarrow$ `Service` $\rightarrow$ `Repository` (持数据库连接)。禁止跨层逆向依赖（如 Repository 依赖 API 契约）。
- **禁止手工修改生成产物**：
  - 前端 `src/api/generated/`
  - 后端 `wire_gen.go`、`internal/dal/query/*.gen.go`、`docs/` (Swagger)
  - 契约或依赖变动时，必须通过对应命令重新生成并验证差异。

## 3. 代码风格与多语言规范
- **多语言（i18n）**：UI 新增或调整文案必须同步维护 `public/locales/zh/*.json` 与 `public/locales/en/*.json`，跨 namespace 引用遵循 `namespace:key`。
- **必要英文注释**：涉及复杂业务逻辑、状态流转、异常分支或关键架构约束处补充简明英文注释，自解释代码不加冗余注释。

## 4. Agent 行为准则
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