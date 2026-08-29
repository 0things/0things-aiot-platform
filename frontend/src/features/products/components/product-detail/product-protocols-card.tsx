import { useTranslation } from 'react-i18next'
import type { AiotBackendApiProductV1ProductProtocolInput } from '@/api/generated/model'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { thingModelTopics, replaceProductKey } from '../../data/topics'

type ProductProtocolsCardProps = {
  protocols: AiotBackendApiProductV1ProductProtocolInput[]
  productKey: string
}

// 产品协议由产品新增、编辑请求统一保存，详情页仅展示当前配置。
export function ProductProtocolsCard({
  protocols,
  productKey,
}: ProductProtocolsCardProps) {
  const { t } = useTranslation('deviceManagement')

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t('productDetail.protocols.title')}</CardTitle>
      </CardHeader>
      <CardContent>
        {protocols.length === 0 ? (
          <p className='text-sm text-muted-foreground'>
            {t('productDetail.protocols.empty')}
          </p>
        ) : (
          <div className='space-y-4'>
            {protocols.map((item) => {
              const transport = item.transportProtocol.toLowerCase()
              const topics =
                transport === 'mqtt'
                  ? thingModelTopics.map((topic) =>
                      replaceProductKey(topic.topic, productKey)
                    )
                  : []
              return (
                <div
                  key={`${item.transportProtocol}-${item.applicationProtocol}`}
                  className='rounded-xl border bg-muted/20 p-4'
                >
                  <div className='flex items-center gap-2'>
                    <span className='font-semibold'>
                      {item.transportProtocol.toUpperCase()}
                    </span>
                    <span className='rounded-full bg-background px-2 py-0.5 text-xs text-muted-foreground'>
                      {item.applicationProtocol.toUpperCase()}
                    </span>
                  </div>
                  <p className='mt-3 text-xs font-medium tracking-wide text-muted-foreground uppercase'>
                    {t('productDetail.protocols.template', {
                      defaultValue: '接入模板',
                    })}
                  </p>
                  <div className='mt-2 space-y-1'>
                    {topics.length > 0 ? (
                      topics.slice(0, 3).map((topic) => (
                        <code
                          key={topic}
                          className='block truncate rounded bg-background px-2 py-1 text-xs text-muted-foreground'
                        >
                          {topic}
                        </code>
                      ))
                    ) : (
                      <code className='block rounded bg-background px-2 py-1 text-xs text-muted-foreground'>
                        {transport}://{`{deviceKey}`}/telemetry
                      </code>
                    )}
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
