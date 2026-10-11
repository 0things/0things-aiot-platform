'use client'

import { Ban } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { useRevokeOrganizationInvitation } from '../api/queries'
import { type InvitationItem } from '../data/schema'

type UsersRevokeDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  currentInvitation: InvitationItem
}

export function UsersRevokeDialog({
  open,
  onOpenChange,
  currentInvitation,
}: UsersRevokeDialogProps) {
  const { t } = useTranslation('users')
  const revokeMutation = useRevokeOrganizationInvitation()

  const handleRevoke = async () => {
    if (!currentInvitation.id) return
    try {
      await revokeMutation.mutateAsync({ id: currentInvitation.id })
      toast.success(t('revokeSuccess'))
      onOpenChange(false)
    } catch {
      toast.error(t('revokeFailed'))
    }
  }

  return (
    <ConfirmDialog
      open={open}
      onOpenChange={onOpenChange}
      handleConfirm={handleRevoke}
      isLoading={revokeMutation.isPending}
      title={
        <span className='text-destructive'>
          <Ban className='me-1 inline-block stroke-destructive' size={18} />{' '}
          {t('revokeInviteTitle')}
        </span>
      }
      desc={
        <p>
          {t('revokeInviteDesc', {
            email: currentInvitation.invitee,
          })}
        </p>
      }
      confirmText={t('revoke')}
      destructive
    />
  )
}
