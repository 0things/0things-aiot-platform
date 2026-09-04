import { useState } from 'react'
import { format, subDays } from 'date-fns'
import {
  functionalUpdate,
  getCoreRowModel,
  useReactTable,
} from '@tanstack/react-table'
import { RefreshCw, Search } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { DataTablePagination } from '@/components/data-table'
import {
  DateTimeRangePicker,
  type DateTimeRangeValue,
} from '@/components/date-time-range-picker'
import {
  type DeviceServiceInvocation,
  type GetDevicesDeviceKeyThingModelServiceInvocationsParams,
  useServiceInvocations,
} from '../../api/service-invocations'

type TimeRange = DateTimeRangeValue

const defaultRange = (): TimeRange => ({
  startAt: subDays(new Date(), 7),
  endAt: new Date(),
})
const toFilters = (
  serviceIdentifier = '',
  range = defaultRange(),
  page = 1,
  pageSize = 20
): GetDevicesDeviceKeyThingModelServiceInvocationsParams => ({
  serviceIdentifier: serviceIdentifier || undefined,
  startAt: range.startAt
    ? format(range.startAt, 'yyyy-MM-dd HH:mm:ss')
    : undefined,
  endAt: range.endAt ? format(range.endAt, 'yyyy-MM-dd HH:mm:ss') : undefined,
  page,
  pageSize,
})

export function ServiceInvocationsTab({ deviceKey }: { deviceKey: string }) {
  const { t } = useTranslation('deviceManagement')
  const [identifier, setIdentifier] = useState('')
  const [range, setRange] = useState<TimeRange>(defaultRange)
  const [filters, setFilters] = useState(() => toFilters())
  const [selected, setSelected] = useState<{
    title: string
    value?: string
  } | null>(null)
  const { data, isLoading, isError, refetch } = useServiceInvocations(
    deviceKey,
    filters
  )
  const table = useReactTable({
    data: data?.invocations || [],
    columns: [],
    state: {
      pagination: {
        pageIndex: (filters.page || 1) - 1,
        pageSize: filters.pageSize || 20,
      },
    },
    pageCount: Math.max(
      1,
      Math.ceil((data?.total || 0) / (filters.pageSize || 20))
    ),
    manualPagination: true,
    onPaginationChange: (updater) => {
      const next = functionalUpdate(updater, {
        pageIndex: (filters.page || 1) - 1,
        pageSize: filters.pageSize || 20,
      })
      setFilters((current) => ({
        ...current,
        page: next.pageSize !== current.pageSize ? 1 : next.pageIndex + 1,
        pageSize: next.pageSize,
      }))
    },
    getCoreRowModel: getCoreRowModel(),
  })
  const apply = () =>
    setFilters(toFilters(identifier, range, 1, filters.pageSize))
  const reset = () => {
    const nextRange = defaultRange()
    setIdentifier('')
    setRange(nextRange)
    setFilters(toFilters('', nextRange))
  }
  return (
    <div className='flex flex-1 flex-col gap-4'>
      <div>
        <h2 className='text-2xl font-bold tracking-tight'>
          {t('deviceDetail.serviceInvocations.title')}
        </h2>
        <p className='text-muted-foreground'>
          {t('deviceDetail.serviceInvocations.description')}
        </p>
      </div>
      <div className='flex flex-wrap items-center gap-2'>
        <Input
          value={identifier}
          onChange={(event) => setIdentifier(event.target.value)}
          onKeyDown={(event) => event.key === 'Enter' && apply()}
          placeholder={t(
            'deviceDetail.serviceInvocations.filters.serviceIdentifier'
          )}
          className='h-9 w-56'
        />
        <DateTimeRangePicker
          value={range}
          onChange={setRange}
          placeholder={t('deviceDetail.serviceInvocations.filters.dateRange')}
          startAtLabel={t('deviceDetail.serviceInvocations.filters.startAt')}
          endAtLabel={t('deviceDetail.serviceInvocations.filters.endAt')}
          timePrecisionLabel={t(
            'deviceDetail.serviceInvocations.filters.timePrecision'
          )}
          className='h-9'
        />
        <Button size='sm' onClick={apply}>
          <Search className='size-4' />
          {t('common:search')}
        </Button>
        <Button variant='outline' size='sm' onClick={() => refetch()}>
          <RefreshCw className='size-4' />
          {t('common:refresh')}
        </Button>
        <Button variant='ghost' size='sm' onClick={reset}>
          {t('common:reset')}
        </Button>
      </div>
      <div className='overflow-hidden rounded-md border'>
        <Table>
          <TableHeader>
            <TableRow>
              {[
                'invokedAt',
                'serviceIdentifier',
                'serviceName',
                'inputParams',
                'outputParams',
              ].map((key) => (
                <TableHead key={key}>
                  {t(`deviceDetail.serviceInvocations.columns.${key}`)}
                </TableHead>
              ))}
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              Array.from({ length: filters.pageSize || 20 }).map((_, index) => (
                <TableRow key={index}>
                  <TableCell colSpan={5}>
                    <Skeleton className='h-5 w-full' />
                  </TableCell>
                </TableRow>
              ))
            ) : isError ? (
              <TableRow>
                <TableCell
                  colSpan={5}
                  className='h-24 text-center text-destructive'
                >
                  {t('deviceDetail.serviceInvocations.error')}
                </TableCell>
              </TableRow>
            ) : (data?.invocations || []).length ? (
              data!.invocations!.map((item) => (
                <InvocationRow
                  key={item.uuid}
                  item={item}
                  onView={setSelected}
                />
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={5}
                  className='h-24 text-center text-muted-foreground'
                >
                  {t('deviceDetail.serviceInvocations.empty')}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <DataTablePagination table={table} className='mt-auto' />
      <Dialog
        open={!!selected}
        onOpenChange={(open) => !open && setSelected(null)}
      >
        <DialogContent className='max-w-2xl'>
          <DialogHeader>
            <DialogTitle>{selected?.title}</DialogTitle>
          </DialogHeader>
          <pre className='max-h-[60vh] overflow-auto rounded-md bg-muted p-4 text-xs'>
            {formatJSON(selected?.value)}
          </pre>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function InvocationRow({
  item,
  onView,
}: {
  item: DeviceServiceInvocation
  onView: (value: { title: string; value?: string }) => void
}) {
  const { t } = useTranslation('deviceManagement')
  const jsonCell = (title: string, value?: string) => (
    <button
      type='button'
      className='block max-w-72 truncate font-mono text-sm text-muted-foreground hover:underline'
      onClick={() => onView({ title, value })}
    >
      {value || '-'}
    </button>
  )
  return (
    <TableRow>
      <TableCell className='whitespace-nowrap'>
        {item.invokedAt || '-'}
      </TableCell>
      <TableCell className='font-mono text-sm'>
        {item.serviceIdentifier}
      </TableCell>
      <TableCell>{item.serviceName}</TableCell>
      <TableCell>
        {jsonCell(
          t('deviceDetail.serviceInvocations.columns.inputParams'),
          item.inputParams
        )}
      </TableCell>
      <TableCell>
        {jsonCell(
          t('deviceDetail.serviceInvocations.columns.outputParams'),
          item.outputParams
        )}
      </TableCell>
    </TableRow>
  )
}

function formatJSON(value?: string) {
  if (!value) return '-'
  try {
    return JSON.stringify(JSON.parse(value), null, 2)
  } catch {
    return value
  }
}
