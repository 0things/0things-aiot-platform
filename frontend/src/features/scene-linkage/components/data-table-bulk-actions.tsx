import { useState } from 'react'
import { type Table } from '@tanstack/react-table'
import { Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { DataTableBulkActions as BulkActionsToolbar } from '@/components/data-table'
import { SceneLinkageMultiDeleteDialog } from './scene-linkage-multi-delete-dialog'

type DataTableBulkActionsProps<TData> = {
  table: Table<TData>
  onConfirm: (ids: number[]) => void
}

export function DataTableBulkActions<TData>({
  table,
  onConfirm,
}: DataTableBulkActionsProps<TData>) {
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false)

  return (
    <>
      <BulkActionsToolbar table={table} entityName='scene'>
        <Tooltip>
          <TooltipTrigger asChild>
            <Button
              variant='destructive'
              size='icon'
              onClick={() => setShowDeleteConfirm(true)}
              className='size-8'
              aria-label='Delete selected scenes'
              title='Delete selected scenes'
            >
              <Trash2 />
              <span className='sr-only'>Delete selected scenes</span>
            </Button>
          </TooltipTrigger>
          <TooltipContent>
            <p>Delete selected scenes</p>
          </TooltipContent>
        </Tooltip>
      </BulkActionsToolbar>

      <SceneLinkageMultiDeleteDialog
        table={table}
        open={showDeleteConfirm}
        onOpenChange={setShowDeleteConfirm}
        onConfirm={onConfirm}
      />
    </>
  )
}
