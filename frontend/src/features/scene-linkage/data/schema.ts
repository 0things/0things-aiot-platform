import { z } from 'zod'

export const sceneStatusSchema = z.union([
  z.literal('enabled'),
  z.literal('disabled'),
])
export type SceneStatus = z.infer<typeof sceneStatusSchema>

export const sceneStatusStyles = new Map<SceneStatus, string>([
  [
    'enabled',
    'bg-teal-100/30 text-teal-900 dark:text-teal-200 border-teal-200',
  ],
  ['disabled', 'bg-neutral-300/40 border-neutral-300'],
])

export const sceneSchema = z.object({
  id: z.number(),
  name: z.string(),
  description: z.string().optional(),
  status: sceneStatusSchema,
  createdAt: z.string(),
})

export type Scene = z.infer<typeof sceneSchema>
