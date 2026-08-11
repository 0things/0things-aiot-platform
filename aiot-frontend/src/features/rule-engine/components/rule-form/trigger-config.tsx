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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import type { CreateRuleFormData } from '../../data/schema'
import { triggerTypeLabels } from '../../data/schema'

interface TriggerConfigProps {
  form: UseFormReturn<CreateRuleFormData>
}

export function TriggerConfig({ form }: TriggerConfigProps) {
  const triggerType = form.watch('trigger.type')

  return (
    <div className='space-y-4'>
      {/* 触发器类型 */}
      <FormField
        control={form.control}
        name='trigger.type'
        render={({ field }) => (
          <FormItem>
            <FormLabel>触发器类型</FormLabel>
            <Select onValueChange={field.onChange} defaultValue={field.value}>
              <FormControl>
                <SelectTrigger>
                  <SelectValue placeholder='选择触发器类型' />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                {Object.entries(triggerTypeLabels).map(([value, label]) => (
                  <SelectItem key={value} value={value}>
                    {label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <FormDescription>选择何时触发此规则</FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      {/* 产品ID（仅数据相关触发器显示） */}
      {(triggerType === 'device_data' || triggerType === 'device_status') && (
        <FormField
          control={form.control}
          name='trigger.productIds'
          render={({ field }) => (
            <FormItem>
              <FormLabel>产品ID（可选）</FormLabel>
              <FormControl>
                <Input
                  placeholder='输入产品ID，多个用逗号分隔'
                  value={field.value?.join(',') || ''}
                  onChange={(e) => {
                    const value = e.target.value
                    field.onChange(
                      value ? value.split(',').map((id) => id.trim()) : []
                    )
                  }}
                />
              </FormControl>
              <FormDescription>
                限定触发规则的产品范围，留空表示所有产品
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
      )}

      {/* Topic（仅数据上报触发器显示） */}
      {triggerType === 'device_data' && (
        <FormField
          control={form.control}
          name='trigger.topic'
          render={({ field }) => (
            <FormItem>
              <FormLabel>MQTT Topic（可选）</FormLabel>
              <FormControl>
                <Input
                  placeholder='例如: device/+/data'
                  {...field}
                  value={field.value || ''}
                />
              </FormControl>
              <FormDescription>
                指定监听的MQTT topic，支持通配符（+ 和 #）
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
      )}

      {/* Cron表达式（仅定时触发器显示） */}
      {triggerType === 'schedule' && (
        <FormField
          control={form.control}
          name='trigger.schedule'
          render={({ field }) => (
            <FormItem>
              <FormLabel>Cron表达式</FormLabel>
              <FormControl>
                <Input
                  placeholder='例如: 0 0 * * * (每天0点)'
                  {...field}
                  value={field.value || ''}
                />
              </FormControl>
              <FormDescription>
                使用标准cron表达式设置定时任务。
                <a
                  href='https://crontab.guru/'
                  target='_blank'
                  rel='noopener noreferrer'
                  className='ml-1 text-primary underline'
                >
                  查看Cron语法帮助
                </a>
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
      )}

      {/* 触发器说明 */}
      <div className='rounded-lg bg-muted p-4'>
        <h4 className='mb-2 text-sm font-medium'>触发器说明</h4>
        <p className='text-sm text-muted-foreground'>
          {triggerType === 'device_data' &&
            '当设备上报数据时触发，可以通过产品ID和Topic进一步过滤'}
          {triggerType === 'device_status' && '当设备状态发生变化时触发'}
          {triggerType === 'device_online' && '当设备上线时触发'}
          {triggerType === 'device_offline' && '当设备离线时触发'}
          {triggerType === 'schedule' &&
            '按照Cron表达式定时触发，适合周期性任务'}
          {triggerType === 'manual' && '仅手动触发，不会自动执行'}
        </p>
      </div>
    </div>
  )
}
