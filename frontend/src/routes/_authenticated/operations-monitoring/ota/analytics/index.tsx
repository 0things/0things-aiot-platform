import { createFileRoute } from '@tanstack/react-router'
import { OTAAnalytics } from '@/features/operations-monitoring/ota'

export const Route = createFileRoute(
  '/_authenticated/operations-monitoring/ota/analytics/'
)({
  component: OTAAnalytics,
})
