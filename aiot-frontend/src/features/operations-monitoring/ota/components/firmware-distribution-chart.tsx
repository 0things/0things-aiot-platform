import { useTranslation } from 'react-i18next'
import {
  Bar,
  BarChart,
  CartesianGrid,
  Cell,
  ResponsiveContainer,
  Tooltip,
  type TooltipProps,
  XAxis,
  YAxis,
} from 'recharts'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import type { OTAFirmwareDistribution } from '../data/analytics-schema'

interface FirmwareDistributionChartProps {
  data: OTAFirmwareDistribution[]
  isLoading?: boolean
}

// Map each row to a theme token so light/dark both render consistently.
const FIRMWARE_TOKENS = [
  'var(--chart-1)',
  'var(--chart-2)',
  'var(--chart-3)',
  'var(--chart-4)',
  'var(--chart-5)',
]

interface FirmwareRow {
  version: string
  devices: number
  percentage: number
  fill: string
}

function FirmwareTooltip({
  active,
  payload,
}: TooltipProps<number, string>) {
  if (!active || !payload?.length) return null
  const row = payload[0].payload as FirmwareRow
  return (
    <div className='rounded-md border bg-popover px-3 py-2 text-xs text-popover-foreground shadow-sm'>
      <div className='font-mono font-medium'>v{row.version}</div>
      <div className='tabular-nums text-muted-foreground'>
        {row.devices} devices · {row.percentage.toFixed(1)}%
      </div>
    </div>
  )
}

export function FirmwareDistributionChart({
  data,
  isLoading,
}: FirmwareDistributionChartProps) {
  const { t } = useTranslation('operationsMonitoring')

  const chartData: FirmwareRow[] = data.map((item, index) => ({
    version: item.version,
    devices: item.deviceCount,
    percentage: item.percentage,
    fill: FIRMWARE_TOKENS[index % FIRMWARE_TOKENS.length],
  }))

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t('ota.analytics.charts.firmwareDistribution')}</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className='flex h-[300px] items-center justify-center'>
            <Skeleton className='h-full w-full rounded-md' />
          </div>
        ) : (
          <ResponsiveContainer width='100%' height={300}>
            <BarChart
              data={chartData}
              margin={{ top: 10, right: 8, left: -16, bottom: 0 }}
            >
              <CartesianGrid
                strokeDasharray='2 4'
                stroke='var(--border)'
                vertical={false}
              />
              <XAxis
                dataKey='version'
                stroke='var(--muted-foreground)'
                fontSize={11}
                tickLine={false}
                axisLine={false}
                tickFormatter={(v: string) => `v${v}`}
              />
              <YAxis
                stroke='var(--muted-foreground)'
                fontSize={11}
                tickLine={false}
                axisLine={false}
                allowDecimals={false}
              />
              <Tooltip
                content={<FirmwareTooltip />}
                cursor={{ fill: 'var(--muted)', opacity: 0.4 }}
              />
              <Bar dataKey='devices' radius={[6, 6, 0, 0]}>
                {chartData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.fill} />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        )}
      </CardContent>
    </Card>
  )
}