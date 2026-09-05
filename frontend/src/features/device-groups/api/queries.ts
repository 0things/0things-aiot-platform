import { useMutation, useQueryClient } from '@tanstack/react-query'
import {
  deleteDeviceGroupsGroupUuid,
  deleteDeviceGroupsGroupUuidDevices,
  getGetDeviceGroupsGroupUuidDevicesQueryKey,
  getGetDeviceGroupsGroupUuidQueryKey,
  getGetDeviceGroupsQueryKey,
  postDeviceGroups,
  postDeviceGroupsGroupUuidDevices,
  putDeviceGroupsGroupUuid,
  useGetDeviceGroups,
  useGetDeviceGroupsGroupUuid,
  useGetDeviceGroupsGroupUuidDevices,
} from '@/api/generated'
import type {
  AiotBackendApiDeviceGroupV1CreateDeviceGroupRequest,
  AiotBackendApiDeviceGroupV1DeviceKeysRequest,
  AiotBackendApiDeviceGroupV1UpdateDeviceGroupRequest,
  GetDeviceGroupsGroupUuidDevicesParams,
  GetDeviceGroupsParams,
} from '@/api/generated/model'

export type { GetDeviceGroupsGroupUuidDevicesParams, GetDeviceGroupsParams }

export const deviceGroupKeys = {
  all: ['device-groups'] as const,
  lists: () => [...deviceGroupKeys.all, 'list'] as const,
  list: (params?: GetDeviceGroupsParams) => getGetDeviceGroupsQueryKey(params),
  details: () => [...deviceGroupKeys.all, 'detail'] as const,
  detail: (groupUuid: string) => getGetDeviceGroupsGroupUuidQueryKey(groupUuid),
  devices: (
    groupUuid: string,
    params?: GetDeviceGroupsGroupUuidDevicesParams
  ) => getGetDeviceGroupsGroupUuidDevicesQueryKey(groupUuid, params),
}

export function useDeviceGroups(params?: GetDeviceGroupsParams) {
  return useGetDeviceGroups(params, {
    query: {
      select: (res) => res?.data,
    },
  })
}

export function useDeviceGroup(groupUuid: string) {
  return useGetDeviceGroupsGroupUuid(groupUuid, {
    query: {
      select: (res) => res?.data,
      enabled: !!groupUuid,
    },
  })
}

export function useDeviceGroupDevices(
  groupUuid: string,
  params?: GetDeviceGroupsGroupUuidDevicesParams
) {
  return useGetDeviceGroupsGroupUuidDevices(groupUuid, params, {
    query: {
      select: (res) => res?.data,
      enabled: !!groupUuid,
    },
  })
}

export function useCreateDeviceGroup() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: AiotBackendApiDeviceGroupV1CreateDeviceGroupRequest) =>
      postDeviceGroups(data as never),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: deviceGroupKeys.lists() })
    },
  })
}

export function useUpdateDeviceGroup() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({
      groupUuid,
      data,
    }: {
      groupUuid: string
      data: AiotBackendApiDeviceGroupV1UpdateDeviceGroupRequest
    }) => putDeviceGroupsGroupUuid(groupUuid, data as never),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: deviceGroupKeys.detail(variables.groupUuid),
      })
      queryClient.invalidateQueries({ queryKey: deviceGroupKeys.lists() })
    },
  })
}

export function useDeleteDeviceGroup() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (groupUuid: string) => deleteDeviceGroupsGroupUuid(groupUuid),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: deviceGroupKeys.lists() })
    },
  })
}

export function useAddDevicesToGroup() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({
      groupUuid,
      data,
    }: {
      groupUuid: string
      data: AiotBackendApiDeviceGroupV1DeviceKeysRequest
    }) => postDeviceGroupsGroupUuidDevices(groupUuid, data as never),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: getGetDeviceGroupsGroupUuidDevicesQueryKey(
          variables.groupUuid
        ),
      })
      queryClient.invalidateQueries({
        queryKey: deviceGroupKeys.detail(variables.groupUuid),
      })
    },
  })
}

export function useRemoveDevicesFromGroup() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({
      groupUuid,
      data,
    }: {
      groupUuid: string
      data: AiotBackendApiDeviceGroupV1DeviceKeysRequest
    }) => deleteDeviceGroupsGroupUuidDevices(groupUuid, data as never),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: getGetDeviceGroupsGroupUuidDevicesQueryKey(
          variables.groupUuid
        ),
      })
      queryClient.invalidateQueries({
        queryKey: deviceGroupKeys.detail(variables.groupUuid),
      })
    },
  })
}
