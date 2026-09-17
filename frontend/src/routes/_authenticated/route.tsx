import { useEffect, useRef, useState } from 'react'
import { createFileRoute, useLocation } from '@tanstack/react-router'
import { useLogto } from '@logto/react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { setAuthToken } from '@/api/clients'
import { useAuthStore } from '@/stores/auth-store'
import { AuthenticatedLayout } from '@/components/layout/authenticated-layout'
import { ensureCurrentUserOrganization } from '@/features/org/api/organization'

export const Route = createFileRoute('/_authenticated')({
  component: AuthenticatedRoute,
})

function extractAvailableOrganizations(
  userInfo: unknown,
  claims: unknown
): string[] {
  const info = userInfo as {
    organization_data?: Array<{ id?: string }>
    organizations?: string[]
  } | null
  const cl = claims as { organizations?: string[] } | null

  const orgsFromInfo = Array.isArray(info?.organization_data)
    ? info.organization_data
        .map((org) => org?.id)
        .filter((id): id is string => typeof id === 'string' && id.length > 0)
    : Array.isArray(info?.organizations)
      ? info.organizations.filter(
          (id): id is string => typeof id === 'string' && id.length > 0
        )
      : []

  const orgsFromClaims = Array.isArray(cl?.organizations)
    ? cl.organizations.filter(
        (id): id is string => typeof id === 'string' && id.length > 0
      )
    : []

  return Array.from(new Set([...orgsFromInfo, ...orgsFromClaims]))
}

function AuthenticatedRoute() {
  const { t } = useTranslation('auth')
  const location = useLocation()
  const {
    isAuthenticated,
    isLoading,
    signIn,
    getAccessToken,
    fetchUserInfo,
    getIdTokenClaims,
  } = useLogto()
  const accessToken = useAuthStore((state) => state.auth.accessToken)
  const setAccessToken = useAuthStore((state) => state.auth.setAccessToken)
  const organizationID = useAuthStore((state) => state.auth.organizationID)
  const setOrganizationID = useAuthStore(
    (state) => state.auth.setOrganizationID
  )
  const resetAuth = useAuthStore((state) => state.auth.reset)
  const bootstrapping = useRef(false)
  const [bootstrapError, setBootstrapError] = useState<Error | null>(null)

  useEffect(() => {
    if (isLoading && !isAuthenticated) return
    if (!isAuthenticated) {
      sessionStorage.setItem(
        'logto_redirect',
        `${location.pathname}${location.search}${location.hash}`
      )
      void signIn(`${window.location.origin}/callback`).catch(
        (error: unknown) => {
          toast.error(error instanceof Error ? error.message : 'Sign-in failed')
        }
      )
      return
    }
    if (organizationID && accessToken) return
    if (bootstrapError || bootstrapping.current) return
    bootstrapping.current = true

    const resource =
      import.meta.env.VITE_LOGTO_RESOURCE || 'http://localhost:8000'

    async function bootstrapTenancy() {
      // 1. Check if user already has assigned organizations in Logto
      const [userInfo, claims] = await Promise.all([
        fetchUserInfo().catch(() => null),
        getIdTokenClaims().catch(() => null),
      ])

      const availableOrgs = extractAvailableOrganizations(userInfo, claims)
      let targetOrgId = organizationID

      if (availableOrgs.length > 0) {
        if (!targetOrgId || !availableOrgs.includes(targetOrgId)) {
          targetOrgId = availableOrgs[0]
        }
      } else {
        // 2. User has no organizations -> bootstrap via un-scoped identity token
        const identityToken = await getAccessToken(resource)
        if (!identityToken) {
          throw new Error('Logto access token is unavailable')
        }
        targetOrgId = await ensureCurrentUserOrganization(identityToken)
      }

      if (!targetOrgId) {
        throw new Error('Failed to resolve organization identity')
      }

      // 3. Acquire organization-scoped access token from Logto
      const organizationToken = await getAccessToken(resource, targetOrgId)
      if (!organizationToken) {
        throw new Error('Logto organization token is unavailable')
      }

      setAuthToken(organizationToken)
      setAccessToken(organizationToken)
      setOrganizationID(targetOrgId)
    }

    void bootstrapTenancy().catch((error: unknown) => {
      bootstrapping.current = false
      const normalizedError =
        error instanceof Error
          ? error
          : new Error('Organization initialization failed')
      resetAuth()
      setBootstrapError(normalizedError)
      toast.error(normalizedError.message)
    })
  }, [
    accessToken,
    bootstrapError,
    fetchUserInfo,
    getAccessToken,
    getIdTokenClaims,
    isAuthenticated,
    isLoading,
    location.hash,
    location.pathname,
    location.search,
    organizationID,
    resetAuth,
    setAccessToken,
    setOrganizationID,
    signIn,
  ])

  if (bootstrapError) {
    return (
      <div className='flex min-h-svh items-center justify-center'>
        {bootstrapError.message}
      </div>
    )
  }

  if ((isLoading && !isAuthenticated) || !isAuthenticated || !accessToken) {
    return (
      <div className='flex min-h-svh items-center justify-center'>
        {t('signIn.loading')}
      </div>
    )
  }

  return <AuthenticatedLayout />
}
