import { z } from 'zod'

export const groupTypeSchema = z.union([
  z.literal('manual'),
  z.literal('dynamic'),
])
export type GroupType = z.infer<typeof groupTypeSchema>

export const deviceGroupSchema = z.object({
  groupUuid: z.string(),
  name: z.string(),
  type: groupTypeSchema,
  description: z.string().optional(),
  rule: z.string().optional(),
  createdAt: z.string().optional(),
  updatedAt: z.string().optional(),
})

export type DeviceGroup = z.infer<typeof deviceGroupSchema>

export const groupFormSchema = z.object({
  name: z
    .string()
    .min(1, 'Group name is required')
    .max(128, 'Group name too long'),
  type: groupTypeSchema,
  rule: z.string().optional(),
  description: z.string().optional(),
})

export type GroupFormData = z.infer<typeof groupFormSchema>
