'use client'

import { useState } from 'react'
import { AlertTriangle } from 'lucide-react'
import { toast } from 'sonner'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { useDeleteDevice } from '../api/queries'
import { type Device } from '../data/schema'

type DeviceDeleteDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  currentRow: Device
}

export function DevicesDeleteDialog({
  open,
  onOpenChange,
  currentRow,
}: DeviceDeleteDialogProps) {
  const [value, setValue] = useState('')
  const deleteDevice = useDeleteDevice()

  const handleDelete = async () => {
    if (value.trim() !== currentRow.deviceKey) return

    try {
      await deleteDevice.mutateAsync(currentRow.id)
      toast.success('Device deleted successfully!')
      setValue('')
      onOpenChange(false)
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Failed to delete device:', error)
      const apiError = error as {
        response?: { data?: { message?: string; error?: string } }
        message?: string
      }
      const errorMessage =
        apiError?.response?.data?.message ||
        apiError?.response?.data?.error ||
        apiError?.message ||
        'Failed to delete device'
      toast.error(errorMessage)
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
      isLoading={deleteDevice.isPending}
      disabled={value.trim() !== currentRow.deviceKey}
      title={
        <span className='text-destructive'>
          <AlertTriangle
            className='me-1 inline-block stroke-destructive'
            size={18}
          />{' '}
          Delete Device
        </span>
      }
      desc={
        <div className='space-y-4'>
          <p className='mb-2'>
            Are you sure you want to delete{' '}
            <span className='font-bold'>{currentRow.name}</span>?
            <br />
            This action will permanently remove the device with key{' '}
            <span className='font-mono font-bold'>
              {currentRow.deviceKey}
            </span>{' '}
            from the system. This cannot be undone.
          </p>

          <Label className='my-2'>
            Device Key:
            <Input
              value={value}
              onChange={(e) => setValue(e.target.value)}
              placeholder='Enter device key to confirm deletion.'
              className='font-mono'
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
