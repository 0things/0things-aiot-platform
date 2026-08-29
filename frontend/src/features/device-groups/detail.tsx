import { getRouteApi } from '@tanstack/react-router'
import { GroupDetailPage } from './components/group-detail/group-detail-page'

const route = getRouteApi('/_authenticated/device-management/groups/$uuid/')

export function DeviceGroupDetail() {
  const { uuid } = route.useParams()
  return <GroupDetailPage uuid={uuid} />
}
