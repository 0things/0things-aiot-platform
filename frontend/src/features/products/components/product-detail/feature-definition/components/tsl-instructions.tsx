import { useTranslation } from 'react-i18next'

export function TslInstructions() {
  const { t } = useTranslation('deviceManagement')

  return (
    <div className='rounded-lg border bg-muted/30 p-4'>
      <h4 className='mb-2 text-sm font-semibold'>
        {t('productDetail.featureDefinition.instructions.title')}
      </h4>
      <ul className='list-disc space-y-1 pl-5 text-xs text-muted-foreground'>
        <li>{t('productDetail.featureDefinition.instructions.properties')}</li>
        <li>{t('productDetail.featureDefinition.instructions.events')}</li>
        <li>{t('productDetail.featureDefinition.instructions.services')}</li>
        <li>{t('productDetail.featureDefinition.instructions.dualMode')}</li>
        <li>{t('productDetail.featureDefinition.instructions.visualMode')}</li>
      </ul>
    </div>
  )
}
