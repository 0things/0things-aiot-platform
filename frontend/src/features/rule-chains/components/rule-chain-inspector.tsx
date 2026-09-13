import { Settings2, Trash2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import type { RuleFlowNode } from '../types/rule-chain-types'

type RuleChainInspectorProps = {
  open: boolean
  selected?: RuleFlowNode
  onOpenChange: (open: boolean) => void
  onDeleteNode: (nodeId: string) => void
  onConfigChange: (key: string, value: unknown) => void
}

export function RuleChainInspector({
  open,
  selected,
  onOpenChange,
  onDeleteNode,
  onConfigChange,
}: RuleChainInspectorProps) {
  const { t } = useTranslation('ruleChain')

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className='flex w-full flex-col gap-0 sm:max-w-sm'>
        <SheetHeader className='border-b px-6 py-5 text-left'>
          <SheetTitle className='flex items-center gap-2'>
            <Settings2 className='size-4' />
            {t('inspector.title')}
          </SheetTitle>
          <SheetDescription>
            {selected
              ? selected.data.definition.description
              : t('inspector.selectHint')}
          </SheetDescription>
        </SheetHeader>
        {selected ? (
          <ScrollArea className='min-h-0 flex-1'>
            <div className='space-y-5 p-6'>
              <div className='flex items-center justify-between'>
                <p className='text-sm font-semibold'>
                  {selected.data.definition.name}
                </p>
                {selected.id !== 'entry' && (
                  <Button
                    variant='ghost'
                    size='icon'
                    className='text-destructive'
                    onClick={() => onDeleteNode(selected.id)}
                    aria-label={t('inspector.deleteNode')}
                  >
                    <Trash2 className='size-4' />
                  </Button>
                )}
              </div>
              <Separator />
              {selected.data.definition.configFields.length ? (
                <div className='space-y-4'>
                  {selected.data.definition.configFields.map((field) => (
                    <div key={field.key} className='space-y-2'>
                      <Label htmlFor={'field-' + field.key}>
                        {field.title}
                      </Label>
                      <Input
                        id={'field-' + field.key}
                        type={field.type === 'number' ? 'number' : 'text'}
                        value={String(selected.data.config[field.key] ?? '')}
                        onChange={(event) =>
                          onConfigChange(
                            field.key,
                            field.type === 'number'
                              ? Number(event.target.value)
                              : event.target.value
                          )
                        }
                      />
                    </div>
                  ))}
                </div>
              ) : (
                <p className='text-sm text-muted-foreground'>
                  {t('inspector.noConfig')}
                </p>
              )}
              <div className='rounded-lg bg-muted/50 p-3 text-xs text-muted-foreground'>
                {t('inspector.tip')}
              </div>
            </div>
          </ScrollArea>
        ) : (
          <p className='p-6 text-sm text-muted-foreground'>
            {t('inspector.selectHint')}
          </p>
        )}
      </SheetContent>
    </Sheet>
  )
}
