import { createFileRoute } from '@tanstack/react-router'
import { DeviceEvents } from '@/features/operations-monitoring/events'

export const Route = createFileRoute(
  '/_authenticated/operations-monitoring/events/'
)({ component: DeviceEvents })
