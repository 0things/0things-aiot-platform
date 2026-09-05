import { Copy } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'

interface CopyButtonProps {
  value: string
  className?: string
}

export function CopyButton({ value, className }: CopyButtonProps) {
  const { t } = useTranslation()
  const label = t('common:copySuccess')

  const handleCopy = async () => {
    await navigator.clipboard.writeText(value)
    toast.success(label)
  }

  return (
    <Button
      variant='ghost'
      size='icon'
      className={cn('h-6 w-6', className)}
      onClick={handleCopy}
      title={label}
      aria-label={label}
    >
      <Copy className='h-3.5 w-3.5 text-primary' />
    </Button>
  )
}
