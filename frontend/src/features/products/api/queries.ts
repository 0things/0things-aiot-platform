import { useMutation, useQueryClient } from '@tanstack/react-query'
import {
  deleteProductsProductKey,
  getGetProductsProductKeyQueryKey,
  getGetProductsQueryKey,
  postProducts,
  putProductsProductKey,
  useGetProducts,
  useGetProductsProductKey,
} from '@/api/generated'
import type {
  GetProductsParams,
  ProductCreateProductRequest as ProductV1CreateProductRequest,
  ProductListProductsResponse as ProductV1ListProductsResponse,
  ProductUpdateProductRequest as ProductV1UpdateProductRequest,
  ProductUpdateProductResponse as ProductV1UpdateProductResponse,
} from '@/api/generated/model'

export type { GetProductsParams }
export type ProductListResponse = ProductV1ListProductsResponse

// ============================================================================
// Query Keys
// ============================================================================

export const productKeys = {
  all: ['products'] as const,
  lists: () => [...productKeys.all, 'list'] as const,
  list: (params?: GetProductsParams) => getGetProductsQueryKey(params),
  details: () => [...productKeys.all, 'detail'] as const,
  detail: (productKey: string) => getGetProductsProductKeyQueryKey(productKey),
}

// ============================================================================
// Queries
// ============================================================================

/**
 * Hook to fetch products list with pagination and filtering
 */
export function useProducts(params?: GetProductsParams) {
  return useGetProducts(params, {
    query: {
      select: (res) => res?.data,
    },
  })
}

/**
 * Hook to fetch all active products (for dropdowns/selects)
 */
export function useAllProducts() {
  return useGetProducts(
    { page: 1, pageSize: 1000, status: 'active' },
    {
      query: {
        select: (res) => res?.data,
        staleTime: 5 * 60 * 1000, // Cache for 5 minutes
      },
    }
  )
}

/**
 * Hook to fetch a single product by product key
 */
export function useProduct(productKey: string) {
  return useGetProductsProductKey(productKey, {
    query: {
      select: (res) => res?.data,
      enabled: !!productKey,
    },
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
    mutationFn: (data: ProductV1CreateProductRequest) =>
      postProducts(data as never),
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
      const res = await putProductsProductKey(productKey, data as never)
      return res?.data as unknown as ProductV1UpdateProductResponse
    },
    onSuccess: (response) => {
      if (response.product?.productKey) {
        queryClient.invalidateQueries({
          queryKey: productKeys.detail(response.product.productKey),
        })
      }
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
    mutationFn: (productKey: string) => deleteProductsProductKey(productKey),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: productKeys.lists() })
    },
  })
}
