import { useQueryClient } from '@tanstack/react-query'
import {
  useGetOrganizationsMembers,
  useGetOrganizationsInvitations,
  useGetOrganizationsInvitationsId,
  usePostOrganizationsInvitations,
  usePostOrganizationsInvitationsIdAccept,
  useDeleteOrganizationsInvitationsId,
  useDeleteOrganizationsMembersUserId,
  getGetOrganizationsMembersQueryKey,
  getGetOrganizationsInvitationsQueryKey,
} from '@/api/generated'
import type {
  GetOrganizationsInvitationsParams,
  GetOrganizationsMembersParams,
} from '@/api/generated/model'

export function useOrganizationMembers(params: GetOrganizationsMembersParams) {
  const query = useGetOrganizationsMembers(params)
  return {
    ...query,
    data: query.data?.data,
  }
}

export function useOrganizationInvitations(
  params: GetOrganizationsInvitationsParams
) {
  const query = useGetOrganizationsInvitations(params)
  return {
    ...query,
    data: query.data?.data,
  }
}

export function useInvitation(id: string) {
  const query = useGetOrganizationsInvitationsId(id, {
    query: { enabled: !!id },
  })
  return {
    ...query,
    data: query.data?.data,
  }
}

export function useAcceptOrganizationInvitation() {
  const queryClient = useQueryClient()
  return usePostOrganizationsInvitationsIdAccept({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: getGetOrganizationsMembersQueryKey(),
        })
        queryClient.invalidateQueries({
          queryKey: getGetOrganizationsInvitationsQueryKey(),
        })
      },
    },
  })
}

export function useCreateOrganizationInvitation() {
  const queryClient = useQueryClient()
  return usePostOrganizationsInvitations({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: getGetOrganizationsInvitationsQueryKey(),
        })
      },
    },
  })
}

export function useRevokeOrganizationInvitation() {
  const queryClient = useQueryClient()
  return useDeleteOrganizationsInvitationsId({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: getGetOrganizationsInvitationsQueryKey(),
        })
      },
    },
  })
}

export function useRemoveOrganizationMember() {
  const queryClient = useQueryClient()
  return useDeleteOrganizationsMembersUserId({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: getGetOrganizationsMembersQueryKey(),
        })
      },
    },
  })
}
