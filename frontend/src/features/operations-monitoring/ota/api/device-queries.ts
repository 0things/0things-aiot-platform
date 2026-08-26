import { getDevices } from '@/api/generated'
import type { Device } from '@/api/generated/model/device'

// 分页获取产品下全部启用设备，供批量升级弹窗复用。
export async function getAllEnabledDevices(
  productId: string
): Promise<Device[]> {
  const pageSize = 100
  const devices: Device[] = []
  let page = 1
  let total = 0
  do {
    const response = await getDevices({
      productId: Number(productId),
      page,
      pageSize,
      enabled: true,
    })
    const current = response.data?.devices || []
    devices.push(...current)
    total = response.data?.total || devices.length
    page += 1
    if (current.length === 0) break
  } while (devices.length < total)
  return devices
}
