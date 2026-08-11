import { useEffect, useMemo, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { deviceServiceApi } from '@/api/clients'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { ScrollArea } from '@/components/ui/scroll-area'
import { useUpdateOTAPackage } from '../api/queries'
import { useOTAPackagesContext } from '../hooks/use-ota-packages-context'

export function DeployPackageDialog() {
  const { t } = useTranslation('operationsMonitoring')
  const { openDialog, setOpenDialog, selectedPackage } = useOTAPackagesContext()
  const updatePackage = useUpdateOTAPackage()
  const [targetDevices, setTargetDevices] = useState<'all' | 'specific'>('all')
  const [selectedDeviceIds, setSelectedDeviceIds] = useState<number[]>([])
  const [isDeploying, setIsDeploying] = useState(false)

  const productId = selectedPackage?.productId
    ? String(selectedPackage.productId)
    : undefined

  const allDevicesQuery = useQuery({
    queryKey: ['ota-deploy-all-devices', productId],
    queryFn: async () => {
      const response = await deviceServiceApi.deviceServiceListDevices({
        productId,
        page: 1,
        pageSize: 1,
        enabled: true,
      })
      return response.data.total || 0
    },
    enabled: !!productId && openDialog === 'deploy',
  })

  const specificDevicesQuery = useQuery({
    queryKey: ['ota-deploy-specific-devices', productId],
    queryFn: async () => {
      const response = await deviceServiceApi.deviceServiceListDevices({
        productId,
        page: 1,
        pageSize: 200,
        enabled: true,
      })
      return response.data.devices || []
    },
    enabled: !!productId && openDialog === 'deploy' && targetDevices === 'specific',
  })

  useEffect(() => {
    if (openDialog === 'deploy') {
      setTargetDevices('all')
      setSelectedDeviceIds([])
    }
  }, [openDialog, selectedPackage?.id])

  const estimatedDeviceCount = useMemo(() => {
    if (targetDevices === 'all') {
      return allDevicesQuery.data ?? 0
    }
    return selectedDeviceIds.length
  }, [allDevicesQuery.data, selectedDeviceIds.length, targetDevices])

  const toggleDeviceSelection = (deviceId: number, checked: boolean) => {
    setSelectedDeviceIds((current) => {
      if (checked) {
        return current.includes(deviceId) ? current : [...current, deviceId]
      }
      return current.filter((id) => id !== deviceId)
    })
  }

  const handleDeploy = async () => {
    if (!selectedPackage?.id) return
    if (targetDevices === 'specific' && selectedDeviceIds.length === 0) {
      toast.error('Select at least one device')
      return
    }

    setIsDeploying(true)

    try {
      await updatePackage.mutateAsync({
        id: selectedPackage.id,
        data: {
          status: 'deploying',
          metadata: JSON.stringify({
            deploy_target_type: targetDevices,
            deploy_device_ids: selectedDeviceIds,
          }),
        },
      })

      toast.success(t('ota.notifications.deploymentStarted'))
      setOpenDialog(null)
    } catch (error) {
      console.error('Failed to deploy package:', error)
      toast.error(t('ota.notifications.error'))
    } finally {
      setIsDeploying(false)
    }
  }

  return (
    <Dialog
      open={openDialog === 'deploy'}
      onOpenChange={(open) => !open && setOpenDialog(null)}
    >
      <DialogContent className='sm:max-w-[560px]'>
        <DialogHeader>
          <DialogTitle>{t('ota.deployForm.title')}</DialogTitle>
          <DialogDescription>
            {t('ota.deployForm.description')}
          </DialogDescription>
        </DialogHeader>

        <div className='space-y-6 py-4'>
          <div className='space-y-3'>
            <Label className='text-base font-semibold'>
              {t('ota.deployForm.fields.targetDevices')}
            </Label>
            <RadioGroup
              value={targetDevices}
              onValueChange={(value) =>
                setTargetDevices(value as 'all' | 'specific')
              }
            >
              <div className='flex items-center space-x-2'>
                <RadioGroupItem value='all' id='all-devices' />
                <Label
                  htmlFor='all-devices'
                  className='cursor-pointer font-normal'
                >
                  {t('ota.deployForm.fields.allDevices')}
                </Label>
              </div>
              <div className='flex items-center space-x-2'>
                <RadioGroupItem value='specific' id='specific-devices' />
                <Label
                  htmlFor='specific-devices'
                  className='cursor-pointer font-normal'
                >
                  {t('ota.deployForm.fields.specificDevices')}
                </Label>
              </div>
            </RadioGroup>
          </div>

          {targetDevices === 'specific' && (
            <div className='space-y-3'>
              <Label className='text-base font-semibold'>Devices</Label>
              <div className='rounded-lg border'>
                <ScrollArea className='h-56'>
                  <div className='space-y-2 p-3'>
                    {specificDevicesQuery.isLoading && (
                      <div className='text-sm text-muted-foreground'>
                        Loading devices...
                      </div>
                    )}
                    {!specificDevicesQuery.isLoading &&
                      specificDevicesQuery.data?.length === 0 && (
                        <div className='text-sm text-muted-foreground'>
                          No enabled devices available for this product.
                        </div>
                      )}
                    {specificDevicesQuery.data?.map((device) => {
                      const deviceId = Number(device.id)
                      const checked = selectedDeviceIds.includes(deviceId)
                      return (
                        <label
                          key={device.id}
                          className='flex cursor-pointer items-start gap-3 rounded-md border p-3'
                        >
                          <Checkbox
                            checked={checked}
                            onCheckedChange={(value) =>
                              toggleDeviceSelection(deviceId, Boolean(value))
                            }
                          />
                          <div className='min-w-0 text-sm'>
                            <div className='font-medium'>
                              {device.name || device.deviceKey}
                            </div>
                            <div className='text-muted-foreground font-mono text-xs'>
                              ID {device.id} · {device.deviceKey}
                            </div>
                          </div>
                        </label>
                      )
                    })}
                  </div>
                </ScrollArea>
              </div>
            </div>
          )}

          {selectedPackage && (
            <div className='space-y-2 rounded-lg border p-4'>
              <div className='text-sm font-medium'>Package Details</div>
              <div className='space-y-1 text-sm text-muted-foreground'>
                <div>
                  <span className='font-medium'>Name:</span>{' '}
                  {selectedPackage.packageName}
                </div>
                <div>
                  <span className='font-medium'>Version:</span>{' '}
                  {selectedPackage.version}
                </div>
                <div>
                  <span className='font-medium'>Target Product:</span>{' '}
                  {selectedPackage.productName}
                </div>
                <div>
                  <span className='font-medium'>
                    {t('ota.deployForm.fields.estimatedDevices')}:
                  </span>{' '}
                  {estimatedDeviceCount} devices
                </div>
              </div>
            </div>
          )}

          <div className='rounded-lg bg-muted p-4'>
            <p className='text-sm'>
              {t('ota.deployForm.confirmation.message', {
                count: estimatedDeviceCount,
              })}
            </p>
            <p className='mt-2 text-sm font-medium text-destructive'>
              {t('ota.deployForm.confirmation.warning')}
            </p>
          </div>
        </div>

        <DialogFooter>
          <Button
            type='button'
            variant='outline'
            onClick={() => setOpenDialog(null)}
            disabled={isDeploying}
          >
            {t('ota.deployForm.cancel')}
          </Button>
          <Button
            onClick={handleDeploy}
            disabled={
              isDeploying ||
              (targetDevices === 'specific' && selectedDeviceIds.length === 0)
            }
          >
            {isDeploying
              ? t('common:deploying', { defaultValue: 'Deploying...' })
              : t('ota.deployForm.submit')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
