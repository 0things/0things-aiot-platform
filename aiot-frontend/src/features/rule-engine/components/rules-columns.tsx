import type { ColumnDef } from '@tanstack/react-table'
import {
  MoreHorizontal,
  Play,
  Edit,
  Trash2,
  Eye,
  Power,
  PowerOff,
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import type { Rule } from '../data/schema'
import {
  ruleTypeLabels,
  ruleStatusLabels,
  triggerTypeLabels,
} from '../data/schema'

interface ColumnsProps {
  onEdit: (rule: Rule) => void
  onView: (rule: Rule) => void
  onDelete: (rule: Rule) => void
  onToggleStatus: (rule: Rule) => void
  onTrigger: (rule: Rule) => void
}

export const createColumns = ({
  onEdit,
  onView,
  onDelete,
  onToggleStatus,
  onTrigger,
}: ColumnsProps): ColumnDef<Rule>[] => [
  {
    id: 'select',
    header: ({ table }) => (
      <Checkbox
        checked={table.getIsAllPageRowsSelected()}
        onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
        aria-label='Select all'
      />
    ),
    cell: ({ row }) => (
      <Checkbox
        checked={row.getIsSelected()}
        onCheckedChange={(value) => row.toggleSelected(!!value)}
        aria-label='Select row'
      />
    ),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: 'name',
    header: '规则名称',
    cell: ({ row }) => {
      const name = row.getValue('name') as string
      const description = row.original.description
      return (
        <div className='flex flex-col'>
          <span className='font-medium'>{name}</span>
          {description && (
            <span className='text-xs text-muted-foreground'>{description}</span>
          )}
        </div>
      )
    },
  },
  {
    accessorKey: 'type',
    header: '类型',
    cell: ({ row }) => {
      const type = row.getValue('type') as Rule['type']
      const label = ruleTypeLabels[type]

      const variantMap: Record<
        Rule['type'],
        'default' | 'secondary' | 'outline' | 'destructive'
      > = {
        device_linkage: 'default',
        data_forwarding: 'secondary',
        alert: 'destructive',
        sql: 'outline',
      }

      return <Badge variant={variantMap[type]}>{label}</Badge>
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id))
    },
  },
  {
    accessorKey: 'status',
    header: '状态',
    cell: ({ row }) => {
      const status = row.getValue('status') as Rule['status']
      const label = ruleStatusLabels[status]

      const variantMap: Record<
        Rule['status'],
        'default' | 'secondary' | 'outline'
      > = {
        enabled: 'default',
        disabled: 'secondary',
        draft: 'outline',
      }

      return <Badge variant={variantMap[status]}>{label}</Badge>
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id))
    },
  },
  {
    accessorKey: 'trigger',
    header: '触发器',
    cell: ({ row }) => {
      const trigger = row.original.trigger
      return (
        <div className='text-sm'>
          {triggerTypeLabels[trigger.type]}
          {trigger.productIds && trigger.productIds.length > 0 && (
            <div className='text-xs text-muted-foreground'>
              产品: {trigger.productIds.length} 个
            </div>
          )}
        </div>
      )
    },
  },
  {
    accessorKey: 'executionCount',
    header: '执行次数',
    cell: ({ row }) => {
      const count = row.getValue('executionCount') as number
      const successCount = row.original.successCount
      const successRate =
        count > 0 ? ((successCount / count) * 100).toFixed(1) : '0'

      return (
        <div className='flex flex-col'>
          <span className='font-mono text-sm'>{count.toLocaleString()}</span>
          <span className='text-xs text-muted-foreground'>
            成功率: {successRate}%
          </span>
        </div>
      )
    },
  },
  {
    accessorKey: 'lastExecutedAt',
    header: '最后执行',
    cell: ({ row }) => {
      const lastExecutedAt = row.getValue('lastExecutedAt') as
        | string
        | undefined
      const lastStatus = row.original.lastExecutionStatus

      if (!lastExecutedAt) {
        return <span className='text-sm text-muted-foreground'>未执行</span>
      }

      const date = new Date(lastExecutedAt)
      const statusColors = {
        success: 'text-green-600',
        failure: 'text-red-600',
        pending: 'text-gray-600',
      }

      return (
        <div className='flex flex-col'>
          <span className='text-sm'>{date.toLocaleString('zh-CN')}</span>
          {lastStatus && (
            <span className={`text-xs ${statusColors[lastStatus]}`}>
              {lastStatus === 'success'
                ? '成功'
                : lastStatus === 'failure'
                  ? '失败'
                  : '待执行'}
            </span>
          )}
        </div>
      )
    },
  },
  {
    accessorKey: 'priority',
    header: '优先级',
    cell: ({ row }) => {
      const priority = row.getValue('priority') as number
      return <span className='font-mono text-sm'>{priority}</span>
    },
  },
  {
    accessorKey: 'tags',
    header: '标签',
    cell: ({ row }) => {
      const tags = row.original.tags
      if (!tags || tags.length === 0) {
        return <span className='text-sm text-muted-foreground'>-</span>
      }
      return (
        <div className='flex flex-wrap gap-1'>
          {tags.slice(0, 2).map((tag, index) => (
            <Badge key={index} variant='outline' className='text-xs'>
              {tag}
            </Badge>
          ))}
          {tags.length > 2 && (
            <Badge variant='outline' className='text-xs'>
              +{tags.length - 2}
            </Badge>
          )}
        </div>
      )
    },
  },
  {
    id: 'actions',
    header: '操作',
    cell: ({ row }) => {
      const rule = row.original

      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant='ghost' className='h-8 w-8 p-0'>
              <span className='sr-only'>打开菜单</span>
              <MoreHorizontal className='h-4 w-4' />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align='end'>
            <DropdownMenuLabel>操作</DropdownMenuLabel>
            <DropdownMenuItem onClick={() => onView(rule)}>
              <Eye className='mr-2 h-4 w-4' />
              查看详情
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => onEdit(rule)}>
              <Edit className='mr-2 h-4 w-4' />
              编辑
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={() => onToggleStatus(rule)}>
              {rule.status === 'enabled' ? (
                <>
                  <PowerOff className='mr-2 h-4 w-4' />
                  禁用规则
                </>
              ) : (
                <>
                  <Power className='mr-2 h-4 w-4' />
                  启用规则
                </>
              )}
            </DropdownMenuItem>
            <DropdownMenuItem
              onClick={() => onTrigger(rule)}
              disabled={rule.status !== 'enabled'}
            >
              <Play className='mr-2 h-4 w-4' />
              手动触发
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              onClick={() => onDelete(rule)}
              className='text-destructive'
            >
              <Trash2 className='mr-2 h-4 w-4' />
              删除
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      )
    },
    meta: {
      className: cn(
        'sticky end-0 z-10 bg-background',
        'shadow-[-4px_0_6px_-2px_rgb(0_0_0_/_0.05)] dark:shadow-[-4px_0_6px_-2px_rgb(0_0_0_/_0.3)]'
      ),
    },
    enableHiding: false,
  },
]
