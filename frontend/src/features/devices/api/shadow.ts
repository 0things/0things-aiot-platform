import { useMutation, useQueryClient } from '@tanstack/react-query'
import {
  getGetDevicesDeviceKeyShadowHistoryQueryKey,
  getGetDevicesDeviceKeyShadowQueryKey,
  putDevicesDeviceKeyShadowDesired,
  useGetDevicesDeviceKeyShadow,
  useGetDevicesDeviceKeyShadowHistory,
} from '@/api/generated'
import type { DeviceShadow, DeviceShadowHistory } from '@/api/generated/model'

export type ShadowDoc = DeviceShadow
export type ShadowHistoryItem = DeviceShadowHistory

export const shadowKeys = {
  all: ['device-shadow'] as const,
  detail: (deviceKey: string) =>
    getGetDevicesDeviceKeyShadowQueryKey(deviceKey),
  history: (deviceKey: string) =>
    getGetDevicesDeviceKeyShadowHistoryQueryKey(deviceKey),
}

export function useDeviceShadow(deviceKey: string) {
  return useGetDevicesDeviceKeyShadow(deviceKey, {
    query: {
      select: (res) => res?.data as unknown as ShadowDoc,
      enabled: !!deviceKey,
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
      const res = await putDevicesDeviceKeyShadowDesired(
        deviceKey,
        input as never
      )
      return res?.data as unknown as ShadowDoc
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: shadowKeys.detail(deviceKey) })
    },
  })
}

export function useShadowHistory(deviceKey: string) {
  return useGetDevicesDeviceKeyShadowHistory(deviceKey, {
    query: {
      select: (res) => (res?.data as unknown as ShadowHistoryItem[]) || [],
      enabled: !!deviceKey,
    },
  })
}
