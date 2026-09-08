import { Code, FileJson, Shield, Zap } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Badge } from '@/components/ui/badge'

export function TslHeader() {
  const { t } = useTranslation('deviceManagement')

  return (
    <div className='rounded-lg border bg-muted/50 p-4'>
      <div className='flex items-start space-x-3'>
        <FileJson className='mt-1 h-5 w-5 text-primary' />
        <div className='space-y-1'>
          <h3 className='text-sm font-semibold'>
            {t('productDetail.featureDefinition.title')}
          </h3>
          <p className='text-xs text-muted-foreground'>
            {t('productDetail.featureDefinition.description')}
          </p>
          <div className='flex flex-wrap gap-2 pt-2'>
            <Badge variant='outline' className='text-xs'>
              <Code className='mr-1 h-3 w-3' />
              {t('productDetail.featureDefinition.badge.jsonFormat')}
            </Badge>
            <Badge variant='outline' className='text-xs'>
              <Zap className='mr-1 h-3 w-3' />
              {t('productDetail.featureDefinition.badge.realtimeValidation')}
            </Badge>
            <Badge variant='outline' className='text-xs'>
              <Shield className='mr-1 h-3 w-3' />
              {t('productDetail.featureDefinition.badge.standardized')}
            </Badge>
          </div>
        </div>
      </div>
    </div>
  )
}
