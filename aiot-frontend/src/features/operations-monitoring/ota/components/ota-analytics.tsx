import { AlertTriangle, RefreshCw } from 'lucide-react'
import { formatDistanceToNow } from 'date-fns'
import { enUS, zhCN } from 'date-fns/locale'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { useOTAAnalytics } from '../hooks/use-ota-analytics'
import { AnalyticsSummaryCards } from './analytics-summary-cards'
import { DeploymentStatusChart } from './deployment-status-chart'
import { DeploymentTimelineChart } from './deployment-timeline-chart'
import { FirmwareDistributionChart } from './firmware-distribution-chart'
import { RecentActivityTable } from './recent-activity-table'

export function OTAAnalytics() {
  const { data, isLoading, isError, error, refetch, isRefetching } =
    useOTAAnalytics()
  const { t, i18n } = useTranslation('operationsMonitoring')
  const locale = i18n.language?.startsWith('zh') ? zhCN : enUS

  if (isError) {
    return (
      <div className='flex h-[400px] items-center justify-center p-4'>
        <Card className='w-full max-w-md'>
          <CardContent className='space-y-4 pt-6'>
            <div className='flex flex-col items-center gap-2 text-center'>
              <AlertTriangle className='h-10 w-10 text-destructive' />
              <h3 className='text-lg font-semibold'>
                {t('ota.analytics.title')}
              </h3>
              <p className='text-sm text-muted-foreground'>
                {error instanceof Error
                  ? error.message
                  : t('ota.analytics.errorFallback')}
              </p>
            </div>
            <Button
              onClick={() => refetch()}
              className='w-full'
              variant='outline'
            >
              <RefreshCw className='mr-2 h-4 w-4' />
              {t('common:refresh')}
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  const lastUpdated = data?.meta?.lastUpdatedAt
    ? formatDistanceToNow(new Date(data.meta.lastUpdatedAt), {
        addSuffix: true,
        locale,
      })
    : null

  return (
    <div className='flex flex-col gap-6'>
      {/* Page header */}
      <div className='flex flex-wrap items-end justify-between gap-3'>
        <div>
          <h1 className='text-2xl font-semibold tracking-tight'>
            {t('ota.analytics.title')}
          </h1>
          <p className='mt-1 text-sm text-muted-foreground'>
            {t('ota.analytics.subtitle')}
          </p>
        </div>
        <div className='flex items-center gap-3'>
          {lastUpdated && (
            <span className='text-xs text-muted-foreground tabular-nums'>
              {t('ota.analytics.lastUpdated', { time: lastUpdated })}
            </span>
          )}
          <Button
            variant='outline'
            size='sm'
            onClick={() => refetch()}
            disabled={isRefetching}
          >
            <RefreshCw
              className={`mr-2 h-4 w-4 ${isRefetching ? 'animate-spin' : ''}`}
            />
            {t('common:refresh')}
          </Button>
        </div>
      </div>

      {/* Snapshot: KPIs + status */}
      <AnalyticsSummaryCards
        summary={
          data?.summary || {
            totalPackages: 0,
            activeDeployments: 0,
            successRate: 0,
            failedDeployments: 0,
          }
        }
        isLoading={isLoading}
      />

      <div className='grid gap-4 md:grid-cols-2 lg:grid-cols-3'>
        <DeploymentStatusChart
          data={data?.deploymentStatus || []}
          isLoading={isLoading}
        />
        <div className='lg:col-span-2'>
          <DeploymentTimelineChart
            data={data?.timeline || []}
            isLoading={isLoading}
          />
        </div>
      </div>

      <Separator />

      {/* History: firmware + activity */}
      <FirmwareDistributionChart
        data={data?.firmwareDistribution || []}
        isLoading={isLoading}
      />

      <RecentActivityTable
        data={data?.recentActivity || []}
        isLoading={isLoading}
      />
    </div>
  )
}