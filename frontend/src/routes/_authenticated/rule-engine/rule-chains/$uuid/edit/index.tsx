import { createFileRoute } from '@tanstack/react-router'
import { RuleChainEditorPage } from '@/features/rule-chains/components/rule-chain-editor-page'

export const Route = createFileRoute(
  '/_authenticated/rule-engine/rule-chains/$uuid/edit/'
)({ component: RuleChainEditorPage })
