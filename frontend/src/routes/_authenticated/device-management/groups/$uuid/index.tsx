import { createFileRoute } from '@tanstack/react-router'
import { DeviceGroupDetail } from '@/features/device-groups/detail'

export const Route = createFileRoute(
  '/_authenticated/device-management/groups/$uuid/'
)({
  component: DeviceGroupDetail,
})
