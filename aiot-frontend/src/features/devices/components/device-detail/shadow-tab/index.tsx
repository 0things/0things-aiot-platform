import { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { toast } from 'sonner'
import { useDeviceShadow, useUpdateDesired } from '@/features/devices/api/shadow'
import { JsonPane } from './json-pane'
import { DesiredEditor } from './desired-editor'

type Props = {
  deviceKey: string
}

export function ShadowTab({ deviceKey }: Props) {
  const { data, isLoading, isError } = useDeviceShadow(deviceKey)
  const update = useUpdateDesired(deviceKey)
  const [editing, setEditing] = useState(false)

  if (isLoading) return <div className='p-4 text-muted-foreground'>Loading shadow…</div>
  if (isError || !data) return <div className='p-4 text-destructive'>Failed to load shadow.</div>

  const hasDelta = Object.keys(data.delta ?? {}).length > 0

  return (
    <div className='space-y-4'>
      <div className='flex items-center gap-2'>
        <Badge variant='outline'>v{data.version}</Badge>
        {hasDelta ? (
          <Badge variant='destructive'>Pending sync</Badge>
        ) : (
          <Badge variant='secondary'>In sync</Badge>
        )}
        <div className='ml-auto'>
          {!editing ? (
            <Button size='sm' onClick={() => setEditing(true)}>
              Edit desired
            </Button>
          ) : (
            <Button size='sm' variant='ghost' onClick={() => setEditing(false)}>
              Cancel
            </Button>
          )}
        </div>
      </div>

      <div className='grid grid-cols-1 gap-4 lg:grid-cols-3'>
        <Card>
          <CardHeader>
            <CardTitle>Desired</CardTitle>
          </CardHeader>
          <CardContent>
            {editing ? (
              <DesiredEditor
                initial={data.desired}
                onSave={async (next) => {
                  await update.mutateAsync({ desired: next, version: data.version })
                  toast.success('Desired state updated')
                  setEditing(false)
                }}
              />
            ) : (
              <JsonPane value={data.desired} />
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Reported</CardTitle>
          </CardHeader>
          <CardContent>
            <JsonPane value={data.reported} />
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Delta</CardTitle>
          </CardHeader>
          <CardContent>
            <JsonPane value={data.delta} emptyHint='No pending changes' />
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
