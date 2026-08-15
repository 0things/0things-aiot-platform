'use client'

import { Copy } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import type { DeviceMQTTParametersResponse as DeviceV1GetMqttParametersResponse } from '@/api/generated/model'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

interface MqttParametersDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  mqttParams: DeviceV1GetMqttParametersResponse | null
  isLoading?: boolean
}

export function MqttParametersDialog({
  open,
  onOpenChange,
  mqttParams,
  isLoading,
}: MqttParametersDialogProps) {
  const { t } = useTranslation('deviceManagement')

  const handleCopyAll = () => {
    if (!mqttParams) return

    const text = `clientId: ${mqttParams.clientId || '-'}
username: ${mqttParams.username || '-'}
passwd: ${mqttParams.password || '-'}
mqttHostUrl: ${mqttParams.mqttHostUrl || '-'}
port: ${mqttParams.port || '-'}`

    navigator.clipboard.writeText(text)
    toast.success(t('deviceDetail.info.mqttDialog.copyAllSuccess'))
  }

  const handleCopyField = (
    value: string | number | undefined,
    label: string
  ) => {
    if (!value) return
    navigator.clipboard.writeText(String(value))
    toast.success(
      t('deviceDetail.info.copySuccess', {
        field: label,
        defaultValue: `${label} copied`,
      })
    )
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-2xl'>
        <DialogHeader>
          <DialogTitle>{t('deviceDetail.info.mqttDialog.title')}</DialogTitle>
        </DialogHeader>

        {isLoading ? (
          <div className='flex items-center justify-center py-8'>
            <div className='text-sm text-muted-foreground'>
              {t('deviceDetail.info.mqttDialog.loading')}
            </div>
          </div>
        ) : mqttParams ? (
          <div className='space-y-4'>
            <div className='space-y-3'>
              <div className='flex items-start gap-4'>
                <div className='min-w-32 text-sm text-muted-foreground'>
                  clientId
                </div>
                <div className='flex-1 overflow-hidden'>
                  <div className='flex items-start gap-2'>
                    <p className='font-mono text-sm break-all'>
                      {mqttParams.clientId || '-'}
                    </p>
                    {mqttParams.clientId && (
                      <Button
                        variant='ghost'
                        size='icon'
                        className='h-6 w-6 shrink-0'
                        onClick={() =>
                          handleCopyField(mqttParams.clientId, 'clientId')
                        }
                      >
                        <Copy className='h-3 w-3' />
                      </Button>
                    )}
                  </div>
                </div>
              </div>

              <div className='flex items-start gap-4'>
                <div className='min-w-32 text-sm text-muted-foreground'>
                  username
                </div>
                <div className='flex-1 overflow-hidden'>
                  <div className='flex items-start gap-2'>
                    <p className='font-mono text-sm break-all'>
                      {mqttParams.username || '-'}
                    </p>
                    {mqttParams.username && (
                      <Button
                        variant='ghost'
                        size='icon'
                        className='h-6 w-6 shrink-0'
                        onClick={() =>
                          handleCopyField(mqttParams.username, 'username')
                        }
                      >
                        <Copy className='h-3 w-3' />
                      </Button>
                    )}
                  </div>
                </div>
              </div>

              <div className='flex items-start gap-4'>
                <div className='min-w-32 text-sm text-muted-foreground'>
                  passwd
                </div>
                <div className='flex-1 overflow-hidden'>
                  <div className='flex items-start gap-2'>
                    <p className='font-mono text-sm break-all'>
                      {mqttParams.password || '-'}
                    </p>
                    {mqttParams.password && (
                      <Button
                        variant='ghost'
                        size='icon'
                        className='h-6 w-6 shrink-0'
                        onClick={() =>
                          handleCopyField(mqttParams.password, 'passwd')
                        }
                      >
                        <Copy className='h-3 w-3' />
                      </Button>
                    )}
                  </div>
                </div>
              </div>

              <div className='flex items-start gap-4'>
                <div className='min-w-32 text-sm text-muted-foreground'>
                  mqttHostUrl
                </div>
                <div className='flex-1 overflow-hidden'>
                  <div className='flex items-start gap-2'>
                    <p className='font-mono text-sm break-all'>
                      {mqttParams.mqttHostUrl || '-'}
                    </p>
                    {mqttParams.mqttHostUrl && (
                      <Button
                        variant='ghost'
                        size='icon'
                        className='h-6 w-6 shrink-0'
                        onClick={() =>
                          handleCopyField(mqttParams.mqttHostUrl, 'mqttHostUrl')
                        }
                      >
                        <Copy className='h-3 w-3' />
                      </Button>
                    )}
                  </div>
                </div>
              </div>

              <div className='flex items-start gap-4'>
                <div className='min-w-32 text-sm text-muted-foreground'>
                  port
                </div>
                <div className='flex-1 overflow-hidden'>
                  <div className='flex items-start gap-2'>
                    <p className='font-mono text-sm break-all'>
                      {mqttParams.port || '-'}
                    </p>
                    {mqttParams.port && (
                      <Button
                        variant='ghost'
                        size='icon'
                        className='h-6 w-6 shrink-0'
                        onClick={() => handleCopyField(mqttParams.port, 'port')}
                      >
                        <Copy className='h-3 w-3' />
                      </Button>
                    )}
                  </div>
                </div>
              </div>
            </div>
          </div>
        ) : (
          <div className='flex items-center justify-center py-8'>
            <div className='text-sm text-muted-foreground'>
              {t('deviceDetail.info.mqttDialog.noData')}
            </div>
          </div>
        )}

        <DialogFooter>
          <Button
            variant='default'
            onClick={handleCopyAll}
            disabled={!mqttParams}
          >
            <Copy className='mr-2 h-4 w-4' />
            {t('deviceDetail.info.mqttDialog.copyAll')}
          </Button>
          <Button variant='outline' onClick={() => onOpenChange(false)}>
            {t('deviceDetail.info.mqttDialog.close')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
