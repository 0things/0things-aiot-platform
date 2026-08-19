import { createFileRoute } from '@tanstack/react-router'
import { OTAPackageDetailPage } from '@/features/operations-monitoring/ota/components/ota-package-detail-page'

export const Route = createFileRoute(
  '/_authenticated/operations-monitoring/ota/packages/$id/'
)({
  component: OTAPackageDetailPage,
})
