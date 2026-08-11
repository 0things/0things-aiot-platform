import { useState } from 'react'
import { Plus } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import {
  useRules,
  useDeleteRules,
  useUpdateRuleStatus,
  useTriggerRule,
} from './api/queries'
import { CreateRuleDialog } from './components/dialogs/create-rule-dialog'
import { DeleteRuleDialog } from './components/dialogs/delete-rule-dialog'
import { EditRuleDialog } from './components/dialogs/edit-rule-dialog'
import { ViewRuleDetailsDialog } from './components/dialogs/view-rule-details-dialog'
import { RulesTable } from './components/rules-table'
import type { Rule } from './data/schema'

export default function RuleEnginePage() {
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [editDialogOpen, setEditDialogOpen] = useState(false)
  const [viewDialogOpen, setViewDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [selectedRule, setSelectedRule] = useState<Rule | null>(null)

  // 查询规则列表
  const { data: rulesData, isLoading } = useRules({
    page: 1,
    pageSize: 100,
  })

  // Mutations
  const deleteRulesMutation = useDeleteRules()
  const updateRuleStatusMutation = useUpdateRuleStatus()
  const triggerRuleMutation = useTriggerRule()

  // 处理编辑
  const handleEdit = (rule: Rule) => {
    setSelectedRule(rule)
    setEditDialogOpen(true)
  }

  // 处理查看
  const handleView = (rule: Rule) => {
    setSelectedRule(rule)
    setViewDialogOpen(true)
  }

  // 处理删除
  const handleDelete = (rule: Rule) => {
    setSelectedRule(rule)
    setDeleteDialogOpen(true)
  }

  // 处理批量删除
  const handleBulkDelete = async (ruleIds: string[]) => {
    try {
      await deleteRulesMutation.mutateAsync(ruleIds)
      toast.success(`已删除 ${ruleIds.length} 条规则`)
    } catch (error: any) {
      toast.error('批量删除失败: ' + error.message)
    }
  }

  // 处理启用/禁用
  const handleToggleStatus = async (rule: Rule) => {
    const newStatus = rule.status === 'enabled' ? 'disabled' : 'enabled'
    try {
      await updateRuleStatusMutation.mutateAsync({
        id: rule.id,
        status: newStatus,
      })
      toast.success(`规则已${newStatus === 'enabled' ? '启用' : '禁用'}`)
    } catch (error: any) {
      toast.error('状态更新失败: ' + error.message)
    }
  }

  // 处理手动触发
  const handleTrigger = async (rule: Rule) => {
    try {
      const execution = await triggerRuleMutation.mutateAsync({
        id: rule.id,
        data: {
          // 这里可以传入测试数据
          timestamp: new Date().toISOString(),
        },
      })
      toast.success(`规则已触发，执行ID: ${execution.id}`)
    } catch (error: any) {
      toast.error('触发失败: ' + error.message)
    }
  }

  return (
    <>
      <Header>
        <Search />
        <div className='ms-auto flex items-center space-x-4'>
          <ThemeSwitch />
          <ConfigDrawer />
          <ProfileDropdown />
        </div>
      </Header>

      <Main fixed>
        <div className='flex h-full flex-col'>
          {/* 页面标题栏 */}
          <div className='mb-6 flex items-center justify-between'>
            <div>
              <h1 className='text-3xl font-bold tracking-tight'>规则引擎</h1>
              <p className='mt-1 text-muted-foreground'>
                管理设备联动、数据流转和告警规则
              </p>
            </div>
            <Button onClick={() => setCreateDialogOpen(true)}>
              <Plus className='mr-2 h-4 w-4' />
              创建规则
            </Button>
          </div>

          {/* 统计卡片 */}
          <div className='mb-6 grid gap-4 md:grid-cols-4'>
            <div className='rounded-lg border bg-card p-4'>
              <div className='flex items-center justify-between'>
                <p className='text-sm font-medium text-muted-foreground'>
                  总规则数
                </p>
              </div>
              <div className='mt-2'>
                <p className='text-2xl font-bold'>{rulesData?.total || 0}</p>
              </div>
            </div>
            <div className='rounded-lg border bg-card p-4'>
              <div className='flex items-center justify-between'>
                <p className='text-sm font-medium text-muted-foreground'>
                  已启用
                </p>
              </div>
              <div className='mt-2'>
                <p className='text-2xl font-bold text-green-600'>
                  {rulesData?.items.filter((r: Rule) => r.status === 'enabled')
                    .length || 0}
                </p>
              </div>
            </div>
            <div className='rounded-lg border bg-card p-4'>
              <div className='flex items-center justify-between'>
                <p className='text-sm font-medium text-muted-foreground'>
                  已禁用
                </p>
              </div>
              <div className='mt-2'>
                <p className='text-2xl font-bold text-gray-600'>
                  {rulesData?.items.filter((r: Rule) => r.status === 'disabled')
                    .length || 0}
                </p>
              </div>
            </div>
            <div className='rounded-lg border bg-card p-4'>
              <div className='flex items-center justify-between'>
                <p className='text-sm font-medium text-muted-foreground'>
                  草稿
                </p>
              </div>
              <div className='mt-2'>
                <p className='text-2xl font-bold text-blue-600'>
                  {rulesData?.items.filter((r: Rule) => r.status === 'draft')
                    .length || 0}
                </p>
              </div>
            </div>
          </div>

          {/* 规则列表表格 */}
          <div className='flex-1'>
            <RulesTable
              data={rulesData?.items || []}
              isLoading={isLoading}
              onEdit={handleEdit}
              onView={handleView}
              onDelete={handleDelete}
              onToggleStatus={handleToggleStatus}
              onTrigger={handleTrigger}
              onBulkDelete={handleBulkDelete}
            />
          </div>

          {/* 对话框 */}
          <CreateRuleDialog
            open={createDialogOpen}
            onOpenChange={setCreateDialogOpen}
          />

          <EditRuleDialog
            rule={selectedRule}
            open={editDialogOpen}
            onOpenChange={setEditDialogOpen}
          />

          <ViewRuleDetailsDialog
            rule={selectedRule}
            open={viewDialogOpen}
            onOpenChange={setViewDialogOpen}
            onEdit={handleEdit}
            onTrigger={handleTrigger}
          />

          <DeleteRuleDialog
            rule={selectedRule}
            open={deleteDialogOpen}
            onOpenChange={setDeleteDialogOpen}
          />
        </div>
      </Main>
    </>
  )
}
