import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'

export type UpgradeStatus = 'pending' | 'in_progress' | 'success' | 'failed' | 'cancelled'
export type VerificationStatus = 'completed' | 'in_progress' | 'failed' | 'pending'

interface StatusBadgeProps {
  status: UpgradeStatus
  className?: string
}

/**
 * StatusBadge component for displaying upgrade status with color coding
 * - pending: gray outline
 * - in_progress: blue secondary
 * - success: green default
 * - failed: red destructive
 * - cancelled: gray outline
 */
export function StatusBadge({ status, className }: StatusBadgeProps) {
  const getVariant = (status: UpgradeStatus) => {
    switch (status) {
      case 'success':
        return 'default'
      case 'failed':
        return 'destructive'
      case 'in_progress':
        return 'secondary'
      case 'pending':
      case 'cancelled':
      default:
        return 'outline'
    }
  }

  return (
    <Badge
      variant={getVariant(status)}
      className={cn('capitalize text-xs', className)}
    >
      {status}
    </Badge>
  )
}

interface VerificationStatusBadgeProps {
  status: VerificationStatus
  progress?: number
  className?: string
}

/**
 * VerificationStatusBadge component for displaying verification status
 * - pending: gray outline
 * - in_progress: blue secondary with progress percentage
 * - completed: green default
 * - failed: red destructive
 */
export function VerificationStatusBadge({
  status,
  progress,
  className,
}: VerificationStatusBadgeProps) {
  const getVariant = (status: VerificationStatus) => {
    switch (status) {
      case 'completed':
        return 'default'
      case 'failed':
        return 'destructive'
      case 'in_progress':
        return 'secondary'
      case 'pending':
      default:
        return 'outline'
    }
  }

  return (
    <Badge
      variant={getVariant(status)}
      className={cn('capitalize text-xs', className)}
    >
      {status}
      {status === 'in_progress' && progress !== undefined && (
        <span className='ml-1'>({progress}%)</span>
      )}
    </Badge>
  )
}

// Batch status badge variants
export type BatchStatus = 'pending' | 'in_progress' | 'paused' | 'completed' | 'failed'

interface BatchStatusBadgeProps {
  status: BatchStatus
  className?: string
}

/**
 * BatchStatusBadge component for displaying upgrade batch status
 */
export function BatchStatusBadge({ status, className }: BatchStatusBadgeProps) {
  const getVariant = (status: BatchStatus) => {
    switch (status) {
      case 'completed':
        return 'default'
      case 'failed':
        return 'destructive'
      case 'in_progress':
        return 'secondary'
      case 'pending':
      case 'paused':
      default:
        return 'outline'
    }
  }

  return (
    <Badge
      variant={getVariant(status)}
      className={cn('capitalize text-xs', className)}
    >
      {status}
    </Badge>
  )
}
