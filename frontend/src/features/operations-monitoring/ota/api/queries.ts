import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type {
  OtaCreateOTAPackageRequest,
  OtaOTAPackageRequest as OtaV1UpdateOTAPackageRequest,
} from '@/api/generated/model'

interface ApiError {
  response?: {
    data?: {
      message?: string
    }
  }
}
import {
  deleteOtaPackagesId,
  getOtaPackages,
  getOtaPackagesId,
  postFilesOta,
  postOtaPackages,
  putOtaPackagesId,
} from '@/api/generated'
import { orvalAxios } from '@/api/orval-mutator'
import type { OTAPackage } from '../data/schema'

type CreateOTAPackagePayload = Omit<OtaCreateOTAPackageRequest, 'productId'> & {
  product_key: string
}

/**
 * Query key factory for OTA packages
 */
export const otaPackageKeys = {
  all: ['ota-packages'] as const,
  lists: () => [...otaPackageKeys.all, 'list'] as const,
  list: (filters?: {
    page?: number
    pageSize?: number
    productId?: string
    status?: string
    packageType?: string
    uploadType?: string
  }) => [...otaPackageKeys.lists(), filters] as const,
  details: () => [...otaPackageKeys.all, 'detail'] as const,
  detail: (id: string) => [...otaPackageKeys.details(), id] as const,
}

/**
 * Hook to fetch OTA packages list
 */
export function useOTAPackages(filters?: {
  page?: number
  pageSize?: number
  productId?: string
  status?: string
  packageType?: string
  uploadType?: string
  searchText?: string
}) {
  return useQuery<{ packages: OTAPackage[]; total: number }>({
    queryKey: otaPackageKeys.list(filters),
    queryFn: async () => {
      const response = (await getOtaPackages({
        page: filters?.page,
        pageSize: filters?.pageSize,
        productId: filters?.productId ? Number(filters.productId) : undefined,
        status: filters?.status,
        packageType: filters?.packageType,
        uploadType: filters?.uploadType,
        searchText: filters?.searchText,
      } as never))?.data ?? {}

      // Transform API response to match UI schema
      const packages: OTAPackage[] = (response.otaPackages || []).map(
        (pkg) => ({
          id: String(pkg.id || ''),
          uuid: pkg.uuid || '',
          packageName: pkg.packageName || '',
          version: pkg.version || '',
          packageType: (pkg.packageType || 'upgrade') as OTAPackage['packageType'],
          productId: pkg.productId?.toString(),
          productName: pkg.productName || '',
          description: pkg.description,
          fileSize: parseInt(String(pkg.fileSize || '0'), 10),
          fileUrl: pkg.fileUrl || '',
          checksum: pkg.checksum || '',
          status: (pkg.status || 'draft') as OTAPackage['status'],
          // Default values for deployment fields not in API yet
          deploymentProgress: 0,
          targetDeviceCount: 0,
          successCount: 0,
          failureCount: 0,
          createdAt: pkg.createdAt || new Date().toISOString(),
          updatedAt: pkg.updatedAt || new Date().toISOString(),
          deployedAt: pkg.releasedAt,
          createdBy: '',
        })
      )

      return { packages, total: response.total || 0 }
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

/**
 * Hook to fetch a single OTA package by ID
 */
export function useOTAPackage(id: string) {
  return useQuery({
    queryKey: otaPackageKeys.detail(id),
    queryFn: async () => {
      const response = await getOtaPackagesId(id as unknown as number)
      return response.data?.otaPackage
    },
    enabled: !!id,
  })
}

/**
 * Hook to create an OTA package
 */
export function useCreateOTAPackage() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (data: CreateOTAPackagePayload) => {
      return postOtaPackages(data)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: otaPackageKeys.lists() })
      toast.success('OTA package created successfully')
    },
    onError: (error) => {
      const message =
        (error as unknown as ApiError).response?.data?.message ||
        'Failed to create OTA package'
      toast.error(message)
    },
  })
}

export function useUploadOTAFile() {
  return useMutation({
    mutationFn: async (file: File) => {
      const response = await postFilesOta({ file })
      if (!response.data?.fileUrl) {
        throw new Error('The upload response did not include a file URL')
      }
      return response.data
    },
    onError: (error) => {
      const message =
        (error as unknown as ApiError).response?.data?.message ||
        'Failed to upload OTA package file'
      toast.error(message)
    },
  })
}

/**
 * Hook to update an OTA package
 */
export function useUpdateOTAPackage() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({
      id,
      data,
    }: {
      id: string
      data: OtaV1UpdateOTAPackageRequest
    }) => {
      return putOtaPackagesId(id as unknown as number, data as never)
    },
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: otaPackageKeys.lists() })
      queryClient.invalidateQueries({
        queryKey: otaPackageKeys.detail(variables.id),
      })
      toast.success('OTA package updated successfully')
    },
    onError: (error) => {
      const message =
        (error as unknown as ApiError).response?.data?.message ||
        'Failed to update OTA package'
      toast.error(message)
    },
  })
}

/**
 * Hook to delete an OTA package
 */
export function useDeleteOTAPackage() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => deleteOtaPackagesId(id as unknown as number),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: otaPackageKeys.lists() })
      toast.success('OTA package deleted successfully')
    },
    onError: (error) => {
      const message =
        (error as unknown as ApiError).response?.data?.message ||
        'Failed to delete OTA package'
      toast.error(message)
    },
  })
}

/**
 * Hook to delete multiple OTA packages
 */
export function useDeleteOTAPackages() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (ids: string[]) => {
      // Delete packages in parallel
      const deletePromises = ids.map((id) => deleteOtaPackagesId(id as unknown as number))
      await Promise.all(deletePromises)
      return { deletedCount: ids.length }
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: otaPackageKeys.lists() })
      toast.success(`${data.deletedCount} OTA packages deleted successfully`)
    },
    onError: (error) => {
      const message =
        (error as unknown as ApiError).response?.data?.message ||
        'Failed to delete OTA packages'
      toast.error(message)
    },
  })
}

interface DeployOTAPackagePayload {
  id: string
  deviceKeys: string[]
}

/**
 * Deploy an OTA package to the given target devices. Backed by
 * POST /v1/ota-packages/:id/deploy.
 */
export async function deployOtaPackage({
  id,
  deviceKeys,
}: DeployOTAPackagePayload) {
  return orvalAxios<{ code: number; message: string; data: { success: boolean } }>({
    method: 'POST',
    url: `/ota-packages/${id}/deploy`,
    data: { deviceKeys },
  })
}

/**
 * Hook to deploy an OTA package
 */
export function useDeployOTAPackage() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, deviceKeys }: DeployOTAPackagePayload) =>
      deployOtaPackage({ id, deviceKeys }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: otaPackageKeys.lists() })
      toast.success('OTA package deployment started')
    },
    onError: (error) => {
      const message =
        (error as unknown as ApiError).response?.data?.message ||
        'Failed to deploy OTA package'
      toast.error(message)
    },
  })
}
