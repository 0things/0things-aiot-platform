import { Copy, Pencil } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import type { AiotBackendApiDeviceGroupV1DeviceGroup as DeviceGroupV1Group } from '@/api/generated/model'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

interface GroupInfoTabProps {
  group: DeviceGroupV1Group
  totalDevices: number
  activeDevices: number
  onlineDevices: number
  onEdit: () => void
}

export function GroupInfoTab({
  group,
  totalDevices,
  activeDevices,
  onlineDevices,
  onEdit,
}: GroupInfoTabProps) {
  const { t } = useTranslation('deviceGroup')

  const handleCopyUuid = () => {
    if (group.groupUuid) {
      navigator.clipboard.writeText(group.groupUuid)
      toast.success(t('copySuccess'))
    }
  }

  return (
    <Card>
      <CardHeader className='flex flex-row items-center justify-between border-b pb-4'>
        <CardTitle className='text-lg font-semibold'>{t('info')}</CardTitle>
        <Button
          variant='outline'
          size='sm'
          className='flex items-center gap-1.5'
          onClick={onEdit}
        >
          <Pencil className='h-3.5 w-3.5' />
          <span>{t('common:edit')}</span>
        </Button>
      </CardHeader>
      <CardContent className='pt-6'>
        <div className='grid grid-cols-1 gap-y-6 text-sm sm:grid-cols-2 lg:grid-cols-3'>
          <div>
            <p className='text-muted-foreground'>{t('name')}</p>
            <p className='mt-1 font-medium text-foreground'>
              {group.name || '-'}
            </p>
          </div>

          <div>
            <p className='text-muted-foreground'>{t('groupUuid')}</p>
            <div className='mt-1 flex items-center gap-2'>
              <span className='font-mono font-medium text-foreground'>
                {group.groupUuid || '-'}
              </span>
              {group.groupUuid && (
                <Button
                  variant='ghost'
                  size='icon'
                  className='h-6 w-6'
                  onClick={handleCopyUuid}
                  title={t('copySuccess')}
                >
                  <Copy className='h-3.5 w-3.5 text-primary' />
                </Button>
              )}
            </div>
          </div>

          <div>
            <p className='text-muted-foreground'>{t('totalDevices')}</p>
            <p className='mt-1 font-medium text-foreground'>{totalDevices}</p>
          </div>

          <div>
            <p className='text-muted-foreground'>{t('activeDevices')}</p>
            <p className='mt-1 font-medium text-foreground'>{activeDevices}</p>
          </div>

          <div>
            <p className='text-muted-foreground'>{t('currentOnline')}</p>
            <p className='mt-1 font-medium text-emerald-600 dark:text-emerald-400'>
              {onlineDevices}
            </p>
          </div>

          <div>
            <p className='text-muted-foreground'>{t('common:createdAt')}</p>
            <p className='mt-1 font-medium text-foreground'>
              {group.createdAt || '-'}
            </p>
          </div>

          <div>
            <p className='text-muted-foreground'>{t('type.label')}</p>
            <p className='mt-1 font-medium text-foreground'>
              {group.type === 'dynamic' ? t('type.dynamic') : t('type.manual')}
            </p>
          </div>

          {group.type === 'dynamic' && (
            <div>
              <p className='text-muted-foreground'>{t('rule')}</p>
              <p className='mt-1 font-mono text-xs font-medium text-foreground'>
                {group.rule || '-'}
              </p>
            </div>
          )}

          <div className='sm:col-span-2 lg:col-span-3'>
            <p className='text-muted-foreground'>{t('descriptionLabel')}</p>
            <p className='mt-1 text-foreground'>{group.description || '-'}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
