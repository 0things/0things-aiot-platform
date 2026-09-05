import { useMutation, useQueryClient } from '@tanstack/react-query'
import {
  deleteDevicesDeviceKeyTags,
  getGetDevicesDeviceKeyTagsQueryKey,
  postDevicesDeviceKeyTags,
  putDevicesDeviceKeyTags,
  useGetDevicesDeviceKeyTags,
} from '@/api/generated'
import type { DeviceTag } from '@/api/generated/model'

export type { DeviceTag }

export const tagKeys = {
  all: ['device-tags'] as const,
  list: (deviceKey: string) => getGetDevicesDeviceKeyTagsQueryKey(deviceKey),
}

export function useDeviceTags(deviceKey: string) {
  return useGetDevicesDeviceKeyTags(deviceKey, {
    query: {
      select: (res) => (res?.data?.tags ?? []) as DeviceTag[],
      enabled: !!deviceKey,
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
