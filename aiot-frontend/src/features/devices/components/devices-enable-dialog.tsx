'use client'

import { Power } from 'lucide-react'
import { toast } from 'sonner'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { useSetDeviceEnabled } from '../api/queries'
import { type Device } from '../data/schema'

type DeviceEnableDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  currentRow: Device
}

export function DevicesEnableDialog({
  open,
  onOpenChange,
  currentRow,
}: DeviceEnableDialogProps) {
  const isEnabled = currentRow.enabled
  const action = isEnabled ? 'disable' : 'enable'
  const setDeviceEnabled = useSetDeviceEnabled()

  const handleToggle = async () => {
    try {
      await setDeviceEnabled.mutateAsync({
        id: currentRow.id,
        data: {
          enabled: !isEnabled,
        },
      })
      toast.success(`Device ${action}d successfully!`)
      onOpenChange(false)
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error(`Failed to ${action} device:`, error)
      const apiError = error as {
        response?: { data?: { message?: string; error?: string } }
        message?: string
      }
      const errorMessage =
        apiError?.response?.data?.message ||
        apiError?.response?.data?.error ||
        apiError?.message ||
        `Failed to ${action} device`
      toast.error(errorMessage)
    }
  }

  return (
    <ConfirmDialog
      open={open}
      onOpenChange={onOpenChange}
      handleConfirm={handleToggle}
      isLoading={setDeviceEnabled.isPending}
      title={
        <span>
          <Power className='me-1 inline-block' size={18} />
          {isEnabled ? 'Disable' : 'Enable'} Device
        </span>
      }
      desc={
        <div className='space-y-4'>
          <p className='mb-2'>
            Are you sure you want to {action}{' '}
            <span className='font-bold'>{currentRow.name}</span>?
            <br />
            {isEnabled
              ? 'The device will not be able to communicate with the system while disabled.'
              : 'The device will be able to resume normal operations once enabled.'}
          </p>
        </div>
      }
      confirmText={isEnabled ? 'Disable' : 'Enable'}
      destructive={isEnabled}
    />
  )
}
