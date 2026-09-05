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
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from 'recharts'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  type ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from '@/components/ui/chart'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useTelemetryHistory } from '@/features/devices/api/telemetry'

interface PropertyHistoryChartProps {
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

export function PropertyHistoryChart({
  deviceKey,
  availableProperties,
}: PropertyHistoryChartProps) {
  const { t } = useTranslation('deviceManagement')
  const [selectedProperty, setSelectedProperty] = useState<string>('')
  const [timeRange, setTimeRange] = useState<TimeRangeKey>('24h')
  const [historyEndTime, setHistoryEndTime] = useState(() => Date.now())

  const activeProperty = availableProperties.some(
    (p) => p.identifier === selectedProperty
  )
    ? selectedProperty
    : availableProperties[0]?.identifier || ''

  const {
    data: points = [],
    isLoading,
    isRefetching,
  } = useTelemetryHistory(
    {
      deviceKey,
      property: activeProperty,
      durationMs: TIME_RANGES[timeRange].durationMs,
      endTime: historyEndTime,
      limit: 500,
    },
    Boolean(activeProperty)
  )

  const handleManualRefresh = () => {
    setHistoryEndTime(Date.now())
  }

  // Format numeric points and calculate aggregate statistics.
  const { chartData, stats } = useMemo(() => {
    const timeFmt = timeRange === '7d' ? 'MM-dd HH:mm' : 'HH:mm:ss'
    const numericPoints = points
      .map((p) => {
        if (p.timestamp == null || !Number.isFinite(p.timestamp)) return null
        const val =
          typeof p.value === 'number'
            ? p.value
            : typeof p.value === 'string'
              ? Number(p.value)
              : Number.NaN
        if (!Number.isFinite(val)) return null
        const ts = p.timestamp
        return {
          timestamp: ts,
          formattedTime: format(new Date(ts), timeFmt),
          fullTime: format(new Date(ts), 'yyyy-MM-dd HH:mm:ss'),
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
  }, [points, timeRange])

  const currentPropMeta = availableProperties.find(
    (p) => p.identifier === activeProperty
  )
  const unit = currentPropMeta?.unit || ''
  const propName = currentPropMeta?.name || activeProperty
  const chartConfig = {
    value: {
      label: propName,
      color: '#9AC4FE',
    },
  } satisfies ChartConfig

  return (
    <Card className='shadow-xs'>
      <CardHeader className='pb-4'>
        <div className='flex flex-wrap items-center justify-between gap-4'>
          <div>
            <div className='flex items-center gap-2'>
              <TrendingUp className='size-5 text-primary' />
              <CardTitle className='text-lg font-semibold'>
                {t('deviceDetail.propertyTab.historyTitle', {
                  property: propName,
                })}
              </CardTitle>
            </div>
            <CardDescription className='mt-1'>
              {t('deviceDetail.propertyTab.historyDescription')}
            </CardDescription>
          </div>

          {/* Chart controls. */}
          <div className='flex flex-wrap items-center gap-2'>
            {/* Property selector. */}
            <Select
              value={activeProperty}
              onValueChange={(val) => {
                setSelectedProperty(val)
                setHistoryEndTime(Date.now())
              }}
            >
              <SelectTrigger className='h-8 w-36 text-xs'>
                <Layers className='mr-1.5 size-3.5 text-muted-foreground' />
                <SelectValue
                  placeholder={t('deviceDetail.propertyTab.selectProperty')}
                />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  {availableProperties.map((p) => (
                    <SelectItem
                      key={p.identifier}
                      value={p.identifier}
                      className='text-xs'
                    >
                      {p.name || p.identifier}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>

            <Select
              value={timeRange}
              onValueChange={(value) => {
                setTimeRange(value as TimeRangeKey)
                setHistoryEndTime(Date.now())
              }}
            >
              <SelectTrigger
                className='w-[160px] rounded-lg'
                aria-label={t('deviceDetail.propertyTab.timeRange')}
              >
                <SelectValue />
              </SelectTrigger>
              <SelectContent className='rounded-xl'>
                {(Object.keys(TIME_RANGES) as TimeRangeKey[]).map((key) => (
                  <SelectItem key={key} value={key} className='rounded-lg'>
                    {t(`deviceDetail.propertyTab.timeRanges.${key}`)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>

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
              {t('deviceDetail.propertyTab.latestValue')}
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
              <span>{t('deviceDetail.propertyTab.maxValue')}</span>
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
              <span>{t('deviceDetail.propertyTab.minValue')}</span>
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
              <span>{t('deviceDetail.propertyTab.avgValue')}</span>
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
              <span>{t('deviceDetail.propertyTab.sampleCount')}</span>
            </div>
            <div className='mt-1 text-xl font-bold tracking-tight'>
              {stats.count}{' '}
              <span className='text-xs font-normal text-muted-foreground'>
                {t('deviceDetail.propertyTab.pointsUnit')}
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
              {t('deviceDetail.propertyTab.noHistoryData')}
            </p>
            <p className='mt-1 text-xs text-muted-foreground/70'>
              {t('deviceDetail.propertyTab.noHistoryDataDescription')}
            </p>
          </div>
        ) : (
          <ChartContainer
            config={chartConfig}
            className='aspect-auto h-[300px] w-full'
          >
            <AreaChart
              accessibilityLayer
              data={chartData}
              margin={{ left: 12, right: 12, top: 10, bottom: 4 }}
            >
              <defs>
                <linearGradient id='fillValue' x1='0' y1='0' x2='0' y2='1'>
                  <stop
                    offset='5%'
                    stopColor='var(--color-value)'
                    stopOpacity={0.4}
                  />
                  <stop
                    offset='95%'
                    stopColor='var(--color-value)'
                    stopOpacity={0.0}
                  />
                </linearGradient>
              </defs>
              <CartesianGrid
                vertical={false}
                strokeDasharray='3 3'
                className='stroke-muted/40'
              />
              <XAxis
                dataKey='formattedTime'
                tickLine={false}
                axisLine={false}
                tickMargin={8}
                minTickGap={40}
                stroke='var(--muted-foreground)'
                fontSize={11}
              />
              <YAxis
                tickLine={false}
                axisLine={false}
                tickMargin={8}
                stroke='var(--muted-foreground)'
                fontSize={11}
                domain={['auto', 'auto']}
              />
              <ChartTooltip
                content={
                  <ChartTooltipContent
                    className='w-[160px]'
                    labelFormatter={(_, payload) =>
                      payload?.[0]?.payload?.fullTime || ''
                    }
                    formatter={(val) => `${val}${unit ? ` ${unit}` : ''}`}
                  />
                }
              />
              <Area
                type='monotone'
                dataKey='value'
                stroke='var(--color-value)'
                fill='url(#fillValue)'
                strokeWidth={2}
                dot={
                  chartData.length <= 25
                    ? { r: 3, fill: 'var(--color-value)' }
                    : false
                }
                activeDot={{ r: 5 }}
              />
            </AreaChart>
          </ChartContainer>
        )}
      </CardContent>
    </Card>
  )
}
