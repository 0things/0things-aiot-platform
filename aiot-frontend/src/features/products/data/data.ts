import type { ProductStatus } from './schema'

export const productStatusStyles = new Map<ProductStatus, string>([
  ['active', 'bg-teal-100/30 text-teal-900 dark:text-teal-200 border-teal-200'],
  ['inactive', 'bg-neutral-300/40 border-neutral-300'],
  ['archived', 'bg-sky-200/40 text-sky-900 dark:text-sky-100 border-sky-300'],
])

export const statuses = [
  {
    value: 'active',
    label: 'Active',
    variant: 'success' as const,
  },
  {
    value: 'inactive',
    label: 'Inactive',
    variant: 'secondary' as const,
  },
  {
    value: 'archived',
    label: 'Archived',
    variant: 'outline' as const,
  },
]

export const categories = [
  {
    value: 'sensor',
    label: 'Sensor',
  },
  {
    value: 'actuator',
    label: 'Actuator',
  },
  {
    value: 'gateway',
    label: 'Gateway',
  },
  {
    value: 'controller',
    label: 'Controller',
  },
  {
    value: 'display',
    label: 'Display',
  },
  {
    value: 'camera',
    label: 'Camera',
  },
  {
    value: 'other',
    label: 'Other',
  },
]

export const nodeTypes = [
  {
    value: 'direct',
    label: '直连设备',
  },
  {
    value: 'gateway',
    label: '网关设备',
  },
  {
    value: 'gateway-sub',
    label: '网关子设备',
  },
]

export const connectivityMethods = [
  {
    value: 'wifi',
    label: 'Wi-Fi',
  },
  {
    value: 'cellular',
    label: '蜂窝(2G/3G/4G/5G)',
  },
  {
    value: 'ethernet',
    label: '以太网',
  },
  {
    value: 'other',
    label: '其他',
  },
]

// Access protocols for direct and gateway devices
export const directGatewayProtocols = [
  {
    value: 'http',
    label: 'HTTP',
  },
  {
    value: 'mqtt',
    label: 'MQTT',
  },
  {
    value: 'other',
    label: '其他',
  },
]

// Access protocols for gateway sub-devices
export const gatewaySubProtocols = [
  {
    value: 'custom',
    label: '自定义',
  },
  {
    value: 'modbus',
    label: 'Modbus',
  },
  {
    value: 'opc-ua',
    label: 'OPC UA',
  },
  {
    value: 'zigbee',
    label: 'ZigBee',
  },
  {
    value: 'ble',
    label: 'BLE',
  },
]
