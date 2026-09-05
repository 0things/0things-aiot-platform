import { useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  ChevronLeft,
  ChevronRight,
  Inbox,
  Plus,
  RefreshCw,
  Search,
  Trash2,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import {
  deleteDeviceGroupsGroupUuidDevices,
  getProducts,
} from '@/api/generated'
import type {
  Device as DeviceV1Device,
  AiotBackendApiDeviceGroupV1DeviceGroup as DeviceGroupV1Group,
} from '@/api/generated/model'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { AddDevicesDialog } from './add-devices-dialog'

interface GroupDevicesTabProps {
  group: DeviceGroupV1Group
  devices: DeviceV1Device[]
  isLoading: boolean
  onRefresh: () => void
  total: number
  page: number
  pageSize: number
  onPageChange: (page: number) => void
  onPageSizeChange: (size: number) => void
  onSearch: (productKey: string, search: string) => void
}

export function GroupDevicesTab({
  group,
  devices,
  isLoading,
  onRefresh,
  total,
  page,
  pageSize,
  onPageChange,
  onPageSizeChange,
  onSearch,
}: GroupDevicesTabProps) {
  const { t } = useTranslation('deviceGroup')
  const queryClient = useQueryClient()
  const [productFilter, setProductFilter] = useState<string>('all')
  const [searchInput, setSearchInput] = useState<string>('')
  const [selectedKeys, setSelectedKeys] = useState<string[]>([])
  const [addDialogOpen, setAddDialogOpen] = useState<boolean>(false)
  const [removeConfirmOpen, setRemoveConfirmOpen] = useState<boolean>(false)
  const [targetKeysToRemove, setTargetKeysToRemove] = useState<string[]>([])

  // 获取产品列表供筛选
  const { data: productsData } = useQuery({
    queryKey: ['products-for-group-table-filter'],
    queryFn: () => getProducts({ page: 1, pageSize: 100 }),
  })

  const products = productsData?.data?.products || []
  // 移除设备 Mutation
  const removeMutation = useMutation({
    mutationFn: (keys: string[]) =>
      deleteDeviceGroupsGroupUuidDevices(group.groupUuid || '', {
        deviceKeys: keys,
      }),
    onSuccess: () => {
      toast.success(t('removeSuccess'))
      queryClient.invalidateQueries({
        queryKey: ['device-group', group.groupUuid, 'devices'],
      })
      queryClient.invalidateQueries({
        queryKey: ['device-group', group.groupUuid],
      })
      setSelectedKeys((prev) =>
        prev.filter((k) => !targetKeysToRemove.includes(k))
      )
      setTargetKeysToRemove([])
      setRemoveConfirmOpen(false)
    },
    onError: (err: unknown) => {
      toast.error(
        err instanceof Error ? err.message : 'Failed to remove devices'
      )
    },
  })

  const handleSelectAll = (checked: boolean) => {
    if (checked) {
      const keys = devices.map((d) => d.deviceKey).filter(Boolean) as string[]
      setSelectedKeys(Array.from(new Set([...selectedKeys, ...keys])))
    } else {
      const deviceSet = new Set(devices.map((d) => d.deviceKey))
      setSelectedKeys(selectedKeys.filter((k) => !deviceSet.has(k)))
    }
  }

  const handleToggleDevice = (deviceKey: string) => {
    setSelectedKeys((prev) =>
      prev.includes(deviceKey)
        ? prev.filter((k) => k !== deviceKey)
        : [...prev, deviceKey]
    )
  }

  const handleSingleRemove = (deviceKey: string) => {
    setTargetKeysToRemove([deviceKey])
    setRemoveConfirmOpen(true)
  }

  const handleBatchRemove = () => {
    if (selectedKeys.length === 0) return
    setTargetKeysToRemove(selectedKeys)
    setRemoveConfirmOpen(true)
  }

  const isAllSelected =
    devices.length > 0 &&
    devices.every((d) => d.deviceKey && selectedKeys.includes(d.deviceKey))

  const existingKeys = useMemo(
    () => devices.map((d) => d.deviceKey).filter(Boolean) as string[],
    [devices]
  )
  const totalPages = Math.max(1, Math.ceil(total / pageSize))

  return (
    <div className='space-y-4'>
      {/* 顶部操作与筛选栏 */}
      <div className='flex flex-wrap items-center justify-between gap-3'>
        <div className='flex flex-wrap items-center gap-3'>
          {group.type !== 'dynamic' && (
            <Button
              className='flex items-center gap-1.5'
              onClick={() => setAddDialogOpen(true)}
            >
              <Plus className='h-4 w-4' />
              <span>{t('addDevice')}</span>
            </Button>
          )}

          <Select value={productFilter} onValueChange={setProductFilter}>
            <SelectTrigger className='w-[160px]'>
              <SelectValue placeholder={t('allProducts')} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value='all'>{t('allProducts')}</SelectItem>
              {products.map((p) => (
                <SelectItem key={p.id} value={p.productKey || String(p.id)}>
                  {p.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <div className='relative w-[240px]'>
            <Search className='absolute top-2.5 left-2.5 h-4 w-4 text-muted-foreground' />
            <Input
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              placeholder={t('searchDevicePlaceholder')}
              className='pl-8'
            />
          </div>

          <Button
            onClick={() => {
              onSearch(
                productFilter === 'all' ? '' : productFilter,
                searchInput
              )
            }}
          >
            <Search className='mr-1.5 h-4 w-4' />
            {t('common:search')}
          </Button>
          <Button
            variant='outline'
            onClick={() => {
              setProductFilter('all')
              setSearchInput('')
              onSearch('', '')
            }}
          >
            {t('common:reset')}
          </Button>

          {selectedKeys.length > 0 && group.type !== 'dynamic' && (
            <Button
              variant='destructive'
              size='sm'
              className='flex items-center gap-1.5'
              onClick={handleBatchRemove}
            >
              <Trash2 className='h-3.5 w-3.5' />
              <span>
                {t('batchRemoveFromGroup')} ({selectedKeys.length})
              </span>
            </Button>
          )}
        </div>

        <Button
          variant='outline'
          size='icon'
          onClick={onRefresh}
          title={t('common:refresh')}
        >
          <RefreshCw className={`h-4 w-4 ${isLoading ? 'animate-spin' : ''}`} />
        </Button>
      </div>

      {/* 设备数据表格 */}
      <Card>
        <CardContent className='p-0'>
          <div className='overflow-hidden rounded-md border'>
            <Table>
              <TableHeader className='bg-muted/40'>
                <TableRow>
                  {group.type !== 'dynamic' && (
                    <TableHead className='w-12 text-center'>
                      <Checkbox
                        checked={isAllSelected}
                        onCheckedChange={(checked) =>
                          handleSelectAll(!!checked)
                        }
                        aria-label='Select all'
                      />
                    </TableHead>
                  )}
                  <TableHead>{t('deviceName')}</TableHead>
                  <TableHead>{t('productName')}</TableHead>
                  <TableHead>{t('nodeType')}</TableHead>
                  <TableHead>{t('statusAndState')}</TableHead>
                  <TableHead>{t('lastOnlineTime')}</TableHead>
                  {group.type !== 'dynamic' && (
                    <TableHead className='text-right'>
                      {t('common:actions')}
                    </TableHead>
                  )}
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell
                      colSpan={group.type !== 'dynamic' ? 7 : 5}
                      className='h-32 text-center text-muted-foreground'
                    >
                      {t('common:loading')}
                    </TableCell>
                  </TableRow>
                ) : devices.length === 0 ? (
                  <TableRow>
                    <TableCell
                      colSpan={group.type !== 'dynamic' ? 7 : 5}
                      className='h-64 text-center'
                    >
                      <div className='flex flex-col items-center justify-center space-y-3 py-6'>
                        <div className='rounded-full bg-muted/60 p-4 text-muted-foreground'>
                          <Inbox className='h-8 w-8' />
                        </div>
                        <p className='text-sm text-muted-foreground'>
                          {t('noDevicesFound')}
                        </p>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : (
                  devices.map((device) => {
                    const isChecked = device.deviceKey
                      ? selectedKeys.includes(device.deviceKey)
                      : false
                    return (
                      <TableRow
                        key={device.deviceKey}
                        className='hover:bg-muted/50'
                      >
                        {group.type !== 'dynamic' && (
                          <TableCell className='text-center'>
                            <Checkbox
                              checked={isChecked}
                              onCheckedChange={() =>
                                device.deviceKey &&
                                handleToggleDevice(device.deviceKey)
                              }
                            />
                          </TableCell>
                        )}
                        <TableCell>
                          <div className='flex flex-col'>
                            <span className='font-medium text-foreground'>
                              {device.name || '-'}
                            </span>
                            <span className='font-mono text-xs text-muted-foreground'>
                              {device.deviceKey}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell>
                          <span className='text-sm'>
                            {device.productName || '-'}
                          </span>
                        </TableCell>
                        <TableCell>
                          <Badge
                            variant='outline'
                            className='text-xs capitalize'
                          >
                            direct
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <div className='flex items-center gap-1.5'>
                            <Badge
                              variant={
                                device.state === 'online'
                                  ? 'default'
                                  : 'secondary'
                              }
                              className='text-xs capitalize'
                            >
                              {device.state === 'online'
                                ? t('statusOnline')
                                : t('statusOffline')}
                            </Badge>
                            <Badge
                              variant={
                                device.enabled ? 'outline' : 'destructive'
                              }
                              className='text-xs'
                            >
                              {device.enabled
                                ? t('enabledYes')
                                : t('enabledNo')}
                            </Badge>
                          </div>
                        </TableCell>
                        <TableCell className='text-xs text-muted-foreground'>
                          {device.lastOnlineTime || '-'}
                        </TableCell>
                        {group.type !== 'dynamic' && (
                          <TableCell className='text-right'>
                            <Button
                              variant='ghost'
                              size='sm'
                              className='text-destructive hover:text-destructive'
                              onClick={() =>
                                device.deviceKey &&
                                handleSingleRemove(device.deviceKey)
                              }
                            >
                              {t('removeFromGroup')}
                            </Button>
                          </TableCell>
                        )}
                      </TableRow>
                    )
                  })
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      <div className='flex items-center justify-between gap-3 border-t px-2 pt-4'>
        <span className='mr-auto text-sm text-muted-foreground'>
          {t('common:totalCount', { count: total })}
        </span>
        <span className='text-sm text-muted-foreground'>{t('pageSize')}</span>
        <Select
          value={String(pageSize)}
          onValueChange={(value) => onPageSizeChange(Number(value))}
        >
          <SelectTrigger className='h-8 w-[82px]'>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {[10, 20, 50].map((size) => (
              <SelectItem key={size} value={String(size)}>
                {size}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        <Button
          variant='outline'
          size='icon'
          className='size-8'
          disabled={page <= 1}
          onClick={() => onPageChange(page - 1)}
        >
          <ChevronLeft className='size-4' />
        </Button>
        <span className='text-sm font-medium'>
          {t('pageOf', { page, totalPages })}
        </span>
        <Button
          variant='outline'
          size='icon'
          className='size-8'
          disabled={page >= totalPages}
          onClick={() => onPageChange(page + 1)}
        >
          <ChevronRight className='size-4' />
        </Button>
      </div>

      {/* 确认移除对话框 */}
      <ConfirmDialog
        open={removeConfirmOpen}
        onOpenChange={setRemoveConfirmOpen}
        title={t('removeFromGroup')}
        desc={t('removeConfirm')}
        confirmText={t('common:delete')}
        cancelBtnText={t('common:cancel')}
        destructive
        handleConfirm={() => removeMutation.mutate(targetKeysToRemove)}
        isLoading={removeMutation.isPending}
      />

      {/* 添加设备弹窗 */}
      {group.groupUuid && (
        <AddDevicesDialog
          open={addDialogOpen}
          onOpenChange={setAddDialogOpen}
          groupUuid={group.groupUuid}
          existingDeviceKeys={existingKeys}
        />
      )}
    </div>
  )
}
