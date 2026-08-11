import { createFileRoute } from '@tanstack/react-router'
import { DeviceDetailPage } from '@/features/devices/components/device-detail/device-detail-page'

export const Route = createFileRoute(
  '/_authenticated/device-management/devices/$deviceKey/'
)({
  component: DeviceDetailPage,
})
