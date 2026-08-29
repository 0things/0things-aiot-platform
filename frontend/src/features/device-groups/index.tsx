import { useTranslation } from 'react-i18next'
import { GroupsDialogs } from './components/groups-dialogs'
import { GroupsPrimaryButtons } from './components/groups-primary-buttons'
import { GroupsProvider } from './components/groups-provider'
import { GroupsTable } from './components/groups-table'

export function DeviceGroups() {
  const { t } = useTranslation('deviceGroup')

  return (
    <GroupsProvider>
      <div className='flex flex-1 flex-col gap-4'>
        <div className='flex flex-wrap items-end justify-between gap-2'>
          <div>
            <h2 className='text-2xl font-bold tracking-tight'>{t('title')}</h2>
            <p className='text-muted-foreground'>{t('description')}</p>
          </div>
          <GroupsPrimaryButtons />
        </div>
        <GroupsTable />
      </div>

      <GroupsDialogs />
    </GroupsProvider>
  )
}
