/* eslint-disable react-refresh/only-export-components -- action plugin registry exports configuration components and metadata together. */
import type { ComponentType } from 'react'
import { Bell, Braces, Gauge, Send, Webhook } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'

export type SceneActionType =
  | 'set_property'
  | 'invoke_service'
  | 'send_notification'
  | 'webhook'
  | 'trigger_scene'

export type SceneAction = {
  id: string
  type: SceneActionType
  config: Record<string, string>
}

type ActionPlugin = {
  type: SceneActionType
  labelKey: string
  descriptionKey: string
  icon: ComponentType<{ className?: string }>
  defaultConfig: Record<string, string>
  Config: ({ value, onChange }: ActionConfigProps) => React.ReactNode
}

type ActionConfigProps = {
  value: Record<string, string>
  onChange: (config: Record<string, string>) => void
  t: (key: string) => string
}

function Field({
  label,
  value,
  placeholder,
  onChange,
}: {
  label: string
  value: string
  placeholder: string
  onChange: (value: string) => void
}) {
  return (
    <div className='space-y-1.5'>
      <Label className='text-xs text-muted-foreground'>{label}</Label>
      <Input
        value={value}
        placeholder={placeholder}
        onChange={(event) => onChange(event.target.value)}
      />
    </div>
  )
}

function PropertyConfig({ value, onChange, t }: ActionConfigProps) {
  const update = (key: string, next: string) =>
    onChange({ ...value, [key]: next })
  return (
    <div className='grid gap-3 sm:grid-cols-2'>
      <Field
        label={t('fields.deviceKey')}
        value={value.deviceKey ?? ''}
        placeholder={t('fields.deviceKeyPlaceholder')}
        onChange={(next) => update('deviceKey', next)}
      />
      <Field
        label={t('fields.propertyKey')}
        value={value.property ?? ''}
        placeholder={t('fields.propertyKeyPlaceholder')}
        onChange={(next) => update('property', next)}
      />
      <div className='sm:col-span-2'>
        <Field
          label={t('fields.value')}
          value={value.value ?? ''}
          placeholder={t('fields.valuePlaceholder')}
          onChange={(next) => update('value', next)}
        />
      </div>
    </div>
  )
}

function ServiceConfig({ value, onChange, t }: ActionConfigProps) {
  const update = (key: string, next: string) =>
    onChange({ ...value, [key]: next })
  return (
    <div className='grid gap-3 sm:grid-cols-2'>
      <Field
        label={t('fields.deviceKey')}
        value={value.deviceKey ?? ''}
        placeholder={t('fields.deviceKeyPlaceholder')}
        onChange={(next) => update('deviceKey', next)}
      />
      <Field
        label={t('fields.serviceKey')}
        value={value.service ?? ''}
        placeholder={t('fields.serviceKeyPlaceholder')}
        onChange={(next) => update('service', next)}
      />
    </div>
  )
}

function NotificationConfig({ value, onChange, t }: ActionConfigProps) {
  return (
    <div className='space-y-1.5'>
      <Label className='text-xs text-muted-foreground'>
        {t('fields.message')}
      </Label>
      <Textarea
        value={value.message ?? ''}
        placeholder={t('fields.messagePlaceholder')}
        onChange={(event) =>
          onChange({ ...value, message: event.target.value })
        }
      />
    </div>
  )
}

function WebhookConfig({ value, onChange, t }: ActionConfigProps) {
  const update = (key: string, next: string) =>
    onChange({ ...value, [key]: next })
  return (
    <div className='grid gap-3 sm:grid-cols-2'>
      <Field
        label={t('fields.url')}
        value={value.url ?? ''}
        placeholder={t('fields.urlPlaceholder')}
        onChange={(next) => update('url', next)}
      />
      <Field
        label={t('fields.method')}
        value={value.method ?? ''}
        placeholder={t('fields.methodPlaceholder')}
        onChange={(next) => update('method', next)}
      />
    </div>
  )
}

function SceneConfig({ value, onChange, t }: ActionConfigProps) {
  return (
    <Field
      label={t('fields.targetScene')}
      value={value.sceneId ?? ''}
      placeholder={t('fields.targetScenePlaceholder')}
      onChange={(next) => onChange({ ...value, sceneId: next })}
    />
  )
}

export const actionPlugins: ActionPlugin[] = [
  {
    type: 'set_property',
    labelKey: 'actions.setProperty',
    descriptionKey: 'actions.setPropertyDescription',
    icon: Gauge,
    defaultConfig: { deviceKey: '', property: '', value: '' },
    Config: PropertyConfig,
  },
  {
    type: 'invoke_service',
    labelKey: 'actions.invokeService',
    descriptionKey: 'actions.invokeServiceDescription',
    icon: Send,
    defaultConfig: { deviceKey: '', service: '' },
    Config: ServiceConfig,
  },
  {
    type: 'send_notification',
    labelKey: 'actions.notification',
    descriptionKey: 'actions.notificationDescription',
    icon: Bell,
    defaultConfig: { message: '' },
    Config: NotificationConfig,
  },
  {
    type: 'webhook',
    labelKey: 'actions.webhook',
    descriptionKey: 'actions.webhookDescription',
    icon: Webhook,
    defaultConfig: { url: '', method: 'POST' },
    Config: WebhookConfig,
  },
  {
    type: 'trigger_scene',
    labelKey: 'actions.triggerScene',
    descriptionKey: 'actions.triggerSceneDescription',
    icon: Braces,
    defaultConfig: { sceneId: '' },
    Config: SceneConfig,
  },
]

export function createSceneAction(type: SceneActionType): SceneAction {
  const plugin = actionPlugins.find((item) => item.type === type)
  if (!plugin) throw new Error(`Unknown scene action plugin: ${type}`)

  return {
    id: crypto.randomUUID(),
    type,
    config: { ...plugin.defaultConfig },
  }
}
