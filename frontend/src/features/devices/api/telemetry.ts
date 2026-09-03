import { keepPreviousData, useQuery } from '@tanstack/react-query'
import {
  getDevicesDeviceKeyShadow,
  getV1DevicesDeviceKeyTelemetryHistory,
} from '@/api/generated'
import type {
  DeviceShadow,
  GetV1DevicesDeviceKeyTelemetryHistoryParams,
} from '@/api/generated/model'

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
 * Fetch the live device shadow through the generated API client.
 */
export function useDeviceShadow(deviceKey: string, enabled = true) {
  return useQuery({
    queryKey: ['devices', deviceKey, 'shadow'],
    queryFn: async ({ signal }) => {
      const resp = await getDevicesDeviceKeyShadow(deviceKey, signal)
      return resp.data as DeviceShadow | undefined
    },
    enabled: Boolean(deviceKey) && enabled,
    placeholderData: keepPreviousData,
    refetchInterval: 3000, // Poll the shadow every three seconds.
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

      const queryParams: GetV1DevicesDeviceKeyTelemetryHistoryParams = {
        property: params.property,
        limit: params.limit || 100,
      }
      if (startTime) queryParams.start_time = startTime
      if (endTime) queryParams.end_time = endTime

      const resp = await getV1DevicesDeviceKeyTelemetryHistory(
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
