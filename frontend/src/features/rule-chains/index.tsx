import { useTranslation } from 'react-i18next'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { RuleChainsPrimaryButtons } from './components/rule-chains-primary-buttons'
import { RuleChainsTable } from './components/rule-chains-table'

export function RuleChains() {
  const { t } = useTranslation('ruleChain')

  return (
    <>
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
            <h2 className='text-2xl font-bold tracking-tight'>
              {t('list.title')}
            </h2>
            <p className='text-muted-foreground'>{t('list.description')}</p>
          </div>
          <RuleChainsPrimaryButtons />
        </div>
        <RuleChainsTable />
      </Main>
    </>
  )
}
