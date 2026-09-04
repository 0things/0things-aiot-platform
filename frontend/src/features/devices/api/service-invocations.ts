import {
  getGetDevicesDeviceKeyThingModelServiceInvocationsQueryKey,
  useGetDevicesDeviceKeyThingModelServiceInvocations,
} from '@/api/generated'
import type {
  DeviceListServiceInvocationsResponse,
  DeviceServiceInvocation,
  GetDevicesDeviceKeyThingModelServiceInvocationsParams,
} from '@/api/generated/model'

export type {
  DeviceServiceInvocation,
  GetDevicesDeviceKeyThingModelServiceInvocationsParams,
}

export const serviceInvocationKeys = {
  all: ['thing-model-service-invocations'] as const,
  list: (
    deviceKey: string,
    params?: GetDevicesDeviceKeyThingModelServiceInvocationsParams
  ) =>
    getGetDevicesDeviceKeyThingModelServiceInvocationsQueryKey(
      deviceKey,
      params
    ),
}

export function useServiceInvocations(
  deviceKey: string,
  params: GetDevicesDeviceKeyThingModelServiceInvocationsParams
) {
  return useGetDevicesDeviceKeyThingModelServiceInvocations(deviceKey, params, {
    query: {
      select: (res) =>
        res?.data as DeviceListServiceInvocationsResponse | undefined,
    },
  })
}
