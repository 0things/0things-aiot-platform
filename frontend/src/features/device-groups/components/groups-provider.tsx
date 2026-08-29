import React, { useState } from 'react'
import useDialogState from '@/hooks/use-dialog-state'
import { type DeviceGroup } from '../data/schema'

type GroupsDialogType = 'create' | 'edit' | 'delete'

type GroupsContextType = {
  open: GroupsDialogType | null
  setOpen: (str: GroupsDialogType | null) => void
  currentRow: DeviceGroup | null
  setCurrentRow: React.Dispatch<React.SetStateAction<DeviceGroup | null>>
}

const GroupsContext = React.createContext<GroupsContextType | null>(null)

export function GroupsProvider({ children }: { children: React.ReactNode }) {
  const [open, setOpen] = useDialogState<GroupsDialogType>(null)
  const [currentRow, setCurrentRow] = useState<DeviceGroup | null>(null)

  return (
    <GroupsContext value={{ open, setOpen, currentRow, setCurrentRow }}>
      {children}
    </GroupsContext>
  )
}

// eslint-disable-next-line react-refresh/only-export-components
export const useGroups = () => {
  const groupsContext = React.useContext(GroupsContext)

  if (!groupsContext) {
    throw new Error('useGroups has to be used within <GroupsProvider>')
  }

  return groupsContext
}
