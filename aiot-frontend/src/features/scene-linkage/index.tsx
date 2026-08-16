import { useState } from 'react'
import { Plus, Save, Trash2, Workflow } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import {
  actionPlugins,
  createSceneAction,
  type SceneAction,
  type SceneActionType,
} from './action-plugins'

export function SceneLinkageDetailPage() {
  const { t } = useTranslation('sceneLinkage')
  const [enabled, setEnabled] = useState(true)
  const [actions, setActions] = useState<SceneAction[]>([
    createSceneAction('set_property'),
  ])

  const updateAction = (id: string, config: Record<string, string>) => {
    setActions((current) =>
      current.map((action) =>
        action.id === id ? { ...action, config } : action
      )
    )
  }

  const addAction = (type: SceneActionType) =>
    setActions((current) => [...current, createSceneAction(type)])

  return (
    <>
      <Header fixed>
        <Search />
        <div className='ms-auto flex items-center space-x-4'>
          <ThemeSwitch />
          <ConfigDrawer />
          <ProfileDropdown />
        </div>
      </Header>

      <Main fixed className='flex flex-1 flex-col gap-6'>
        <div className='flex flex-wrap items-start justify-between gap-4'>
          <div>
            <div className='mb-2 flex items-center gap-2 text-primary'>
              <Workflow className='size-5' />
              <span className='text-sm font-medium'>{t('eyebrow')}</span>
            </div>
            <h1 className='text-3xl font-bold tracking-tight'>{t('title')}</h1>
            <p className='mt-1 text-muted-foreground'>{t('description')}</p>
          </div>
          <Button onClick={() => toast.success(t('savedLocally'))}>
            <Save className='mr-2 size-4' />
            {t('saveDraft')}
          </Button>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>{t('basic.title')}</CardTitle>
            <CardDescription>{t('basic.description')}</CardDescription>
          </CardHeader>
          <CardContent className='grid gap-4 md:grid-cols-[minmax(0,1fr)_minmax(0,2fr)_auto] md:items-end'>
            <div className='space-y-1.5'>
              <Label>{t('basic.name')}</Label>
              <Input defaultValue={t('basic.defaultName')} />
            </div>
            <div className='space-y-1.5'>
              <Label>{t('basic.note')}</Label>
              <Input placeholder={t('basic.notePlaceholder')} />
            </div>
            <div className='flex items-center gap-3 pb-2'>
              <Switch checked={enabled} onCheckedChange={setEnabled} />
              <Label>
                {enabled ? t('basic.enabled') : t('basic.disabled')}
              </Label>
            </div>
          </CardContent>
        </Card>

        <section className='space-y-3'>
          <div className='flex items-center gap-3'>
            <div className='flex size-7 items-center justify-center rounded-full bg-primary text-sm font-semibold text-primary-foreground'>
              1
            </div>
            <div>
              <h2 className='font-semibold'>{t('trigger.title')}</h2>
              <p className='text-sm text-muted-foreground'>
                {t('trigger.description')}
              </p>
            </div>
          </div>
          <Card className='gap-0 py-0'>
            <CardContent className='grid gap-4 py-5 md:grid-cols-3'>
              <div className='space-y-1.5'>
                <Label>{t('trigger.eventType')}</Label>
                <Select defaultValue='property-report'>
                  <SelectTrigger className='w-full'>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value='property-report'>
                      {t('trigger.propertyReport')}
                    </SelectItem>
                    <SelectItem value='device-status'>
                      {t('trigger.deviceStatus')}
                    </SelectItem>
                    <SelectItem value='device-event'>
                      {t('trigger.deviceEvent')}
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className='space-y-1.5'>
                <Label>{t('trigger.property')}</Label>
                <Input placeholder='temperature' />
              </div>
              <div className='space-y-1.5'>
                <Label>{t('trigger.condition')}</Label>
                <Select defaultValue='gt'>
                  <SelectTrigger className='w-full'>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value='gt'>&gt; 30</SelectItem>
                    <SelectItem value='eq'>= true</SelectItem>
                    <SelectItem value='changed'>
                      {t('trigger.changed')}
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </CardContent>
          </Card>
        </section>

        <section className='space-y-3'>
          <div className='flex items-center gap-3'>
            <div className='flex size-7 items-center justify-center rounded-full bg-primary text-sm font-semibold text-primary-foreground'>
              2
            </div>
            <div>
              <h2 className='font-semibold'>{t('actions.title')}</h2>
              <p className='text-sm text-muted-foreground'>
                {t('actions.description')}
              </p>
            </div>
          </div>

          <div className='space-y-3'>
            {actions.map((action, index) => {
              const plugin = actionPlugins.find(
                (item) => item.type === action.type
              )
              if (!plugin) return null
              const Icon = plugin.icon
              const Config = plugin.Config
              return (
                <Card key={action.id} className='gap-0 py-0'>
                  <CardContent className='grid gap-4 py-5 lg:grid-cols-[220px_minmax(0,1fr)_auto] lg:items-start'>
                    <div className='space-y-1.5'>
                      <Label>{t('actions.action', { index: index + 1 })}</Label>
                      <Select
                        value={action.type}
                        onValueChange={(type: SceneActionType) =>
                          setActions((current) =>
                            current.map((item) =>
                              item.id === action.id
                                ? { ...createSceneAction(type), id: item.id }
                                : item
                            )
                          )
                        }
                      >
                        <SelectTrigger className='w-full'>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          {actionPlugins.map((item) => (
                            <SelectItem key={item.type} value={item.type}>
                              {t(item.labelKey)}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <div className='rounded-lg bg-muted/45 p-4'>
                      <div className='mb-3 flex items-center gap-2'>
                        <Icon className='size-4 text-primary' />
                        <span className='text-sm font-medium'>
                          {t(plugin.descriptionKey)}
                        </span>
                      </div>
                      <Config
                        value={action.config}
                        onChange={(config) => updateAction(action.id, config)}
                        t={t}
                      />
                    </div>
                    <Button
                      variant='ghost'
                      size='icon'
                      className='text-muted-foreground hover:text-destructive'
                      disabled={actions.length === 1}
                      onClick={() =>
                        setActions((current) =>
                          current.filter((item) => item.id !== action.id)
                        )
                      }
                      aria-label={t('actions.remove')}
                    >
                      <Trash2 className='size-4' />
                    </Button>
                  </CardContent>
                </Card>
              )
            })}
          </div>

          <Select onValueChange={(type: SceneActionType) => addAction(type)}>
            <SelectTrigger className='w-fit'>
              <Plus className='size-4' />
              <SelectValue placeholder={t('actions.add')} />
            </SelectTrigger>
            <SelectContent>
              {actionPlugins.map((plugin) => (
                <SelectItem key={plugin.type} value={plugin.type}>
                  {t(plugin.labelKey)}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </section>

        <div className='rounded-lg border border-dashed bg-muted/20 px-4 py-3 text-sm text-muted-foreground'>
          {t('futureNote')}
        </div>
      </Main>
    </>
  )
}
