import { Link, Outlet, useRouterState } from '@tanstack/react-router'
import { PackageOpen, Smartphone } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'

export function DeviceManagement() {
  const { t } = useTranslation('deviceManagement')
  const router = useRouterState()

  // Determine active tab from current route
  const activeTab = router.location.pathname.includes('/products')
    ? 'products'
    : router.location.pathname.includes('/devices')
      ? 'devices'
      : 'products'

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
            <TabsTrigger value='products' asChild>
              <Link
                to='/device-management/products'
                className='flex items-center gap-2'
              >
                <PackageOpen size={16} />
                {t('tabs.products')}
              </Link>
            </TabsTrigger>
            <TabsTrigger value='devices' asChild>
              <Link
                to='/device-management/devices'
                className='flex items-center gap-2'
              >
                <Smartphone size={16} />
                {t('tabs.devices')}
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
