import type {
  DeviceSimulatePushRequest,
} from '@/api/generated/model'
import {
  deleteDevicesIdPushRecords,
  getDevicesIdPushRecords,
  getDevicesPushRecordsPushRecordId,
  postDevicesIdSimulatePush,
} from '@/api/generated'
import {
  getProductsIdTsl,
  getProductsKeyProductKey,
} from '@/api/generated'
import { getDeviceId } from './device-id'

/**
 * Simulate a device push
 */
export async function simulatePush(params: {
  deviceKey: string
  request: DeviceSimulatePushRequest
}) {
  return postDevicesIdSimulatePush(await getDeviceId(params.deviceKey), params.request)
}

/**
 * List push records for a device
 */
export async function listPushRecords(params: {
  deviceKey: string
  page?: number
  pageSize?: number
  operationType?: string
  status?: string
}) {
  return getDevicesIdPushRecords(await getDeviceId(params.deviceKey), {
    page: params.page,
    pageSize: params.pageSize,
    operationType: params.operationType,
    status: params.status,
  })
}

/**
 * Get a specific push record
 */
export async function getPushRecord(deviceKey: string, pushRecordId: string) {
  await getDeviceId(deviceKey)
  return getDevicesPushRecordsPushRecordId(Number(pushRecordId))
}

/**
 * Clear push records for a device
 */
export async function clearPushRecords(params: {
  deviceKey: string
  beforeTimestamp?: string
}) {
  return deleteDevicesIdPushRecords(await getDeviceId(params.deviceKey), {
    beforeTimestamp: params.beforeTimestamp ? Number(params.beforeTimestamp) : undefined,
  })
}

/**
 * Get product TSL (Thing Specification Language)
 */
export async function getProductTSL(productKey: string) {
  const product = await getProductsKeyProductKey(productKey)
  if (product.product?.id === undefined) throw new Error('Product not found')
  return getProductsIdTsl(product.product.id)
}
