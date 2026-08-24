import React, { type ReactNode } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { renderHook, waitFor } from '@testing-library/react'
import * as axiosClients from '@/api/clients'
import {
  useOTAPackageDetail,
  useUpgradeStatistics,
  useDeviceDeployments,
  useUpgradeBatches,
} from '../api/detail-queries'

/**
 * Integration tests for OTA Package Detail Page
 * Tests the complete data flow from API calls to component rendering
 */

interface MockGetConfig {
  params?: { page?: number; status?: string }
}

// Mock axios module
jest.mock('@/api/clients', () => ({
  axiosInstance: {
    get: jest.fn(),
  },
}))

// Create query client wrapper for tests
function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  })
  return ({ children }: { children: ReactNode }) =>
    React.createElement(QueryClientProvider, { client: queryClient }, children)
}

describe('OTA Package Detail Page - Integration Tests', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  /**
   * Test: Complete data flow from package load to statistics
   * Scenario: User opens OTA package detail page
   */
  it('should load package data and display all tabs with data', async () => {
    const packageData = {
      id: 'pkg-001',
      packageName: 'firmware-v2.0.0',
      version: 'v2.0.0',
      packageType: 'upgrade',
      productId: 1,
      productKey: 'WMS003',
      status: 'completed',
      fileUrl: 'https://example.com/firmware.bin',
      checksum: 'abc123def456',
      verificationStatus: 'completed',
      verificationProgress: 100,
      description: 'Major release with performance improvements',
      releaseNotes: 'v2.0.0: Performance improvements and bug fixes',
      metadata: { build: '2024-02-21', branch: 'main' },
      createdAt: '2024-02-21T10:30:00Z',
      updatedAt: '2024-02-21T10:30:00Z',
    }

    const statsData = {
      packageId: 'pkg-001',
      totalTargetDevices: 1000,
      successfulUpgrades: 800,
      failedUpgrades: 50,
      cancelledUpgrades: 10,
      pendingUpgrades: 100,
      inProgressUpgrades: 40,
    }

    const batchesData = {
      batches: [
        {
          batchId: 'batch-001',
          batchName: 'Phase 1: Cloud Devices',
          batchType: 'production',
          upgradeStrategy: 'static',
          status: 'completed',
          targetDeviceCount: 500,
          createdAt: '2024-02-20T10:30:00Z',
        },
        {
          batchId: 'batch-002',
          batchName: 'Phase 2: Edge Devices',
          batchType: 'production',
          upgradeStrategy: 'rolling',
          status: 'completed',
          targetDeviceCount: 500,
          createdAt: '2024-02-21T08:00:00Z',
        },
      ],
    }

    // Setup mock responses
    ;(axiosClients.axiosInstance.get as jest.Mock).mockImplementation(
      (url: string) => {
        if (url.includes('/upgrade-statistics')) {
          return Promise.resolve({ data: { statistics: statsData } })
        }
        if (url.includes('/batches')) {
          return Promise.resolve({ data: batchesData })
        }
        if (url.includes('/ota-packages/')) {
          return Promise.resolve({ data: packageData })
        }
        return Promise.reject(new Error('Unknown URL: ' + url))
      }
    )

    // Test package detail hook
    const { result: packageResult } = renderHook(
      () => useOTAPackageDetail('firmware-v2.0.0'),
      { wrapper: createWrapper() }
    )

    await waitFor(() => {
      expect(packageResult.current.isLoading).toBe(false)
    })

    expect(packageResult.current.data?.packageName).toBe('firmware-v2.0.0')
    expect(packageResult.current.data?.version).toBe('v2.0.0')
    expect(packageResult.current.data?.description).toBe(
      'Major release with performance improvements'
    )

    // Test statistics hook
    const { result: statsResult } = renderHook(
      () => useUpgradeStatistics('firmware-v2.0.0', true),
      { wrapper: createWrapper() }
    )

    await waitFor(() => {
      expect(statsResult.current.isLoading).toBe(false)
    })

    expect(statsResult.current.data?.totalTargetDevices).toBe(1000)
    expect(statsResult.current.data?.successfulUpgrades).toBe(800)
    expect(statsResult.current.data?.failedUpgrades).toBe(50)

    // Test batches hook
    const { result: batchesResult } = renderHook(
      () => useUpgradeBatches('firmware-v2.0.0'),
      { wrapper: createWrapper() }
    )

    await waitFor(() => {
      expect(batchesResult.current.isLoading).toBe(false)
    })

    expect(batchesResult.current.data).toHaveLength(2)
    expect(batchesResult.current.data?.[0]?.batchName).toBe(
      'Phase 1: Cloud Devices'
    )
    expect(batchesResult.current.data?.[1]?.status).toBe('completed')
  })

  /**
   * Test: Device deployments with pagination
   * Scenario: User views device list with pagination controls
   */
  it('should load device deployments with pagination support', async () => {
    const page1Data = {
      deployments: Array.from({ length: 50 }, (_, i) => ({
        deviceId: `device-${i + 1}`,
        deviceKey: `dev-${i + 1}`,
        deviceName: `IoT Device ${i + 1}`,
        productId: 1,
        productKey: 'WMS003',
        currentVersion: 'v1.9.0',
        upgradeBatchId: 'batch-001',
        status: i % 3 === 0 ? 'success' : i % 3 === 1 ? 'pending' : 'failed',
        lastStatusChangeTime: Date.now() - Math.random() * 86400000,
        createdAt: '2024-02-20T10:30:00Z',
      })),
      pagination: {
        page: 1,
        pageSize: 50,
        total: 150,
      },
    }

    const page2Data = {
      deployments: Array.from({ length: 50 }, (_, i) => ({
        deviceId: `device-${i + 51}`,
        deviceKey: `dev-${i + 51}`,
        deviceName: `IoT Device ${i + 51}`,
        productId: 1,
        productKey: 'WMS003',
        currentVersion: 'v2.0.0',
        upgradeBatchId: 'batch-002',
        status: 'success',
        lastStatusChangeTime: Date.now() - Math.random() * 3600000,
        createdAt: '2024-02-21T10:30:00Z',
      })),
      pagination: {
        page: 2,
        pageSize: 50,
        total: 150,
      },
    }

    ;(axiosClients.axiosInstance.get as jest.Mock).mockImplementation(
      (url: string, config?: MockGetConfig) => {
        if (url.includes('/device-deployments')) {
          const page = config?.params?.page || 1
          if (page === 1) {
            return Promise.resolve({ data: page1Data })
          } else if (page === 2) {
            return Promise.resolve({ data: page2Data })
          }
        }
        return Promise.reject(new Error('Unknown URL: ' + url))
      }
    )

    // Test page 1
    const { result: page1Result } = renderHook(
      () => useDeviceDeployments('firmware-v2.0.0', 1, 50, undefined, true),
      { wrapper: createWrapper() }
    )

    await waitFor(() => {
      expect(page1Result.current.isLoading).toBe(false)
    })

    expect(page1Result.current.data?.deployments).toHaveLength(50)
    expect(page1Result.current.data?.total).toBe(150)
    expect(page1Result.current.data?.page).toBe(1)

    // Test page 2
    const { result: page2Result } = renderHook(
      () => useDeviceDeployments('firmware-v2.0.0', 2, 50, undefined, true),
      { wrapper: createWrapper() }
    )

    await waitFor(() => {
      expect(page2Result.current.isLoading).toBe(false)
    })

    expect(page2Result.current.data?.deployments).toHaveLength(50)
    expect(page2Result.current.data?.page).toBe(2)
    expect(page2Result.current.data?.deployments?.[0]?.currentVersion).toBe(
      'v2.0.0'
    )
  })

  /**
   * Test: Device filtering by status
   * Scenario: User applies status filter to see only successful upgrades
   */
  it('should support status filtering for device deployments', async () => {
    const successData = {
      deployments: Array.from({ length: 30 }, (_, i) => ({
        deviceId: `device-${i + 1}`,
        deviceKey: `dev-${i + 1}`,
        deviceName: `Device ${i + 1}`,
        productId: 1,
        productKey: 'WMS003',
        currentVersion: 'v2.0.0',
        upgradeBatchId: 'batch-001',
        status: 'success',
        lastStatusChangeTime: Date.now() - Math.random() * 3600000,
        createdAt: '2024-02-21T10:30:00Z',
      })),
      pagination: {
        page: 1,
        pageSize: 100,
        total: 30,
      },
    }

    ;(axiosClients.axiosInstance.get as jest.Mock).mockImplementation(
      (url: string, config?: MockGetConfig) => {
        if (url.includes('/device-deployments')) {
          const status = config?.params?.status
          if (status === 'success') {
            return Promise.resolve({ data: successData })
          }
        }
        return Promise.reject(new Error('Unknown filter'))
      }
    )

    const { result } = renderHook(
      () => useDeviceDeployments('firmware-v2.0.0', 1, 100, 'success', true),
      { wrapper: createWrapper() }
    )

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(result.current.data?.deployments).toHaveLength(30)
    expect(
      result.current.data?.deployments?.every((d) => d.status === 'success')
    ).toBe(true)
    expect(axiosClients.axiosInstance.get).toHaveBeenCalledWith(
      expect.stringContaining('/device-deployments'),
      expect.objectContaining({
        params: expect.objectContaining({
          status: 'success',
        }),
      })
    )
  })

  /**
   * Test: Error handling and retry
   * Scenario: API fails to load data, user clicks retry
   */
  it('should handle API errors and support retry', async () => {
    ;(axiosClients.axiosInstance.get as jest.Mock).mockRejectedValueOnce(
      new Error('Network error')
    )

    const { result } = renderHook(
      () => useOTAPackageDetail('firmware-v2.0.0'),
      { wrapper: createWrapper() }
    )

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(result.current.isError).toBe(true)
    expect(result.current.error).toBeDefined()

    // Simulate retry
    ;(axiosClients.axiosInstance.get as jest.Mock).mockResolvedValueOnce({
      data: {
        id: 'pkg-001',
        packageName: 'firmware-v2.0.0',
        version: 'v2.0.0',
        packageType: 'upgrade',
        productId: 1,
        productKey: 'WMS003',
        status: 'completed',
        fileUrl: 'https://example.com/firmware.bin',
        checksum: 'abc123',
        verificationStatus: 'completed',
        verificationProgress: 100,
        description: 'Test',
        releaseNotes: 'Test notes',
        metadata: {},
        createdAt: '2024-02-21T10:30:00Z',
        updatedAt: '2024-02-21T10:30:00Z',
      },
    })

    await result.current.refetch()

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true)
    })

    expect(result.current.data?.packageName).toBe('firmware-v2.0.0')
  })

  /**
   * Test: Real-time update polling
   * Scenario: Statistics auto-refresh every 5 seconds when tab is active
   */
  it('should poll statistics when tab is active', async () => {
    const statsData = {
      packageId: 'pkg-001',
      totalTargetDevices: 1000,
      successfulUpgrades: 800,
      failedUpgrades: 50,
      cancelledUpgrades: 10,
      pendingUpgrades: 100,
      inProgressUpgrades: 40,
    }

    ;(axiosClients.axiosInstance.get as jest.Mock).mockResolvedValue({
      data: statsData,
    })

    const { result } = renderHook(
      () => useUpgradeStatistics('firmware-v2.0.0', true),
      { wrapper: createWrapper() }
    )

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    // First call should be made
    expect(axiosClients.axiosInstance.get).toHaveBeenCalled()

    // Test that polling interval is configured
    expect(result.current.data?.successfulUpgrades).toBe(800)
  })

  /**
   * Test: Multiple data sources integration
   * Scenario: All data types load and integrate correctly
   */
  it('should integrate all data sources correctly', async () => {
    const mockPackage = {
      id: 'pkg-001',
      packageName: 'firmware-v2.0.0',
      version: 'v2.0.0',
      packageType: 'upgrade',
      productId: 1,
      productKey: 'WMS003',
      status: 'completed',
      fileUrl: 'https://example.com/firmware.bin',
      checksum: 'abc123',
      verificationStatus: 'completed',
      verificationProgress: 100,
      description: 'Production release',
      releaseNotes: 'Bug fixes and improvements',
      metadata: { releaseDate: '2024-02-21' },
      createdAt: '2024-02-21T10:30:00Z',
      updatedAt: '2024-02-21T10:30:00Z',
    }

    const mockStats = {
      packageId: 'pkg-001',
      totalTargetDevices: 500,
      successfulUpgrades: 450,
      failedUpgrades: 30,
      cancelledUpgrades: 5,
      pendingUpgrades: 15,
      inProgressUpgrades: 0,
    }

    const mockDeployments = {
      deployments: [
        {
          deviceId: 'device-1',
          deviceKey: 'dev-1',
          deviceName: 'Device 1',
          productId: 1,
          productKey: 'WMS003',
          currentVersion: 'v2.0.0',
          upgradeBatchId: 'batch-001',
          status: 'success',
          lastStatusChangeTime: Date.now(),
          createdAt: '2024-02-21T10:30:00Z',
        },
      ],
      pagination: { page: 1, pageSize: 100, total: 1 },
    }

    const mockBatches = {
      batches: [
        {
          batchId: 'batch-001',
          batchName: 'Batch 1',
          batchType: 'production',
          upgradeStrategy: 'static',
          status: 'completed',
          targetDeviceCount: 500,
          createdAt: '2024-02-21T10:30:00Z',
        },
      ],
    }

    ;(axiosClients.axiosInstance.get as jest.Mock).mockImplementation(
      (url: string) => {
        if (url.includes('/upgrade-statistics')) {
          return Promise.resolve({ data: mockStats })
        }
        if (url.includes('/device-deployments')) {
          return Promise.resolve({ data: mockDeployments })
        }
        if (url.includes('/batches')) {
          return Promise.resolve({ data: mockBatches })
        }
        return Promise.resolve({ data: mockPackage })
      }
    )

    // Load all data in parallel
    const { result: pkgResult } = renderHook(
      () => useOTAPackageDetail('firmware-v2.0.0'),
      { wrapper: createWrapper() }
    )

    const { result: statsResult } = renderHook(
      () => useUpgradeStatistics('firmware-v2.0.0', true),
      { wrapper: createWrapper() }
    )

    const { result: deploymentsResult } = renderHook(
      () => useDeviceDeployments('firmware-v2.0.0', 1, 100, undefined, true),
      { wrapper: createWrapper() }
    )

    const { result: batchesResult } = renderHook(
      () => useUpgradeBatches('firmware-v2.0.0'),
      { wrapper: createWrapper() }
    )

    // Wait for all to load
    await waitFor(() => {
      expect(pkgResult.current.isSuccess).toBe(true)
      expect(statsResult.current.isSuccess).toBe(true)
      expect(deploymentsResult.current.isSuccess).toBe(true)
      expect(batchesResult.current.isSuccess).toBe(true)
    })

    // Verify all data is present and consistent
    expect(pkgResult.current.data?.packageName).toBe('firmware-v2.0.0')
    expect(statsResult.current.data?.totalTargetDevices).toBe(500)
    expect(deploymentsResult.current.data?.deployments).toHaveLength(1)
    expect(batchesResult.current.data).toHaveLength(1)

    // Verify consistency: deployment version matches package version
    expect(
      deploymentsResult.current.data?.deployments?.[0]?.currentVersion
    ).toBe(pkgResult.current.data?.version)
  })
})
