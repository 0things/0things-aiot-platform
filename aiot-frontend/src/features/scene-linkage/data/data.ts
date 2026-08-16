import { type SceneStatus } from './schema'

export const statuses: { value: SceneStatus; label: string }[] = [
  { value: 'enabled', label: 'Enabled' },
  { value: 'disabled', label: 'Disabled' },
]
