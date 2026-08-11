import type { UseFormReturn } from 'react-hook-form'
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import type { CreateRuleFormData } from '../../data/schema'

interface SqlConfigProps {
  form: UseFormReturn<CreateRuleFormData>
}

export function SqlConfig({ form }: SqlConfigProps) {
  return (
    <div className='space-y-4'>
      <FormField
        control={form.control}
        name='sqlConfig.sql'
        render={({ field }) => (
          <FormItem>
            <FormLabel>SQL 查询语句</FormLabel>
            <FormControl>
              <Textarea
                rows={6}
                placeholder='SELECT * FROM telemetry_stream WHERE temperature > 30'
                {...field}
                value={field.value || ''}
              />
            </FormControl>
            <FormDescription>
              用于处理设备数据流的 SQL 查询语句
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={form.control}
        name='sqlConfig.dataSource'
        render={({ field }) => (
          <FormItem>
            <FormLabel>数据源</FormLabel>
            <FormControl>
              <Input
                placeholder='device_stream'
                {...field}
                value={field.value || ''}
              />
            </FormControl>
            <FormDescription>可选，指定 SQL 执行的数据源名称</FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={form.control}
        name='sqlConfig.outputTopic'
        render={({ field }) => (
          <FormItem>
            <FormLabel>输出 Topic</FormLabel>
            <FormControl>
              <Input
                placeholder='analytics/temperature/hourly'
                {...field}
                value={field.value || ''}
              />
            </FormControl>
            <FormDescription>可选，指定 SQL 结果输出的 Topic</FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />
    </div>
  )
}
