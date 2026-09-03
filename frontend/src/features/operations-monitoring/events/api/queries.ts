import { getGetDeviceEventsQueryKey, useGetDeviceEvents } from '@/api/generated'
import type {
  DeviceEvent,
  DeviceEventListDeviceEventsResponse,
  GetDeviceEventsParams,
} from '@/api/generated/model'

export type { DeviceEvent, GetDeviceEventsParams }
export type DeviceEventsResponse = DeviceEventListDeviceEventsResponse

export const deviceEventKeys = {
  all: ['/device-events'] as const,
  list: (params?: GetDeviceEventsParams) => getGetDeviceEventsQueryKey(params),
}

export function useDeviceEvents(params?: GetDeviceEventsParams) {
  return useGetDeviceEvents(params, {
    query: {
      select: (res) => res?.data,
    },
  })
}
