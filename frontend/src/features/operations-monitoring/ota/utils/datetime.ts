type OtaTimeValue = string | number | null | undefined

export function formatOtaDateTime(value?: OtaTimeValue): string {
  if (value === null || value === undefined || value === '') return '-'
  if (typeof value === 'number') return String(value)

  const match =
    /^(\d{4})-(\d{2})-(\d{2})[Tt ](\d{2}):(\d{2}):(\d{2})(?:\.\d+)?(Z|[+-]\d{2}:?\d{2})?$/.exec(
      value
    )
  if (!match) return value

  const [, year, month, day, hour, minute, second, rawOffset] = match
  let offset = rawOffset ?? ''
  if (offset && offset !== 'Z' && !offset.includes(':')) {
    offset = `${offset.slice(0, 3)}:${offset.slice(3)}`
  }
  return `${year}-${month}-${day} ${hour}:${minute}:${second}${
    offset ? ` ${offset}` : ''
  }`
}

export function formatOtaDate(value?: OtaTimeValue): string {
  if (value === null || value === undefined || value === '') return '-'
  if (typeof value === 'number') return String(value)

  const match = /^(\d{4})-(\d{2})-(\d{2})/.exec(value)
  return match ? `${match[1]}-${match[2]}-${match[3]}` : value
}
