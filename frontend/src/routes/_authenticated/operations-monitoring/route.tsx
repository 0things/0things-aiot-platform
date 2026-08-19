import { createFileRoute } from '@tanstack/react-router'
import { OperationsMonitoring } from '@/features/operations-monitoring'

export const Route = createFileRoute('/_authenticated/operations-monitoring')({
  component: OperationsMonitoring,
})
