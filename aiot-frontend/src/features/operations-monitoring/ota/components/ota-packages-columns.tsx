import { type ColumnDef } from '@tanstack/react-table'
import { useTranslation } from 'react-i18next'
import { Link } from '@tanstack/react-router'
import { cn } from '@/lib/utils'
import { Badge } from '@/components/ui/badge'
import { Checkbox } from '@/components/ui/checkbox'
import { DataTableColumnHeader } from '@/components/data-table'
import { LongText } from '@/components/long-text'
import { statuses } from '../data/data'
import { type OTAPackage } from '../data/schema'
import { DataTableRowActions } from './data-table-row-actions'

export function useOTAPackagesColumns(): ColumnDef<OTAPackage>[] {
  const { t } = useTranslation('operationsMonitoring')

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
      accessorKey: 'packageName',
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t('ota.packageList.columns.packageName')}
        />
      ),
      cell: ({ row }) => {
        const packageName = row.getValue('packageName') as string
        const id = row.original.id as string
        return (
          <Link
            to='/operations-monitoring/ota/packages/$id'
            params={{ id }}
            className='block'
          >
            <LongText className='max-w-48 ps-3 font-mono text-sm cursor-pointer hover:text-primary hover:underline transition-colors'>
              {packageName}
            </LongText>
          </Link>
        )
      },
      meta: {
        className: cn(
          'drop-shadow-[0_1px_2px_rgb(0_0_0_/_0.1)] dark:drop-shadow-[0_1px_2px_rgb(255_255_255_/_0.1)]',
          'ps-0.5 max-md:sticky start-6 @4xl/content:table-cell @4xl/content:drop-shadow-none'
        ),
      },
      enableSorting: false,
      enableHiding: false,
    },
    {
      accessorKey: 'version',
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t('ota.packageList.columns.version')}
        />
      ),
      cell: ({ row }) => (
        <div className='font-mono text-sm'>{row.getValue('version')}</div>
      ),
      enableSorting: false,
    },
    {
      accessorKey: 'packageType',
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t('ota.packageList.columns.packageType')}
        />
      ),
      cell: ({ row }) => {
        const packageType = row.getValue('packageType') as string
        return (
          <Badge variant='secondary' className='capitalize'>
            {t(`ota.packageList.packageTypes.${packageType}`)}
          </Badge>
        )
      },
      filterFn: (row, id, value) => {
        return value.includes(row.getValue(id))
      },
      enableSorting: false,
    },
    {
      accessorKey: 'productName',
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t('ota.packageList.columns.productName')}
        />
      ),
      cell: ({ row }) => <div>{row.getValue('productName')}</div>,
      enableSorting: false,
    },
    {
      accessorKey: 'status',
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t('ota.packageList.columns.status')}
        />
      ),
      cell: ({ row }) => {
        const status = row.getValue('status') as string
        const statusInfo = statuses.find((s) => s.value === status)
        const variant =
          statusInfo?.variant === 'default' && status === 'completed'
            ? 'default'
            : statusInfo?.variant
        return (
          <Badge variant={variant || 'outline'} className='capitalize'>
            {t(`ota.packageList.statuses.${status}`)}
          </Badge>
        )
      },
      filterFn: (row, id, value) => {
        return value.includes(row.getValue(id))
      },
      enableSorting: false,
    },
    {
      accessorKey: 'createdAt',
      header: ({ column }) => (
        <DataTableColumnHeader
          column={column}
          title={t('ota.packageList.columns.createdAt')}
        />
      ),
      cell: ({ row }) => {
        const date = new Date(row.getValue('createdAt'))
        const formatted = date.toISOString().slice(0, 19).replace('T', ' ')
        return <div className='text-nowrap'>{formatted} UTC</div>
      },
      enableSorting: false,
    },
    {
      id: 'actions',
      header: () => (
        <div className='text-center'>
          {t('ota.packageList.columns.actions')}
        </div>
      ),
      cell: DataTableRowActions,
      meta: {
        className: cn(
          'sticky end-0 z-10 bg-background',
          'shadow-[-4px_0_6px_-2px_rgb(0_0_0_/_0.05)] dark:shadow-[-4px_0_6px_-2px_rgb(0_0_0_/_0.3)]'
        ),
      },
      enableHiding: false,
    },
  ]
}
