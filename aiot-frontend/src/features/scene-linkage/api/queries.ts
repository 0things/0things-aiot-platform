import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type {
  ListSceneLinkagesResponse,
  SceneLinkage,
  SceneLinkageRequest,
  SceneLinkageDetail,
  SceneLinkageDetailRequest,
} from '@/api/generated/model'
import {
  getSceneLinkages,
  getSceneLinkagesId,
  getSceneLinkagesIdDetail,
  postSceneLinkages,
  putSceneLinkagesId,
  deleteSceneLinkagesId,
  postSceneLinkagesIdDetail,
  putSceneLinkagesIdDetail,
} from '@/api/generated'

export const sceneLinkageKeys = {
  all: ['scene-linkages'] as const,
  lists: () => [...sceneLinkageKeys.all, 'list'] as const,
  list: (filters: {
    page?: number
    pageSize?: number
    search?: string
    enable?: number
  }) => [...sceneLinkageKeys.lists(), filters] as const,
  details: () => [...sceneLinkageKeys.all, 'detail'] as const,
  detail: (id: number) => [...sceneLinkageKeys.details(), id] as const,
  detailConfig: (id: number) =>
    [...sceneLinkageKeys.all, 'detail-config', id] as const,
}

export function useSceneLinkages(params: {
  page?: number
  pageSize?: number
  search?: string
  enable?: number
}) {
  return useQuery({
    queryKey: sceneLinkageKeys.list(params),
    queryFn: async () => {
      return getSceneLinkages(params) as unknown as ListSceneLinkagesResponse
    },
  })
}

export function useSceneLinkage(id: number) {
  return useQuery({
    queryKey: sceneLinkageKeys.detail(id),
    queryFn: async () =>
      (await getSceneLinkagesId(id)).sceneLinkage as SceneLinkage,
    enabled: Number.isFinite(id) && id > 0,
  })
}

export function useSceneLinkageDetail(id: number) {
  return useQuery({
    queryKey: sceneLinkageKeys.detailConfig(id),
    queryFn: async () =>
      (await getSceneLinkagesIdDetail(id)).detail as SceneLinkageDetail,
    enabled: Number.isFinite(id) && id > 0,
  })
}

export function useCreateSceneLinkage() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (data: SceneLinkageRequest) =>
      (await postSceneLinkages(data)).sceneLinkage,
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: sceneLinkageKeys.lists() }),
  })
}

export function useUpdateSceneLinkage() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, data }: { id: number; data: SceneLinkageRequest }) =>
      (await putSceneLinkagesId(id, data)).sceneLinkage,
    onSuccess: () => qc.invalidateQueries({ queryKey: sceneLinkageKeys.all }),
  })
}

export function useDeleteSceneLinkage() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (id: number) => (await deleteSceneLinkagesId(id)).success,
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
    }) => (await postSceneLinkagesIdDetail(id, data)).detail,
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
    }) => (await putSceneLinkagesIdDetail(id, data)).detail,
    onSuccess: (_r, v) =>
      qc.invalidateQueries({ queryKey: sceneLinkageKeys.detailConfig(v.id) }),
  })
}
