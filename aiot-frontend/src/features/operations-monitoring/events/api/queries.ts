import { useQuery } from '@tanstack/react-query'
import { axiosInstance } from '@/api/clients'
import { DEVICE_SERVICE_BASE_URL } from '@/api/config'

export type DeviceEvent = {
  id: number
  deviceKey: string
  deviceName: string
  productName: string
  eventType: string
  eventAt: string
  data: string
}

type DeviceEventsResponse = {
  events: DeviceEvent[]
  total: number
  page: number
  pageSize: number
}

export type DeviceEventFilters = {
  page: number
  pageSize: number
  keyword?: string
  eventType?: string
  startAt?: string
  endAt?: string
}

export const deviceEventKeys = {
  all: ['device-events'] as const,
  list: (filters: DeviceEventFilters) =>
    [...deviceEventKeys.all, filters] as const,
}

export function useDeviceEvents(filters: DeviceEventFilters) {
  return useQuery({
    queryKey: deviceEventKeys.list(filters),
    queryFn: async () => {
      const response = await axiosInstance.get<DeviceEventsResponse>(
        `${DEVICE_SERVICE_BASE_URL}/v1/device-events`,
        {
          params: {
            page: filters.page,
            pageSize: filters.pageSize,
            keyword: filters.keyword || undefined,
            event_type: filters.eventType || undefined,
            start_at: filters.startAt || undefined,
            end_at: filters.endAt || undefined,
          },
        }
      )
      return response.data
    },
  })
}
