import { useMemo, useState } from 'react'
import { format } from 'date-fns'
import {
  Activity,
  ArrowDown,
  ArrowUp,
  Calendar,
  Layers,
  Loader2,
  RotateCw,
  TrendingUp,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import {
  Area,
  AreaChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useTelemetryHistory } from '@/features/devices/api/telemetry'

interface TelemetryHistoryChartProps {
  deviceKey: string
  availableProperties: Array<{
    identifier: string
    name?: string
    unit?: string
  }>
}

type TimeRangeKey = '1h' | '6h' | '24h' | '7d'

const TIME_RANGES: Record<TimeRangeKey, { durationMs: number }> = {
  '1h': { durationMs: 1 * 60 * 60 * 1000 },
  '6h': { durationMs: 6 * 60 * 60 * 1000 },
  '24h': { durationMs: 24 * 60 * 60 * 1000 },
  '7d': { durationMs: 7 * 24 * 60 * 60 * 1000 },
}

export function TelemetryHistoryChart({
  deviceKey,
  availableProperties,
}: TelemetryHistoryChartProps) {
  const { t } = useTranslation('deviceManagement')
  const [selectedProperty, setSelectedProperty] = useState<string>(
    availableProperties[0]?.identifier || 'temperature'
  )
  const [timeRange, setTimeRange] = useState<TimeRangeKey>('24h')

  const activeProperty = availableProperties.some(
    (p) => p.identifier === selectedProperty
  )
    ? selectedProperty
    : availableProperties[0]?.identifier || selectedProperty

  const {
    data: points = [],
    isLoading,
    isRefetching,
    refetch,
  } = useTelemetryHistory({
    deviceKey,
    property: activeProperty,
    durationMs: TIME_RANGES[timeRange].durationMs,
    limit: 500,
  })

  const handleManualRefresh = () => {
    refetch()
  }

  const handleRangeChange = (range: TimeRangeKey) => {
    setTimeRange(range)
  }

  // Format numeric points and calculate aggregate statistics.
  const { chartData, stats } = useMemo(() => {
    const numericPoints = points
      .map((p) => {
        const val =
          typeof p.value === 'number'
            ? p.value
            : typeof p.value === 'string'
              ? Number(p.value)
              : Number.NaN
        if (!Number.isFinite(val)) return null
        return {
          timestamp: p.timestamp,
          formattedTime: format(new Date(p.timestamp), 'HH:mm:ss'),
          fullTime: format(new Date(p.timestamp), 'yyyy-MM-dd HH:mm:ss'),
          value: val,
        }
      })
      .filter((point): point is NonNullable<typeof point> => point !== null)
      .sort((a, b) => a.timestamp - b.timestamp)

    if (numericPoints.length === 0) {
      return {
        chartData: [],
        stats: {
          min: '--',
          max: '--',
          avg: '--',
          latest: '--',
          count: 0,
        },
      }
    }

    const values = numericPoints.map((p) => p.value)
    const min = Math.min(...values)
    const max = Math.max(...values)
    const avg = values.reduce((sum, v) => sum + v, 0) / values.length
    const latest = values[values.length - 1]

    return {
      chartData: numericPoints,
      stats: {
        min: Number(min.toFixed(2)),
        max: Number(max.toFixed(2)),
        avg: Number(avg.toFixed(2)),
        latest: Number(latest.toFixed(2)),
        count: values.length,
      },
    }
  }, [points])

  const currentPropMeta = availableProperties.find(
    (p) => p.identifier === activeProperty
  )
  const unit = currentPropMeta?.unit || ''
  const propName = currentPropMeta?.name || activeProperty

  return (
    <Card className='shadow-xs'>
      <CardHeader className='pb-4'>
        <div className='flex flex-wrap items-center justify-between gap-4'>
          <div>
            <div className='flex items-center gap-2'>
              <TrendingUp className='size-5 text-primary' />
              <CardTitle className='text-lg font-semibold'>
                {t('deviceDetail.telemetryTab.historyTitle', {
                  property: propName,
                })}
              </CardTitle>
            </div>
            <CardDescription className='mt-1'>
              {t('deviceDetail.telemetryTab.historyDescription')}
            </CardDescription>
          </div>

          {/* Chart controls. */}
          <div className='flex flex-wrap items-center gap-2'>
            {/* Property selector. */}
            <Select
              value={activeProperty}
              onValueChange={(val) => setSelectedProperty(val)}
            >
              <SelectTrigger className='h-8 w-36 text-xs'>
                <Layers className='mr-1.5 size-3.5 text-muted-foreground' />
                <SelectValue
                  placeholder={t('deviceDetail.telemetryTab.selectProperty')}
                />
              </SelectTrigger>
              <SelectContent>
                {availableProperties.map((p) => (
                  <SelectItem
                    key={p.identifier}
                    value={p.identifier}
                    className='text-xs'
                  >
                    {p.name || p.identifier}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>

            {/* Time-range selector. */}
            <div className='flex items-center rounded-md border bg-muted/40 p-0.5'>
              {(Object.keys(TIME_RANGES) as TimeRangeKey[]).map((key) => (
                <Button
                  key={key}
                  size='sm'
                  variant={timeRange === key ? 'secondary' : 'ghost'}
                  onClick={() => handleRangeChange(key)}
                  className='h-7 px-2.5 text-xs font-normal shadow-none'
                >
                  {t(`deviceDetail.telemetryTab.timeRanges.${key}`)}
                </Button>
              ))}
            </div>

            {/* Manual refresh. */}
            <Button
              size='icon'
              variant='outline'
              onClick={handleManualRefresh}
              disabled={isLoading || isRefetching}
              className='size-8'
            >
              <RotateCw
                className={`size-3.5 ${isRefetching ? 'animate-spin' : ''}`}
              />
            </Button>
          </div>
        </div>

        {/* Keep statistic cards mounted to avoid layout shift while loading. */}
        <div className='mt-4 grid grid-cols-2 gap-3 sm:grid-cols-5'>
          <div className='rounded-lg border bg-card p-3'>
            <div className='text-xs font-medium text-muted-foreground'>
              {t('deviceDetail.telemetryTab.latestValue')}
            </div>
            <div className='mt-1 text-xl font-bold tracking-tight text-primary'>
              {stats.latest}{' '}
              {stats.latest !== '--' && (
                <span className='text-xs font-normal text-muted-foreground'>
                  {unit}
                </span>
              )}
            </div>
          </div>

          <div className='rounded-lg border bg-card p-3'>
            <div className='flex items-center gap-1 text-xs font-medium text-muted-foreground'>
              <ArrowUp className='size-3 text-red-500' />
              <span>{t('deviceDetail.telemetryTab.maxValue')}</span>
            </div>
            <div className='mt-1 text-xl font-bold tracking-tight'>
              {stats.max}{' '}
              {stats.max !== '--' && (
                <span className='text-xs font-normal text-muted-foreground'>
                  {unit}
                </span>
              )}
            </div>
          </div>

          <div className='rounded-lg border bg-card p-3'>
            <div className='flex items-center gap-1 text-xs font-medium text-muted-foreground'>
              <ArrowDown className='size-3 text-emerald-500' />
              <span>{t('deviceDetail.telemetryTab.minValue')}</span>
            </div>
            <div className='mt-1 text-xl font-bold tracking-tight'>
              {stats.min}{' '}
              {stats.min !== '--' && (
                <span className='text-xs font-normal text-muted-foreground'>
                  {unit}
                </span>
              )}
            </div>
          </div>

          <div className='rounded-lg border bg-card p-3'>
            <div className='flex items-center gap-1 text-xs font-medium text-muted-foreground'>
              <Activity className='size-3 text-blue-500' />
              <span>{t('deviceDetail.telemetryTab.avgValue')}</span>
            </div>
            <div className='mt-1 text-xl font-bold tracking-tight'>
              {stats.avg}{' '}
              {stats.avg !== '--' && (
                <span className='text-xs font-normal text-muted-foreground'>
                  {unit}
                </span>
              )}
            </div>
          </div>

          <div className='rounded-lg border bg-card p-3'>
            <div className='flex items-center gap-1 text-xs font-medium text-muted-foreground'>
              <Calendar className='size-3 text-muted-foreground' />
              <span>{t('deviceDetail.telemetryTab.sampleCount')}</span>
            </div>
            <div className='mt-1 text-xl font-bold tracking-tight'>
              {stats.count}{' '}
              <span className='text-xs font-normal text-muted-foreground'>
                {t('deviceDetail.telemetryTab.pointsUnit')}
              </span>
            </div>
          </div>
        </div>
      </CardHeader>

      <CardContent>
        {isLoading && chartData.length === 0 ? (
          <div className='flex h-[320px] w-full items-center justify-center'>
            <Loader2 className='size-6 animate-spin text-muted-foreground' />
          </div>
        ) : chartData.length === 0 ? (
          <div className='flex h-[320px] w-full flex-col items-center justify-center rounded-md border border-dashed p-8 text-center'>
            <TrendingUp className='size-10 text-muted-foreground/40' />
            <p className='mt-2 font-medium text-muted-foreground'>
              {t('deviceDetail.telemetryTab.noHistoryData')}
            </p>
            <p className='mt-1 text-xs text-muted-foreground/70'>
              {t('deviceDetail.telemetryTab.noHistoryDataDescription')}
            </p>
          </div>
        ) : (
          <div className='h-[320px] w-full'>
            <ResponsiveContainer width='100%' height={320}>
              <AreaChart
                data={chartData}
                margin={{ top: 10, right: 10, left: -20, bottom: 0 }}
              >
                <defs>
                  <linearGradient
                    id='telemetryColor'
                    x1='0'
                    y1='0'
                    x2='0'
                    y2='1'
                  >
                    <stop
                      offset='5%'
                      stopColor='var(--primary)'
                      stopOpacity={0.35}
                    />
                    <stop
                      offset='95%'
                      stopColor='var(--primary)'
                      stopOpacity={0.0}
                    />
                  </linearGradient>
                </defs>
                <CartesianGrid
                  strokeDasharray='3 3'
                  className='stroke-muted/40'
                  vertical={false}
                />
                <XAxis
                  dataKey='formattedTime'
                  stroke='var(--muted-foreground)'
                  fontSize={11}
                  tickLine={false}
                  axisLine={false}
                />
                <YAxis
                  stroke='var(--muted-foreground)'
                  fontSize={11}
                  tickLine={false}
                  axisLine={false}
                  domain={['auto', 'auto']}
                />
                <Tooltip
                  content={({ active, payload }) => {
                    if (active && payload && payload.length) {
                      const data = payload[0].payload
                      return (
                        <div className='rounded-lg border bg-popover/95 p-2.5 text-popover-foreground shadow-md backdrop-blur-sm'>
                          <div className='text-xs text-muted-foreground'>
                            {data.fullTime}
                          </div>
                          <div className='mt-1 flex items-center gap-2 text-sm font-semibold'>
                            <span>{propName}:</span>
                            <Badge variant='secondary' className='font-mono'>
                              {data.value} {unit}
                            </Badge>
                          </div>
                        </div>
                      )
                    }
                    return null
                  }}
                />
                <Area
                  type='monotone'
                  dataKey='value'
                  stroke='var(--primary)'
                  strokeWidth={2}
                  fillOpacity={1}
                  fill='url(#telemetryColor)'
                  isAnimationActive={false}
                  dot={
                    chartData.length < 30
                      ? { r: 3, fill: 'var(--primary)' }
                      : false
                  }
                  activeDot={{ r: 5 }}
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
