import { useState } from 'react'
import {
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
  type ColumnFiltersState,
  type SortingState,
  type VisibilityState,
} from '@tanstack/react-table'
import { Settings2, Search } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Input } from '@/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import type { Rule } from '../data/schema'
import { createColumns } from './rules-columns'

interface RulesTableProps {
  data: Rule[]
  isLoading?: boolean
  onEdit: (rule: Rule) => void
  onView: (rule: Rule) => void
  onDelete: (rule: Rule) => void
  onToggleStatus: (rule: Rule) => void
  onTrigger: (rule: Rule) => void
  onBulkDelete?: (ruleIds: string[]) => void
}

export function RulesTable({
  data,
  isLoading,
  onEdit,
  onView,
  onDelete,
  onToggleStatus,
  onTrigger,
  onBulkDelete,
}: RulesTableProps) {
  const [sorting, setSorting] = useState<SortingState>([])
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([])
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({})
  const [rowSelection, setRowSelection] = useState({})
  const [globalFilter, setGlobalFilter] = useState('')

  const columns = createColumns({
    onEdit,
    onView,
    onDelete,
    onToggleStatus,
    onTrigger,
  })

  const table = useReactTable({
    data,
    columns,
    state: {
      sorting,
      columnFilters,
      columnVisibility,
      rowSelection,
      globalFilter,
    },
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    onColumnVisibilityChange: setColumnVisibility,
    onRowSelectionChange: setRowSelection,
    onGlobalFilterChange: setGlobalFilter,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
  })

  const selectedRows = table.getFilteredSelectedRowModel().rows
  const selectedRuleIds = selectedRows.map((row) => row.original.id)

  return (
    <div className='space-y-4'>
      {/* 工具栏 */}
      <div className='flex items-center justify-between'>
        <div className='flex flex-1 items-center gap-2'>
          {/* 搜索框 */}
          <div className='relative max-w-sm'>
            <Search className='absolute top-2.5 left-2.5 h-4 w-4 text-muted-foreground' />
            <Input
              placeholder='搜索规则名称或描述...'
              value={globalFilter}
              onChange={(e) => setGlobalFilter(e.target.value)}
              className='pl-8'
            />
          </div>

          {/* 批量操作 */}
          {selectedRows.length > 0 && (
            <div className='flex items-center gap-2'>
              <span className='text-sm text-muted-foreground'>
                已选择 {selectedRows.length} 项
              </span>
              {onBulkDelete && (
                <Button
                  variant='destructive'
                  size='sm'
                  onClick={() => {
                    onBulkDelete(selectedRuleIds)
                    setRowSelection({})
                  }}
                >
                  批量删除
                </Button>
              )}
            </div>
          )}
        </div>

        {/* 列可见性控制 */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant='outline' size='sm' className='ml-auto'>
              <Settings2 className='mr-2 h-4 w-4' />
              列显示
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align='end'>
            {table
              .getAllColumns()
              .filter((column) => column.getCanHide())
              .map((column) => {
                return (
                  <DropdownMenuCheckboxItem
                    key={column.id}
                    className='capitalize'
                    checked={column.getIsVisible()}
                    onCheckedChange={(value) =>
                      column.toggleVisibility(!!value)
                    }
                  >
                    {column.id}
                  </DropdownMenuCheckboxItem>
                )
              })}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      {/* 表格 */}
      <div className='rounded-md border'>
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead key={header.id}>
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
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className='h-24 text-center'
                >
                  加载中...
                </TableCell>
              </TableRow>
            ) : table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && 'selected'}
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
                  colSpan={columns.length}
                  className='h-24 text-center'
                >
                  暂无数据
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      {/* 分页 */}
      <div className='flex items-center justify-between'>
        <div className='text-sm text-muted-foreground'>
          共 {table.getFilteredRowModel().rows.length} 条记录
          {selectedRows.length > 0 && ` (已选择 ${selectedRows.length} 条)`}
        </div>
        <div className='flex items-center space-x-2'>
          <Button
            variant='outline'
            size='sm'
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
          >
            上一页
          </Button>
          <div className='flex items-center gap-1 text-sm'>
            <span>
              第 {table.getState().pagination.pageIndex + 1} 页，共{' '}
              {table.getPageCount()} 页
            </span>
          </div>
          <Button
            variant='outline'
            size='sm'
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
          >
            下一页
          </Button>
        </div>
      </div>
    </div>
  )
}
