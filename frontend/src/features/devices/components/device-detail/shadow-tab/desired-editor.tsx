import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'

type Props = {
  initial: Record<string, unknown>
  onSave: (next: Record<string, unknown>) => void | Promise<void>
}

export function DesiredEditor({ initial, onSave }: Props) {
  const [text, setText] = useState(() => JSON.stringify(initial ?? {}, null, 2))
  const [error, setError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)

  const handleSave = async () => {
    setError(null)
    let parsed: Record<string, unknown>
    try {
      parsed = JSON.parse(text)
      if (
        typeof parsed !== 'object' ||
        parsed === null ||
        Array.isArray(parsed)
      ) {
        setError('Desired state must be a JSON object')
        return
      }
    } catch (e) {
      setError(`Invalid JSON: ${(e as Error).message}`)
      return
    }
    setSaving(true)
    try {
      await onSave(parsed)
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className='space-y-2'>
      <Textarea
        className='h-64 font-mono text-xs'
        value={text}
        onChange={(e) => setText(e.target.value)}
      />
      {error ? <div className='text-xs text-destructive'>{error}</div> : null}
      <div className='flex justify-end'>
        <Button size='sm' onClick={handleSave} disabled={saving}>
          {saving ? 'Saving…' : 'Save'}
        </Button>
      </div>
    </div>
  )
}
