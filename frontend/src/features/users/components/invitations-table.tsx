import { useState } from 'react'
import {
  type ColumnDef,
  type PaginationState,
  type SortingState,
  type Updater,
  flexRender,
  getCoreRowModel,
  getSortedRowModel,
  useReactTable,
} from '@tanstack/react-table'
import { Ban, Mail } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { type NavigateFn } from '@/hooks/use-table-url-state'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { DataTablePagination } from '@/components/data-table'
import { type InvitationItem } from '../data/schema'
import { useUsers } from './users-provider'

type InvitationsTableProps = {
  data: InvitationItem[]
  page: number
  pageSize: number
  total: number
  navigate: NavigateFn
  isLoading?: boolean
}

export function InvitationsTable({
  data,
  page,
  pageSize,
  total,
  navigate,
  isLoading,
}: InvitationsTableProps) {
  const { t } = useTranslation(['users', 'common'])
  const { setOpen, setCurrentInvitation } = useUsers()
  const [sorting, setSorting] = useState<SortingState>([])
  const pagination: PaginationState = { pageIndex: page - 1, pageSize }
  const onPaginationChange = (updater: Updater<PaginationState>) => {
    const next = typeof updater === 'function' ? updater(pagination) : updater
    navigate({
      search: (previous) => ({
        ...(previous as Record<string, unknown>),
        page: next.pageIndex === 0 ? undefined : next.pageIndex + 1,
        pageSize: next.pageSize === 10 ? undefined : next.pageSize,
      }),
    })
  }

  const getStatusBadge = (status?: string) => {
    switch ((status || 'pending').toLowerCase()) {
      case 'pending':
        return (
          <Badge
            variant='outline'
            className='border-amber-500/30 bg-amber-500/10 text-amber-600 dark:text-amber-400'
          >
            {t('statusPending')}
          </Badge>
        )
      case 'accepted':
        return (
          <Badge
            variant='outline'
            className='border-emerald-500/30 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400'
          >
            {t('statusAccepted')}
          </Badge>
        )
      case 'revoked':
        return (
          <Badge
            variant='outline'
            className='border-zinc-500/30 bg-zinc-500/10 text-zinc-500'
          >
            {t('statusRevoked')}
          </Badge>
        )
      case 'expired':
        return (
          <Badge
            variant='outline'
            className='border-rose-500/30 bg-rose-500/10 text-rose-600 dark:text-rose-400'
          >
            {t('statusExpired')}
          </Badge>
        )
      default:
        return <Badge variant='outline'>{status}</Badge>
    }
  }

  const columns: ColumnDef<InvitationItem>[] = [
    {
      accessorKey: 'invitee',
      header: t('common:email'),
      cell: ({ row }) => {
        const invitee = row.original.invitee
        return (
          <div className='flex items-center gap-2'>
            <div className='flex h-8 w-8 items-center justify-center rounded-full bg-primary/10 text-primary'>
              <Mail className='h-4 w-4' />
            </div>
            <span className='font-medium text-foreground'>{invitee}</span>
          </div>
        )
      },
    },
    {
      accessorKey: 'status',
      header: t('common:status'),
      cell: ({ row }) => getStatusBadge(row.original.status),
    },
    {
      accessorKey: 'createdAt',
      header: t('createdAt'),
      cell: ({ row }) => {
        const createdAt = row.original.createdAt
        if (!createdAt) return '-'
        return (
          <span className='text-sm text-muted-foreground'>{createdAt}</span>
        )
      },
    },
    {
      accessorKey: 'expiresAt',
      header: t('expiresAt'),
      cell: ({ row }) => {
        const expiresAt = row.original.expiresAt
        if (!expiresAt) return '-'
        const isExpired = row.original.isExpired
        return (
          <span
            className={`text-sm ${
              isExpired ? 'text-rose-500' : 'text-muted-foreground'
            }`}
          >
            {expiresAt}
          </span>
        )
      },
    },
    {
      id: 'actions',
      cell: ({ row }) => {
        const invitation = row.original
        const isPending = invitation.status?.toLowerCase() === 'pending'
        if (!isPending) return null

        return (
          <div className='flex justify-end'>
            <Button
              variant='ghost'
              size='sm'
              className='text-destructive hover:bg-destructive/10 hover:text-destructive'
              onClick={() => {
                setCurrentInvitation(invitation)
                setOpen('revoke')
              }}
            >
              <Ban className='mr-1 h-4 w-4' />
              {t('revokeInvite')}
            </Button>
          </div>
        )
      },
    },
  ]

  // eslint-disable-next-line react-hooks/incompatible-library
  const table = useReactTable({
    data,
    columns,
    state: {
      sorting,
      pagination,
    },
    pageCount: Math.max(1, Math.ceil(total / pageSize)),
    manualPagination: true,
    onPaginationChange,
    onSortingChange: setSorting,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  })

  return (
    <div className='flex flex-1 flex-col gap-4'>
      <div className='overflow-hidden rounded-md border'>
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <TableHead key={header.id}>
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
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className='h-24 text-center text-muted-foreground'
                >
                  {t('invitationLoading')}
                </TableCell>
              </TableRow>
            ) : table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow key={row.id}>
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
                  className='h-24 text-center text-muted-foreground'
                >
                  {t('common:noResults')}
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
