import { useOTAPackages as useOTAPackagesQuery } from '../api/queries'

/**
 * Hook to fetch OTA packages
 * This is a wrapper around the API query hook
 */
export function useOTAPackages() {
  return useOTAPackagesQuery()
}
