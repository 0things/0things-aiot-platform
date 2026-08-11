import { useTranslation } from 'react-i18next'
import {
  Bar,
  BarChart,
  CartesianGrid,
  Cell,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import type { OTAFirmwareDistribution } from '../data/analytics-schema'

interface FirmwareDistributionChartProps {
  data: OTAFirmwareDistribution[]
  isLoading?: boolean
}

const COLORS = ['#3b82f6', '#8b5cf6', '#ec4899', '#f59e0b', '#10b981']

export function FirmwareDistributionChart({
  data,
  isLoading,
}: FirmwareDistributionChartProps) {
  const { t } = useTranslation('operationsMonitoring')

  const chartData = data.map((item, index) => ({
    version: item.version,
    devices: item.deviceCount,
    percentage: item.percentage,
    fill: COLORS[index % COLORS.length],
  }))

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t('ota.analytics.charts.firmwareDistribution')}</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className='flex h-[300px] items-center justify-center'>
            <div className='h-full w-full animate-pulse rounded bg-muted' />
          </div>
        ) : (
          <ResponsiveContainer width='100%' height={300}>
            <BarChart data={chartData}>
              <CartesianGrid
                strokeDasharray='3 3'
                stroke='hsl(var(--border))'
              />
              <XAxis
                dataKey='version'
                stroke='hsl(var(--muted-foreground))'
                fontSize={12}
              />
              <YAxis stroke='hsl(var(--muted-foreground))' fontSize={12} />
              <Tooltip
                formatter={(value) => `${value ?? 0} devices`}
                contentStyle={{
                  backgroundColor: 'hsl(var(--background))',
                  border: '1px solid hsl(var(--border))',
                  borderRadius: '6px',
                }}
              />
              <Bar dataKey='devices' radius={[8, 8, 0, 0]}>
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
