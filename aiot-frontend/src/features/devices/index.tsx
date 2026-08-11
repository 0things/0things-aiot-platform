import { useTranslation } from 'react-i18next'
import { DevicesDialogs } from './components/devices-dialogs'
import { DevicesPrimaryButtons } from './components/devices-primary-buttons'
import { DevicesProvider } from './components/devices-provider'
import { DevicesTable } from './components/devices-table'

export function Devices() {
  const { t } = useTranslation('deviceManagement')

  return (
    <DevicesProvider>
      <div className='flex flex-1 flex-col gap-4'>
        <div className='flex flex-wrap items-end justify-between gap-2'>
          <div>
            <h2 className='text-2xl font-bold tracking-tight'>
              {t('tabs.devices')}
            </h2>
            <p className='text-muted-foreground'>
              {t('pages.devices.description')}
            </p>
          </div>
          <DevicesPrimaryButtons />
        </div>
        <DevicesTable />
      </div>

      <DevicesDialogs />
    </DevicesProvider>
  )
}
