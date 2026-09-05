import { useMutation, useQueryClient } from '@tanstack/react-query'
import {
  deleteDevicesDeviceKeyPushRecords,
  getGetDevicesDeviceKeyPushRecordsQueryKey,
  getGetDevicesPushRecordsPushRecordIdQueryKey,
  getGetProductsProductKeyTslQueryKey,
  postDevicesDeviceKeySimulatePush,
  useGetDevicesDeviceKeyPushRecords,
  useGetDevicesPushRecordsPushRecordId,
  useGetProductsProductKeyTsl,
} from '@/api/generated'
import type {
  DeviceListPushRecordsResponse,
  DeviceSimulatePushRequest,
  GetDevicesDeviceKeyPushRecordsParams,
} from '@/api/generated/model'

export type PushRecordFilters = GetDevicesDeviceKeyPushRecordsParams

export const pushRecordKeys = {
  all: ['push-records'] as const,
  lists: () => [...pushRecordKeys.all, 'list'] as const,
  list: (deviceKey: string, filters?: PushRecordFilters) =>
    getGetDevicesDeviceKeyPushRecordsQueryKey(deviceKey, filters),
  details: () => [...pushRecordKeys.all, 'detail'] as const,
  detail: (recordId: string) =>
    getGetDevicesPushRecordsPushRecordIdQueryKey(Number(recordId)),
  tsl: (productKey: string) => getGetProductsProductKeyTslQueryKey(productKey),
}

/**
 * Hook for simulating a device push
 */
export function useSimulatePush() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (params: { deviceKey: string; payload: string }) => {
      const req: DeviceSimulatePushRequest = {
        payload: params.payload,
      }
      return postDevicesDeviceKeySimulatePush(params.deviceKey, req)
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ['/devices', variables.deviceKey, 'push-records'],
      })
    },
  })
}

/**
 * Hook for listing push records for a device
 */
export function usePushRecords(deviceKey: string, options?: PushRecordFilters) {
  return useGetDevicesDeviceKeyPushRecords(deviceKey, options, {
    query: {
      select: (res) => res?.data as unknown as DeviceListPushRecordsResponse,
      enabled: !!deviceKey,
    },
  })
}

/**
 * Hook for getting a specific push record
 */
export function usePushRecord(_deviceKey: string, recordId: string) {
  return useGetDevicesPushRecordsPushRecordId(Number(recordId), {
    query: {
      select: (res) => res?.data,
      enabled: !!recordId,
    },
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
      return deleteDevicesDeviceKeyPushRecords(params.deviceKey, {
        beforeTimestamp: params.beforeTimestamp
          ? Number(params.beforeTimestamp)
          : undefined,
      })
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ['/devices', variables.deviceKey, 'push-records'],
      })
    },
  })
}

/**
 * Hook for getting product TSL (Thing Specification Language)
 */
export function useProductTSL(productKey: string) {
  return useGetProductsProductKeyTsl(productKey, {
    query: {
      select: (res) => res?.data,
      enabled: !!productKey,
    },
  })
}
