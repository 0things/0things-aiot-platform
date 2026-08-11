import { z } from 'zod'

// ==================== 枚举定义 ====================

// 规则类型
export const ruleTypeEnum = z.enum([
  'device_linkage', // 设备联动
  'data_forwarding', // 数据流转
  'alert', // 告警规则
  'sql', // SQL规则
])
export type RuleType = z.infer<typeof ruleTypeEnum>

// 规则状态
export const ruleStatusEnum = z.enum([
  'enabled', // 启用
  'disabled', // 禁用
  'draft', // 草稿
])
export type RuleStatus = z.infer<typeof ruleStatusEnum>

// 触发器类型
export const triggerTypeEnum = z.enum([
  'device_data', // 设备数据上报
  'device_status', // 设备状态变化
  'device_online', // 设备上线
  'device_offline', // 设备离线
  'schedule', // 定时触发
  'manual', // 手动触发
])
export type TriggerType = z.infer<typeof triggerTypeEnum>

// 条件运算符
export const operatorEnum = z.enum([
  'eq', // 等于 =
  'ne', // 不等于 !=
  'gt', // 大于 >
  'gte', // 大于等于 >=
  'lt', // 小于 <
  'lte', // 小于等于 <=
  'in', // 包含于
  'nin', // 不包含于
  'contains', // 字符串包含
  'between', // 区间
])
export type Operator = z.infer<typeof operatorEnum>

// 逻辑运算符
export const logicOperatorEnum = z.enum(['AND', 'OR'])
export type LogicOperator = z.infer<typeof logicOperatorEnum>

// 动作类型
export const actionTypeEnum = z.enum([
  'device_control', // 设备控制
  'http_request', // HTTP请求
  'email', // 邮件通知
  'sms', // 短信通知
  'webhook', // Webhook
  'kafka', // Kafka消息
  'database', // 数据库写入
  'function', // 函数计算
])
export type ActionType = z.infer<typeof actionTypeEnum>

// HTTP方法
export const httpMethodEnum = z.enum(['GET', 'POST', 'PUT', 'DELETE', 'PATCH'])
export type HttpMethod = z.infer<typeof httpMethodEnum>

// ==================== 触发器定义 ====================

// 触发器配置
export const triggerConfigSchema = z.object({
  type: triggerTypeEnum,
  // 数据源过滤
  productIds: z.array(z.string()).optional(), // 产品ID列表
  deviceIds: z.array(z.string()).optional(), // 设备ID列表
  deviceTags: z.record(z.string(), z.string()).optional(), // 设备标签
  // 定时配置（cron表达式）
  schedule: z.string().optional(), // 如: "0 0 * * *" 每天0点
  // 数据topic（MQTT topic）
  topic: z.string().optional(),
})
export type TriggerConfig = z.infer<typeof triggerConfigSchema>

// ==================== 条件定义 ====================

// 单个条件
export const conditionSchema = z.object({
  field: z.string(), // 字段名，如 "temperature"
  operator: operatorEnum, // 运算符
  value: z.any(), // 比较值
  valueType: z.enum(['number', 'string', 'boolean', 'json']).optional(),
})
export type Condition = z.infer<typeof conditionSchema>

// 条件组（支持嵌套）
export const conditionGroupSchema: z.ZodType<any> = z.lazy(() =>
  z.object({
    logic: logicOperatorEnum, // AND/OR
    conditions: z.array(z.union([conditionSchema, conditionGroupSchema])),
  })
)
export type ConditionGroup = z.infer<typeof conditionGroupSchema>

// ==================== 动作定义 ====================

// 设备控制动作参数
export const deviceControlActionSchema = z.object({
  targetDeviceId: z.string(),
  command: z.string(), // 指令名称
  params: z.record(z.string(), z.any()), // 指令参数
})

// HTTP请求动作参数
export const httpRequestActionSchema = z.object({
  url: z.string().url(),
  method: httpMethodEnum,
  headers: z.record(z.string(), z.string()).optional(),
  body: z.string().optional(), // 支持模板变量，如 ${temperature}
  timeout: z.number().optional(),
})

// 邮件通知动作参数
export const emailActionSchema = z.object({
  to: z.array(z.string().email()),
  cc: z.array(z.string().email()).optional(),
  subject: z.string(),
  body: z.string(), // 支持HTML和模板变量
})

// Webhook动作参数
export const webhookActionSchema = z.object({
  url: z.string().url(),
  method: httpMethodEnum.default('POST'),
  headers: z.record(z.string(), z.string()).optional(),
  bodyTemplate: z.string().optional(),
})

// 动作配置（联合类型）
export const actionConfigSchema = z.union([
  z.object({
    type: z.literal('device_control'),
    params: deviceControlActionSchema,
  }),
  z.object({
    type: z.literal('http_request'),
    params: httpRequestActionSchema,
  }),
  z.object({ type: z.literal('email'), params: emailActionSchema }),
  z.object({ type: z.literal('webhook'), params: webhookActionSchema }),
  z.object({ type: z.literal('sms'), params: z.record(z.string(), z.any()) }),
  z.object({ type: z.literal('kafka'), params: z.record(z.string(), z.any()) }),
  z.object({ type: z.literal('database'), params: z.record(z.string(), z.any()) }),
  z.object({ type: z.literal('function'), params: z.record(z.string(), z.any()) }),
])
export type ActionConfig = z.infer<typeof actionConfigSchema>

// ==================== SQL规则定义 ====================

export const sqlRuleConfigSchema = z.object({
  sql: z.string(), // SQL查询语句
  dataSource: z.string().optional(), // 数据源名称
  outputTopic: z.string().optional(), // 输出topic
})
export type SqlRuleConfig = z.infer<typeof sqlRuleConfigSchema>

// ==================== 完整规则Schema ====================

export const ruleSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string().optional(),
  type: ruleTypeEnum,
  status: ruleStatusEnum,

  // 触发器配置
  trigger: triggerConfigSchema,

  // 条件配置（可选，SQL规则不需要）
  condition: conditionGroupSchema.optional(),

  // 动作配置（可以有多个动作）
  actions: z.array(actionConfigSchema).optional(),

  // SQL配置（仅SQL规则需要）
  sqlConfig: sqlRuleConfigSchema.optional(),

  // 执行统计
  executionCount: z.number().default(0), // 总执行次数
  successCount: z.number().default(0), // 成功次数
  failureCount: z.number().default(0), // 失败次数
  lastExecutedAt: z.string().optional(), // 最后执行时间
  lastExecutionStatus: z.enum(['success', 'failure', 'pending']).optional(),

  // 元数据
  createdBy: z.string(),
  createdAt: z.string(),
  updatedAt: z.string(),
  priority: z.number().default(0), // 优先级，数字越大优先级越高
  tags: z.array(z.string()).optional(),
})

export type Rule = z.infer<typeof ruleSchema>

// ==================== 表单Schema ====================

export const createRuleFormSchema = z
  .object({
    name: z.string().min(1, '规则名称不能为空').max(100, '规则名称过长'),
    description: z.string().max(500).optional(),
    type: ruleTypeEnum,
    status: ruleStatusEnum.default('draft'),
    trigger: triggerConfigSchema,
    condition: conditionGroupSchema.optional(),
    actions: z.array(actionConfigSchema).min(1, '至少需要一个动作').optional(),
    sqlConfig: sqlRuleConfigSchema.optional(),
    priority: z.number().min(0).max(100).default(0),
    tags: z.array(z.string()).optional(),
  })
  .refine(
    (data) => {
      // SQL规则必须有sqlConfig
      if (data.type === 'sql') {
        return !!data.sqlConfig
      }
      // 其他规则必须有actions
      return !!data.actions && data.actions.length > 0
    },
    {
      message: '请配置规则动作或SQL',
    }
  )

export type CreateRuleFormData = z.infer<typeof createRuleFormSchema>

// ==================== 执行历史Schema ====================

export const ruleExecutionSchema = z.object({
  id: z.string(),
  ruleId: z.string(),
  ruleName: z.string(),
  status: z.enum(['success', 'failure', 'timeout', 'skipped']),

  // 触发信息
  triggeredAt: z.string(),
  triggerData: z.record(z.string(), z.any()), // 触发时的原始数据

  // 执行信息
  conditionResult: z.boolean().optional(), // 条件判断结果
  actionsExecuted: z
    .array(
      z.object({
        actionType: actionTypeEnum,
        status: z.enum(['success', 'failure']),
        startTime: z.string(),
        endTime: z.string(),
        duration: z.number(), // 毫秒
        error: z.string().optional(),
        result: z.any().optional(),
      })
    )
    .optional(),

  // 错误信息
  error: z.string().optional(),
  errorStack: z.string().optional(),

  // 性能指标
  duration: z.number(), // 总执行时长（毫秒）

  createdAt: z.string(),
})

export type RuleExecution = z.infer<typeof ruleExecutionSchema>

// ==================== 统计Schema ====================

export const ruleStatisticsSchema = z.object({
  totalRules: z.number(),
  enabledRules: z.number(),
  disabledRules: z.number(),
  draftRules: z.number(),

  totalExecutions: z.number(),
  successExecutions: z.number(),
  failureExecutions: z.number(),

  avgExecutionTime: z.number(), // 平均执行时间（毫秒）

  rulesByType: z.record(z.string(), z.number()), // 按类型分组统计

  executionTrend: z.array(
    z.object({
      date: z.string(),
      total: z.number(),
      success: z.number(),
      failure: z.number(),
    })
  ),
})

export type RuleStatistics = z.infer<typeof ruleStatisticsSchema>

// ==================== 辅助类型 ====================

// 可用字段定义（用于条件构建器）
export interface AvailableField {
  field: string
  type: 'number' | 'string' | 'boolean' | 'json'
  label: string
  description?: string
}

// 规则类型标签映射
export const ruleTypeLabels: Record<RuleType, string> = {
  device_linkage: '设备联动',
  data_forwarding: '数据流转',
  alert: '告警规则',
  sql: 'SQL规则',
}

// 规则状态标签映射
export const ruleStatusLabels: Record<RuleStatus, string> = {
  enabled: '已启用',
  disabled: '已禁用',
  draft: '草稿',
}

// 触发器类型标签映射
export const triggerTypeLabels: Record<TriggerType, string> = {
  device_data: '设备数据上报',
  device_status: '设备状态变化',
  device_online: '设备上线',
  device_offline: '设备离线',
  schedule: '定时触发',
  manual: '手动触发',
}

// 运算符标签映射
export const operatorLabels: Record<Operator, string> = {
  eq: '等于 (=)',
  ne: '不等于 (!=)',
  gt: '大于 (>)',
  gte: '大于等于 (>=)',
  lt: '小于 (<)',
  lte: '小于等于 (<=)',
  in: '包含于',
  nin: '不包含于',
  contains: '包含',
  between: '区间',
}

// 动作类型标签映射
export const actionTypeLabels: Record<ActionType, string> = {
  device_control: '设备控制',
  http_request: 'HTTP请求',
  email: '邮件通知',
  sms: '短信通知',
  webhook: 'Webhook',
  kafka: 'Kafka消息',
  database: '数据库写入',
  function: '函数计算',
}
