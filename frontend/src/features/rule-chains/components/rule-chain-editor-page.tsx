import { getRouteApi } from '@tanstack/react-router'
import { RuleChainCanvasPage } from '../rule-chain-canvas'

const route = getRouteApi('/_authenticated/rule-engine/rule-chains/$uuid/edit/')

export function RuleChainEditorPage() {
  const { uuid } = route.useParams()
  return <RuleChainCanvasPage uuid={uuid} />
}
