import { Link } from '@tanstack/react-router'
import { type ColumnDef } from '@tanstack/react-table'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'
import { Badge } from '@/components/ui/badge'
import { DataTableColumnHeader } from '@/components/data-table'
import { LongText } from '@/components/long-text'
import { type DeviceGroup } from '../data/schema'
import { DataTableRowActions } from './data-table-row-actions'

export const useGroupsColumns = (): ColumnDef<DeviceGroup>[] => {
  const { t } = useTranslation('deviceGroup')
  const { t: tCommon } = useTranslation('common')

  return [
    {
      accessorKey: 'groupUuid',
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t('groupUuid')} />
      ),
      cell: ({ row }) => (
        <Link
          to='/device-management/groups/$uuid'
          params={{ uuid: row.original.groupUuid }}
          className='hover:underline'
        >
          <LongText className='max-w-44 ps-1 font-mono text-sm'>
            {row.getValue('groupUuid')}
          </LongText>
        </Link>
      ),
      meta: {
        className: cn('max-md:sticky start-0 z-10 rounded-tl-[inherit]'),
      },
      enableSorting: false,
      enableHiding: false,
    },
    {
      accessorKey: 'name',
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t('name')} />
      ),
      cell: ({ row }) => (
        <Link
          to='/device-management/groups/$uuid'
          params={{ uuid: row.original.groupUuid }}
          className='hover:underline'
        >
          <LongText className='max-w-48 font-medium text-foreground'>
            {row.getValue('name')}
          </LongText>
        </Link>
      ),
      meta: { className: 'w-48' },
      enableSorting: false,
    },
    {
      accessorKey: 'type',
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t('type.label')} />
      ),
      cell: ({ row }) => {
        const type = row.getValue('type') as string
        return (
          <Badge
            variant={type === 'dynamic' ? 'default' : 'secondary'}
            className='text-xs capitalize'
          >
            {type === 'dynamic' ? t('type.dynamic') : t('type.manual')}
          </Badge>
        )
      },
      filterFn: (row, id, value) => {
        return value.includes(row.getValue(id))
      },
      meta: { className: 'w-28' },
      enableSorting: false,
    },
    {
      accessorKey: 'rule',
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t('rule')} />
      ),
      cell: ({ row }) => {
        const rule = row.getValue('rule') as string | undefined
        return (
          <LongText className='max-w-64 font-mono text-xs text-muted-foreground'>
            {rule || '-'}
          </LongText>
        )
      },
      meta: { className: 'min-w-[160px]' },
      enableSorting: false,
    },
    {
      accessorKey: 'description',
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t('descriptionLabel')} />
      ),
      cell: ({ row }) => {
        const desc = row.getValue('description') as string | undefined
        return (
          <LongText className='max-w-72 text-sm text-muted-foreground'>
            {desc || '-'}
          </LongText>
        )
      },
      meta: { className: 'min-w-[180px]' },
      enableSorting: false,
    },
    {
      accessorKey: 'createdAt',
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={tCommon('createdAt')} />
      ),
      cell: ({ row }) => {
        const raw = row.getValue('createdAt') as string | undefined
        return <div className='text-nowrap'>{raw || '-'}</div>
      },
      enableSorting: false,
    },
    {
      id: 'actions',
      header: () => <div className='text-center'>{tCommon('actions')}</div>,
      cell: ({ row }) => <DataTableRowActions row={row} />,
      meta: {
        className: cn(
          'sticky end-0 z-10 bg-background',
          'shadow-[-4px_0_6px_-2px_rgb(0_0_0_/_0.05)] dark:shadow-[-4px_0_6px_-2px_rgb(0_0_0_/_0.3)]'
        ),
      },
      enableSorting: false,
      enableHiding: false,
    },
  ]
}
