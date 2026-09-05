import { Route } from '@/routes/_authenticated/device-management/devices/$deviceKey'
import { useTranslation } from 'react-i18next'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { DeviceEvents } from '@/features/operations-monitoring/events'
import { PropertyTab } from './property-tab'
import { ServiceInvocationsTab } from './service-invocations-tab'

export function ThingModelDataTab() {
  const { t } = useTranslation('deviceManagement')
  const { deviceKey } = Route.useParams()

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
        <PropertyTab />
      </TabsContent>
      <TabsContent value='eventManagement' className='mt-0'>
        <DeviceEvents deviceKey={deviceKey} />
      </TabsContent>
      <TabsContent value='serviceInvocation' className='mt-0'>
        <ServiceInvocationsTab deviceKey={deviceKey} />
      </TabsContent>
    </Tabs>
  )
}
