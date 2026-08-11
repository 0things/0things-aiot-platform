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
import { Textarea } from '@/components/ui/textarea'
import type { CreateRuleFormData } from '../../data/schema'
import { ruleTypeLabels, ruleStatusLabels } from '../../data/schema'

interface RuleBasicInfoProps {
  form: UseFormReturn<CreateRuleFormData>
}

export function RuleBasicInfo({ form }: RuleBasicInfoProps) {
  return (
    <div className='space-y-4'>
      {/* 规则名称 */}
      <FormField
        control={form.control}
        name='name'
        render={({ field }) => (
          <FormItem>
            <FormLabel>规则名称</FormLabel>
            <FormControl>
              <Input placeholder='输入规则名称' {...field} />
            </FormControl>
            <FormDescription>
              建议使用简洁明了的名称，例如：温度过高自动开空调
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      {/* 规则描述 */}
      <FormField
        control={form.control}
        name='description'
        render={({ field }) => (
          <FormItem>
            <FormLabel>规则描述（可选）</FormLabel>
            <FormControl>
              <Textarea
                placeholder='描述此规则的用途和工作方式'
                rows={3}
                {...field}
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />

      {/* 规则类型 */}
      <FormField
        control={form.control}
        name='type'
        render={({ field }) => (
          <FormItem>
            <FormLabel>规则类型</FormLabel>
            <Select onValueChange={field.onChange} defaultValue={field.value}>
              <FormControl>
                <SelectTrigger>
                  <SelectValue placeholder='选择规则类型' />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                {Object.entries(ruleTypeLabels).map(([value, label]) => (
                  <SelectItem key={value} value={value}>
                    {label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <FormDescription>
              {field.value === 'device_linkage' &&
                '设备之间的联动控制，例如：传感器触发设备动作'}
              {field.value === 'data_forwarding' &&
                '将设备数据转发到外部系统或服务'}
              {field.value === 'alert' &&
                '基于条件触发告警通知（邮件、短信等）'}
              {field.value === 'sql' && '使用SQL语句处理和分析设备数据流'}
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      {/* 规则状态 */}
      <FormField
        control={form.control}
        name='status'
        render={({ field }) => (
          <FormItem>
            <FormLabel>规则状态</FormLabel>
            <Select onValueChange={field.onChange} defaultValue={field.value}>
              <FormControl>
                <SelectTrigger>
                  <SelectValue placeholder='选择规则状态' />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                {Object.entries(ruleStatusLabels).map(([value, label]) => (
                  <SelectItem key={value} value={value}>
                    {label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <FormDescription>
              草稿状态的规则不会执行，启用后才会生效
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      {/* 优先级 */}
      <FormField
        control={form.control}
        name='priority'
        render={({ field }) => (
          <FormItem>
            <FormLabel>优先级</FormLabel>
            <FormControl>
              <Input
                type='number'
                min={0}
                max={100}
                placeholder='0-100，数字越大优先级越高'
                {...field}
                onChange={(e) => field.onChange(parseInt(e.target.value))}
              />
            </FormControl>
            <FormDescription>
              当多个规则同时触发时，优先级高的规则先执行
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />
    </div>
  )
}
