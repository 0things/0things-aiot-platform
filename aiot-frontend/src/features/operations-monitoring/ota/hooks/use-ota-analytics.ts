import { useQuery } from '@tanstack/react-query'
import type { OTAAnalyticsData } from '../data/analytics-schema'

// Mock data generator for analytics
function generateMockAnalyticsData(): OTAAnalyticsData {
  // Generate timeline data for the last 30 days
  const timeline = Array.from({ length: 30 }, (_, i) => {
    const date = new Date()
    date.setDate(date.getDate() - (29 - i))
    const deployments = Math.floor(Math.random() * 15) + 5
    const successes = Math.floor(deployments * (0.7 + Math.random() * 0.25))
    const failures = deployments - successes

    return {
      date: date.toISOString().split('T')[0],
      deployments,
      successes,
      failures,
    }
  })

  const totalDeployments = timeline.reduce(
    (sum, day) => sum + day.deployments,
    0
  )
  const totalSuccesses = timeline.reduce((sum, day) => sum + day.successes, 0)
  const totalFailures = timeline.reduce((sum, day) => sum + day.failures, 0)

  return {
    summary: {
      totalPackages: 24,
      activeDeployments: 8,
      successRate: Math.round((totalSuccesses / totalDeployments) * 100),
      failedDeployments: totalFailures,
    },
    deploymentStatus: [
      {
        status: 'success',
        count: totalSuccesses,
        percentage: Math.round((totalSuccesses / totalDeployments) * 100),
      },
      {
        status: 'failed',
        count: totalFailures,
        percentage: Math.round((totalFailures / totalDeployments) * 100),
      },
      {
        status: 'in_progress',
        count: 8,
        percentage: Math.round((8 / (totalDeployments + 8)) * 100),
      },
    ],
    timeline,
    firmwareDistribution: [
      { version: '2.1.0', deviceCount: 450, percentage: 45 },
      { version: '2.0.5', deviceCount: 280, percentage: 28 },
      { version: '1.9.8', deviceCount: 150, percentage: 15 },
      { version: '1.8.2', deviceCount: 80, percentage: 8 },
      { version: '1.7.5', deviceCount: 40, percentage: 4 },
    ],
    recentActivity: [
      {
        id: '1',
        packageName: 'firmware-v2.1.0',
        version: '2.1.0',
        action: 'completed',
        productName: 'WMS003',
        timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
        status: 'success',
      },
      {
        id: '2',
        packageName: 'camera-firmware-v1.5.2',
        version: '1.5.2',
        action: 'deployed',
        productName: 'Camera-001',
        timestamp: new Date(Date.now() - 4 * 60 * 60 * 1000).toISOString(),
        status: 'in_progress',
      },
      {
        id: '3',
        packageName: 'sensor-patch-v1.2.3',
        version: '1.2.3',
        action: 'failed',
        productName: 'WMS003',
        timestamp: new Date(Date.now() - 6 * 60 * 60 * 1000).toISOString(),
        status: 'failed',
      },
      {
        id: '4',
        packageName: 'firmware-v2.0.9',
        version: '2.0.9',
        action: 'created',
        productName: 'WMS003',
        timestamp: new Date(Date.now() - 8 * 60 * 60 * 1000).toISOString(),
        status: 'pending',
      },
      {
        id: '5',
        packageName: 'bootloader-v3.1.0',
        version: '3.1.0',
        action: 'completed',
        productName: 'WMS003',
        timestamp: new Date(Date.now() - 12 * 60 * 60 * 1000).toISOString(),
        status: 'success',
      },
    ],
  }
}

async function fetchOTAAnalytics(): Promise<OTAAnalyticsData> {
  // Simulate API delay
  await new Promise((resolve) => setTimeout(resolve, 800))
  return generateMockAnalyticsData()
}

export function useOTAAnalytics() {
  return useQuery({
    queryKey: ['ota-analytics'],
    queryFn: fetchOTAAnalytics,
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchInterval: 30 * 1000, // Refetch every 30 seconds for real-time feel
  })
}
