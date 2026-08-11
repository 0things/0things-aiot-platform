import { useTranslation } from 'react-i18next'
import {
  Area,
  AreaChart,
  CartesianGrid,
  Legend,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import type { OTATimelineDataPoint } from '../data/analytics-schema'

interface DeploymentTimelineChartProps {
  data: OTATimelineDataPoint[]
  isLoading?: boolean
}

export function DeploymentTimelineChart({
  data,
  isLoading,
}: DeploymentTimelineChartProps) {
  const { t } = useTranslation('operationsMonitoring')

  const chartData = data.map((item) => ({
    date: new Date(item.date).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
    }),
    [t('ota.analytics.timeline.deployments')]: item.deployments,
    [t('ota.analytics.timeline.successes')]: item.successes,
    [t('ota.analytics.timeline.failures')]: item.failures,
  }))

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t('ota.analytics.charts.deploymentTimeline')}</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className='flex h-[300px] items-center justify-center'>
            <div className='h-full w-full animate-pulse rounded bg-muted' />
          </div>
        ) : (
          <ResponsiveContainer width='100%' height={300}>
            <AreaChart data={chartData}>
              <defs>
                <linearGradient
                  id='colorDeployments'
                  x1='0'
                  y1='0'
                  x2='0'
                  y2='1'
                >
                  <stop offset='5%' stopColor='#3b82f6' stopOpacity={0.8} />
                  <stop offset='95%' stopColor='#3b82f6' stopOpacity={0} />
                </linearGradient>
                <linearGradient id='colorSuccesses' x1='0' y1='0' x2='0' y2='1'>
                  <stop offset='5%' stopColor='#22c55e' stopOpacity={0.8} />
                  <stop offset='95%' stopColor='#22c55e' stopOpacity={0} />
                </linearGradient>
                <linearGradient id='colorFailures' x1='0' y1='0' x2='0' y2='1'>
                  <stop offset='5%' stopColor='#ef4444' stopOpacity={0.8} />
                  <stop offset='95%' stopColor='#ef4444' stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid
                strokeDasharray='3 3'
                stroke='hsl(var(--border))'
              />
              <XAxis
                dataKey='date'
                stroke='hsl(var(--muted-foreground))'
                fontSize={12}
              />
              <YAxis stroke='hsl(var(--muted-foreground))' fontSize={12} />
              <Tooltip
                contentStyle={{
                  backgroundColor: 'hsl(var(--background))',
                  border: '1px solid hsl(var(--border))',
                  borderRadius: '6px',
                }}
              />
              <Legend />
              <Area
                type='monotone'
                dataKey={t('ota.analytics.timeline.deployments')}
                stroke='#3b82f6'
                fillOpacity={1}
                fill='url(#colorDeployments)'
              />
              <Area
                type='monotone'
                dataKey={t('ota.analytics.timeline.successes')}
                stroke='#22c55e'
                fillOpacity={1}
                fill='url(#colorSuccesses)'
              />
              <Area
                type='monotone'
                dataKey={t('ota.analytics.timeline.failures')}
                stroke='#ef4444'
                fillOpacity={1}
                fill='url(#colorFailures)'
              />
            </AreaChart>
          </ResponsiveContainer>
        )}
      </CardContent>
    </Card>
  )
}
