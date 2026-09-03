import { useState } from 'react'
import { Link, useNavigate } from '@tanstack/react-router'
import { ArrowLeft, Copy, Eye, EyeOff } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import type { Device as DeviceV1Device } from '@/api/generated/model'
import { Button } from '@/components/ui/button'

interface DeviceHeaderProps {
  device: DeviceV1Device
  onBack: () => void
}

export function DeviceHeader({ device, onBack }: DeviceHeaderProps) {
  const { t } = useTranslation('deviceManagement')
  const navigate = useNavigate()
  const [secretVisible, setSecretVisible] = useState(false)

  const handleCopyProductKey = () => {
    navigator.clipboard.writeText(device.productKey || '')
    toast.success(t('deviceDetail.header.copySuccess'))
  }

  const toggleSecretVisibility = () => {
    setSecretVisible(!secretVisible)
  }

  const handleProductClick = () => {
    if (device.productKey) {
      navigate({
        to: '/device-management/products/$productKey',
        params: { productKey: device.productKey },
      })
    }
  }

  return (
    <div className='mb-6 space-y-4'>
      <div className='flex items-center gap-4'>
        <Button
          variant='ghost'
          size='icon'
          onClick={onBack}
          aria-label='Back to device list'
        >
          <ArrowLeft className='h-4 w-4' />
        </Button>
        <h1 className='text-2xl font-bold'>
          {device.name || 'Unnamed Device'}
        </h1>
      </div>

      <div className='flex flex-wrap items-center gap-4 text-sm'>
        <div className='flex items-center gap-2'>
          <span className='text-muted-foreground'>
            {t('deviceDetail.header.product')}:
          </span>
          <span>{device.productName || '-'}</span>
          {device.productKey && (
            <Button
              variant='link'
              size='sm'
              className='h-auto p-0 text-sm'
              onClick={handleProductClick}
            >
              {t('deviceDetail.header.view')}
            </Button>
          )}
        </div>

        <div className='flex items-center gap-2'>
          <span className='text-muted-foreground'>
            {t('deviceDetail.header.deviceSecret')}:
          </span>
          <span className='font-mono'>
            {secretVisible ? device.deviceKey || '-' : '********'}
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
        </div>

        <div className='flex items-center gap-2'>
          <span className='text-muted-foreground'>
            {t('deviceDetail.header.productKey')}:
          </span>
          {device.productKey ? (
            <Link
              to='/device-management/products/$productKey'
              params={{ productKey: device.productKey }}
              aria-label={t('deviceDetail.header.viewProduct')}
              className='font-mono text-primary underline-offset-4 hover:underline focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:outline-none'
            >
              {device.productKey}
            </Link>
          ) : (
            <span className='font-mono'>-</span>
          )}
          <Button
            variant='ghost'
            size='icon'
            className='h-6 w-6'
            onClick={handleCopyProductKey}
          >
            <Copy className='h-3 w-3' />
          </Button>
        </div>
      </div>
    </div>
  )
}
