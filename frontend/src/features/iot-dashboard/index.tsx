import { useTranslation } from 'react-i18next'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { TopNav } from '@/components/layout/top-nav'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { IotOverview } from './components/iot-overview'

export function IotDashboard() {
  const { t } = useTranslation('iotDashboard')

  const topNav = [
    {
      title: t('topNav.overview'),
      href: 'iot-dashboard/overview',
      isActive: true,
      disabled: false,
    },
    {
      title: t('topNav.devices'),
      href: 'iot-dashboard/devices',
      isActive: false,
      disabled: false,
    },
    {
      title: t('topNav.analytics'),
      href: 'iot-dashboard/analytics',
      isActive: false,
      disabled: false,
    },
  ]

  return (
    <>
      {/* ===== Top Heading ===== */}
      <Header>
        <TopNav links={topNav} />
        <div className='ms-auto flex items-center space-x-4'>
          <Search />
          <ThemeSwitch />
          <ConfigDrawer />
          <ProfileDropdown />
        </div>
      </Header>

      {/* ===== Main ===== */}
      <Main>
        <div className='mb-2 flex items-center justify-between space-y-2'>
          <h1 className='text-2xl font-bold tracking-tight'>
            {t('tabs.overview')}
          </h1>
        </div>
        <Tabs
          orientation='vertical'
          defaultValue='overview'
          className='space-y-4'
        >
          <div className='w-full overflow-x-auto pb-2'>
            <TabsList>
              <TabsTrigger value='overview'>{t('tabs.overview')}</TabsTrigger>
              <TabsTrigger value='devices'>{t('tabs.devices')}</TabsTrigger>
              <TabsTrigger value='analytics'>{t('tabs.analytics')}</TabsTrigger>
            </TabsList>
          </div>
          <TabsContent value='overview' className='space-y-4'>
            <IotOverview />
          </TabsContent>

          <TabsContent value='devices' className='space-y-4'>
            <Card>
              <CardHeader>
                <CardTitle>{t('placeholders.devicesTitle')}</CardTitle>
                <CardDescription>
                  {t('placeholders.devicesDescription')}
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className='py-8 text-center'>
                  <p className='text-muted-foreground'>
                    {t('placeholders.devicesBody')}
                  </p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value='analytics' className='space-y-4'>
            <Card>
              <CardHeader>
                <CardTitle>{t('placeholders.analyticsTitle')}</CardTitle>
                <CardDescription>
                  {t('placeholders.analyticsDescription')}
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className='py-8 text-center'>
                  <p className='text-muted-foreground'>
                    {t('placeholders.analyticsBody')}
                  </p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </Main>
    </>
  )
}
