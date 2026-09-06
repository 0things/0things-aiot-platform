import { createFileRoute } from '@tanstack/react-router'
import { RuleChainCanvasPage } from '@/features/rule-chains/rule-chain-canvas'

export const Route = createFileRoute(
  '/_authenticated/rule-engine/rule-chains/new/'
)({ component: RuleChainCanvasPage })
