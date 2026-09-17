import { setAuthToken } from '@/api/clients'
import { postMeOrganization } from '@/api/generated'

export async function ensureCurrentUserOrganization(
  identityToken: string
): Promise<string> {
  setAuthToken(identityToken)
  const response = await postMeOrganization()
  const organizationId = response.data?.organizationId
  if (!organizationId) {
    throw new Error('Failed to retrieve organizationId from response')
  }
  return organizationId
}
