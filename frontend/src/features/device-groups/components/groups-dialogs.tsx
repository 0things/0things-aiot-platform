import { GroupsActionDialog } from './groups-action-dialog'
import { GroupsDeleteDialog } from './groups-delete-dialog'
import { useGroups } from './groups-provider'

export function GroupsDialogs() {
  const { open, setOpen, currentRow, setCurrentRow } = useGroups()

  return (
    <>
      <GroupsActionDialog
        key='group-create'
        open={open === 'create'}
        onOpenChange={() => setOpen('create')}
      />

      {currentRow && (
        <>
          <GroupsActionDialog
            key={`group-edit-${currentRow.groupUuid}`}
            open={open === 'edit'}
            onOpenChange={() => {
              setOpen('edit')
              setTimeout(() => {
                setCurrentRow(null)
              }, 500)
            }}
            currentRow={currentRow}
          />

          <GroupsDeleteDialog
            key={`group-delete-${currentRow.groupUuid}`}
            open={open === 'delete'}
            onOpenChange={() => {
              setOpen('delete')
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
