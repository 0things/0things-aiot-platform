import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  getDevicesDeviceKeyShadow,
  getDevicesDeviceKeyShadowHistory,
  putDevicesDeviceKeyShadowDesired,
} from '@/api/generated'

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
  history: (deviceKey: string) =>
    [...shadowKeys.all, deviceKey, 'history'] as const,
}

export function useDeviceShadow(deviceKey: string) {
  return useQuery({
    queryKey: shadowKeys.detail(deviceKey),
    enabled: !!deviceKey,
    queryFn: async (): Promise<ShadowDoc> => {
      const res = await getDevicesDeviceKeyShadow(deviceKey)
      return (res?.data ?? res) as ShadowDoc
    },
  })
}

export function useUpdateDesired(deviceKey: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (input: {
      desired: Record<string, unknown>
      version: number
    }) => {
      const res = await putDevicesDeviceKeyShadowDesired(deviceKey, input)
      return (res?.data ?? res) as ShadowDoc
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
      const res = await getDevicesDeviceKeyShadowHistory(deviceKey)
      return (res?.data ?? res) as Array<{
        version: number
        updatedAt: string
        source: string
        desired: unknown
        reported: unknown
      }>
    },
  })
}
