import { createFileRoute } from '@tanstack/react-router'
import { DeviceGroups } from '@/features/device-groups'

export const Route = createFileRoute(
  '/_authenticated/device-management/groups/'
)({
  component: DeviceGroups,
})
