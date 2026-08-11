# 规则引擎模块

IoT 平台的规则引擎功能模块，支持设备联动、数据流转、告警通知等场景。

## 📁 目录结构

```
src/features/rule-engine/
├── api/
│   └── queries.ts                    # React Query API hooks (使用Mock数据)
├── components/
│   ├── dialogs/
│   │   ├── create-rule-dialog.tsx    # 创建规则对话框
│   │   ├── delete-rule-dialog.tsx    # 删除确认对话框
│   │   └── view-rule-details-dialog.tsx  # 查看详情对话框
│   ├── rule-form/
│   │   ├── action-config.tsx         # 动作配置组件
│   │   ├── condition-builder.tsx     # 条件构建器（核心组件）
│   │   ├── rule-basic-info.tsx       # 基本信息表单
│   │   └── trigger-config.tsx        # 触发器配置
│   ├── rules-columns.tsx             # 表格列定义
│   └── rules-table.tsx               # 规则列表表格
├── data/
│   ├── mock-data.ts                  # Mock数据
│   └── schema.ts                     # Zod数据模型定义
└── index.tsx                         # 主页面

路由:
src/routes/_authenticated/rule-engine/index.tsx
```

## 🚀 功能特性

### 1. 规则类型

- ✅ **设备联动** - 设备之间的自动化控制
- ✅ **数据流转** - 数据转发到外部系统
- ✅ **告警规则** - 条件触发告警通知
- ⏳ **SQL规则** - SQL数据处理（框架已就绪）

### 2. 触发器类型

- 设备数据上报
- 设备状态变化
- 设备上线/离线
- 定时触发（Cron）
- 手动触发

### 3. 条件构建器

- 支持多条件组合（AND/OR）
- 支持嵌套条件组
- 10种运算符（=, !=, >, >=, <, <=, in, nin, contains, between）
- 可视化条件编辑

### 4. 动作类型

- ✅ **设备控制** - 发送指令到目标设备
- ✅ **HTTP请求** - 调用外部API
- ✅ **Webhook** - 推送到Webhook地址
- ✅ **邮件通知** - 发送告警邮件
- ⏳ 短信通知（框架已就绪）
- ⏳ Kafka消息（框架已就绪）

### 5. 规则管理

- 规则列表查看（支持搜索、筛选、排序）
- 创建/编辑/删除规则
- 启用/禁用规则
- 手动触发规则
- 批量操作
- 查看规则详情

### 6. 执行监控

- 执行次数统计
- 成功率统计
- 最后执行时间
- 执行历史（Mock数据中已包含）

## 🎯 快速开始

### 1. 访问规则引擎页面

启动项目后，访问：

```
http://localhost:5173/rule-engine
```

### 2. 创建第一条规则

1. 点击 **"创建规则"** 按钮
2. 填写基本信息：
   - 规则名称：例如 "温度过高自动开空调"
   - 规则类型：选择 "设备联动"
   - 规则状态：选择 "草稿" 或 "启用"
3. 配置触发器：
   - 触发类型：选择 "设备数据上报"
   - 产品ID：输入目标产品ID（可选）
4. 配置条件：
   - 点击 "添加条件"
   - 字段：选择 "temperature"（温度）
   - 运算符：选择 "大于 (>)"
   - 值：输入 "30"
5. 配置动作：
   - 点击 "设备控制"
   - 目标设备ID：输入空调设备ID
   - 指令名称：输入 "setTemperature"
   - 指令参数：输入 `{"temperature": 26}`
6. 点击 "创建规则"

### 3. 测试规则

- 在规则列表中，点击操作菜单 (⋮)
- 选择 "手动触发"
- 系统会模拟执行规则并显示结果

## 💡 使用示例

### 示例1：温度告警

```typescript
// 规则配置
{
  name: "温度过高告警",
  type: "alert",
  trigger: {
    type: "device_data",
    productIds: ["sensor-product-001"]
  },
  condition: {
    logic: "AND",
    conditions: [
      { field: "temperature", operator: "gt", value: 35 }
    ]
  },
  actions: [
    {
      type: "email",
      params: {
        to: ["admin@example.com"],
        subject: "温度告警",
        body: "设备 ${deviceId} 温度达到 ${temperature}℃"
      }
    }
  ]
}
```

### 示例2：复杂条件联动

```typescript
// 嵌套条件：(温度>35 AND 湿度<30) OR 告警等级=critical
{
  condition: {
    logic: "OR",
    conditions: [
      {
        logic: "AND",
        conditions: [
          { field: "temperature", operator: "gt", value: 35 },
          { field: "humidity", operator: "lt", value: 30 }
        ]
      },
      { field: "alertLevel", operator: "eq", value: "critical" }
    ]
  }
}
```

### 示例3：数据转发

```typescript
{
  name: "数据转发到分析平台",
  type: "data_forwarding",
  trigger: {
    type: "device_data"
  },
  actions: [
    {
      type: "http_request",
      params: {
        url: "https://analytics.example.com/api/ingest",
        method: "POST",
        body: '{"deviceId":"${deviceId}","data":${data},"timestamp":"${timestamp}"}'
      }
    }
  ]
}
```

## 🔧 开发说明

### Mock 数据模式

当前所有 API 请求都使用 Mock 数据，不会发送实际的网络请求：

- 数据存储在内存中（`rulesData` 变量）
- 自动模拟 API 延迟（300ms）
- 支持所有 CRUD 操作
- 数据在页面刷新后重置

### 切换到真实 API

修改 `src/features/rule-engine/api/queries.ts`：

```typescript
// 替换 Mock 实现为真实 API 调用
export function useRules(params: RulesQueryParams = {}) {
  return useQuery({
    queryKey: ruleKeys.list(params),
    queryFn: async () => {
      // 使用真实的 API 客户端
      const { data } = await ruleEngineClient.get('/v1/rules', { params })
      return data
    },
  })
}
```

### 自定义可用字段

修改 `src/features/rule-engine/data/mock-data.ts` 中的 `mockAvailableFields`：

```typescript
export const mockAvailableFields = [
  {
    field: 'customField',
    type: 'number',
    label: '自定义字段',
    description: '这是一个自定义字段',
  },
  // ... 添加更多字段
]
```

## 📊 数据模型

### Rule（规则）

```typescript
interface Rule {
  id: string
  name: string
  description?: string
  type: 'device_linkage' | 'data_forwarding' | 'alert' | 'sql'
  status: 'enabled' | 'disabled' | 'draft'
  trigger: TriggerConfig
  condition?: ConditionGroup
  actions?: ActionConfig[]
  sqlConfig?: SqlRuleConfig
  executionCount: number
  successCount: number
  failureCount: number
  lastExecutedAt?: string
  lastExecutionStatus?: 'success' | 'failure' | 'pending'
  createdBy: string
  createdAt: string
  updatedAt: string
  priority: number
  tags?: string[]
}
```

### ConditionGroup（条件组）

```typescript
interface ConditionGroup {
  logic: 'AND' | 'OR'
  conditions: Array<Condition | ConditionGroup> // 支持嵌套
}

interface Condition {
  field: string
  operator:
    | 'eq'
    | 'ne'
    | 'gt'
    | 'gte'
    | 'lt'
    | 'lte'
    | 'in'
    | 'nin'
    | 'contains'
    | 'between'
  value: any
  valueType?: 'number' | 'string' | 'boolean' | 'json'
}
```

## 🎨 UI组件

### 核心组件

1. **ConditionBuilder** - 可视化条件构建器
   - 支持拖拽排序（未实现）
   - 嵌套条件组
   - 实时预览

2. **ActionConfig** - 动作配置器
   - 多动作支持
   - 模板变量支持
   - 折叠/展开

3. **RulesTable** - 规则列表表格
   - 列排序
   - 全局搜索
   - 批量选择
   - 列显示控制

## 🚧 待实现功能

### 高优先级

- [ ] 编辑规则功能
- [ ] SQL规则编辑器
- [ ] 规则执行历史页面
- [ ] 规则统计分析页面

### 中优先级

- [ ] 规则导入/导出
- [ ] 规则复制功能
- [ ] 规则测试调试工具
- [ ] 规则模板

### 低优先级

- [ ] 规则版本管理
- [ ] 规则审批流程
- [ ] 规则调度策略
- [ ] 性能优化（虚拟滚动）

## 📝 注意事项

1. **数据验证**：所有表单都使用 Zod 进行运行时验证
2. **类型安全**：完整的 TypeScript 类型定义
3. **错误处理**：使用 toast 显示操作结果
4. **性能**：React Query 自动缓存和优化请求
5. **Mock数据**：开发环境使用Mock数据，生产环境需切换到真实API

## 🔗 相关链接

- [Zod文档](https://zod.dev/)
- [React Query文档](https://tanstack.com/query/latest)
- [TanStack Table文档](https://tanstack.com/table/latest)
- [Shadcn UI文档](https://ui.shadcn.com/)

## 📞 支持

如有问题或建议，请联系开发团队。
