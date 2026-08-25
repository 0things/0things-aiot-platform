import { useTranslation } from 'react-i18next'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
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

type StatusKey = OTARecentActivity['status']

const STATUS_VARIANT: Record<
  StatusKey,
  'default' | 'secondary' | 'destructive' | 'outline'
> = {
  success: 'default',
  failed: 'destructive',
  in_progress: 'secondary',
  pending: 'outline',
}

export function RecentActivityTable({
  data,
  isLoading,
}: RecentActivityTableProps) {
  const { t } = useTranslation('ota')

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t('analytics.recentActivity.title')}</CardTitle>
      </CardHeader>
      <CardContent className='px-0 pt-0'>
        {isLoading ? (
          <div className='space-y-2 px-6 pb-2'>
            {Array.from({ length: 5 }).map((_, i) => (
              <Skeleton key={i} className='h-10 w-full' />
            ))}
          </div>
        ) : data.length === 0 ? (
          <div className='px-6 pt-2 pb-6 text-sm text-muted-foreground'>
            {t('analytics.recentActivity.noData')}
          </div>
        ) : (
          <Table>
            <TableHeader>
              <TableRow className='hover:bg-transparent'>
                <TableHead className='text-xs tracking-wide text-muted-foreground uppercase'>
                  {t('analytics.recentActivity.packageName')}
                </TableHead>
                <TableHead className='text-xs tracking-wide text-muted-foreground uppercase'>
                  {t('analytics.recentActivity.version')}
                </TableHead>
                <TableHead className='text-xs tracking-wide text-muted-foreground uppercase'>
                  {t('analytics.recentActivity.action')}
                </TableHead>
                <TableHead className='text-xs tracking-wide text-muted-foreground uppercase'>
                  {t('analytics.recentActivity.productName')}
                </TableHead>
                <TableHead className='text-xs tracking-wide text-muted-foreground uppercase'>
                  {t('analytics.recentActivity.status')}
                </TableHead>
                <TableHead className='text-right text-xs tracking-wide text-muted-foreground uppercase'>
                  {t('analytics.recentActivity.time')}
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {data.map((item) => (
                <TableRow key={item.id}>
                  <TableCell className='font-medium'>
                    {item.packageName}
                  </TableCell>
                  <TableCell className='font-mono text-xs text-muted-foreground'>
                    v{item.version}
                  </TableCell>
                  <TableCell>
                    <Badge variant='outline' className='font-normal'>
                      {t(`analytics.actions.${item.action}` as const)}
                    </Badge>
                  </TableCell>
                  <TableCell className='text-muted-foreground'>
                    {item.productName}
                  </TableCell>
                  <TableCell>
                    <Badge variant={STATUS_VARIANT[item.status as StatusKey]}>
                      {t(`analytics.statuses.${item.status}` as const)}
                    </Badge>
                  </TableCell>
                  <TableCell className='text-right font-mono text-xs text-muted-foreground tabular-nums'>
                    {item.timestamp || '-'}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </CardContent>
    </Card>
  )
}
