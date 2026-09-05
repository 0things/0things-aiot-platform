import { ArrowLeft } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import type { AiotBackendApiDeviceGroupV1DeviceGroup as DeviceGroupV1Group } from '@/api/generated/model'
import { Button } from '@/components/ui/button'
import { CopyButton } from '@/components/copy-button'

interface GroupHeaderProps {
  group: DeviceGroupV1Group
  totalDevices: number
  onBack: () => void
}

export function GroupHeader({ group, totalDevices, onBack }: GroupHeaderProps) {
  const { t } = useTranslation('deviceGroup')

  return (
    <div className='mb-6 space-y-4 rounded-lg border bg-card p-4 text-card-foreground shadow-sm'>
      <div className='flex items-center gap-3 border-b pb-3'>
        <Button
          variant='ghost'
          size='icon'
          onClick={onBack}
          aria-label='Back to groups list'
        >
          <ArrowLeft className='h-4 w-4' />
        </Button>
        <h1 className='text-2xl font-bold tracking-tight'>
          {group.name || '-'}
        </h1>
      </div>

      <div className='grid grid-cols-1 gap-y-3 text-sm sm:grid-cols-2 lg:grid-cols-4'>
        <div className='flex items-center gap-2'>
          <span className='text-muted-foreground'>{t('groupUuid')}: </span>
          <span className='font-mono text-foreground'>
            {group.groupUuid || '-'}
          </span>
          {group.groupUuid && <CopyButton value={group.groupUuid} />}
        </div>

        <div>
          <span className='text-muted-foreground'>{t('totalDevices')}: </span>
          <span className='font-semibold text-foreground'>{totalDevices}</span>
        </div>
      </div>
    </div>
  )
}
