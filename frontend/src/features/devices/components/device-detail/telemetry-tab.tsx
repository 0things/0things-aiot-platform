import { Route } from '@/routes/_authenticated/device-management/devices/$deviceKey'
import { Loader2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useDeviceTelemetry } from '@/features/devices/api/queries'

export function TelemetryTab() {
  const { t } = useTranslation('deviceManagement')
  const { deviceKey } = Route.useParams()
  const {
    data: telemetryResponse,
    isLoading,
    error,
  } = useDeviceTelemetry(deviceKey)

  const telemetry = telemetryResponse?.telemetry

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t('deviceDetail.telemetryTab.title')}</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className='flex min-h-[400px] items-center justify-center'>
            <Loader2 className='h-8 w-8 animate-spin text-muted-foreground' />
          </div>
        ) : error ? (
          <div className='flex min-h-[400px] flex-col items-center justify-center space-y-2'>
            <p className='text-lg text-muted-foreground'>
              {t('deviceDetail.telemetryTab.error')}
            </p>
          </div>
        ) : !telemetry ? (
          <div className='flex min-h-[400px] items-center justify-center'>
            <p className='text-lg text-muted-foreground'>
              {t('deviceDetail.telemetryTab.noData')}
            </p>
          </div>
        ) : (
          <div className='rounded-md bg-muted p-4'>
            <pre className='overflow-auto text-sm'>
              <code>{JSON.stringify(JSON.parse(telemetry), null, 2)}</code>
            </pre>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
