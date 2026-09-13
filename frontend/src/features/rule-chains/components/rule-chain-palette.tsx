import { ChevronRight, PanelLeftClose, PanelLeftOpen } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible'
import { Input } from '@/components/ui/input'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import type { RuleNodeDefinition } from '../types/rule-chain-types'

type RuleChainPaletteProps = {
  open: boolean
  search: string
  groups: Record<string, RuleNodeDefinition[]>
  onSearchChange: (value: string) => void
  onOpenChange: (open: boolean) => void
  onAddNode: (definition: RuleNodeDefinition) => void
}

export function RuleChainPalette({
  open,
  search,
  groups,
  onSearchChange,
  onOpenChange,
  onAddNode,
}: RuleChainPaletteProps) {
  const { t } = useTranslation('ruleChain')

  if (!open) {
    return (
      <div className='flex w-11 shrink-0 items-start justify-center border-r border-border bg-card/70 pt-3'>
        <Button
          variant='ghost'
          size='icon'
          onClick={() => onOpenChange(true)}
          aria-label={t('actions.expandPalette')}
        >
          <PanelLeftOpen />
        </Button>
      </div>
    )
  }

  return (
    <aside className='flex w-72 shrink-0 flex-col border-r border-border bg-card/70'>
      <div className='flex items-start justify-between gap-2 border-b border-border p-4'>
        <div>
          <h2 className='font-semibold text-card-foreground'>
            {t('palette.title')}
          </h2>
          <p className='mt-1 text-xs text-muted-foreground'>
            {t('palette.hint')}
          </p>
        </div>
        <Button
          variant='ghost'
          size='icon'
          onClick={() => onOpenChange(false)}
          aria-label={t('actions.collapsePalette')}
        >
          <PanelLeftClose />
        </Button>
      </div>
      <div className='p-3'>
        <Input
          value={search}
          onChange={(event) => onSearchChange(event.target.value)}
          placeholder={t('palette.searchPlaceholder')}
          aria-label={t('palette.searchPlaceholder')}
        />
      </div>
      <Separator />
      <ScrollArea className='min-h-0 flex-1'>
        <div className='flex flex-col gap-3 p-3'>
          {Object.keys(groups).length ? (
            Object.entries(groups).map(([category, items]) => (
              <Collapsible
                key={category}
                defaultOpen
                className='flex flex-col gap-2'
              >
                <CollapsibleTrigger asChild>
                  <Button
                    variant='ghost'
                    className='w-full justify-between px-2 text-sm font-medium text-card-foreground data-[state=open]:[&>svg]:rotate-90'
                  >
                    <span className='flex items-center gap-2'>
                      <ChevronRight data-icon='inline-start' />
                      {category}
                    </span>
                    <Badge variant='secondary'>{items.length}</Badge>
                  </Button>
                </CollapsibleTrigger>
                <CollapsibleContent className='flex flex-col gap-1'>
                  {items.map((definition) => {
                    const Icon = definition.icon
                    return (
                      <Button
                        key={definition.key}
                        variant='ghost'
                        draggable
                        onDragStart={(event) => {
                          event.dataTransfer.setData(
                            'application/rule-node',
                            definition.key
                          )
                          event.dataTransfer.effectAllowed = 'move'
                        }}
                        onClick={() => onAddNode(definition)}
                        className='h-auto cursor-grab items-start justify-start gap-2 px-2 py-2 text-left whitespace-normal text-card-foreground active:cursor-grabbing'
                      >
                        <Icon className='mt-0.5 size-4 shrink-0 text-amber-300' />
                        <span className='flex min-w-0 flex-col gap-0.5'>
                          <span className='text-sm font-medium'>
                            {definition.name}
                          </span>
                          <span className='line-clamp-2 text-xs text-muted-foreground'>
                            {definition.description}
                          </span>
                        </span>
                      </Button>
                    )
                  })}
                </CollapsibleContent>
              </Collapsible>
            ))
          ) : (
            <p className='px-2 py-6 text-center text-sm text-muted-foreground'>
              {t('palette.empty')}
            </p>
          )}
        </div>
      </ScrollArea>
    </aside>
  )
}
