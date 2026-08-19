'use client'

import { useState } from 'react'
import { type Table } from '@tanstack/react-table'
import { AlertTriangle } from 'lucide-react'
import { toast } from 'sonner'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { useDeleteProduct } from '../api/queries'
import { type Product } from '../data/schema'

type ProductMultiDeleteDialogProps<TData> = {
  open: boolean
  onOpenChange: (open: boolean) => void
  table: Table<TData>
}

const CONFIRM_WORD = 'DELETE'

export function ProductsMultiDeleteDialog<TData>({
  open,
  onOpenChange,
  table,
}: ProductMultiDeleteDialogProps<TData>) {
  const [value, setValue] = useState('')
  const [isDeleting, setIsDeleting] = useState(false)
  const deleteProduct = useDeleteProduct()

  const selectedRows = table.getFilteredSelectedRowModel().rows

  const handleDelete = async () => {
    if (value.trim() !== CONFIRM_WORD) {
      toast.error(`Please type "${CONFIRM_WORD}" to confirm.`)
      return
    }

    setIsDeleting(true)
    const products = selectedRows.map((row) => row.original as Product)
    let successCount = 0
    let failCount = 0

    try {
      // Delete products one by one
      for (const product of products) {
        try {
          await deleteProduct.mutateAsync(product.id)
          successCount++
        } catch (error) {
          // eslint-disable-next-line no-console
          console.error(
            `Failed to delete product ${product.productKey}:`,
            error
          )
          failCount++
        }
      }

      // Show result
      if (successCount > 0) {
        toast.success(
          `Successfully deleted ${successCount} ${successCount > 1 ? 'products' : 'product'}`
        )
      }
      if (failCount > 0) {
        toast.error(
          `Failed to delete ${failCount} ${failCount > 1 ? 'products' : 'product'}`
        )
      }

      setValue('')
      table.resetRowSelection()
      onOpenChange(false)
    } finally {
      setIsDeleting(false)
    }
  }

  return (
    <ConfirmDialog
      open={open}
      onOpenChange={(open) => {
        if (!open) setValue('')
        onOpenChange(open)
      }}
      handleConfirm={handleDelete}
      isLoading={isDeleting}
      disabled={value.trim() !== CONFIRM_WORD}
      title={
        <span className='text-destructive'>
          <AlertTriangle
            className='me-1 inline-block stroke-destructive'
            size={18}
          />{' '}
          Delete {selectedRows.length}{' '}
          {selectedRows.length > 1 ? 'products' : 'product'}
        </span>
      }
      desc={
        <div className='space-y-4'>
          <p className='mb-2'>
            Are you sure you want to delete the selected products? <br />
            This action cannot be undone.
          </p>

          <Label className='my-4 flex flex-col items-start gap-1.5'>
            <span className=''>Confirm by typing "{CONFIRM_WORD}":</span>
            <Input
              value={value}
              onChange={(e) => setValue(e.target.value)}
              placeholder={`Type "${CONFIRM_WORD}" to confirm.`}
            />
          </Label>

          <Alert variant='destructive'>
            <AlertTitle>Warning!</AlertTitle>
            <AlertDescription>
              Please be careful, this operation can not be rolled back.
            </AlertDescription>
          </Alert>
        </div>
      }
      confirmText='Delete'
      destructive
    />
  )
}
