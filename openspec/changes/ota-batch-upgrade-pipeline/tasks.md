## 1. 协议与基础设施

- [ ] 1.1 固化 `ota.upgrade.command.v1`、`ota.upgrade.report.v1` 和 DLQ 的 JSON schema，覆盖 batch_id、设备标识、默认模块、目标版本、下载校验和错误字段，并用 schema 示例校验通过
- [ ] 1.2 部署 EMQX MQTT Broker，配置 TLS、设备认证和按 product/device 隔离的 ACL，使用测试客户端验证未授权 topic 不能发布或订阅
- [ ] 1.3 增加 backend 与 OTA transport worker 的 Kafka、EMQX、TLS 和重试配置，启动配置检查并验证缺失凭据时服务安全失败
- [ ] 1.4 引入并锁定 Eclipse Paho Go MQTT 客户端，完成 QoS 1 发布、订阅、重连和优雅关闭测试

## 2. OTA 数据模型与批次事务

- [ ] 2.1 扩展 OTA 设备任务模型，增加 batch_id、默认模块、目标版本、发送尝试、超时、错误和最后上报字段，并添加 `(batch_id, device_id)` 唯一约束及状态/重试索引；执行迁移测试确认旧数据可读
- [ ] 2.2 为 OTA 设备任务增加 `(batch_id, device_id, target_version)` 幂等约束及 Kafka 发布结果字段，验证批次和设备任务在同一事务中提交或整体回滚
- [ ] 2.3 改造 BatchUpgrade，为每个设备创建独立批次任务，保留同设备跨批次历史，并验证重复设备输入不会生成重复任务
- [ ] 2.4 实现批次和设备任务状态聚合，覆盖成功、部分成功、失败、超时和取消场景；运行后端 OTA service 测试验证聚合结果

## 3. Kafka 发布与 MQTT 传输

- [ ] 3.1 在 BatchUpgrade 事务提交后直接发布 Kafka OTA 命令，记录发布结果并返回失败信息；使用 Kafka 测试容器或 mock 验证成功和失败路径
- [ ] 3.2 实现 OTA transport worker，消费命令并向 `/ota/device/upgrade/{productKey}/{deviceName}` 发布升级信息；用 EMQX 测试客户端验证在线设备收到完整 payload
- [ ] 3.3 实现 `(batch_id, device_id, target_version)` 幂等处理和 DLQ，重复消费不得重复发布或创建任务；执行重复消息集成测试
- [ ] 3.4 记录设备离线时的未完成任务、有效期和失败原因，暂不实现后台自动扫描或自动重试；通过查询接口验证任务仍可被显式重试

## 4. 设备上报与生命周期

- [ ] 4.1 消费 `/ota/device/progress/{productKey}/{deviceName}` 并转换为 report 事件，按 batch_id + device_id 更新进度、阶段和错误；验证下载、校验、烧写失败码映射
- [ ] 4.2 消费 `/ota/device/inform/{productKey}/{deviceName}` 的默认模块版本上报，只有版本等于目标版本时才标记任务成功；验证 100% 进度但版本不匹配仍非成功
- [ ] 4.3 记录首次进度时间、设备升级超时、最大重试次数和取消任务，不实现自动重试；使用时间控制测试验证超时和取消状态
- [ ] 4.4 停止旧的“定时任务直接 pending → in_progress”语义，移除 OTA 定时扫描；运行 OTA 回归测试确认数据库状态只在实际 Kafka 发布或设备上报后变化

## 5. API、前端与观测

- [ ] 5.1 扩展 BatchUpgrade、report、批次详情、设备列表和统计 API，增加 batch_id 及批次过滤；运行 Swagger 生成和 handler 测试验证请求/响应兼容性
- [ ] 5.2 更新前端 OTA 批次详情和设备部署列表，按 batch_id 展示状态、进度、错误、当前版本和目标版本；执行 `pnpm build` 和本地页面验证
- [ ] 5.3 增加批次取消、失败重试和未完成任务查询能力，维护中英文 locale；执行 API、前端 lint 和格式检查
- [ ] 5.4 增加 batch_id、device_key、Kafka offset 和 MQTT packet id 的结构化日志、指标与告警；通过测试发布验证链路可追踪
- [ ] 5.5 执行端到端测试：创建批次 → Kafka → EMQX → 设备收到命令 → 进度上报 → 版本确认 → 批次成功，并保存测试报告
