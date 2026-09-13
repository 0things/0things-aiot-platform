import type { Node } from '@xyflow/react'
import type { LucideIcon } from 'lucide-react'

export type ConfigField = {
  key: string
  type: string
  title: string
}

export type RuleNodeDefinition = {
  key: string
  name: string
  description: string
  category: string
  icon: LucideIcon
  isSystem: boolean
  inputs: string[]
  outputs: string[]
  configFields: ConfigField[]
  defaultConfig: Record<string, unknown>
}

export type RuleNodeData = {
  definition: RuleNodeDefinition
  config: Record<string, unknown>
}

export type RuleFlowNode = Node<RuleNodeData>
