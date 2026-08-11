type Props = {
  value: unknown
  emptyHint?: string
}

export function JsonPane({ value, emptyHint = 'Empty' }: Props) {
  const isEmpty =
    value == null ||
    (typeof value === 'object' && value !== null && Object.keys(value as object).length === 0)
  if (isEmpty) {
    return <div className='text-sm text-muted-foreground italic'>{emptyHint}</div>
  }
  return (
    <pre className='max-h-80 overflow-auto rounded border bg-muted/40 p-3 text-xs'>
      {JSON.stringify(value, null, 2)}
    </pre>
  )
}
