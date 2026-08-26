import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  getOtaPackagesUuid,
  getOtaPackagesUuidBatches,
  getOtaPackagesUuidDeviceDeployments,
  getOtaPackagesUuidUpgradeStatistics,
  postOtaPackagesUuidBatchUpgrade,
} from '@/api/generated'
import type {
  OtaUpgradeStatistics,
  OtaDeviceDeployment,
  OtaUpgradeBatch,
} from '@/api/generated/model'

interface OTAPackageDetail {
  id: string
  uuid: string
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
  createdAt: string
  updatedAt: string
}

export interface UpgradeBatch {
  batchId: string
  upgradeStrategy: string
  status: string
  targetDeviceCount: number
  createdAt: string
}

interface DeviceDeployment {
  deviceId: string
  deviceKey: string
  deviceName: string
  productId: string
  productKey: string
  currentVersion: string
  status: string
  upgradeBatchId: string
  lastStatusChangeTime: string | number
  createdAt: string
}

interface DeviceDeploymentList {
  deployments: DeviceDeployment[]
  total: number
  page: number
  pageSize: number
}

/**
 * Hook to fetch OTA package details by UUID
 * Calls: GET /v1/ota-packages/{uuid}
 */
export function useOTAPackageDetail(uuid: string) {
  return useQuery({
    queryKey: ['ota-package-detail', uuid],
    queryFn: async (): Promise<OTAPackageDetail> => {
      const response = await getOtaPackagesUuid(uuid)
      const data = response.data
      return {
        id: data?.id?.toString() || uuid,
        uuid: data?.uuid || uuid,
        packageName: data?.packageName || uuid,
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
        createdAt: data?.createdAt ?? '',
        updatedAt: data?.updatedAt ?? '',
      }
    },
    staleTime: 0, // No caching
  })
}

/**
 * Hook to fetch upgrade statistics for an OTA package
 * Calls: GET /v1/ota-packages/{uuid}/upgrade-statistics
 */
export function useUpgradeStatistics(uuid: string) {
  return useQuery({
    queryKey: ['upgrade-statistics', uuid],
    queryFn: async () => {
      const response = await getOtaPackagesUuidUpgradeStatistics(uuid)
      const data: OtaUpgradeStatistics = response.data?.statistics ?? {}
      return {
        packageId: data.packageId || uuid,
        totalTargetDevices: data.totalTargetDevices || 0,
        successfulUpgrades: data.successfulUpgrades || 0,
        failedUpgrades: data.failedUpgrades || 0,
        cancelledUpgrades: data.cancelledUpgrades || 0,
        pendingUpgrades: data.pendingUpgrades || 0,
        inProgressUpgrades: data.inProgressUpgrades || 0,
      }
    },
    staleTime: 0, // No caching
    enabled: !!uuid, // Only run query if uuid is provided
  })
}

/**
 * Hook to fetch device deployment status for an OTA package
 * Supports pagination and filtering
 * Calls: GET /v1/ota-packages/{packageName}/device-deployments?page=X&pageSize=Y&status=Z
 */
export function useDeviceDeployments(
  uuid: string,
  page = 1,
  pageSize = 100,
  status?: string,
  refreshKey = 0
) {
  return useQuery<DeviceDeploymentList>({
    queryKey: ['device-deployments', uuid, page, pageSize, status, refreshKey],
    queryFn: async () => {
      const data =
        (
          await getOtaPackagesUuidDeviceDeployments(uuid, {
            page,
            pageSize,
            status,
          })
        )?.data ?? {}
      return {
        deployments: (data.items || []).map((d: OtaDeviceDeployment) => ({
          deviceId: String(d.deviceId ?? ''),
          deviceKey: d.deviceKey ?? '',
          deviceName: d.deviceName ?? '',
          productId: String(d.productId ?? ''),
          productKey: d.productKey ?? '',
          currentVersion: d.currentVersion ?? '',
          status: d.status ?? '',
          upgradeBatchId: d.upgradeBatchId ?? '',
          lastStatusChangeTime: d.lastStatusChangeTime ?? '',
          createdAt: d.createdAt ?? '',
        })),
        total: data.total || 0,
        page: data.page || page,
        pageSize: data.pageSize || pageSize,
      }
    },
    staleTime: 0, // No caching
    enabled: !!uuid, // Only run query if uuid is provided
  })
}

/**
 * Hook to fetch upgrade batches for an OTA package
 * Calls: GET /v1/ota-packages/{uuid}/batches
 */
export function useUpgradeBatches(uuid: string) {
  return useQuery<UpgradeBatch[]>({
    queryKey: ['upgrade-batches', uuid],
    queryFn: async () => {
      const data = (await getOtaPackagesUuidBatches(uuid))?.data ?? {}
      return (data.items || []).map((b: OtaUpgradeBatch) => ({
        batchId: b.batchId ?? '',
        upgradeStrategy: b.upgradeStrategy ?? '',
        status: b.status ?? '',
        targetDeviceCount: b.targetDeviceCount ?? 0,
        createdAt: b.createdAt ?? '',
      }))
    },
    staleTime: 0, // No caching
    enabled: !!uuid, // Only run query if uuid is provided
  })
}

/**
 * Mutation to create a static upgrade batch for an OTA package and enqueue the
 * selected devices for upgrade.
 * Calls: POST /v1/ota-packages/{uuid}/batch-upgrade
 */
export function useBatchUpgrade(uuid: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (deviceKeys: string[]) => {
      const response = await postOtaPackagesUuidBatchUpgrade(uuid, {
        deviceKeys,
      })
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['upgrade-batches', uuid],
      })
      queryClient.invalidateQueries({
        queryKey: ['device-deployments', uuid],
      })
      queryClient.invalidateQueries({
        queryKey: ['upgrade-statistics', uuid],
      })
    },
  })
}
