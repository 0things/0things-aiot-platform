import {
  deleteDevicesDeviceKeyPushRecords,
  getDevicesDeviceKeyPushRecords,
  getDevicesPushRecordsPushRecordId,
  getProductsProductKeyTsl,
  postDevicesDeviceKeySimulatePush,
} from '@/api/generated'
import type { DeviceSimulatePushRequest } from '@/api/generated/model'

/**
 * Simulate a device push
 */
export async function simulatePush(params: {
  deviceKey: string
  request: DeviceSimulatePushRequest
}) {
  return postDevicesDeviceKeySimulatePush(params.deviceKey, params.request)
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
  return getDevicesDeviceKeyPushRecords(params.deviceKey, {
    page: params.page,
    pageSize: params.pageSize,
    operationType: params.operationType,
    status: params.status,
  })
}

/**
 * Get a specific push record
 */
export async function getPushRecord(_deviceKey: string, pushRecordId: string) {
  return getDevicesPushRecordsPushRecordId(Number(pushRecordId))
}

/**
 * Clear push records for a device
 */
export async function clearPushRecords(params: {
  deviceKey: string
  beforeTimestamp?: string
}) {
  return deleteDevicesDeviceKeyPushRecords(params.deviceKey, {
    beforeTimestamp: params.beforeTimestamp
      ? Number(params.beforeTimestamp)
      : undefined,
  })
}

/**
 * Get product TSL (Thing Specification Language)
 */
export async function getProductTSL(productKey: string) {
  return getProductsProductKeyTsl(productKey)
}
