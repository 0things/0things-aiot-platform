import { useState } from 'react'
import { endOfDay, format, startOfDay, subDays } from 'date-fns'
import {
  type ColumnDef,
  flexRender,
  functionalUpdate,
  getCoreRowModel,
  getFilteredRowModel,
  useReactTable,
} from '@tanstack/react-table'
import { AlertTriangle, CalendarDays, RefreshCw } from 'lucide-react'
import type { DateRange } from 'react-day-picker'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Calendar } from '@/components/ui/calendar'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  DataTablePagination,
  DataTableViewOptions,
} from '@/components/data-table'
import {
  type DeviceEvent,
  type DeviceEventFilters,
  useDeviceEvents,
} from '../api/queries'
import { EventsBulkActions } from './events-bulk-actions'

const pageSize = 20
const defaultDateRange = () => {
  const today = new Date()
  return {
    from: startOfDay(subDays(today, 6)),
    to: endOfDay(today),
  }
}

const defaultFilters = (): DeviceEventFilters => {
  const range = defaultDateRange()
  return {
    page: 1,
    pageSize,
    startAt: range.from.toISOString(),
    endAt: range.to.toISOString(),
  }
}

export function DeviceEvents() {
  const { t } = useTranslation('operationsMonitoring')
  const [filters, setFilters] = useState<DeviceEventFilters>(defaultFilters)
  const [selectedEvent, setSelectedEvent] = useState<DeviceEvent | null>(null)
  const [dateRange, setDateRange] = useState<DateRange | undefined>(
    defaultDateRange
  )
  const { data, isLoading, isError, refetch } = useDeviceEvents(filters)
  const [rowSelection, setRowSelection] = useState<Record<string, boolean>>({})
  const columns: ColumnDef<DeviceEvent>[] = [
    {
      id: 'select',
      header: ({ table }) => (
        <Checkbox
          checked={
            table.getIsAllPageRowsSelected() ||
            (table.getIsSomePageRowsSelected() && 'indeterminate')
          }
          onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
          aria-label='Select all'
          className='translate-y-[2px]'
        />
      ),
      cell: ({ row }) => (
        <Checkbox
          checked={row.getIsSelected()}
          onCheckedChange={(value) => row.toggleSelected(!!value)}
          onClick={(event) => event.stopPropagation()}
          aria-label='Select row'
          className='translate-y-[2px]'
        />
      ),
      enableSorting: false,
      enableHiding: false,
    },
    {
      accessorKey: 'eventAt',
      header: t('events.columns.eventAt'),
      cell: ({ row }) => (
        <span className='whitespace-nowrap'>
          {format(new Date(row.original.eventAt), 'yyyy-MM-dd HH:mm:ss')}
        </span>
      ),
    },
    {
      accessorKey: 'deviceName',
      header: t('events.columns.deviceName'),
      cell: ({ row }) => row.original.deviceName || row.original.deviceKey,
    },
    {
      accessorKey: 'deviceKey',
      header: t('events.columns.deviceKey'),
      cell: ({ row }) => (
        <span className='font-mono text-xs'>{row.original.deviceKey}</span>
      ),
    },
    {
      accessorKey: 'eventType',
      header: t('events.columns.eventType'),
      cell: ({ row }) => (
        <span className='rounded-full bg-muted px-2 py-1 font-mono text-xs'>
          {row.original.eventType}
        </span>
      ),
    },
    {
      accessorKey: 'data',
      header: t('events.columns.data'),
      cell: ({ row }) => (
        <span className='block max-w-80 truncate font-mono text-xs text-muted-foreground'>
          {row.original.data || '{}'}
        </span>
      ),
    },
  ]
  const table = useReactTable({
    data: data?.events || [],
    columns,
    state: {
      pagination: { pageIndex: filters.page - 1, pageSize: filters.pageSize },
      rowSelection,
    },
    enableRowSelection: true,
    onRowSelectionChange: setRowSelection,
    manualPagination: true,
    pageCount: Math.max(1, Math.ceil((data?.total || 0) / filters.pageSize)),
    onPaginationChange: (updater) => {
      const pagination = functionalUpdate(updater, {
        pageIndex: filters.page - 1,
        pageSize: filters.pageSize,
      })
      setFilters((current) => ({
        ...current,
        page: pagination.pageIndex + 1,
        pageSize: pagination.pageSize,
      }))
    },
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
  })
  const updateFilters = (values: Partial<DeviceEventFilters>) => {
    setFilters((current) => ({ ...current, ...values, page: 1 }))
  }
  return (
    <div className='flex flex-1 flex-col gap-4'>
      <div>
        <h2 className='text-2xl font-bold tracking-tight'>
          {t('events.title')}
        </h2>
        <p className='text-muted-foreground'>{t('events.description')}</p>
      </div>

      <div className='flex items-start justify-between gap-2'>
        <div className='flex flex-wrap items-center gap-2'>
          <Input
            value={filters.keyword || ''}
            onChange={(event) => updateFilters({ keyword: event.target.value })}
            placeholder={t('events.filters.keyword')}
            className='h-9 w-full sm:w-56'
          />
          <Input
            value={filters.eventType || ''}
            onChange={(event) =>
              updateFilters({ eventType: event.target.value })
            }
            placeholder={t('events.filters.eventType')}
            className='h-9 w-full sm:w-44'
          />
          <Popover>
            <PopoverTrigger asChild>
              <Button
                variant='outline'
                className='h-9 w-full justify-start font-normal sm:w-64'
              >
                <CalendarDays className='size-4 text-muted-foreground' />
                {dateRange?.from
                  ? dateRange.to
                    ? `${format(dateRange.from, 'yyyy-MM-dd')} ~ ${format(dateRange.to, 'yyyy-MM-dd')}`
                    : format(dateRange.from, 'yyyy-MM-dd')
                  : t('events.filters.dateRange')}
              </Button>
            </PopoverTrigger>
            <PopoverContent className='w-auto p-0' align='start'>
              <Calendar
                mode='range'
                selected={dateRange}
                onSelect={(range) => {
                  setDateRange(range)
                  updateFilters({
                    startAt: range?.from
                      ? startOfDay(range.from).toISOString()
                      : undefined,
                    endAt: range?.to
                      ? endOfDay(range.to).toISOString()
                      : undefined,
                  })
                }}
                numberOfMonths={2}
              />
            </PopoverContent>
          </Popover>
          <Button variant='outline' size='sm' onClick={() => refetch()}>
            <RefreshCw className='size-4' />
            {t('common:refresh')}
          </Button>
          <Button
            variant='ghost'
            size='sm'
            onClick={() => {
              setDateRange(defaultDateRange())
              setFilters(defaultFilters())
            }}
          >
            {t('common:reset')}
          </Button>
        </div>
        <DataTableViewOptions
          table={table}
          labels={{
            view: t('events.viewOptions.view'),
            toggleColumns: t('events.viewOptions.toggleColumns'),
            columns: {
              eventAt: t('events.columns.eventAt'),
              deviceName: t('events.columns.deviceName'),
              deviceKey: t('events.columns.deviceKey'),
              eventType: t('events.columns.eventType'),
              data: t('events.columns.data'),
            },
          }}
        />
      </div>

      <div className='overflow-hidden rounded-md border'>
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <TableHead key={header.id} colSpan={header.colSpan}>
                    {header.isPlaceholder
                      ? null
                      : flexRender(
                          header.column.columnDef.header,
                          header.getContext()
                        )}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {isLoading ? (
              Array.from({ length: pageSize }).map((_, index) => (
                <TableRow key={index}>
                  <TableCell colSpan={table.getVisibleLeafColumns().length}>
                    <Skeleton className='h-5 w-full' />
                  </TableCell>
                </TableRow>
              ))
            ) : isError ? (
              <TableRow>
                <TableCell
                  colSpan={table.getVisibleLeafColumns().length}
                  className='h-32 text-center'
                >
                  <div className='flex flex-col items-center gap-2'>
                    <AlertTriangle className='size-6 text-destructive' />
                    <span>{t('events.error')}</span>
                    <Button
                      variant='outline'
                      size='sm'
                      onClick={() => refetch()}
                    >
                      {t('events.retry')}
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ) : table.getRowModel().rows.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  className='cursor-pointer'
                  onClick={() => setSelectedEvent(row.original)}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={table.getVisibleLeafColumns().length}
                  className='h-24 text-center text-muted-foreground'
                >
                  {t('events.empty')}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      <DataTablePagination table={table} className='mt-auto' />

      <EventsBulkActions table={table} />

      <Dialog
        open={!!selectedEvent}
        onOpenChange={(open) => !open && setSelectedEvent(null)}
      >
        <DialogContent className='max-w-2xl'>
          <DialogHeader>
            <DialogTitle>{t('events.detail.title')}</DialogTitle>
            <DialogDescription>
              {selectedEvent?.eventType} ·{' '}
              {selectedEvent?.deviceName || selectedEvent?.deviceKey}
            </DialogDescription>
          </DialogHeader>
          <pre className='max-h-[60vh] overflow-auto rounded-md bg-muted p-4 text-xs'>
            {formatEventData(selectedEvent?.data)}
          </pre>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function formatEventData(data?: string) {
  if (!data) return '{}'
  try {
    return JSON.stringify(JSON.parse(data), null, 2)
  } catch {
    return data
  }
}
