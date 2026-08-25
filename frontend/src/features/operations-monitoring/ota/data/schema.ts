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
  uuid: z.string(),
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

export const otaPackageFileExtensions = [
  '.tar.gz',
  '.tar.xz',
  '.gzip',
  '.pack',
  '.bin',
  '.dav',
  '.tar',
  '.gz',
  '.zip',
  '.apk',
] as const

function hasSupportedOTAFileExtension(filename: string) {
  const normalizedName = filename.toLowerCase()
  return otaPackageFileExtensions.some((extension) =>
    normalizedName.endsWith(extension)
  )
}

export const createPackageFormSchema = z
  .object({
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
    productKey: z.string().min(1, 'packageForm.validation.productNameRequired'),
    description: z.string().optional(),
    uploadMethod: z.enum(['file', 'url']),
    file: z.instanceof(File).optional(),
    fileUrl: z.string().optional(),
  })
  .superRefine((data, ctx) => {
    if (data.uploadMethod === 'file') {
      if (!data.file) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ['file'],
          message: 'packageForm.validation.fileOrUrlRequired',
        })
      } else if (data.file.size > 100 * 1024 * 1024) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ['file'],
          message: 'packageForm.validation.fileTooLarge',
        })
      } else if (!hasSupportedOTAFileExtension(data.file.name)) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ['file'],
          message: 'packageForm.validation.invalidFileType',
        })
      }
      return
    }

    if (!data.fileUrl) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ['fileUrl'],
        message: 'packageForm.validation.fileOrUrlRequired',
      })
    } else if (!z.string().url().safeParse(data.fileUrl).success) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ['fileUrl'],
        message: 'packageForm.validation.invalidUrl',
      })
    } else if (!hasSupportedOTAFileExtension(new URL(data.fileUrl).pathname)) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ['fileUrl'],
        message: 'packageForm.validation.invalidFileType',
      })
    }
  })

export type CreatePackageFormData = z.infer<typeof createPackageFormSchema>

// Edit form schema - same as create but packageName and version are not required for validation
export const editPackageFormSchema = z.object({
  packageName: z.string(),
  version: z.string(),
  packageType: otaPackageTypeEnum,
  productId: z.string(),
  description: z.string().optional(),
  file: z
    .instanceof(File)
    .optional()
    .refine(
      (file) => !file || hasSupportedOTAFileExtension(file.name),
      'packageForm.validation.invalidFileType'
    ),
})

export type EditPackageFormData = z.infer<typeof editPackageFormSchema>
