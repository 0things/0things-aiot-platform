import { deviceServiceApi, productTSLServiceApi } from '@/api/clients'
import type {
  DeviceV1SimulatePushRequest,
} from '@/api/generated/device-service'

/**
 * Simulate a device push
 */
export async function simulatePush(params: {
  deviceKey: string
  request: DeviceV1SimulatePushRequest
}) {
  return deviceServiceApi.deviceServiceSimulatePush({
    deviceKey: params.deviceKey,
    deviceV1SimulatePushRequest: params.request,
  })
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
  return deviceServiceApi.deviceServiceListPushRecords({
    deviceKey: params.deviceKey,
    page: params.page,
    pageSize: params.pageSize,
    operationType: params.operationType,
    status: params.status,
  })
}

/**
 * Get a specific push record
 */
export async function getPushRecord(pushRecordId: string) {
  return deviceServiceApi.deviceServiceGetPushRecord({
    pushRecordId,
  })
}

/**
 * Clear push records for a device
 */
export async function clearPushRecords(params: {
  deviceKey: string
  beforeTimestamp?: string
}) {
  return deviceServiceApi.deviceServiceClearPushRecords({
    deviceKey: params.deviceKey,
    beforeTimestamp: params.beforeTimestamp,
  })
}

/**
 * Get product TSL (Thing Specification Language)
 */
export async function getProductTSL(productKey: string) {
  return productTSLServiceApi.productTSLServiceGetProductTSL({
    productKey,
  })
}
