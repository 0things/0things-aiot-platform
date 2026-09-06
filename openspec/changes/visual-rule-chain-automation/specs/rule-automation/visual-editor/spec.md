## Purpose

提供基于 React Flow 的现代规则链工作区，复现 ThingsBoard 的节点、关系、配置和调试心智模型，同时保持 0things 的 shadcn/ui 体验。

## ADDED Requirements

### Requirement: 提供可编辑画布
系统 SHALL 提供平移、缩放、MiniMap、Fit View 和分类节点库。用户 SHALL 能拖入、连接兼容端口、移动、框选、复制粘贴和删除节点/边；画布 SHALL 保存布局并支持撤销/重做和自动布局。

#### Scenario: 条件分支
- **WHEN** 用户把条件的 `True` 和 `False` 输出连到不同节点
- **THEN** 画布 SHALL 显示带关系名的独立边并拒绝不兼容连线

### Requirement: 以统一定义配置节点
画布 SHALL 从后端数据库节点定义目录读取已启用的非系统定义，并严格按 `Filter`、`Enrichment`、`Transformation`、`Action`、`External`、`Flow`、`Analytics` 七类分组。每条画布 SHALL 自动显示不可删除的系统入口节点；选中节点时 SHALL 打开 shadcn 配置抽屉，展示说明、端口、动态表单和校验；前端 MUST NOT 自行硬编码节点目录。

#### Scenario: 物模型条件
- **WHEN** 用户选择数值物模型属性
- **THEN** 表单 SHALL 只展示允许的运算符和阈值输入

### Requirement: Webhook 为首期外部动作
画布 SHALL 提供一个“调用 Webhook”动作节点；用户 SHALL 配置 URL、方法、受限请求体映射、超时、重试和幂等策略，并可填写受保护凭据。短信、邮件及 IM 通知节点不属于首期范围。

#### Scenario: Webhook 配置不完整
- **WHEN** Webhook 节点缺少 URL 或必填凭据
- **THEN** 编辑器 SHALL 标记节点无效并阻止发布

### Requirement: 校验、测试、发布与回放
编辑页 SHALL 提供保存、服务端校验、模拟测试和发布。用户 SHALL 在画布查看节点状态、关系路径、耗时和错误，并进入执行、告警或投递详情。

#### Scenario: 回放 Webhook 失败
- **WHEN** 用户打开一次 Webhook 失败执行
- **THEN** 画布 SHALL 高亮失败 Webhook 节点及 `Failure` 路径，并可跳转到脱敏投递尝试

### Requirement: 可访问性和多语言
编辑器 SHALL 支持键盘选择、删除、可见焦点和可读节点/关系标签；所有界面文字 SHALL 同步提供中文和英文。

#### Scenario: 键盘删除节点
- **WHEN** 键盘用户聚焦节点并按删除键
- **THEN** 系统 SHALL 经确认删除节点和相连边，并将焦点移至可预测画布位置
