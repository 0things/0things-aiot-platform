import { toast } from 'sonner'
import { useCreateRule } from '../../api/queries'
import type { CreateRuleFormData } from '../../data/schema'
import { RuleFormDialog } from './rule-form-dialog'

interface CreateRuleDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function CreateRuleDialog({
  open,
  onOpenChange,
}: CreateRuleDialogProps) {
  const createRuleMutation = useCreateRule()

  const handleSubmit = async (data: CreateRuleFormData) => {
    try {
      await createRuleMutation.mutateAsync(data)
      toast.success('规则创建成功')
      onOpenChange(false)
    } catch (error: any) {
      toast.error('创建失败: ' + error.message)
    }
  }

  return (
    <RuleFormDialog
      open={open}
      onOpenChange={onOpenChange}
      mode='create'
      pending={createRuleMutation.isPending}
      onSubmit={handleSubmit}
    />
  )
}
