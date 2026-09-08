import { z } from 'zod'
import { TSL_TEMPLATES } from './constants'
import type { DataType, TSLModel } from './types'

// Zod Schema for TSL validation - 使用 lazy 支持递归
export const dataTypeSchema: z.ZodType<DataType> = z.lazy(() =>
  z.object({
    type: z.enum([
      'int',
      'float',
      'double',
      'bool',
      'string',
      'enum',
      'struct',
      'array',
    ]),
    specs: z.any().optional(), // specs 结构在不同类型下差异很大，使用 any 更灵活
  })
)

// 基础验证 schema（不包含自定义错误信息）
export const propertySchema = z.object({
  identifier: z.string().min(1),
  name: z.string().min(1),
  accessMode: z.enum(['r', 'rw']),
  required: z.boolean(),
  dataType: dataTypeSchema,
  desc: z.string().optional(),
  description: z.string().optional(),
})

export const serviceParamSchema = z.object({
  identifier: z.string().min(1),
  name: z.string().min(1),
  dataType: dataTypeSchema,
})

export const serviceSchema = z.object({
  identifier: z.string().min(1),
  name: z.string().min(1),
  required: z.boolean(),
  callType: z.enum(['async', 'sync']),
  inputData: z.array(serviceParamSchema),
  outputData: z.array(serviceParamSchema),
  desc: z.string().optional(),
  description: z.string().optional(),
})

export const eventParamSchema = z.object({
  identifier: z.string().min(1),
  name: z.string().min(1),
  dataType: dataTypeSchema,
})

export const eventSchema = z.object({
  identifier: z.string().min(1),
  name: z.string().min(1),
  type: z.enum(['info', 'alert', 'error']),
  required: z.boolean(),
  outputData: z.array(eventParamSchema),
  desc: z.string().optional(),
  description: z.string().optional(),
})

export const tslSchema = z.object({
  schema: z.string().optional().default('schema.json'),
  profile: z
    .object({
      productKey: z.string().min(1),
    })
    .optional(),
  properties: z.array(propertySchema).optional().default([]),
  events: z.array(eventSchema).optional().default([]),
  services: z.array(serviceSchema).optional().default([]),
})

export function normalizeTSLModel(
  model: Partial<TSLModel> | null | undefined
): TSLModel {
  if (!model) {
    return TSL_TEMPLATES.empty.value as unknown as TSLModel
  }
  return {
    schema: model.schema || 'schema.json',
    version: model.version || '1.0.0',
    profile: {
      productKey: model.profile?.productKey || '',
    },
    properties: Array.isArray(model.properties) ? model.properties : [],
    events: Array.isArray(model.events) ? model.events : [],
    services: Array.isArray(model.services) ? model.services : [],
  }
}
