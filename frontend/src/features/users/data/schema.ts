import { z } from 'zod'
import type {
  AiotBackendApiV1OrganizationMemberResp,
  AiotBackendApiV1OrganizationInvitationResp,
} from '@/api/generated/model'

export type MemberUser = AiotBackendApiV1OrganizationMemberResp
export type InvitationItem = AiotBackendApiV1OrganizationInvitationResp

export const userStatusSchema = z.union([
  z.literal('active'),
  z.literal('inactive'),
  z.literal('invited'),
  z.literal('suspended'),
])
export type UserStatus = z.infer<typeof userStatusSchema>

export const userSchema = z.object({
  id: z.string(),
  firstName: z.string(),
  lastName: z.string(),
  username: z.string(),
  email: z.string(),
  phoneNumber: z.string(),
  status: userStatusSchema,
  role: z.string().optional(),
  createdAt: z.coerce.date(),
  updatedAt: z.coerce.date(),
})
export type User = z.infer<typeof userSchema>

export const userListSchema = z.array(userSchema)

/**
 * Maps raw organization member DTO from Logto/Backend into the UI User domain model.
 */
export function mapMemberToUser(member: MemberUser): User {
  const name =
    member.name || member.username || member.primaryEmail || member.id || 'User'
  const nameParts = name.trim().split(/\s+/)
  const firstName = nameParts[0] || 'User'
  const lastName = nameParts.slice(1).join(' ') || ''
  const username =
    member.username ||
    member.name ||
    (member.primaryEmail ? member.primaryEmail.split('@')[0] : member.id || '')

  return {
    id: member.id || '',
    firstName,
    lastName,
    username,
    email: member.primaryEmail || '',
    phoneNumber: member.primaryPhone || '-',
    status: 'active',
    createdAt: member.createdAt ? new Date(member.createdAt) : new Date(),
    updatedAt: member.createdAt ? new Date(member.createdAt) : new Date(),
  }
}
