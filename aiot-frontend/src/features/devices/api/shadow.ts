import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { axiosInstance } from '@/api/clients'
import { DEVICE_SERVICE_BASE_URL } from '@/api/config'

export type ShadowDoc = {
  desired: Record<string, unknown>
  reported: Record<string, unknown>
  delta: Record<string, unknown>
  metadata?: Record<string, unknown>
  version: number
  updatedAt?: string
}

export const shadowKeys = {
  all: ['device-shadow'] as const,
  detail: (deviceKey: string) => [...shadowKeys.all, deviceKey] as const,
  history: (deviceKey: string) => [...shadowKeys.all, deviceKey, 'history'] as const,
}

export function useDeviceShadow(deviceKey: string) {
  return useQuery({
    queryKey: shadowKeys.detail(deviceKey),
    enabled: !!deviceKey,
    queryFn: async (): Promise<ShadowDoc> => {
      const { data } = await axiosInstance.get(
        `${DEVICE_SERVICE_BASE_URL}/v1/devices/${encodeURIComponent(deviceKey)}/shadow`
      )
      return data
    },
  })
}

export function useUpdateDesired(deviceKey: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (input: { desired: Record<string, unknown>; version: number }) => {
      const { data } = await axiosInstance.put(
        `${DEVICE_SERVICE_BASE_URL}/v1/devices/${encodeURIComponent(deviceKey)}/shadow/desired`,
        input
      )
      return data as ShadowDoc
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: shadowKeys.detail(deviceKey) })
    },
  })
}

export function useShadowHistory(deviceKey: string) {
  return useQuery({
    queryKey: shadowKeys.history(deviceKey),
    enabled: !!deviceKey,
    queryFn: async () => {
      const { data } = await axiosInstance.get(
        `${DEVICE_SERVICE_BASE_URL}/v1/devices/${encodeURIComponent(deviceKey)}/shadow/history`
      )
      return data as Array<{ version: number; updatedAt: string; source: string; desired: unknown; reported: unknown }>
    },
  })
}
