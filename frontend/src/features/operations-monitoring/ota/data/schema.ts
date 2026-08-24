import { z } from 'zod'

export const otaPackageStatusEnum = z.enum([
  'draft',
  'deploying',
  'success',
  'partial',
  'completed',
  'failed',
  'cancelled',
  'released',
  'archived',
])

export const otaPackageTypeEnum = z.enum([
  'upgrade',
  'security',
  'patch',
  'firmware',
  'config',
  'full',
])

export const otaPackageSchema = z.object({
  id: z.string(),
  packageName: z.string(),
  version: z.string(),
  packageType: otaPackageTypeEnum,
  productId: z.string().optional(),
  productName: z.string(),
  description: z.string().optional(),
  fileSize: z.number(),
  fileUrl: z.string(),
  checksum: z.string(),
  status: otaPackageStatusEnum,
  deploymentProgress: z.number().min(0).max(100),
  targetDeviceCount: z.number(),
  successCount: z.number(),
  failureCount: z.number(),
  createdAt: z.string(),
  updatedAt: z.string(),
  deployedAt: z.string().optional(),
  createdBy: z.string(),
})

export type OTAPackage = z.infer<typeof otaPackageSchema>
export type OTAPackageStatus = z.infer<typeof otaPackageStatusEnum>
export type OTAPackageType = z.infer<typeof otaPackageTypeEnum>

export const createPackageFormSchema = z.object({
  packageName: z
    .string()
    .min(1, 'packageForm.validation.packageNameRequired'),
  version: z
    .string()
    .min(1, 'packageForm.validation.versionRequired')
    .regex(
      /^\d+\.\d+\.\d+(\.\d+)?$/,
      'packageForm.validation.versionInvalid'
    ),
  packageType: otaPackageTypeEnum,
  productKey: z
    .string()
    .min(1, 'packageForm.validation.productNameRequired'),
  description: z.string().optional(),
  file: z.instanceof(File).optional(),
})

export type CreatePackageFormData = z.infer<typeof createPackageFormSchema>

// Edit form schema - same as create but packageName and version are not required for validation
export const editPackageFormSchema = z.object({
  packageName: z.string(),
  version: z.string(),
  packageType: otaPackageTypeEnum,
  productId: z.string(),
  description: z.string().optional(),
  file: z.instanceof(File).optional(),
})

export type EditPackageFormData = z.infer<typeof editPackageFormSchema>
