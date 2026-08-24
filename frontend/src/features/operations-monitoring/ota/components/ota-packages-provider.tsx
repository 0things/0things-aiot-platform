import { createContext, useState, type ReactNode } from 'react'
import { type OTAPackage } from '../data/schema'

type DialogType = 'create' | 'edit' | 'delete' | 'deploy' | null

interface OTAPackagesContextType {
  openDialog: DialogType
  setOpenDialog: (dialog: DialogType) => void
  selectedPackage: OTAPackage | null
  setSelectedPackage: (pkg: OTAPackage | null) => void
  selectedRows: Record<string, boolean>
  setSelectedRows: (rows: Record<string, boolean>) => void
}

// eslint-disable-next-line react-refresh/only-export-components
export const OTAPackagesContext = createContext<
  OTAPackagesContextType | undefined
>(undefined)

export function OTAPackagesProvider({ children }: { children: ReactNode }) {
  const [openDialog, setOpenDialog] = useState<DialogType>(null)
  const [selectedPackage, setSelectedPackage] = useState<OTAPackage | null>(
    null
  )
  const [selectedRows, setSelectedRows] = useState<Record<string, boolean>>({})

  return (
    <OTAPackagesContext.Provider
      value={{
        openDialog,
        setOpenDialog,
        selectedPackage,
        setSelectedPackage,
        selectedRows,
        setSelectedRows,
      }}
    >
      {children}
    </OTAPackagesContext.Provider>
  )
}
