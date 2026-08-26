## Context

当前 backend 已有 OTA 批次、设备升级记录和 Kafka 基础能力，但 OTA 定时任务只修改数据库状态，缺少实际消息下发和设备上报闭环。`telemetry-service` 计划迁移，不作为本设计的运行时依赖。详细目标和行为约束见 `proposal.md` 以及两个 capability spec。

## Goals / Non-Goals

**Goals:**

- 建立批次、设备任务、传输消息和设备回报之间的稳定关联。
- 在批次事务提交后直接发布 Kafka 命令，保持首期实现简单。
- 通过 MQTT 适配器支持在线推送、离线补发、幂等和失败重试。
- 以设备最终版本上报作为成功判定依据，并提供批次级查询和统计。
- 保留现有 MQTT topic 约定，默认采用 MQTT 下发元数据、HTTPS 下载固件。

**Non-Goals:**

- 本变更不实现设备固件、Bootloader 或具体 MCU 下载代码。
- 首期不实现动态分组的持续匹配，只为后续扩展保留策略字段。
- 不把 `telemetry-service` 改造成 OTA 服务，也不依赖其数据库或 API。

## Decisions

### 1. 使用设备任务作为状态聚合单元

扩展现有 OTA 设备升级记录，使每条记录绑定 `batch_id` 和设备，并保存目标版本、模块、发送尝试、超时和错误信息。通过 `batch_id + device_id` 唯一定位批次内任务，批次状态由设备任务聚合得出。

备选方案是继续按 `ota_package_id + device_id` 保存最新状态，但无法保留同一设备多批次历史，也无法处理并行任务，因此不采用。

### 2. 使用现有设备状态表驱动直接发布

批次和设备任务在一个数据库事务中提交。事务成功后，BatchUpgrade 直接向 `ota.upgrade.command.v1` 发布每台设备的升级命令；Kafka 成功后将设备任务标记为已发送，失败则保留待发送/失败状态并返回错误。首期不新增 outbox 表，也不使用后台扫描器自动重试。

备选方案是事务 outbox，可靠性更高但会新增表、发布器和运维复杂度；当前首期选择直接发布，后续需要可靠重试时再引入 outbox。

### 3. 使用成熟 Broker，并将 Kafka-MQTT 转换隔离为 OTA transport worker

部署 EMQX 作为 MQTT Broker，使用其 TLS、认证和 ACL 能力隔离产品与设备 topic；OTA transport worker 使用成熟的 Eclipse Paho Go MQTT 客户端，消费 Kafka OTA 命令并向 `/ota/device/upgrade/{productKey}/{deviceName}` 发布 MQTT 消息。设备的 progress/inform 消息由同一适配层转换为 `ota.upgrade.report.v1`，由 backend OTA consumer 处理。该 worker 可以作为 backend 的独立 command 部署，也可以独立服务部署，但不得依赖 `telemetry-service`。

EMQX 的认证/授权支持用户名密码、JWT、X.509 等方式，并提供细粒度 ACL；Paho Go 支持 MQTT 发布和订阅。备选方案是让 backend HTTP 进程直接维护 MQTT 长连接，或自研 Broker；前者会耦合 API 与设备连接，后者会引入不必要的协议和运维风险，均不采用。

### 4. 以批次和设备组合做幂等键

Kafka 消息包含 `batch_id` 和设备标识；生产和消费均以 `batch_id + device_id + target_version` 去重，Kafka key 使用 `batch_id + device_id` 以保证同一批次设备顺序。设备收到重复命令时按批次和设备返回当前状态，不重复安装。

### 5. 采用 MQTT 通知、HTTPS 固件下载

升级通知中携带版本、大小、URL、摘要/签名和过期时间，设备通过 HTTPS 下载并在安装前校验。MQTT 文件下载仅作为设备明确支持且满足单文件大小限制时的可选协议。

### 6. 用版本上报完成闭环

进度消息只更新阶段和百分比；设备重启后在 inform topic 上报模块版本。只有上报版本等于任务目标版本时，任务才转为 success。超时扫描从首次进度上报开始计时，避免单纯“消息已发送”被误判为升级成功。

### 7. 查询接口增加批次范围

批次详情、设备部署列表和统计接口统一接收 `batch_id`，旧的按 OTA 包聚合接口保留兼容读取，但新前端页面必须使用批次维度，避免历史批次混入当前统计。

## Risks / Trade-offs

- [设备协议尚未实现] → 先冻结 JSON schema、topic、错误码和幂等规则，再联调设备；旧设备上报缺少 batch_id 时只能进入兼容/无法关联状态。
- [数据库已提交但 Kafka 发布失败] → 保留设备任务为待发送/失败，记录错误，并通过后续显式重试接口补发；若需要自动最终一致，再引入 outbox。
- [离线设备长期积压] → 为任务设置有效期和最大重试次数，设备上线时只补发未过期任务。
- [大规模批次冲击 broker] → 通过每分钟推送速率、批量领取和 backoff 控制发送压力，并监控 Kafka lag 与 MQTT 发布失败率。
- [升级包 URL 过期] → 发布前生成有效 URL，发送失败重试时重新生成 URL；设备端过期后允许主动拉取最新升级信息。
- [旧数据没有 batch_id] → 迁移时保留旧记录可读性，将空 batch_id 视为历史记录，不纳入新调度和新批次统计。

## Migration Plan

1. 增加任务关联、目标固件、重试/超时字段及索引，不删除旧记录。
2. 先部署只读兼容版本，继续支持旧 report 请求，同时上线批次维度查询。
3. 部署 Kafka producer 和 OTA transport worker，使用测试设备验证在线、离线和重复消息。
4. 设备端上线带 batch_id 的 progress/inform 协议后，切换新 report 消费路径。
5. 停止“定时任务直接 pending → in_progress”的旧逻辑，保留超时扫描和重试调度。
6. 观察一段完整发布周期后，再移除旧 report 的模糊更新兼容逻辑。

回滚时停止新 worker 和新批次创建入口；不回滚已写入的设备任务历史，避免破坏审计记录。

## Resolved Constraints

- 当前没有现成 MQTT Broker，首期部署 EMQX；EMQX 的认证、授权和 TLS 配置纳入部署任务。
- 设备 progress/inform 上报携带 `batch_id`，服务端按批次和设备关联，不再依赖模糊的设备最新记录。
- 首期只支持默认 OTA 模块，暂不实现多模块并行升级。
