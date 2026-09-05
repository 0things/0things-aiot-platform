import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import type { Alert, GetAlertsParams } from '@/api/generated/model'
import { orvalAxios } from '@/api/orval-mutator'

export type AlertSeverity = 'info' | 'warning' | 'critical'
export type AlertStatus = 'open' | 'acknowledged' | 'resolved' | 'snoozed'

export type { Alert, GetAlertsParams }

export const alertKeys = {
  all: ['alerts'] as const,
  lists: () => [...alertKeys.all, 'list'] as const,
  list: (params?: Record<string, unknown>) =>
    [...alertKeys.lists(), params] as const,
  detail: (id: string) => [...alertKeys.all, 'detail', id] as const,
  openCount: () => [...alertKeys.all, 'open-count'] as const,
}

export function useAlerts(filters?: {
  severity?: string
  status?: string
  device_key?: string
  deviceKey?: string
  page?: number
  pageSize?: number
}) {
  const queryParams: GetAlertsParams = {
    page: filters?.page,
    pageSize: filters?.pageSize,
    severity: filters?.severity,
    status: filters?.status,
    device_key: filters?.device_key ?? filters?.deviceKey,
  }

  return useQuery<{ alerts: Alert[]; total: number }>({
    queryKey: alertKeys.list(queryParams as Record<string, unknown>),
    queryFn: async () => {
      const res = await orvalAxios<{
        data: { alerts: Alert[]; total: number }
      }>({
        url: '/alerts',
        method: 'GET',
        params: queryParams,
      })
      return (res?.data ?? res) as unknown as { alerts: Alert[]; total: number }
    },
  })
}

export function useOpenAlertCount() {
  return useQuery<number>({
    queryKey: alertKeys.openCount(),
    queryFn: async () => {
      const res = await orvalAxios<{ data: { total: number } }>({
        url: '/alerts',
        method: 'GET',
        params: { status: 'open', pageSize: 0 },
      })
      return res?.data?.total ?? 0
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
