import { useQuery } from '@tanstack/react-query'
import type {
  OtaOTAPackage,
  OtaUpgradeStatistics,
  OtaDeviceDeployment,
  OtaUpgradeBatch,
} from '@/api/generated/model'
import {
  getOtaPackagesId,
  getOtaPackagesIdBatches,
  getOtaPackagesIdDeviceDeployments,
  getOtaPackagesIdUpgradeStatistics,
} from '@/api/generated'

interface OTAPackageDetail {
  id: string
  packageName: string
  version?: string
  packageType?: string
  productId?: string
  productKey?: string
  status?: string
  fileUrl?: string
  checksum?: string
  signature: string
  verificationStatus: string
  verificationProgress: number
  description?: string
  releaseNotes?: string
  createdAt: Date
  updatedAt: Date
}

interface UpgradeBatch {
  batchId: string
  batchName: string
  batchType: string
  upgradeStrategy: string
  status: string
  targetDeviceCount: number
  createdAt: Date
}

interface DeviceDeployment {
  deviceId: string
  deviceKey: string
  deviceName: string
  productId: string
  productKey: string
  currentVersion: string
  upgradeBatchId: string
  status: string
  lastStatusChangeTime: string | number
  createdAt: Date
}

interface DeviceDeploymentList {
  deployments: DeviceDeployment[]
  total: number
  page: number
  pageSize: number
}

/**
 * Hook to fetch OTA package details by ID
 * Calls: GET /v1/ota-packages/{id}
 */
export function useOTAPackageDetail(id: string) {
  return useQuery({
    queryKey: ['ota-package-detail', id],
    queryFn: async (): Promise<OTAPackageDetail> => {
      const response = await getOtaPackagesId(Number(id))
      const data = response.data?.otaPackage as OtaOTAPackage | undefined
      return {
        id: data?.id?.toString() || id,
        packageName: data?.packageName || id,
        version: data?.version,
        packageType: data?.packageType,
        productId: data?.productId?.toString(),
        productKey: data?.productKey,
        status: data?.status,
        fileUrl: data?.fileUrl,
        checksum: data?.checksum,
        signature: 'md5',
        verificationStatus: 'completed',
        verificationProgress: 100,
        description: data?.description,
        releaseNotes: data?.releaseNotes,
        createdAt: data?.createdAt ? new Date(data.createdAt) : new Date(),
        updatedAt: data?.updatedAt ? new Date(data.updatedAt) : new Date(),
      }
    },
    staleTime: 0, // No caching
  })
}

/**
 * Hook to fetch upgrade statistics for an OTA package
 * Calls: GET /v1/ota-packages/{packageName}/upgrade-statistics
 */
export function useUpgradeStatistics(packageName: string) {
  return useQuery({
    queryKey: ['upgrade-statistics', packageName],
    queryFn: async () => {
      const response = await getOtaPackagesIdUpgradeStatistics(Number(packageName))
      const data: OtaUpgradeStatistics =
        (response.data as unknown as OtaUpgradeStatistics) ?? {}
        return {
          packageId: data.packageId || packageName,
          totalTargetDevices: data.totalTargetDevices || 0,
          successfulUpgrades: data.successfulUpgrades || 0,
          failedUpgrades: data.failedUpgrades || 0,
          cancelledUpgrades: data.cancelledUpgrades || 0,
          pendingUpgrades: data.pendingUpgrades || 0,
          inProgressUpgrades: data.inProgressUpgrades || 0,
        }
    },
    staleTime: 0, // No caching
    enabled: !!packageName, // Only run query if packageName is provided
  })
}

/**
 * Hook to fetch device deployment status for an OTA package
 * Supports pagination and filtering
 * Calls: GET /v1/ota-packages/{packageName}/device-deployments?page=X&pageSize=Y&status=Z
 */
export function useDeviceDeployments(
  packageName: string,
  page = 1,
  pageSize = 100,
  status?: string
) {
  return useQuery<DeviceDeploymentList>({
    queryKey: ['device-deployments', packageName, page, pageSize, status],
    queryFn: async () => {
      const data = (await getOtaPackagesIdDeviceDeployments(Number(packageName), {
        page,
        pageSize,
        status,
      }))?.data ?? {}
        return {
          deployments: (data.deployments || []).map((d: OtaDeviceDeployment) => ({
            deviceId: String(d.deviceId ?? ''),
            deviceKey: d.deviceKey ?? '',
            deviceName: d.deviceName ?? '',
            productId: String(d.productId ?? ''),
            productKey: d.productKey ?? '',
            currentVersion: d.currentVersion ?? '',
            upgradeBatchId: d.upgradeBatchId ?? '',
            status: d.status ?? '',
            lastStatusChangeTime: d.lastStatusChangeTime ?? 0,
            createdAt: d.createdAt ? new Date(d.createdAt) : new Date(),
          })),
          total: data.total || 0,
          page: data.page || page,
          pageSize: data.pageSize || pageSize,
        }
    },
    staleTime: 0, // No caching
    enabled: !!packageName, // Only run query if packageName is provided
  })
}

/**
 * Hook to fetch upgrade batches for an OTA package
 * Calls: GET /v1/ota-packages/{packageName}/batches
 */
export function useUpgradeBatches(packageName: string) {
  return useQuery<UpgradeBatch[]>({
    queryKey: ['upgrade-batches', packageName],
    queryFn: async () => {
      const data = (await getOtaPackagesIdBatches(Number(packageName)))?.data ?? {}
      return (data.batches || []).map((b: OtaUpgradeBatch) => ({
          batchId: b.batchId ?? '',
          batchName: b.batchName ?? '',
          batchType: b.batchType ?? '',
          upgradeStrategy: b.upgradeStrategy ?? '',
          status: b.status ?? '',
          targetDeviceCount: b.targetDeviceCount ?? 0,
          createdAt: b.createdAt ? new Date(b.createdAt) : new Date(),
        }))
    },
    staleTime: 0, // No caching
    enabled: !!packageName, // Only run query if packageName is provided
  })
}
