import { useQuery } from '@tanstack/react-query'
import type { OtaOTAPackage } from '@/api/generated/model'
import {
  getOtaPackagesId,
  getOtaPackagesIdBatches,
  getOtaPackagesIdDeviceDeployments,
  getOtaPackagesIdUpgradeStatistics,
} from '@/api/generated'

/**
 * Hook to fetch OTA package details by ID
 * Calls: GET /v1/ota-packages/{id}
 */
export function useOTAPackageDetail(id: string) {
  return useQuery({
    queryKey: ['ota-package-detail', id],
    queryFn: async () => {
      try {
        const response = await getOtaPackagesId(Number(id))
        const data = response.otaPackage as OtaOTAPackage | undefined
        return {
          id: data?.id?.toString() || id,
          packageName: data?.packageName || id,
          version: data?.version,
          packageType: data?.packageType,
          productId: data?.productId,
          productKey: data?.productKey,
          status: data?.status,
          fileUrl: data?.fileUrl,
          checksum: data?.checksum,
          signature: 'md5',
          verificationStatus: 'completed',
          verificationProgress: 100,
          description: data?.description,
          releaseNotes: data?.releaseNotes,
          metadata: data?.metadata ? JSON.parse(data.metadata) : {},
          createdAt: data?.createdAt ? new Date(data.createdAt) : new Date(),
          updatedAt: data?.updatedAt ? new Date(data.updatedAt) : new Date(),
        }
      } catch (error) {
        console.error('Failed to fetch OTA package detail:', error)
        throw error
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
      try {
        const data = (await getOtaPackagesIdUpgradeStatistics(Number(packageName))) as any
        return {
          packageId: data.packageId || packageName,
          totalTargetDevices: data.totalTargetDevices || 0,
          successfulUpgrades: data.successfulUpgrades || 0,
          failedUpgrades: data.failedUpgrades || 0,
          cancelledUpgrades: data.cancelledUpgrades || 0,
          pendingUpgrades: data.pendingUpgrades || 0,
          inProgressUpgrades: data.inProgressUpgrades || 0,
        }
      } catch (error) {
        console.error('Failed to fetch upgrade statistics:', error)
        throw error
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
  return useQuery({
    queryKey: ['device-deployments', packageName, page, pageSize, status],
    queryFn: async () => {
      try {
        const data = await getOtaPackagesIdDeviceDeployments(Number(packageName), {
          page,
          pageSize,
          status,
        })
        return {
          deployments: (data.deployments || []).map((d: any) => ({
            deviceId: d.deviceId,
            deviceKey: d.deviceKey,
            deviceName: d.deviceName,
            productId: d.productId,
            productKey: d.productKey,
            currentVersion: d.currentVersion,
            upgradeBatchId: d.upgradeBatchId,
            status: d.status,
            lastStatusChangeTime: d.lastStatusChangeTime,
            createdAt: d.createdAt ? new Date(d.createdAt) : new Date(),
          })),
          total: data.total || 0,
          page: data.page || page,
          pageSize: data.pageSize || pageSize,
        }
      } catch (error) {
        console.error('Failed to fetch device deployments:', error)
        throw error
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
  return useQuery({
    queryKey: ['upgrade-batches', packageName],
    queryFn: async () => {
      try {
        const data = await getOtaPackagesIdBatches(Number(packageName))
        return (data.batches || []).map((b: any) => ({
          batchId: b.batchId,
          batchName: b.batchName,
          batchType: b.batchType,
          upgradeStrategy: b.upgradeStrategy,
          status: b.status,
          targetDeviceCount: b.targetDeviceCount,
          createdAt: b.createdAt ? new Date(b.createdAt) : new Date(),
        }))
      } catch (error) {
        console.error('Failed to fetch upgrade batches:', error)
        throw error
      }
    },
    staleTime: 0, // No caching
    enabled: !!packageName, // Only run query if packageName is provided
  })
}
