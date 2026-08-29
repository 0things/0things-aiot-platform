import { useState, useMemo } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import {
  type ColumnFiltersState,
  type PaginationState,
  type SortingState,
  type Updater,
  type VisibilityState,
  flexRender,
  getCoreRowModel,
  getFacetedRowModel,
  getFacetedUniqueValues,
  getSortedRowModel,
  useReactTable,
} from '@tanstack/react-table'
import { useTranslation } from 'react-i18next'
import { getDeviceGroups } from '@/api/generated'
import { cn } from '@/lib/utils'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { DataTablePagination, DataTableToolbar } from '@/components/data-table'
import { type DeviceGroup } from '../data/schema'
import { useGroupsColumns } from './groups-columns'

export function GroupsTable() {
  const { t } = useTranslation('deviceGroup')
  const queryClient = useQueryClient()
  const columns = useGroupsColumns()

  const [rowSelection, setRowSelection] = useState({})
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({})
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([])
  const [appliedFilters, setAppliedFilters] = useState<ColumnFiltersState>([])
  const [sorting, setSorting] = useState<SortingState>([])
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 10,
  })

  const searchText = appliedFilters.find((f) => f.id === 'name')?.value as
    string | undefined
  const typeFilters = appliedFilters.find((f) => f.id === 'type')?.value as
    string[] | undefined
  const typeQuery =
    typeFilters && typeFilters.length === 1 ? typeFilters[0] : undefined

  const handleSearch = () => {
    setAppliedFilters(columnFilters)
    setPagination((prev) => ({ ...prev, pageIndex: 0 }))
  }

  const handleRefresh = () => {
    setColumnFilters([])
    setAppliedFilters([])
    setPagination((prev) => ({ ...prev, pageIndex: 0 }))
    queryClient.invalidateQueries({ queryKey: ['device-groups'] })
  }

  const {
    data: response,
    isLoading,
    isError,
  } = useQuery({
    queryKey: [
      'device-groups',
      pagination.pageIndex + 1,
      pagination.pageSize,
      searchText,
      typeQuery,
    ],
    queryFn: () =>
      getDeviceGroups({
        page: pagination.pageIndex + 1,
        pageSize: pagination.pageSize,
        search: searchText,
        type: typeQuery,
      }),
  })

  const data = useMemo(() => {
    const rawItems = response?.data?.items || []
    let result: DeviceGroup[] = rawItems.map((g) => ({
      groupUuid: g.groupUuid || '',
      name: g.name || '',
      type: (g.type as 'manual' | 'dynamic') || 'manual',
      description: g.description,
      rule: g.rule,
      createdAt: g.createdAt,
      updatedAt: g.updatedAt,
    }))

    if (typeFilters && typeFilters.length > 0 && !typeQuery) {
      result = result.filter((g) => typeFilters.includes(g.type))
    }

    return result
  }, [response, typeFilters, typeQuery])

  const totalCount = response?.data?.total || 0
  const pageCount = Math.ceil(totalCount / pagination.pageSize)

  const handlePaginationChange = (updater: Updater<PaginationState>) => {
    setPagination((prev) => {
      const next = typeof updater === 'function' ? updater(prev) : updater
      return {
        ...next,
        pageIndex: Math.max(0, next.pageIndex),
        pageSize: Math.max(1, next.pageSize),
      }
    })
  }

  // eslint-disable-next-line react-hooks/incompatible-library
  const table = useReactTable({
    data,
    columns,
    pageCount,
    state: {
      sorting,
      pagination,
      rowSelection,
      columnFilters,
      columnVisibility,
    },
    enableRowSelection: true,
    manualPagination: true,
    manualFiltering: true,
    onPaginationChange: handlePaginationChange,
    onColumnFiltersChange: setColumnFilters,
    onRowSelectionChange: setRowSelection,
    onSortingChange: setSorting,
    onColumnVisibilityChange: setColumnVisibility,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFacetedRowModel: getFacetedRowModel(),
    getFacetedUniqueValues: getFacetedUniqueValues(),
  })

  return (
    <div
      className={cn(
        'max-sm:has-[div[role="toolbar"]]:mb-16',
        'flex flex-1 flex-col gap-4'
      )}
    >
      <DataTableToolbar
        table={table}
        searchPlaceholder={t('search')}
        searchKey='name'
        filters={[
          {
            columnId: 'type',
            title: t('type.label'),
            options: [
              { label: t('type.manual'), value: 'manual' },
              { label: t('type.dynamic'), value: 'dynamic' },
            ],
          },
        ]}
        onSearch={handleSearch}
        onRefresh={handleRefresh}
      />
      <div className='overflow-hidden rounded-md border'>
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id} className='group/row'>
                {headerGroup.headers.map((header) => (
                  <TableHead
                    key={header.id}
                    colSpan={header.colSpan}
                    className={cn(
                      'bg-background group-hover/row:bg-muted group-data-[state=selected]/row:bg-muted',
                      header.column.columnDef.meta?.className,
                      header.column.columnDef.meta?.thClassName
                    )}
                  >
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
              Array.from({ length: pagination.pageSize }).map((_, index) => (
                <TableRow key={index}>
                  {columns.map((_, cellIndex) => (
                    <TableCell key={cellIndex}>
                      <Skeleton className='h-6 w-full' />
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : isError ? (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className='h-24 text-center text-destructive'
                >
                  {t('common:error')}
                </TableCell>
              </TableRow>
            ) : table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && 'selected'}
                  className='group/row'
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell
                      key={cell.id}
                      className={cn(
                        'bg-background group-hover/row:bg-muted group-data-[state=selected]/row:bg-muted',
                        cell.column.columnDef.meta?.className,
                        cell.column.columnDef.meta?.tdClassName
                      )}
                    >
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
                  colSpan={columns.length}
                  className='h-24 text-center text-muted-foreground'
                >
                  {t('empty')}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <DataTablePagination table={table} className='mt-auto' />
    </div>
  )
}
