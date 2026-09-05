import { keepPreviousData } from '@tanstack/react-query'
import {
  getGetDevicesDeviceKeyTelemetryHistoryQueryKey,
  getGetDevicesDeviceKeyThingModelPropertiesQueryKey,
  useGetDevicesDeviceKeyTelemetryHistory,
  useGetDevicesDeviceKeyThingModelProperties,
} from '@/api/generated'
import type {
  GetDevicesDeviceKeyTelemetryHistoryParams,
  TelemetryPoint,
  ThingModelProperty,
} from '@/api/generated/model'

export type { TelemetryPoint }

export interface TelemetryHistoryQueryParams {
  deviceKey: string
  property: string
  durationMs?: number
  startTime?: number
  endTime?: number
  limit?: number
}

export const telemetryKeys = {
  all: ['device-telemetry'] as const,
  properties: (deviceKey: string) =>
    getGetDevicesDeviceKeyThingModelPropertiesQueryKey(deviceKey),
  history: (
    deviceKey: string,
    params?: GetDevicesDeviceKeyTelemetryHistoryParams
  ) => getGetDevicesDeviceKeyTelemetryHistoryQueryKey(deviceKey, params),
}

/**
 * Fetch the device's latest reported property values.
 */
export function useThingModelProperties(deviceKey: string, enabled = true) {
  return useGetDevicesDeviceKeyThingModelProperties(deviceKey, {
    query: {
      select: (res) => (res?.data?.items || []) as ThingModelProperty[],
      enabled: Boolean(deviceKey) && enabled,
      placeholderData: keepPreviousData,
    },
  })
}

/**
 * Query telemetry history through the generated API client.
 */
export function useTelemetryHistory(
  params: TelemetryHistoryQueryParams,
  enabled = true
) {
  let startTime = params.startTime
  let endTime = params.endTime
  if (params.durationMs) {
    endTime = Date.now()
    startTime = endTime - params.durationMs
  }

  const queryParams: GetDevicesDeviceKeyTelemetryHistoryParams = {
    property: params.property,
    limit: params.limit || 100,
    ...(startTime ? { start_time: startTime } : {}),
    ...(endTime ? { end_time: endTime } : {}),
  }

  return useGetDevicesDeviceKeyTelemetryHistory(params.deviceKey, queryParams, {
    query: {
      select: (res) => (res?.data as unknown as TelemetryPoint[]) || [],
      enabled: Boolean(params.deviceKey && params.property) && enabled,
      placeholderData: keepPreviousData,
      staleTime: 5000,
    },
  })
}
