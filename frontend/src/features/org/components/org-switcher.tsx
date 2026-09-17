import { useEffect, useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { useLogto } from '@logto/react'
import { Building2, Check, ChevronsUpDown, Loader2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { setAuthToken } from '@/api/clients'
import { useAuthStore } from '@/stores/auth-store'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from '@/components/ui/sidebar'

export function OrgSwitcher() {
  const { t } = useTranslation()
  const { isMobile } = useSidebar()
  const { fetchUserInfo, getAccessToken } = useLogto()
  const queryClient = useQueryClient()
  const organizationID = useAuthStore((state) => state.auth.organizationID)
  const setAccessToken = useAuthStore((state) => state.auth.setAccessToken)
  const setOrganizationID = useAuthStore(
    (state) => state.auth.setOrganizationID
  )
  const [organizations, setOrganizations] = useState<
    Array<{ id: string; name: string }>
  >([])
  const [isLoading, setIsLoading] = useState(true)
  const [isSwitching, setIsSwitching] = useState(false)

  useEffect(() => {
    let active = true
    void fetchUserInfo()
      .then((userInfo) => {
        if (!active) return
        const items = Array.isArray(userInfo?.organization_data)
          ? userInfo.organization_data
              .filter(
                (organization) =>
                  typeof organization?.id === 'string' &&
                  typeof organization?.name === 'string'
              )
              .map(({ id, name }) => ({ id, name }))
          : []
        setOrganizations(items)
      })
      .finally(() => {
        if (active) setIsLoading(false)
      })
    return () => {
      active = false
    }
  }, [fetchUserInfo])

  const handleSwitch = async (nextOrganizationID: string) => {
    if (nextOrganizationID === organizationID || isSwitching) return
    setIsSwitching(true)
    try {
      const resource =
        import.meta.env.VITE_LOGTO_RESOURCE || 'http://localhost:8000'
      const token = await getAccessToken(resource, nextOrganizationID)
      if (!token) throw new Error('Organization access token is unavailable')
      setAuthToken(token)
      setAccessToken(token)
      setOrganizationID(nextOrganizationID)
      await queryClient.invalidateQueries()
      toast.success(
        t('switchOrgSuccess', {
          name: nextOrganizationID,
          defaultValue: `Switched to organization: ${nextOrganizationID}`,
        })
      )
    } catch {
      toast.error(
        t('switchOrgFailed', { defaultValue: 'Failed to switch organization' })
      )
    } finally {
      setIsSwitching(false)
    }
  }

  const displayOrganizationID =
    organizationID ||
    (isLoading
      ? t('loading', { defaultValue: 'Loading...' })
      : t('noOrganizations', { defaultValue: 'No organizations' }))
  const activeOrganization = organizations.find(
    (organization) => organization.id === organizationID
  )
  const displayOrganizationName =
    activeOrganization?.name || displayOrganizationID

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size='lg'
              disabled={isLoading || isSwitching || organizations.length < 2}
              className='data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground'
            >
              <div className='flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground'>
                {isSwitching ? (
                  <Loader2 className='size-4 animate-spin' />
                ) : (
                  <Building2 className='size-4' />
                )}
              </div>
              <div className='grid flex-1 text-start text-sm leading-tight'>
                <span className='truncate font-semibold'>
                  {displayOrganizationName}
                </span>
                <span className='truncate text-xs text-muted-foreground'>
                  {displayOrganizationID}
                </span>
              </div>
              {organizations.length > 1 && (
                <ChevronsUpDown className='ms-auto' />
              )}
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className='w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg'
            align='start'
            side={isMobile ? 'bottom' : 'right'}
            sideOffset={4}
          >
            <DropdownMenuLabel className='text-xs text-muted-foreground'>
              {t('switchOrg', { defaultValue: 'Switch Organization' })}
            </DropdownMenuLabel>
            {organizations.map((organization) => (
              <DropdownMenuItem
                key={organization.id}
                onClick={() => void handleSwitch(organization.id)}
                className='flex items-center justify-between gap-2 p-2'
              >
                <div className='flex items-center gap-2 truncate'>
                  <div className='flex size-6 shrink-0 items-center justify-center rounded-sm border'>
                    <Building2 className='size-3.5' />
                  </div>
                  <span className='truncate font-medium'>
                    {organization.name}
                  </span>
                </div>
                {organization.id === organizationID && (
                  <Check className='size-4 shrink-0 text-primary' />
                )}
              </DropdownMenuItem>
            ))}
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
