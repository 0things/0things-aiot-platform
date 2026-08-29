import { z } from 'zod'
import { protocolLabels } from './protocols'

// Product status schema
const productStatusSchema = z.union([
  z.literal('active'),
  z.literal('inactive'),
  z.literal('archived'),
])
export type ProductStatus = z.infer<typeof productStatusSchema>

// Node type schema
const nodeTypeSchema = z.union([
  z.literal('direct'),
  z.literal('gateway'),
  z.literal('gateway-sub'),
])
export type NodeType = z.infer<typeof nodeTypeSchema>

// Connectivity method schema
const connectivityMethodSchema = z.union([
  z.literal('wifi'),
  z.literal('cellular'),
  z.literal('ethernet'),
  z.literal('other'),
])
export type ConnectivityMethod = z.infer<typeof connectivityMethodSchema>

// Access protocol schemas - different protocols for different node types
// For direct and gateway devices
const directGatewayProtocolSchema = z.union([
  z.literal('default'),
  z.literal('http'),
  z.literal('mqtt'),
  z.literal('coap'),
  z.literal('tcp'),
  z.literal('udp'),
  z.literal('gb28181'),
])

// For gateway sub-devices
const gatewaySubProtocolSchema = z.union([
  z.literal('custom'),
  z.literal('modbus'),
  z.literal('opc-ua'),
  z.literal('zigbee'),
  z.literal('ble'),
  z.literal('jt808'),
  z.literal('jt1078'),
])

// Combined access protocol schema (union of all protocols)
const accessProtocolSchema = z.union([
  directGatewayProtocolSchema,
  gatewaySubProtocolSchema,
])
export type AccessProtocol = z.infer<typeof accessProtocolSchema>
export type DirectGatewayProtocol = z.infer<typeof directGatewayProtocolSchema>
export type GatewaySubProtocol = z.infer<typeof gatewaySubProtocolSchema>

// Full product schema (from API)
export const productSchema = z.object({
  id: z.string(),
  productKey: z.string(),
  name: z.string(),
  description: z.string().optional(),
  categoryId: z.number().optional(),
  status: productStatusSchema,
  nodeType: nodeTypeSchema.optional(),
  connectivityMethod: connectivityMethodSchema.optional(),
  accessProtocol: accessProtocolSchema.optional(),
  deviceCount: z.number().optional(),
  createdAt: z.string(),
  updatedAt: z.string(),
})

export type Product = z.infer<typeof productSchema>

// Schema for product form validation
export const productFormSchema = z
  .object({
    name: z.string().min(1, 'Name is required').max(100, 'Name too long'),
    description: z.string().max(500, 'Description too long').optional(),
    categoryId: z.string().min(1, 'Category is required'),
    nodeType: nodeTypeSchema,
    connectivityMethod: connectivityMethodSchema.optional(),
    accessProtocol: accessProtocolSchema,
  })
  .superRefine((data, ctx) => {
    // Connectivity method is required for direct and gateway devices
    if (
      (data.nodeType === 'direct' || data.nodeType === 'gateway') &&
      !data.connectivityMethod
    ) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: 'Connectivity method is required for this node type',
        path: ['connectivityMethod'],
      })
    }

    // Validate access protocol matches node type
    if (data.nodeType === 'direct' || data.nodeType === 'gateway') {
      // Direct and gateway devices should use http, mqtt, or other
      if (
        data.accessProtocol !== 'default' &&
        data.accessProtocol !== 'http' &&
        data.accessProtocol !== 'mqtt' &&
        data.accessProtocol !== 'coap' &&
        data.accessProtocol !== 'tcp' &&
        data.accessProtocol !== 'udp' &&
        data.accessProtocol !== 'gb28181'
      ) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message:
            'Direct and gateway devices must use a supported transport protocol',
          path: ['accessProtocol'],
        })
      }
    } else if (data.nodeType === 'gateway-sub') {
      // Gateway sub-devices should use gateway protocols
      if (
        data.accessProtocol === 'default' ||
        data.accessProtocol === 'http' ||
        data.accessProtocol === 'mqtt' ||
        data.accessProtocol === 'coap' ||
        data.accessProtocol === 'tcp' ||
        data.accessProtocol === 'udp' ||
        data.accessProtocol === 'gb28181'
      ) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message:
            'Gateway sub-devices must use gateway protocols (Custom, Modbus, OPC UA, ZigBee, BLE, JT808, JT1078)',
          path: ['accessProtocol'],
        })
      }
    }
  })

export type ProductFormData = z.infer<typeof productFormSchema>

// Label mapping helpers
export const nodeTypeLabels: Record<NodeType, string> = {
  direct: '直连设备',
  gateway: '网关设备',
  'gateway-sub': '网关子设备',
}

export const connectivityMethodLabels: Record<ConnectivityMethod, string> = {
  wifi: 'Wi-Fi',
  cellular: '蜂窝(2G/3G/4G/5G)',
  ethernet: '以太网',
  other: '其他',
}

export const accessProtocolLabels: Record<AccessProtocol, string> =
  protocolLabels as unknown as Record<AccessProtocol, string>
