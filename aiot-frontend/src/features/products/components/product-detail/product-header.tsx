import { useState } from 'react'
import { ArrowLeft, Copy, Eye, EyeOff } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import type { ProductV1Product } from '@/api/generated/device-service'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'

interface ProductHeaderProps {
  product: ProductV1Product
  onBack: () => void
}

export function ProductHeader({ product, onBack }: ProductHeaderProps) {
  const { t } = useTranslation('deviceManagement')
  const [secretVisible, setSecretVisible] = useState(false)

  const handleCopySecret = () => {
    // In real implementation, get secret from secure source
    navigator.clipboard.writeText(product.productKey || '')
    toast.success(t('productDetail.header.copySuccess'))
  }

  const toggleSecretVisibility = () => {
    setSecretVisible(!secretVisible)
  }

  // Get status variant
  const statusVariant = product.status === 'active' ? 'default' : 'secondary'

  return (
    <div className='mb-6 space-y-4'>
      <div className='flex items-center gap-4'>
        <Button
          variant='ghost'
          size='icon'
          onClick={onBack}
          aria-label='Back to product list'
        >
          <ArrowLeft className='h-4 w-4' />
        </Button>
        <h1 className='text-3xl font-bold'>
          {product.name || 'Unnamed Product'}
        </h1>
        <Badge variant={statusVariant}>
          {t(`productDetail.status.${product.status || 'unknown'}`)}
        </Badge>
      </div>

      <div className='flex flex-wrap items-center gap-4 text-sm'>
        <div className='flex items-center gap-2'>
          <span className='text-muted-foreground'>
            {t('productDetail.header.productKey')}:
          </span>
          <span className='font-mono'>{product.productKey || '-'}</span>
        </div>

        <div className='flex items-center gap-2'>
          <span className='text-muted-foreground'>
            {t('productDetail.header.productSecret')}:
          </span>
          <span className='font-mono'>
            {secretVisible ? product.productKey || '-' : '********'}
          </span>
          <Button
            variant='ghost'
            size='icon'
            className='h-6 w-6'
            onClick={toggleSecretVisibility}
          >
            {secretVisible ? (
              <EyeOff className='h-3 w-3' />
            ) : (
              <Eye className='h-3 w-3' />
            )}
          </Button>
          <Button
            variant='ghost'
            size='icon'
            className='h-6 w-6'
            onClick={handleCopySecret}
          >
            <Copy className='h-3 w-3' />
          </Button>
        </div>

        <div className='flex items-center gap-2'>
          <span className='text-muted-foreground'>
            {t('productDetail.header.deviceCount')}:
          </span>
          <span className='font-semibold'>{product.deviceCount || 0}</span>
        </div>
      </div>
    </div>
  )
}
