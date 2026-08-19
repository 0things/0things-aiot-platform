import { useParams } from '@tanstack/react-router'
import { SceneLinkageForm } from './index'

export function SceneLinkageEditPage() {
  const { sceneId } = useParams({
    from: '/_authenticated/rule-engine/scene-linkage/$sceneId/',
  })
  return <SceneLinkageForm mode='edit' sceneId={Number(sceneId)} />
}
