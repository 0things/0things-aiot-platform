import { useMutation, useQueryClient } from '@tanstack/react-query'
import {
  deleteSceneLinkagesId,
  getGetSceneLinkagesIdDetailQueryKey,
  getGetSceneLinkagesIdQueryKey,
  getGetSceneLinkagesQueryKey,
  postSceneLinkages,
  postSceneLinkagesIdDetail,
  putSceneLinkagesId,
  putSceneLinkagesIdDetail,
  useGetSceneLinkages,
  useGetSceneLinkagesId,
  useGetSceneLinkagesIdDetail,
} from '@/api/generated'
import type {
  CreateSceneLinkageDetailResponse,
  CreateSceneLinkageResponse,
  GetSceneLinkagesParams,
  ListSceneLinkagesResponse,
  SceneLinkage,
  SceneLinkageDetail,
  SceneLinkageDetailRequest,
  SceneLinkageRequest,
  SceneLinkageSuccessResponse,
  UpdateSceneLinkageDetailResponse,
  UpdateSceneLinkageResponse,
} from '@/api/generated/model'

export type { GetSceneLinkagesParams }

export const sceneLinkageKeys = {
  all: ['scene-linkages'] as const,
  lists: () => [...sceneLinkageKeys.all, 'list'] as const,
  list: (params?: GetSceneLinkagesParams) =>
    getGetSceneLinkagesQueryKey(params),
  details: () => [...sceneLinkageKeys.all, 'detail'] as const,
  detail: (id: number) => getGetSceneLinkagesIdQueryKey(id),
  detailConfig: (id: number) => getGetSceneLinkagesIdDetailQueryKey(id),
}

export function useSceneLinkages(params?: GetSceneLinkagesParams) {
  return useGetSceneLinkages(params, {
    query: {
      select: (res) => (res?.data ?? {}) as ListSceneLinkagesResponse,
    },
  })
}

export function useSceneLinkage(id: number) {
  return useGetSceneLinkagesId(id, {
    query: {
      select: (res) => res?.data?.sceneLinkage as SceneLinkage,
      enabled: Number.isFinite(id) && id > 0,
    },
  })
}

export function useSceneLinkageDetail(id: number) {
  return useGetSceneLinkagesIdDetail(id, {
    query: {
      select: (res) => res?.data?.detail as SceneLinkageDetail,
      enabled: Number.isFinite(id) && id > 0,
    },
  })
}

export function useCreateSceneLinkage() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (data: SceneLinkageRequest) => {
      const res = await postSceneLinkages(data)
      return (res?.data as CreateSceneLinkageResponse)?.sceneLinkage
    },
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: sceneLinkageKeys.lists() }),
  })
}

export function useUpdateSceneLinkage() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async ({
      id,
      data,
    }: {
      id: number
      data: SceneLinkageRequest
    }) => {
      const res = await putSceneLinkagesId(id, data)
      return (res?.data as UpdateSceneLinkageResponse)?.sceneLinkage
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: sceneLinkageKeys.all }),
  })
}

export function useDeleteSceneLinkage() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (id: number) => {
      const res = await deleteSceneLinkagesId(id)
      return (res?.data as SceneLinkageSuccessResponse)?.success
    },
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: sceneLinkageKeys.lists() }),
  })
}

export function useCreateSceneLinkageDetail() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async ({
      id,
      data,
    }: {
      id: number
      data: SceneLinkageDetailRequest
    }) => {
      const res = await postSceneLinkagesIdDetail(id, data)
      return (res?.data as CreateSceneLinkageDetailResponse)?.detail
    },
    onSuccess: (_r, v) =>
      qc.invalidateQueries({ queryKey: sceneLinkageKeys.detailConfig(v.id) }),
  })
}

export function useUpdateSceneLinkageDetail() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async ({
      id,
      data,
    }: {
      id: number
      data: SceneLinkageDetailRequest
    }) => {
      const res = await putSceneLinkagesIdDetail(id, data)
      return (res?.data as UpdateSceneLinkageDetailResponse)?.detail
    },
    onSuccess: (_r, v) =>
      qc.invalidateQueries({ queryKey: sceneLinkageKeys.detailConfig(v.id) }),
  })
}
