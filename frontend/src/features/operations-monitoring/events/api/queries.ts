import { useQuery } from '@tanstack/react-query'
import { getDeviceEvents } from '@/api/generated'

export type DeviceEvent = {
  id: number
  deviceKey: string
  deviceName: string
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
      const res = await getDeviceEvents({
        page: filters.page,
        pageSize: filters.pageSize,
        keyword: filters.keyword,
        event_type: filters.eventType,
        start_at: filters.startAt,
        end_at: filters.endAt,
      })
      return (res?.data ?? res) as Promise<DeviceEventsResponse>
    },
  })
}
