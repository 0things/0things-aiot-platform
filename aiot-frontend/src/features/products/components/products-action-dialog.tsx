'use client'

import React from 'react'
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { zodResolver } from '@hookform/resolvers/zod'
import { toast } from 'sonner'
import type {
  ProductCreateProductRequest as ProductV1CreateProductRequest,
  ProductUpdateProductRequest as ProductV1UpdateProductRequest,
} from '@/api/generated/model'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { SelectDropdown } from '@/components/select-dropdown'
import { useCreateProduct, useUpdateProduct } from '../api/queries'
import {
  categories,
  connectivityMethods,
  directGatewayProtocols,
  gatewaySubProtocols,
  nodeTypes,
  statuses,
} from '../data/data'
import {
  type Product,
  productFormSchema,
  type ProductFormData,
} from '../data/schema'

type ProductActionDialogProps = {
  currentRow?: Product
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function ProductsActionDialog({
  currentRow,
  open,
  onOpenChange,
}: ProductActionDialogProps) {
  const { t } = useTranslation('deviceManagement')
  const isEdit = !!currentRow

  // Create and update mutations
  const createProduct = useCreateProduct()
  const updateProduct = useUpdateProduct()

  const form = useForm<ProductFormData>({
    resolver: zodResolver(productFormSchema),
    defaultValues: isEdit
      ? {
          name: currentRow.name,
          description: currentRow.description || '',
          category: currentRow.category || '',
          status: currentRow.status,
          nodeType:
            (currentRow as Product & { nodeType?: string }).nodeType ||
            'direct',
          connectivityMethod: (
            currentRow as Product & { connectivityMethod?: string }
          ).connectivityMethod,
          accessProtocol:
            (currentRow as Product & { accessProtocol?: string })
              .accessProtocol || 'http',
          metadata: currentRow.metadata || '',
        }
      : {
          name: '',
          description: '',
          category: '',
          status: 'active',
          nodeType: 'direct',
          connectivityMethod: undefined,
          accessProtocol: 'http',
          metadata: '',
        },
  })

  // Watch nodeType to control connectivity method field visibility and protocol options
  const nodeType = form.watch('nodeType')
  const showConnectivityMethod = nodeType === 'direct' || nodeType === 'gateway'

  // Determine which protocols to show based on node type
  const accessProtocolOptions =
    nodeType === 'gateway-sub' ? gatewaySubProtocols : directGatewayProtocols

  // Reset access protocol when node type changes to ensure valid protocol for the type
  const prevNodeType = React.useRef(nodeType)
  React.useEffect(() => {
    if (prevNodeType.current !== nodeType) {
      // Reset to first valid protocol for the new node type
      const defaultProtocol = nodeType === 'gateway-sub' ? 'custom' : 'http'
      form.setValue('accessProtocol', defaultProtocol)
      prevNodeType.current = nodeType
    }
  }, [nodeType, form])

  const onSubmit = async (values: ProductFormData) => {
    try {
      if (isEdit && currentRow) {
        // Update existing product
        await updateProduct.mutateAsync({
          productKey: currentRow.productKey,
          data: {
            name: values.name,
            description: values.description || undefined,
            category: values.category,
            status: values.status,
            metadata: values.metadata || undefined,
            ...(values.nodeType && { nodeType: values.nodeType }),
            ...(values.connectivityMethod && {
              connectivityMethod: values.connectivityMethod,
            }),
            ...(values.accessProtocol && {
              accessProtocol: values.accessProtocol,
            }),
          } as unknown as ProductV1UpdateProductRequest,
        })
        toast.success('Product updated successfully!')
      } else {
        // Create new product
        await createProduct.mutateAsync({
          name: values.name,
          description: values.description || undefined,
          category: values.category,
          status: values.status,
          metadata: values.metadata || undefined,
          ...(values.nodeType && { nodeType: values.nodeType }),
          ...(values.connectivityMethod && {
            connectivityMethod: values.connectivityMethod,
          }),
          ...(values.accessProtocol && {
            accessProtocol: values.accessProtocol,
          }),
        } as unknown as ProductV1CreateProductRequest)
        toast.success('Product created successfully!')
      }
      form.reset()
      onOpenChange(false)
    } catch (error) {
      const apiError = error as {
        response?: { data?: { message?: string; error?: string } }
        message?: string
      }
      const errorMessage =
        apiError?.response?.data?.message ||
        apiError?.response?.data?.error ||
        apiError?.message ||
        (isEdit ? 'Failed to update product' : 'Failed to create product')
      toast.error(errorMessage)
    }
  }

  const isSubmitting = createProduct.isPending || updateProduct.isPending

  return (
    <Dialog
      open={open}
      onOpenChange={(state) => {
        form.reset()
        onOpenChange(state)
      }}
    >
      <DialogContent className='sm:max-w-lg'>
        <DialogHeader className='text-start'>
          <DialogTitle>
            {isEdit ? 'Edit Product' : 'Create New Product'}
          </DialogTitle>
          <DialogDescription>
            {isEdit
              ? 'Update the product details here. '
              : 'Add a new product to your inventory. '}
            Click save when you&apos;re done.
          </DialogDescription>
        </DialogHeader>
        <div className='h-105 w-[calc(100%+0.75rem)] overflow-y-auto py-1 pe-3'>
          <Form {...form}>
            <form
              id='product-form'
              onSubmit={form.handleSubmit(onSubmit)}
              className='space-y-4 px-0.5'
            >
              <FormField
                control={form.control}
                name='name'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 text-end'>Name</FormLabel>
                    <FormControl>
                      <Input
                        placeholder='Smart Sensor X1'
                        className='col-span-4'
                        autoComplete='off'
                        {...field}
                      />
                    </FormControl>
                    <FormMessage className='col-span-4 col-start-3' />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='nodeType'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 text-end'>
                      {t('productDetail.info.fields.nodeType')}
                    </FormLabel>
                    <SelectDropdown
                      defaultValue={field.value}
                      onValueChange={field.onChange}
                      placeholder='Select node type'
                      className='col-span-4'
                      items={nodeTypes.map(({ value }) => ({
                        label: t(`productDetail.nodeTypes.${value}`),
                        value,
                      }))}
                    />
                    <FormMessage className='col-span-4 col-start-3' />
                  </FormItem>
                )}
              />
              {showConnectivityMethod && (
                <FormField
                  control={form.control}
                  name='connectivityMethod'
                  render={({ field }) => (
                    <FormItem className='grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1'>
                      <FormLabel className='col-span-2 text-end'>
                        {t('productDetail.info.fields.connectivityMethod')}
                      </FormLabel>
                      <SelectDropdown
                        defaultValue={field.value}
                        onValueChange={field.onChange}
                        placeholder='Select connectivity method'
                        className='col-span-4'
                        items={connectivityMethods.map(({ value }) => ({
                          label: t(`productDetail.connectivityMethods.${value}`),
                          value,
                        }))}
                      />
                      <FormMessage className='col-span-4 col-start-3' />
                    </FormItem>
                  )}
                />
              )}
              <FormField
                control={form.control}
                name='accessProtocol'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 text-end'>
                      {nodeType === 'gateway-sub'
                        ? t('productDetail.info.fields.gatewayAccessProtocol')
                        : t('productDetail.info.fields.accessProtocol')}
                    </FormLabel>
                    <SelectDropdown
                      defaultValue={field.value}
                      onValueChange={field.onChange}
                      placeholder='Select access protocol'
                      className='col-span-4'
                      items={accessProtocolOptions.map(({ value }) => ({
                        label: t(`productDetail.accessProtocols.${value}`),
                        value,
                      }))}
                    />
                    <FormMessage className='col-span-4 col-start-3' />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='category'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 text-end'>
                      Category
                    </FormLabel>
                    <SelectDropdown
                      defaultValue={field.value}
                      onValueChange={field.onChange}
                      placeholder='Select a category'
                      className='col-span-4'
                      items={categories.map(({ label, value }) => ({
                        label,
                        value,
                      }))}
                    />
                    <FormMessage className='col-span-4 col-start-3' />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='status'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 text-end'>
                      Status
                    </FormLabel>
                    <SelectDropdown
                      defaultValue={field.value}
                      onValueChange={field.onChange}
                      placeholder='Select status'
                      className='col-span-4'
                      items={statuses.map(({ label, value }) => ({
                        label,
                        value,
                      }))}
                    />
                    <FormMessage className='col-span-4 col-start-3' />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='description'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-6 items-start space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 pt-2 text-end'>
                      Description
                    </FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder='Product description...'
                        className='col-span-4 resize-none'
                        rows={3}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage className='col-span-4 col-start-3' />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='metadata'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-6 items-start space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 pt-2 text-end'>
                      Metadata
                    </FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder='{"key": "value"}'
                        className='col-span-4 resize-none font-mono text-sm'
                        rows={3}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage className='col-span-4 col-start-3' />
                  </FormItem>
                )}
              />
            </form>
          </Form>
        </div>
        <DialogFooter>
          <Button type='submit' form='product-form' disabled={isSubmitting}>
            {isSubmitting && (
              <span className='mr-2 inline-block h-4 w-4 animate-spin rounded-full border-2 border-solid border-current border-r-transparent' />
            )}
            Save changes
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
