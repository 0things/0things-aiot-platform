import { Button } from '@/components/ui/button'
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

export function IotDashboard() {
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
          <h1 className='text-2xl font-bold tracking-tight'>IoT Dashboard</h1>
          <div className='flex items-center space-x-2'>
            <Button>Refresh</Button>
          </div>
        </div>
        <Tabs
          orientation='vertical'
          defaultValue='overview'
          className='space-y-4'
        >
          <div className='w-full overflow-x-auto pb-2'>
            <TabsList>
              <TabsTrigger value='overview'>Overview</TabsTrigger>
              <TabsTrigger value='devices'>Devices</TabsTrigger>
              <TabsTrigger value='analytics'>Analytics</TabsTrigger>
            </TabsList>
          </div>
          <TabsContent value='overview' className='space-y-4'>
            <div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-3'>
              <Card>
                <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
                  <CardTitle className='text-sm font-medium'>
                    Products
                  </CardTitle>
                  <svg
                    xmlns='http://www.w3.org/2000/svg'
                    viewBox='0 0 24 24'
                    fill='none'
                    stroke='currentColor'
                    strokeLinecap='round'
                    strokeLinejoin='round'
                    strokeWidth='2'
                    className='h-4 w-4 text-muted-foreground'
                  >
                    <path d='M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z' />
                    <polyline points='9 22 9 12 15 12 15 22' />
                  </svg>
                </CardHeader>
                <CardContent>
                  <div className='text-2xl font-bold'>8</div>
                  <p className='text-xs text-muted-foreground'>
                    +2 new products
                  </p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
                  <CardTitle className='text-sm font-medium'>Devices</CardTitle>
                  <svg
                    xmlns='http://www.w3.org/2000/svg'
                    viewBox='0 0 24 24'
                    fill='none'
                    stroke='currentColor'
                    strokeLinecap='round'
                    strokeLinejoin='round'
                    strokeWidth='2'
                    className='h-4 w-4 text-muted-foreground'
                  >
                    <rect x='4' y='4' width='16' height='16' rx='2' />
                    <rect x='9' y='9' width='6' height='6' />
                    <path d='M9 1v6M15 1v6M9 17v6M15 17v6M1 9h6M1 15h6M17 9h6M17 15h6' />
                  </svg>
                </CardHeader>
                <CardContent>
                  <div className='text-2xl font-bold'>1,234</div>
                  <p className='text-xs text-muted-foreground'>
                    +15% from last week
                  </p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
                  <CardTitle className='text-sm font-medium'>
                    Active Devices
                  </CardTitle>
                  <svg
                    xmlns='http://www.w3.org/2000/svg'
                    viewBox='0 0 24 24'
                    fill='none'
                    stroke='currentColor'
                    strokeLinecap='round'
                    strokeLinejoin='round'
                    strokeWidth='2'
                    className='h-4 w-4 text-muted-foreground'
                  >
                    <circle cx='12' cy='12' r='3' />
                    <path d='M12 1v6M12 17v6M4.22 4.22l4.24 4.24M15.54 15.54l4.24 4.24M1 12h6M17 12h6M4.22 19.78l4.24-4.24M15.54 8.46l4.24-4.24' />
                  </svg>
                </CardHeader>
                <CardContent>
                  <div className='text-2xl font-bold'>892</div>
                  <p className='text-xs text-muted-foreground'>
                    72% of total devices
                  </p>
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value='devices' className='space-y-4'>
            <Card>
              <CardHeader>
                <CardTitle>Device Overview</CardTitle>
                <CardDescription>
                  Connected devices and their status
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className='py-8 text-center'>
                  <p className='text-muted-foreground'>
                    Device management interface will be implemented here
                  </p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
          <TabsContent value='analytics' className='space-y-4'>
            <Card>
              <CardHeader>
                <CardTitle>Analytics Dashboard</CardTitle>
                <CardDescription>
                  IoT system analytics and metrics
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className='py-8 text-center'>
                  <p className='text-muted-foreground'>
                    Analytics interface will be implemented here
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

const topNav = [
  {
    title: 'Overview',
    href: 'iot-dashboard/overview',
    isActive: true,
    disabled: false,
  },
  {
    title: 'Devices',
    href: 'iot-dashboard/devices',
    isActive: false,
    disabled: false,
  },
  {
    title: 'Analytics',
    href: 'iot-dashboard/analytics',
    isActive: false,
    disabled: false,
  },
]
