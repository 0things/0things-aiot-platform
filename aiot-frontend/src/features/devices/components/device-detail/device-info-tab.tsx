import { useState } from 'react'
import { Copy } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import type { Device as DeviceV1Device } from '@/api/generated/model'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Switch } from '@/components/ui/switch'
import { useDeviceMqttParameters } from '../../api/queries'
import { MqttParametersDialog } from './mqtt-parameters-dialog'

interface DeviceInfoTabProps {
  device: DeviceV1Device
}

export function DeviceInfoTab({ device }: DeviceInfoTabProps) {
  const { t } = useTranslation('deviceManagement')
  const [logUploadEnabled, setLogUploadEnabled] = useState(false)
  const [mqttDialogOpen, setMqttDialogOpen] = useState(false)

  const {
    data: mqttParams,
    isLoading: isMqttParamsLoading,
    refetch: refetchMqttParams,
  } = useDeviceMqttParameters(device.deviceKey || '', mqttDialogOpen)

  const handleCopy = (text: string, label: string) => {
    navigator.clipboard.writeText(text)
    toast.success(
      t('deviceDetail.info.copySuccess', {
        field: label,
        defaultValue: 'Copied',
      })
    )
  }

  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return '-'
    const date = new Date(dateString)
    return date.toISOString().slice(0, 19).replace('T', ' ')
  }

  const getStatusVariant = () => {
    if (device.state === 'online') return 'default'
    if (device.state === 'offline') return 'secondary'
    return 'destructive'
  }

  const handleLogUploadToggle = (checked: boolean) => {
    setLogUploadEnabled(checked)
    toast.success(
      checked
        ? t('deviceDetail.info.logUploadEnabled')
        : t('deviceDetail.info.logUploadDisabled')
    )
  }

  const handleTestLatency = () => {
    toast.info(t('deviceDetail.info.testingLatency'))
    // Placeholder for latency test
    setTimeout(() => {
      toast.success(t('deviceDetail.info.latencyResult', { latency: '45ms' }))
    }, 1000)
  }

  const handleViewMqttParams = () => {
    setMqttDialogOpen(true)
    if (!mqttParams) {
      refetchMqttParams()
    }
  }

  return (
    <div className='space-y-6'>
      {/* Main Device Information Card */}
      <Card>
        <CardHeader>
          <CardTitle>{t('deviceDetail.info.title')}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className='grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3'>
            {/* Left Column */}
            <div className='space-y-4'>
              <div className='space-y-1'>
                <p className='text-sm text-muted-foreground'>
                  {t('deviceDetail.info.fields.productName')}
                </p>
                <p className='font-medium'>{device.productName || '-'}</p>
              </div>

              <div className='space-y-1'>
                <p className='text-sm text-muted-foreground'>
                  {t('deviceDetail.info.fields.notes')}
                </p>
                <div className='flex items-center gap-2'>
                  <p className='font-medium'>{device.metadata || '-'}</p>
                  <Button variant='link' size='sm' className='h-auto p-0'>
                    {t('deviceDetail.info.edit')}
                  </Button>
                </div>
              </div>

              <div className='space-y-1'>
                <p className='text-sm text-muted-foreground'>
                  {t('deviceDetail.info.fields.createdAt')}
                </p>
                <p className='font-medium'>{formatDate(device.createdAt)}</p>
              </div>

              <div className='space-y-1'>
                <p className='text-sm text-muted-foreground'>
                  {t('deviceDetail.info.fields.currentStatus')}
                </p>
                <Badge variant={getStatusVariant()}>
                  {t(`devices.state.${device.state}`, {
                    defaultValue: device.state || '-',
                  })}
                </Badge>
              </div>

              <div className='space-y-1'>
                <p className='text-sm text-muted-foreground'>
                  {t('deviceDetail.info.fields.mqttParams')}
                </p>
                <Button
                  variant='link'
                  size='sm'
                  className='h-auto p-0'
                  onClick={handleViewMqttParams}
                >
                  {t('deviceDetail.header.view')}
                </Button>
              </div>
            </div>

            {/* Middle Column */}
            <div className='space-y-4'>
              <div className='space-y-1'>
                <p className='text-sm text-muted-foreground'>ProductKey</p>
                <div className='flex items-center gap-2'>
                  <p className='font-mono text-sm'>
                    {device.productKey || '-'}
                  </p>
                  {device.productKey && (
                    <Button
                      variant='ghost'
                      size='icon'
                      className='h-6 w-6'
                      onClick={() =>
                        handleCopy(device.productKey || '', 'ProductKey')
                      }
                    >
                      <Copy className='h-3 w-3' />
                    </Button>
                  )}
                </div>
              </div>

              <div className='space-y-1'>
                <p className='text-sm text-muted-foreground'>DeviceKey</p>
                <div className='flex items-center gap-2'>
                  <p className='font-mono text-sm'>{device.deviceKey || '-'}</p>
                  {device.deviceKey && (
                    <Button
                      variant='ghost'
                      size='icon'
                      className='h-6 w-6'
                      onClick={() =>
                        handleCopy(device.deviceKey || '', 'DeviceKey')
                      }
                    >
                      <Copy className='h-3 w-3' />
                    </Button>
                  )}
                </div>
              </div>

              <div className='space-y-1'>
                <p className='text-sm text-muted-foreground'>
                  {t('deviceDetail.info.fields.ipAddress')}
                </p>
                <p className='font-medium'>-</p>
              </div>

              <div className='space-y-1'>
                <p className='text-sm text-muted-foreground'>
                  {t('deviceDetail.info.fields.activationTime')}
                </p>
                <p className='font-medium'>-</p>
              </div>

              <div className='space-y-1'>
                <p className='text-sm text-muted-foreground'>
                  {t('deviceDetail.info.fields.realTimeLatency')}
                </p>
                <Button variant='outline' size='sm' onClick={handleTestLatency}>
                  {t('deviceDetail.info.test')}
                </Button>
              </div>

              <div className='space-y-1'>
                <p className='text-sm text-muted-foreground'>
                  {t('deviceDetail.info.fields.lastOfflineTime')}
                </p>
                <p className='font-medium'>
                  {device.lastOfflineTime
                    ? formatDate(String(device.lastOfflineTime))
                    : '-'}
                </p>
              </div>
            </div>

            {/* Right Column */}
            <div className='space-y-4'>
              <div className='space-y-1'>
                <p className='text-sm text-muted-foreground'>
                  {t('deviceDetail.info.fields.authMethod')}
                </p>
                <p className='font-medium'>-</p>
              </div>

              <div className='space-y-1'>
                <p className='text-sm text-muted-foreground'>
                  {t('deviceDetail.info.fields.firmwareVersion')}
                </p>
                <p className='font-medium'>-</p>
              </div>

              <div className='space-y-1'>
                <p className='text-sm text-muted-foreground'>
                  {t('deviceDetail.info.fields.lastOnlineTime')}
                </p>
                <p className='font-medium'>
                  {device.lastOnlineTime ?? '-'}
                </p>
              </div>

              <div className='space-y-1'>
                <p className='text-sm text-muted-foreground'>
                  {t('deviceDetail.info.fields.deviceLogUpload')}
                </p>
                <Switch
                  checked={logUploadEnabled}
                  onCheckedChange={handleLogUploadToggle}
                />
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Device Extended Information Card */}
      <Card>
        <CardHeader>
          <CardTitle>{t('deviceDetail.extendedInfo.title')}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className='grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4'>
            <div className='space-y-1'>
              <p className='text-sm text-muted-foreground'>
                {t('deviceDetail.extendedInfo.fields.sdkLanguage')}
              </p>
              <p className='font-medium'>-</p>
            </div>
            <div className='space-y-1'>
              <p className='text-sm text-muted-foreground'>
                {t('deviceDetail.extendedInfo.fields.version')}
              </p>
              <p className='font-medium'>-</p>
            </div>
            <div className='space-y-1'>
              <p className='text-sm text-muted-foreground'>
                {t('deviceDetail.extendedInfo.fields.moduleVendor')}
              </p>
              <p className='font-medium'>-</p>
            </div>
            <div className='space-y-1'>
              <p className='text-sm text-muted-foreground'>
                {t('deviceDetail.extendedInfo.fields.moduleInfo')}
              </p>
              <p className='font-medium'>-</p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Tag Information Card */}
      <Card>
        <CardHeader className='flex flex-row items-center justify-between'>
          <CardTitle>{t('deviceDetail.tags.title')}</CardTitle>
          <Button variant='outline' size='sm'>
            {t('deviceDetail.info.edit')}
          </Button>
        </CardHeader>
        <CardContent>
          <p className='text-sm text-muted-foreground'>
            {t('deviceDetail.tags.empty')}
          </p>
        </CardContent>
      </Card>

      {/* MQTT Parameters Dialog */}
      <MqttParametersDialog
        open={mqttDialogOpen}
        onOpenChange={setMqttDialogOpen}
        mqttParams={mqttParams || null}
        isLoading={isMqttParamsLoading}
      />
    </div>
  )
}
