import {
  getGetDevicesDeviceKeyEndpointsQueryKey,
  useGetDevicesDeviceKeyEndpoints,
} from '@/api/generated'

export type DeviceEndpoints = {
  http?: { http: string; rpcSubscribe: string }
  mqtt?: {
    host: string
    port: string
    telemetryTopic: string
    attributesTopic: string
    attributesSubscribeTopic: string
    rpcSubscribeTopic: string
  }
  coap?: { coap: string; docker?: { coap: string }; rpcSubscribe: string }
}

export const endpointKeys = {
  all: ['device-endpoints'] as const,
  detail: (deviceKey: string) =>
    getGetDevicesDeviceKeyEndpointsQueryKey(deviceKey),
}

export function useDeviceEndpoints(deviceKey: string) {
  return useGetDevicesDeviceKeyEndpoints(deviceKey, {
    query: {
      select: (res) => res?.data as unknown as DeviceEndpoints,
      enabled: !!deviceKey,
    },
  })
}
