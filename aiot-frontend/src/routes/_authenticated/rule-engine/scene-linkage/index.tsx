import { createFileRoute } from '@tanstack/react-router'
import { SceneLinkageListPage } from '@/features/scene-linkage/list'

export const Route = createFileRoute(
  '/_authenticated/rule-engine/scene-linkage/'
)({ component: SceneLinkageListPage })
