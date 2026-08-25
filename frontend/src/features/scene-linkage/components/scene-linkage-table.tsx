import { useMemo, useState } from 'react'
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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { DataTablePagination, DataTableToolbar } from '@/components/data-table'
import {
  useSceneLinkages,
  useDeleteSceneLinkage,
  sceneLinkageKeys,
} from '../api/queries'
import { statuses } from '../data/data'
import { type Scene } from '../data/schema'
import { DataTableBulkActions } from './data-table-bulk-actions'
import { sceneLinkageColumns as columns } from './scene-linkage-columns'
import { SceneLinkageDeleteDialog } from './scene-linkage-delete-dialog'
import { useSceneLinkage } from './scene-linkage-provider'

export function SceneLinkageTable() {
  const { t } = useTranslation('sceneLinkage')
  const { open, setOpen, currentRow, setCurrentRow } = useSceneLinkage()
  const queryClient = useQueryClient()
  const deleteSceneLinkage = useDeleteSceneLinkage()

  const [rowSelection, setRowSelection] = useState({})
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({})
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([])
  const [appliedFilters, setAppliedFilters] = useState<ColumnFiltersState>([])
  const [sorting, setSorting] = useState<SortingState>([])
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 10,
  })

  const status = appliedFilters.find((f) => f.id === 'status')?.value as
    string | undefined
  const search = appliedFilters.find((f) => f.id === 'name')?.value as
    string | undefined
  const enable =
    status === 'enabled' ? 1 : status === 'disabled' ? 0 : undefined

  const handleSearch = () => {
    setAppliedFilters(columnFilters)
    setPagination((prev) => ({ ...prev, pageIndex: 0 }))
  }

  const handleRefresh = () => {
    setColumnFilters([])
    setAppliedFilters([])
    setPagination((prev) => ({ ...prev, pageIndex: 0 }))
    queryClient.invalidateQueries({ queryKey: sceneLinkageKeys.lists() })
  }

  const { data: response, isLoading } = useSceneLinkages({
    page: pagination.pageIndex + 1,
    pageSize: pagination.pageSize,
    search,
    enable,
  })

  const data = useMemo(() => {
    if (!response?.items) return []
    return response.items.map((item) => ({
      id: item.id ?? 0,
      name: item.name ?? '',
      description: item.description ?? '',
      status: item.enable === 1 ? 'enabled' : 'disabled',
      createdAt: item.createdAt ?? '',
    })) as Scene[]
  }, [response])

  const pageCount = useMemo(() => {
    if (!response?.total) return data.length > 0 ? 1 : 1
    return Math.ceil(response.total / pagination.pageSize)
  }, [response?.total, pagination.pageSize, data.length])

  const handlePaginationChange = (updater: Updater<PaginationState>) => {
    setPagination((old) => {
      const next = typeof updater === 'function' ? updater(old) : updater
      if (next.pageSize !== old.pageSize && response?.total) {
        const newPageCount = Math.ceil(response.total / next.pageSize)
        if (next.pageIndex >= newPageCount) {
          return { ...next, pageIndex: 0 }
        }
      }
      return next
    })
  }

  const statusFilterOptions = useMemo(
    () =>
      statuses.map((status) => ({
        label: t(status.value === 'enabled' ? 'list.enabled' : 'list.disabled'),
        value: status.value,
      })),
    [t]
  )

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

  const handleDelete = (id: number) => {
    deleteSceneLinkage.mutate(id)
    setCurrentRow(null)
  }

  const handleMultiDelete = (ids: number[]) => {
    ids.forEach((id) => deleteSceneLinkage.mutate(id))
  }

  return (
    <div
      className={cn(
        'max-sm:has-[div[role="toolbar"]]:mb-16',
        'flex flex-1 flex-col gap-4'
      )}
    >
      <DataTableToolbar
        table={table}
        searchPlaceholder={t('list.searchPlaceholder')}
        searchKey='name'
        filters={[
          {
            columnId: 'status',
            title: 'Status',
            options: statusFilterOptions,
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
                {headerGroup.headers.map((header) => {
                  return (
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
                  )
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {isLoading ? (
              Array.from({ length: pagination.pageSize }).map((_, index) => (
                <TableRow key={index}>
                  <TableCell
                    colSpan={columns.length}
                    className='h-24 text-center'
                  >
                    Loading...
                  </TableCell>
                </TableRow>
              ))
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
                  className='h-24 text-center'
                >
                  No scenes found.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <DataTablePagination table={table} className='mt-auto' />
      <DataTableBulkActions table={table} onConfirm={handleMultiDelete} />

      {currentRow && (
        <SceneLinkageDeleteDialog
          open={open === 'delete'}
          currentRow={currentRow}
          onConfirm={handleDelete}
          onOpenChange={(next) => {
            setOpen(next ? 'delete' : null)
            if (!next) setTimeout(() => setCurrentRow(null), 500)
          }}
        />
      )}
    </div>
  )
}
