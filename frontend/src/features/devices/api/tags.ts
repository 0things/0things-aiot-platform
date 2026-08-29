import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  deleteDevicesDeviceKeyTags,
  getDevicesDeviceKeyTags,
  postDevicesDeviceKeyTags,
  putDevicesDeviceKeyTags,
} from '@/api/generated'

export type DeviceTag = { key: string; value: string; source?: string }

export const tagKeys = {
  all: ['device-tags'] as const,
  list: (deviceKey: string) => [...tagKeys.all, deviceKey] as const,
}

export function useDeviceTags(deviceKey: string) {
  return useQuery({
    queryKey: tagKeys.list(deviceKey),
    enabled: !!deviceKey,
    queryFn: async (): Promise<DeviceTag[]> => {
      const data = await getDevicesDeviceKeyTags(deviceKey)
      return (data?.data?.tags ?? []) as DeviceTag[]
    },
  })
}

export function useSetTags(deviceKey: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (tags: Record<string, string>) => {
      return putDevicesDeviceKeyTags(deviceKey, { tags })
    },
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: tagKeys.list(deviceKey) }),
  })
}

export function useAddTags(deviceKey: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (tags: Record<string, string>) => {
      return postDevicesDeviceKeyTags(deviceKey, { tags })
    },
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: tagKeys.list(deviceKey) }),
  })
}

export function useRemoveTags(deviceKey: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (keys: string[]) => {
      return deleteDevicesDeviceKeyTags(deviceKey, { keys })
    },
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: tagKeys.list(deviceKey) }),
  })
}
