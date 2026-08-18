import { useTranslation } from 'react-i18next'
import {
  Area,
  AreaChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  type TooltipProps,
  XAxis,
  YAxis,
} from 'recharts'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import type { OTATimelineDataPoint } from '../data/analytics-schema'

interface DeploymentTimelineChartProps {
  data: OTATimelineDataPoint[]
  isLoading?: boolean
}

interface TimelineRow {
  date: string
  deployments: number
  successes: number
  failures: number
}

function TimelineTooltip({
  active,
  payload,
  label,
  t,
}: TooltipProps<number, string> & {
  payload?: Array<{ payload: TimelineRow }>
  label?: string | number
  t: ReturnType<typeof useTranslation>['t']
}) {
  if (!active || !payload?.length) return null
  const row = payload[0].payload as TimelineRow
  return (
    <div className='rounded-md border bg-popover px-3 py-2 text-xs text-popover-foreground shadow-sm'>
      <div className='mb-1 font-medium tabular-nums'>{label}</div>
      <div className='flex items-center gap-2'>
        <span className='h-2 w-2 rounded-full bg-chart-1' />
        <span className='text-muted-foreground'>
          {t('ota.analytics.timeline.deployments')}
        </span>
        <span className='ml-auto tabular-nums'>{row.deployments}</span>
      </div>
      <div className='flex items-center gap-2'>
        <span className='h-2 w-2 rounded-full bg-chart-2' />
        <span className='text-muted-foreground'>
          {t('ota.analytics.timeline.successes')}
        </span>
        <span className='ml-auto tabular-nums'>{row.successes}</span>
      </div>
      <div className='flex items-center gap-2'>
        <span className='h-2 w-2 rounded-full bg-chart-5' />
        <span className='text-muted-foreground'>
          {t('ota.analytics.timeline.failures')}
        </span>
        <span className='ml-auto tabular-nums'>{row.failures}</span>
      </div>
    </div>
  )
}

export function DeploymentTimelineChart({
  data,
  isLoading,
}: DeploymentTimelineChartProps) {
  const { t } = useTranslation('operationsMonitoring')

  const chartData: TimelineRow[] = data.map((item) => ({
    date: new Date(item.date).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
    }),
    deployments: item.deployments,
    successes: item.successes,
    failures: item.failures,
  }))

  const series: Array<{
    key: keyof TimelineRow
    label: string
    token: string
  }> = [
    {
      key: 'deployments',
      label: t('ota.analytics.timeline.deployments'),
      token: 'var(--chart-1)',
    },
    {
      key: 'successes',
      label: t('ota.analytics.timeline.successes'),
      token: 'var(--chart-2)',
    },
    {
      key: 'failures',
      label: t('ota.analytics.timeline.failures'),
      token: 'var(--chart-5)',
    },
  ]

  return (
    <Card className='h-full'>
      <CardHeader className='flex flex-row items-center justify-between space-y-0'>
        <CardTitle>{t('ota.analytics.charts.deploymentTimeline')}</CardTitle>
        <div className='flex items-center gap-3 text-xs text-muted-foreground'>
          {series.map((s) => (
            <span key={s.key} className='flex items-center gap-1.5'>
              <span
                className='h-2 w-2 rounded-full'
                style={{ backgroundColor: s.token }}
              />
              {s.label}
            </span>
          ))}
        </div>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className='flex h-[280px] items-center justify-center'>
            <Skeleton className='h-full w-full rounded-md' />
          </div>
        ) : (
          <ResponsiveContainer width='100%' height={280}>
            <AreaChart
              data={chartData}
              margin={{ top: 10, right: 8, left: -16, bottom: 0 }}
            >
              <defs>
                {series.map((s) => (
                  <linearGradient
                    key={s.key}
                    id={`ota-${s.key}-gradient`}
                    x1='0'
                    y1='0'
                    x2='0'
                    y2='1'
                  >
                    <stop offset='0%' stopColor={s.token} stopOpacity={0.32} />
                    <stop offset='100%' stopColor={s.token} stopOpacity={0} />
                  </linearGradient>
                ))}
              </defs>
              <CartesianGrid
                strokeDasharray='2 4'
                stroke='var(--border)'
                vertical={false}
              />
              <XAxis
                dataKey='date'
                stroke='var(--muted-foreground)'
                fontSize={11}
                tickLine={false}
                axisLine={false}
                dy={4}
              />
              <YAxis
                stroke='var(--muted-foreground)'
                fontSize={11}
                tickLine={false}
                axisLine={false}
                allowDecimals={false}
              />
              <Tooltip content={<TimelineTooltip t={t} />} cursor={{ stroke: 'var(--border)', strokeDasharray: '2 4' }} />
              {series.map((s) => (
                <Area
                  key={s.key}
                  type='monotone'
                  dataKey={s.key}
                  stroke={s.token}
                  strokeWidth={1.75}
                  fill={`url(#ota-${s.key}-gradient)`}
                  fillOpacity={1}
                />
              ))}
            </AreaChart>
          </ResponsiveContainer>
        )}
      </CardContent>
    </Card>
  )
}
