import { useState, useMemo } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import {
  type ColumnFiltersState,
  type PaginationState,
  type SortingState,
  type VisibilityState,
  flexRender,
  getCoreRowModel,
  getFacetedRowModel,
  getFacetedUniqueValues,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
} from '@tanstack/react-table'
import { AlertTriangle, RefreshCw } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
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
import { otaPackageKeys, useOTAPackages } from '../api/queries'
import { statuses, packageTypes } from '../data/data'
import { DataTableBulkActions } from './data-table-bulk-actions'
import { useOTAPackagesColumns } from './ota-packages-columns'

export function OTAPackagesTable() {
  const { t } = useTranslation('ota')
  const queryClient = useQueryClient()

  const [rowSelection, setRowSelection] = useState({})
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({})
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([])
  const [appliedFilters, setAppliedFilters] = useState<ColumnFiltersState>([])
  const [sorting, setSorting] = useState<SortingState>([])
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 20,
  })

  // Extract filter values from appliedFilters (not columnFilters)
  const status = appliedFilters.find((f) => f.id === 'status')?.value as
    string | undefined
  const packageType = appliedFilters.find((f) => f.id === 'packageType')
    ?.value as string | undefined
  const searchText = appliedFilters.find((f) => f.id === 'packageName')
    ?.value as string | undefined

  // Handler to apply filters and trigger the API call
  const handleSearch = () => {
    setAppliedFilters(columnFilters)
    setPagination((prev) => ({ ...prev, pageIndex: 0 }))
  }

  // Handler to reset filters and refresh data
  const handleRefresh = () => {
    setColumnFilters([])
    setAppliedFilters([])
    setPagination((prev) => ({ ...prev, pageIndex: 0 }))
    queryClient.invalidateQueries({ queryKey: otaPackageKeys.lists() })
  }

  const { data, isLoading, isError, refetch } = useOTAPackages({
    page: pagination.pageIndex + 1, // API pages are 1-indexed
    pageSize: pagination.pageSize,
    status,
    packageType,
    searchText,
  })
  const columns = useOTAPackagesColumns()

  // Filter the fetched data by appliedFilters (search button) so typing does
  // not affect the displayed rows until the search action is triggered.
  const tableData = useMemo(() => {
    const packages = data?.packages ?? []
    let result = packages
    if (status) result = result.filter((p) => p.status === status)
    if (packageType)
      result = result.filter((p) => p.packageType === packageType)
    if (searchText) {
      const q = searchText.toLowerCase()
      result = result.filter((p) =>
        (p.packageName || '').toLowerCase().includes(q)
      )
    }
    return result
  }, [data, status, packageType, searchText])

  // eslint-disable-next-line react-hooks/incompatible-library
  const table = useReactTable({
    data: tableData,
    columns,
    state: {
      sorting,
      pagination,
      rowSelection,
      columnFilters,
      columnVisibility,
    },
    enableRowSelection: true,
    manualFiltering: true, // Filtering is applied via appliedFilters on search
    onPaginationChange: setPagination,
    onColumnFiltersChange: setColumnFilters,
    onRowSelectionChange: setRowSelection,
    onSortingChange: setSorting,
    onColumnVisibilityChange: setColumnVisibility,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(), // Client-side sorting for current page
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
        searchPlaceholder={t('packageList.filters.searchPlaceholder')}
        searchKey='packageName'
        onSearch={handleSearch}
        onRefresh={handleRefresh}
        filters={[
          {
            columnId: 'packageType',
            title: t('packageList.columns.packageType'),
            options: packageTypes.map((type) => ({
              label: t(`packageList.packageTypes.${type.value}`),
              value: type.value,
            })),
          },
          {
            columnId: 'status',
            title: t('common:status'),
            options: statuses.map((status) => ({
              label: t(`packageList.statuses.${status.value}`),
              value: status.value,
            })),
          },
        ]}
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
                      header.column.columnDef.meta?.className
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
                  className='h-32 text-center'
                >
                  <div className='flex flex-col items-center gap-3'>
                    <AlertTriangle className='h-8 w-8 text-destructive' />
                    <p className='text-sm text-muted-foreground'>
                      {t('common:error', {
                        defaultValue: 'Error loading data. Please try again.',
                      })}
                    </p>
                    <Button
                      variant='outline'
                      size='sm'
                      onClick={() => refetch()}
                    >
                      <RefreshCw className='mr-2 h-4 w-4' />
                      {t('common:retry', { defaultValue: 'Retry' })}
                    </Button>
                  </div>
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
                        cell.column.columnDef.meta?.className
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
                  className='h-24 text-center'
                >
                  {t('packageList.emptyState.description')}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <DataTablePagination table={table} className='mt-auto' />
      <DataTableBulkActions table={table} />
    </div>
  )
}
