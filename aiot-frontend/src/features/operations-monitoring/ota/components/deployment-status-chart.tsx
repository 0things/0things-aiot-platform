import { useTranslation } from 'react-i18next'
import {
  Cell,
  Pie,
  PieChart,
  ResponsiveContainer,
  Legend,
  Tooltip,
} from 'recharts'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import type { OTADeploymentStatus } from '../data/analytics-schema'

interface DeploymentStatusChartProps {
  data: OTADeploymentStatus[]
  isLoading?: boolean
}

const COLORS = {
  success: '#22c55e',
  failed: '#ef4444',
  in_progress: '#f59e0b',
}

export function DeploymentStatusChart({
  data,
  isLoading,
}: DeploymentStatusChartProps) {
  const { t } = useTranslation('operationsMonitoring')

  const chartData = data.map((item) => ({
    name: t(`ota.analytics.deploymentStatus.${item.status}`),
    value: item.count,
    percentage: item.percentage,
    fill: COLORS[item.status],
  }))

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t('ota.analytics.charts.deploymentStatus')}</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className='flex h-[300px] items-center justify-center'>
            <div className='h-40 w-40 animate-pulse rounded-full bg-muted' />
          </div>
        ) : (
          <ResponsiveContainer width='100%' height={300}>
            <PieChart>
              <Pie
                data={chartData}
                cx='50%'
                cy='50%'
                labelLine={false}
                label={(entry) => `${Math.round((entry.percent ?? 0) * 100)}%`}
                outerRadius={100}
                fill='#8884d8'
                dataKey='value'
              >
                {chartData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.fill} />
                ))}
              </Pie>
              <Tooltip
                formatter={(value) => value ?? 0}
                contentStyle={{
                  backgroundColor: 'hsl(var(--background))',
                  border: '1px solid hsl(var(--border))',
                  borderRadius: '6px',
                }}
              />
              <Legend />
            </PieChart>
          </ResponsiveContainer>
        )}
      </CardContent>
    </Card>
  )
}
