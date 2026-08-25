import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
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
import type {
  ListSceneLinkagesResponse,
  SceneLinkage,
  SceneLinkageRequest,
  SceneLinkageDetail,
  SceneLinkageDetailRequest,
  GetSceneLinkageResponse,
  CreateSceneLinkageResponse,
  UpdateSceneLinkageResponse,
  GetSceneLinkageDetailResponse,
  CreateSceneLinkageDetailResponse,
  UpdateSceneLinkageDetailResponse,
  SceneLinkageSuccessResponse,
} from '@/api/generated/model'

// The backend returns a generic { code, message, data } envelope; the generated
// client types each endpoint's response as its flat payload, so we unwrap .data.
type ApiResponse<T> = { code?: number; message?: string; data?: T }

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
      const res = await getSceneLinkages(params)
      return ((res as unknown as ApiResponse<ListSceneLinkagesResponse>).data ??
        {}) as ListSceneLinkagesResponse
    },
  })
}

export function useSceneLinkage(id: number) {
  return useQuery({
    queryKey: sceneLinkageKeys.detail(id),
    queryFn: async () =>
      (
        (await getSceneLinkagesId(
          id
        )) as unknown as ApiResponse<GetSceneLinkageResponse>
      ).data?.sceneLinkage as SceneLinkage,
    enabled: Number.isFinite(id) && id > 0,
  })
}

export function useSceneLinkageDetail(id: number) {
  return useQuery({
    queryKey: sceneLinkageKeys.detailConfig(id),
    queryFn: async () =>
      (
        (await getSceneLinkagesIdDetail(
          id
        )) as unknown as ApiResponse<GetSceneLinkageDetailResponse>
      ).data?.detail as SceneLinkageDetail,
    enabled: Number.isFinite(id) && id > 0,
  })
}

export function useCreateSceneLinkage() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (data: SceneLinkageRequest) =>
      (
        (await postSceneLinkages(
          data
        )) as unknown as ApiResponse<CreateSceneLinkageResponse>
      ).data?.sceneLinkage,
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
    }) =>
      (
        (await putSceneLinkagesId(
          id,
          data
        )) as unknown as ApiResponse<UpdateSceneLinkageResponse>
      ).data?.sceneLinkage,
    onSuccess: () => qc.invalidateQueries({ queryKey: sceneLinkageKeys.all }),
  })
}

export function useDeleteSceneLinkage() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (id: number) =>
      (
        (await deleteSceneLinkagesId(
          id
        )) as unknown as ApiResponse<SceneLinkageSuccessResponse>
      ).data?.success,
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
    }) =>
      (
        (await postSceneLinkagesIdDetail(
          id,
          data
        )) as unknown as ApiResponse<CreateSceneLinkageDetailResponse>
      ).data?.detail,
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
    }) =>
      (
        (await putSceneLinkagesIdDetail(
          id,
          data
        )) as unknown as ApiResponse<UpdateSceneLinkageDetailResponse>
      ).data?.detail,
    onSuccess: (_r, v) =>
      qc.invalidateQueries({ queryKey: sceneLinkageKeys.detailConfig(v.id) }),
  })
}
