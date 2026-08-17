import { createFileRoute } from '@tanstack/react-router'
import { SceneLinkageEditPage } from '@/features/scene-linkage/scene-linkage-edit-page'

export const Route = createFileRoute(
  '/_authenticated/rule-engine/scene-linkage/$sceneId/'
)({ component: SceneLinkageEditPage })
