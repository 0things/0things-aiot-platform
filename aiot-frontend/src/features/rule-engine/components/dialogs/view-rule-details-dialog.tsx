import { Play, Edit } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Separator } from '@/components/ui/separator'
import { useRuleExecutions } from '../../api/queries'
import type { Rule } from '../../data/schema'
import {
  ruleTypeLabels,
  ruleStatusLabels,
  triggerTypeLabels,
  actionTypeLabels,
} from '../../data/schema'

interface ViewRuleDetailsDialogProps {
  rule: Rule | null
  open: boolean
  onOpenChange: (open: boolean) => void
  onEdit?: (rule: Rule) => void
  onTrigger?: (rule: Rule) => void
}

export function ViewRuleDetailsDialog({
  rule,
  open,
  onOpenChange,
  onEdit,
  onTrigger,
}: ViewRuleDetailsDialogProps) {
  const { data: executionsData, isLoading: executionsLoading } =
    useRuleExecutions(rule?.id, {
      page: 1,
      pageSize: 10,
    })

  if (!rule) return null

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='max-h-[90vh] max-w-3xl overflow-y-auto'>
        <DialogHeader>
          <div className='flex items-start justify-between'>
            <div className='flex-1'>
              <DialogTitle className='text-2xl'>{rule.name}</DialogTitle>
              {rule.description && (
                <p className='mt-2 text-sm text-muted-foreground'>
                  {rule.description}
                </p>
              )}
            </div>
            <div className='flex gap-2'>
              {onEdit && (
                <Button
                  variant='outline'
                  size='sm'
                  onClick={() => {
                    onEdit(rule)
                    onOpenChange(false)
                  }}
                >
                  <Edit className='mr-1 h-4 w-4' />
                  编辑
                </Button>
              )}
              {onTrigger && rule.status === 'enabled' && (
                <Button
                  variant='outline'
                  size='sm'
                  onClick={() => {
                    onTrigger(rule)
                    onOpenChange(false)
                  }}
                >
                  <Play className='mr-1 h-4 w-4' />
                  触发
                </Button>
              )}
            </div>
          </div>
        </DialogHeader>

        <div className='space-y-6'>
          {/* 基本信息 */}
          <div>
            <h3 className='mb-3 font-semibold'>基本信息</h3>
            <div className='grid grid-cols-2 gap-4 text-sm'>
              <div>
                <span className='text-muted-foreground'>规则类型:</span>
                <Badge className='ml-2' variant='outline'>
                  {ruleTypeLabels[rule.type]}
                </Badge>
              </div>
              <div>
                <span className='text-muted-foreground'>规则状态:</span>
                <Badge className='ml-2' variant='outline'>
                  {ruleStatusLabels[rule.status]}
                </Badge>
              </div>
              <div>
                <span className='text-muted-foreground'>优先级:</span>
                <span className='ml-2 font-mono'>{rule.priority}</span>
              </div>
              <div>
                <span className='text-muted-foreground'>创建时间:</span>
                <span className='ml-2'>
                  {new Date(rule.createdAt).toLocaleString('zh-CN')}
                </span>
              </div>
            </div>
          </div>

          <Separator />

          {/* 触发器配置 */}
          <div>
            <h3 className='mb-3 font-semibold'>触发器配置</h3>
            <div className='space-y-2 text-sm'>
              <div>
                <span className='text-muted-foreground'>触发类型:</span>
                <span className='ml-2'>
                  {triggerTypeLabels[rule.trigger.type]}
                </span>
              </div>
              {rule.trigger.productIds &&
                rule.trigger.productIds.length > 0 && (
                  <div>
                    <span className='text-muted-foreground'>产品ID:</span>
                    <div className='mt-1 ml-2 flex flex-wrap gap-1'>
                      {rule.trigger.productIds.map((id) => (
                        <Badge key={id} variant='secondary'>
                          {id}
                        </Badge>
                      ))}
                    </div>
                  </div>
                )}
              {rule.trigger.topic && (
                <div>
                  <span className='text-muted-foreground'>MQTT Topic:</span>
                  <code className='ml-2 rounded bg-muted px-2 py-1'>
                    {rule.trigger.topic}
                  </code>
                </div>
              )}
              {rule.trigger.schedule && (
                <div>
                  <span className='text-muted-foreground'>Cron表达式:</span>
                  <code className='ml-2 rounded bg-muted px-2 py-1'>
                    {rule.trigger.schedule}
                  </code>
                </div>
              )}
            </div>
          </div>

          <Separator />

          {/* 条件配置 */}
          {rule.condition && rule.condition.conditions.length > 0 && (
            <>
              <div>
                <h3 className='mb-3 font-semibold'>条件配置</h3>
                <pre className='overflow-auto rounded-lg bg-muted p-3 text-xs'>
                  {JSON.stringify(rule.condition, null, 2)}
                </pre>
              </div>
              <Separator />
            </>
          )}

          {/* 动作配置 */}
          {rule.actions && rule.actions.length > 0 && (
            <>
              <div>
                <h3 className='mb-3 font-semibold'>动作配置</h3>
                <div className='space-y-3'>
                  {rule.actions.map((action, index) => (
                    <div
                      key={index}
                      className='rounded-lg border bg-muted/30 p-3'
                    >
                      <div className='mb-2 font-medium'>
                        动作 {index + 1}: {actionTypeLabels[action.type]}
                      </div>
                      <pre className='overflow-auto rounded bg-muted p-2 text-xs'>
                        {JSON.stringify(action.params, null, 2)}
                      </pre>
                    </div>
                  ))}
                </div>
              </div>
              <Separator />
            </>
          )}

          {/* SQL配置 */}
          {rule.sqlConfig && (
            <>
              <div>
                <h3 className='mb-3 font-semibold'>SQL配置</h3>
                <div className='space-y-2'>
                  <div>
                    <span className='text-sm text-muted-foreground'>
                      SQL查询:
                    </span>
                    <pre className='mt-1 overflow-auto rounded-lg bg-muted p-3 text-xs'>
                      {rule.sqlConfig.sql}
                    </pre>
                  </div>
                  {rule.sqlConfig.dataSource && (
                    <div className='text-sm'>
                      <span className='text-muted-foreground'>数据源:</span>
                      <span className='ml-2'>{rule.sqlConfig.dataSource}</span>
                    </div>
                  )}
                  {rule.sqlConfig.outputTopic && (
                    <div className='text-sm'>
                      <span className='text-muted-foreground'>输出Topic:</span>
                      <span className='ml-2'>{rule.sqlConfig.outputTopic}</span>
                    </div>
                  )}
                </div>
              </div>
              <Separator />
            </>
          )}

          {/* 执行统计 */}
          <div>
            <h3 className='mb-3 font-semibold'>执行统计</h3>
            <div className='grid grid-cols-3 gap-4'>
              <div className='rounded-lg bg-muted p-3 text-center'>
                <div className='text-2xl font-bold'>{rule.executionCount}</div>
                <div className='text-xs text-muted-foreground'>总执行次数</div>
              </div>
              <div className='rounded-lg bg-green-100 p-3 text-center dark:bg-green-900/20'>
                <div className='text-2xl font-bold text-green-600 dark:text-green-400'>
                  {rule.successCount}
                </div>
                <div className='text-xs text-muted-foreground'>成功次数</div>
              </div>
              <div className='rounded-lg bg-red-100 p-3 text-center dark:bg-red-900/20'>
                <div className='text-2xl font-bold text-red-600 dark:text-red-400'>
                  {rule.failureCount}
                </div>
                <div className='text-xs text-muted-foreground'>失败次数</div>
              </div>
            </div>
            {rule.lastExecutedAt && (
              <div className='mt-2 text-sm text-muted-foreground'>
                最后执行时间:{' '}
                {new Date(rule.lastExecutedAt).toLocaleString('zh-CN')}
              </div>
            )}
          </div>

          <Separator />

          <div>
            <h3 className='mb-3 font-semibold'>最近执行记录</h3>
            <div className='rounded-lg border'>
              <div className='grid grid-cols-[160px_96px_1fr] gap-3 border-b bg-muted/40 px-4 py-2 text-xs font-medium text-muted-foreground'>
                <span>触发时间</span>
                <span>状态</span>
                <span>错误信息</span>
              </div>
              {executionsLoading ? (
                <div className='px-4 py-6 text-sm text-muted-foreground'>
                  加载中...
                </div>
              ) : executionsData?.items?.length ? (
                executionsData.items.map((execution: typeof executionsData.items[number]) => (
                  <div
                    key={execution.id}
                    className='grid grid-cols-[160px_96px_1fr] gap-3 border-b px-4 py-3 text-sm last:border-b-0'
                  >
                    <span>
                      {new Date(execution.triggeredAt).toLocaleString('zh-CN')}
                    </span>
                    <span>
                      <Badge
                        variant={
                          execution.status === 'success'
                            ? 'default'
                            : execution.status === 'failure'
                              ? 'destructive'
                              : 'secondary'
                        }
                      >
                        {execution.status === 'success'
                          ? '成功'
                          : execution.status === 'failure'
                            ? '失败'
                            : '跳过'}
                      </Badge>
                    </span>
                    <span className='text-muted-foreground'>
                      {execution.error || '-'}
                    </span>
                  </div>
                ))
              ) : (
                <div className='px-4 py-6 text-sm text-muted-foreground'>
                  暂无执行记录
                </div>
              )}
            </div>
          </div>

          {/* 标签 */}
          {rule.tags && rule.tags.length > 0 && (
            <>
              <Separator />
              <div>
                <h3 className='mb-3 font-semibold'>标签</h3>
                <div className='flex flex-wrap gap-2'>
                  {rule.tags.map((tag, index) => (
                    <Badge key={index} variant='outline'>
                      {tag}
                    </Badge>
                  ))}
                </div>
              </div>
            </>
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}
