import React from 'react'
import { Check, Languages } from 'lucide-react'
import { useTranslation } from 'react-i18next'
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

const languages = [
  {
    value: 'en',
    label: 'English',
    flag: '🇺🇸',
  },
  {
    value: 'zh',
    label: '中文',
    flag: '🇨🇳',
  },
]

export function LanguageSwitcher() {
  const { i18n } = useTranslation()
  const [open, setOpen] = React.useState(false)

  const currentLanguage = languages.find((lang) => lang.value === i18n.language)

  const handleLanguageChange = (langValue: string) => {
    i18n.changeLanguage(langValue)
    setOpen(false)
  }

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant='outline'
          role='combobox'
          aria-expanded={open}
          className='w-full justify-between'
        >
          <span className='flex items-center gap-2'>
            <Languages className='size-4' />
            <span>
              {currentLanguage?.flag} {currentLanguage?.label}
            </span>
          </span>
        </Button>
      </PopoverTrigger>
      <PopoverContent className='w-[200px] p-0'>
        <Command>
          <CommandInput placeholder='Search language...' />
          <CommandList>
            <CommandEmpty>No language found.</CommandEmpty>
            <CommandGroup>
              {languages.map((language) => (
                <CommandItem
                  key={language.value}
                  value={language.value}
                  onSelect={handleLanguageChange}
                >
                  <Check
                    className={cn(
                      'me-2 size-4',
                      i18n.language === language.value
                        ? 'opacity-100'
                        : 'opacity-0'
                    )}
                  />
                  <span className='flex items-center gap-2'>
                    <span>{language.flag}</span>
                    <span>{language.label}</span>
                  </span>
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  )
}
