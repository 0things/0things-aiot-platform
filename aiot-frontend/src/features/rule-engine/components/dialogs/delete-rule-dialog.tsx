import { toast } from 'sonner'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { useDeleteRule } from '../../api/queries'
import type { Rule } from '../../data/schema'

interface DeleteRuleDialogProps {
  rule: Rule | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function DeleteRuleDialog({
  rule,
  open,
  onOpenChange,
}: DeleteRuleDialogProps) {
  const deleteRuleMutation = useDeleteRule()

  const handleDelete = async () => {
    if (!rule) return

    try {
      await deleteRuleMutation.mutateAsync(rule.id)
      toast.success('规则已删除')
      onOpenChange(false)
    } catch (error: any) {
      toast.error('删除失败: ' + error.message)
    }
  }

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>确认删除规则？</AlertDialogTitle>
          <AlertDialogDescription>
            你正在删除规则 <strong>{rule?.name}</strong>
            。此操作不可撤销，该规则的所有配置和执行历史将被永久删除。
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>取消</AlertDialogCancel>
          <AlertDialogAction
            onClick={handleDelete}
            disabled={deleteRuleMutation.isPending}
            className='text-destructive-foreground bg-destructive hover:bg-destructive/90'
          >
            {deleteRuleMutation.isPending ? '删除中...' : '确认删除'}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
