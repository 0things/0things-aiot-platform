import { Link } from '@tanstack/react-router'
import { Plus } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'

export function SceneLinkagePrimaryButtons() {
  const { t } = useTranslation('sceneLinkage')
  return (
    <div className='flex gap-2'>
      <Button asChild className='space-x-1'>
        <Link to='/rule-engine/scene-linkage/new'>
          <span>{t('list.create')}</span> <Plus size={18} />
        </Link>
      </Button>
    </div>
  )
}
