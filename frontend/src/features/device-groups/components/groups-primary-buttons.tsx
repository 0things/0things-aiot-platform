import { Plus } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { useGroups } from './groups-provider'

export function GroupsPrimaryButtons() {
  const { t } = useTranslation('common')
  const { setOpen } = useGroups()

  return (
    <div className='flex gap-2'>
      <Button className='gap-2' onClick={() => setOpen('create')}>
        <Plus className='h-4 w-4' />
        <span>{t('create')}</span>
      </Button>
    </div>
  )
}
