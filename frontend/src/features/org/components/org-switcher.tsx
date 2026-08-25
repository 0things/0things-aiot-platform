import * as React from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { Building2, Check, ChevronsUpDown, Loader2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { setAuthToken } from '@/api/clients'
import { useAuthStore } from '@/stores/auth-store'
import { useOrgStore } from '@/stores/org-store'
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
import { useGetOrganizations, usePostAuthSwitchOrg } from '@/features/org/api'

export function OrgSwitcher() {
  const { t } = useTranslation()
  const { isMobile } = useSidebar()
  const queryClient = useQueryClient()
  const accessToken = useAuthStore((s) => s.auth.accessToken)
  const setAccessToken = useAuthStore((s) => s.auth.setAccessToken)
  const { currentOrgId, organizations, setCurrentOrgId, setOrganizations } =
    useOrgStore()

  const { data: orgsData, isLoading } = useGetOrganizations({
    query: {
      enabled: !!accessToken,
    },
  })
  const { mutateAsync: switchOrg, isPending: isSwitching } =
    usePostAuthSwitchOrg()

  React.useEffect(() => {
    if (orgsData?.data && Array.isArray(orgsData.data)) {
      setOrganizations(orgsData.data)
    }
  }, [orgsData, setOrganizations])

  const activeOrg = React.useMemo(() => {
    if (currentOrgId && organizations.length > 0) {
      const found = organizations.find((o) => o.id === currentOrgId)
      if (found) return found
    }
    return organizations.find((o) => o.is_current) || organizations[0] || null
  }, [currentOrgId, organizations])

  const handleSwitch = async (orgId: number, orgName: string) => {
    if (orgId === activeOrg?.id) return
    try {
      const res = await switchOrg({ data: { org_id: orgId } })
      const newToken = res.data?.accessToken
      if (newToken) {
        setAccessToken(newToken)
        setAuthToken(newToken)
        setCurrentOrgId(orgId)
        await queryClient.invalidateQueries()
        toast.success(
          t('switchOrgSuccess', {
            name: orgName,
            defaultValue: `已切换到组织: ${orgName}`,
          })
        )
      }
    } catch {
      toast.error(t('switchOrgFailed', { defaultValue: '切换组织失败' }))
    }
  }

  const orgDisplayName =
    activeOrg?.name ||
    (isLoading
      ? t('loading', { defaultValue: '加载中...' })
      : t('noOrganizations', { defaultValue: '暂无组织' }))

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size='lg'
              className='data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground'
              disabled={isLoading || isSwitching}
            >
              <div className='flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground'>
                {isSwitching ? (
                  <Loader2 className='size-4 animate-spin' />
                ) : (
                  <Building2 className='size-4' />
                )}
              </div>
              <div className='grid flex-1 text-start text-sm leading-tight'>
                <span className='truncate font-semibold'>{orgDisplayName}</span>
                <span className='truncate text-xs text-muted-foreground'>
                  {t('organizations', { defaultValue: '组织' })}
                </span>
              </div>
              <ChevronsUpDown className='ms-auto' />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className='w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg'
            align='start'
            side={isMobile ? 'bottom' : 'right'}
            sideOffset={4}
          >
            <DropdownMenuLabel className='text-xs text-muted-foreground'>
              {t('organizations', { defaultValue: '组织列表' })}
            </DropdownMenuLabel>
            {organizations.map((org) => {
              const isSelected = org.id === activeOrg?.id
              return (
                <DropdownMenuItem
                  key={org.id}
                  onClick={() => org.id && handleSwitch(org.id, org.name || '')}
                  className='flex items-center justify-between gap-2 p-2'
                >
                  <div className='flex items-center gap-2 truncate'>
                    <div className='flex size-6 shrink-0 items-center justify-center rounded-sm border'>
                      <Building2 className='size-3.5' />
                    </div>
                    <span className='truncate font-medium'>{org.name}</span>
                  </div>
                  {isSelected && (
                    <Check className='size-4 shrink-0 text-primary' />
                  )}
                </DropdownMenuItem>
              )
            })}
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
