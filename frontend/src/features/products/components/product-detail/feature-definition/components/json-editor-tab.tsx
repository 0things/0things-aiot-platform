import { AlertCircle, CheckCheck, Code, RefreshCw, XCircle } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'

interface JsonEditorTabProps {
  tslText: string
  onChange: (value: string) => void
  onFormat: () => void
  jsonValid: boolean | null
  validationError: string
}

export function JsonEditorTab({
  tslText,
  onChange,
  onFormat,
  jsonValid,
  validationError,
}: JsonEditorTabProps) {
  const { t } = useTranslation('deviceManagement')

  return (
    <div className='rounded-lg border border-border bg-background shadow-sm'>
      <div className='flex items-center justify-between border-b border-border bg-muted/50 px-4 py-3'>
        <div className='flex items-center gap-3'>
          <Code className='h-4 w-4 text-muted-foreground' />
          <span className='text-sm font-medium'>
            TSL {t('productDetail.featureDefinition.sections.properties')}
          </span>
          {/* JSON 验证状态指示器 */}
          {jsonValid === true && (
            <Badge
              variant='outline'
              className='gap-1.5 border-green-500/30 bg-green-500/10 text-green-700 dark:text-green-400'
            >
              <CheckCheck className='h-3 w-3' />
              {t('productDetail.featureDefinition.statusBadge.valid')}
            </Badge>
          )}
          {jsonValid === false && (
            <Badge
              variant='outline'
              className='gap-1.5 border-destructive/30 bg-destructive/10 text-destructive'
            >
              <XCircle className='h-3 w-3' />
              {t('productDetail.featureDefinition.statusBadge.invalid')}
            </Badge>
          )}
          {jsonValid === null && (
            <Badge variant='outline' className='gap-1.5 text-muted-foreground'>
              <AlertCircle className='h-3 w-3' />
              {t('productDetail.featureDefinition.statusBadge.unknown')}
            </Badge>
          )}
        </div>
        <div className='flex items-center gap-2'>
          <Button
            variant='outline'
            size='sm'
            onClick={onFormat}
            className='h-8 gap-1.5 text-xs'
            disabled={jsonValid === false}
          >
            <RefreshCw className='h-3.5 w-3.5' />
            {t('productDetail.featureDefinition.buttons.format')}
          </Button>
          <Badge variant='secondary' className='font-mono text-xs'>
            JSON
          </Badge>
        </div>
      </div>
      {/* 验证错误提示 */}
      {validationError && jsonValid === false && (
        <div className='border-b border-destructive/30 bg-destructive/5 px-4 py-2'>
          <div className='flex items-start gap-2'>
            <AlertCircle className='mt-0.5 h-3.5 w-3.5 flex-shrink-0 text-destructive' />
            <p className='text-xs text-destructive'>{validationError}</p>
          </div>
        </div>
      )}
      <div className='max-h-[600px] overflow-y-auto'>
        <Textarea
          value={tslText}
          onChange={(e) => onChange(e.target.value)}
          placeholder={`${t('productDetail.featureDefinition.title')}\n{\n  "schema": "schema.json",\n  "profile": {\n    "productKey": "PRODUCT_KEY"\n  },\n  "properties": [],\n  "events": [],\n  "services": []\n}`}
          className='min-h-[400px] resize-none border-0 bg-background font-mono text-sm leading-relaxed focus-visible:ring-0 focus-visible:ring-offset-0'
        />
      </div>
    </div>
  )
}
