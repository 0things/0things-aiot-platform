import { Edit2, Plus, Trash2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import type { Property } from '../types'

interface PropertyListProps {
  properties: Property[]
  onAdd: () => void
  onEdit: (property: Property, index: number) => void
  onDelete: (index: number) => void
}

export function PropertyList({
  properties,
  onAdd,
  onEdit,
  onDelete,
}: PropertyListProps) {
  const { t } = useTranslation('deviceManagement')

  return (
    <div className='rounded-lg border'>
      <div className='flex items-center justify-between border-b bg-muted/30 px-4 py-3'>
        <div className='flex items-center gap-2'>
          <h4 className='text-sm font-semibold'>
            {t('productDetail.featureDefinition.sections.properties')}
          </h4>
          <Badge variant='secondary'>{properties.length}</Badge>
        </div>
        <Button size='sm' onClick={onAdd}>
          <Plus className='mr-1 h-3 w-3' />
          {t('productDetail.featureDefinition.buttons.addProperty')}
        </Button>
      </div>
      <div className='p-4'>
        {properties.length === 0 ? (
          <div className='py-8 text-center text-sm text-muted-foreground'>
            {t('productDetail.featureDefinition.empty.properties')}
          </div>
        ) : (
          <div className='space-y-2'>
            {properties.map((property, index) => (
              <div
                key={index}
                className='flex items-center justify-between rounded-md border p-3 hover:bg-muted/50'
              >
                <div className='flex-1'>
                  <div className='flex items-center gap-2'>
                    <span className='font-mono text-sm font-medium'>
                      {property.identifier}
                    </span>
                    <Badge
                      variant='outline'
                      className='font-mono text-xs uppercase'
                    >
                      {property.accessMode}
                    </Badge>
                    <Badge
                      variant='outline'
                      className='font-mono text-xs text-primary'
                    >
                      {property.dataType.type}
                    </Badge>
                  </div>
                  <p className='mt-1 text-xs text-muted-foreground'>
                    {property.name}
                  </p>
                </div>
                <div className='flex gap-2'>
                  <Button
                    variant='ghost'
                    size='sm'
                    onClick={() => onEdit(property, index)}
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
