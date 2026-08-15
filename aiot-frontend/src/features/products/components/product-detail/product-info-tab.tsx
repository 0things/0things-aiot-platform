import { useTranslation } from 'react-i18next'
import type { Product as ProductV1Product } from '@/api/generated/model'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

interface ProductInfoTabProps {
  product: ProductV1Product
}

export function ProductInfoTab({ product }: ProductInfoTabProps) {
  const { t } = useTranslation('deviceManagement')

  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return '-'
    const date = new Date(dateString)
    return date.toISOString().slice(0, 19).replace('T', ' ') + ' UTC'
  }

  // Helper to get label for node type
  const getNodeTypeLabel = () => {
    const nodeType = (product as ProductV1Product & { nodeType?: string })
      .nodeType
    if (!nodeType) return '-'
    return t(`productDetail.nodeTypes.${nodeType}`, { defaultValue: '-' })
  }

  // Helper to get label for connectivity method
  const getConnectivityMethodLabel = () => {
    const connectivityMethod = (
      product as ProductV1Product & { connectivityMethod?: string }
    ).connectivityMethod
    if (!connectivityMethod) return '-'
    return t(`productDetail.connectivityMethods.${connectivityMethod}`, {
      defaultValue: '-',
    })
  }

  // Helper to get label for access protocol
  const getAccessProtocolLabel = () => {
    const accessProtocol = (
      product as ProductV1Product & { accessProtocol?: string }
    ).accessProtocol
    if (!accessProtocol) return '-'
    return t(`productDetail.accessProtocols.${accessProtocol}`, {
      defaultValue: '-',
    })
  }

  // Check if should show connectivity method (only for direct and gateway)
  const nodeType = (product as ProductV1Product & { nodeType?: string })
    .nodeType
  const showConnectivityMethod = nodeType === 'direct' || nodeType === 'gateway'

  // Dynamic label for access protocol
  const accessProtocolFieldLabel =
    nodeType === 'gateway-sub'
      ? t('productDetail.info.fields.gatewayAccessProtocol')
      : t('productDetail.info.fields.accessProtocol')

  // Get status variant and label
  const statusVariant = product.status === 'active' ? 'default' : 'secondary'
  const statusLabel = t(`productDetail.status.${product.status || 'unknown'}`)

  return (
    <div className='space-y-6'>
      <Card>
        <CardHeader className='flex flex-row items-center justify-between'>
          <CardTitle>{t('productDetail.info.title')}</CardTitle>
          <Button variant='outline' size='sm'>
            {t('productDetail.info.editButton')}
          </Button>
        </CardHeader>
        <CardContent>
          <div className='grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3'>
            <div className='space-y-1'>
              <p className='text-sm text-muted-foreground'>
                {t('productDetail.info.fields.name')}
              </p>
              <p className='font-medium'>{product.name || '-'}</p>
            </div>

            <div className='space-y-1'>
              <p className='text-sm text-muted-foreground'>
                {t('productDetail.info.fields.category')}
              </p>
              <Badge variant='secondary' className='capitalize'>
                {product.category || '-'}
              </Badge>
            </div>

            <div className='space-y-1'>
              <p className='text-sm text-muted-foreground'>
                {t('productDetail.info.fields.status')}
              </p>
              <Badge variant={statusVariant}>{statusLabel}</Badge>
            </div>

            <div className='space-y-1'>
              <p className='text-sm text-muted-foreground'>
                {t('productDetail.info.fields.nodeType')}
              </p>
              <Badge variant='outline'>{getNodeTypeLabel()}</Badge>
            </div>

            {showConnectivityMethod && (
              <div className='space-y-1'>
                <p className='text-sm text-muted-foreground'>
                  {t('productDetail.info.fields.connectivityMethod')}
                </p>
                <p className='font-medium'>{getConnectivityMethodLabel()}</p>
              </div>
            )}

            <div className='space-y-1'>
              <p className='text-sm text-muted-foreground'>
                {accessProtocolFieldLabel}
              </p>
              <p className='font-medium'>{getAccessProtocolLabel()}</p>
            </div>

            <div className='space-y-1'>
              <p className='text-sm text-muted-foreground'>
                {t('productDetail.info.fields.createdAt')}
              </p>
              <p className='font-medium'>{formatDate(product.createdAt)}</p>
            </div>

            <div className='space-y-1 md:col-span-2 lg:col-span-3'>
              <p className='text-sm text-muted-foreground'>
                {t('productDetail.info.fields.description')}
              </p>
              <p className='font-medium'>{product.description || '-'}</p>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>{t('productDetail.info.labels.title')}</CardTitle>
        </CardHeader>
        <CardContent>
          <p className='text-sm text-muted-foreground'>
            {t('productDetail.info.labels.empty')}
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
