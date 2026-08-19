import { createFileRoute } from '@tanstack/react-router'
import { Devices } from '@/features/devices'

export const Route = createFileRoute(
  '/_authenticated/device-management/devices/'
)({
  component: Devices,
})
