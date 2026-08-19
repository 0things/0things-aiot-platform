import { ProductsDialogs } from './components/products-dialogs'
import { ProductsPrimaryButtons } from './components/products-primary-buttons'
import { ProductsProvider } from './components/products-provider'
import { ProductsTable } from './components/products-table'

export function Products() {
  return (
    <ProductsProvider>
      <div className='flex flex-1 flex-col gap-4'>
        <div className='flex flex-wrap items-end justify-between gap-2'>
          <div>
            <h2 className='text-2xl font-bold tracking-tight'>
              Product Management
            </h2>
            <p className='text-muted-foreground'>
              Create and manage your product catalog.
            </p>
          </div>
          <ProductsPrimaryButtons />
        </div>
        <ProductsTable />
      </div>

      <ProductsDialogs />
    </ProductsProvider>
  )
}
