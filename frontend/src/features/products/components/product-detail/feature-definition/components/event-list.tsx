import { Edit2, Plus, Trash2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import type { Event } from '../types'

interface EventListProps {
  events: Event[]
  onAdd: () => void
  onEdit: (event: Event, index: number) => void
  onDelete: (index: number) => void
}

export function EventList({ events, onAdd, onEdit, onDelete }: EventListProps) {
  const { t } = useTranslation('deviceManagement')

  return (
    <div className='rounded-lg border'>
      <div className='flex items-center justify-between border-b bg-muted/30 px-4 py-3'>
        <div className='flex items-center gap-2'>
          <h4 className='text-sm font-semibold'>
            {t('productDetail.featureDefinition.sections.events')}
          </h4>
          <Badge variant='secondary'>{events.length}</Badge>
        </div>
        <Button size='sm' onClick={onAdd}>
          <Plus className='mr-1 h-3 w-3' />
          {t('productDetail.featureDefinition.buttons.addEvent')}
        </Button>
      </div>
      <div className='p-4'>
        {events.length === 0 ? (
          <div className='py-8 text-center text-sm text-muted-foreground'>
            {t('productDetail.featureDefinition.empty.events')}
          </div>
        ) : (
          <div className='space-y-2'>
            {events.map((event, index) => (
              <div
                key={index}
                className='flex items-center justify-between rounded-md border p-3 hover:bg-muted/50'
              >
                <div className='flex-1'>
                  <div className='flex items-center gap-2'>
                    <span className='font-mono text-sm font-medium'>
                      {event.identifier}
                    </span>
                    <Badge
                      variant='outline'
                      className='font-mono text-xs uppercase'
                    >
                      {event.type}
                    </Badge>
                    <Badge variant='outline' className='text-xs'>
                      {t('productDetail.featureDefinition.statusBadge.output')}:{' '}
                      {event.outputData?.length ?? 0}
                    </Badge>
                  </div>
                  <p className='mt-1 text-xs text-muted-foreground'>
                    {event.name}
                  </p>
                </div>
                <div className='flex gap-2'>
                  <Button
                    variant='ghost'
                    size='sm'
                    onClick={() => onEdit(event, index)}
                  >
                    <Edit2 className='h-3 w-3' />
                  </Button>
                  <Button
                    variant='ghost'
                    size='sm'
                    onClick={() => onDelete(index)}
                  >
                    <Trash2 className='h-3 w-3' />
                  </Button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
