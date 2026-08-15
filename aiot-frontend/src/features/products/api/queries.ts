import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type {
  ProductCreateProductRequest as ProductV1CreateProductRequest,
  ProductGetProductByKeyResponse as ProductV1GetProductByKeyResponse,
  ProductListProductsResponse as ProductV1ListProductsResponse,
  ProductUpdateProductRequest as ProductV1UpdateProductRequest,
  ProductUpdateProductResponse as ProductV1UpdateProductResponse,
} from '@/api/generated/model'
import {
  deleteProductsId,
  getProducts,
  getProductsKeyProductKey,
  postProducts,
  putProductsKeyProductKey,
} from '@/api/generated'

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
      return getProducts(params) as unknown as ProductV1ListProductsResponse
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
      return getProducts({
        page: 1,
        pageSize: 1000,
        status: 'active',
      }) as unknown as ProductV1ListProductsResponse
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
      return getProductsKeyProductKey(productKey) as unknown as ProductV1GetProductByKeyResponse
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
    mutationFn: (data: ProductV1CreateProductRequest) => postProducts(data as never),
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
      return putProductsKeyProductKey(productKey, data as never) as unknown as ProductV1UpdateProductResponse
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
    mutationFn: (id: string) => deleteProductsId(Number(id)),
    onSuccess: () => {
      // Invalidate all product lists to refetch
      queryClient.invalidateQueries({ queryKey: productKeys.lists() })
    },
  })
}
