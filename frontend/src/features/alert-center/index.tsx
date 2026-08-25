import { useState } from 'react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  useAlerts,
  useAckAlert,
  useResolveAlert,
  type AlertStatus,
} from './api/queries'

export function AlertCenter() {
  const [status, setStatus] = useState<AlertStatus | undefined>('open')
  const { data, isLoading } = useAlerts({ status, page: 1, pageSize: 50 })
  const ack = useAckAlert()
  const resolve = useResolveAlert()

  return (
    <div className='space-y-4 p-6'>
      <div className='flex items-center justify-between'>
        <h1 className='text-2xl font-semibold'>Alert Center</h1>
        <div className='flex gap-2'>
          {(['open', 'acknowledged', 'resolved'] as AlertStatus[]).map((s) => (
            <Button
              key={s}
              size='sm'
              variant={status === s ? 'default' : 'outline'}
              onClick={() => setStatus(s)}
            >
              {s}
            </Button>
          ))}
        </div>
      </div>

      {isLoading ? (
        <div className='text-muted-foreground'>Loading…</div>
      ) : !data || data.alerts.length === 0 ? (
        <div className='text-muted-foreground italic'>No alerts.</div>
      ) : (
        <div className='space-y-2'>
          {data.alerts.map((a) => (
            <Card key={a.id}>
              <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
                <div className='flex items-center gap-2'>
                  <Badge
                    variant={
                      a.severity === 'critical'
                        ? 'destructive'
                        : a.severity === 'warning'
                          ? 'default'
                          : 'secondary'
                    }
                  >
                    {a.severity}
                  </Badge>
                  <CardTitle className='text-base'>{a.summary}</CardTitle>
                  {a.count > 1 ? (
                    <Badge variant='outline'>×{a.count}</Badge>
                  ) : null}
                </div>
                <div className='flex gap-2'>
                  {a.status === 'open' ? (
                    <Button
                      size='sm'
                      variant='outline'
                      onClick={() => ack.mutate(a.id)}
                    >
                      Ack
                    </Button>
                  ) : null}
                  {a.status !== 'resolved' ? (
                    <Button size='sm' onClick={() => resolve.mutate(a.id)}>
                      Resolve
                    </Button>
                  ) : null}
                </div>
              </CardHeader>
              <CardContent className='text-xs text-muted-foreground'>
                <div>Device: {a.deviceKey}</div>
                <div>Raised: {new Date(a.raisedAt).toLocaleString()}</div>
                {a.lastRaisedAt !== a.raisedAt ? (
                  <div>Last: {new Date(a.lastRaisedAt).toLocaleString()}</div>
                ) : null}
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}
