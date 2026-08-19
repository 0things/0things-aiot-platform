import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useTranslation } from 'react-i18next'
import { showSubmittedData } from '@/lib/show-submitted-data'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'

const items = [
  {
    id: 'recents',
    translationKey: 'display.sidebar.items.recents',
  },
  {
    id: 'home',
    translationKey: 'display.sidebar.items.home',
  },
  {
    id: 'applications',
    translationKey: 'display.sidebar.items.applications',
  },
  {
    id: 'desktop',
    translationKey: 'display.sidebar.items.desktop',
  },
  {
    id: 'downloads',
    translationKey: 'display.sidebar.items.downloads',
  },
  {
    id: 'documents',
    translationKey: 'display.sidebar.items.documents',
  },
] as const

export function DisplayForm() {
  const { t } = useTranslation('settings')

  const displayFormSchema = z.object({
    items: z.array(z.string()).refine((value) => value.some((item) => item), {
      message: t('display.validation.selectAtLeastOne'),
    }),
  })

  type DisplayFormValues = z.infer<typeof displayFormSchema>

  const defaultValues: Partial<DisplayFormValues> = {
    items: ['recents', 'home'],
  }

  const form = useForm<DisplayFormValues>({
    resolver: zodResolver(displayFormSchema),
    defaultValues,
  })

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit((data) => showSubmittedData(data))}
        className='space-y-8'
      >
        <FormField
          control={form.control}
          name='items'
          render={() => (
            <FormItem>
              <div className='mb-4'>
                <FormLabel className='text-base'>
                  {t('display.sidebar.label')}
                </FormLabel>
                <FormDescription>
                  {t('display.sidebar.description')}
                </FormDescription>
              </div>
              {items.map((item) => (
                <FormField
                  key={item.id}
                  control={form.control}
                  name='items'
                  render={({ field }) => {
                    return (
                      <FormItem
                        key={item.id}
                        className='flex flex-row items-start'
                      >
                        <FormControl>
                          <Checkbox
                            checked={field.value?.includes(item.id)}
                            onCheckedChange={(checked) => {
                              return checked
                                ? field.onChange([...field.value, item.id])
                                : field.onChange(
                                    field.value?.filter(
                                      (value) => value !== item.id
                                    )
                                  )
                            }}
                          />
                        </FormControl>
                        <FormLabel className='font-normal'>
                          {t(item.translationKey)}
                        </FormLabel>
                      </FormItem>
                    )
                  }}
                />
              ))}
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type='submit'>{t('display.updateButton')}</Button>
      </form>
    </Form>
  )
}
