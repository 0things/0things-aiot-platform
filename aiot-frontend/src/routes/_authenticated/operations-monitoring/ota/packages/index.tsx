import { createFileRoute } from '@tanstack/react-router'
import { OTAPackages } from '@/features/operations-monitoring/ota'

export const Route = createFileRoute(
  '/_authenticated/operations-monitoring/ota/packages/'
)({
  component: OTAPackages,
})
