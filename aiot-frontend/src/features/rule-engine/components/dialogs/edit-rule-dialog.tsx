import { toast } from 'sonner'
import { useUpdateRule } from '../../api/queries'
import type { CreateRuleFormData, Rule } from '../../data/schema'
import { RuleFormDialog } from './rule-form-dialog'

interface EditRuleDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  rule: Rule | null
}

export function EditRuleDialog({
  open,
  onOpenChange,
  rule,
}: EditRuleDialogProps) {
  const updateRuleMutation = useUpdateRule()

  const handleSubmit = async (data: CreateRuleFormData) => {
    if (!rule) return
    try {
      await updateRuleMutation.mutateAsync({
        id: rule.id,
        data,
      })
      toast.success('规则更新成功')
      onOpenChange(false)
    } catch (error: any) {
      toast.error('更新失败: ' + error.message)
    }
  }

  return (
    <RuleFormDialog
      open={open}
      onOpenChange={onOpenChange}
      mode='edit'
      initialRule={rule}
      pending={updateRuleMutation.isPending}
      onSubmit={handleSubmit}
    />
  )
}
