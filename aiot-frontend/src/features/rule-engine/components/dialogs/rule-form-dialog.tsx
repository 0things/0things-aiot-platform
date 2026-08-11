import { useEffect, useMemo, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Form } from '@/components/ui/form'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useAvailableFields } from '../../api/queries'
import {
  createRuleFormSchema,
  type CreateRuleFormData,
  type Rule,
} from '../../data/schema'
import { ActionConfig } from '../rule-form/action-config'
import { ConditionBuilder } from '../rule-form/condition-builder'
import { RuleBasicInfo } from '../rule-form/rule-basic-info'
import { SqlConfig } from '../rule-form/sql-config'
import { TriggerConfig } from '../rule-form/trigger-config'

interface RuleFormDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  mode: 'create' | 'edit'
  initialRule?: Rule | null
  pending?: boolean
  onSubmit: (data: CreateRuleFormData) => Promise<void>
}

function getDefaultValues(rule?: Rule | null): CreateRuleFormData {
  if (!rule) {
    return {
      name: '',
      description: '',
      type: 'device_linkage',
      status: 'draft',
      trigger: {
        type: 'device_data',
      },
      condition: {
        logic: 'AND',
        conditions: [],
      },
      actions: [],
      priority: 0,
      tags: [],
    }
  }

  return {
    name: rule.name,
    description: rule.description || '',
    type: rule.type,
    status: rule.status,
    trigger: rule.trigger,
    condition: rule.condition || {
      logic: 'AND',
      conditions: [],
    },
    actions: rule.actions || [],
    sqlConfig: rule.sqlConfig,
    priority: rule.priority || 0,
    tags: rule.tags || [],
  }
}

export function RuleFormDialog({
  open,
  onOpenChange,
  mode,
  initialRule,
  pending,
  onSubmit,
}: RuleFormDialogProps) {
  const [activeTab, setActiveTab] = useState('basic')
  const defaultValues = useMemo(
    () => getDefaultValues(initialRule),
    [initialRule]
  )

  const form = useForm<CreateRuleFormData>({
    resolver: zodResolver(createRuleFormSchema) as never,
    defaultValues,
  })

  useEffect(() => {
    if (open) {
      form.reset(defaultValues)
      setActiveTab('basic')
    }
  }, [defaultValues, form, open])

  const ruleType = form.watch('type')
  const productIds = form.watch('trigger.productIds')
  const primaryProductId = productIds?.[0]
  const { data: availableFields = [] } = useAvailableFields(primaryProductId)

  const tabs =
    ruleType === 'sql'
      ? ['basic', 'trigger', 'sql']
      : ['basic', 'trigger', 'condition', 'actions']

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='flex max-h-[90vh] w-[95vw] max-w-5xl flex-col'>
        <DialogHeader className='flex-shrink-0'>
          <DialogTitle>{mode === 'create' ? '创建规则' : '编辑规则'}</DialogTitle>
          <DialogDescription>
            配置规则引擎规则，支持设备联动、数据流转和告警能力
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit((data) => onSubmit(data as CreateRuleFormData))}
            className='flex flex-1 flex-col overflow-hidden'
          >
            <Tabs
              value={activeTab}
              onValueChange={setActiveTab}
              className='flex flex-1 flex-col overflow-hidden'
            >
              <TabsList
                className={
                  tabs.length === 2
                    ? 'grid w-full flex-shrink-0 grid-cols-2'
                    : tabs.length === 3
                      ? 'grid w-full flex-shrink-0 grid-cols-3'
                    : 'grid w-full flex-shrink-0 grid-cols-4'
                }
              >
                <TabsTrigger value='basic'>基本信息</TabsTrigger>
                <TabsTrigger value='trigger'>触发器</TabsTrigger>
                {ruleType === 'sql' && <TabsTrigger value='sql'>SQL配置</TabsTrigger>}
                {ruleType !== 'sql' && (
                  <TabsTrigger value='condition'>条件</TabsTrigger>
                )}
                {ruleType !== 'sql' && (
                  <TabsTrigger value='actions'>动作</TabsTrigger>
                )}
              </TabsList>

              <div className='flex-1 overflow-y-auto py-4'>
                <TabsContent value='basic' className='mt-0 space-y-4'>
                  <RuleBasicInfo form={form} />
                </TabsContent>

                <TabsContent value='trigger' className='mt-0 space-y-4'>
                  <TriggerConfig form={form} />
                </TabsContent>

                {ruleType === 'sql' && (
                  <TabsContent value='sql' className='mt-0 space-y-4'>
                    <SqlConfig form={form} />
                  </TabsContent>
                )}

                {ruleType !== 'sql' && (
                  <TabsContent value='condition' className='mt-0 space-y-4'>
                    <div>
                      <h3 className='mb-2 text-lg font-medium'>条件配置</h3>
                      <p className='mb-4 text-sm text-muted-foreground'>
                        配置触发规则的条件，支持多条件组合和嵌套
                      </p>
                      <ConditionBuilder
                        value={
                          form.watch('condition') || {
                            logic: 'AND',
                            conditions: [],
                          }
                        }
                        onChange={(value) => form.setValue('condition', value)}
                        availableFields={availableFields}
                      />
                    </div>
                  </TabsContent>
                )}

                {ruleType !== 'sql' && (
                  <TabsContent value='actions' className='mt-0 space-y-4'>
                    <ActionConfig form={form} />
                  </TabsContent>
                )}
              </div>
            </Tabs>

            <div className='flex flex-shrink-0 justify-between border-t pt-4'>
              <Button
                type='button'
                variant='outline'
                onClick={() => {
                  onOpenChange(false)
                  form.reset(defaultValues)
                }}
              >
                取消
              </Button>
              <div className='flex gap-2'>
                <Button
                  type='button'
                  variant='outline'
                  onClick={() => {
                    const currentIndex = tabs.indexOf(activeTab)
                    if (currentIndex > 0) {
                      setActiveTab(tabs[currentIndex - 1])
                    }
                  }}
                  disabled={activeTab === tabs[0]}
                >
                  上一步
                </Button>
                <Button
                  type='button'
                  variant='outline'
                  onClick={() => {
                    const currentIndex = tabs.indexOf(activeTab)
                    if (currentIndex < tabs.length - 1) {
                      setActiveTab(tabs[currentIndex + 1])
                    }
                  }}
                  disabled={activeTab === tabs[tabs.length - 1]}
                >
                  下一步
                </Button>
                <Button type='submit' disabled={pending}>
                  {pending
                    ? mode === 'create'
                      ? '创建中...'
                      : '保存中...'
                    : mode === 'create'
                      ? '创建规则'
                      : '保存修改'}
                </Button>
              </div>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
