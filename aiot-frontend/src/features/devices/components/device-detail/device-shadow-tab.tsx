import { useTranslation } from 'react-i18next'
import { Card, CardContent } from '@/components/ui/card'

export function DeviceShadowTab() {
  const { t } = useTranslation('deviceManagement')

  return (
    <Card>
      <CardContent className='flex min-h-[400px] items-center justify-center'>
        <div className='text-center'>
          <p className='text-lg text-muted-foreground'>
            {t('deviceDetail.deviceShadowTab.comingSoon')}
          </p>
        </div>
      </CardContent>
    </Card>
  )
}
