import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  getAlerts,
  postAlertsIdAck,
  postAlertsIdResolve,
} from '@/api/generated'

export type AlertSeverity = 'info' | 'warning' | 'critical'
export type AlertStatus = 'open' | 'acknowledged' | 'resolved' | 'snoozed'

export type Alert = {
  id: string
  tenantId: string
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
  list: (filters: Record<string, unknown>) => [...alertKeys.all, 'list', filters] as const,
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
      const res = await getAlerts({
        severity: filters.severity,
        status: filters.status,
        device_key: filters.deviceKey,
        page: filters.page,
        pageSize: filters.pageSize,
      })
      return (res?.data ?? res) as unknown as Promise<{ alerts: Alert[]; total: number }>
    },
  })
}

export function useOpenAlertCount() {
  return useQuery({
    queryKey: alertKeys.openCount(),
    queryFn: async () => {
      const data = await getAlerts({ status: 'open', pageSize: 0 })
      return ((data?.data ?? data) as { total: number }).total
    },
    refetchInterval: 15_000,
  })
}

export function useAckAlert() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (id: string) => {
      return postAlertsIdAck(Number(id))
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: alertKeys.all }),
  })
}

export function useResolveAlert() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (id: string) => {
      return postAlertsIdResolve(Number(id))
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: alertKeys.all }),
  })
}
