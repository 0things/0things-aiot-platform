import { useState } from 'react'
import { X, Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { toast } from 'sonner'
import { useDeviceTags, useAddTags, useRemoveTags } from '@/features/devices/api/tags'

type Props = { deviceKey: string }

export function TagsEditor({ deviceKey }: Props) {
  const { data: tags = [], isLoading } = useDeviceTags(deviceKey)
  const add = useAddTags(deviceKey)
  const remove = useRemoveTags(deviceKey)
  const [newKey, setNewKey] = useState('')
  const [newValue, setNewValue] = useState('')

  const handleAdd = async () => {
    if (!newKey.trim()) return
    try {
      await add.mutateAsync({ [newKey.trim()]: newValue })
      setNewKey('')
      setNewValue('')
      toast.success('Tag added')
    } catch (e) {
      toast.error((e as Error).message)
    }
  }

  return (
    <div className='space-y-3'>
      <div className='flex flex-wrap gap-2'>
        {isLoading ? (
          <span className='text-sm text-muted-foreground'>Loading…</span>
        ) : tags.length === 0 ? (
          <span className='text-sm text-muted-foreground italic'>No tags</span>
        ) : (
          tags.map((t) => (
            <Badge key={t.key} variant='secondary' className='gap-1 pr-1'>
              <span>
                {t.key}={t.value}
              </span>
              <button
                aria-label={`Remove ${t.key}`}
                onClick={() => remove.mutate([t.key])}
                className='ml-1 rounded hover:bg-destructive/20'
              >
                <X className='size-3' />
              </button>
            </Badge>
          ))
        )}
      </div>
      <div className='flex gap-2'>
        <Input
          placeholder='key'
          value={newKey}
          onChange={(e) => setNewKey(e.target.value)}
          className='max-w-[180px]'
        />
        <Input
          placeholder='value'
          value={newValue}
          onChange={(e) => setNewValue(e.target.value)}
          className='max-w-[240px]'
        />
        <Button size='sm' onClick={handleAdd} disabled={!newKey.trim() || add.isPending}>
          <Plus className='size-4' /> Add
        </Button>
      </div>
    </div>
  )
}
