import { useParams } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ArrowLeft, Loader2, AlertCircle, RotateCw } from 'lucide-react'
import { useNavigate } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import {
  useOTAPackageDetail,
  useUpgradeStatistics,
  useDeviceDeployments,
  useUpgradeBatches,
} from '../api/detail-queries'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

/**
 * Type definitions for API response data
 */
interface UpgradeBatch {
  batchId: string
  batchName: string
  batchType: string
  upgradeStrategy: string
  status: string
  targetDeviceCount: number
  createdAt: Date
}

interface DeviceDeployment {
  deviceId: string
  deviceKey: string
  deviceName: string
  productId: string
  productKey: string
  currentVersion: string
  upgradeBatchId: string
  status: string
  lastStatusChangeTime: string | number
  createdAt: Date
}

/**
 * OTA Package Detail Page
 * Displays comprehensive details of an OTA package including:
 * - Package metadata
 * - Upgrade statistics (target devices, success/fail counts)
 * - Device deployment status
 * - Upgrade batches
 *
 * Responsive design:
 * - Mobile (< 640px): Single column layout, stacked cards
 * - Tablet (640px - 1024px): 2-column layout for statistics
 * - Desktop (> 1024px): 4-column layout for statistics
 */
export function OTAPackageDetailPage() {
  const { t } = useTranslation('operationsMonitoring')
  const navigate = useNavigate()
  const params = useParams({
    from: '/_authenticated/operations-monitoring/ota/packages/$id/',
  })

  const [currentPage, setCurrentPage] = useState(1)
  const [pageSize] = useState(100)
  const [statusFilter, setStatusFilter] = useState<string>()
  const [lastUpdated, setLastUpdated] = useState<Date>(new Date())

  // Data fetching hooks
  const packageQuery = useOTAPackageDetail(params.id)
  const packageName = packageQuery.data?.packageName

  const statisticsQuery = useUpgradeStatistics(packageName || '')
  const deploymentsQuery = useDeviceDeployments(
    packageName || '',
    currentPage,
    pageSize,
    statusFilter
  )
  const batchesQuery = useUpgradeBatches(packageName || '')

  useEffect(() => {
    if (
      !statisticsQuery.isLoading &&
      !deploymentsQuery.isLoading &&
      !batchesQuery.isLoading
    ) {
      setLastUpdated(new Date())
    }
  }, [statisticsQuery.isLoading, deploymentsQuery.isLoading, batchesQuery.isLoading])

  const formatLastUpdated = () => {
    const now = new Date()
    const diff = now.getTime() - lastUpdated.getTime()
    const seconds = Math.floor(diff / 1000)
    const minutes = Math.floor(seconds / 60)
    const hours = Math.floor(minutes / 60)

    if (seconds < 60) return 'Just now'
    if (minutes < 60) return `${minutes}m ago`
    if (hours < 24) return `${hours}h ago`
    return lastUpdated.toLocaleTimeString()
  }

  const handleBack = () => {
    navigate({
      to: '/operations-monitoring/ota/packages',
    })
  }

  // Error component
  const ErrorAlert = ({
    title,
    message,
    onRetry
  }: {
    title: string
    message: string
    onRetry?: () => void
  }) => (
    <div className='rounded-lg border border-red-200 bg-red-50 p-4 flex items-start gap-3'>
      <AlertCircle className='h-5 w-5 text-red-600 flex-shrink-0 mt-0.5' />
      <div className='flex-1'>
        <h3 className='font-medium text-red-900'>{title}</h3>
        <p className='text-sm text-red-700 mt-1'>{message}</p>
      </div>
      {onRetry && (
        <Button
          variant='outline'
          size='sm'
          onClick={onRetry}
          className='flex-shrink-0'
        >
          <RotateCw className='h-4 w-4 mr-1' />
          Retry
        </Button>
      )}
    </div>
  )

  // Loading skeleton
  const SkeletonCard = () => (
    <div className='rounded-lg border p-4 bg-card animate-pulse'>
      <div className='h-4 bg-muted rounded w-24 mb-3' />
      <div className='h-8 bg-muted rounded w-16' />
    </div>
  )

  // Main error state
  if (packageQuery.isError) {
    return (
      <div className='flex flex-1 flex-col gap-4 p-6'>
        <Button
          variant='ghost'
          size='sm'
          onClick={handleBack}
          className='mb-2 w-fit'
        >
          <ArrowLeft className='mr-2 h-4 w-4' />
          {t('common.back')}
        </Button>
        <ErrorAlert
          title='Failed to load package details'
          message='The OTA package could not be loaded. Please check the package ID and try again.'
          onRetry={() => { void (packageQuery as any).refetch?.() }}
        />
      </div>
    )
  }

  const pkg = packageQuery.data
  const stats = statisticsQuery.data

  return (
    <div className='flex flex-1 flex-col gap-4 px-2 sm:px-4 py-2'>
      {/* Header with back button */}
      <div className='flex items-center justify-between gap-2'>
        <Button
          variant='ghost'
          size='sm'
          onClick={handleBack}
          className='mb-2'
        >
          <ArrowLeft className='mr-2 h-4 w-4' />
          {t('common.back')}
        </Button>
      </div>

      {/* Main content */}
      <div className='flex flex-col gap-6'>
        {/* Title section */}
        <div>
          <h1 className='text-2xl sm:text-3xl font-bold tracking-tight'>
            {packageQuery.isLoading ? (
              <span className='animate-pulse'>Loading...</span>
            ) : (
              pkg?.packageName || 'Unknown Package'
            )}
          </h1>
          {!packageQuery.isLoading && (
            <p className='text-sm sm:text-base text-muted-foreground mt-1'>
              Version: {pkg?.version} • Type:{' '}
              <Badge className='ml-1 text-xs sm:text-sm'>{pkg?.packageType}</Badge>
            </p>
          )}
        </div>

        {/* Last Updated info */}
        <div className='text-xs sm:text-sm text-muted-foreground px-2 py-1 bg-muted rounded'>
          Last updated: {formatLastUpdated()}
        </div>

        {/* Statistics cards section */}
        <div className='grid gap-3 sm:gap-4 grid-cols-1 sm:grid-cols-2 lg:grid-cols-4'>
          <div className='rounded-lg border p-3 sm:p-4 bg-card'>
            <div className='text-xs sm:text-sm font-medium text-muted-foreground'>
              Target Devices
            </div>
            <div className='text-xl sm:text-2xl font-bold mt-2'>
              {statisticsQuery.isLoading ? (
                <Loader2 className='h-4 w-4 animate-spin' />
              ) : statisticsQuery.isError ? (
                <span className='text-red-600 text-sm'>Error</span>
              ) : (
                stats?.totalTargetDevices ?? '-'
              )}
            </div>
          </div>
          <div className='rounded-lg border p-3 sm:p-4 bg-card'>
            <div className='text-xs sm:text-sm font-medium text-muted-foreground'>
              Successful
            </div>
            <div className='text-xl sm:text-2xl font-bold text-green-600 mt-2'>
              {statisticsQuery.isLoading ? (
                <Loader2 className='h-4 w-4 animate-spin' />
              ) : statisticsQuery.isError ? (
                <span className='text-red-600 text-sm'>Error</span>
              ) : (
                stats?.successfulUpgrades ?? '-'
              )}
            </div>
          </div>
          <div className='rounded-lg border p-3 sm:p-4 bg-card'>
            <div className='text-xs sm:text-sm font-medium text-muted-foreground'>
              Failed
            </div>
            <div className='text-xl sm:text-2xl font-bold text-red-600 mt-2'>
              {statisticsQuery.isLoading ? (
                <Loader2 className='h-4 w-4 animate-spin' />
              ) : statisticsQuery.isError ? (
                <span className='text-red-600 text-sm'>Error</span>
              ) : (
                stats?.failedUpgrades ?? '-'
              )}
            </div>
          </div>
          <div className='rounded-lg border p-3 sm:p-4 bg-card'>
            <div className='text-xs sm:text-sm font-medium text-muted-foreground'>
              In Progress
            </div>
            <div className='text-xl sm:text-2xl font-bold text-blue-600 mt-2'>
              {statisticsQuery.isLoading ? (
                <Loader2 className='h-4 w-4 animate-spin' />
              ) : statisticsQuery.isError ? (
                <span className='text-red-600 text-sm'>Error</span>
              ) : (
                stats?.inProgressUpgrades ?? '-'
              )}
            </div>
          </div>
        </div>

        {/* Tabs section */}
        <Tabs defaultValue='batches' className='w-full'>
          <TabsList>
            <TabsTrigger value='batches'>Batch Management</TabsTrigger>
            <TabsTrigger value='devices'>Device List</TabsTrigger>
            <TabsTrigger value='info'>Package Info</TabsTrigger>
          </TabsList>

          {/* Batch Management Tab */}
          <TabsContent value='batches' className='space-y-4'>
            {batchesQuery.isError && (
              <ErrorAlert
                title='Failed to load batches'
                message='Could not load upgrade batches. Please try again.'
                onRetry={() => { void (batchesQuery as any).refetch?.() }}
              />
            )}
            <div className='rounded-lg border bg-card'>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Batch ID</TableHead>
                    <TableHead>Batch Name</TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead>Strategy</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Devices</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {batchesQuery.isLoading ? (
                    <TableRow>
                      <TableCell colSpan={6} className='text-center py-8'>
                        <Loader2 className='h-4 w-4 animate-spin mx-auto' />
                      </TableCell>
                    </TableRow>
                  ) : batchesQuery.data?.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={6} className='text-center py-8 text-muted-foreground'>
                        No batches found
                      </TableCell>
                    </TableRow>
                  ) : (
                    batchesQuery.data?.map((batch: UpgradeBatch) => (
                      <TableRow key={batch.batchId}>
                        <TableCell className='font-mono text-sm'>
                          {batch.batchId}
                        </TableCell>
                        <TableCell>{batch.batchName}</TableCell>
                        <TableCell>
                          <Badge variant='outline' className='capitalize'>
                            {batch.batchType}
                          </Badge>
                        </TableCell>
                        <TableCell className='capitalize'>
                          {batch.upgradeStrategy}
                        </TableCell>
                        <TableCell>
                          <Badge
                            className='capitalize'
                            variant={
                              batch.status === 'completed'
                                ? 'default'
                                : batch.status === 'in_progress'
                                  ? 'secondary'
                                  : 'outline'
                            }
                          >
                            {batch.status}
                          </Badge>
                        </TableCell>
                        <TableCell>{batch.targetDeviceCount}</TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>
          </TabsContent>

          {/* Device List Tab */}
          <TabsContent value='devices' className='space-y-4'>
            {deploymentsQuery.isError && (
              <ErrorAlert
                title='Failed to load devices'
                message='Could not load device deployment status. Please try again.'
                onRetry={() => { void (deploymentsQuery as any).refetch?.() }}
              />
            )}
            <div className='flex flex-col sm:flex-row gap-2 mb-4'>
              <select
                value={statusFilter || ''}
                onChange={(e) => {
                  setStatusFilter(e.target.value || undefined)
                  setCurrentPage(1)
                }}
                className='px-3 py-2 border rounded-md text-sm flex-1 sm:flex-none'
              >
                <option value=''>All Status</option>
                <option value='pending'>Pending</option>
                <option value='in_progress'>In Progress</option>
                <option value='success'>Success</option>
                <option value='failed'>Failed</option>
                <option value='cancelled'>Cancelled</option>
              </select>
            </div>

            <div className='rounded-lg border bg-card overflow-x-auto'>
              <Table className='text-sm sm:text-base'>
                <TableHeader>
                  <TableRow>
                    <TableHead className='text-xs sm:text-sm'>Device</TableHead>
                    <TableHead className='text-xs sm:text-sm'>Product</TableHead>
                    <TableHead className='text-xs sm:text-sm'>Version</TableHead>
                    <TableHead className='text-xs sm:text-sm'>Batch ID</TableHead>
                    <TableHead className='text-xs sm:text-sm'>Status</TableHead>
                    <TableHead className='text-xs sm:text-sm'>Last Update</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {deploymentsQuery.isLoading ? (
                    <TableRow>
                      <TableCell colSpan={6} className='text-center py-8'>
                        <Loader2 className='h-4 w-4 animate-spin mx-auto' />
                      </TableCell>
                    </TableRow>
                  ) : deploymentsQuery.data?.deployments?.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={6} className='text-center py-8 text-muted-foreground text-xs sm:text-sm'>
                        No devices found
                      </TableCell>
                    </TableRow>
                  ) : (
                    deploymentsQuery.data?.deployments?.map((d: DeviceDeployment) => (
                      <TableRow key={d.deviceId}>
                        <TableCell className='font-mono text-xs sm:text-sm'>
                          {d.deviceName}
                        </TableCell>
                        <TableCell className='text-xs sm:text-sm'>{d.productKey}</TableCell>
                        <TableCell className='text-xs sm:text-sm'>{d.currentVersion}</TableCell>
                        <TableCell className='font-mono text-xs'>
                          {d.upgradeBatchId}
                        </TableCell>
                        <TableCell>
                          <Badge
                            className='capitalize text-xs'
                            variant={
                              d.status === 'success'
                                ? 'default'
                                : d.status === 'failed'
                                  ? 'destructive'
                                  : 'outline'
                            }
                          >
                            {d.status}
                          </Badge>
                        </TableCell>
                        <TableCell className='text-xs sm:text-sm text-muted-foreground'>
                          {new Date(d.lastStatusChangeTime).toLocaleString()}
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>

            {/* Pagination */}
            <div className='flex flex-col sm:flex-row justify-between items-center gap-4'>
              <div className='text-xs sm:text-sm text-muted-foreground'>
                Page {currentPage} of{' '}
                {Math.ceil((deploymentsQuery.data?.total || 0) / pageSize)}
              </div>
              <div className='flex gap-2'>
                <Button
                  variant='outline'
                  size='sm'
                  onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                  disabled={currentPage === 1}
                  className='text-xs'
                >
                  Previous
                </Button>
                <Button
                  variant='outline'
                  size='sm'
                  onClick={() => setCurrentPage((p) => p + 1)}
                  disabled={
                    currentPage >=
                    Math.ceil((deploymentsQuery.data?.total || 0) / pageSize)
                  }
                  className='text-xs'
                >
                  Next
                </Button>
              </div>
            </div>
          </TabsContent>

          {/* Package Info Tab */}
          <TabsContent value='info' className='space-y-4'>
            {packageQuery.isError && (
              <ErrorAlert
                title='Failed to load package information'
                message='Could not load the package details. Please try again.'
                onRetry={() => { void (packageQuery as any).refetch?.() }}
              />
            )}
            {packageQuery.isLoading ? (
              <div className='rounded-lg border p-6 bg-card space-y-4'>
                <div className='grid grid-cols-2 gap-4'>
                  {[1, 2, 3, 4].map((i) => (
                    <SkeletonCard key={i} />
                  ))}
                </div>
              </div>
            ) : (
              <div className='space-y-4'>
                {/* Basic Info Section */}
                <div className='rounded-lg border bg-card overflow-hidden'>
                  <div className='bg-muted px-6 py-3 border-b'>
                    <h3 className='font-semibold text-sm'>Basic Information</h3>
                  </div>
                  <div className='p-6'>
                    <div className='grid grid-cols-1 md:grid-cols-2 gap-6'>
                      <div>
                        <label className='text-xs uppercase tracking-wide font-semibold text-muted-foreground mb-2 block'>
                          Package Name
                        </label>
                        <p className='text-base font-mono break-all'>{pkg?.packageName}</p>
                      </div>
                      <div>
                        <label className='text-xs uppercase tracking-wide font-semibold text-muted-foreground mb-2 block'>
                          Version
                        </label>
                        <p className='text-base font-mono'>{pkg?.version}</p>
                      </div>
                      <div>
                        <label className='text-xs uppercase tracking-wide font-semibold text-muted-foreground mb-2 block'>
                          Package Type
                        </label>
                        <div className='flex items-center gap-2'>
                          <Badge variant='secondary' className='capitalize'>{pkg?.packageType}</Badge>
                        </div>
                      </div>
                      <div>
                        <label className='text-xs uppercase tracking-wide font-semibold text-muted-foreground mb-2 block'>
                          Status
                        </label>
                        <Badge className='capitalize'>{pkg?.status}</Badge>
                      </div>
                    </div>
                  </div>
                </div>

                {/* File Information Section */}
                <div className='rounded-lg border bg-card overflow-hidden'>
                  <div className='bg-muted px-6 py-3 border-b'>
                    <h3 className='font-semibold text-sm'>File Information</h3>
                  </div>
                  <div className='p-6 space-y-4'>
                    <div>
                      <label className='text-xs uppercase tracking-wide font-semibold text-muted-foreground mb-2 block'>
                        File URL
                      </label>
                      <p className='text-sm break-all bg-muted p-3 rounded font-mono'>
                        {pkg?.fileUrl || '-'}
                      </p>
                    </div>
                    <div>
                      <label className='text-xs uppercase tracking-wide font-semibold text-muted-foreground mb-2 block'>
                        Checksum (SHA256)
                      </label>
                      <p className='text-sm break-all bg-muted p-3 rounded font-mono'>
                        {pkg?.checksum || '-'}
                      </p>
                    </div>
                  </div>
                </div>

                {/* Description Section */}
                <div className='rounded-lg border bg-card overflow-hidden'>
                  <div className='bg-muted px-6 py-3 border-b'>
                    <h3 className='font-semibold text-sm'>Description & Notes</h3>
                  </div>
                  <div className='p-6 space-y-4'>
                    <div>
                      <label className='text-xs uppercase tracking-wide font-semibold text-muted-foreground mb-2 block'>
                        Description
                      </label>
                      <p className='text-sm leading-relaxed'>
                        {pkg?.description || <span className='text-muted-foreground italic'>No description provided</span>}
                      </p>
                    </div>
                    <div>
                      <label className='text-xs uppercase tracking-wide font-semibold text-muted-foreground mb-2 block'>
                        Release Notes
                      </label>
                      <p className='text-sm whitespace-pre-wrap leading-relaxed bg-muted p-3 rounded max-h-48 overflow-y-auto'>
                        {pkg?.releaseNotes || <span className='text-muted-foreground italic'>No release notes provided</span>}
                      </p>
                    </div>
                  </div>
                </div>

                {/* Actions */}
                <div className='flex gap-2 pt-2'>
                  <Button variant='outline' size='sm'>
                    Edit Package
                  </Button>
                  <Button variant='outline' size='sm'>
                    Download
                  </Button>
                </div>
              </div>
            )}
          </TabsContent>
        </Tabs>
      </div>
    </div>
  )
}
