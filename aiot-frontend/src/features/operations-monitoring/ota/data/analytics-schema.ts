import { z } from 'zod'

export const otaAnalyticsSummarySchema = z.object({
  totalPackages: z.number(),
  activeDeployments: z.number(),
  successRate: z.number(), // percentage
  failedDeployments: z.number(),
})

export const otaDeploymentStatusSchema = z.object({
  status: z.enum(['success', 'failed', 'in_progress']),
  count: z.number(),
  percentage: z.number(),
})

export const otaTimelineDataPointSchema = z.object({
  date: z.string(),
  deployments: z.number(),
  successes: z.number(),
  failures: z.number(),
})

export const otaFirmwareDistributionSchema = z.object({
  version: z.string(),
  deviceCount: z.number(),
  percentage: z.number(),
})

export const otaRecentActivitySchema = z.object({
  id: z.string(),
  packageName: z.string(),
  version: z.string(),
  action: z.enum(['created', 'deployed', 'completed', 'failed']),
  productName: z.string(),
  timestamp: z.string(),
  status: z.enum(['success', 'failed', 'in_progress', 'pending']),
})

export const otaAnalyticsDataSchema = z.object({
  summary: otaAnalyticsSummarySchema,
  deploymentStatus: z.array(otaDeploymentStatusSchema),
  timeline: z.array(otaTimelineDataPointSchema),
  firmwareDistribution: z.array(otaFirmwareDistributionSchema),
  recentActivity: z.array(otaRecentActivitySchema),
})

export type OTAAnalyticsSummary = z.infer<typeof otaAnalyticsSummarySchema>
export type OTADeploymentStatus = z.infer<typeof otaDeploymentStatusSchema>
export type OTATimelineDataPoint = z.infer<typeof otaTimelineDataPointSchema>
export type OTAFirmwareDistribution = z.infer<
  typeof otaFirmwareDistributionSchema
>
export type OTARecentActivity = z.infer<typeof otaRecentActivitySchema>
export type OTAAnalyticsData = z.infer<typeof otaAnalyticsDataSchema>
