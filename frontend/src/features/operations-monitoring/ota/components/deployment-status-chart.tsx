import { useTranslation } from 'react-i18next'
import {
  Cell,
  Pie,
  PieChart,
  ResponsiveContainer,
  Tooltip,
  type TooltipProps,
} from 'recharts'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import type { OTADeploymentStatus } from '../data/analytics-schema'

interface DeploymentStatusChartProps {
  data: OTADeploymentStatus[]
  isLoading?: boolean
}

// Mapped to theme.css chart tokens (1..5) so light/dark both work
const STATUS_TOKEN: Record<OTADeploymentStatus['status'], string> = {
  success: 'var(--chart-2)',
  failed: 'var(--chart-5)',
  in_progress: 'var(--chart-4)',
}

function StatusTooltip({
  active,
  payload,
}: TooltipProps<number, string> & {
  payload?: Array<{ payload: { name: string; value: number; percentage: number } }>
}) {
  if (!active || !payload?.length) return null
  const item = payload[0]
  const original = item.payload as {
    name: string
    value: number
    percentage: number
  }
  return (
    <div className='rounded-md border bg-popover px-3 py-2 text-xs text-popover-foreground shadow-sm'>
      <div className='font-medium'>{original.name}</div>
      <div className='tabular-nums text-muted-foreground'>
        {original.value} · {original.percentage.toFixed(1)}%
      </div>
    </div>
  )
}

export function DeploymentStatusChart({
  data,
  isLoading,
}: DeploymentStatusChartProps) {
  const { t } = useTranslation('ota')

  const chartData = data.map((item) => ({
    name: t(`analytics.deploymentStatus.${item.status}`),
    value: item.count,
    percentage: item.percentage,
    fill: STATUS_TOKEN[item.status],
  }))

  const total = chartData.reduce((acc, c) => acc + c.value, 0)
  const successItem = chartData.find((c) => c.name === t('analytics.deploymentStatus.success'))
  const successRate = successItem?.percentage ?? 0

  return (
    <Card className='h-full'>
      <CardHeader>
        <CardTitle>{t('analytics.charts.deploymentStatus')}</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className='flex h-[280px] items-center justify-center'>
            <Skeleton className='h-40 w-40 rounded-full' />
          </div>
        ) : (
          <div className='relative h-[280px] w-full'>
            <ResponsiveContainer width='100%' height='100%'>
              <PieChart>
                <Pie
                  data={chartData}
                  cx='50%'
                  cy='50%'
                  innerRadius={70}
                  outerRadius={100}
                  paddingAngle={2}
                  dataKey='value'
                  stroke='var(--background)'
                  strokeWidth={2}
                >
                  {chartData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.fill} />
                  ))}
                </Pie>
                <Tooltip content={<StatusTooltip />} />
              </PieChart>
            </ResponsiveContainer>
            <div className='pointer-events-none absolute inset-0 flex flex-col items-center justify-center'>
              <span className='text-3xl font-semibold tracking-tight tabular-nums'>
                {successRate.toFixed(0)}
                <span className='ml-0.5 text-base font-medium text-muted-foreground'>
                  %
                </span>
              </span>
              <span className='text-xs text-muted-foreground'>
                {t('analytics.deploymentStatus.success')} ·{' '}
                <span className='tabular-nums'>{total}</span>
              </span>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
