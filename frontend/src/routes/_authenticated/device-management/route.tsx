import { createFileRoute } from '@tanstack/react-router'
import { DeviceManagement } from '@/features/device-management'

export const Route = createFileRoute('/_authenticated/device-management')({
  component: DeviceManagement,
})
