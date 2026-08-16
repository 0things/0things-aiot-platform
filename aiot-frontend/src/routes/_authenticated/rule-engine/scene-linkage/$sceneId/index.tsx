import { createFileRoute } from '@tanstack/react-router'
import { SceneLinkageDetailPage } from '@/features/scene-linkage'

export const Route = createFileRoute(
  '/_authenticated/rule-engine/scene-linkage/$sceneId/'
)({ component: SceneLinkageDetailPage })
