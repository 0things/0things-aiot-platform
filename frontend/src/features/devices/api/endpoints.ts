import { useQuery } from '@tanstack/react-query'
import { axiosInstance } from '@/api/clients'
import { DEVICE_SERVICE_BASE_URL } from '@/api/config'

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

export function useDeviceEndpoints(deviceKey: string) {
  return useQuery({
    queryKey: ['devices', deviceKey, 'endpoints'],
    queryFn: async () => {
      const response = await axiosInstance.get(
        `${DEVICE_SERVICE_BASE_URL}/v1/devices/${deviceKey}/endpoints`
      )
      return response.data.data as DeviceEndpoints
    },
    enabled: !!deviceKey,
  })
}
