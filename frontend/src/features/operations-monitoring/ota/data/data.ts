import { type OTAPackageStatus, type OTAPackageType } from './schema'

export const statuses: {
  value: OTAPackageStatus
  label: string
  variant: 'default' | 'secondary' | 'destructive' | 'outline'
}[] = [
  {
    value: 'draft',
    label: 'Draft',
    variant: 'secondary',
  },
  {
    value: 'deploying',
    label: 'Deploying',
    variant: 'default',
  },
  {
    value: 'completed',
    label: 'Completed',
    variant: 'default',
  },
  {
    value: 'failed',
    label: 'Failed',
    variant: 'destructive',
  },
  {
    value: 'cancelled',
    label: 'Cancelled',
    variant: 'outline',
  },
  {
    value: 'released',
    label: 'Released',
    variant: 'default',
  },
  {
    value: 'archived',
    label: 'Archived',
    variant: 'outline',
  },
]

export const packageTypes: {
  value: OTAPackageType
  label: string
}[] = [
  {
    value: 'upgrade',
    label: 'Upgrade',
  },
  {
    value: 'security',
    label: 'Security',
  },
  {
    value: 'patch',
    label: 'Patch',
  },
  {
    value: 'firmware',
    label: 'Firmware',
  },
  {
    value: 'config',
    label: 'Config',
  },
  {
    value: 'full',
    label: 'Full',
  },
]
