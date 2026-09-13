import { createFileRoute } from '@tanstack/react-router'
import { RuleChains } from '@/features/rule-chains'

export const Route = createFileRoute(
  '/_authenticated/rule-engine/rule-chains/'
)({
  component: RuleChains,
})
