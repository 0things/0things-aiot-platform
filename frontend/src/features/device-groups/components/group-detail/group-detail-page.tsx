import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import {
  getDeviceGroupsGroupUuid,
  getDeviceGroupsGroupUuidDevices,
} from '@/api/generated'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Main } from '@/components/layout/main'
import { EditGroupDialog } from './edit-group-dialog'
import { GroupDevicesTab } from './group-devices-tab'
import { GroupHeader } from './group-header'
import { GroupInfoTab } from './group-info-tab'

interface GroupDetailPageProps {
  uuid: string
}

export function GroupDetailPage({ uuid }: GroupDetailPageProps) {
  const { t } = useTranslation('deviceGroup')
  const navigate = useNavigate()
  const [editDialogOpen, setEditDialogOpen] = useState(false)
  const [devicePage, setDevicePage] = useState(1)
  const [devicePageSize, setDevicePageSize] = useState(20)
  const [deviceProductKey, setDeviceProductKey] = useState('')
  const [deviceSearch, setDeviceSearch] = useState('')

  const groupQuery = useQuery({
    queryKey: ['device-group', uuid],
    queryFn: () => getDeviceGroupsGroupUuid(uuid),
    enabled: !!uuid,
  })

  const devicesQuery = useQuery({
    queryKey: [
      'device-group',
      uuid,
      'devices',
      devicePage,
      devicePageSize,
      deviceProductKey,
      deviceSearch,
    ],
    queryFn: () =>
      getDeviceGroupsGroupUuidDevices(uuid, {
        page: devicePage,
        pageSize: devicePageSize,
        productKey: deviceProductKey || undefined,
        search: deviceSearch || undefined,
      }),
    enabled: !!uuid,
  })

  const group = groupQuery.data?.data
  const devices = devicesQuery.data?.data?.devices ?? []
  const totalDevices = devicesQuery.data?.data?.total ?? 0

  const handleBack = () => {
    navigate({ to: '/device-management/groups' })
  }

  // 统计数据聚合为单个对象，避免三字段成群传递
  const deviceStats = {
    total: totalDevices,
    active: devices.filter((d) => d.enabled !== false).length,
    online: devices.filter((d) => d.state === 'online').length,
  }

  if (groupQuery.isLoading) {
    return (
      <Main fixed>
        <div className='flex h-full flex-col'>
          <div className='animate-pulse space-y-4'>
            <div className='h-8 w-48 rounded bg-muted'></div>
            <div className='h-24 rounded bg-muted'></div>
            <div className='h-96 rounded bg-muted'></div>
          </div>
        </div>
      </Main>
    )
  }

  if (groupQuery.error || !group) {
    return (
      <Main fixed>
        <div className='flex h-full flex-col items-center justify-center space-y-4'>
          <p className='text-lg text-muted-foreground'>{t('empty')}</p>
          <Button onClick={handleBack}>{t('common:back')}</Button>
        </div>
      </Main>
    )
  }

  return (
    <Main fixed>
      <div className='flex h-full min-w-0 flex-col'>
        <GroupHeader
          group={group}
          totalDevices={deviceStats.total}
          activeDevices={deviceStats.active}
          onlineDevices={deviceStats.online}
          onBack={handleBack}
        />

        <Tabs defaultValue='info' className='flex min-w-0 flex-1 flex-col'>
          <TabsList className='mb-4 w-fit'>
            <TabsTrigger value='info' className='px-4'>
              {t('info')}
            </TabsTrigger>
            <TabsTrigger value='devices' className='px-4'>
              {t('devices')}
            </TabsTrigger>
          </TabsList>

          <div className='min-w-0 flex-1 overflow-x-hidden overflow-y-auto'>
            <TabsContent value='info' className='mt-0'>
              <GroupInfoTab
                group={group}
                totalDevices={deviceStats.total}
                activeDevices={deviceStats.active}
                onlineDevices={deviceStats.online}
                onEdit={() => setEditDialogOpen(true)}
              />
            </TabsContent>

            <TabsContent value='devices' className='mt-0'>
              <GroupDevicesTab
                group={group}
                devices={devices}
                isLoading={devicesQuery.isLoading}
                onRefresh={() => devicesQuery.refetch()}
                total={totalDevices}
                page={devicePage}
                pageSize={devicePageSize}
                onPageChange={setDevicePage}
                onPageSizeChange={(size) => {
                  setDevicePageSize(size)
                  setDevicePage(1)
                }}
                onSearch={(productKey, search) => {
                  setDeviceProductKey(productKey)
                  setDeviceSearch(search)
                  setDevicePage(1)
                }}
              />
            </TabsContent>
          </div>
        </Tabs>

        <EditGroupDialog
          open={editDialogOpen}
          onOpenChange={setEditDialogOpen}
          group={group}
        />
      </div>
    </Main>
  )
}
