import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { axiosInstance } from '@/api/clients'
import { DEVICE_SERVICE_BASE_URL } from '@/api/config'

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
      const { data } = await axiosInstance.get(`${DEVICE_SERVICE_BASE_URL}/v1/alerts`, {
        params: filters,
      })
      return data as { alerts: Alert[]; total: number }
    },
  })
}

export function useOpenAlertCount() {
  return useQuery({
    queryKey: alertKeys.openCount(),
    queryFn: async () => {
      const { data } = await axiosInstance.get(`${DEVICE_SERVICE_BASE_URL}/v1/alerts`, {
        params: { status: 'open', pageSize: 0 },
      })
      return (data as { total: number }).total
    },
    refetchInterval: 15_000,
  })
}

export function useAckAlert() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (id: string) => {
      const { data } = await axiosInstance.post(
        `${DEVICE_SERVICE_BASE_URL}/v1/alerts/${id}/ack`,
        {}
      )
      return data
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: alertKeys.all }),
  })
}

export function useResolveAlert() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (id: string) => {
      const { data } = await axiosInstance.post(
        `${DEVICE_SERVICE_BASE_URL}/v1/alerts/${id}/resolve`,
        {}
      )
      return data
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: alertKeys.all }),
  })
}
