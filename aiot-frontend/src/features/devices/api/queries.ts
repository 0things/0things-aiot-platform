import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { deviceServiceApi, axiosInstance } from '@/api/clients'
import { DEVICE_SERVICE_BASE_URL } from '@/api/config'
import type {
  DeviceV1CreateDeviceRequest,
  DeviceV1UpdateDeviceRequest,
  DeviceV1ActivateDeviceRequest,
  DeviceV1SetDeviceEnabledRequest,
  DeviceV1BatchUploadDevicesResponse,
} from '@/api/generated/device-service'

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
      // Axios will serialize array to repeated query params: ?states=online&states=offline
      const response = await deviceServiceApi.deviceServiceListDevices(
        {
          page: params.page,
          pageSize: params.pageSize,
          productId: params.productId,
          states: params.states,
          enabled: params.enabled,
        },
        {
          // Pass searchText as additional query parameter via axios options
          params: {
            ...(params.searchText && { searchText: params.searchText }),
          },
        }
      )
      return response.data
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
      const response = await deviceServiceApi.deviceServiceGetDeviceByKey({
        deviceKey,
      })
      return response.data
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
      const response = await deviceServiceApi.deviceServiceGetDeviceTelemetry({
        deviceKey,
      })
      return response.data
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
      const response = await deviceServiceApi.deviceServiceGetDeviceStatistics()
      return response.data
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
      const response = await deviceServiceApi.deviceServiceGetMqttParameters({
        deviceKey,
      })
      return response.data
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
    mutationFn: async (data: DeviceV1CreateDeviceRequest) => {
      const response = await deviceServiceApi.deviceServiceCreateDevice({
        deviceV1CreateDeviceRequest: data,
      })
      return response.data
    },
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
      const response = await deviceServiceApi.deviceServiceUpdateDevice({
        id,
        deviceV1UpdateDeviceRequest: data,
      })
      return response.data
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
    mutationFn: async (id: string) => {
      const response = await deviceServiceApi.deviceServiceDeleteDevice({ id })
      return response.data
    },
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
      data,
    }: {
      id: string
      data: DeviceV1ActivateDeviceRequest
    }) => {
      const response = await deviceServiceApi.deviceServiceActivateDevice({
        id,
        deviceV1ActivateDeviceRequest: data,
      })
      return response.data
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
      const response = await deviceServiceApi.deviceServiceSetDeviceEnabled({
        id,
        deviceV1SetDeviceEnabledRequest: data,
      })
      return response.data
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
      const formData = new FormData()
      formData.append('file', file)

      const response = await axiosInstance.post<DeviceV1BatchUploadDevicesResponse>(
        `${DEVICE_SERVICE_BASE_URL}/v1/devices/batch/upload`,
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        }
      )
      return response.data
    },
    onSuccess: (response) => {
      // Only invalidate if some devices were created successfully
      if (response.successCount && response.successCount > 0) {
        queryClient.invalidateQueries({ queryKey: deviceKeys.lists() })
      }
    },
  })
}
