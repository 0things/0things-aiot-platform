import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { orvalAxios } from '@/api/orval-mutator'

export type AlertSeverity = 'info' | 'warning' | 'critical'
export type AlertStatus = 'open' | 'acknowledged' | 'resolved' | 'snoozed'

export type Alert = {
  id: string
  organizationId: string
  ruleId: string
  ruleName?: string
  deviceKey: string
  severity: AlertSeverity
  status: AlertStatus
  summary: string
  payload?: Record<string, unknown>
  count: number
  raisedAt: string
  lastRaisedAt: string
  ackAt?: string | null
  resolvedAt?: string | null
}

export const alertKeys = {
  all: ['alerts'] as const,
  list: (filters: Record<string, unknown>) =>
    [...alertKeys.all, 'list', filters] as const,
  detail: (id: string) => [...alertKeys.all, 'detail', id] as const,
  openCount: () => [...alertKeys.all, 'open-count'] as const,
}

export function useAlerts(filters: {
  severity?: AlertSeverity
  status?: AlertStatus
  deviceKey?: string
  ruleId?: string
  start?: string
  end?: string
  page?: number
  pageSize?: number
}) {
  return useQuery({
    queryKey: alertKeys.list(filters),
    queryFn: async () => {
      const res = await orvalAxios<{
        data: { alerts: Alert[]; total: number }
      }>({
        url: '/alerts',
        method: 'GET',
        params: {
          severity: filters.severity,
          status: filters.status,
          device_key: filters.deviceKey,
          page: filters.page,
          pageSize: filters.pageSize,
        },
      })
      return (res?.data ?? res) as unknown as { alerts: Alert[]; total: number }
    },
  })
}

export function useOpenAlertCount() {
  return useQuery({
    queryKey: alertKeys.openCount(),
    queryFn: async () => {
      const data = await orvalAxios<{ data: { total: number } }>({
        url: '/alerts',
        method: 'GET',
        params: { status: 'open', pageSize: 0 },
      })
      return data?.data?.total ?? 0
    },
    refetchInterval: 15_000,
  })
}

export function useAckAlert() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (id: string) => {
      return orvalAxios({ url: `/alerts/${id}/ack`, method: 'POST' })
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: alertKeys.all }),
  })
}

export function useResolveAlert() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (id: string) => {
      return orvalAxios({ url: `/alerts/${id}/resolve`, method: 'POST' })
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: alertKeys.all }),
  })
}
