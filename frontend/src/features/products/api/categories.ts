import { useQuery } from '@tanstack/react-query'
import { getCategoriesTree } from '@/api/generated'

export type CategoryNode = {
  id?: number
  name?: string
  children?: CategoryNode[]
}

export function useProductCategories() {
  return useQuery({
    queryKey: ['product-categories', 'tree'],
    queryFn: async () => {
      const response = await getCategoriesTree()
      return (response?.data ?? []) as CategoryNode[]
    },
    staleTime: 5 * 60 * 1000,
  })
}
