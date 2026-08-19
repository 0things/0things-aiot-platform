import { useContext } from 'react'
import { OTAPackagesContext } from '../components/ota-packages-provider'

export function useOTAPackagesContext() {
  const context = useContext(OTAPackagesContext)
  if (context === undefined) {
    throw new Error(
      'useOTAPackagesContext must be used within OTAPackagesProvider'
    )
  }
  return context
}
