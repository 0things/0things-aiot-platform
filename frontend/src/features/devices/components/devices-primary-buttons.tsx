import { Plus, Upload } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { useDevicesDialog } from './devices-provider'

export function DevicesPrimaryButtons() {
  const { t } = useTranslation('deviceManagement')
  const { setOpen } = useDevicesDialog()

  return (
    <div className='flex gap-2'>
      <Button variant='outline' onClick={() => setOpen('batch-upload')}>
        <Upload className='mr-2 size-4' />
        {t('devices.buttons.batchUpload')}
      </Button>
      <Button onClick={() => setOpen('create')}>
        <Plus className='mr-2 size-4' />
        {t('devices.buttons.createDevice')}
      </Button>
    </div>
  )
}
