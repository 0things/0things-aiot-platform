import { type Row } from '@tanstack/react-table'
import { MoreHorizontal, Pencil, Trash2, Upload, Eye } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { type OTAPackage } from '../data/schema'
import { useOTAPackagesContext } from '../hooks/use-ota-packages-context'

interface DataTableRowActionsProps {
  row: Row<OTAPackage>
}

export function DataTableRowActions({ row }: DataTableRowActionsProps) {
  const { t } = useTranslation('operationsMonitoring')
  const { setOpenDialog, setSelectedPackage } = useOTAPackagesContext()
  const pkg = row.original

  const canEdit = pkg.status === 'draft'
  const canDelete = pkg.status === 'draft'
  const canDeploy = pkg.status === 'draft'

  const handleAction = (action: 'edit' | 'delete' | 'deploy' | 'view') => {
    setSelectedPackage(pkg)
    setOpenDialog(action)
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant='ghost'
          className='flex h-8 w-8 p-0 data-[state=open]:bg-muted'
        >
          <MoreHorizontal className='h-4 w-4' />
          <span className='sr-only'>Open menu</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align='end' className='w-[160px]'>
        <DropdownMenuItem onSelect={() => handleAction('view')}>
          <Eye className='me-2 h-4 w-4' />
          {t('ota.packageList.actions.viewDetails')}
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem
          disabled={!canEdit}
          onSelect={() => canEdit && handleAction('edit')}
        >
          <Pencil className='me-2 h-4 w-4' />
          {t('ota.packageList.actions.edit')}
        </DropdownMenuItem>
        <DropdownMenuItem
          disabled={!canDeploy}
          onSelect={() => canDeploy && handleAction('deploy')}
        >
          <Upload className='me-2 h-4 w-4' />
          {t('ota.packageList.actions.deploy')}
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem
          disabled={!canDelete}
          onSelect={() => canDelete && handleAction('delete')}
          className='text-destructive'
        >
          <Trash2 className='me-2 h-4 w-4' />
          {t('ota.packageList.actions.delete')}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
