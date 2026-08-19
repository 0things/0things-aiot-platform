import * as React from 'react'
import { Check, ChevronsUpDown, Loader2 } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'

type SearchableSelectOption = {
  label: string
  value: string
}

type SearchableSelectProps = {
  options: SearchableSelectOption[]
  value?: string
  onValueChange: (value: string) => void
  placeholder: string
  searchPlaceholder: string
  emptyText: string
  disabled?: boolean
  isLoading?: boolean
  className?: string
}

export const SearchableSelect = React.forwardRef<
  HTMLButtonElement,
  SearchableSelectProps
>(
  (
    {
      options,
      value,
      onValueChange,
      placeholder,
      searchPlaceholder,
      emptyText,
      disabled = false,
      isLoading = false,
      className,
    },
    ref
  ) => {
    const [open, setOpen] = React.useState(false)
    const selectedOption = options.find((option) => option.value === value)

    return (
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            ref={ref}
            type='button'
            variant='outline'
            role='combobox'
            aria-expanded={open}
            disabled={disabled || isLoading}
            className={cn('w-full justify-between font-normal', className)}
          >
            <span className='truncate'>
              {selectedOption?.label || placeholder}
            </span>
            {isLoading ? (
              <Loader2 className='size-4 animate-spin text-muted-foreground' />
            ) : (
              <ChevronsUpDown className='size-4 text-muted-foreground' />
            )}
          </Button>
        </PopoverTrigger>
        <PopoverContent
          className='w-(--radix-popover-trigger-width) p-0'
          align='start'
        >
          <Command>
            <CommandInput placeholder={searchPlaceholder} />
            <CommandList>
              <CommandEmpty>{emptyText}</CommandEmpty>
              <CommandGroup>
                {options.map((option) => (
                  <CommandItem
                    key={option.value}
                    value={`${option.label} ${option.value}`}
                    onSelect={() => {
                      onValueChange(option.value)
                      setOpen(false)
                    }}
                  >
                    <Check
                      className={cn(
                        'size-4',
                        value === option.value ? 'opacity-100' : 'opacity-0'
                      )}
                    />
                    {option.label}
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
    )
  }
)

SearchableSelect.displayName = 'SearchableSelect'
