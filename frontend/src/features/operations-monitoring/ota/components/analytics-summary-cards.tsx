import { Package, Upload, CheckCircle2, XCircle } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import type { OTAAnalyticsSummary } from '../data/analytics-schema'

interface AnalyticsSummaryCardsProps {
  summary: OTAAnalyticsSummary
  isLoading?: boolean
}

export function AnalyticsSummaryCards({
  summary,
  isLoading,
}: AnalyticsSummaryCardsProps) {
  const { t } = useTranslation('operationsMonitoring')

  const cards = [
    {
      title: t('ota.analytics.summary.totalPackages'),
      value: summary.totalPackages,
      icon: Package,
      color: 'text-blue-600',
      bgColor: 'bg-blue-100',
    },
    {
      title: t('ota.analytics.summary.activeDeployments'),
      value: summary.activeDeployments,
      icon: Upload,
      color: 'text-amber-600',
      bgColor: 'bg-amber-100',
    },
    {
      title: t('ota.analytics.summary.successRate'),
      value: `${summary.successRate}%`,
      icon: CheckCircle2,
      color: 'text-green-600',
      bgColor: 'bg-green-100',
    },
    {
      title: t('ota.analytics.summary.failedDeployments'),
      value: summary.failedDeployments,
      icon: XCircle,
      color: 'text-red-600',
      bgColor: 'bg-red-100',
    },
  ]

  return (
    <div className='grid gap-4 md:grid-cols-2 lg:grid-cols-4'>
      {cards.map((card) => {
        const Icon = card.icon
        return (
          <Card key={card.title}>
            <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
              <CardTitle className='text-sm font-medium'>
                {card.title}
              </CardTitle>
              <div className={`rounded-full p-2 ${card.bgColor}`}>
                <Icon className={`h-4 w-4 ${card.color}`} />
              </div>
            </CardHeader>
            <CardContent>
              {isLoading ? (
                <div className='h-8 w-20 animate-pulse rounded bg-muted' />
              ) : (
                <div className='text-2xl font-bold'>{card.value}</div>
              )}
            </CardContent>
          </Card>
        )
      })}
    </div>
  )
}
