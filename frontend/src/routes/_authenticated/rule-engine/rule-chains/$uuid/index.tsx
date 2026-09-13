import { createFileRoute, getRouteApi } from '@tanstack/react-router'
import { RuleChainDetailPage } from '@/features/rule-chains/components/rule-chain-detail-page'

export const Route = createFileRoute(
  '/_authenticated/rule-engine/rule-chains/$uuid/'
)({
  component: () => {
    const { uuid } = getRouteApi(
      '/_authenticated/rule-engine/rule-chains/$uuid/'
    ).useParams()
    return <RuleChainDetailPage uuid={uuid} />
  },
})
