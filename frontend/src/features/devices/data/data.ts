import type { DeviceState } from './schema'

export const deviceStateStyles = new Map<DeviceState, string>([
  ['online', 'bg-teal-100/30 text-teal-900 dark:text-teal-200 border-teal-200'],
  ['offline', 'bg-neutral-300/40 border-neutral-300'],
  [
    'inactive',
    'bg-amber-100/40 text-amber-900 dark:text-amber-200 border-amber-300',
  ],
])

export const states = [
  {
    value: 'online',
    label: 'Online',
    variant: 'success' as const,
  },
  {
    value: 'offline',
    label: 'Offline',
    variant: 'secondary' as const,
  },
  {
    value: 'inactive',
    label: 'Inactive',
    variant: 'outline' as const,
  },
]
