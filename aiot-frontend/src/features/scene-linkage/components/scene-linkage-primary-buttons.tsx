import { Plus } from 'lucide-react'
import { Link } from '@tanstack/react-router'
import { Button } from '@/components/ui/button'
import { useTranslation } from 'react-i18next'

export function SceneLinkagePrimaryButtons() {
  const { t } = useTranslation('sceneLinkage')
  return (
    <div className='flex gap-2'>
      <Button asChild className='space-x-1'>
        <Link
          to='/rule-engine/scene-linkage/$sceneId'
          params={{ sceneId: 'new' }}
        >
          <span>{t('list.create')}</span> <Plus size={18} />
        </Link>
      </Button>
    </div>
  )
}
