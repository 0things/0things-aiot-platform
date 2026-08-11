import { formatDistanceToNow } from 'date-fns'
import { CheckCircle2, XCircle, Loader2, Clock } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import type { OTARecentActivity } from '../data/analytics-schema'

interface RecentActivityTableProps {
  data: OTARecentActivity[]
  isLoading?: boolean
}

const statusConfig = {
  success: {
    icon: CheckCircle2,
    variant: 'default' as const,
    className: 'bg-green-100 text-green-800 hover:bg-green-100',
  },
  failed: {
    icon: XCircle,
    variant: 'destructive' as const,
    className: 'bg-red-100 text-red-800 hover:bg-red-100',
  },
  in_progress: {
    icon: Loader2,
    variant: 'secondary' as const,
    className: 'bg-amber-100 text-amber-800 hover:bg-amber-100',
  },
  pending: {
    icon: Clock,
    variant: 'outline' as const,
    className: 'bg-gray-100 text-gray-800 hover:bg-gray-100',
  },
}

const actionConfig = {
  created: 'default' as const,
  deployed: 'secondary' as const,
  completed: 'default' as const,
  failed: 'destructive' as const,
}

export function RecentActivityTable({
  data,
  isLoading,
}: RecentActivityTableProps) {
  const { t } = useTranslation('operationsMonitoring')

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t('ota.analytics.recentActivity.title')}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className='overflow-auto rounded-md border'>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>
                  {t('ota.analytics.recentActivity.packageName')}
                </TableHead>
                <TableHead>
                  {t('ota.analytics.recentActivity.version')}
                </TableHead>
                <TableHead>
                  {t('ota.analytics.recentActivity.action')}
                </TableHead>
                <TableHead>
                  {t('ota.analytics.recentActivity.productName')}
                </TableHead>
                <TableHead>
                  {t('ota.analytics.recentActivity.status')}
                </TableHead>
                <TableHead>{t('ota.analytics.recentActivity.time')}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading ? (
                Array.from({ length: 5 }).map((_, index) => (
                  <TableRow key={index}>
                    {Array.from({ length: 6 }).map((_, cellIndex) => (
                      <TableCell key={cellIndex}>
                        <div className='h-5 w-full animate-pulse rounded bg-muted' />
                      </TableCell>
                    ))}
                  </TableRow>
                ))
              ) : data.length > 0 ? (
                data.map((activity) => {
                  const StatusIcon = statusConfig[activity.status].icon
                  return (
                    <TableRow key={activity.id}>
                      <TableCell className='font-medium'>
                        {activity.packageName}
                      </TableCell>
                      <TableCell>{activity.version}</TableCell>
                      <TableCell>
                        <Badge variant={actionConfig[activity.action]}>
                          {t(`ota.analytics.actions.${activity.action}`)}
                        </Badge>
                      </TableCell>
                      <TableCell>{activity.productName}</TableCell>
                      <TableCell>
                        <div className='flex items-center gap-2'>
                          <Badge
                            variant={statusConfig[activity.status].variant}
                            className={statusConfig[activity.status].className}
                          >
                            <StatusIcon className='mr-1 h-3 w-3' />
                            {t(`ota.analytics.statuses.${activity.status}`)}
                          </Badge>
                        </div>
                      </TableCell>
                      <TableCell className='text-muted-foreground'>
                        {formatDistanceToNow(new Date(activity.timestamp), {
                          addSuffix: true,
                        })}
                      </TableCell>
                    </TableRow>
                  )
                })
              ) : (
                <TableRow>
                  <TableCell colSpan={6} className='h-24 text-center'>
                    {t('ota.analytics.recentActivity.noData')}
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  )
}
