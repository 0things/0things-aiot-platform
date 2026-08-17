import { getDevicesKeyDeviceKey } from '@/api/generated'

export async function getDeviceId(deviceKey: string) {
  const response = await getDevicesKeyDeviceKey(deviceKey)
  const id = response.data?.device?.id
  if (id === undefined) {
    throw new Error(`Device not found: ${deviceKey}`)
  }
  return id
}
