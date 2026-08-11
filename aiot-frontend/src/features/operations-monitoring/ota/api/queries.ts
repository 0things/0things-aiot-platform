import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { otaPackageServiceApi } from '@/api/clients'
import type {
  OtaV1CreateOTAPackageRequest,
  OtaV1UpdateOTAPackageRequest,
} from '@/api/generated/device-service'
import type { OTAPackage } from '../data/schema'

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
}) {
  return useQuery<OTAPackage[]>({
    queryKey: otaPackageKeys.list(filters),
    queryFn: async () => {
      const response =
        await otaPackageServiceApi.oTAPackageServiceListOTAPackages({
          page: filters?.page,
          pageSize: filters?.pageSize,
          productId: filters?.productId,
          status: filters?.status,
          packageType: filters?.packageType,
          uploadType: filters?.uploadType,
        })

      // Transform API response to match UI schema
      const packages: OTAPackage[] = (response.data.otaPackages || []).map(
        (pkg) => ({
          id: pkg.id || '',
          packageName: pkg.packageName || '',
          version: pkg.version || '',
          packageType: (pkg.packageType as any) || 'upgrade',
          productId: pkg.productId,
          productName: pkg.productName || '',
          description: pkg.description,
          fileSize: parseInt(pkg.fileSize || '0', 10),
          fileUrl: pkg.fileUrl || '',
          checksum: pkg.checksum || '',
          status: (pkg.status as any) || 'draft',
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

      return packages
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
      const response =
        await otaPackageServiceApi.oTAPackageServiceGetOTAPackage({ id })
      return response.data.otaPackage
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
    mutationFn: async (data: OtaV1CreateOTAPackageRequest) => {
      const response =
        await otaPackageServiceApi.oTAPackageServiceCreateOTAPackage({
          otaV1CreateOTAPackageRequest: data,
        })
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: otaPackageKeys.lists() })
      toast.success('OTA package created successfully')
    },
    onError: (error: any) => {
      const message =
        error.response?.data?.message || 'Failed to create OTA package'
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
      const response =
        await otaPackageServiceApi.oTAPackageServiceUpdateOTAPackage({
          id,
          otaV1UpdateOTAPackageRequest: data,
        })
      return response.data
    },
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: otaPackageKeys.lists() })
      queryClient.invalidateQueries({
        queryKey: otaPackageKeys.detail(variables.id),
      })
      toast.success('OTA package updated successfully')
    },
    onError: (error: any) => {
      const message =
        error.response?.data?.message || 'Failed to update OTA package'
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
    mutationFn: async (id: string) => {
      const response =
        await otaPackageServiceApi.oTAPackageServiceDeleteOTAPackage({ id })
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: otaPackageKeys.lists() })
      toast.success('OTA package deleted successfully')
    },
    onError: (error: any) => {
      const message =
        error.response?.data?.message || 'Failed to delete OTA package'
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
      const deletePromises = ids.map((id) =>
        otaPackageServiceApi.oTAPackageServiceDeleteOTAPackage({ id })
      )
      await Promise.all(deletePromises)
      return { deletedCount: ids.length }
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: otaPackageKeys.lists() })
      toast.success(`${data.deletedCount} OTA packages deleted successfully`)
    },
    onError: (error: any) => {
      const message =
        error.response?.data?.message || 'Failed to delete OTA packages'
      toast.error(message)
    },
  })
}
