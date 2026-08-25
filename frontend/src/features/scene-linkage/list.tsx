import { useTranslation } from 'react-i18next'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { SceneLinkagePrimaryButtons } from './components/scene-linkage-primary-buttons'
import { SceneLinkageProvider } from './components/scene-linkage-provider'
import { SceneLinkageTable } from './components/scene-linkage-table'

export function SceneLinkageListPage() {
  const { t } = useTranslation('sceneLinkage')

  return (
    <SceneLinkageProvider>
      <Header fixed>
        <Search />
        <div className='ms-auto flex items-center space-x-4'>
          <ThemeSwitch />
          <ConfigDrawer />
          <ProfileDropdown />
        </div>
      </Header>

      <Main fixed className='flex flex-1 flex-col gap-4 sm:gap-6'>
        <div className='flex flex-wrap items-end justify-between gap-2'>
          <div>
            <h2 className='text-2xl font-bold tracking-tight'>{t('title')}</h2>
            <p className='text-muted-foreground'>{t('list.description')}</p>
          </div>
          <SceneLinkagePrimaryButtons />
        </div>
        <SceneLinkageTable />
      </Main>
    </SceneLinkageProvider>
  )
}
