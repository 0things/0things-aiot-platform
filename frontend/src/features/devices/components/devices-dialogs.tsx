import { DevicesActionDialog } from './devices-action-dialog'
import { DevicesActivateDialog } from './devices-activate-dialog'
import { DevicesBatchUploadDialog } from './devices-batch-upload-dialog'
import { DevicesDeleteDialog } from './devices-delete-dialog'
import { DevicesEnableDialog } from './devices-enable-dialog'
import { useDevices } from './devices-provider'

export function DevicesDialogs() {
  const { open, setOpen, currentRow, setCurrentRow } = useDevices()
  return (
    <>
      <DevicesActionDialog
        key='device-create'
        open={open === 'create'}
        onOpenChange={() => setOpen('create')}
      />

      <DevicesBatchUploadDialog
        key='device-batch-upload'
        open={open === 'batch-upload'}
        onOpenChange={() => setOpen('batch-upload')}
      />

      {currentRow && (
        <>
          <DevicesActionDialog
            key={`device-edit-${currentRow.id}`}
            open={open === 'edit'}
            onOpenChange={() => {
              setOpen('edit')
              setTimeout(() => {
                setCurrentRow(null)
              }, 500)
            }}
            currentRow={currentRow}
          />

          <DevicesDeleteDialog
            key={`device-delete-${currentRow.id}`}
            open={open === 'delete'}
            onOpenChange={() => {
              setOpen('delete')
              setTimeout(() => {
                setCurrentRow(null)
              }, 500)
            }}
            currentRow={currentRow}
          />

          <DevicesActivateDialog
            key={`device-activate-${currentRow.id}`}
            open={open === 'activate'}
            onOpenChange={() => {
              setOpen('activate')
              setTimeout(() => {
                setCurrentRow(null)
              }, 500)
            }}
            currentRow={currentRow}
          />

          <DevicesEnableDialog
            key={`device-enable-${currentRow.id}`}
            open={open === 'enable'}
            onOpenChange={() => {
              setOpen('enable')
              setTimeout(() => {
                setCurrentRow(null)
              }, 500)
            }}
            currentRow={currentRow}
          />
        </>
      )}
    </>
  )
}
