import { AlertTriangle, RefreshCw } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { useOTAAnalytics } from '../hooks/use-ota-analytics'
import { AnalyticsSummaryCards } from './analytics-summary-cards'
import { DeploymentStatusChart } from './deployment-status-chart'
import { DeploymentTimelineChart } from './deployment-timeline-chart'
import { FirmwareDistributionChart } from './firmware-distribution-chart'
import { RecentActivityTable } from './recent-activity-table'

export function OTAAnalytics() {
  const { data, isLoading, isError, error, refetch } = useOTAAnalytics()

  if (isError) {
    return (
      <div className='flex h-[400px] items-center justify-center p-4'>
        <Card className='w-full max-w-md'>
          <CardContent className='space-y-4 pt-6'>
            <div className='flex flex-col items-center gap-2 text-center'>
              <AlertTriangle className='h-10 w-10 text-destructive' />
              <h3 className='text-lg font-semibold'>
                Failed to load analytics
              </h3>
              <p className='text-sm text-muted-foreground'>
                {error instanceof Error
                  ? error.message
                  : 'An error occurred while loading analytics data. Please try again.'}
              </p>
            </div>
            <Button
              onClick={() => refetch()}
              className='w-full'
              variant='outline'
            >
              <RefreshCw className='mr-2 h-4 w-4' />
              Retry
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className='flex flex-col gap-4'>
      {/* Summary Cards */}
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

      {/* Charts Row */}
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

      {/* Firmware Distribution */}
      <FirmwareDistributionChart
        data={data?.firmwareDistribution || []}
        isLoading={isLoading}
      />

      {/* Recent Activity */}
      <RecentActivityTable
        data={data?.recentActivity || []}
        isLoading={isLoading}
      />
    </div>
  )
}
