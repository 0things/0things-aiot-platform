import type { ComponentProps } from 'react'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'
import { useFormField } from '@/components/ui/form'

export function OTAFormMessage({ className, ...props }: ComponentProps<'p'>) {
  const { t } = useTranslation('ota')
  const { error, formMessageId } = useFormField()
  const body = error ? t(String(error.message ?? '')) : props.children

  if (!body) {
    return null
  }

  return (
    <p
      data-slot='form-message'
      id={formMessageId}
      className={cn('text-sm text-destructive', className)}
      {...props}
    >
      {body}
    </p>
  )
}
