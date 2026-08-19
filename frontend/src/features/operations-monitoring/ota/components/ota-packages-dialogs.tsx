import { useTranslation } from 'react-i18next'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { useDeleteOTAPackage } from '../api/queries'
import { useOTAPackagesContext } from '../hooks/use-ota-packages-context'
import { CreatePackageDialog } from './create-package-dialog'
import { DeployPackageDialog } from './deploy-package-dialog'
import { EditPackageDialog } from './edit-package-dialog'
import { ViewPackageDetailsDialog } from './view-package-details-dialog'

export function OTAPackagesDialogs() {
  const { t } = useTranslation('operationsMonitoring')
  const { openDialog, setOpenDialog, selectedPackage } = useOTAPackagesContext()
  const deleteMutation = useDeleteOTAPackage()

  const handleDelete = async () => {
    if (!selectedPackage?.id) return

    await deleteMutation.mutateAsync(selectedPackage.id, {
      onSuccess: () => {
        setOpenDialog(null)
      },
    })
  }

  return (
    <>
      {/* Create Package Dialog */}
      <CreatePackageDialog />

      {/* Edit Package Dialog */}
      <EditPackageDialog />

      {/* View Package Details Dialog */}
      <ViewPackageDetailsDialog />

      {/* Deploy Package Dialog */}
      <DeployPackageDialog />

      {/* Delete Dialog */}
      <AlertDialog
        open={openDialog === 'delete'}
        onOpenChange={(open) => !open && setOpenDialog(null)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t('ota.deleteDialog.title')}</AlertDialogTitle>
            <AlertDialogDescription>
              {t('ota.deleteDialog.description', {
                packageName: selectedPackage?.packageName || '',
              })}
              <br />
              <span className='font-medium text-destructive'>
                {t('ota.deleteDialog.warning')}
              </span>
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={deleteMutation.isPending}>
              {t('common:cancel')}
            </AlertDialogCancel>
            <AlertDialogAction
              className='text-destructive-foreground bg-destructive hover:bg-destructive/90'
              onClick={handleDelete}
              disabled={deleteMutation.isPending}
            >
              {deleteMutation.isPending
                ? t('common:deleting', { defaultValue: 'Deleting...' })
                : t('common:delete')}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  )
}
