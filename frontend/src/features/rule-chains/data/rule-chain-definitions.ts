import type { Edge } from '@xyflow/react'
import {
  Activity,
  CircleHelp,
  Clock3,
  Database,
  Filter,
  GitBranch,
  ListFilter,
  Play,
  Save,
  Send,
  Waypoints,
  Webhook,
} from 'lucide-react'
import type { AiotBackendApiV1RuleNodeDefinition } from '@/api/generated/model'
import type {
  ConfigField,
  RuleFlowNode,
  RuleNodeDefinition,
} from '../types/rule-chain-types'

const iconMap = {
  Activity,
  Clock3,
  Database,
  Filter,
  GitBranch,
  ListFilter,
  Play,
  Save,
  Send,
  Waypoints,
  Webhook,
}

function ports(value?: string) {
  try {
    return (JSON.parse(value ?? '[]') as { label?: string }[])
      .map((item) => item.label)
      .filter((item): item is string => !!item)
  } catch {
    return []
  }
}

function schemaFields(
  definition: AiotBackendApiV1RuleNodeDefinition
): ConfigField[] {
  try {
    const schema = JSON.parse(definition.configSchema ?? '{}') as {
      properties?: Record<string, { type?: string; title?: string }>
    }
    return Object.entries(schema.properties ?? {}).map(([key, value]) => ({
      key,
      type: value.type ?? 'string',
      title: value.title ?? key,
    }))
  } catch {
    return []
  }
}

export function toDefinition(
  definition: AiotBackendApiV1RuleNodeDefinition
): RuleNodeDefinition {
  return {
    key: definition.key ?? '',
    name: definition.name ?? '',
    description: definition.description ?? '',
    category: definition.category ?? 'Action',
    icon: iconMap[definition.icon as keyof typeof iconMap] ?? CircleHelp,
    isSystem: definition.isSystem ?? false,
    inputs: ports(definition.inputPorts),
    outputs: ports(definition.outputPorts),
    configFields: schemaFields(definition),
    defaultConfig: parseDefaultConfig(definition.defaultConfig),
  }
}

function parseDefaultConfig(value?: string): Record<string, unknown> {
  try {
    const parsed = JSON.parse(value ?? '{}')
    return parsed && typeof parsed === 'object' && !Array.isArray(parsed)
      ? (parsed as Record<string, unknown>)
      : {}
  } catch {
    return {}
  }
}

export function defaultConfig(definition: RuleNodeDefinition) {
  return { ...definition.defaultConfig }
}

export function entryNode(definition: RuleNodeDefinition): RuleFlowNode {
  return {
    id: 'entry',
    type: 'ruleNode',
    position: { x: 80, y: 220 },
    data: { definition, config: defaultConfig(definition) },
  }
}

export function serializeGraph(nodes: RuleFlowNode[], edges: Edge[]) {
  return {
    nodes: nodes.map((node) => ({
      id: node.id,
      type: node.type,
      position: node.position,
      data: {
        definitionKey: node.data.definition.key,
        config: node.data.config,
      },
    })),
    edges: edges.map((edge) => ({
      id: edge.id,
      source: edge.source,
      target: edge.target,
      sourceHandle: edge.sourceHandle,
      targetHandle: edge.targetHandle,
      label: edge.label,
    })),
  }
}
