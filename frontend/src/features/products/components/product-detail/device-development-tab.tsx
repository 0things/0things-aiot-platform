import { useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { deviceSDKs } from '@/config/device-sdks'
import {
  BookOpen,
  ChevronRight,
  CircleHelp,
  Database,
  FileCode,
  CheckCircle2,
  KeyRound,
  Package,
  ExternalLink,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
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
    <div className='space-y-8 pb-8'>
      <nav
        aria-label={t('productDetail.deviceDevelopment.title')}
        className='grid gap-3 border-b pb-6 sm:grid-cols-2 xl:grid-cols-4'
      >
        {steps.map((step, index) => {
          const Icon = step.icon
          const isCurrent = currentStep === step.id
          const isComplete = currentStep > step.id

          return (
            <button
              key={step.id}
              type='button'
              onClick={() => setCurrentStep(step.id)}
              className={`group flex min-w-0 items-start gap-3 rounded-lg px-3 py-2 text-left transition-colors hover:bg-muted/60 ${
                isCurrent ? 'bg-primary/5' : ''
              }`}
            >
              <span
                className={`mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-full border text-sm font-semibold ${
                  isCurrent
                    ? 'border-primary bg-primary text-primary-foreground'
                    : isComplete
                      ? 'border-primary/40 bg-primary/10 text-primary'
                      : 'border-muted-foreground/30 text-muted-foreground'
                }`}
              >
                {isComplete ? <CheckCircle2 className='h-4 w-4' /> : step.id}
              </span>
              <span className='min-w-0'>
                <span
                  className={`block text-sm font-semibold ${
                    isCurrent ? 'text-foreground' : 'text-muted-foreground'
                  }`}
                >
                  {t(`productDetail.deviceDevelopment.steps.${step.key}`)}
                </span>
                <span className='mt-1 flex items-center gap-1 text-xs text-muted-foreground'>
                  <Icon className='h-3.5 w-3.5' />
                  {t(`productDetail.deviceDevelopment.stepHints.${step.key}`)}
                </span>
              </span>
              {index < steps.length - 1 && (
                <ChevronRight className='ml-auto mt-1 hidden h-4 w-4 text-muted-foreground/50 xl:block' />
              )}
            </button>
          )
        })}
      </nav>

      <div className='grid min-w-0 gap-8 xl:grid-cols-[minmax(0,1.65fr)_minmax(260px,0.85fr)]'>
        <section className='min-w-0'>
          {currentStep === 1 && <RegisterStep onNext={handleNextStep} />}
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
          {currentStep === 4 && <ProductionStep onPrev={handlePrevStep} />}
        </section>
        <DevelopmentGuide currentStep={currentStep} />
      </div>
    </div>
  )
}

function DevelopmentGuide({ currentStep }: { currentStep: DeviceStep }) {
  const { t } = useTranslation('deviceManagement')
  const guideKey = steps.find((step) => step.id === currentStep)?.key

  return (
    <aside className='min-w-0 border-t pt-6 xl:border-t-0 xl:border-l xl:pl-8 xl:pt-0'>
      <div className='flex items-center gap-2'>
        <BookOpen className='h-4 w-4 text-primary' />
        <h3 className='font-semibold'>
          {t('productDetail.deviceDevelopment.guide.title')}
        </h3>
      </div>
      <p className='mt-2 text-sm leading-6 text-muted-foreground'>
        {t(`productDetail.deviceDevelopment.guide.${guideKey}`)}
      </p>

      <div className='mt-6 space-y-3 border-t pt-5'>
        <p className='text-xs font-medium uppercase tracking-[0.16em] text-muted-foreground'>
          {t('productDetail.deviceDevelopment.guide.checklistTitle')}
        </p>
        {['product', 'network', 'model'].map((item) => (
          <div key={item} className='flex gap-3 text-sm text-muted-foreground'>
            <CheckCircle2 className='mt-0.5 h-4 w-4 shrink-0 text-primary' />
            <span>{t(`productDetail.deviceDevelopment.guide.${item}`)}</span>
          </div>
        ))}
      </div>

      <a
        href='https://github.com/0things'
        target='_blank'
        rel='noopener noreferrer'
        className='mt-6 inline-flex items-center gap-2 text-sm font-medium text-primary hover:underline'
      >
        <CircleHelp className='h-4 w-4' />
        {t('productDetail.deviceDevelopment.guide.help')}
        <ExternalLink className='h-3.5 w-3.5' />
      </a>
    </aside>
  )
}

// 步骤1: 设备注册
function RegisterStep({ onNext }: { onNext: () => void }) {
  const { t } = useTranslation('deviceManagement')
  const navigate = useNavigate()

  const handleAddDevice = () => {
    navigate({ to: '/device-management/devices' })
  }

  return (
    <div className='space-y-8'>
      <div>
        <p className='text-sm font-medium text-primary'>
          {t('productDetail.deviceDevelopment.currentStep', { step: 1 })}
        </p>
        <h3 className='mt-2 text-2xl font-semibold tracking-tight'>
          {t('productDetail.deviceDevelopment.register.title')}
        </h3>
        <p className='mt-3 max-w-2xl text-sm leading-6 text-muted-foreground'>
          {t('productDetail.deviceDevelopment.register.note')}
        </p>
      </div>

      <button
        type='button'
        className='group flex w-full max-w-2xl items-center justify-between rounded-xl border bg-card p-5 text-left shadow-sm transition-all hover:-translate-y-0.5 hover:border-primary/60 hover:shadow-md'
        onClick={handleAddDevice}
      >
        <span className='flex items-center gap-4'>
          <span className='flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10'>
            <Database className='h-6 w-6 text-primary' />
          </span>
          <span>
            <span className='block font-semibold'>
            {t('productDetail.deviceDevelopment.register.addDevice')}
            </span>
            <span className='mt-1 block text-sm text-muted-foreground'>
              {t('productDetail.deviceDevelopment.register.addDeviceHint')}
            </span>
          </span>
        </span>
        <ChevronRight className='h-5 w-5 text-muted-foreground transition-transform group-hover:translate-x-1 group-hover:text-primary' />
      </button>

      <div className='max-w-2xl border-t pt-6'>
        <div className='flex items-center gap-2'>
          <KeyRound className='h-4 w-4 text-primary' />
          <h4 className='font-semibold'>
            {t('productDetail.deviceDevelopment.register.resultTitle')}
          </h4>
        </div>
        <div className='mt-4 grid gap-3 sm:grid-cols-2'>
          {['deviceName', 'deviceSecret'].map((item) => (
            <div key={item} className='rounded-lg bg-muted/60 p-4'>
              <p className='text-sm font-medium'>
                {t(`productDetail.deviceDevelopment.register.${item}`)}
              </p>
              <p className='mt-1 text-xs leading-5 text-muted-foreground'>
                {t(`productDetail.deviceDevelopment.register.${item}Hint`)}
              </p>
            </div>
          ))}
        </div>
      </div>

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
    <div className='space-y-8'>
      <div>
        <p className='text-sm font-medium text-primary'>
          {t('productDetail.deviceDevelopment.currentStep', { step: 2 })}
        </p>
        <h3 className='mt-2 text-2xl font-semibold tracking-tight'>
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
        <div className='grid grid-cols-2 gap-3 sm:grid-cols-4'>
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
            <div className='flex flex-wrap gap-3'>
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
              <SelectTrigger className='w-full max-w-80'>
                <SelectValue
                  placeholder={t(
                    'productDetail.deviceDevelopment.access.selectDevice'
                  )}
                />
              </SelectTrigger>
              <SelectContent>
                {devices.map((device) => (
                  <SelectItem
                    key={device.id}
                    value={device.id != null ? String(device.id) : ''}
                  >
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
    <div className='space-y-8'>
      <div>
        <p className='text-sm font-medium text-primary'>
          {t('productDetail.deviceDevelopment.currentStep', { step: 3 })}
        </p>
        <h3 className='mt-2 text-2xl font-semibold tracking-tight'>
          {t('productDetail.deviceDevelopment.verify.title')}
        </h3>
        <p className='mt-2 text-sm text-muted-foreground'>
          {t('productDetail.deviceDevelopment.verify.description')}
        </p>
      </div>

      <div className='max-w-2xl rounded-xl border bg-card p-6'>
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
      </div>

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
function ProductionStep({ onPrev }: { onPrev: () => void }) {
  const { t } = useTranslation('deviceManagement')
  const navigate = useNavigate()

  const handleBatchRegistration = () => {
    navigate({ to: '/device-management/devices' })
  }

  return (
    <div className='space-y-8'>
      <div>
        <p className='text-sm font-medium text-primary'>
          {t('productDetail.deviceDevelopment.currentStep', { step: 4 })}
        </p>
        <h3 className='mt-2 text-2xl font-semibold tracking-tight'>
          {t('productDetail.deviceDevelopment.production.title')}
        </h3>
        <p className='mt-2 text-sm text-muted-foreground'>
          {t('productDetail.deviceDevelopment.production.description')}
        </p>
      </div>

      <div className='space-y-4'>
        <button
          type='button'
          className='w-full max-w-2xl rounded-xl border bg-card p-6 text-left transition-colors hover:border-primary/60'
          onClick={handleBatchRegistration}
        >
            <h4 className='mb-2 font-semibold'>
              {t('productDetail.deviceDevelopment.production.batchTitle')}
            </h4>
            <p className='text-sm text-muted-foreground'>
              {t('productDetail.deviceDevelopment.production.batchDescription')}
            </p>
        </button>
      </div>

      <div className='flex justify-start'>
        <Button variant='outline' onClick={onPrev}>
          {t('productDetail.deviceDevelopment.prevButton')}
        </Button>
      </div>
    </div>
  )
}
