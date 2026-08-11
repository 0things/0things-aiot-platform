import { createFileRoute } from '@tanstack/react-router'
import { IotDashboard } from '@/features/iot-dashboard'

export const Route = createFileRoute('/_authenticated/iot-dashboard/')({
  component: IotDashboard,
})
