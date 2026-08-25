import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import * as apiClient from './push-records-client'

// Query Key Factory Pattern
export interface PushRecordFilters {
  page?: number
  pageSize?: number
  operationType?: string
  status?: string
}

export const pushRecordKeys = {
  all: ['push-records'] as const,
  lists: () => [...pushRecordKeys.all, 'list'] as const,
  list: (deviceKey: string, filters?: PushRecordFilters) =>
    [...pushRecordKeys.lists(), deviceKey, filters] as const,
  details: () => [...pushRecordKeys.all, 'detail'] as const,
  detail: (recordId: string, deviceKey: string) =>
    [...pushRecordKeys.details(), deviceKey, recordId] as const,
}

/**
 * Hook for simulating a device push
 */
export function useSimulatePush() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (params: { deviceKey: string; payload: string }) => {
      const response = await apiClient.simulatePush({
        deviceKey: params.deviceKey,
        request: {
          payload: params.payload,
        },
      })
      return response
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: pushRecordKeys.list(variables.deviceKey),
      })
    },
  })
}

/**
 * Hook for listing push records for a device
 */
export function usePushRecords(
  deviceKey: string,
  options?: {
    page?: number
    pageSize?: number
    operationType?: string
    status?: string
  }
) {
  return useQuery({
    queryKey: pushRecordKeys.list(deviceKey, options),
    queryFn: async () => {
      const response = await apiClient.listPushRecords({
        deviceKey,
        page: options?.page || 1,
        pageSize: options?.pageSize || 20,
        operationType: options?.operationType,
        status: options?.status,
      })
      return response
    },
    enabled: !!deviceKey,
  })
}

/**
 * Hook for getting a specific push record
 */
export function usePushRecord(deviceKey: string, recordId: string) {
  return useQuery({
    queryKey: pushRecordKeys.detail(recordId, deviceKey),
    queryFn: async () => {
      const response = await apiClient.getPushRecord(deviceKey, recordId)
      return response
    },
    enabled: !!deviceKey && !!recordId,
  })
}

/**
 * Hook for clearing old push records
 */
export function useClearPushRecords() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (params: {
      deviceKey: string
      beforeTimestamp?: string
    }) => {
      const response = await apiClient.clearPushRecords({
        deviceKey: params.deviceKey,
        beforeTimestamp: params.beforeTimestamp,
      })
      return response
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: pushRecordKeys.list(variables.deviceKey),
      })
    },
  })
}

/**
 * Hook for getting product TSL (Thing Specification Language)
 */
export function useProductTSL(productKey: string) {
  return useQuery({
    queryKey: ['product-tsl', productKey],
    queryFn: async () => {
      const response = await apiClient.getProductTSL(productKey)
      return response
    },
    enabled: !!productKey,
  })
}
