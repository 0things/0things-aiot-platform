import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { productServiceApi } from '@/api/clients'
import type {
  ProductV1CreateProductRequest,
  ProductV1UpdateProductRequest,
} from '@/api/generated/device-service'

// ============================================================================
// Query Keys
// ============================================================================

export const productKeys = {
  all: ['products'] as const,
  lists: () => [...productKeys.all, 'list'] as const,
  list: (filters: {
    page?: number
    pageSize?: number
    category?: string
    status?: string
    searchText?: string
  }) => [...productKeys.lists(), filters] as const,
  details: () => [...productKeys.all, 'detail'] as const,
  detail: (productKey: string) =>
    [...productKeys.details(), productKey] as const,
}

// ============================================================================
// Queries
// ============================================================================

/**
 * Hook to fetch products list with pagination and filtering
 */
export function useProducts(params: {
  page?: number
  pageSize?: number
  category?: string
  status?: string
  searchText?: string
}) {
  return useQuery({
    queryKey: productKeys.list(params),
    queryFn: async () => {
      const response = await productServiceApi.productServiceListProducts(
        {
          page: params.page,
          pageSize: params.pageSize,
          category: params.category,
          status: params.status,
        },
        {
          // Pass searchText as additional query parameter via axios options
          params: {
            ...(params.searchText && { searchText: params.searchText }),
          },
        }
      )
      return response.data
    },
  })
}

/**
 * Hook to fetch all active products (for dropdowns/selects)
 */
export function useAllProducts() {
  return useQuery({
    queryKey: productKeys.list({ pageSize: 1000, status: 'active' }),
    queryFn: async () => {
      const response = await productServiceApi.productServiceListProducts({
        page: 1,
        pageSize: 1000,
        status: 'active',
      })
      return response.data
    },
    staleTime: 5 * 60 * 1000, // Cache for 5 minutes
  })
}

/**
 * Hook to fetch a single product by product key
 */
export function useProduct(productKey: string) {
  return useQuery({
    queryKey: productKeys.detail(productKey),
    queryFn: async () => {
      const response = await productServiceApi.productServiceGetProductByKey({
        productKey,
      })
      return response.data
    },
    enabled: !!productKey,
  })
}

// ============================================================================
// Mutations
// ============================================================================

/**
 * Hook to create a new product
 */
export function useCreateProduct() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (data: ProductV1CreateProductRequest) => {
      const response = await productServiceApi.productServiceCreateProduct({
        productV1CreateProductRequest: data,
      })
      return response.data
    },
    onSuccess: () => {
      // Invalidate all product lists to refetch
      queryClient.invalidateQueries({ queryKey: productKeys.lists() })
    },
  })
}

/**
 * Hook to update an existing product
 */
export function useUpdateProduct() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({
      productKey,
      data,
    }: {
      productKey: string
      data: ProductV1UpdateProductRequest
    }) => {
      const response = await productServiceApi.productServiceUpdateProduct({
        productKey,
        productV1UpdateProductRequest: data,
      })
      return response.data
    },
    onSuccess: (response) => {
      // Invalidate the specific product detail using productKey
      if (response.product?.productKey) {
        queryClient.invalidateQueries({
          queryKey: productKeys.detail(response.product.productKey),
        })
      }
      // Invalidate all product lists to refetch
      queryClient.invalidateQueries({ queryKey: productKeys.lists() })
    },
  })
}

/**
 * Hook to delete a product
 */
export function useDeleteProduct() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (id: string) => {
      const response = await productServiceApi.productServiceDeleteProduct({
        id,
      })
      return response.data
    },
    onSuccess: () => {
      // Invalidate all product lists to refetch
      queryClient.invalidateQueries({ queryKey: productKeys.lists() })
    },
  })
}
