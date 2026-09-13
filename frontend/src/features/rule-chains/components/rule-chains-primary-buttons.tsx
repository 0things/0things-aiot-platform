import { Link } from '@tanstack/react-router'
import { Plus } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'

export function RuleChainsPrimaryButtons() {
  const { t } = useTranslation('ruleChain')

  return (
    <Button asChild>
      <Link to='/rule-engine/rule-chains/new'>
        <Plus data-icon='inline-start' />
        {t('list.create')}
      </Link>
    </Button>
  )
}
