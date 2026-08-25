import { create } from 'zustand'
import type { ApiOrganizationItem } from '@/api/generated/model'

const CURRENT_ORG_ID_KEY = 'aiot_current_org_id'

interface OrgState {
  currentOrgId: number | null
  organizations: ApiOrganizationItem[]
  setCurrentOrgId: (id: number) => void
  setOrganizations: (orgs: ApiOrganizationItem[]) => void
  reset: () => void
}

export const useOrgStore = create<OrgState>()((set) => {
  const storedId = localStorage.getItem(CURRENT_ORG_ID_KEY)
  const initialId = storedId ? Number(storedId) : null

  return {
    currentOrgId: initialId,
    organizations: [],
    setCurrentOrgId: (id) =>
      set((state) => {
        localStorage.setItem(CURRENT_ORG_ID_KEY, String(id))
        return { ...state, currentOrgId: id }
      }),
    setOrganizations: (organizations) =>
      set((state) => {
        const current = organizations.find((o) => o.is_current)
        const newOrgId = current?.id ?? (organizations[0]?.id || null)
        if (
          newOrgId &&
          (!state.currentOrgId ||
            !organizations.some((o) => o.id === state.currentOrgId))
        ) {
          localStorage.setItem(CURRENT_ORG_ID_KEY, String(newOrgId))
          return { ...state, organizations, currentOrgId: newOrgId }
        }
        return { ...state, organizations }
      }),
    reset: () =>
      set((state) => {
        localStorage.removeItem(CURRENT_ORG_ID_KEY)
        return { ...state, currentOrgId: null, organizations: [] }
      }),
  }
})
