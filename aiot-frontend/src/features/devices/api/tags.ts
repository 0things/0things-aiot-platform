import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { axiosInstance } from '@/api/clients'
import { DEVICE_SERVICE_BASE_URL } from '@/api/config'

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
      const { data } = await axiosInstance.get(
        `${DEVICE_SERVICE_BASE_URL}/v1/devices/${encodeURIComponent(deviceKey)}/tags`
      )
      return data?.tags ?? []
    },
  })
}

export function useSetTags(deviceKey: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (tags: Record<string, string>) => {
      const { data } = await axiosInstance.put(
        `${DEVICE_SERVICE_BASE_URL}/v1/devices/${encodeURIComponent(deviceKey)}/tags`,
        { tags }
      )
      return data
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: tagKeys.list(deviceKey) }),
  })
}

export function useAddTags(deviceKey: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (tags: Record<string, string>) => {
      const { data } = await axiosInstance.post(
        `${DEVICE_SERVICE_BASE_URL}/v1/devices/${encodeURIComponent(deviceKey)}/tags`,
        { tags }
      )
      return data
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: tagKeys.list(deviceKey) }),
  })
}

export function useRemoveTags(deviceKey: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (keys: string[]) => {
      const { data } = await axiosInstance.delete(
        `${DEVICE_SERVICE_BASE_URL}/v1/devices/${encodeURIComponent(deviceKey)}/tags`,
        { data: { keys } }
      )
      return data
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: tagKeys.list(deviceKey) }),
  })
}
