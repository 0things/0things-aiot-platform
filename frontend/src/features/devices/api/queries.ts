import { useMutation, useQueryClient } from '@tanstack/react-query'
import {
  deleteDevicesDeviceKey,
  getGetDevicesDeviceKeyQueryKey,
  getGetDevicesQueryKey,
  getGetDeviceStatisticsQueryKey,
  postDevices,
  postDevicesBatchUpload,
  postDevicesDeviceKeyActivate,
  postDevicesDeviceKeyEnabled,
  putDevicesDeviceKey,
  useGetDevices,
  useGetDevicesDeviceKey,
  useGetDevicesDeviceKeyTelemetry,
  useGetDeviceStatistics,
} from '@/api/generated'
import type {
  DeviceBatchUploadDevicesResponse as DeviceV1BatchUploadDevicesResponse,
  DeviceCreateDeviceRequest as DeviceV1CreateDeviceRequest,
  DeviceGetDeviceResponse as DeviceV1GetDeviceResponse,
  DeviceListDevicesResponse as DeviceV1ListDevicesResponse,
  DeviceSetDeviceEnabledRequest as DeviceV1SetDeviceEnabledRequest,
  DeviceUpdateDeviceRequest as DeviceV1UpdateDeviceRequest,
  GetDevicesParams,
} from '@/api/generated/model'

export type { GetDevicesParams }
export type DeviceListResponse = DeviceV1ListDevicesResponse

// ============================================================================
// Query Keys
// ============================================================================

export const deviceKeys = {
  all: ['devices'] as const,
  lists: () => [...deviceKeys.all, 'list'] as const,
  list: (params?: GetDevicesParams) => getGetDevicesQueryKey(params),
  details: () => [...deviceKeys.all, 'detail'] as const,
  detail: (deviceKey: string) => getGetDevicesDeviceKeyQueryKey(deviceKey),
  statistics: () => getGetDeviceStatisticsQueryKey(),
}

// ============================================================================
// Queries
// ============================================================================

/**
 * Hook to fetch devices list with pagination and filtering
 */
export function useDevices(params?: GetDevicesParams) {
  return useGetDevices(params, {
    query: {
      select: (res) => res?.data,
    },
  })
}

/**
 * Hook to fetch a single device by deviceKey
 */
export function useDeviceByKey(deviceKey: string) {
  return useGetDevicesDeviceKey(deviceKey, {
    query: {
      select: (res) => res?.data,
      enabled: !!deviceKey,
    },
  })
}

/**
 * Hook to fetch device telemetry
 */
export function useDeviceTelemetry(deviceKey: string) {
  return useGetDevicesDeviceKeyTelemetry(deviceKey, {
    query: {
      select: (res) => res?.data,
      enabled: !!deviceKey,
    },
  })
}

/**
 * Hook to fetch device statistics
 */
export function useDeviceStatistics() {
  return useGetDeviceStatistics({
    query: {
      select: (res) => res?.data,
    },
  })
}

// ============================================================================
// Mutations
// ============================================================================

/**
 * Hook to create a new device
 */
export function useCreateDevice() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: DeviceV1CreateDeviceRequest) =>
      postDevices(data as never),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: deviceKeys.lists() })
    },
  })
}

/**
 * Hook to update an existing device
 */
export function useUpdateDevice() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({
      id,
      data,
    }: {
      id: string
      data: DeviceV1UpdateDeviceRequest
    }) => {
      return putDevicesDeviceKey(id, data as never) as never
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: deviceKeys.detail(variables.id),
      })
      queryClient.invalidateQueries({ queryKey: deviceKeys.lists() })
    },
  })
}

/**
 * Hook to delete a device
 */
export function useDeleteDevice() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (deviceKey: string) => deleteDevicesDeviceKey(deviceKey),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: deviceKeys.lists() })
    },
  })
}

/**
 * Hook to activate a device
 */
export function useActivateDevice() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({
      id,
      data: _data,
    }: {
      id: string
      data: Record<string, never>
    }) => {
      return postDevicesDeviceKeyActivate(
        id
      ) as unknown as DeviceV1GetDeviceResponse
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: deviceKeys.detail(variables.id),
      })
      queryClient.invalidateQueries({ queryKey: deviceKeys.lists() })
    },
  })
}

/**
 * Hook to enable or disable a device
 */
export function useSetDeviceEnabled() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({
      deviceKey,
      data,
    }: {
      deviceKey: string
      data: DeviceV1SetDeviceEnabledRequest
    }) => {
      return postDevicesDeviceKeyEnabled(deviceKey, data as never) as never
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: deviceKeys.detail(variables.deviceKey),
      })
      queryClient.invalidateQueries({ queryKey: deviceKeys.lists() })
    },
  })
}

/**
 * Hook to batch upload devices from Excel/CSV file
 */
export function useBatchUploadDevices() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (file: File) => {
      const res = await postDevicesBatchUpload({ file })
      return (res?.data ?? res) as unknown as DeviceV1BatchUploadDevicesResponse
    },
    onSuccess: (response) => {
      if (response.successCount && response.successCount > 0) {
        queryClient.invalidateQueries({ queryKey: deviceKeys.lists() })
      }
    },
  })
}
