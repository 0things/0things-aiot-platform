import React, { useState } from 'react'
import useDialogState from '@/hooks/use-dialog-state'
import { type Scene } from '../data/schema'

type SceneLinkageDialogType = 'create' | 'edit' | 'delete'

type SceneLinkageContextType = {
  open: SceneLinkageDialogType | null
  setOpen: (str: SceneLinkageDialogType | null) => void
  currentRow: Scene | null
  setCurrentRow: React.Dispatch<React.SetStateAction<Scene | null>>
}

const SceneLinkageContext = React.createContext<SceneLinkageContextType | null>(
  null
)

export function SceneLinkageProvider({
  children,
}: {
  children: React.ReactNode
}) {
  const [open, setOpen] = useDialogState<SceneLinkageDialogType>(null)
  const [currentRow, setCurrentRow] = useState<Scene | null>(null)

  return (
    <SceneLinkageContext
      value={{ open, setOpen, currentRow, setCurrentRow }}
    >
      {children}
    </SceneLinkageContext>
  )
}

// eslint-disable-next-line react-refresh/only-export-components
export const useSceneLinkage = () => {
  const context = React.useContext(SceneLinkageContext)

  if (!context) {
    throw new Error('useSceneLinkage has to be used within <SceneLinkageContext>')
  }

  return context
}
