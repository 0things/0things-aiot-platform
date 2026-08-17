import { createFileRoute } from '@tanstack/react-router'
import { SceneLinkageForm } from '@/features/scene-linkage'

export const Route = createFileRoute(
  '/_authenticated/rule-engine/scene-linkage/new/'
)({
  component: () => <SceneLinkageForm mode='create' />,
})
