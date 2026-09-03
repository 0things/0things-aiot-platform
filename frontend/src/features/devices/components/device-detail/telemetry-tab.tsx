import { useMemo } from 'react'
import { format } from 'date-fns'
import { Route } from '@/routes/_authenticated/device-management/devices/$deviceKey'
import { Activity, Database, Loader2, Radio, RotateCw } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { useDeviceShadow } from '@/features/devices/api/telemetry'
import { TelemetryHistoryChart } from './telemetry/telemetry-history-chart'

export function TelemetryTab() {
  const { deviceKey } = Route.useParams()
  const { t } = useTranslation('deviceManagement')

  const {
    data: shadow,
    isLoading: isShadowLoading,
    isRefetching,
    refetch,
  } = useDeviceShadow(deviceKey)

  const lastSeen = shadow?.updatedAt
    ? format(new Date(shadow.updatedAt), 'yyyy-MM-dd HH:mm:ss')
    : null

  // Build the property list from reported values without imposing units.
  const propertyList = useMemo(() => {
    const rawAttrs =
      (shadow?.state?.reported as Record<string, unknown>) ||
      (shadow?.reported as Record<string, unknown>) ||
      (shadow as unknown as { attributes?: Record<string, unknown> })
        ?.attributes ||
      {}
    const keys = Object.keys(rawAttrs)
    if (keys.length > 0) {
      return keys.map((key) => {
        const val = rawAttrs[key]
        return {
          identifier: key,
          name: key,
          value: val,
        }
      })
    }
    return []
  }, [shadow])

  return (
    <div className='flex flex-col gap-6'>
      {/* Live telemetry property cards. */}
      <Card className='shadow-xs'>
        <CardHeader className='pb-4'>
          <div className='flex flex-wrap items-center justify-between gap-4'>
            <div>
              <div className='flex items-center gap-2'>
                <Radio className='size-5 animate-pulse text-emerald-500' />
                <CardTitle className='text-lg font-semibold'>
                  {t('deviceDetail.telemetryTab.snapshotTitle')}
                </CardTitle>
              </div>
              <CardDescription className='mt-1'>
                {t('deviceDetail.telemetryTab.snapshotDescription')}
              </CardDescription>
            </div>

            <div className='flex items-center gap-3'>
              {lastSeen && (
                <span className='text-xs text-muted-foreground'>
                  {t('deviceDetail.telemetryTab.lastSeen', { time: lastSeen })}
                </span>
              )}
              <Button
                size='sm'
                variant='outline'
                onClick={() => refetch()}
                disabled={isShadowLoading || isRefetching}
                className='h-8 text-xs'
              >
                <RotateCw
                  className={`mr-1.5 size-3.5 ${isRefetching ? 'animate-spin' : ''}`}
                />
                {t('deviceDetail.telemetryTab.refreshShadow')}
              </Button>
            </div>
          </div>
        </CardHeader>

        <CardContent>
          {isShadowLoading && !shadow ? (
            <div className='flex h-32 items-center justify-center'>
              <Loader2 className='size-6 animate-spin text-muted-foreground' />
            </div>
          ) : propertyList.length === 0 ? (
            <div className='flex h-32 flex-col items-center justify-center rounded-md border border-dashed text-muted-foreground'>
              <Database className='size-8 text-muted-foreground/40' />
              <p className='mt-2 text-sm'>
                {t('deviceDetail.telemetryTab.noShadowData')}
              </p>
            </div>
          ) : (
            <div className='grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4'>
              {propertyList.map((prop) => {
                const isBool = typeof prop.value === 'boolean'
                const isObj =
                  typeof prop.value === 'object' && prop.value !== null

                return (
                  <div
                    key={prop.identifier}
                    className='relative flex flex-col justify-between rounded-xl border bg-card/60 p-4 transition-all hover:border-primary/40 hover:shadow-xs'
                  >
                    <div className='flex items-center justify-between'>
                      <span className='text-xs font-medium text-muted-foreground'>
                        {prop.name}
                      </span>
                      <Badge
                        variant='outline'
                        className='font-mono text-[10px] text-muted-foreground'
                      >
                        {prop.identifier}
                      </Badge>
                    </div>

                    <div className='my-3 flex items-baseline gap-1.5'>
                      {isBool ? (
                        <Badge
                          variant={prop.value ? 'default' : 'secondary'}
                          className='px-2 py-0.5 text-xs'
                        >
                          {prop.value ? 'true' : 'false'}
                        </Badge>
                      ) : isObj ? (
                        <pre className='max-h-16 overflow-auto font-mono text-xs text-muted-foreground'>
                          {JSON.stringify(prop.value)}
                        </pre>
                      ) : (
                        <span className='text-2xl font-bold tracking-tight text-foreground'>
                          {String(prop.value)}
                        </span>
                      )}
                    </div>

                    <div className='flex items-center justify-between text-[11px] text-muted-foreground/70'>
                      <div className='flex items-center gap-1'>
                        <Activity className='size-3 text-primary/70' />
                        <span>{t('deviceDetail.telemetryTab.liveReport')}</span>
                      </div>
                    </div>
                  </div>
                )
              })}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Historical time-series chart. */}
      <TelemetryHistoryChart
        deviceKey={deviceKey}
        availableProperties={propertyList}
      />
    </div>
  )
}
