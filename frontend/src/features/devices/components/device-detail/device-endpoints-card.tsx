import type { ReactNode } from 'react'
import { Copy, Globe, Radio, Wifi } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useDeviceEndpoints } from '../../api/endpoints'

export function DeviceEndpointsCard({ deviceKey }: { deviceKey: string }) {
  const { t } = useTranslation('deviceManagement')
  const { data: connectionData = {}, isLoading } = useDeviceEndpoints(deviceKey)
  const copy = async (value: string, label: string) => {
    await navigator.clipboard.writeText(value)
    toast.success(
      t('deviceDetail.header.copySuccess', { defaultValue: `${label} copied` })
    )
  }

  return (
    <Card>
      <CardHeader>
        <div className='flex items-start justify-between gap-4'>
          <div>
            <CardTitle>
              {t('deviceDetail.tabs.connection', { defaultValue: '连接信息' })}
            </CardTitle>
            <p className='mt-1 text-sm text-muted-foreground'>
              {t('deviceDetail.endpoints.connectionHint', {
                defaultValue: '按设备协议生成，可直接复制使用',
              })}
            </p>
          </div>
          <span className='rounded-full border bg-muted/40 px-2.5 py-1 text-xs font-medium text-muted-foreground'>
            {Object.keys(connectionData).length}{' '}
            {t('deviceDetail.endpoints.protocols', { defaultValue: '个协议' })}
          </span>
        </div>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className='grid gap-4 lg:grid-cols-2'>
            {[0, 1].map((item) => (
              <div
                key={item}
                className='h-48 animate-pulse rounded-xl border bg-muted/30'
              />
            ))}
          </div>
        ) : Object.keys(connectionData).length > 0 ? (
          <div className='grid gap-4 lg:grid-cols-2'>
            {connectionData.http && (
              <ProtocolCard
                title='HTTP'
                icon={<Globe className='h-4 w-4' />}
                tone='blue'
              >
                <ConnectionRow
                  label='HTTP'
                  value={connectionData.http.http}
                  onCopy={copy}
                />
                <ConnectionRow
                  label='RPC Subscribe'
                  value={connectionData.http.rpcSubscribe}
                  onCopy={copy}
                />
              </ProtocolCard>
            )}
            {connectionData.mqtt && (
              <ProtocolCard
                title='MQTT'
                icon={<Radio className='h-4 w-4' />}
                tone='violet'
              >
                <ConnectionRow
                  label='Host'
                  value={connectionData.mqtt.host}
                  onCopy={copy}
                />
                <ConnectionRow
                  label='Port'
                  value={connectionData.mqtt.port}
                  onCopy={copy}
                />
                <ConnectionRow
                  label='Telemetry Topic'
                  value={connectionData.mqtt.telemetryTopic}
                  onCopy={copy}
                />
                <ConnectionRow
                  label='Attributes Topic'
                  value={connectionData.mqtt.attributesTopic}
                  onCopy={copy}
                />
                <ConnectionRow
                  label='Attributes Subscribe Topic'
                  value={connectionData.mqtt.attributesSubscribeTopic}
                  onCopy={copy}
                />
                <ConnectionRow
                  label='RPC Subscribe Topic'
                  value={connectionData.mqtt.rpcSubscribeTopic}
                  onCopy={copy}
                />
              </ProtocolCard>
            )}
            {connectionData.coap && (
              <ProtocolCard
                title='CoAP'
                icon={<Wifi className='h-4 w-4' />}
                tone='emerald'
              >
                <ConnectionRow
                  label='CoAP'
                  value={connectionData.coap.coap}
                  onCopy={copy}
                />
                {connectionData.coap.docker && (
                  <ConnectionRow
                    label='Docker'
                    value={connectionData.coap.docker.coap}
                    onCopy={copy}
                  />
                )}
                <ConnectionRow
                  label='RPC Subscribe'
                  value={connectionData.coap.rpcSubscribe}
                  onCopy={copy}
                />
              </ProtocolCard>
            )}
          </div>
        ) : (
          <p className='text-sm text-muted-foreground'>
            {t('deviceDetail.endpoints.empty')}
          </p>
        )}
      </CardContent>
    </Card>
  )
}

function ProtocolCard({
  title,
  icon,
  tone,
  children,
}: {
  title: string
  icon: ReactNode
  tone: 'blue' | 'violet' | 'emerald'
  children: ReactNode
}) {
  const tones = {
    blue: 'border-l-blue-500 [&_svg]:text-blue-500',
    violet: 'border-l-violet-500 [&_svg]:text-violet-500',
    emerald: 'border-l-emerald-500 [&_svg]:text-emerald-500',
  }
  return (
    <div
      className={`overflow-hidden rounded-xl border border-l-4 bg-card shadow-sm ${tones[tone]}`}
    >
      <div className='flex items-center gap-2 border-b bg-muted/20 px-4 py-3'>
        {icon}
        <p className='font-semibold tracking-tight'>{title}</p>
      </div>
      <div className='divide-y'>{children}</div>
    </div>
  )
}

function ConnectionRow({
  label,
  value,
  onCopy,
}: {
  label: string
  value: string
  onCopy: (value: string, label: string) => void
}) {
  return (
    <div className='group grid grid-cols-[8rem_minmax(0,1fr)_2rem] items-start gap-3 px-4 py-3 text-sm'>
      <span className='pt-0.5 text-xs font-medium tracking-wide text-muted-foreground uppercase'>
        {label}
      </span>
      <code className='rounded bg-muted/50 px-2 py-1 font-mono text-xs leading-5 break-all'>
        {value}
      </code>
      <Button
        variant='ghost'
        size='icon'
        className='h-7 w-7'
        onClick={() => onCopy(value, label)}
        aria-label={label}
      >
        <Copy className='h-3.5 w-3.5' />
      </Button>
    </div>
  )
}
