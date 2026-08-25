import {
  getOrganizations,
  postAuthSwitchOrg,
  useGetOrganizations,
  usePostAuthSwitchOrg,
} from '@/api/generated'
import type {
  ApiOrganizationItem,
  ApiSwitchOrgRequest,
  ApiSwitchOrgResponseData,
} from '@/api/generated/model'

export const getMyOrganizations = () => getOrganizations()

export const switchOrganization = (data: ApiSwitchOrgRequest) =>
  postAuthSwitchOrg(data)

export { useGetOrganizations, usePostAuthSwitchOrg }
export type {
  ApiOrganizationItem,
  ApiSwitchOrgRequest,
  ApiSwitchOrgResponseData,
}
