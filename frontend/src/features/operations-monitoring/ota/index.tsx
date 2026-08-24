import { useTranslation } from 'react-i18next'
import { OTAAnalytics as OTAAnalyticsComponent } from './components/ota-analytics'
import { OTAErrorBoundary } from './components/ota-error-boundary'
import { OTAPackagesDialogs } from './components/ota-packages-dialogs'
import { OTAPackagesPrimaryButtons } from './components/ota-packages-primary-buttons'
import { OTAPackagesProvider } from './components/ota-packages-provider'
import { OTAPackagesTable } from './components/ota-packages-table'

export function OTAPackages() {
  const { t } = useTranslation('ota')

  return (
    <OTAErrorBoundary>
      <OTAPackagesProvider>
        <div className='flex flex-1 flex-col gap-4'>
          <div className='flex flex-wrap items-end justify-between gap-2'>
            <div>
              <h2 className='text-2xl font-bold tracking-tight'>
                {t('tabs.packages')}
              </h2>
              <p className='text-muted-foreground'>
                {t('packageList.description')}
              </p>
            </div>
            <OTAPackagesPrimaryButtons />
          </div>
          <OTAPackagesTable />
        </div>
        <OTAPackagesDialogs />
      </OTAPackagesProvider>
    </OTAErrorBoundary>
  )
}

export function OTAAnalytics() {
  return (
    <OTAErrorBoundary>
      <OTAAnalyticsComponent />
    </OTAErrorBoundary>
  )
}
