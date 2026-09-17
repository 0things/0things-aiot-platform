import { useNavigate } from '@tanstack/react-router'
import { useLogto } from '@logto/react'
import { useTranslation } from 'react-i18next'
import { useAuthStore } from '@/stores/auth-store'
import { ConfirmDialog } from '@/components/confirm-dialog'

interface SignOutDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function SignOutDialog({ open, onOpenChange }: SignOutDialogProps) {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { auth } = useAuthStore()
  const { signOut } = useLogto()

  const handleSignOut = () => {
    auth.reset()
    void signOut(`${window.location.origin}/`).catch(() => {
      navigate({ to: '/sign-in', replace: true })
    })
  }

  return (
    <ConfirmDialog
      open={open}
      onOpenChange={onOpenChange}
      title={t('signOutConfirmTitle', { defaultValue: 'Sign out' })}
      desc={t('signOutConfirmDesc', {
        defaultValue:
          'Are you sure you want to sign out? You will need to sign in again to access your account.',
      })}
      confirmText={t('signOut', { defaultValue: 'Sign out' })}
      cancelBtnText={t('cancel', { defaultValue: 'Cancel' })}
      destructive
      handleConfirm={handleSignOut}
      className='sm:max-w-sm'
    />
  )
}
