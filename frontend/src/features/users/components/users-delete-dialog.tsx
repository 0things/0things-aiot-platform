'use client'

import { useState } from 'react'
import { AlertTriangle } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { useRemoveOrganizationMember } from '../api/queries'
import { type User } from '../data/schema'

type UserDeleteDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  currentRow: User
}

export function UsersDeleteDialog({
  open,
  onOpenChange,
  currentRow,
}: UserDeleteDialogProps) {
  const { t } = useTranslation('users')
  const [value, setValue] = useState('')
  const removeMutation = useRemoveOrganizationMember()

  const handleDelete = async () => {
    if (value.trim() !== currentRow.username) return
    try {
      await removeMutation.mutateAsync({ userId: currentRow.id })
      toast.success(t('removeMemberSuccess'))
      onOpenChange(false)
    } catch {
      toast.error(t('removeMemberFailed'))
    }
  }

  return (
    <ConfirmDialog
      open={open}
      onOpenChange={onOpenChange}
      handleConfirm={handleDelete}
      disabled={
        value.trim() !== currentRow.username || removeMutation.isPending
      }
      isLoading={removeMutation.isPending}
      title={
        <span className='text-destructive'>
          <AlertTriangle
            className='me-1 inline-block stroke-destructive'
            size={18}
          />{' '}
          {t('removeMemberTitle')}
        </span>
      }
      desc={
        <div className='space-y-4'>
          <p className='mb-2'>
            {t('removeMemberDesc', { name: currentRow.username })}
          </p>

          <Label className='my-2'>
            {t('confirmUsernameLabel')}
            <Input
              value={value}
              onChange={(e) => setValue(e.target.value)}
              placeholder={t('confirmUsernamePlaceholder')}
            />
          </Label>

          <Alert variant='destructive'>
            <AlertTitle>{t('warning')}</AlertTitle>
            <AlertDescription>{t('warningDesc')}</AlertDescription>
          </Alert>
        </div>
      }
      confirmText={t('delete')}
      destructive
    />
  )
}
