import { createFileRoute } from '@tanstack/react-router'
import { ProductDetailPage } from '@/features/products/components/product-detail/product-detail-page'

export const Route = createFileRoute(
  '/_authenticated/device-management/products/$productKey/'
)({
  component: ProductDetailPage,
})
