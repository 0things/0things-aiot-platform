import { useState, useMemo } from 'react'
import { useQueryClient } from '@tanstack/react-query'
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
import { deviceGroupKeys, useDeviceGroups } from '../api/queries'
import { type DeviceGroup } from '../data/schema'
import { useGroupsColumns } from './groups-columns'

export function GroupsTable() {
  const { t } = useTranslation('deviceGroup')
  const queryClient = useQueryClient()
  const columns = useGroupsColumns()

  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({})
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([])
  const [appliedFilters, setAppliedFilters] = useState<ColumnFiltersState>([])
  const [sorting, setSorting] = useState<SortingState>([])
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 10,
  })

  const typeFilter = appliedFilters.find((f) => f.id === 'type')
  const typeFilters = typeFilter?.value as string[] | undefined
  const typeQuery =
    typeFilters && typeFilters.length === 1 ? typeFilters[0] : undefined

  const searchFilter = appliedFilters.find((f) => f.id === 'name')
  const searchText = (searchFilter?.value as string) || undefined

  const handleSearch = () => {
    setAppliedFilters(columnFilters)
    setPagination((prev) => ({ ...prev, pageIndex: 0 }))
  }

  const handleRefresh = () => {
    setColumnFilters([])
    setAppliedFilters([])
    setPagination((prev) => ({ ...prev, pageIndex: 0 }))
    queryClient.invalidateQueries({ queryKey: deviceGroupKeys.lists() })
  }

  const {
    data: response,
    isLoading,
    isError,
  } = useDeviceGroups({
    page: pagination.pageIndex + 1,
    pageSize: pagination.pageSize,
    search: searchText,
    type: typeQuery,
  })

  const data = useMemo(() => {
    const rawItems = response?.items || []
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

  const totalCount = response?.total || 0
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
      columnFilters,
      columnVisibility,
    },
    enableRowSelection: false,
    manualPagination: true,
    manualFiltering: true,
    onPaginationChange: handlePaginationChange,
    onColumnFiltersChange: setColumnFilters,
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
              Array.from({ length: 5 }).map((_, index) => (
                <TableRow key={index}>
                  {columns.map((_column, colIndex) => (
                    <TableCell key={colIndex}>
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
                  {t('table.error', {
                    defaultValue: 'Failed to load device groups.',
                  })}
                </TableCell>
              </TableRow>
            ) : data.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className='h-24 text-center text-muted-foreground'
                >
                  {t('noResults')}
                </TableCell>
              </TableRow>
            ) : (
              table.getRowModel().rows.map((row) => (
                <TableRow key={row.id}>
                  {row.getVisibleCells().map((cell) => (
                    <TableCell
                      key={cell.id}
                      className={cell.column.columnDef.meta?.className}
                    >
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
      <DataTablePagination table={table} />
    </div>
  )
}
