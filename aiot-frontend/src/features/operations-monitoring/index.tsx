import { Link, Outlet, useRouterState } from '@tanstack/react-router'
import { Upload, BarChart3 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'

export function OperationsMonitoring() {
  const { t } = useTranslation('operationsMonitoring')
  const router = useRouterState()

  const activeTab = router.location.pathname.includes('/analytics')
    ? 'analytics'
    : router.location.pathname.includes('/packages')
      ? 'packages'
      : 'packages'

  return (
    <>
      <Header>
        <Search />
        <div className='ms-auto flex items-center space-x-4'>
          <ThemeSwitch />
          <ConfigDrawer />
          <ProfileDropdown />
        </div>
      </Header>

      <Main fixed>
        <Tabs value={activeTab} className='space-y-4'>
          <TabsList>
            <TabsTrigger value='packages' asChild>
              <Link
                to='/operations-monitoring/ota/packages'
                className='flex items-center gap-2'
              >
                <Upload size={16} />
                {t('ota.tabs.packages')}
              </Link>
            </TabsTrigger>
            <TabsTrigger value='analytics' asChild>
              <Link
                to='/operations-monitoring/ota/analytics'
                className='flex items-center gap-2'
              >
                <BarChart3 size={16} />
                {t('ota.tabs.analytics')}
              </Link>
            </TabsTrigger>
          </TabsList>

          <div className='flex flex-1 flex-col overflow-hidden'>
            <Outlet />
          </div>
        </Tabs>
      </Main>
    </>
  )
}
