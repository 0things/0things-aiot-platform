import { useState, useRef, useEffect } from 'react'
import { Route } from '@/routes/_authenticated/device-management/devices/$deviceKey'
import { Send, Trash2, Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useSimulatePush, usePushRecords } from '../../api'

interface LogEntry {
  id: string
  timestamp: Date
  request: string
  response?: string
  error?: boolean
  duration?: number
}

export function OnlineDebugTab() {
  const { deviceKey } = Route.useParams()
  const [jsonPayload, setJsonPayload] = useState('{}')
  const [logs, setLogs] = useState<LogEntry[]>([])
  const logsEndRef = useRef<HTMLDivElement>(null)

  const { mutate: simulatePush, isPending: isSimulating } = useSimulatePush()
  const { refetch: refetchRecords } = usePushRecords(deviceKey || '', {
    page: 1,
    pageSize: 10,
  })

  useEffect(() => {
    logsEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [logs])

  const handleExecute = async () => {
    try {
      JSON.parse(jsonPayload)
    } catch (e) {
      toast.error(`Invalid JSON: ${(e as Error).message}`)
      return
    }

    if (!deviceKey) return

    const startTime = Date.now()

    simulatePush(
      {
        deviceKey,
        payload: jsonPayload,
      },
      {
        onSuccess: (resp) => {
          const inner = resp.data
          setLogs((prev) => [
            ...prev,
            {
              id: `log-${Date.now()}`,
              timestamp: new Date(),
              request: jsonPayload,
              response: inner?.message || JSON.stringify(resp, null, 2),
              error: inner?.status !== 'success',
              duration: Date.now() - startTime,
            },
          ])
          if (inner?.status === 'success') {
            toast.success('Push simulation succeeded')
            refetchRecords()
          } else {
            toast.error(`Push simulation failed: ${inner?.message}`)
          }
        },
        onError: (error) => {
          setLogs((prev) => [
            ...prev,
            {
              id: `log-${Date.now()}`,
              timestamp: new Date(),
              request: jsonPayload,
              response: (error as Error).message,
              error: true,
              duration: Date.now() - startTime,
            },
          ])
          toast.error(`Error: ${(error as Error).message}`)
        },
      }
    )
  }

  return (
    <div className='space-y-6'>
      {/* Input */}
      <Card>
        <CardHeader>
          <CardTitle className='text-lg'>Command Console</CardTitle>
        </CardHeader>
        <CardContent className='space-y-4'>
          <div>
            <label className='mb-2 block text-xs font-semibold tracking-wide text-muted-foreground uppercase'>
              JSON Payload
            </label>
            <textarea
              value={jsonPayload}
              onChange={(e) => setJsonPayload(e.target.value)}
              className='h-40 w-full resize-none rounded-lg border bg-muted p-3 font-mono text-sm'
              placeholder='Enter JSON payload...'
              spellCheck='false'
            />
          </div>

          <Button
            onClick={handleExecute}
            disabled={isSimulating}
            className='w-full'
          >
            {isSimulating ? (
              <>
                <Loader2 className='mr-2 h-4 w-4 animate-spin' />
                Executing...
              </>
            ) : (
              <>
                <Send className='mr-2 h-4 w-4' />
                Execute
              </>
            )}
          </Button>
        </CardContent>
      </Card>

      {/* Result Log */}
      <Card>
        <CardHeader className='flex flex-row items-center justify-between'>
          <CardTitle className='text-lg'>Result Log</CardTitle>
          <Button
            variant='ghost'
            size='sm'
            onClick={() => setLogs([])}
            disabled={logs.length === 0}
          >
            <Trash2 className='mr-1 h-4 w-4' />
            Clear
          </Button>
        </CardHeader>
        <CardContent>
          <div className='h-96 space-y-2 overflow-y-auto rounded-lg bg-muted p-4 font-mono text-xs'>
            {logs.length === 0 ? (
              <div className='flex h-full items-center justify-center text-gray-500'>
                No logs yet. Execute a command to see results here.
              </div>
            ) : (
              <>
                {logs.map((log) => (
                  <div
                    key={log.id}
                    className='space-y-2 rounded border bg-white p-3'
                  >
                    <div className='flex items-center justify-between text-gray-500'>
                      <span>{log.timestamp.toLocaleTimeString()}</span>
                      {log.duration != null && <span>{log.duration}ms</span>}
                    </div>
                    <div>
                      <div className='mb-1 text-gray-500'>Request:</div>
                      <pre className='overflow-x-auto rounded bg-gray-50 p-2'>
                        {log.request}
                      </pre>
                    </div>
                    {log.response && (
                      <div>
                        <div className='mb-1 text-gray-500'>Response:</div>
                        <pre
                          className={`overflow-x-auto rounded p-2 ${log.error ? 'bg-red-50 text-red-700' : 'bg-gray-50'}`}
                        >
                          {log.response}
                        </pre>
                      </div>
                    )}
                  </div>
                ))}
                <div ref={logsEndRef} />
              </>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
