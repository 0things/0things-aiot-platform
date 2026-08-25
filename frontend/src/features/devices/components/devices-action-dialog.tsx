'use client'

import React from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import { useAllProducts } from '@/features/products/api/queries'
import { useCreateDevice, useUpdateDevice } from '../api/queries'
import {
  type Device,
  deviceFormSchema,
  type DeviceFormData,
} from '../data/schema'

type DeviceActionDialogProps = {
  currentRow?: Device
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function DevicesActionDialog({
  currentRow,
  open,
  onOpenChange,
}: DeviceActionDialogProps) {
  const isEdit = !!currentRow
  const { t } = useTranslation('deviceManagement')

  // Fetch products list
  const { data: productsResponse, isLoading: isLoadingProducts } =
    useAllProducts()
  const products = React.useMemo(
    () => productsResponse?.products || [],
    [productsResponse]
  )

  // Create and update mutations
  const createDevice = useCreateDevice()
  const updateDevice = useUpdateDevice()

  const form = useForm<DeviceFormData>({
    resolver: zodResolver(deviceFormSchema),
    defaultValues: {
      name: '',
      productId: '',
      metadata: '',
    },
  })

  // Reset form when dialog opens or currentRow changes
  React.useEffect(() => {
    if (open) {
      if (isEdit && currentRow) {
        form.reset({
          name: currentRow.name,
          productId: currentRow.productId,
          metadata: currentRow.metadata || '',
        })
      } else {
        form.reset({
          name: '',
          productId: '',
          metadata: '',
        })
      }
    }
  }, [open, isEdit, currentRow, products, form])

  const onSubmit = async (values: DeviceFormData) => {
    try {
      if (isEdit && currentRow) {
        // Update existing device
        await updateDevice.mutateAsync({
          id: currentRow.id,
          data: {
            name: values.name,
            metadata: values.metadata || undefined,
          },
        })
        toast.success(
          t('dialog.edit.successMessage', 'Device updated successfully!')
        )
      } else {
        // Create new device
        await createDevice.mutateAsync({
          name: values.name,
          productId: Number(values.productId),
          metadata: values.metadata || undefined,
        })
        toast.success(
          t('dialog.create.successMessage', 'Device created successfully!')
        )
      }
      onOpenChange(false)
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Failed to submit device form:', error)

      // Extract error message from API response
      // API error format: { code, reason, message, metadata }
      const apiError = error as {
        response?: {
          data?: {
            message?: string
            error?: string
          }
        }
        message?: string
      }

      const errorMessage =
        apiError?.response?.data?.message ||
        apiError?.response?.data?.error ||
        apiError?.message ||
        (isEdit
          ? t('dialog.edit.errorMessage', 'Failed to update device')
          : t('dialog.create.errorMessage', 'Failed to create device'))

      toast.error(errorMessage)
    }
  }

  const isSubmitting = createDevice.isPending || updateDevice.isPending

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-lg'>
        <DialogHeader className='text-start'>
          <DialogTitle>
            {isEdit ? t('dialog.edit.title') : t('dialog.create.title')}
          </DialogTitle>
          <DialogDescription>
            {isEdit
              ? t('dialog.edit.description')
              : t('dialog.create.description')}
          </DialogDescription>
        </DialogHeader>
        <div className='h-105 w-[calc(100%+0.75rem)] overflow-y-auto py-1 pe-3'>
          <Form {...form}>
            <form
              id='device-form'
              onSubmit={form.handleSubmit(onSubmit)}
              className='space-y-4 px-0.5'
            >
              <FormField
                control={form.control}
                name='name'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 text-end'>
                      {t('dialog.form.name')}
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t('dialog.form.placeholder.name')}
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
                name='productId'
                render={({ field }) => (
                  <FormItem className='grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1'>
                    <FormLabel className='col-span-2 text-end'>
                      {t('dialog.form.product')}
                    </FormLabel>
                    {isEdit ? (
                      <FormControl>
                        <Input
                          value={currentRow?.productName || ''}
                          disabled
                          className='col-span-4'
                        />
                      </FormControl>
                    ) : (
                      <Select
                        onValueChange={field.onChange}
                        value={field.value}
                        disabled={isLoadingProducts}
                      >
                        <FormControl>
                          <SelectTrigger className='col-span-4'>
                            <SelectValue
                              placeholder={t('dialog.form.placeholder.product')}
                            />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {products.map((product) => (
                            <SelectItem
                              key={product.id}
                              value={
                                product.id != null ? String(product.id) : ''
                              }
                            >
                              {product.name}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    )}
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
                      {t('dialog.form.metadata')}
                    </FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder={t('dialog.form.placeholder.metadata')}
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
          <Button type='submit' form='device-form' disabled={isSubmitting}>
            {isSubmitting && (
              <span className='mr-2 inline-block h-4 w-4 animate-spin rounded-full border-2 border-solid border-current border-r-transparent' />
            )}
            {t('dialog.form.saveButton')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
