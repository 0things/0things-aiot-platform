import { useEffect } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useHandleSignInCallback } from '@logto/react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { useAuthStore } from '@/stores/auth-store'

export const Route = createFileRoute('/(auth)/callback')({
  component: LogtoCallback,
})

function LogtoCallback() {
  const { t } = useTranslation('auth')
  const navigate = useNavigate()
  const { auth } = useAuthStore()
  const { isLoading, isAuthenticated, error } = useHandleSignInCallback()

  useEffect(() => {
    if (error) {
      auth.reset()
      toast.error(error.message)
      void navigate({ to: '/sign-in', replace: true })
      return
    }
    if (isLoading) return
    if (!isAuthenticated) {
      void navigate({ to: '/sign-in', replace: true })
      return
    }

    const redirect = sessionStorage.getItem('logto_redirect') || '/'
    sessionStorage.removeItem('logto_redirect')
    const url = new URL(redirect, window.location.origin)
    void navigate({
      to: `${url.pathname}${url.search}${url.hash}`,
      replace: true,
    })
  }, [auth, error, isAuthenticated, isLoading, navigate])

  return (
    <div className='flex min-h-svh items-center justify-center'>
      {t('signIn.loading', { defaultValue: 'Signing in…' })}
    </div>
  )
}
