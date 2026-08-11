import React, { ReactNode } from 'react'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import {
  useOTAPackageDetail,
  useUpgradeStatistics,
  useDeviceDeployments,
  useUpgradeBatches,
} from './detail-queries'
import * as axiosClients from '@/api/clients'

// Mock axios module
jest.mock('@/api/clients', () => ({
  axiosInstance: {
    get: jest.fn(),
  },
}))

// Mock react-query provider
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

describe('OTA Package Detail Hooks', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('useOTAPackageDetail should fetch data from API', async () => {
    const mockData = {
      id: 'pkg-001',
      packageName: 'test-package',
      version: 'v1.0.0',
      packageType: 'upgrade',
      productId: 1,
      productKey: 'WMS003',
      status: 'completed',
      fileUrl: 'https://example.com/firmware.bin',
      checksum: 'abc123',
      verificationStatus: 'completed',
      verificationProgress: 100,
      description: 'Test description',
      releaseNotes: 'Test notes',
      metadata: {},
      createdAt: '2024-02-21T10:30:00Z',
      updatedAt: '2024-02-21T10:30:00Z',
    }

    ;(axiosClients.axiosInstance.get as jest.Mock).mockResolvedValueOnce({
      data: mockData,
    })

    const { result } = renderHook(
      () => useOTAPackageDetail('test-package'),
      { wrapper: createWrapper() }
    )

    // Initially loading
    expect(result.current.isLoading).toBe(true)

    // Wait for data to load
    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    // Check data structure
    expect(result.current.data).toBeDefined()
    expect(result.current.data?.packageName).toBe('test-package')
    expect(result.current.data?.version).toBe('v1.0.0')
    expect(result.current.data?.packageType).toBe('upgrade')
    expect(axiosClients.axiosInstance.get).toHaveBeenCalledWith('/v1/ota-packages/test-package')
  })

  it('useUpgradeStatistics should fetch statistics from API', async () => {
    const mockStats = {
      packageId: 'test-package',
      totalTargetDevices: 100,
      successfulUpgrades: 45,
      failedUpgrades: 5,
      cancelledUpgrades: 2,
      pendingUpgrades: 30,
      inProgressUpgrades: 18,
    }

    ;(axiosClients.axiosInstance.get as jest.Mock).mockResolvedValueOnce({
      data: mockStats,
    })

    const { result } = renderHook(
      () => useUpgradeStatistics('test-package', true),
      { wrapper: createWrapper() }
    )

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(result.current.data).toBeDefined()
    expect(result.current.data?.totalTargetDevices).toBe(100)
    expect(result.current.data?.successfulUpgrades).toBe(45)
    expect(result.current.data?.failedUpgrades).toBe(5)
    expect(axiosClients.axiosInstance.get).toHaveBeenCalledWith('/v1/ota-packages/test-package/upgrade-statistics')
  })

  it('useDeviceDeployments should support pagination', async () => {
    const mockDeployments = {
      deployments: Array.from({ length: 50 }, (_, i) => ({
        deviceId: i + 1,
        deviceKey: `device-${i + 1}`,
        deviceName: `Device ${i + 1}`,
        productId: 1,
        productKey: 'WMS003',
        currentVersion: 'v0.9.0',
        upgradeBatchId: 'batch-001',
        status: 'pending',
        lastStatusChangeTime: 1708517400,
        createdAt: '2024-02-20T10:30:00Z',
      })),
      pagination: {
        page: 1,
        pageSize: 50,
        total: 100,
      },
    }

    ;(axiosClients.axiosInstance.get as jest.Mock).mockResolvedValueOnce({
      data: mockDeployments,
    })

    const { result } = renderHook(
      () => useDeviceDeployments('test-package', 1, 50, undefined, true),
      { wrapper: createWrapper() }
    )

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(result.current.data?.deployments).toBeDefined()
    expect(result.current.data?.deployments?.length).toBe(50)
    expect(result.current.data?.total).toBe(100)
    expect(result.current.data?.page).toBe(1)
    expect(axiosClients.axiosInstance.get).toHaveBeenCalledWith(
      '/v1/ota-packages/test-package/device-deployments',
      expect.objectContaining({
        params: expect.objectContaining({
          page: 1,
          pageSize: 50,
        }),
      })
    )
  })

  it('useDeviceDeployments should support filtering', async () => {
    const mockDeployments = {
      deployments: [
        {
          deviceId: 1,
          deviceKey: 'device-1',
          deviceName: 'Device 1',
          productId: 1,
          productKey: 'WMS003',
          currentVersion: 'v1.0.0',
          upgradeBatchId: 'batch-001',
          status: 'success',
          lastStatusChangeTime: 1708517400,
          createdAt: '2024-02-20T10:30:00Z',
        },
      ],
      pagination: {
        page: 1,
        pageSize: 100,
        total: 1,
      },
    }

    ;(axiosClients.axiosInstance.get as jest.Mock).mockResolvedValueOnce({
      data: mockDeployments,
    })

    const { result } = renderHook(
      () => useDeviceDeployments('test-package', 1, 100, 'success', true),
      { wrapper: createWrapper() }
    )

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(result.current.data?.deployments).toBeDefined()
    expect(result.current.data?.deployments?.length).toBe(1)
    expect(axiosClients.axiosInstance.get).toHaveBeenCalledWith(
      '/v1/ota-packages/test-package/device-deployments',
      expect.objectContaining({
        params: expect.objectContaining({
          status: 'success',
        }),
      })
    )
  })

  it('useUpgradeBatches should fetch batches from API', async () => {
    const mockBatches = {
      batches: [
        {
          batchId: 'batch-001',
          batchName: 'Verification Batch',
          batchType: 'verification',
          upgradeStrategy: 'static',
          status: 'completed',
          targetDeviceCount: 50,
          createdAt: '2024-02-15T10:30:00Z',
        },
      ],
    }

    ;(axiosClients.axiosInstance.get as jest.Mock).mockResolvedValueOnce({
      data: mockBatches,
    })

    const { result } = renderHook(
      () => useUpgradeBatches('test-package'),
      { wrapper: createWrapper() }
    )

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(result.current.data).toBeDefined()
    expect(result.current.data?.length).toBe(1)
    expect(result.current.data?.[0]?.batchId).toBe('batch-001')
    expect(axiosClients.axiosInstance.get).toHaveBeenCalledWith('/v1/ota-packages/test-package/batches')
  })
})
