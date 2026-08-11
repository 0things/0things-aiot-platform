import { Plus, Trash2, ChevronRight } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import type { ConditionGroup, Condition, Operator } from '../../data/schema'
import {
  conditionGroupSchema,
  conditionSchema,
  operatorLabels,
} from '../../data/schema'
import type { AvailableField } from '../../data/schema'

interface ConditionBuilderProps {
  value: ConditionGroup
  onChange: (value: ConditionGroup) => void
  availableFields: AvailableField[]
}

export function ConditionBuilder({
  value,
  onChange,
  availableFields,
}: ConditionBuilderProps) {
  type ConditionNode = Condition | ConditionGroup

  const addCondition = () => {
    onChange({
      ...value,
      conditions: [
        ...value.conditions,
        {
          field: availableFields[0]?.field || '',
          operator: 'eq',
          value: '',
        },
      ],
    })
  }

  const addGroup = () => {
    onChange({
      ...value,
      conditions: [
        ...value.conditions,
        {
          logic: 'AND',
          conditions: [],
        },
      ],
    })
  }

  const updateCondition = (
    index: number,
    newCondition: ConditionNode
  ) => {
    const newConditions = [...value.conditions]
    newConditions[index] = newCondition
    onChange({ ...value, conditions: newConditions })
  }

  const removeCondition = (index: number) => {
    onChange({
      ...value,
      conditions: value.conditions.filter(
        (_: ConditionNode, i: number) => i !== index
      ),
    })
  }

  const toggleLogic = () => {
    onChange({
      ...value,
      logic: value.logic === 'AND' ? 'OR' : 'AND',
    })
  }

  const isConditionGroup = (cond: ConditionNode): cond is ConditionGroup => {
    return conditionGroupSchema.safeParse(cond).success
  }

  const isCondition = (cond: ConditionNode): cond is Condition => {
    return conditionSchema.safeParse(cond).success
  }

  return (
    <div className='space-y-4 rounded-lg border bg-muted/30 p-4'>
      {/* 逻辑运算符切换 */}
      <div className='flex items-center gap-2'>
        <span className='text-sm text-muted-foreground'>条件组合方式:</span>
        <Button
          type='button'
          variant='outline'
          size='sm'
          onClick={toggleLogic}
          className='font-mono'
        >
          {value.logic}
        </Button>
        <span className='text-xs text-muted-foreground'>
          ({value.logic === 'AND' ? '所有条件必须满足' : '任一条件满足即可'})
        </span>
      </div>

      {/* 条件列表 */}
      <div className='space-y-2'>
        {value.conditions.map((condition: ConditionNode, index: number) => (
          <div key={index} className='flex min-w-0 items-start gap-2'>
            {/* 逻辑连接词显示 */}
            {index > 0 && (
              <Badge
                variant='secondary'
                className='mt-2 shrink-0 px-2 py-1 text-xs'
              >
                {value.logic}
              </Badge>
            )}

            {/* 条件是嵌套组 */}
            {isConditionGroup(condition) ? (
              <div className='min-w-0 flex-1'>
                <ConditionBuilder
                  value={condition}
                  onChange={(newGroup) => updateCondition(index, newGroup)}
                  availableFields={availableFields}
                />
              </div>
            ) : isCondition(condition) ? (
              /* 单条件 */
              <div className='flex min-w-0 flex-1 items-start gap-2'>
                {/* 字段选择 */}
                <Select
                  value={(condition as Condition).field}
                  onValueChange={(field) =>
                    updateCondition(index, {
                      ...(condition as Condition),
                      field,
                    })
                  }
                >
                  <SelectTrigger className='w-[140px] min-w-[120px] shrink-0'>
                    <SelectValue placeholder='选择字段' />
                  </SelectTrigger>
                  <SelectContent>
                    {availableFields.map((f) => (
                      <SelectItem key={f.field} value={f.field}>
                        <div className='flex flex-col'>
                          <span>{f.label}</span>
                          <span className='text-xs text-muted-foreground'>
                            {f.field} ({f.type})
                          </span>
                        </div>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>

                {/* 运算符选择 */}
                <Select
                  value={(condition as Condition).operator}
                  onValueChange={(operator: Operator) =>
                    updateCondition(index, {
                      ...(condition as Condition),
                      operator,
                    })
                  }
                >
                  <SelectTrigger className='w-[120px] min-w-[100px] shrink-0'>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {Object.entries(operatorLabels).map(([op, label]) => (
                      <SelectItem key={op} value={op}>
                        {label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>

                {/* 值输入 */}
                <Input
                  value={(condition as Condition).value}
                  onChange={(e) =>
                    updateCondition(index, {
                      ...(condition as Condition),
                      value: e.target.value,
                    })
                  }
                  placeholder='输入值'
                  className='min-w-[80px] flex-1'
                />
              </div>
            ) : null}

            {/* 删除按钮 */}
            <Button
              type='button'
              variant='ghost'
              size='icon'
              onClick={() => removeCondition(index)}
              className='mt-1 shrink-0'
            >
              <Trash2 className='h-4 w-4 text-destructive' />
            </Button>
          </div>
        ))}

        {value.conditions.length === 0 && (
          <div className='py-4 text-center text-sm text-muted-foreground'>
            暂无条件，点击下方按钮添加
          </div>
        )}
      </div>

      {/* 添加按钮 */}
      <div className='flex gap-2'>
        <Button
          type='button'
          variant='outline'
          size='sm'
          onClick={addCondition}
          className='gap-1'
        >
          <Plus className='h-3 w-3' />
          添加条件
        </Button>
        <Button
          type='button'
          variant='outline'
          size='sm'
          onClick={addGroup}
          className='gap-1'
        >
          <ChevronRight className='h-3 w-3' />
          添加条件组
        </Button>
      </div>

      {/* 条件预览（JSON） */}
      {value.conditions.length > 0 && (
        <details className='text-xs'>
          <summary className='cursor-pointer text-muted-foreground hover:text-foreground'>
            查看JSON (调试)
          </summary>
          <pre className='mt-2 overflow-auto rounded bg-muted p-2'>
            {JSON.stringify(value, null, 2)}
          </pre>
        </details>
      )}
    </div>
  )
}
