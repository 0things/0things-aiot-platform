import React, { useState } from 'react'
import useDialogState from '@/hooks/use-dialog-state'
import { type User, type InvitationItem } from '../data/schema'

type UsersDialogType = 'invite' | 'add' | 'edit' | 'delete' | 'revoke'

type UsersContextType = {
  open: UsersDialogType | null
  setOpen: (str: UsersDialogType | null) => void
  currentRow: User | null
  setCurrentRow: React.Dispatch<React.SetStateAction<User | null>>
  currentInvitation: InvitationItem | null
  setCurrentInvitation: React.Dispatch<
    React.SetStateAction<InvitationItem | null>
  >
}

const UsersContext = React.createContext<UsersContextType | null>(null)

export function UsersProvider({ children }: { children: React.ReactNode }) {
  const [open, setOpen] = useDialogState<UsersDialogType>(null)
  const [currentRow, setCurrentRow] = useState<User | null>(null)
  const [currentInvitation, setCurrentInvitation] =
    useState<InvitationItem | null>(null)

  return (
    <UsersContext
      value={{
        open,
        setOpen,
        currentRow,
        setCurrentRow,
        currentInvitation,
        setCurrentInvitation,
      }}
    >
      {children}
    </UsersContext>
  )
}

// eslint-disable-next-line react-refresh/only-export-components
export const useUsers = () => {
  const usersContext = React.useContext(UsersContext)

  if (!usersContext) {
    throw new Error('useUsers has to be used within <UsersContext>')
  }

  return usersContext
}
