import { useCallback, useMemo, useState } from 'react'
import {
  addEdge,
  Background,
  Controls,
  Handle,
  MiniMap,
  Position,
  ReactFlow,
  ReactFlowProvider,
  useEdgesState,
  useNodesState,
  useReactFlow,
  type Connection,
  type Edge,
  type Node,
  type NodeProps,
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import {
  Activity,
  ChevronRight,
  Clock3,
  Database,
  Filter,
  GitBranch,
  ListFilter,
  PanelLeftClose,
  PanelLeftOpen,
  Play,
  Save,
  SlidersHorizontal,
  Webhook,
  Waypoints,
  type LucideIcon,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { useGetRuleNodeDefinitions } from '@/api/generated'
import type { AiotBackendApiV1RuleNodeDefinition } from '@/api/generated/model'
import { cn } from '@/lib/utils'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible'
import { Input } from '@/components/ui/input'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'

type RuleNodeDefinition = {
  key: string
  name: string
  description: string
  category: string
  icon: LucideIcon
  inputs: string[]
  outputs: string[]
  configFields: string[]
}

type RuleNodeData = {
  definition: RuleNodeDefinition
}

type RuleFlowNode = Node<RuleNodeData>

const entryDefinition: RuleNodeDefinition = {
  key: 'system.entry',
  name: '消息入口',
  description: '规则链的系统消息入口',
  category: 'System',
  icon: Play,
  inputs: [],
  outputs: ['Success'],
  configFields: [],
}

const categoryIcons: Record<string, LucideIcon> = {
  System: Play,
  Filter,
  Enrichment: Database,
  Transformation: Waypoints,
  Action: Activity,
  External: Webhook,
  Flow: Clock3,
  Analytics: GitBranch,
}

const legacyNodeDefinitions: RuleNodeDefinition[] = [
  {
    key: 'system.entry',
    name: '消息入口',
    description: '规则链的系统消息入口',
    category: 'System',
    icon: Play,
    inputs: [],
    outputs: ['Success'],
    configFields: [],
  },
  {
    key: 'filter.message-type',
    name: '消息类型',
    description: '按规则消息类型分流',
    category: 'Filter',
    icon: ListFilter,
    inputs: ['输入'],
    outputs: ['True', 'False', 'Failure'],
    configFields: ['messageTypes'],
  },
  {
    key: 'filter.property',
    name: '属性条件',
    description: '按属性值判断消息是否通过',
    category: 'Filter',
    icon: Filter,
    inputs: ['输入'],
    outputs: ['True', 'False', 'Failure'],
    configFields: ['property', 'operator', 'value'],
  },
  {
    key: 'filter.value-condition',
    name: '数值条件',
    description: '按数值阈值决定 True 或 False 分支',
    category: 'Filter',
    icon: GitBranch,
    inputs: ['输入'],
    outputs: ['True', 'False', 'Failure'],
    configFields: ['path', 'operator', 'value'],
  },
  {
    key: 'transformation.map',
    name: '字段映射',
    description: '映射消息字段和元数据',
    category: 'Transformation',
    icon: Waypoints,
    inputs: ['输入'],
    outputs: ['Success', 'Failure'],
    configFields: ['mappings'],
  },
  {
    key: 'action.save-attributes',
    name: '保存属性',
    description: '将字段保存为实体属性',
    category: 'Action',
    icon: Database,
    inputs: ['输入'],
    outputs: ['Success', 'Failure'],
    configFields: ['scope', 'attributes'],
  },
  {
    key: 'action.save-telemetry',
    name: '保存遥测',
    description: '将字段保存为遥测数据',
    category: 'Action',
    icon: Activity,
    inputs: ['输入'],
    outputs: ['Success', 'Failure'],
    configFields: ['telemetry'],
  },
  {
    key: 'external.webhook',
    name: '调用 Webhook',
    description: '向 HTTPS 公网地址发送 HTTP 请求',
    category: 'External',
    icon: Webhook,
    inputs: ['输入'],
    outputs: ['Success', 'Failure'],
    configFields: ['url', 'method', 'headers', 'bodyTemplate', 'timeoutMs'],
  },
  {
    key: 'flow.delay',
    name: '延迟',
    description: '延迟后继续执行',
    category: 'Flow',
    icon: Clock3,
    inputs: ['输入'],
    outputs: ['Success', 'Failure'],
    configFields: ['durationSeconds'],
  },
]

function parseLabels(value: string | undefined) {
  try {
    return (JSON.parse(value ?? '[]') as { label?: string }[])
      .map((port) => port.label)
      .filter((label): label is string => !!label)
  } catch {
    return []
  }
}

function toRuleNodeDefinition(
  definition: AiotBackendApiV1RuleNodeDefinition
): RuleNodeDefinition {
  let configFields: string[] = []
  try {
    configFields = Object.keys(
      (
        JSON.parse(definition.configSchema ?? '{}') as {
          properties?: Record<string, unknown>
        }
      ).properties ?? {}
    )
  } catch {
    // Invalid schema should not prevent the rest of the catalog from rendering.
  }

  return {
    key: definition.key ?? '',
    name: definition.name ?? '',
    description: definition.description ?? '',
    category: definition.category ?? 'Action',
    icon: categoryIcons[definition.category ?? ''] ?? Activity,
    inputs: parseLabels(definition.inputPorts),
    outputs: parseLabels(definition.outputPorts),
    configFields,
  }
}

const initialNodes: RuleFlowNode[] = [
  {
    id: 'entry',
    type: 'ruleNode',
    position: { x: 80, y: 220 },
    data: { definition: entryDefinition },
  },
]

function RuleCanvasNode({ data, selected }: NodeProps<RuleFlowNode>) {
  const { definition } = data
  const Icon = definition.icon

  return (
    <div
      className={cn(
        'min-w-36 rounded-md border bg-card p-1.5 shadow-sm transition-shadow',
        selected ? 'border-ring ring-2 ring-ring/30' : 'border-border'
      )}
    >
      {definition.inputs.length > 0 && (
        <Handle type='target' position={Position.Left} id='input' />
      )}
      <div className='flex items-center gap-1'>
        <Icon className='size-3 text-muted-foreground' />
        <span className='text-[11px] font-medium'>{definition.name}</span>
      </div>
      <p className='mt-0.5 text-[10px] text-muted-foreground'>
        {definition.description}
      </p>
      <div className='mt-1.5 flex flex-wrap gap-1'>
        {definition.outputs.map((output) => (
          <div key={output} className='relative'>
            <Badge variant='secondary' className='px-1 py-0 text-[10px]'>
              {output}
            </Badge>
            <Handle
              type='source'
              position={Position.Right}
              id={output}
              className='!bg-muted-foreground'
            />
          </div>
        ))}
      </div>
    </div>
  )
}

const nodeTypes = { ruleNode: RuleCanvasNode }

function RuleChainCanvasInner() {
  const { t } = useTranslation('ruleChain')
  const catalogQuery = useGetRuleNodeDefinitions({
    query: { select: (response) => response.data?.items ?? [] },
  })
  const [nodes, setNodes, onNodesChange] =
    useNodesState<RuleFlowNode>(initialNodes)
  const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([])
  const [selectedNodeId, setSelectedNodeId] = useState<string>()
  const [paletteOpen, setPaletteOpen] = useState(true)
  const [inspectorOpen, setInspectorOpen] = useState(false)
  const [paletteSearch, setPaletteSearch] = useState('')
  const { screenToFlowPosition } = useReactFlow()

  const selectedNode = nodes.find((node) => node.id === selectedNodeId)
  const nodeDefinitions = useMemo(
    () =>
      catalogQuery.data?.length
        ? catalogQuery.data.map(toRuleNodeDefinition)
        : legacyNodeDefinitions,
    [catalogQuery.data]
  )
  const groups = useMemo(
    () =>
      nodeDefinitions
        .filter((definition) => definition.key !== 'system.entry')
        .filter((definition) => {
          const query = paletteSearch.trim().toLocaleLowerCase()
          return (
            !query ||
            definition.name.toLocaleLowerCase().includes(query) ||
            definition.description.toLocaleLowerCase().includes(query) ||
            definition.category.toLocaleLowerCase().includes(query)
          )
        })
        .reduce<Record<string, RuleNodeDefinition[]>>((result, definition) => {
          result[definition.category] ??= []
          result[definition.category].push(definition)
          return result
        }, {}),
    [nodeDefinitions, paletteSearch]
  )

  const onConnect = useCallback(
    (connection: Connection) =>
      setEdges((current) => addEdge(connection, current)),
    [setEdges]
  )

  const addNode = useCallback(
    (definition: RuleNodeDefinition, position = { x: 320, y: 200 }) => {
      const id = crypto.randomUUID()
      setNodes((current) => [
        ...current,
        { id, type: 'ruleNode', position, data: { definition } },
      ])
      setSelectedNodeId(id)
      setInspectorOpen(true)
    },
    [setNodes]
  )

  const onDrop = useCallback(
    (event: React.DragEvent<HTMLDivElement>) => {
      event.preventDefault()
      const key = event.dataTransfer.getData('application/rule-node')
      const definition = nodeDefinitions.find((item) => item.key === key)
      if (!definition) return
      addNode(
        definition,
        screenToFlowPosition({ x: event.clientX, y: event.clientY })
      )
    },
    [addNode, nodeDefinitions, screenToFlowPosition]
  )

  return (
    <div className='flex min-h-0 flex-1 overflow-hidden rounded-xl border bg-background'>
      {paletteOpen ? (
        <aside className='flex w-64 shrink-0 flex-col border-r bg-card'>
          <div className='flex items-start justify-between gap-2 border-b p-4'>
            <div className='flex flex-col gap-1'>
              <h2 className='font-semibold'>{t('palette.title')}</h2>
              <p className='text-xs text-muted-foreground'>
                {t('palette.hint')}
              </p>
            </div>
            <Button
              variant='ghost'
              size='icon'
              onClick={() => setPaletteOpen(false)}
              aria-label={t('actions.collapsePalette')}
            >
              <PanelLeftClose />
            </Button>
          </div>
          <div className='p-3'>
            <Input
              value={paletteSearch}
              onChange={(event) => setPaletteSearch(event.target.value)}
              placeholder={t('palette.searchPlaceholder')}
              aria-label={t('palette.searchPlaceholder')}
            />
          </div>
          <Separator />
          <ScrollArea className='min-h-0 flex-1'>
            <div className='flex flex-col gap-3 p-3'>
              {Object.entries(groups).map(([category, definitions]) => (
                <Collapsible key={category} className='flex flex-col gap-2'>
                  <CollapsibleTrigger asChild>
                    <Button
                      variant='ghost'
                      className='w-full justify-between px-2 text-sm font-medium data-[state=open]:[&>svg]:rotate-90'
                    >
                      <span className='flex items-center gap-2'>
                        <ChevronRight data-icon='inline-start' />
                        {category}
                      </span>
                      <Badge variant='secondary'>{definitions.length}</Badge>
                    </Button>
                  </CollapsibleTrigger>
                  <CollapsibleContent className='flex flex-col gap-1'>
                    {definitions.map((definition) => {
                      const Icon = definition.icon
                      return (
                        <Button
                          key={definition.key}
                          variant='ghost'
                          draggable
                          onDragStart={(event) => {
                            event.dataTransfer.setData(
                              'application/rule-node',
                              definition.key
                            )
                            event.dataTransfer.effectAllowed = 'move'
                          }}
                          onClick={() => addNode(definition)}
                          className='h-auto cursor-grab items-start justify-start gap-2 px-2 py-2 text-left whitespace-normal active:cursor-grabbing'
                        >
                          <Icon className='mt-0.5 size-4 shrink-0 text-muted-foreground' />
                          <span className='flex min-w-0 flex-col gap-0.5'>
                            <span className='text-sm font-medium'>
                              {definition.name}
                            </span>
                            <span className='line-clamp-2 text-xs text-muted-foreground'>
                              {definition.description}
                            </span>
                          </span>
                        </Button>
                      )
                    })}
                  </CollapsibleContent>
                </Collapsible>
              ))}
            </div>
          </ScrollArea>
        </aside>
      ) : (
        <div className='flex w-11 shrink-0 items-start justify-center border-r bg-card pt-3'>
          <Button
            variant='ghost'
            size='icon'
            onClick={() => setPaletteOpen(true)}
            aria-label={t('actions.expandPalette')}
          >
            <PanelLeftOpen />
          </Button>
        </div>
      )}

      <div
        className='relative min-w-0 flex-1'
        onDrop={onDrop}
        onDragOver={(event) => event.preventDefault()}
      >
        <ReactFlow
          nodes={nodes}
          edges={edges}
          nodeTypes={nodeTypes}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          onNodeClick={(_, node) => {
            setSelectedNodeId(node.id)
            setInspectorOpen(true)
          }}
          fitView
          fitViewOptions={{ maxZoom: 1, padding: 0.3 }}
          proOptions={{ hideAttribution: true }}
        >
          <Background gap={18} />
          <Controls showInteractive={false} />
          <MiniMap pannable zoomable />
        </ReactFlow>
      </div>

      <Sheet open={inspectorOpen} onOpenChange={setInspectorOpen}>
        <SheetContent className='flex w-full flex-col gap-0 sm:max-w-sm'>
          <SheetHeader className='border-b px-6 py-5 text-left'>
            <SheetTitle>{t('inspector.title')}</SheetTitle>
            <SheetDescription>
              {selectedNode
                ? selectedNode.data.definition.description
                : t('inspector.selectHint')}
            </SheetDescription>
          </SheetHeader>
          {selectedNode ? (
            <div className='flex flex-col gap-5 p-6'>
              <div className='flex items-center gap-2'>
                <SlidersHorizontal className='size-4 text-muted-foreground' />
                <p className='text-sm font-medium'>
                  {selectedNode.data.definition.name}
                </p>
              </div>
              <Separator />
              <div className='flex flex-col gap-2'>
                <p className='text-xs font-medium text-muted-foreground'>
                  {t('inspector.configFields')}
                </p>
                {selectedNode.data.definition.configFields.length > 0 ? (
                  selectedNode.data.definition.configFields.map((field) => (
                    <Badge key={field} variant='outline'>
                      {field}
                    </Badge>
                  ))
                ) : (
                  <p className='text-sm text-muted-foreground'>
                    {t('inspector.noConfig')}
                  </p>
                )}
              </div>
            </div>
          ) : null}
        </SheetContent>
      </Sheet>
    </div>
  )
}

export function RuleChainCanvasPage() {
  const { t } = useTranslation('ruleChain')

  return (
    <>
      <Header fixed>
        <Search />
        <div className='ms-auto flex items-center gap-4'>
          <ThemeSwitch />
          <ConfigDrawer />
          <ProfileDropdown />
        </div>
      </Header>
      <Main fixed fluid className='flex min-h-0 flex-1 flex-col gap-3 py-3'>
        <div className='flex flex-wrap items-center justify-between gap-3 border-b pb-3'>
          <div className='flex flex-col gap-1'>
            <div className='flex items-center gap-2'>
              <h1 className='text-xl font-semibold tracking-tight'>
                {t('title')}
              </h1>
              <Badge variant='secondary'>{t('draft')}</Badge>
            </div>
            <p className='text-sm text-muted-foreground'>{t('description')}</p>
          </div>
          <div className='flex items-center gap-2'>
            <Button variant='outline' size='sm' disabled>
              <Play data-icon='inline-start' />
              {t('testComingSoon')}
            </Button>
            <Button size='sm' disabled>
              <Save data-icon='inline-start' />
              {t('saveComingSoon')}
            </Button>
          </div>
        </div>
        <ReactFlowProvider>
          <RuleChainCanvasInner />
        </ReactFlowProvider>
      </Main>
    </>
  )
}
