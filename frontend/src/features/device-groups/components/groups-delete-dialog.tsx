import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { deleteDeviceGroupsGroupUuid } from '@/api/generated'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { type DeviceGroup } from '../data/schema'

interface GroupsDeleteDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  currentRow: DeviceGroup | null
}

export function GroupsDeleteDialog({
  open,
  onOpenChange,
  currentRow,
}: GroupsDeleteDialogProps) {
  const { t } = useTranslation('deviceGroup')
  const { t: tCommon } = useTranslation('common')
  const queryClient = useQueryClient()

  const deleteMutation = useMutation({
    mutationFn: () => deleteDeviceGroupsGroupUuid(currentRow?.groupUuid || ''),
    onSuccess: () => {
      toast.success(tCommon('success'))
      queryClient.invalidateQueries({ queryKey: ['device-groups'] })
      onOpenChange(false)
    },
    onError: (err: unknown) => {
      const msg =
        (err as { response?: { data?: { message?: string } } })?.response?.data
          ?.message ||
        (err as Error)?.message ||
        'Delete failed'
      toast.error(msg)
    },
  })

  return (
    <ConfirmDialog
      open={open}
      onOpenChange={onOpenChange}
      handleConfirm={() => deleteMutation.mutate()}
      disabled={deleteMutation.isPending}
      title={tCommon('delete')}
      desc={t('deleteConfirm')}
      confirmText={tCommon('delete')}
      cancelBtnText={tCommon('cancel')}
      destructive
    />
  )
}
