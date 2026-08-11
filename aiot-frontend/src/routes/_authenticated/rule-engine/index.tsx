import { createFileRoute } from '@tanstack/react-router'
import RuleEnginePage from '@/features/rule-engine'

export const Route = createFileRoute('/_authenticated/rule-engine/')({
  component: RuleEnginePage,
})
