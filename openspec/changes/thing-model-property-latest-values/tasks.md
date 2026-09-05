## 1. 后端属性最后值读模型

- [x] 1.1 定义物模型属性最后值 API DTO（属性元数据、可空最后值和可空最后上报时间），并通过 Swagger 生成验证契约可见。
- [x] 1.2 在时序访问层实现按设备和属性标识符集合查询最新点的能力，并用单元测试验证每个属性返回时间最新的点。
- [x] 1.3 实现设备、产品 TSL 与时序最后点的服务端聚合，验证已上报、未上报、未定义键、设备不存在和无效/缺失 TSL 场景。
- [x] 1.4 注册 `GET /devices/:deviceKey/thing-model/properties` 路由与依赖注入，并用 Handler 测试验证成功响应及错误响应。

## 2. 契约生成与前端属性页

- [x] 2.1 生成 Swagger 文档和前端 API 客户端，验证生成客户端包含属性最后值接口与模型。
- [x] 2.2 将设备详情 `PropertyTab` 改为调用专用属性最后值接口，移除该页面对 Shadow 和原始最新 telemetry 接口的引用，并通过前端类型检查验证。
- [x] 2.3 按 TSL 顺序展示属性名称、标识符、类型、单位、读写方式、当前值和最后上报时间；无值时展示未上报状态，并维护中英文文案。
- [x] 2.4 让属性历史趋势选择器使用属性最后值接口返回的定义属性，验证历史请求仍只调用 `/v1/devices/:deviceKey/telemetry/history`。

## 3. 集成验证

- [x] 3.1 添加或更新后端服务与 Handler 测试，执行 `make -C backend test` 和 `make -C backend build`。
- [x] 3.2 执行 `pnpm -C frontend format` 与 `pnpm -C frontend build`，验证属性页不再请求 `/devices/:deviceKey/shadow` 或 `/devices/:deviceKey/telemetry`。

## 4. 物模型数据模块统一

- [x] 4.1 将属性最后值和服务调用记录统一到 `ThingModelDataHandler` 与 `ThingModelDataService`，并更新路由、依赖注入及测试。
