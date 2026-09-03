import { useTranslation } from 'react-i18next'
import { Card, CardContent } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { TelemetryTab } from './telemetry-tab'

export function ThingModelDataTab() {
  const { t } = useTranslation('deviceManagement')

  return (
    <Tabs defaultValue='properties' className='gap-4'>
      <TabsList className='mb-1 w-fit'>
        <TabsTrigger value='properties' className='px-4'>
          {t('deviceDetail.thingModelDataTab.properties')}
        </TabsTrigger>
        <TabsTrigger value='eventManagement' className='px-4'>
          {t('deviceDetail.thingModelDataTab.eventManagement')}
        </TabsTrigger>
        <TabsTrigger value='serviceInvocation' className='px-4'>
          {t('deviceDetail.thingModelDataTab.serviceInvocation')}
        </TabsTrigger>
      </TabsList>

      <TabsContent value='properties' className='mt-0'>
        <TelemetryTab />
      </TabsContent>
      <TabsContent value='eventManagement' className='mt-0'>
        <ThingModelSubPage
          title={t('deviceDetail.thingModelDataTab.eventManagement')}
          description={t(
            'deviceDetail.thingModelDataTab.eventManagementComingSoon'
          )}
        />
      </TabsContent>
      <TabsContent value='serviceInvocation' className='mt-0'>
        <ThingModelSubPage
          title={t('deviceDetail.thingModelDataTab.serviceInvocation')}
          description={t(
            'deviceDetail.thingModelDataTab.serviceInvocationComingSoon'
          )}
        />
      </TabsContent>
    </Tabs>
  )
}

function ThingModelSubPage({
  title,
  description,
}: {
  title: string
  description: string
}) {
  return (
    <Card>
      <CardContent className='flex min-h-[400px] items-center justify-center'>
        <div className='text-center'>
          <p className='text-lg text-muted-foreground'>{title}</p>
          <p className='mt-2 text-sm text-muted-foreground'>{description}</p>
        </div>
      </CardContent>
    </Card>
  )
}
