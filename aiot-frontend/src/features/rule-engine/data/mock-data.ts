import type { Rule, RuleExecution, RuleStatistics } from './schema'

// ==================== Mock 规则数据 ====================

export const mockRules: Rule[] = [
  {
    id: 'rule-001',
    name: '温度过高自动开启空调',
    description: '当温度传感器检测到温度超过30度时，自动开启空调',
    type: 'device_linkage',
    status: 'enabled',
    trigger: {
      type: 'device_data',
      productIds: ['product-001'],
      topic: 'device/+/data',
    },
    condition: {
      logic: 'AND',
      conditions: [
        {
          field: 'temperature',
          operator: 'gt',
          value: 30,
          valueType: 'number',
        },
      ],
    },
    actions: [
      {
        type: 'device_control',
        params: {
          targetDeviceId: 'device-ac-001',
          command: 'setTemperature',
          params: {
            temperature: 26,
            mode: 'cooling',
          },
        },
      },
      {
        type: 'webhook',
        params: {
          url: 'https://api.example.com/webhook',
          method: 'POST',
          bodyTemplate:
            '{"event":"high_temperature","value":${temperature},"deviceId":"${deviceId}"}',
        },
      },
    ],
    executionCount: 142,
    successCount: 138,
    failureCount: 4,
    lastExecutedAt: '2026-01-06T10:30:00Z',
    lastExecutionStatus: 'success',
    createdBy: 'admin',
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-06T10:30:00Z',
    priority: 10,
    tags: ['温控', '自动化'],
  },
  {
    id: 'rule-002',
    name: '设备离线告警',
    description: '当设备离线超过5分钟时发送告警邮件',
    type: 'alert',
    status: 'enabled',
    trigger: {
      type: 'device_offline',
    },
    condition: {
      logic: 'AND',
      conditions: [
        {
          field: 'offlineDuration',
          operator: 'gte',
          value: 300,
          valueType: 'number',
        },
      ],
    },
    actions: [
      {
        type: 'email',
        params: {
          to: ['admin@example.com', 'ops@example.com'],
          subject: '设备离线告警',
          body: '<p>设备 ${deviceName} 已离线超过5分钟</p><p>设备ID: ${deviceId}</p>',
        },
      },
    ],
    executionCount: 23,
    successCount: 23,
    failureCount: 0,
    lastExecutedAt: '2026-01-05T14:20:00Z',
    lastExecutionStatus: 'success',
    createdBy: 'admin',
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-05T14:20:00Z',
    priority: 20,
    tags: ['告警', '监控'],
  },
  {
    id: 'rule-003',
    name: '设备数据转发到外部系统',
    description: '将所有传感器数据转发到外部数据分析平台',
    type: 'data_forwarding',
    status: 'enabled',
    trigger: {
      type: 'device_data',
      productIds: ['product-001', 'product-002'],
    },
    actions: [
      {
        type: 'http_request',
        params: {
          url: 'https://analytics.example.com/api/ingest',
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: 'Bearer token123',
          },
          body: '{"deviceId":"${deviceId}","data":${data},"timestamp":"${timestamp}"}',
          timeout: 5000,
        },
      },
    ],
    executionCount: 15234,
    successCount: 15198,
    failureCount: 36,
    lastExecutedAt: '2026-01-06T11:45:00Z',
    lastExecutionStatus: 'success',
    createdBy: 'admin',
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-06T11:45:00Z',
    priority: 5,
    tags: ['数据转发', '集成'],
  },
  {
    id: 'rule-004',
    name: '每小时温度统计',
    description: '使用SQL统计每小时的平均温度',
    type: 'sql',
    status: 'enabled',
    trigger: {
      type: 'device_data',
      topic: 'device/+/temperature',
    },
    sqlConfig: {
      sql: 'SELECT AVG(temperature) as avg_temp, deviceId, HOUR(timestamp) as hour FROM device_data WHERE productId = "product-001" GROUP BY deviceId, HOUR(timestamp)',
      dataSource: 'device_stream',
      outputTopic: 'analytics/temperature/hourly',
    },
    executionCount: 720,
    successCount: 720,
    failureCount: 0,
    lastExecutedAt: '2026-01-06T11:00:00Z',
    lastExecutionStatus: 'success',
    createdBy: 'admin',
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-06T11:00:00Z',
    priority: 0,
    tags: ['SQL', '统计'],
  },
  {
    id: 'rule-005',
    name: '复杂条件测试',
    description: '测试嵌套条件和多动作执行',
    type: 'device_linkage',
    status: 'disabled',
    trigger: {
      type: 'device_data',
      productIds: ['product-003'],
    },
    condition: {
      logic: 'OR',
      conditions: [
        {
          logic: 'AND',
          conditions: [
            {
              field: 'temperature',
              operator: 'gt',
              value: 35,
              valueType: 'number',
            },
            {
              field: 'humidity',
              operator: 'lt',
              value: 30,
              valueType: 'number',
            },
          ],
        },
        {
          field: 'alertLevel',
          operator: 'eq',
          value: 'critical',
          valueType: 'string',
        },
      ],
    },
    actions: [
      {
        type: 'device_control',
        params: {
          targetDeviceId: 'device-fan-001',
          command: 'turnOn',
          params: { speed: 'high' },
        },
      },
      {
        type: 'sms',
        params: {
          phoneNumbers: ['+86 138****1234'],
          message: '告警：环境条件异常',
        },
      },
    ],
    executionCount: 0,
    successCount: 0,
    failureCount: 0,
    lastExecutionStatus: 'pending',
    createdBy: 'admin',
    createdAt: '2026-01-04T00:00:00Z',
    updatedAt: '2026-01-04T00:00:00Z',
    priority: 15,
    tags: ['测试'],
  },
  {
    id: 'rule-006',
    name: '夜间定时关闭照明',
    description: '每晚23:00自动关闭所有照明设备',
    type: 'device_linkage',
    status: 'draft',
    trigger: {
      type: 'schedule',
      schedule: '0 23 * * *', // cron: 每天23:00
    },
    actions: [
      {
        type: 'device_control',
        params: {
          targetDeviceId: 'device-light-001',
          command: 'turnOff',
          params: {},
        },
      },
    ],
    executionCount: 0,
    successCount: 0,
    failureCount: 0,
    lastExecutionStatus: 'pending',
    createdBy: 'admin',
    createdAt: '2026-01-06T09:00:00Z',
    updatedAt: '2026-01-06T09:00:00Z',
    priority: 0,
    tags: ['定时', '节能'],
  },
]

// ==================== Mock 执行历史数据 ====================

export const mockExecutions: RuleExecution[] = [
  {
    id: 'exec-001',
    ruleId: 'rule-001',
    ruleName: '温度过高自动开启空调',
    status: 'success',
    triggeredAt: '2026-01-06T10:30:00Z',
    triggerData: {
      deviceId: 'device-sensor-001',
      temperature: 32.5,
      humidity: 65,
      timestamp: '2026-01-06T10:30:00Z',
    },
    conditionResult: true,
    actionsExecuted: [
      {
        actionType: 'device_control',
        status: 'success',
        startTime: '2026-01-06T10:30:00.100Z',
        endTime: '2026-01-06T10:30:00.250Z',
        duration: 150,
        result: { success: true, message: 'Command sent successfully' },
      },
      {
        actionType: 'webhook',
        status: 'success',
        startTime: '2026-01-06T10:30:00.260Z',
        endTime: '2026-01-06T10:30:00.450Z',
        duration: 190,
        result: { statusCode: 200 },
      },
    ],
    duration: 450,
    createdAt: '2026-01-06T10:30:00Z',
  },
  {
    id: 'exec-002',
    ruleId: 'rule-001',
    ruleName: '温度过高自动开启空调',
    status: 'failure',
    triggeredAt: '2026-01-06T09:15:00Z',
    triggerData: {
      deviceId: 'device-sensor-001',
      temperature: 31.2,
      humidity: 68,
      timestamp: '2026-01-06T09:15:00Z',
    },
    conditionResult: true,
    actionsExecuted: [
      {
        actionType: 'device_control',
        status: 'failure',
        startTime: '2026-01-06T09:15:00.100Z',
        endTime: '2026-01-06T09:15:05.100Z',
        duration: 5000,
        error: 'Timeout: Device not responding',
      },
    ],
    error: 'Action execution failed: device_control',
    duration: 5100,
    createdAt: '2026-01-06T09:15:00Z',
  },
  {
    id: 'exec-003',
    ruleId: 'rule-002',
    ruleName: '设备离线告警',
    status: 'success',
    triggeredAt: '2026-01-05T14:20:00Z',
    triggerData: {
      deviceId: 'device-sensor-003',
      deviceName: '温度传感器-003',
      offlineDuration: 320,
      lastOnlineTime: '2026-01-05T14:14:40Z',
    },
    conditionResult: true,
    actionsExecuted: [
      {
        actionType: 'email',
        status: 'success',
        startTime: '2026-01-05T14:20:00.100Z',
        endTime: '2026-01-05T14:20:01.200Z',
        duration: 1100,
        result: { messageId: 'msg-12345' },
      },
    ],
    duration: 1200,
    createdAt: '2026-01-05T14:20:00Z',
  },
  {
    id: 'exec-004',
    ruleId: 'rule-003',
    ruleName: '设备数据转发到外部系统',
    status: 'success',
    triggeredAt: '2026-01-06T11:45:00Z',
    triggerData: {
      deviceId: 'device-sensor-002',
      data: { temperature: 25.3, humidity: 55 },
      timestamp: '2026-01-06T11:45:00Z',
    },
    conditionResult: undefined,
    actionsExecuted: [
      {
        actionType: 'http_request',
        status: 'success',
        startTime: '2026-01-06T11:45:00.050Z',
        endTime: '2026-01-06T11:45:00.180Z',
        duration: 130,
        result: { statusCode: 200, body: { received: true } },
      },
    ],
    duration: 180,
    createdAt: '2026-01-06T11:45:00Z',
  },
  {
    id: 'exec-005',
    ruleId: 'rule-004',
    ruleName: '每小时温度统计',
    status: 'success',
    triggeredAt: '2026-01-06T11:00:00Z',
    triggerData: {
      aggregatedRecords: 156,
      timeWindow: '2026-01-06T10:00:00Z - 2026-01-06T11:00:00Z',
    },
    conditionResult: undefined,
    actionsExecuted: [],
    duration: 2340,
    createdAt: '2026-01-06T11:00:00Z',
  },
]

// ==================== Mock 统计数据 ====================

export const mockStatistics: RuleStatistics = {
  totalRules: 6,
  enabledRules: 4,
  disabledRules: 1,
  draftRules: 1,

  totalExecutions: 16119,
  successExecutions: 16078,
  failureExecutions: 41,

  avgExecutionTime: 245, // 毫秒

  rulesByType: {
    device_linkage: 3,
    data_forwarding: 1,
    alert: 1,
    sql: 1,
  },

  executionTrend: [
    {
      date: '2026-01-01',
      total: 2150,
      success: 2145,
      failure: 5,
    },
    {
      date: '2026-01-02',
      total: 2380,
      success: 2372,
      failure: 8,
    },
    {
      date: '2026-01-03',
      total: 2290,
      success: 2283,
      failure: 7,
    },
    {
      date: '2026-01-04',
      total: 2450,
      success: 2440,
      failure: 10,
    },
    {
      date: '2026-01-05',
      total: 3120,
      success: 3109,
      failure: 11,
    },
    {
      date: '2026-01-06',
      total: 3729,
      success: 3729,
      failure: 0,
    },
  ],
}

// ==================== Mock 可用字段 ====================

export const mockAvailableFields = [
  {
    field: 'temperature',
    type: 'number' as const,
    label: '温度',
    description: '设备上报的温度值（摄氏度）',
  },
  {
    field: 'humidity',
    type: 'number' as const,
    label: '湿度',
    description: '设备上报的湿度值（百分比）',
  },
  {
    field: 'pressure',
    type: 'number' as const,
    label: '气压',
    description: '气压值（hPa）',
  },
  {
    field: 'status',
    type: 'string' as const,
    label: '设备状态',
    description: '设备的当前状态',
  },
  {
    field: 'online',
    type: 'boolean' as const,
    label: '在线状态',
    description: '设备是否在线',
  },
  {
    field: 'battery',
    type: 'number' as const,
    label: '电池电量',
    description: '电池电量百分比',
  },
  {
    field: 'signalStrength',
    type: 'number' as const,
    label: '信号强度',
    description: '网络信号强度（dBm）',
  },
  {
    field: 'alertLevel',
    type: 'string' as const,
    label: '告警等级',
    description: '告警级别：info, warning, critical',
  },
]
