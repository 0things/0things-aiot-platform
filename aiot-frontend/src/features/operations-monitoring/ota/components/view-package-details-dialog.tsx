import { useTranslation } from 'react-i18next'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import { useOTAPackagesContext } from '../hooks/use-ota-packages-context'

export function ViewPackageDetailsDialog() {
  const { t } = useTranslation('operationsMonitoring')
  const { openDialog, setOpenDialog, selectedPackage } = useOTAPackagesContext()

  if (!selectedPackage) return null

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toISOString().slice(0, 19).replace('T', ' ') + ' UTC'
  }

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
  }

  return (
    <Sheet
      open={openDialog === 'view'}
      onOpenChange={(open) => !open && setOpenDialog(null)}
    >
      <SheetContent className='flex w-full flex-col gap-0 p-0 sm:max-w-[480px] lg:max-w-[560px]'>
        <SheetHeader className='px-6 pt-6'>
          <SheetTitle>{t('ota.viewDialog.title')}</SheetTitle>
          <SheetDescription>{t('ota.viewDialog.description')}</SheetDescription>
        </SheetHeader>
        <div className='flex-1 overflow-y-auto px-6 pb-6'>
          <div className='space-y-6 pt-4'>
            {/* Basic Information */}
            <div className='space-y-3'>
              <h3 className='text-sm font-semibold'>
                {t('ota.viewDialog.sections.basicInfo')}
              </h3>
              <div className='grid gap-3 rounded-lg border p-4'>
                <div className='grid grid-cols-3 gap-2'>
                  <span className='text-sm text-muted-foreground'>
                    {t('ota.packageForm.fields.packageName')}
                  </span>
                  <span className='col-span-2 font-mono text-sm font-medium'>
                    {selectedPackage.packageName}
                  </span>
                </div>
                <div className='grid grid-cols-3 gap-2'>
                  <span className='text-sm text-muted-foreground'>
                    {t('ota.packageForm.fields.version')}
                  </span>
                  <span className='col-span-2 font-mono text-sm font-medium'>
                    {selectedPackage.version}
                  </span>
                </div>
                <div className='grid grid-cols-3 gap-2'>
                  <span className='text-sm text-muted-foreground'>
                    {t('ota.packageForm.fields.packageType')}
                  </span>
                  <div className='col-span-2'>
                    <Badge variant='secondary'>
                      {t(
                        `ota.packageList.packageTypes.${selectedPackage.packageType}`
                      )}
                    </Badge>
                  </div>
                </div>
                <div className='grid grid-cols-3 gap-2'>
                  <span className='text-sm text-muted-foreground'>
                    {t('ota.packageForm.fields.productName')}
                  </span>
                  <span className='col-span-2 text-sm font-medium'>
                    {selectedPackage.productName}
                  </span>
                </div>
                <div className='grid grid-cols-3 gap-2'>
                  <span className='text-sm text-muted-foreground'>
                    {t('ota.packageList.columns.status')}
                  </span>
                  <div className='col-span-2'>
                    <Badge
                      variant={
                        selectedPackage.status === 'completed'
                          ? 'default'
                          : selectedPackage.status === 'failed'
                            ? 'destructive'
                            : 'outline'
                      }
                    >
                      {t(`ota.packageList.statuses.${selectedPackage.status}`)}
                    </Badge>
                  </div>
                </div>
                {selectedPackage.description && (
                  <div className='grid grid-cols-3 gap-2'>
                    <span className='text-sm text-muted-foreground'>
                      {t('ota.packageForm.fields.description')}
                    </span>
                    <span className='col-span-2 text-sm'>
                      {selectedPackage.description}
                    </span>
                  </div>
                )}
              </div>
            </div>

            {/* File Information */}
            <div className='space-y-3'>
              <h3 className='text-sm font-semibold'>
                {t('ota.viewDialog.sections.fileInfo')}
              </h3>
              <div className='grid gap-3 rounded-lg border p-4'>
                <div className='grid grid-cols-3 gap-2'>
                  <span className='text-sm text-muted-foreground'>
                    {t('ota.viewDialog.fields.fileSize')}
                  </span>
                  <span className='col-span-2 text-sm font-medium'>
                    {formatFileSize(selectedPackage.fileSize)}
                  </span>
                </div>
                <div className='grid grid-cols-3 gap-2'>
                  <span className='text-sm text-muted-foreground'>
                    {t('ota.viewDialog.fields.checksum')}
                  </span>
                  <span className='col-span-2 font-mono text-sm break-all'>
                    {selectedPackage.checksum}
                  </span>
                </div>
                <div className='grid grid-cols-3 gap-2'>
                  <span className='text-sm text-muted-foreground'>
                    {t('ota.viewDialog.fields.fileUrl')}
                  </span>
                  <a
                    href={selectedPackage.fileUrl}
                    target='_blank'
                    rel='noopener noreferrer'
                    className='col-span-2 text-sm break-all text-primary hover:underline'
                  >
                    {selectedPackage.fileUrl}
                  </a>
                </div>
              </div>
            </div>

            {/* Deployment Statistics */}
            {selectedPackage.status !== 'draft' && (
              <div className='space-y-3'>
                <h3 className='text-sm font-semibold'>
                  {t('ota.viewDialog.sections.deploymentStats')}
                </h3>
                <div className='grid gap-3 rounded-lg border p-4'>
                  <div className='grid grid-cols-3 gap-2'>
                    <span className='text-sm text-muted-foreground'>
                      {t('ota.viewDialog.fields.targetDeviceCount')}
                    </span>
                    <span className='col-span-2 text-sm font-medium'>
                      {selectedPackage.targetDeviceCount}
                    </span>
                  </div>
                  <div className='grid grid-cols-3 gap-2'>
                    <span className='text-sm text-muted-foreground'>
                      {t('ota.viewDialog.fields.successCount')}
                    </span>
                    <span className='col-span-2 text-sm font-medium text-green-600'>
                      {selectedPackage.successCount}
                    </span>
                  </div>
                  <div className='grid grid-cols-3 gap-2'>
                    <span className='text-sm text-muted-foreground'>
                      {t('ota.viewDialog.fields.failureCount')}
                    </span>
                    <span className='col-span-2 text-sm font-medium text-red-600'>
                      {selectedPackage.failureCount}
                    </span>
                  </div>
                  <div className='grid grid-cols-3 gap-2'>
                    <span className='text-sm text-muted-foreground'>
                      {t('ota.viewDialog.fields.deploymentProgress')}
                    </span>
                    <span className='col-span-2 text-sm font-medium'>
                      {selectedPackage.deploymentProgress}%
                    </span>
                  </div>
                  {selectedPackage.deployedAt && (
                    <div className='grid grid-cols-3 gap-2'>
                      <span className='text-sm text-muted-foreground'>
                        {t('ota.viewDialog.fields.deployedAt')}
                      </span>
                      <span className='col-span-2 text-sm'>
                        {formatDate(selectedPackage.deployedAt)}
                      </span>
                    </div>
                  )}
                </div>
              </div>
            )}

            {/* Timestamps */}
            <div className='space-y-3'>
              <h3 className='text-sm font-semibold'>
                {t('ota.viewDialog.sections.timestamps')}
              </h3>
              <div className='grid gap-3 rounded-lg border p-4'>
                <div className='grid grid-cols-3 gap-2'>
                  <span className='text-sm text-muted-foreground'>
                    {t('ota.viewDialog.fields.createdAt')}
                  </span>
                  <span className='col-span-2 text-sm'>
                    {formatDate(selectedPackage.createdAt)}
                  </span>
                </div>
                <div className='grid grid-cols-3 gap-2'>
                  <span className='text-sm text-muted-foreground'>
                    {t('ota.viewDialog.fields.updatedAt')}
                  </span>
                  <span className='col-span-2 text-sm'>
                    {formatDate(selectedPackage.updatedAt)}
                  </span>
                </div>
                <div className='grid grid-cols-3 gap-2'>
                  <span className='text-sm text-muted-foreground'>
                    {t('ota.viewDialog.fields.createdBy')}
                  </span>
                  <span className='col-span-2 text-sm'>
                    {selectedPackage.createdBy}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <div className='sticky bottom-0 flex gap-2 bg-background pt-4 pb-2'>
            <Button
              type='button'
              variant='outline'
              onClick={() => setOpenDialog(null)}
              className='flex-1'
            >
              {t('ota.viewDialog.close')}
            </Button>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  )
}
