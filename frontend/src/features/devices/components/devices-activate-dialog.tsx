'use client'

import { CheckCircle } from 'lucide-react'
import { toast } from 'sonner'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { useActivateDevice } from '../api/queries'
import { type Device } from '../data/schema'

type DeviceActivateDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  currentRow: Device
}

export function DevicesActivateDialog({
  open,
  onOpenChange,
  currentRow,
}: DeviceActivateDialogProps) {
  const activateDevice = useActivateDevice()

  const handleActivate = async () => {
    try {
      await activateDevice.mutateAsync({
        id: currentRow.id,
        data: {},
      })
      toast.success('Device activated successfully!')
      onOpenChange(false)
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Failed to activate device:', error)
      const apiError = error as {
        response?: { data?: { message?: string; error?: string } }
        message?: string
      }
      const errorMessage =
        apiError?.response?.data?.message ||
        apiError?.response?.data?.error ||
        apiError?.message ||
        'Failed to activate device'
      toast.error(errorMessage)
    }
  }

  return (
    <ConfirmDialog
      open={open}
      onOpenChange={onOpenChange}
      handleConfirm={handleActivate}
      isLoading={activateDevice.isPending}
      title={
        <span>
          <CheckCircle className='me-1 inline-block' size={18} />
          Activate Device
        </span>
      }
      desc={
        <div className='space-y-4'>
          <p className='mb-2'>
            Are you sure you want to activate{' '}
            <span className='font-bold'>{currentRow.name}</span>?
            <br />
            Once activated, the device will be able to connect and communicate
            with the system.
          </p>
        </div>
      }
      confirmText='Activate'
    />
  )
}
