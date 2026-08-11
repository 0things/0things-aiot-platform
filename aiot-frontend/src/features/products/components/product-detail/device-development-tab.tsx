import { useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { deviceSDKs } from '@/config/device-sdks'
import {
  Database,
  FileCode,
  CheckCircle2,
  Package,
  ExternalLink,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useDevices } from '@/features/devices/api/queries'

type DeviceStep = 1 | 2 | 3 | 4

const steps: Array<{ id: DeviceStep; key: string; icon: typeof Database }> = [
  { id: 1, key: 'register', icon: Database },
  { id: 2, key: 'access', icon: FileCode },
  { id: 3, key: 'verify', icon: CheckCircle2 },
  { id: 4, key: 'production', icon: Package },
]

interface DeviceDevelopmentTabProps {
  productId: string
}

export function DeviceDevelopmentTab({ productId }: DeviceDevelopmentTabProps) {
  const { t } = useTranslation('deviceManagement')
  const [currentStep, setCurrentStep] = useState<DeviceStep>(1)

  const handleNextStep = () => {
    if (currentStep < 4) {
      setCurrentStep((currentStep + 1) as DeviceStep)
    }
  }

  const handlePrevStep = () => {
    if (currentStep > 1) {
      setCurrentStep((currentStep - 1) as DeviceStep)
    }
  }

  return (
    <div className='space-y-6'>
      {/* 步骤导航 */}
      <div className='flex items-center justify-between'>
        {steps.map((step, index) => (
          <div key={step.id} className='flex flex-1 items-center'>
            <button
              onClick={() => setCurrentStep(step.id)}
              className='flex items-center'
            >
              <div
                className={`flex h-10 w-10 items-center justify-center rounded-full border-2 ${
                  currentStep === step.id
                    ? 'border-primary bg-primary text-primary-foreground'
                    : currentStep > step.id
                      ? 'border-primary bg-primary/10 text-primary'
                      : 'border-muted bg-background text-muted-foreground'
                }`}
              >
                {step.id}
              </div>
              <div className='ml-3'>
                <p
                  className={`text-sm font-medium ${
                    currentStep === step.id
                      ? 'text-foreground'
                      : 'text-muted-foreground'
                  }`}
                >
                  {t(`productDetail.deviceDevelopment.steps.${step.key}`)}
                </p>
              </div>
            </button>
            {index < steps.length - 1 && (
              <div
                className={`mx-4 h-[2px] flex-1 ${
                  currentStep > step.id ? 'bg-primary' : 'bg-muted'
                }`}
              />
            )}
          </div>
        ))}
      </div>

      {/* 步骤内容 */}
      <div
        className='overflow-y-auto rounded-lg border-2 border-primary p-6'
        style={{ maxHeight: 'calc(100vh - 350px)' }}
      >
        {currentStep === 1 && (
          <RegisterStep productId={productId} onNext={handleNextStep} />
        )}
        {currentStep === 2 && (
          <AccessStep
            productId={productId}
            onNext={handleNextStep}
            onPrev={handlePrevStep}
          />
        )}
        {currentStep === 3 && (
          <VerifyStep onNext={handleNextStep} onPrev={handlePrevStep} />
        )}
        {currentStep === 4 && (
          <ProductionStep productId={productId} onPrev={handlePrevStep} />
        )}
      </div>
    </div>
  )
}

// 步骤1: 设备注册
function RegisterStep({ onNext }: { productId: string; onNext: () => void }) {
  const { t } = useTranslation('deviceManagement')
  const navigate = useNavigate()

  const handleAddDevice = () => {
    navigate({ to: '/device-management/devices' })
  }

  return (
    <div className='space-y-6'>
      <div>
        <h3 className='text-base font-medium'>
          {t('productDetail.deviceDevelopment.register.title')}
        </h3>
        <p className='mt-2 text-sm text-muted-foreground'>
          {t('productDetail.deviceDevelopment.register.note')}
        </p>
      </div>

      <Card
        className='w-64 cursor-pointer transition-colors hover:border-primary'
        onClick={handleAddDevice}
      >
        <CardContent className='flex items-center space-x-4 p-6'>
          <div className='flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10'>
            <Database className='h-6 w-6 text-primary' />
          </div>
          <p className='font-medium'>
            {t('productDetail.deviceDevelopment.register.addDevice')}
          </p>
        </CardContent>
      </Card>

      <p className='text-sm text-muted-foreground'>
        {t('productDetail.deviceDevelopment.register.skipTip')}
      </p>

      <div className='flex justify-start'>
        <Button onClick={onNext}>
          {t('productDetail.deviceDevelopment.register.nextButton')}
        </Button>
      </div>
    </div>
  )
}

// 步骤2: 设备接入
function AccessStep({
  productId,
  onNext,
  onPrev,
}: {
  productId: string
  onNext: () => void
  onPrev: () => void
}) {
  const { t } = useTranslation('deviceManagement')
  const [selectedSDK, setSelectedSDK] = useState(deviceSDKs[0].id)
  const [selectedDevice, setSelectedDevice] = useState<string>('')

  // 获取该产品的设备列表
  const { data: devicesResponse } = useDevices({ productId })
  const devices = devicesResponse?.devices || []

  const currentSDK = deviceSDKs.find((sdk) => sdk.id === selectedSDK)

  return (
    <div className='space-y-6'>
      <div>
        <h3 className='text-base font-medium'>
          {t('productDetail.deviceDevelopment.access.title')}
        </h3>
      </div>

      {/* 开发语言/版本选择 */}
      <div>
        <div className='mb-2 flex items-center'>
          <label className='text-sm font-medium'>
            {t('productDetail.deviceDevelopment.access.sdkSelection')}
          </label>
        </div>
        <p className='mb-4 text-sm text-muted-foreground'>
          {t('productDetail.deviceDevelopment.access.sdkSelectionNote')}
        </p>
        <div className='grid grid-cols-4 gap-3'>
          {deviceSDKs.map((sdk) => (
            <button
              key={sdk.id}
              onClick={() => setSelectedSDK(sdk.id)}
              className={`rounded-lg border-2 px-4 py-3 text-sm font-medium transition-colors ${
                selectedSDK === sdk.id
                  ? 'border-primary bg-primary/5 text-primary'
                  : 'border-muted bg-background text-foreground hover:border-primary/50'
              }`}
            >
              {sdk.name}
            </button>
          ))}
        </div>
      </div>

      {/* 设备开发与配置信息 */}
      <div>
        <h4 className='mb-4 text-sm font-medium'>
          {t('productDetail.deviceDevelopment.access.configTitle')}
        </h4>

        <div className='space-y-4'>
          {/* SDK 文档和下载 */}
          <div>
            <p className='mb-3 text-sm text-muted-foreground'>
              {t('productDetail.deviceDevelopment.access.configNote')}
            </p>
            <div className='flex gap-3'>
              <a
                href={currentSDK?.docUrl}
                target='_blank'
                rel='noopener noreferrer'
                className='inline-flex items-center gap-2 rounded-lg border bg-background px-6 py-4 text-sm font-medium transition-colors hover:bg-accent'
              >
                <FileCode className='h-5 w-5' />
                <span>
                  {t('productDetail.deviceDevelopment.access.devDoc', {
                    sdk: currentSDK?.name,
                  })}
                </span>
                <ExternalLink className='h-4 w-4 opacity-50' />
              </a>
              <a
                href={currentSDK?.packageUrl}
                target='_blank'
                rel='noopener noreferrer'
                className='inline-flex items-center gap-2 rounded-lg border bg-background px-6 py-4 text-sm font-medium transition-colors hover:bg-accent'
              >
                <Package className='h-5 w-5' />
                <span>
                  {t('productDetail.deviceDevelopment.access.downloadPackage', {
                    sdk: currentSDK?.name,
                  })}
                </span>
                <ExternalLink className='h-4 w-4 opacity-50' />
              </a>
            </div>
          </div>

          {/* 设备配置选择 */}
          <div>
            <p className='mb-3 text-sm text-muted-foreground'>
              {t('productDetail.deviceDevelopment.access.deviceConfigNote')}
            </p>
            <Select value={selectedDevice} onValueChange={setSelectedDevice}>
              <SelectTrigger className='w-80'>
                <SelectValue
                  placeholder={t(
                    'productDetail.deviceDevelopment.access.selectDevice'
                  )}
                />
              </SelectTrigger>
              <SelectContent>
                {devices.map((device) => (
                  <SelectItem key={device.id} value={device.id || ''}>
                    {device.deviceKey} - {device.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            {selectedDevice && (
              <p className='mt-2 text-xs text-muted-foreground'>
                {t('productDetail.deviceDevelopment.access.deviceSelected')}
              </p>
            )}
          </div>
        </div>
      </div>

      <div className='flex justify-between'>
        <Button variant='outline' onClick={onPrev}>
          {t('productDetail.deviceDevelopment.prevButton')}
        </Button>
        <Button onClick={onNext}>
          {t('productDetail.deviceDevelopment.access.nextButton')}
        </Button>
      </div>
    </div>
  )
}

// 步骤3: 连接验证
function VerifyStep({
  onNext,
  onPrev,
}: {
  onNext: () => void
  onPrev: () => void
}) {
  const { t } = useTranslation('deviceManagement')

  return (
    <div className='space-y-6'>
      <div>
        <h3 className='text-base font-medium'>
          {t('productDetail.deviceDevelopment.verify.title')}
        </h3>
        <p className='mt-2 text-sm text-muted-foreground'>
          {t('productDetail.deviceDevelopment.verify.description')}
        </p>
      </div>

      <Card>
        <CardContent className='p-6'>
          <div className='space-y-4'>
            <div className='flex items-start space-x-3'>
              <div className='flex h-6 w-6 items-center justify-center rounded-full bg-primary/10 text-sm font-medium text-primary'>
                1
              </div>
              <div className='flex-1'>
                <p className='text-sm'>
                  {t('productDetail.deviceDevelopment.verify.step1')}
                </p>
              </div>
            </div>
            <div className='flex items-start space-x-3'>
              <div className='flex h-6 w-6 items-center justify-center rounded-full bg-primary/10 text-sm font-medium text-primary'>
                2
              </div>
              <div className='flex-1'>
                <p className='text-sm'>
                  {t('productDetail.deviceDevelopment.verify.step2')}
                </p>
              </div>
            </div>
            <div className='flex items-start space-x-3'>
              <div className='flex h-6 w-6 items-center justify-center rounded-full bg-primary/10 text-sm font-medium text-primary'>
                3
              </div>
              <div className='flex-1'>
                <p className='text-sm'>
                  {t('productDetail.deviceDevelopment.verify.step3')}
                </p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <div className='flex justify-between'>
        <Button variant='outline' onClick={onPrev}>
          {t('productDetail.deviceDevelopment.prevButton')}
        </Button>
        <Button onClick={onNext}>
          {t('productDetail.deviceDevelopment.verify.nextButton')}
        </Button>
      </div>
    </div>
  )
}

// 步骤4: 设备量产
function ProductionStep({ onPrev }: { productId: string; onPrev: () => void }) {
  const { t } = useTranslation('deviceManagement')
  const navigate = useNavigate()

  const handleBatchRegistration = () => {
    navigate({ to: '/device-management/devices' })
  }

  return (
    <div className='space-y-6'>
      <div>
        <h3 className='text-base font-medium'>
          {t('productDetail.deviceDevelopment.production.title')}
        </h3>
        <p className='mt-2 text-sm text-muted-foreground'>
          {t('productDetail.deviceDevelopment.production.description')}
        </p>
      </div>

      <div className='space-y-4'>
        <Card
          className='cursor-pointer transition-colors hover:border-primary'
          onClick={handleBatchRegistration}
        >
          <CardContent className='p-6'>
            <h4 className='mb-2 font-medium'>
              {t('productDetail.deviceDevelopment.production.batchTitle')}
            </h4>
            <p className='text-sm text-muted-foreground'>
              {t('productDetail.deviceDevelopment.production.batchDescription')}
            </p>
          </CardContent>
        </Card>
      </div>

      <div className='flex justify-start'>
        <Button variant='outline' onClick={onPrev}>
          {t('productDetail.deviceDevelopment.prevButton')}
        </Button>
      </div>
    </div>
  )
}
