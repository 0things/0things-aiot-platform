import { create } from 'zustand'

export interface OTADetailStore {
  // UI state
  activeTab: 'batches' | 'devices' | 'info'
  setActiveTab: (tab: 'batches' | 'devices' | 'info') => void

  // Device list filters
  statusFilter: string | undefined
  setStatusFilter: (status: string | undefined) => void

  // Pagination
  currentPage: number
  setCurrentPage: (page: number) => void
  pageSize: number
  setPageSize: (size: number) => void

  // Tab visibility state
  isTabActive: boolean
  setIsTabActive: (active: boolean) => void

  // Sorting
  sortBy: 'name' | 'status' | 'time'
  sortOrder: 'asc' | 'desc'
  setSortBy: (field: 'name' | 'status' | 'time', order: 'asc' | 'desc') => void

  // Last updated timestamp
  lastUpdated: Date
  setLastUpdated: (date: Date) => void

  // Reset state
  reset: () => void
}

const defaultState = {
  activeTab: 'batches' as const,
  statusFilter: undefined,
  currentPage: 1,
  pageSize: 100,
  isTabActive: true,
  sortBy: 'time' as const,
  sortOrder: 'desc' as const,
  lastUpdated: new Date(),
}

/**
 * Zustand store for managing OTA detail page state
 * Handles:
 * - Active tab selection
 * - Device list filtering and pagination
 * - Tab visibility for polling control
 * - Sorting preferences
 * - Last updated timestamp
 */
export const useOTADetailStore = create<OTADetailStore>((set) => ({
  ...defaultState,

  setActiveTab: (tab) => set({ activeTab: tab }),

  setStatusFilter: (status) => {
    set({ statusFilter: status, currentPage: 1 })
  },

  setCurrentPage: (page) => set({ currentPage: page }),

  setPageSize: (size) => set({ pageSize: size, currentPage: 1 }),

  setIsTabActive: (active) => set({ isTabActive: active }),

  setSortBy: (field, order) => set({ sortBy: field, sortOrder: order }),

  setLastUpdated: (date) => set({ lastUpdated: date }),

  reset: () => set(defaultState),
}))
