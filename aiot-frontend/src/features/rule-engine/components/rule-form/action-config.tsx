import { useState } from 'react'
import type { UseFormReturn } from 'react-hook-form'
import { Plus, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import type { CreateRuleFormData, ActionType } from '../../data/schema'
import { actionTypeLabels, httpMethodEnum } from '../../data/schema'

interface ActionConfigProps {
  form: UseFormReturn<CreateRuleFormData>
}

export function ActionConfig({ form }: ActionConfigProps) {
  const actions = form.watch('actions') || []
  const [expandedActionIndex, setExpandedActionIndex] = useState<number | null>(
    0
  )

  const addAction = (type: ActionType) => {
    const newAction = createDefaultAction(type)
    form.setValue('actions', [...actions, newAction])
    setExpandedActionIndex(actions.length)
  }

  const removeAction = (index: number) => {
    const newActions = actions.filter((_, i) => i !== index)
    form.setValue('actions', newActions)
    if (expandedActionIndex === index) {
      setExpandedActionIndex(null)
    }
  }

  const updateAction = (index: number, field: string, value: any) => {
    const newActions = [...actions]
    const action = { ...newActions[index] }

    // 使用点号分隔的路径更新嵌套属性
    const keys = field.split('.')
    let current: any = action
    for (let i = 0; i < keys.length - 1; i++) {
      current = current[keys[i]]
    }
    current[keys[keys.length - 1]] = value

    newActions[index] = action
    form.setValue('actions', newActions)
  }

  return (
    <div className='space-y-4'>
      <div className='flex items-center justify-between'>
        <div>
          <h3 className='text-lg font-medium'>动作配置</h3>
          <p className='text-sm text-muted-foreground'>
            配置规则触发后要执行的动作
          </p>
        </div>
      </div>

      {/* 动作列表 */}
      <div className='space-y-3'>
        {actions.map((action, index) => (
          <Card key={index}>
            <CardHeader className='pb-3'>
              <div className='flex items-center justify-between'>
                <div className='flex items-center gap-2'>
                  <CardTitle className='text-base'>
                    动作 {index + 1}: {actionTypeLabels[action.type]}
                  </CardTitle>
                </div>
                <div className='flex items-center gap-1'>
                  <Button
                    type='button'
                    variant='ghost'
                    size='sm'
                    onClick={() =>
                      setExpandedActionIndex(
                        expandedActionIndex === index ? null : index
                      )
                    }
                  >
                    {expandedActionIndex === index ? '收起' : '展开'}
                  </Button>
                  <Button
                    type='button'
                    variant='ghost'
                    size='icon'
                    onClick={() => removeAction(index)}
                  >
                    <Trash2 className='h-4 w-4 text-destructive' />
                  </Button>
                </div>
              </div>
            </CardHeader>

            {expandedActionIndex === index && (
              <CardContent className='space-y-3'>
                {/* 设备控制动作 */}
                {action.type === 'device_control' && (
                  <>
                    <div>
                      <label className='text-sm font-medium'>目标设备ID</label>
                      <Input
                        value={action.params.targetDeviceId || ''}
                        onChange={(e) =>
                          updateAction(
                            index,
                            'params.targetDeviceId',
                            e.target.value
                          )
                        }
                        placeholder='输入设备ID'
                      />
                    </div>
                    <div>
                      <label className='text-sm font-medium'>指令名称</label>
                      <Input
                        value={action.params.command || ''}
                        onChange={(e) =>
                          updateAction(index, 'params.command', e.target.value)
                        }
                        placeholder='例如: turnOn, setTemperature'
                      />
                    </div>
                    <div>
                      <label className='text-sm font-medium'>
                        指令参数（JSON）
                      </label>
                      <Textarea
                        value={JSON.stringify(
                          action.params.params || {},
                          null,
                          2
                        )}
                        onChange={(e) => {
                          try {
                            const parsed = JSON.parse(e.target.value)
                            updateAction(index, 'params.params', parsed)
                          } catch {
                            // 忽略JSON解析错误
                          }
                        }}
                        placeholder='{"temperature": 26}'
                        rows={3}
                      />
                    </div>
                  </>
                )}

                {/* HTTP请求动作 */}
                {action.type === 'http_request' && (
                  <>
                    <div>
                      <label className='text-sm font-medium'>请求URL</label>
                      <Input
                        value={action.params.url || ''}
                        onChange={(e) =>
                          updateAction(index, 'params.url', e.target.value)
                        }
                        placeholder='https://api.example.com/webhook'
                      />
                    </div>
                    <div>
                      <label className='text-sm font-medium'>HTTP方法</label>
                      <Select
                        value={action.params.method || 'POST'}
                        onValueChange={(value) =>
                          updateAction(index, 'params.method', value)
                        }
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          {httpMethodEnum.options.map((method) => (
                            <SelectItem key={method} value={method}>
                              {method}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <div>
                      <label className='text-sm font-medium'>
                        请求体（支持模板变量）
                      </label>
                      <Textarea
                        value={action.params.body || ''}
                        onChange={(e) =>
                          updateAction(index, 'params.body', e.target.value)
                        }
                        placeholder='{"deviceId": "${deviceId}", "value": ${temperature}}'
                        rows={4}
                      />
                      <p className='mt-1 text-xs text-muted-foreground'>
                        使用 ${`{变量名}`} 引用触发数据中的字段
                      </p>
                    </div>
                  </>
                )}

                {/* Webhook动作 */}
                {action.type === 'webhook' && (
                  <>
                    <div>
                      <label className='text-sm font-medium'>Webhook URL</label>
                      <Input
                        value={action.params.url || ''}
                        onChange={(e) =>
                          updateAction(index, 'params.url', e.target.value)
                        }
                        placeholder='https://your-server.com/webhook'
                      />
                    </div>
                    <div>
                      <label className='text-sm font-medium'>
                        消息模板（可选）
                      </label>
                      <Textarea
                        value={action.params.bodyTemplate || ''}
                        onChange={(e) =>
                          updateAction(
                            index,
                            'params.bodyTemplate',
                            e.target.value
                          )
                        }
                        placeholder='自定义消息格式，支持模板变量'
                        rows={3}
                      />
                    </div>
                  </>
                )}

                {/* 邮件通知动作 */}
                {action.type === 'email' && (
                  <>
                    <div>
                      <label className='text-sm font-medium'>
                        收件人（多个用逗号分隔）
                      </label>
                      <Input
                        value={action.params.to?.join(',') || ''}
                        onChange={(e) =>
                          updateAction(
                            index,
                            'params.to',
                            e.target.value
                              .split(',')
                              .map((email) => email.trim())
                          )
                        }
                        placeholder='admin@example.com, ops@example.com'
                      />
                    </div>
                    <div>
                      <label className='text-sm font-medium'>邮件主题</label>
                      <Input
                        value={action.params.subject || ''}
                        onChange={(e) =>
                          updateAction(index, 'params.subject', e.target.value)
                        }
                        placeholder='告警：温度异常'
                      />
                    </div>
                    <div>
                      <label className='text-sm font-medium'>
                        邮件内容（支持HTML）
                      </label>
                      <Textarea
                        value={action.params.body || ''}
                        onChange={(e) =>
                          updateAction(index, 'params.body', e.target.value)
                        }
                        placeholder='<p>设备 ${deviceId} 温度达到 ${temperature}℃</p>'
                        rows={4}
                      />
                    </div>
                  </>
                )}
              </CardContent>
            )}
          </Card>
        ))}

        {actions.length === 0 && (
          <Card>
            <CardContent className='flex flex-col items-center justify-center py-8'>
              <p className='mb-4 text-sm text-muted-foreground'>
                暂无动作配置，请添加至少一个动作
              </p>
            </CardContent>
          </Card>
        )}
      </div>

      {/* 添加动作按钮 */}
      <div className='flex flex-wrap gap-2'>
        <Button
          type='button'
          variant='outline'
          size='sm'
          onClick={() => addAction('device_control')}
        >
          <Plus className='mr-1 h-4 w-4' />
          设备控制
        </Button>
        <Button
          type='button'
          variant='outline'
          size='sm'
          onClick={() => addAction('http_request')}
        >
          <Plus className='mr-1 h-4 w-4' />
          HTTP请求
        </Button>
        <Button
          type='button'
          variant='outline'
          size='sm'
          onClick={() => addAction('webhook')}
        >
          <Plus className='mr-1 h-4 w-4' />
          Webhook
        </Button>
        <Button
          type='button'
          variant='outline'
          size='sm'
          onClick={() => addAction('email')}
        >
          <Plus className='mr-1 h-4 w-4' />
          邮件通知
        </Button>
      </div>
    </div>
  )
}

// 创建默认动作配置
function createDefaultAction(type: ActionType): any {
  switch (type) {
    case 'device_control':
      return {
        type,
        params: {
          targetDeviceId: '',
          command: '',
          params: {},
        },
      }
    case 'http_request':
      return {
        type,
        params: {
          url: '',
          method: 'POST',
          body: '',
        },
      }
    case 'webhook':
      return {
        type,
        params: {
          url: '',
          method: 'POST',
        },
      }
    case 'email':
      return {
        type,
        params: {
          to: [],
          subject: '',
          body: '',
        },
      }
    default:
      return {
        type,
        params: {},
      }
  }
}
