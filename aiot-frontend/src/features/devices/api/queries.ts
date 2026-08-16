import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type {
  DeviceBatchUploadDevicesResponse as DeviceV1BatchUploadDevicesResponse,
  DeviceCreateDeviceRequest as DeviceV1CreateDeviceRequest,
  DeviceGetDeviceResponse as DeviceV1GetDeviceResponse,
  DeviceMQTTParametersResponse as DeviceV1GetMqttParametersResponse,
  DeviceListDevicesResponse as DeviceV1ListDevicesResponse,
  DeviceSetDeviceEnabledRequest as DeviceV1SetDeviceEnabledRequest,
  DeviceTelemetryResponse as DeviceV1GetDeviceTelemetryResponse,
  DeviceUpdateDeviceRequest as DeviceV1UpdateDeviceRequest,
} from '@/api/generated/model'
import {
  deleteDevicesId,
  getDevices,
  getDevicesDeviceKeyMqttParameters,
  getDevicesIdTelemetry,
  getDevicesKeyDeviceKey,
  getDeviceStatistics,
  postDevices,
  postDevicesBatchUpload,
  postDevicesIdActivate,
  putDevicesId,
  postDevicesIdEnabled,
} from '@/api/generated'

// ============================================================================
// Query Keys
// ============================================================================

export const deviceKeys = {
  all: ['devices'] as const,
  lists: () => [...deviceKeys.all, 'list'] as const,
  list: (filters: {
    page?: number
    pageSize?: number
    productId?: string
    states?: string[]
    enabled?: boolean
    searchText?: string
  }) => [...deviceKeys.lists(), filters] as const,
  details: () => [...deviceKeys.all, 'detail'] as const,
  detail: (id: string) => [...deviceKeys.details(), id] as const,
  statistics: () => [...deviceKeys.all, 'statistics'] as const,
}

// ============================================================================
// Queries
// ============================================================================

/**
 * Hook to fetch devices list with pagination and filtering
 */
export function useDevices(params: {
  page?: number
  pageSize?: number
  productId?: string
  states?: string[]
  enabled?: boolean
  searchText?: string
}) {
  return useQuery({
    queryKey: deviceKeys.list(params),
    queryFn: async () => {
      return getDevices({
        ...params,
        productId: params.productId ? Number(params.productId) : undefined,
      }) as unknown as DeviceV1ListDevicesResponse
    },
  })
}

/**
 * Hook to fetch a single device by deviceKey
 */
export function useDeviceByKey(deviceKey: string) {
  return useQuery({
    queryKey: [...deviceKeys.details(), 'key', deviceKey],
    queryFn: async () => {
      return getDevicesKeyDeviceKey(deviceKey) as unknown as DeviceV1GetDeviceResponse
    },
    enabled: !!deviceKey,
  })
}

/**
 * Hook to fetch device telemetry
 */
export function useDeviceTelemetry(deviceKey: string) {
  return useQuery({
    queryKey: [...deviceKeys.details(), 'telemetry', deviceKey],
    queryFn: async () => {
      return getDevicesIdTelemetry(Number(deviceKey)) as unknown as DeviceV1GetDeviceTelemetryResponse
    },
    enabled: !!deviceKey,
  })
}

/**
 * Hook to fetch device statistics
 */
export function useDeviceStatistics() {
  return useQuery({
    queryKey: deviceKeys.statistics(),
    queryFn: async () => {
      return getDeviceStatistics() as never
    },
  })
}

/**
 * Hook to fetch MQTT parameters for a device
 */
export function useDeviceMqttParameters(deviceKey: string, enabled = false) {
  return useQuery({
    queryKey: [...deviceKeys.details(), 'mqtt-parameters', deviceKey],
    queryFn: async () => {
      return getDevicesDeviceKeyMqttParameters(deviceKey) as unknown as DeviceV1GetMqttParametersResponse
    },
    enabled: enabled && !!deviceKey,
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
    mutationFn: (data: DeviceV1CreateDeviceRequest) => postDevices(data as never),
    onSuccess: () => {
      // Invalidate all device lists to refetch
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
      return putDevicesId(Number(id), data as never) as never
    },
    onSuccess: (_, variables) => {
      // Invalidate the specific device detail
      queryClient.invalidateQueries({
        queryKey: deviceKeys.detail(variables.id),
      })
      // Invalidate all device lists to refetch
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
    mutationFn: (id: string) => deleteDevicesId(Number(id)),
    onSuccess: () => {
      // Invalidate all device lists to refetch
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
      return postDevicesIdActivate(Number(id)) as unknown as DeviceV1GetDeviceResponse
    },
    onSuccess: (_, variables) => {
      // Invalidate the specific device detail
      queryClient.invalidateQueries({
        queryKey: deviceKeys.detail(variables.id),
      })
      // Invalidate all device lists to refetch
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
      id,
      data,
    }: {
      id: string
      data: DeviceV1SetDeviceEnabledRequest
    }) => {
      return postDevicesIdEnabled(Number(id), data as never) as never
    },
    onSuccess: (_, variables) => {
      // Invalidate the specific device detail
      queryClient.invalidateQueries({
        queryKey: deviceKeys.detail(variables.id),
      })
      // Invalidate all device lists to refetch
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
      return postDevicesBatchUpload({ file }) as unknown as DeviceV1BatchUploadDevicesResponse
    },
    onSuccess: (response) => {
      // Only invalidate if some devices were created successfully
      if (response.successCount && response.successCount > 0) {
        queryClient.invalidateQueries({ queryKey: deviceKeys.lists() })
      }
    },
  })
}
