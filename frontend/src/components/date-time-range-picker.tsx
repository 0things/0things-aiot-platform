import { format, setHours, setMinutes, setSeconds } from 'date-fns'
import { CalendarDays } from 'lucide-react'
import type { DateRange } from 'react-day-picker'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Calendar } from '@/components/ui/calendar'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'

export type DateTimeRangeValue = {
  startAt?: Date
  endAt?: Date
}

type DateTimeRangePickerProps = {
  value: DateTimeRangeValue
  onChange: (value: DateTimeRangeValue) => void
  placeholder: string
  startAtLabel: string
  endAtLabel: string
  timePrecisionLabel: string
  className?: string
  disabled?: boolean
}

export function DateTimeRangePicker({
  value,
  onChange,
  placeholder,
  startAtLabel,
  endAtLabel,
  timePrecisionLabel,
  className,
  disabled,
}: DateTimeRangePickerProps) {
  const selectDates = (next?: DateRange) => {
    onChange({
      startAt: next?.from ? withTime(next.from, value.startAt) : undefined,
      endAt: next?.to ? withTime(next.to, value.endAt) : undefined,
    })
  }
  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          variant='outline'
          className={cn('min-w-80 justify-start font-normal', className)}
          disabled={disabled}
        >
          <CalendarDays data-icon='inline-start' />
          {value.startAt && value.endAt
            ? `${format(value.startAt, 'yyyy-MM-dd HH:mm:ss')} → ${format(value.endAt, 'yyyy-MM-dd HH:mm:ss')}`
            : placeholder}
        </Button>
      </PopoverTrigger>
      <PopoverContent className='w-auto p-0' align='start'>
        <Calendar
          className='p-3'
          mode='range'
          selected={{ from: value.startAt, to: value.endAt }}
          onSelect={selectDates}
          numberOfMonths={2}
        />
        <Separator />
        <div className='grid grid-cols-2 gap-3 p-3'>
          <TimeSelect
            label={startAtLabel}
            value={value.startAt}
            onChange={(startAt) => onChange({ ...value, startAt })}
          />
          <TimeSelect
            label={endAtLabel}
            value={value.endAt}
            onChange={(endAt) => onChange({ ...value, endAt })}
          />
        </div>
        <p className='sr-only'>{timePrecisionLabel}</p>
      </PopoverContent>
    </Popover>
  )
}

function withTime(date: Date, source?: Date) {
  const time = source || new Date()
  return setSeconds(
    setMinutes(setHours(date, time.getHours()), time.getMinutes()),
    time.getSeconds()
  )
}

function TimeSelect({
  label,
  value,
  onChange,
}: {
  label: string
  value?: Date
  onChange: (value: Date) => void
}) {
  const setPart = (part: 'hours' | 'minutes' | 'seconds', next: string) => {
    const date = value || new Date()
    const number = Number(next)
    onChange(
      part === 'hours'
        ? setHours(date, number)
        : part === 'minutes'
          ? setMinutes(date, number)
          : setSeconds(date, number)
    )
  }
  return (
    <div className='flex flex-col gap-1.5'>
      <p className='text-xs font-medium text-muted-foreground'>{label}</p>
      <div className='grid grid-cols-3 gap-1'>
        {(['hours', 'minutes', 'seconds'] as const).map((part) => (
          <Select
            key={part}
            value={String(getTimePart(value, part)).padStart(2, '0')}
            onValueChange={(next) => setPart(part, next)}
          >
            <SelectTrigger size='sm' className='w-full px-2'>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                {Array.from(
                  { length: part === 'hours' ? 24 : 60 },
                  (_, index) => (
                    <SelectItem
                      key={index}
                      value={String(index).padStart(2, '0')}
                    >
                      {String(index).padStart(2, '0')}
                    </SelectItem>
                  )
                )}
              </SelectGroup>
            </SelectContent>
          </Select>
        ))}
      </div>
    </div>
  )
}

function getTimePart(
  value: Date | undefined,
  part: 'hours' | 'minutes' | 'seconds'
) {
  if (!value) return 0
  if (part === 'hours') return value.getHours()
  if (part === 'minutes') return value.getMinutes()
  return value.getSeconds()
}
