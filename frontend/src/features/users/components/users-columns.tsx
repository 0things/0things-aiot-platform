import { type ColumnDef } from '@tanstack/react-table'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'
import { Badge } from '@/components/ui/badge'
import { Checkbox } from '@/components/ui/checkbox'
import { DataTableColumnHeader } from '@/components/data-table'
import { LongText } from '@/components/long-text'
import { callTypes } from '../data/data'
import { type User } from '../data/schema'
import { DataTableRowActions } from './data-table-row-actions'

export const useUsersColumns = (): ColumnDef<User>[] => {
  const { t } = useTranslation(['users', 'common'])

  return [
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
      meta: {
        className: cn('max-md:sticky start-0 z-10 rounded-tl-[inherit]'),
      },
      cell: ({ row }) => (
        <Checkbox
          checked={row.getIsSelected()}
          onCheckedChange={(value) => row.toggleSelected(!!value)}
          aria-label='Select row'
          className='translate-y-[2px]'
        />
      ),
      enableSorting: false,
      enableHiding: false,
    },
    {
      accessorKey: 'username',
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t('users:username')} />
      ),
      cell: ({ row }) => (
        <LongText className='max-w-36 ps-3'>
          {row.getValue('username')}
        </LongText>
      ),
      meta: {
        className: cn(
          'drop-shadow-[0_1px_2px_rgb(0_0_0_/_0.1)] dark:drop-shadow-[0_1px_2px_rgb(255_255_255_/_0.1)]',
          'ps-0.5 max-md:sticky start-6 @4xl/content:table-cell @4xl/content:drop-shadow-none'
        ),
      },
      enableHiding: false,
    },
    {
      id: 'fullName',
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t('common:name')} />
      ),
      cell: ({ row }) => {
        const { firstName, lastName } = row.original
        const fullName = `${firstName} ${lastName}`.trim() || '-'
        return <LongText className='max-w-36'>{fullName}</LongText>
      },
      meta: { className: 'w-36' },
    },
    {
      accessorKey: 'email',
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t('common:email')} />
      ),
      cell: ({ row }) => (
        <div className='w-fit ps-2 text-nowrap'>
          {row.getValue('email') || '-'}
        </div>
      ),
    },
    {
      accessorKey: 'phoneNumber',
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t('users:phoneNumber')} />
      ),
      cell: ({ row }) => <div>{row.getValue('phoneNumber') || '-'}</div>,
      enableSorting: false,
    },
    {
      accessorKey: 'status',
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t('common:status')} />
      ),
      cell: ({ row }) => {
        const { status } = row.original
        const badgeColor = callTypes.get(status)
        const statusKey = `users:status${status.charAt(0).toUpperCase() + status.slice(1)}`
        return (
          <div className='flex space-x-2'>
            <Badge variant='outline' className={cn('capitalize', badgeColor)}>
              {t(statusKey, { defaultValue: status })}
            </Badge>
          </div>
        )
      },
      filterFn: (row, id, value) => {
        return value.includes(row.getValue(id))
      },
      enableHiding: false,
      enableSorting: false,
    },
    {
      id: 'actions',
      header: () => <div className='text-center'>{t('common:actions')}</div>,
      cell: DataTableRowActions,
    },
  ]
}
