import { keepPreviousData, useQuery } from '@tanstack/react-query'
import {
  getDevicesDeviceKeyThingModelProperties,
  getDevicesDeviceKeyTelemetryHistory,
} from '@/api/generated'
import type { GetDevicesDeviceKeyTelemetryHistoryParams } from '@/api/generated/model'
import type { ThingModelProperty } from '@/api/generated/model'

export interface TelemetryPoint {
  timestamp: number
  property: string
  value: number | string | boolean | Record<string, unknown>
}

export interface TelemetryHistoryQueryParams {
  deviceKey: string
  property: string
  durationMs?: number
  startTime?: number
  endTime?: number
  limit?: number
}

// ============================================================================
// Compose React Query hooks from the generated API client.
// ============================================================================

/**
 * Fetch the device's latest reported property values.
 */
export function useThingModelProperties(deviceKey: string, enabled = true) {
  return useQuery({
    queryKey: ['devices', deviceKey, 'thing-model', 'properties'],
    queryFn: async ({ signal }) => {
      const response = await getDevicesDeviceKeyThingModelProperties(
        deviceKey,
        signal
      )
      return (response.data?.items || []) as ThingModelProperty[]
    },
    enabled: Boolean(deviceKey) && enabled,
    placeholderData: keepPreviousData,
  })
}

/**
 * Query telemetry history through the generated API client.
 */
export function useTelemetryHistory(
  params: TelemetryHistoryQueryParams,
  enabled = true
) {
  return useQuery({
    queryKey: ['devices', params.deviceKey, 'telemetry', 'history', params],
    queryFn: async ({ signal }) => {
      let startTime = params.startTime
      let endTime = params.endTime
      if (params.durationMs) {
        endTime = Date.now()
        startTime = endTime - params.durationMs
      }

      const queryParams: GetDevicesDeviceKeyTelemetryHistoryParams = {
        property: params.property,
        limit: params.limit || 100,
      }
      if (startTime) queryParams.start_time = startTime
      if (endTime) queryParams.end_time = endTime

      const resp = await getDevicesDeviceKeyTelemetryHistory(
        params.deviceKey,
        queryParams,
        signal
      )
      return (resp.data as unknown as TelemetryPoint[]) || []
    },
    enabled: Boolean(params.deviceKey && params.property) && enabled,
    placeholderData: keepPreviousData,
    staleTime: 5000,
  })
}
