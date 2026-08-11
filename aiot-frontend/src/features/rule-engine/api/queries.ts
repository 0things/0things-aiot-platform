import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { axiosInstance } from '@/api/clients'
import { DEVICE_SERVICE_BASE_URL } from '@/api/config'
import type {
  AvailableField,
  CreateRuleFormData,
  Rule,
  RuleExecution,
  RuleStatistics,
  RuleStatus,
} from '../data/schema'

export const ruleKeys = {
  all: ['rules'] as const,
  lists: () => [...ruleKeys.all, 'list'] as const,
  list: (filters: any) => [...ruleKeys.lists(), filters] as const,
  details: () => [...ruleKeys.all, 'detail'] as const,
  detail: (id: string) => [...ruleKeys.details(), id] as const,
  statistics: () => [...ruleKeys.all, 'statistics'] as const,
  executions: (ruleId?: string) =>
    [...ruleKeys.all, 'executions', ruleId] as const,
  availableFields: () => [...ruleKeys.all, 'availableFields'] as const,
}

interface RulesQueryParams {
  page?: number
  pageSize?: number
  type?: string
  status?: RuleStatus
  search?: string
  productId?: string
}

type RuleApi = {
  id: number
  name: string
  description?: string
  type: string
  status: string
  productId?: number
  priority?: number
  triggerConfig?: string
  conditionConfig?: string
  actionConfig?: string
  sqlConfig?: string
  executionCount?: number
  successCount?: number
  failureCount?: number
  lastExecutionStatus?: string
  createdBy?: string
  tags?: string[]
  lastExecutedAt?: string
  createdAt?: string
  updatedAt?: string
}

type RuleExecutionApi = {
  id: number
  ruleId: number
  ruleName: string
  status: string
  triggeredAt?: string
  triggerData?: string
  conditionResult?: boolean
  duration?: number
  error?: string
  createdAt?: string
}

function parseJSON<T>(raw?: string, fallback?: T): T {
  if (!raw) return fallback as T
  try {
    return JSON.parse(raw) as T
  } catch {
    return fallback as T
  }
}

function mapRule(item: RuleApi): Rule {
  return {
    id: String(item.id),
    name: item.name,
    description: item.description,
    type: item.type as Rule['type'],
    status: item.status as Rule['status'],
    trigger: parseJSON(item.triggerConfig, { type: 'device_data' }),
    condition: parseJSON(item.conditionConfig, undefined),
    actions: parseJSON(item.actionConfig, []),
    sqlConfig: parseJSON(item.sqlConfig, undefined),
    executionCount: item.executionCount || 0,
    successCount: item.successCount || 0,
    failureCount: item.failureCount || 0,
    lastExecutedAt: item.lastExecutedAt,
    lastExecutionStatus:
      (item.lastExecutionStatus as Rule['lastExecutionStatus']) || 'pending',
    createdBy: item.createdBy || 'current-user',
    createdAt: item.createdAt || new Date().toISOString(),
    updatedAt: item.updatedAt || new Date().toISOString(),
    priority: item.priority || 0,
    tags: item.tags || [],
  }
}

function mapExecution(item: RuleExecutionApi): RuleExecution {
  return {
    id: String(item.id),
    ruleId: String(item.ruleId),
    ruleName: item.ruleName,
    status: item.status as RuleExecution['status'],
    triggeredAt: item.triggeredAt || item.createdAt || new Date().toISOString(),
    triggerData: parseJSON(item.triggerData, {}),
    conditionResult: item.conditionResult,
    duration: item.duration || 0,
    error: item.error,
    createdAt: item.createdAt || new Date().toISOString(),
  }
}

function buildRulePayload(data: Partial<CreateRuleFormData>) {
  const productId = data.trigger?.productIds?.[0]
  return {
    name: data.name,
    description: data.description,
    type: data.type,
    status: data.status,
    productId: productId ? Number(productId) : 0,
    priority: data.priority ?? 0,
    triggerConfig: JSON.stringify(data.trigger || {}),
    conditionConfig: JSON.stringify(data.condition || {}),
    actionConfig: JSON.stringify(data.actions || []),
    tags: data.tags || [],
    sqlConfig: JSON.stringify(data.sqlConfig || {}),
  }
}

export function useRules(params: RulesQueryParams = {}) {
  return useQuery({
    queryKey: ruleKeys.list(params),
    queryFn: async () => {
      const { data } = await axiosInstance.get(
        `${DEVICE_SERVICE_BASE_URL}/v1/rules`,
        {
          params: {
            page: params.page,
            pageSize: params.pageSize,
            type: params.type,
            status: params.status,
            search: params.search,
            productId: params.productId ? Number(params.productId) : undefined,
          },
        }
      )
      return {
        items: (data.items || []).map(mapRule),
        total: data.total || 0,
        page: data.page || params.page || 1,
        pageSize: data.pageSize || params.pageSize || 20,
      }
    },
  })
}

export function useRule(id: string) {
  return useQuery({
    queryKey: ruleKeys.detail(id),
    queryFn: async () => {
      const { data } = await axiosInstance.get(
        `${DEVICE_SERVICE_BASE_URL}/v1/rules/${id}`
      )
      return mapRule(data.rule)
    },
    enabled: !!id,
  })
}

export function useRuleStatistics() {
  return useQuery({
    queryKey: ruleKeys.statistics(),
    queryFn: async () => {
      const { data } = await axiosInstance.get(
        `${DEVICE_SERVICE_BASE_URL}/v1/rules`,
        {
          params: { page: 1, pageSize: 200 },
        }
      )
      const items = (data.items || []).map(mapRule)
      const stats: RuleStatistics = {
        totalRules: items.length,
        enabledRules: items.filter((r: Rule) => r.status === 'enabled').length,
        disabledRules: items.filter((r: Rule) => r.status === 'disabled').length,
        draftRules: items.filter((r: Rule) => r.status === 'draft').length,
        totalExecutions: items.reduce((sum: number, item: Rule) => sum + item.executionCount, 0),
        successExecutions: items.reduce((sum: number, item: Rule) => sum + item.successCount, 0),
        failureExecutions: items.reduce((sum: number, item: Rule) => sum + item.failureCount, 0),
        avgExecutionTime: 0,
        rulesByType: items.reduce((acc: Record<string, number>, item: Rule) => {
          acc[item.type] = (acc[item.type] || 0) + 1
          return acc
        }, {} as Record<string, number>),
        executionTrend: [],
      }
      return stats
    },
  })
}

export function useAvailableFields(productId?: string) {
  return useQuery({
    queryKey: [...ruleKeys.availableFields(), productId],
    queryFn: async () => {
      if (!productId) return [] as AvailableField[]
      const { data } = await axiosInstance.get(
        `${DEVICE_SERVICE_BASE_URL}/v1/rules/available-fields`,
        {
          params: { productId: Number(productId) },
        }
      )
      return (data.fields || []) as AvailableField[]
    },
    enabled: !!productId,
  })
}

export function useCreateRule() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (data: CreateRuleFormData) => {
      const response = await axiosInstance.post(
        `${DEVICE_SERVICE_BASE_URL}/v1/rules`,
        buildRulePayload(data)
      )
      return mapRule(response.data.rule)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ruleKeys.lists() })
      queryClient.invalidateQueries({ queryKey: ruleKeys.statistics() })
    },
  })
}

export function useUpdateRule() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({
      id,
      data,
    }: {
      id: string
      data: Partial<CreateRuleFormData>
    }) => {
      const response = await axiosInstance.put(
        `${DEVICE_SERVICE_BASE_URL}/v1/rules/${id}`,
        buildRulePayload(data)
      )
      return mapRule(response.data.rule)
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ruleKeys.detail(variables.id) })
      queryClient.invalidateQueries({ queryKey: ruleKeys.lists() })
    },
  })
}

export function useDeleteRule() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (id: string) => {
      await axiosInstance.delete(`${DEVICE_SERVICE_BASE_URL}/v1/rules/${id}`)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ruleKeys.lists() })
      queryClient.invalidateQueries({ queryKey: ruleKeys.statistics() })
    },
  })
}

export function useDeleteRules() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (ids: string[]) => {
      await Promise.all(
        ids.map((id) =>
          axiosInstance.delete(`${DEVICE_SERVICE_BASE_URL}/v1/rules/${id}`)
        )
      )
      return { deleted: ids.length }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ruleKeys.lists() })
      queryClient.invalidateQueries({ queryKey: ruleKeys.statistics() })
    },
  })
}

export function useUpdateRuleStatus() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, status }: { id: string; status: RuleStatus }) => {
      const url =
        status === 'enabled'
          ? `${DEVICE_SERVICE_BASE_URL}/v1/rules/${id}:enable`
          : `${DEVICE_SERVICE_BASE_URL}/v1/rules/${id}:disable`
      const response = await axiosInstance.post(url, {})
      return mapRule(response.data.rule)
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ruleKeys.detail(variables.id) })
      queryClient.invalidateQueries({ queryKey: ruleKeys.lists() })
      queryClient.invalidateQueries({ queryKey: ruleKeys.statistics() })
    },
  })
}

export function useTestRule() {
  return useMutation({
    mutationFn: async ({
      id,
      testData,
    }: {
      id: string
      testData: Record<string, any>
    }) => {
      const { data } = await axiosInstance.post(
        `${DEVICE_SERVICE_BASE_URL}/v1/rules/${id}:evaluate`,
        {
          inputPayload: JSON.stringify(testData),
        }
      )
      return {
        conditionPassed: data.execution?.conditionResult ?? true,
        conditionDetails: {
          testData,
          result: data.execution?.conditionResult ?? true,
        },
        estimatedActions: [],
      }
    },
  })
}

export function useTriggerRule() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({
      id,
      data,
    }: {
      id: string
      data?: Record<string, any>
    }) => {
      const response = await axiosInstance.post(
        `${DEVICE_SERVICE_BASE_URL}/v1/rules/${id}:evaluate`,
        {
          inputPayload: JSON.stringify(data || {}),
        }
      )
      return mapExecution(response.data.execution)
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ruleKeys.executions(variables.id),
      })
      queryClient.invalidateQueries({ queryKey: ruleKeys.detail(variables.id) })
      queryClient.invalidateQueries({ queryKey: ruleKeys.lists() })
    },
  })
}

interface ExecutionsQueryParams {
  page?: number
  pageSize?: number
}

export function useRuleExecutions(
  ruleId?: string,
  params: ExecutionsQueryParams = {}
) {
  return useQuery({
    queryKey: [...ruleKeys.executions(ruleId), params],
    queryFn: async () => {
      if (!ruleId) {
        return { items: [], total: 0 }
      }
      const { data } = await axiosInstance.get(
        `${DEVICE_SERVICE_BASE_URL}/v1/rules/${ruleId}/executions`,
        {
          params: {
            page: params.page,
            pageSize: params.pageSize,
          },
        }
      )
      return {
        items: (data.items || []).map(mapExecution),
        total: data.total || 0,
      }
    },
    enabled: !!ruleId,
  })
}

export function useRuleExecution(executionId: string) {
  return useQuery({
    queryKey: [...ruleKeys.executions(), executionId],
    queryFn: async () => {
      throw new Error('Execution detail endpoint is not implemented')
    },
    enabled: false && !!executionId,
  })
}
