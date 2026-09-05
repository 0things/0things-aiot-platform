import {
  getGetCategoriesTreeQueryKey,
  useGetCategoriesTree,
} from '@/api/generated'
import type { AiotBackendApiCategoryV1Category } from '@/api/generated/model'

export type CategoryNode = AiotBackendApiCategoryV1Category

export const categoryKeys = {
  all: ['product-categories'] as const,
  tree: () => getGetCategoriesTreeQueryKey(),
}

export function useProductCategories() {
  return useGetCategoriesTree({
    query: {
      select: (res) => (res?.data ?? []) as CategoryNode[],
      staleTime: 5 * 60 * 1000,
    },
  })
}
