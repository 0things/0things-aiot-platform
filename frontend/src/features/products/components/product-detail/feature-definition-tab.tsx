import { useState, useEffect } from 'react'
import { z } from 'zod'
import { AxiosError } from 'axios'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  AlertCircle,
  CheckCircle,
  Code,
  FileJson,
  Zap,
  Shield,
  RefreshCw,
  Plus,
  Edit2,
  Trash2,
  List,
  CheckCheck,
  XCircle,
  Loader2,
  HelpCircle,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import {
  deleteProductsProductKeyTsl,
  getGetProductsProductKeyTslQueryKey,
  getProductsProductKeyTsl,
  postProductsProductKeyTsl,
} from '@/api/generated'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { Textarea } from '@/components/ui/textarea'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'

// TSL 类型定义
interface DataTypeSpecs {
  min?: string | number
  max?: string | number
  unit?: string
  step?: string | number
  [key: string]: string | number | undefined
}

interface DataType {
  type:
    'int' | 'float' | 'double' | 'bool' | 'string' | 'enum' | 'struct' | 'array'
  specs?: DataTypeSpecs
}

interface Property {
  identifier: string
  name: string
  accessMode: 'r' | 'rw'
  required: boolean
  desc?: string
  description?: string
  dataType: DataType
}

interface ServiceParam {
  identifier: string
  name: string
  dataType: DataType
}

interface Service {
  identifier: string
  name: string
  required: boolean
  callType: 'async' | 'sync'
  inputData: ServiceParam[]
  outputData: ServiceParam[]
}

interface EventParam {
  identifier: string
  name: string
  dataType: DataType
}

interface Event {
  identifier: string
  name: string
  type: 'info' | 'alert' | 'error'
  required: boolean
  outputData: EventParam[]
}

interface TSLModel {
  schema: string
  version: string
  profile: {
    productKey: string
  }
  properties: Property[]
  events: Event[]
  services: Service[]
}

// TSL模板 - 名称使用翻译键，实际显示名称在渲染时翻译
const TSL_TEMPLATES = {
  empty: {
    nameKey: 'productDetail.featureDefinition.templates.empty', // 翻译键
    value: {
      schema: 'schema.json',
      version: '1.0.0',
      profile: {
        productKey: 'PRODUCT_KEY',
      },
      properties: [],
      events: [],
      services: [],
    },
  },
  sensor: {
    nameKey: 'productDetail.featureDefinition.templates.sensor',
    value: {
      schema: 'schema.json',
      version: '1.0.0',
      profile: {
        productKey: 'PRODUCT_KEY',
      },
      properties: [
        {
          identifier: 'temperature',
          name: 'Temperature',
          accessMode: 'r',
          required: false,
          dataType: {
            type: 'double',
            specs: {
              min: '-50',
              max: '150',
              unit: '°C',
              step: '0.1',
            },
          },
        },
        {
          identifier: 'humidity',
          name: 'Humidity',
          accessMode: 'r',
          required: false,
          dataType: {
            type: 'double',
            specs: {
              min: '0',
              max: '100',
              unit: '%',
              step: '0.1',
            },
          },
        },
      ],
      events: [
        {
          identifier: 'temp_humi_report',
          name: 'Temperature and Humidity Report',
          type: 'info',
          required: false,
          outputData: [
            {
              identifier: 'temperature',
              name: 'Temperature',
              dataType: {
                type: 'double',
                specs: {
                  min: '-50',
                  max: '150',
                  unit: '°C',
                  step: '0.1',
                },
              },
            },
            {
              identifier: 'humidity',
              name: 'Humidity',
              dataType: {
                type: 'double',
                specs: {
                  min: '0',
                  max: '100',
                  unit: '%',
                  step: '0.1',
                },
              },
            },
          ],
        },
      ],
      services: [],
    },
  },
  switch: {
    nameKey: 'productDetail.featureDefinition.templates.switch',
    value: {
      schema: 'schema.json',
      version: '1.0.0',
      profile: {
        productKey: 'PRODUCT_KEY',
      },
      properties: [
        {
          identifier: 'powerstate',
          name: 'Power State',
          accessMode: 'rw',
          required: false,
          dataType: {
            type: 'bool',
            specs: {
              '0': 'Off',
              '1': 'On',
            },
          },
        },
      ],
      events: [],
      services: [
        {
          identifier: 'set',
          name: 'Set',
          required: false,
          callType: 'async',
          inputData: [
            {
              identifier: 'powerstate',
              name: 'Power State',
              dataType: {
                type: 'bool',
                specs: {
                  '0': 'Off',
                  '1': 'On',
                },
              },
            },
          ],
          outputData: [],
        },
      ],
    },
  },
  advanced: {
    nameKey: 'productDetail.featureDefinition.templates.advanced',
    value: {
      schema: 'schema.json',
      version: '1.0.0',
      profile: {
        productKey: 'PRODUCT_KEY',
      },
      properties: [
        {
          identifier: 'location',
          name: 'Location',
          accessMode: 'r',
          required: false,
          dataType: {
            type: 'struct',
            specs: {
              lng: { type: 'double' },
              lat: { type: 'double' },
              speed: { type: 'float' },
            },
          },
        },
        {
          identifier: 'status',
          name: 'Status',
          accessMode: 'r',
          required: false,
          dataType: {
            type: 'enum',
            specs: {
              '0': 'Standby',
              '1': 'Working',
              '2': 'Fault',
            },
          },
        },
        {
          identifier: 'tags',
          name: 'Tags',
          accessMode: 'rw',
          required: false,
          dataType: {
            type: 'array',
            specs: {
              type: 'string',
              size: 10,
            },
          },
        },
      ],
      events: [
        {
          identifier: 'fault',
          name: 'Fault Alert',
          type: 'error',
          required: false,
          outputData: [
            {
              identifier: 'code',
              name: 'Fault Code',
              dataType: {
                type: 'int',
              },
            },
            {
              identifier: 'message',
              name: 'Fault Message',
              dataType: {
                type: 'string',
              },
            },
          ],
        },
      ],
      services: [
        {
          identifier: 'reboot',
          name: 'Reboot',
          required: false,
          callType: 'async',
          inputData: [],
          outputData: [
            {
              identifier: 'result',
              name: 'Result',
              dataType: {
                type: 'bool',
              },
            },
          ],
        },
      ],
    },
  },
}

// Zod Schema for TSL validation - 使用 lazy 支持递归
const dataTypeSchema: z.ZodType<DataType> = z.lazy(() =>
  z.object({
    type: z.enum([
      'int',
      'float',
      'double',
      'bool',
      'string',
      'enum',
      'struct',
      'array',
    ]),
    specs: z.any().optional(), // specs 结构在不同类型下差异很大，使用 any 更灵活
  })
)

// 基础验证 schema（不包含自定义错误信息）
const propertySchema = z.object({
  identifier: z.string().min(1),
  name: z.string().min(1),
  accessMode: z.enum(['r', 'rw']),
  required: z.boolean(),
  dataType: dataTypeSchema,
})

const serviceParamSchema = z.object({
  identifier: z.string().min(1),
  name: z.string().min(1),
  dataType: dataTypeSchema,
})

const serviceSchema = z.object({
  identifier: z.string().min(1),
  name: z.string().min(1),
  required: z.boolean(),
  callType: z.enum(['async', 'sync']),
  inputData: z.array(serviceParamSchema),
  outputData: z.array(serviceParamSchema),
})

const eventParamSchema = z.object({
  identifier: z.string().min(1),
  name: z.string().min(1),
  dataType: dataTypeSchema,
})

const eventSchema = z.object({
  identifier: z.string().min(1),
  name: z.string().min(1),
  type: z.enum(['info', 'alert', 'error']),
  required: z.boolean(),
  outputData: z.array(eventParamSchema),
})

const tslSchema = z.object({
  schema: z.string().optional().default('schema.json'),
  profile: z
    .object({
      productKey: z.string().min(1),
    })
    .optional(),
  properties: z.array(propertySchema).optional().default([]),
  events: z.array(eventSchema).optional().default([]),
  services: z.array(serviceSchema).optional().default([]),
})

function normalizeTSLModel(
  model: Partial<TSLModel> | null | undefined
): TSLModel {
  if (!model) {
    return TSL_TEMPLATES.empty.value as TSLModel
  }
  return {
    schema: model.schema || 'schema.json',
    version: model.version || '1.0.0',
    profile: {
      productKey: model.profile?.productKey || '',
    },
    properties: Array.isArray(model.properties) ? model.properties : [],
    events: Array.isArray(model.events) ? model.events : [],
    services: Array.isArray(model.services) ? model.services : [],
  }
}

interface FeatureDefinitionTabProps {
  productKey: string
}

export function FeatureDefinitionTab({
  productKey,
}: FeatureDefinitionTabProps) {
  const { t } = useTranslation(['deviceManagement', 'common'])
  const queryClient = useQueryClient()

  // Query for fetching TSL data
  const {
    data: tslData,
    isLoading: isLoadingTSL,
    error: tslError,
  } = useQuery({
    queryKey: getGetProductsProductKeyTslQueryKey(productKey),
    queryFn: async () => {
      try {
        const response = await getProductsProductKeyTsl(productKey)
        // Parse the TSL string from API into TSLModel
        if (response.data?.productTsl?.tsl) {
          const parsed = JSON.parse(response.data.productTsl.tsl)
          return normalizeTSLModel(parsed)
        }
        return null
      } catch (error: unknown) {
        // If 404, it means no TSL exists, which is fine
        if (error instanceof AxiosError && error.response?.status === 404) {
          return null
        }
        throw error
      }
    },
    enabled: !!productKey,
  })

  // Mutation for saving/updating TSL
  const saveTSLMutation = useMutation({
    mutationFn: async (tsl: TSLModel) => {
      const tslJsonString = JSON.stringify(tsl)
      return postProductsProductKeyTsl(productKey, { tsl: tslJsonString })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['product-tsl', productKey] })
      toast.success(t('productDetail.featureDefinition.status.saveSuccess'))
    },
    onError: () => {
      toast.error(t('productDetail.featureDefinition.status.saveFailed'))
    },
  })

  // Mutation for deleting TSL
  const deleteTSLMutation = useMutation({
    mutationFn: async () => {
      return deleteProductsProductKeyTsl(productKey)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['product-tsl', productKey] })
      toast.success(t('productDetail.featureDefinition.status.deleteSuccess'))
    },
    onError: () => {
      toast.error(t('productDetail.featureDefinition.status.deleteFailed'))
    },
  })

  const [tslModel, setTslModel] = useState<TSLModel>(
    normalizeTSLModel(TSL_TEMPLATES.empty.value as TSLModel)
  )
  const [tslText, setTslText] = useState<string>(
    JSON.stringify(TSL_TEMPLATES.empty.value, null, 2)
  )
  const [error, setError] = useState<string>('')
  const [success, setSuccess] = useState<string>('')
  const [currentTemplate, setCurrentTemplate] = useState<string>('empty')
  const [editMode, setEditMode] = useState<'visual' | 'json'>('visual')
  const [jsonValid, setJsonValid] = useState<boolean | null>(true) // null = unknown, true = valid, false = invalid
  const [validationError, setValidationError] = useState<string>('')

  // 属性编辑对话框状态
  const [propertyDialogOpen, setPropertyDialogOpen] = useState(false)
  const [editingProperty, setEditingProperty] = useState<Property | null>(null)
  const [editingPropertyIndex, setEditingPropertyIndex] = useState<number>(-1)

  // 服务编辑对话框状态
  const [serviceDialogOpen, setServiceDialogOpen] = useState(false)
  const [editingService, setEditingService] = useState<Service | null>(null)
  const [editingServiceIndex, setEditingServiceIndex] = useState<number>(-1)

  // 事件编辑对话框状态
  const [eventDialogOpen, setEventDialogOpen] = useState(false)
  const [editingEvent, setEditingEvent] = useState<Event | null>(null)
  const [editingEventIndex, setEditingEventIndex] = useState<number>(-1)

  // 删除确认对话框状态
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)

  // Sync fetched TSL data to editor state.
  // The synchronous setState is deferred to a microtask so the effect body itself
  // does not synchronously setState (which would trigger cascading renders).
  useEffect(() => {
    let cancelled = false
    const id = setTimeout(() => {
      if (cancelled) return
      if (tslData) {
        const normalized = normalizeTSLModel(tslData)
        setTslModel(normalized)
        setTslText(JSON.stringify(normalized, null, 2))
        setError('')
        setSuccess('')
        setCurrentTemplate('empty') // Reset template selection
      } else if (!isLoadingTSL && !tslError) {
        // No TSL data exists, use empty template
        const emptyModel = normalizeTSLModel(
          TSL_TEMPLATES.empty.value as TSLModel
        )
        setTslModel(emptyModel)
        setTslText(JSON.stringify(emptyModel, null, 2))
        setCurrentTemplate('empty')
      }
    }, 0)
    return () => {
      cancelled = true
      clearTimeout(id)
    }
  }, [tslData, isLoadingTSL, tslError])

  // 从 JSON 更新模型
  const syncFromJson = (text: string) => {
    try {
      const parsed = JSON.parse(text)
      const normalized = normalizeTSLModel(parsed)
      setTslModel(normalized)
      setError('')
      return true
    } catch (e) {
      const message = e instanceof Error ? e.message : String(e)
      setError(
        `${t('productDetail.featureDefinition.errors.jsonParse')}: ${message}`
      )
      return false
    }
  }

  // 从模型更新 JSON
  const syncToJson = (model: TSLModel) => {
    const normalized = normalizeTSLModel(model)
    const formatted = JSON.stringify(normalized, null, 2)
    setTslText(formatted)
  }

  // 编辑模式切换时的同步逻辑见下方 Tabs 的 onValueChange 处理（用户事件触发，不在 effect 中同步 setState）

  // 验证TSL JSON格式
  const validateTSL = (text: string): boolean => {
    try {
      const parsed = JSON.parse(text)

      // 使用 zod 进行验证
      tslSchema.parse(parsed)

      setJsonValid(true)
      setValidationError('')
      return true
    } catch (e) {
      if (e instanceof z.ZodError) {
        const firstError = e.issues[0]
        // 将 Zod 默认错误信息转换为国际化
        let errorMessage = ''
        if (firstError.code === 'too_small' && firstError.minimum === 1) {
          const field = firstError.path[firstError.path.length - 1]
          if (field === 'identifier') {
            errorMessage = t(
              'productDetail.featureDefinition.errors.identifierRequired'
            )
          } else if (field === 'name') {
            errorMessage = t(
              'productDetail.featureDefinition.errors.nameRequired'
            )
          } else if (field === 'productKey') {
            errorMessage = t(
              'productDetail.featureDefinition.errors.productKeyRequired'
            )
          } else {
            errorMessage = `${firstError.path.join('.')}: ${firstError.message}`
          }
        } else {
          errorMessage = `${firstError.path.join('.')}: ${firstError.message}`
        }
        setValidationError(errorMessage)
        setError(
          `${t('productDetail.featureDefinition.errors.tslValidation')}: ${errorMessage}`
        )
      } else if (e instanceof SyntaxError) {
        setValidationError(
          t('productDetail.featureDefinition.errors.jsonSyntax')
        )
        setError(
          `${t('productDetail.featureDefinition.errors.jsonSyntax')}: ${e.message}`
        )
      } else {
        const message = e instanceof Error ? e.message : String(e)
        setValidationError(message)
        setError(message)
      }
      setJsonValid(false)
      return false
    }
  }

  // 实时验证 JSON（不显示错误，只更新状态）
  const validateJsonSyntax = (text: string): void => {
    if (!text.trim()) {
      setJsonValid(null)
      setValidationError('')
      return
    }

    try {
      const parsed = JSON.parse(text)

      // 使用 zod 进行验证
      tslSchema.parse(parsed)

      setJsonValid(true)
      setValidationError('')
    } catch (e) {
      setJsonValid(false)
      if (e instanceof z.ZodError) {
        const firstError = e.issues[0]
        setValidationError(
          `${firstError.path.join('.')}: ${firstError.message}`
        )
      } else if (e instanceof SyntaxError) {
        setValidationError(
          t('productDetail.featureDefinition.errors.jsonSyntax')
        )
      } else {
        setValidationError(e instanceof Error ? e.message : String(e))
      }
    }
  }

  // 加载模板
  const loadTemplate = (templateKey: string) => {
    const template = TSL_TEMPLATES[templateKey as keyof typeof TSL_TEMPLATES]
    if (template) {
      const newModel = normalizeTSLModel(template.value as TSLModel)
      setTslModel(newModel)
      syncToJson(newModel)
      setCurrentTemplate(templateKey)
      setError('')
      setSuccess(
        t('productDetail.featureDefinition.status.templateLoaded', {
          name: t(template.nameKey),
        })
      )
      setTimeout(() => setSuccess(''), 3000)
    }
  }

  // 格式化JSON
  const formatTSL = () => {
    try {
      const parsed = JSON.parse(tslText)

      // 使用 zod 验证
      tslSchema.parse(parsed)

      const normalized = normalizeTSLModel(parsed)
      const formatted = JSON.stringify(normalized, null, 2)
      setTslText(formatted)
      setTslModel(normalized)
      setJsonValid(true)
      setValidationError('')
      setError('')
      setSuccess(t('productDetail.featureDefinition.status.jsonFormatted'))
      setTimeout(() => setSuccess(''), 3000)
    } catch (e) {
      if (e instanceof z.ZodError) {
        const firstError = e.issues[0]
        // 将 Zod 默认错误信息转换为国际化
        let errorMessage = ''
        if (firstError.code === 'too_small' && firstError.minimum === 1) {
          const field = firstError.path[firstError.path.length - 1]
          if (field === 'identifier') {
            errorMessage = t(
              'productDetail.featureDefinition.errors.identifierRequired'
            )
          } else if (field === 'name') {
            errorMessage = t(
              'productDetail.featureDefinition.errors.nameRequired'
            )
          } else if (field === 'productKey') {
            errorMessage = t(
              'productDetail.featureDefinition.errors.productKeyRequired'
            )
          } else {
            errorMessage = `${firstError.path.join('.')}: ${firstError.message}`
          }
        } else {
          errorMessage = `${firstError.path.join('.')}: ${firstError.message}`
        }
        setError(
          `${t('productDetail.featureDefinition.status.validationFailed')}: ${errorMessage}`
        )
        setValidationError(errorMessage)
      } else if (e instanceof SyntaxError) {
        setError(
          `${t('productDetail.featureDefinition.errors.jsonSyntax')} ${t('productDetail.featureDefinition.status.validationFailed')}: ${e.message}`
        )
        setValidationError(
          t('productDetail.featureDefinition.errors.jsonSyntax')
        )
      } else {
        const message = e instanceof Error ? e.message : String(e)
        setError(
          `${t('productDetail.featureDefinition.buttons.format')} ${t('productDetail.featureDefinition.status.validationFailed')}: ${message}`
        )
        setValidationError(message)
      }
      setJsonValid(false)
    }
  }

  // 保存TSL
  const saveTSL = () => {
    const jsonToValidate =
      editMode === 'json' ? tslText : JSON.stringify(tslModel)
    if (validateTSL(jsonToValidate)) {
      // Use the correct model based on edit mode
      const modelToSave =
        editMode === 'json' ? normalizeTSLModel(JSON.parse(tslText)) : tslModel
      saveTSLMutation.mutate(modelToSave)
    }
  }

  // 清空编辑器
  const clearTSL = () => {
    const emptyModel = normalizeTSLModel(TSL_TEMPLATES.empty.value as TSLModel)
    setTslModel(emptyModel)
    syncToJson(emptyModel)
    setError('')
    setSuccess('')
  }

  // 添加/编辑属性
  const openPropertyDialog = (property?: Property, index?: number) => {
    if (property && index !== undefined) {
      setEditingProperty(property)
      setEditingPropertyIndex(index)
    } else {
      setEditingProperty({
        identifier: '',
        name: '',
        accessMode: 'rw',
        required: false,
        dataType: {
          type: 'int',
          specs: {},
        },
      })
      setEditingPropertyIndex(-1)
    }
    setPropertyDialogOpen(true)
  }

  const saveProperty = (property: Property) => {
    const properties = [...(tslModel.properties || [])]
    if (editingPropertyIndex >= 0) {
      properties[editingPropertyIndex] = property
    } else {
      properties.push(property)
    }
    const newModel = { ...tslModel, properties }
    setTslModel(newModel)
    syncToJson(newModel)
    setPropertyDialogOpen(false)
    setSuccess(
      editingPropertyIndex >= 0
        ? t('productDetail.featureDefinition.status.propertyUpdated')
        : t('productDetail.featureDefinition.status.propertyAdded')
    )
    setTimeout(() => setSuccess(''), 3000)
  }

  const deleteProperty = (index: number) => {
    const properties = [...(tslModel.properties || [])]
    properties.splice(index, 1)
    const newModel = { ...tslModel, properties }
    setTslModel(newModel)
    syncToJson(newModel)
    setSuccess(t('productDetail.featureDefinition.status.propertyDeleted'))
    setTimeout(() => setSuccess(''), 3000)
  }

  // 添加/编辑服务
  const openServiceDialog = (service?: Service, index?: number) => {
    if (service && index !== undefined) {
      setEditingService(service)
      setEditingServiceIndex(index)
    } else {
      setEditingService({
        identifier: '',
        name: '',
        required: false,
        callType: 'async',
        inputData: [],
        outputData: [],
      })
      setEditingServiceIndex(-1)
    }
    setServiceDialogOpen(true)
  }

  const saveService = (service: Service) => {
    const services = [...(tslModel.services || [])]
    if (editingServiceIndex >= 0) {
      services[editingServiceIndex] = service
    } else {
      services.push(service)
    }
    const newModel = { ...tslModel, services }
    setTslModel(newModel)
    syncToJson(newModel)
    setServiceDialogOpen(false)
    setSuccess(
      editingServiceIndex >= 0
        ? t('productDetail.featureDefinition.status.serviceUpdated')
        : t('productDetail.featureDefinition.status.serviceAdded')
    )
    setTimeout(() => setSuccess(''), 3000)
  }

  const deleteService = (index: number) => {
    const services = [...(tslModel.services || [])]
    services.splice(index, 1)
    const newModel = { ...tslModel, services }
    setTslModel(newModel)
    syncToJson(newModel)
    setSuccess(t('productDetail.featureDefinition.status.serviceDeleted'))
    setTimeout(() => setSuccess(''), 3000)
  }

  // 添加/编辑事件
  const openEventDialog = (event?: Event, index?: number) => {
    if (event && index !== undefined) {
      setEditingEvent(event)
      setEditingEventIndex(index)
    } else {
      setEditingEvent({
        identifier: '',
        name: '',
        type: 'info',
        required: false,
        outputData: [],
      })
      setEditingEventIndex(-1)
    }
    setEventDialogOpen(true)
  }

  const saveEvent = (event: Event) => {
    const events = [...(tslModel.events || [])]
    if (editingEventIndex >= 0) {
      events[editingEventIndex] = event
    } else {
      events.push(event)
    }
    const newModel = { ...tslModel, events }
    setTslModel(newModel)
    syncToJson(newModel)
    setEventDialogOpen(false)
    setSuccess(
      editingEventIndex >= 0
        ? t('productDetail.featureDefinition.status.eventUpdated')
        : t('productDetail.featureDefinition.status.eventAdded')
    )
    setTimeout(() => setSuccess(''), 3000)
  }

  const deleteEvent = (index: number) => {
    const events = [...(tslModel.events || [])]
    events.splice(index, 1)
    const newModel = { ...tslModel, events }
    setTslModel(newModel)
    syncToJson(newModel)
    setSuccess(t('productDetail.featureDefinition.status.eventDeleted'))
    setTimeout(() => setSuccess(''), 3000)
  }

  return (
    <div className='w-full max-w-full min-w-0 space-y-4'>
      {/* 头部说明 */}
      <div className='rounded-lg border bg-muted/50 p-4'>
        <div className='flex items-start space-x-3'>
          <FileJson className='mt-1 h-5 w-5 text-primary' />
          <div className='space-y-1'>
            <h3 className='text-sm font-semibold'>
              {t('productDetail.featureDefinition.title')}
            </h3>
            <p className='text-xs text-muted-foreground'>
              {t('productDetail.featureDefinition.description')}
            </p>
            <div className='flex flex-wrap gap-2 pt-2'>
              <Badge variant='outline' className='text-xs'>
                <Code className='mr-1 h-3 w-3' />
                {t('productDetail.featureDefinition.badge.jsonFormat')}
              </Badge>
              <Badge variant='outline' className='text-xs'>
                <Zap className='mr-1 h-3 w-3' />
                {t('productDetail.featureDefinition.badge.realtimeValidation')}
              </Badge>
              <Badge variant='outline' className='text-xs'>
                <Shield className='mr-1 h-3 w-3' />
                {t('productDetail.featureDefinition.badge.standardized')}
              </Badge>
            </div>
          </div>
        </div>
      </div>

      {/* 模板选择 */}
      <div className='flex items-center justify-between'>
        <div className='flex items-center space-x-2'>
          <span className='text-sm font-medium text-muted-foreground'>
            {t('productDetail.featureDefinition.templates.label')}
          </span>
          <div className='flex gap-2'>
            {Object.entries(TSL_TEMPLATES).map(([key, template]) => (
              <Button
                key={key}
                variant={currentTemplate === key ? 'default' : 'outline'}
                size='sm'
                onClick={() => loadTemplate(key)}
                className='h-7 text-xs'
              >
                {t(template.nameKey)}
              </Button>
            ))}
          </div>
        </div>

        <div className='flex gap-2'>
          {editMode === 'json' && (
            <Button
              variant='outline'
              size='sm'
              onClick={formatTSL}
              className='h-7 text-xs'
            >
              <RefreshCw className='mr-1 h-3 w-3' />
              {t('productDetail.featureDefinition.buttons.format')}
            </Button>
          )}
          <Button
            variant='outline'
            size='sm'
            onClick={clearTSL}
            className='h-7 text-xs'
          >
            {t('common:clear')}
          </Button>
        </div>
      </div>

      {/* 状态提示 */}
      {error && (
        <div className='flex items-center gap-2 rounded-md border border-destructive/50 bg-destructive/10 p-3'>
          <AlertCircle className='h-4 w-4 text-destructive' />
          <span className='text-sm text-destructive'>{error}</span>
        </div>
      )}

      {success && (
        <div className='flex items-center gap-2 rounded-md border border-green-500/50 bg-green-500/10 p-3'>
          <CheckCircle className='h-4 w-4 text-green-500' />
          <span className='text-sm text-green-600'>{success}</span>
        </div>
      )}

      {/* 编辑器 */}
      <Tabs
        value={editMode}
        onValueChange={(value) => {
          const next = value as 'visual' | 'json'
          setEditMode(next)
          if (next === 'json') {
            syncToJson(tslModel)
          } else {
            syncFromJson(tslText)
          }
        }}
        className='min-w-0'
      >
        <TabsList>
          <TabsTrigger value='visual'>
            <List className='mr-2 h-4 w-4' />
            {t('productDetail.featureDefinition.editMode.visual')}
          </TabsTrigger>
          <TabsTrigger value='json'>
            <Code className='mr-2 h-4 w-4' />
            {t('productDetail.featureDefinition.editMode.json')}
          </TabsTrigger>
        </TabsList>

        <TabsContent value='visual' className='min-w-0 space-y-4'>
          {isLoadingTSL && (
            <div className='flex items-center justify-center py-8'>
              <Loader2 className='h-6 w-6 animate-spin' />
              <span className='ml-2 text-sm text-muted-foreground'>
                {t('common:loading')}
              </span>
            </div>
          )}

          {tslError && !isLoadingTSL && (
            <div className='flex items-center gap-2 rounded-md border border-destructive/50 bg-destructive/10 p-3'>
              <AlertCircle className='h-4 w-4 text-destructive' />
              <span className='text-sm text-destructive'>
                {t('productDetail.featureDefinition.errors.loadFailed')}
              </span>
            </div>
          )}

          {/* 属性列表 */}
          <div className='rounded-lg border'>
            <div className='flex items-center justify-between border-b bg-muted/30 px-4 py-3'>
              <div className='flex items-center gap-2'>
                <h4 className='text-sm font-semibold'>
                  {t('productDetail.featureDefinition.sections.properties')}
                </h4>
                <Badge variant='secondary'>
                  {(tslModel?.properties || []).length}
                </Badge>
              </div>
              <Button size='sm' onClick={() => openPropertyDialog()}>
                <Plus className='mr-1 h-3 w-3' />
                {t('productDetail.featureDefinition.buttons.addProperty')}
              </Button>
            </div>
            <div className='p-4'>
              {(tslModel?.properties || []).length === 0 ? (
                <div className='py-8 text-center text-sm text-muted-foreground'>
                  {t('productDetail.featureDefinition.empty.properties')}
                </div>
              ) : (
                <div className='space-y-2'>
                  {(tslModel?.properties || []).map((property, index) => (
                    <div
                      key={index}
                      className='flex items-center justify-between rounded-md border p-3 hover:bg-muted/50'
                    >
                      <div className='flex-1'>
                        <div className='flex items-center gap-2'>
                          <span className='font-mono text-sm font-medium'>
                            {property.identifier}
                          </span>
                          <Badge
                            variant='outline'
                            className='font-mono text-xs uppercase'
                          >
                            {property.accessMode}
                          </Badge>
                          <Badge
                            variant='outline'
                            className='font-mono text-xs text-primary'
                          >
                            {property.dataType.type}
                          </Badge>
                        </div>
                        <p className='mt-1 text-xs text-muted-foreground'>
                          {property.name}
                        </p>
                      </div>
                      <div className='flex gap-2'>
                        <Button
                          variant='ghost'
                          size='sm'
                          onClick={() => openPropertyDialog(property, index)}
                        >
                          <Edit2 className='h-3 w-3' />
                        </Button>
                        <Button
                          variant='ghost'
                          size='sm'
                          onClick={() => deleteProperty(index)}
                        >
                          <Trash2 className='h-3 w-3' />
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* 服务列表 */}
          <div className='rounded-lg border'>
            <div className='flex items-center justify-between border-b bg-muted/30 px-4 py-3'>
              <div className='flex items-center gap-2'>
                <h4 className='text-sm font-semibold'>
                  {t('productDetail.featureDefinition.sections.services')}
                </h4>
                <Badge variant='secondary'>
                  {(tslModel?.services || []).length}
                </Badge>
              </div>
              <Button size='sm' onClick={() => openServiceDialog()}>
                <Plus className='mr-1 h-3 w-3' />
                {t('productDetail.featureDefinition.buttons.addService')}
              </Button>
            </div>
            <div className='p-4'>
              {(tslModel?.services || []).length === 0 ? (
                <div className='py-8 text-center text-sm text-muted-foreground'>
                  {t('productDetail.featureDefinition.empty.services')}
                </div>
              ) : (
                <div className='space-y-2'>
                  {(tslModel?.services || []).map((service, index) => (
                    <div
                      key={index}
                      className='flex items-center justify-between rounded-md border p-3 hover:bg-muted/50'
                    >
                      <div className='flex-1'>
                        <div className='flex items-center gap-2'>
                          <span className='font-mono text-sm font-medium'>
                            {service.identifier}
                          </span>
                          <Badge
                            variant='outline'
                            className='font-mono text-xs uppercase'
                          >
                            {service.callType}
                          </Badge>
                          <Badge variant='outline' className='text-xs'>
                            {t(
                              'productDetail.featureDefinition.statusBadge.input'
                            )}
                            : {service.inputData?.length ?? 0}
                          </Badge>
                          <Badge variant='outline' className='text-xs'>
                            {t(
                              'productDetail.featureDefinition.statusBadge.output'
                            )}
                            : {service.outputData?.length ?? 0}
                          </Badge>
                        </div>
                        <p className='mt-1 text-xs text-muted-foreground'>
                          {service.name}
                        </p>
                      </div>
                      <div className='flex gap-2'>
                        <Button
                          variant='ghost'
                          size='sm'
                          onClick={() => openServiceDialog(service, index)}
                        >
                          <Edit2 className='h-3 w-3' />
                        </Button>
                        <Button
                          variant='ghost'
                          size='sm'
                          onClick={() => deleteService(index)}
                        >
                          <Trash2 className='h-3 w-3' />
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* 事件列表 */}
          <div className='rounded-lg border'>
            <div className='flex items-center justify-between border-b bg-muted/30 px-4 py-3'>
              <div className='flex items-center gap-2'>
                <h4 className='text-sm font-semibold'>
                  {t('productDetail.featureDefinition.sections.events')}
                </h4>
                <Badge variant='secondary'>
                  {(tslModel?.events || []).length}
                </Badge>
              </div>
              <Button size='sm' onClick={() => openEventDialog()}>
                <Plus className='mr-1 h-3 w-3' />
                {t('productDetail.featureDefinition.buttons.addEvent')}
              </Button>
            </div>
            <div className='p-4'>
              {(tslModel?.events || []).length === 0 ? (
                <div className='py-8 text-center text-sm text-muted-foreground'>
                  {t('productDetail.featureDefinition.empty.events')}
                </div>
              ) : (
                <div className='space-y-2'>
                  {(tslModel?.events || []).map((event, index) => (
                    <div
                      key={index}
                      className='flex items-center justify-between rounded-md border p-3 hover:bg-muted/50'
                    >
                      <div className='flex-1'>
                        <div className='flex items-center gap-2'>
                          <span className='font-mono text-sm font-medium'>
                            {event.identifier}
                          </span>
                          <Badge
                            variant='outline'
                            className='font-mono text-xs uppercase'
                          >
                            {event.type}
                          </Badge>
                          <Badge variant='outline' className='text-xs'>
                            {t(
                              'productDetail.featureDefinition.statusBadge.output'
                            )}
                            : {event.outputData?.length ?? 0}
                          </Badge>
                        </div>
                        <p className='mt-1 text-xs text-muted-foreground'>
                          {event.name}
                        </p>
                      </div>
                      <div className='flex gap-2'>
                        <Button
                          variant='ghost'
                          size='sm'
                          onClick={() => openEventDialog(event, index)}
                        >
                          <Edit2 className='h-3 w-3' />
                        </Button>
                        <Button
                          variant='ghost'
                          size='sm'
                          onClick={() => deleteEvent(index)}
                        >
                          <Trash2 className='h-3 w-3' />
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </TabsContent>

        <TabsContent value='json' className='min-w-0'>
          <div className='rounded-lg border border-border bg-background shadow-sm'>
            <div className='flex items-center justify-between border-b border-border bg-muted/50 px-4 py-3'>
              <div className='flex items-center gap-3'>
                <Code className='h-4 w-4 text-muted-foreground' />
                <span className='text-sm font-medium'>
                  TSL {t('productDetail.featureDefinition.sections.properties')}
                </span>
                {/* JSON 验证状态指示器 */}
                {jsonValid === true && (
                  <Badge
                    variant='outline'
                    className='gap-1.5 border-green-500/30 bg-green-500/10 text-green-700 dark:text-green-400'
                  >
                    <CheckCheck className='h-3 w-3' />
                    {t('productDetail.featureDefinition.statusBadge.valid')}
                  </Badge>
                )}
                {jsonValid === false && (
                  <Badge
                    variant='outline'
                    className='gap-1.5 border-destructive/30 bg-destructive/10 text-destructive'
                  >
                    <XCircle className='h-3 w-3' />
                    {t('productDetail.featureDefinition.statusBadge.invalid')}
                  </Badge>
                )}
                {jsonValid === null && (
                  <Badge
                    variant='outline'
                    className='gap-1.5 text-muted-foreground'
                  >
                    <AlertCircle className='h-3 w-3' />
                    {t('productDetail.featureDefinition.statusBadge.unknown')}
                  </Badge>
                )}
              </div>
              <div className='flex items-center gap-2'>
                <Button
                  variant='outline'
                  size='sm'
                  onClick={formatTSL}
                  className='h-8 gap-1.5 text-xs'
                  disabled={jsonValid === false}
                >
                  <RefreshCw className='h-3.5 w-3.5' />
                  {t('productDetail.featureDefinition.buttons.format')}
                </Button>
                <Badge variant='secondary' className='font-mono text-xs'>
                  JSON
                </Badge>
              </div>
            </div>
            {/* 验证错误提示 */}
            {validationError && jsonValid === false && (
              <div className='border-b border-destructive/30 bg-destructive/5 px-4 py-2'>
                <div className='flex items-start gap-2'>
                  <AlertCircle className='mt-0.5 h-3.5 w-3.5 flex-shrink-0 text-destructive' />
                  <p className='text-xs text-destructive'>{validationError}</p>
                </div>
              </div>
            )}
            <div className='max-h-[600px] overflow-y-auto'>
              <Textarea
                value={tslText}
                onChange={(e) => {
                  const newValue = e.target.value
                  setTslText(newValue)
                  setError('')
                  setSuccess('')
                  // 实时验证
                  validateJsonSyntax(newValue)
                }}
                placeholder={`${t('productDetail.featureDefinition.title')}\n{\n  "schema": "schema.json",\n  "profile": {\n    "productKey": "PRODUCT_KEY"\n  },\n  "properties": [],\n  "events": [],\n  "services": []\n}`}
                className='min-h-[400px] resize-none border-0 bg-background font-mono text-sm leading-relaxed focus-visible:ring-0 focus-visible:ring-offset-0'
              />
            </div>
          </div>
        </TabsContent>
      </Tabs>

      {/* 保存按钮 */}
      <div className='flex justify-end gap-2'>
        {tslData && (
          <Button
            variant='outline'
            onClick={() => setDeleteDialogOpen(true)}
            disabled={deleteTSLMutation.isPending}
            className='min-w-[100px]'
          >
            {deleteTSLMutation.isPending ? (
              <Loader2 className='mr-2 h-4 w-4 animate-spin' />
            ) : null}
            {t('common:delete')}
          </Button>
        )}
        <Button
          onClick={saveTSL}
          disabled={saveTSLMutation.isPending}
          className='min-w-[100px]'
        >
          {saveTSLMutation.isPending ? (
            <Loader2 className='mr-2 h-4 w-4 animate-spin' />
          ) : null}
          {t('common:save')}
        </Button>
      </div>

      {/* 使用说明 */}
      <div className='rounded-lg border bg-muted/30 p-4'>
        <h4 className='mb-2 text-sm font-semibold'>
          {t('productDetail.featureDefinition.instructions.title')}
        </h4>
        <ul className='list-disc space-y-1 pl-5 text-xs text-muted-foreground'>
          <li>
            {t('productDetail.featureDefinition.instructions.properties')}
          </li>
          <li>{t('productDetail.featureDefinition.instructions.events')}</li>
          <li>{t('productDetail.featureDefinition.instructions.services')}</li>
          <li>{t('productDetail.featureDefinition.instructions.dualMode')}</li>
          <li>
            {t('productDetail.featureDefinition.instructions.visualMode')}
          </li>
        </ul>
      </div>

      {/* 属性编辑对话框 */}
      <PropertyDialog
        open={propertyDialogOpen}
        onOpenChange={setPropertyDialogOpen}
        property={editingProperty}
        onSave={saveProperty}
      />

      {/* 服务编辑对话框 */}
      <ServiceDialog
        open={serviceDialogOpen}
        onOpenChange={setServiceDialogOpen}
        service={editingService}
        onSave={saveService}
      />

      {/* 事件编辑对话框 */}
      <EventDialog
        open={eventDialogOpen}
        onOpenChange={setEventDialogOpen}
        event={editingEvent}
        onSave={saveEvent}
      />

      {/* 删除确认对话框 */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent className='max-w-md'>
          <DialogHeader>
            <DialogTitle>
              {t('productDetail.featureDefinition.deleteDialog.title')}
            </DialogTitle>
            <DialogDescription>
              {t('productDetail.featureDefinition.deleteDialog.description')}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant='outline'
              onClick={() => setDeleteDialogOpen(false)}
              disabled={deleteTSLMutation.isPending}
            >
              {t('common:cancel')}
            </Button>
            <Button
              variant='destructive'
              onClick={() => {
                deleteTSLMutation.mutate()
                setDeleteDialogOpen(false)
              }}
              disabled={deleteTSLMutation.isPending}
            >
              {deleteTSLMutation.isPending ? (
                <Loader2 className='mr-2 h-4 w-4 animate-spin' />
              ) : null}
              {t('common:delete')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

// 常用工程单位
const COMMON_UNITS = [
  { value: 'none', label: '无单位' },
  { value: '°C', label: '℃ (摄氏度)' },
  { value: '°F', label: '℉ (华氏度)' },
  { value: '%', label: '% (百分比)' },
  { value: 'V', label: 'V (伏特)' },
  { value: 'mV', label: 'mV (毫伏)' },
  { value: 'A', label: 'A (安培)' },
  { value: 'mA', label: 'mA (毫安)' },
  { value: 'W', label: 'W (瓦特)' },
  { value: 'kW', label: 'kW (千瓦)' },
  { value: 'kW·h', label: 'kW·h (千瓦·时)' },
  { value: 'Pa', label: 'Pa (帕斯卡)' },
  { value: 'kPa', label: 'kPa (千帕)' },
  { value: 'MPa', label: 'MPa (兆帕)' },
  { value: 'bar', label: 'bar (巴)' },
  { value: 'm/s', label: 'm/s (米/秒)' },
  { value: 'km/h', label: 'km/h (千米/时)' },
  { value: 'rpm', label: 'rpm (转/分)' },
  { value: 'lx', label: 'lx (勒克斯)' },
  { value: 'dB', label: 'dB (分贝)' },
  { value: 'ppm', label: 'ppm (百万分率)' },
  { value: 'mg/L', label: 'mg/L (毫克/升)' },
  { value: 'μg/m³', label: 'μg/m³ (微克/立方米)' },
  { value: 's', label: 's (秒)' },
  { value: 'min', label: 'min (分钟)' },
  { value: 'h', label: 'h (小时)' },
  { value: 'custom', label: '自定义单位...' },
]

// 属性编辑对话框组件
function PropertyDialog({
  open,
  onOpenChange,
  property,
  onSave,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  property: Property | null
  onSave: (property: Property) => void
}) {
  const { t } = useTranslation(['deviceManagement', 'common'])
  const [formData, setFormData] = useState<Property | null>(null)
  const [isCustomUnit, setIsCustomUnit] = useState(false)

  // 只在对话框打开时初始化表单数据
  useEffect(() => {
    if (open && property) {
      const cloned = JSON.parse(JSON.stringify(property))
      setFormData(cloned)
      const currentUnit = cloned.dataType?.specs?.unit
      const isInList = COMMON_UNITS.some(
        (u) => u.value === currentUnit || (u.value === 'none' && !currentUnit)
      )
      setIsCustomUnit(Boolean(currentUnit && !isInList))
    }
  }, [open, property])

  const handleSave = () => {
    if (formData && formData.identifier && formData.name) {
      onSave(formData)
    }
  }

  if (!formData) return null

  const isEdit = Boolean(property && property.identifier)
  const currentUnitValue = isCustomUnit
    ? 'custom'
    : formData.dataType?.specs?.unit || 'none'

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='max-h-[90vh] max-w-lg overflow-y-auto p-5'>
        <DialogHeader className='pb-1'>
          <DialogTitle className='text-base font-semibold'>
            {isEdit
              ? t(
                  'productDetail.featureDefinition.propertyDialog.editTitle',
                  '编辑自定义功能'
                )
              : t(
                  'productDetail.featureDefinition.propertyDialog.addTitle',
                  '添加自定义功能'
                )}
          </DialogTitle>
        </DialogHeader>

        <div className='space-y-3 py-1 text-xs'>
          {/* 功能类型 */}
          <div className='space-y-1.5'>
            <div className='flex items-center gap-1.5'>
              <Label className='flex items-center gap-1 text-xs font-medium'>
                <span className='text-destructive'>*</span>
                {t(
                  'productDetail.featureDefinition.propertyDialog.featureType',
                  '功能类型'
                )}
              </Label>
              <Tooltip>
                <TooltipTrigger asChild>
                  <HelpCircle className='size-3.5 cursor-pointer text-muted-foreground/70 transition-colors hover:text-foreground' />
                </TooltipTrigger>
                <TooltipContent side='top'>
                  <span>
                    {t(
                      'productDetail.featureDefinition.propertyDialog.featureTypeTooltip',
                      '支持属性、服务、事件三种物模型功能类型'
                    )}
                  </span>
                </TooltipContent>
              </Tooltip>
            </div>
            <div className='inline-flex h-8 rounded-sm border border-input bg-muted/20 p-0.5 text-xs'>
              <button
                type='button'
                className='rounded-xs border border-primary bg-primary/5 px-3.5 font-medium text-primary shadow-2xs'
              >
                {t(
                  'productDetail.featureDefinition.sections.propertiesOnly',
                  '属性'
                )}
              </button>
              <button
                type='button'
                disabled
                className='cursor-not-allowed px-3.5 text-muted-foreground/50'
              >
                {t(
                  'productDetail.featureDefinition.sections.servicesOnly',
                  '服务'
                )}
              </button>
              <button
                type='button'
                disabled
                className='cursor-not-allowed px-3.5 text-muted-foreground/50'
              >
                {t(
                  'productDetail.featureDefinition.sections.eventsOnly',
                  '事件'
                )}
              </button>
            </div>
          </div>

          {/* 功能名称 */}
          <div className='space-y-1.5'>
            <div className='flex items-center gap-1.5'>
              <Label className='flex items-center gap-1 text-xs font-medium'>
                <span className='text-destructive'>*</span>
                {t(
                  'productDetail.featureDefinition.propertyDialog.name',
                  '功能名称'
                )}
              </Label>
              <Tooltip>
                <TooltipTrigger asChild>
                  <HelpCircle className='size-3.5 cursor-pointer text-muted-foreground/70 transition-colors hover:text-foreground' />
                </TooltipTrigger>
                <TooltipContent side='top'>
                  <span>
                    {t(
                      'productDetail.featureDefinition.propertyDialog.nameTooltip',
                      '功能展示名称，支持中英文、数字和下划线'
                    )}
                  </span>
                </TooltipContent>
              </Tooltip>
            </div>
            <Input
              className='h-8 text-xs'
              value={formData.name}
              onChange={(e) =>
                setFormData({ ...formData, name: e.target.value })
              }
              placeholder={t(
                'productDetail.featureDefinition.placeholders.propertyName',
                '请输入您的功能名称'
              )}
            />
          </div>

          {/* 标识符 */}
          <div className='space-y-1.5'>
            <div className='flex items-center gap-1.5'>
              <Label className='flex items-center gap-1 text-xs font-medium'>
                <span className='text-destructive'>*</span>
                {t(
                  'productDetail.featureDefinition.propertyDialog.identifier',
                  '标识符'
                )}
              </Label>
              <Tooltip>
                <TooltipTrigger asChild>
                  <HelpCircle className='size-3.5 cursor-pointer text-muted-foreground/70 transition-colors hover:text-foreground' />
                </TooltipTrigger>
                <TooltipContent side='top'>
                  <span>
                    {t(
                      'productDetail.featureDefinition.propertyDialog.identifierTooltip',
                      '属性唯一标识符，英文、数字和下划线，首字符为英文字母'
                    )}
                  </span>
                </TooltipContent>
              </Tooltip>
            </div>
            <Input
              className='h-8 font-mono text-xs'
              value={formData.identifier}
              onChange={(e) =>
                setFormData({ ...formData, identifier: e.target.value })
              }
              placeholder={t(
                'productDetail.featureDefinition.placeholders.propertyId',
                '请输入您的标识符'
              )}
            />
          </div>

          {/* 数据类型 */}
          <div className='space-y-1.5'>
            <Label className='flex items-center gap-1 text-xs font-medium'>
              <span className='text-destructive'>*</span>
              {t(
                'productDetail.featureDefinition.propertyDialog.dataType',
                '数据类型'
              )}
            </Label>
            <Select
              value={formData.dataType.type}
              onValueChange={(
                value:
                  | 'int'
                  | 'float'
                  | 'double'
                  | 'bool'
                  | 'string'
                  | 'enum'
                  | 'struct'
                  | 'array'
              ) =>
                setFormData({
                  ...formData,
                  dataType: { ...formData.dataType, type: value },
                })
              }
            >
              <SelectTrigger className='h-8 w-full text-xs'>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value='int'>int32 (整数型)</SelectItem>
                <SelectItem value='float'>float (单精度浮点型)</SelectItem>
                <SelectItem value='double'>double (双精度浮点型)</SelectItem>
                <SelectItem value='bool'>bool (布尔型)</SelectItem>
                <SelectItem value='string'>text (字符串)</SelectItem>
                <SelectItem value='enum'>enum (枚举型)</SelectItem>
                <SelectItem value='struct'>struct (结构体)</SelectItem>
                <SelectItem value='array'>array (数组)</SelectItem>
              </SelectContent>
            </Select>
          </div>

          {/* 数值类型的规格：取值范围、步长、单位 */}
          {['int', 'float', 'double'].includes(formData.dataType.type) && (
            <>
              {/* 取值范围 */}
              <div className='space-y-1.5'>
                <Label className='text-xs font-medium'>
                  {t(
                    'productDetail.featureDefinition.propertyDialog.range',
                    '取值范围'
                  )}
                </Label>
                <div className='flex items-center gap-2'>
                  <Input
                    className='h-8 text-xs'
                    placeholder={t(
                      'productDetail.featureDefinition.propertyDialog.min',
                      '最小值'
                    )}
                    value={formData.dataType.specs?.min ?? ''}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        dataType: {
                          ...formData.dataType,
                          specs: {
                            ...formData.dataType.specs,
                            min: e.target.value,
                          },
                        },
                      })
                    }
                  />
                  <span className='text-xs text-muted-foreground select-none'>
                    ~
                  </span>
                  <Input
                    className='h-8 text-xs'
                    placeholder={t(
                      'productDetail.featureDefinition.propertyDialog.max',
                      '最大值'
                    )}
                    value={formData.dataType.specs?.max ?? ''}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        dataType: {
                          ...formData.dataType,
                          specs: {
                            ...formData.dataType.specs,
                            max: e.target.value,
                          },
                        },
                      })
                    }
                  />
                </div>
              </div>

              {/* 步长 */}
              <div className='space-y-1.5'>
                <Label className='text-xs font-medium'>
                  {t(
                    'productDetail.featureDefinition.propertyDialog.step',
                    '步长'
                  )}
                </Label>
                <Input
                  className='h-8 text-xs'
                  placeholder={t(
                    'productDetail.featureDefinition.placeholders.step',
                    '请输入步长'
                  )}
                  value={formData.dataType.specs?.step ?? ''}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      dataType: {
                        ...formData.dataType,
                        specs: {
                          ...formData.dataType.specs,
                          step: e.target.value,
                        },
                      },
                    })
                  }
                />
              </div>

              {/* 单位 */}
              <div className='space-y-1.5'>
                <Label className='text-xs font-medium'>
                  {t(
                    'productDetail.featureDefinition.propertyDialog.unit',
                    '单位'
                  )}
                </Label>
                <Select
                  value={currentUnitValue}
                  onValueChange={(val) => {
                    if (val === 'custom') {
                      setIsCustomUnit(true)
                    } else {
                      setIsCustomUnit(false)
                      setFormData({
                        ...formData,
                        dataType: {
                          ...formData.dataType,
                          specs: {
                            ...formData.dataType.specs,
                            unit: val === 'none' ? '' : val,
                          },
                        },
                      })
                    }
                  }}
                >
                  <SelectTrigger className='h-8 w-full text-xs'>
                    <SelectValue
                      placeholder={t(
                        'productDetail.featureDefinition.placeholders.unit',
                        '请选择单位'
                      )}
                    />
                  </SelectTrigger>
                  <SelectContent className='max-h-56'>
                    {COMMON_UNITS.map((unit) => (
                      <SelectItem key={unit.value} value={unit.value}>
                        {unit.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {isCustomUnit && (
                  <div className='flex gap-1.5 pt-1'>
                    <Input
                      placeholder={t(
                        'productDetail.featureDefinition.placeholders.customUnit',
                        '请输入自定义单位'
                      )}
                      value={formData.dataType.specs?.unit ?? ''}
                      onChange={(e) =>
                        setFormData({
                          ...formData,
                          dataType: {
                            ...formData.dataType,
                            specs: {
                              ...formData.dataType.specs,
                              unit: e.target.value,
                            },
                          },
                        })
                      }
                      className='h-8 flex-1 text-xs'
                    />
                    <Button
                      type='button'
                      variant='ghost'
                      size='sm'
                      className='h-8 px-2 text-xs text-muted-foreground'
                      onClick={() => {
                        setIsCustomUnit(false)
                        setFormData({
                          ...formData,
                          dataType: {
                            ...formData.dataType,
                            specs: {
                              ...formData.dataType.specs,
                              unit: '',
                            },
                          },
                        })
                      }}
                    >
                      {t('common:cancel')}
                    </Button>
                  </div>
                )}
              </div>
            </>
          )}

          {/* 枚举类型规格 */}
          {formData.dataType.type === 'enum' && (
            <div className='space-y-1.5'>
              <Label className='text-xs font-medium'>
                {t(
                  'productDetail.featureDefinition.propertyDialog.enumValues',
                  '枚举项'
                )}
              </Label>
              <Textarea
                placeholder='{"0": "关", "1": "开"}'
                value={JSON.stringify(formData.dataType.specs || {}, null, 2)}
                onChange={(e) => {
                  try {
                    const specs = JSON.parse(e.target.value)
                    setFormData({
                      ...formData,
                      dataType: { ...formData.dataType, specs },
                    })
                  } catch {
                    // 允许输入过程中的临时非标准JSON
                  }
                }}
                className='min-h-[64px] font-mono text-xs'
              />
            </div>
          )}

          {/* 结构体规格 */}
          {formData.dataType.type === 'struct' && (
            <div className='space-y-1.5'>
              <Label className='text-xs font-medium'>
                {t(
                  'productDetail.featureDefinition.propertyDialog.structDef',
                  '结构体定义'
                )}
              </Label>
              <Textarea
                placeholder='{"lng": {"type": "double"}, "lat": {"type": "double"}}'
                value={JSON.stringify(formData.dataType.specs || {}, null, 2)}
                onChange={(e) => {
                  try {
                    const specs = JSON.parse(e.target.value)
                    setFormData({
                      ...formData,
                      dataType: { ...formData.dataType, specs },
                    })
                  } catch {
                    // 忽略 JSON 解析错误
                  }
                }}
                className='min-h-[72px] font-mono text-xs'
              />
              <p className='text-[11px] text-muted-foreground'>
                {t(
                  'productDetail.featureDefinition.propertyDialog.structDefDesc',
                  '定义结构体的字段，每个字段包含 type 等属性'
                )}
              </p>
            </div>
          )}

          {/* 数组规格 */}
          {formData.dataType.type === 'array' && (
            <div className='space-y-1.5'>
              <Label className='text-xs font-medium'>
                {t(
                  'productDetail.featureDefinition.propertyDialog.arrayElemType',
                  '数组元素类型'
                )}
              </Label>
              <Textarea
                placeholder='{"type": "int", "size": 128}'
                value={JSON.stringify(formData.dataType.specs || {}, null, 2)}
                onChange={(e) => {
                  try {
                    const specs = JSON.parse(e.target.value)
                    setFormData({
                      ...formData,
                      dataType: { ...formData.dataType, specs },
                    })
                  } catch {
                    // 忽略 JSON 解析错误
                  }
                }}
                className='min-h-[72px] font-mono text-xs'
              />
              <p className='text-[11px] text-muted-foreground'>
                {t(
                  'productDetail.featureDefinition.propertyDialog.arrayElemDesc',
                  '定义数组元素的类型和大小限制'
                )}
              </p>
            </div>
          )}

          {/* 读写类型 */}
          <div className='space-y-1.5'>
            <Label className='flex items-center gap-1 text-xs font-medium'>
              <span className='text-destructive'>*</span>
              {t(
                'productDetail.featureDefinition.propertyDialog.accessMode',
                '读写类型'
              )}
            </Label>
            <RadioGroup
              value={formData.accessMode}
              onValueChange={(value: 'r' | 'rw') =>
                setFormData({ ...formData, accessMode: value })
              }
              className='flex items-center gap-8 pt-0.5'
            >
              <div className='flex items-center space-x-2'>
                <RadioGroupItem value='rw' id='access-rw' />
                <Label
                  htmlFor='access-rw'
                  className='cursor-pointer text-xs font-normal'
                >
                  {t('productDetail.featureDefinition.accessMode.rw', '读写')}
                </Label>
              </div>
              <div className='flex items-center space-x-2'>
                <RadioGroupItem value='r' id='access-r' />
                <Label
                  htmlFor='access-r'
                  className='cursor-pointer text-xs font-normal'
                >
                  {t('productDetail.featureDefinition.accessMode.r', '只读')}
                </Label>
              </div>
            </RadioGroup>
          </div>

          {/* 描述 */}
          <div className='space-y-1.5'>
            <Label className='text-xs font-medium'>
              {t(
                'productDetail.featureDefinition.propertyDialog.descriptionLabel',
                '描述'
              )}
            </Label>
            <div className='relative'>
              <Textarea
                placeholder={t(
                  'productDetail.featureDefinition.placeholders.desc',
                  '请输入描述'
                )}
                maxLength={100}
                rows={2}
                value={formData.desc || formData.description || ''}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    desc: e.target.value,
                    description: e.target.value,
                  })
                }
                className='h-[56px] min-h-[56px] resize-none pb-5 text-xs'
              />
              <div className='pointer-events-none absolute right-2.5 bottom-1 text-[11px] text-muted-foreground select-none'>
                {(formData.desc || formData.description || '').length}/100
              </div>
            </div>
          </div>
        </div>

        <DialogFooter className='flex justify-end gap-2 pt-1'>
          <Button size='sm' onClick={handleSave} className='h-8 text-xs'>
            {t('common:confirm')}
          </Button>
          <Button
            size='sm'
            variant='outline'
            onClick={() => onOpenChange(false)}
            className='h-8 text-xs'
          >
            {t('common:cancel')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

// 服务编辑对话框组件
function ServiceDialog({
  open,
  onOpenChange,
  service,
  onSave,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  service: Service | null
  onSave: (service: Service) => void
}) {
  const { t } = useTranslation(['deviceManagement', 'common'])
  const [formData, setFormData] = useState<Service | null>(null)

  // 只在对话框打开时初始化表单数据
  useEffect(() => {
    if (open && service) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setFormData(JSON.parse(JSON.stringify(service)))
    }
  }, [open, service])

  const handleSave = () => {
    if (formData && formData.identifier && formData.name) {
      onSave(formData)
    }
  }

  if (!formData) return null

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='max-h-[90vh] max-w-3xl overflow-hidden'>
        <DialogHeader>
          <DialogTitle>
            {t('productDetail.featureDefinition.serviceDialog.title', {
              isEdit: formData.identifier
                ? t('productDetail.featureDefinition.serviceDialog.edit')
                : t('productDetail.featureDefinition.serviceDialog.add'),
            })}
          </DialogTitle>
          <DialogDescription>
            {t('productDetail.featureDefinition.serviceDialog.description')}
          </DialogDescription>
        </DialogHeader>
        <div className='max-h-[calc(90vh-200px)] overflow-y-auto'>
          <div className='space-y-4'>
            <div className='grid grid-cols-2 gap-4'>
              <div className='space-y-2'>
                <Label>
                  {t(
                    'productDetail.featureDefinition.serviceDialog.identifier'
                  )}
                </Label>
                <Input
                  value={formData.identifier}
                  onChange={(e) =>
                    setFormData({ ...formData, identifier: e.target.value })
                  }
                  placeholder={t(
                    'productDetail.featureDefinition.placeholders.serviceId'
                  )}
                />
              </div>
              <div className='space-y-2'>
                <Label>
                  {t('productDetail.featureDefinition.serviceDialog.name')}
                </Label>
                <Input
                  value={formData.name}
                  onChange={(e) =>
                    setFormData({ ...formData, name: e.target.value })
                  }
                  placeholder={t(
                    'productDetail.featureDefinition.placeholders.serviceName'
                  )}
                />
              </div>
            </div>
            <div className='space-y-2'>
              <Label>
                {t('productDetail.featureDefinition.serviceDialog.callType')}
              </Label>
              <Select
                value={formData.callType}
                onValueChange={(value: 'async' | 'sync') =>
                  setFormData({ ...formData, callType: value })
                }
              >
                <SelectTrigger className='w-full'>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value='async'>
                    {t('productDetail.featureDefinition.serviceDialog.async')}
                  </SelectItem>
                  <SelectItem value='sync'>
                    {t('productDetail.featureDefinition.serviceDialog.sync')}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className='space-y-2'>
              <Label>
                {t('productDetail.featureDefinition.serviceDialog.inputParams')}
              </Label>
              <Textarea
                placeholder='[{"identifier": "speed", "name": "Speed", "dataType": {"type": "int"}}]'
                value={JSON.stringify(formData.inputData, null, 2)}
                onChange={(e) => {
                  try {
                    const inputData = JSON.parse(e.target.value)
                    setFormData({ ...formData, inputData })
                  } catch {
                    // 忽略 JSON 解析错误，用户可能还在输入
                  }
                }}
                className='min-h-[120px] font-mono text-sm'
              />
            </div>

            <div className='space-y-2'>
              <Label>
                {t(
                  'productDetail.featureDefinition.serviceDialog.outputParams'
                )}
              </Label>
              <Textarea
                placeholder='[{"identifier": "result", "name": "Result", "dataType": {"type": "bool"}}]'
                value={JSON.stringify(formData.outputData, null, 2)}
                onChange={(e) => {
                  try {
                    const outputData = JSON.parse(e.target.value)
                    setFormData({ ...formData, outputData })
                  } catch {
                    // 忽略 JSON 解析错误，用户可能还在输入
                  }
                }}
                className='min-h-[120px] font-mono text-sm'
              />
            </div>
          </div>
        </div>
        <DialogFooter>
          <Button variant='outline' onClick={() => onOpenChange(false)}>
            {t('common:cancel')}
          </Button>
          <Button onClick={handleSave}>{t('common:confirm')}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

// 事件编辑对话框组件
function EventDialog({
  open,
  onOpenChange,
  event,
  onSave,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  event: Event | null
  onSave: (event: Event) => void
}) {
  const { t } = useTranslation(['deviceManagement', 'common'])
  const [formData, setFormData] = useState<Event | null>(null)

  // 只在对话框打开时初始化表单数据
  useEffect(() => {
    if (open && event) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setFormData(JSON.parse(JSON.stringify(event)))
    }
  }, [open, event])

  const handleSave = () => {
    if (formData && formData.identifier && formData.name) {
      onSave(formData)
    }
  }

  if (!formData) return null

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='max-h-[90vh] max-w-3xl overflow-hidden'>
        <DialogHeader>
          <DialogTitle>
            {t('productDetail.featureDefinition.eventDialog.title', {
              isEdit: formData.identifier
                ? t('productDetail.featureDefinition.eventDialog.edit')
                : t('productDetail.featureDefinition.eventDialog.add'),
            })}
          </DialogTitle>
          <DialogDescription>
            {t('productDetail.featureDefinition.eventDialog.description')}
          </DialogDescription>
        </DialogHeader>
        <div className='max-h-[calc(90vh-200px)] overflow-y-auto'>
          <div className='space-y-4'>
            <div className='grid grid-cols-2 gap-4'>
              <div className='space-y-2'>
                <Label>
                  {t('productDetail.featureDefinition.eventDialog.identifier')}
                </Label>
                <Input
                  value={formData.identifier}
                  onChange={(e) =>
                    setFormData({ ...formData, identifier: e.target.value })
                  }
                  placeholder={t(
                    'productDetail.featureDefinition.placeholders.eventId'
                  )}
                />
              </div>
              <div className='space-y-2'>
                <Label>
                  {t('productDetail.featureDefinition.eventDialog.name')}
                </Label>
                <Input
                  value={formData.name}
                  onChange={(e) =>
                    setFormData({ ...formData, name: e.target.value })
                  }
                  placeholder={t(
                    'productDetail.featureDefinition.placeholders.eventName'
                  )}
                />
              </div>
            </div>
            <div className='space-y-2'>
              <Label>
                {t('productDetail.featureDefinition.eventDialog.eventType')}
              </Label>
              <Select
                value={formData.type}
                onValueChange={(value: 'info' | 'alert' | 'error') =>
                  setFormData({ ...formData, type: value })
                }
              >
                <SelectTrigger className='w-full'>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value='info'>
                    {t('productDetail.featureDefinition.eventTypes.info')}
                  </SelectItem>
                  <SelectItem value='alert'>
                    {t('productDetail.featureDefinition.eventTypes.alert')}
                  </SelectItem>
                  <SelectItem value='error'>
                    {t('productDetail.featureDefinition.eventTypes.error')}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className='space-y-2'>
              <Label>
                {t('productDetail.featureDefinition.eventDialog.outputData')}
              </Label>
              <Textarea
                placeholder='[{"identifier": "value", "name": "Value", "dataType": {"type": "int"}}]'
                value={JSON.stringify(formData.outputData, null, 2)}
                onChange={(e) => {
                  try {
                    const outputData = JSON.parse(e.target.value)
                    setFormData({ ...formData, outputData })
                  } catch {
                    // 忽略 JSON 解析错误，用户可能还在输入
                  }
                }}
                className='min-h-[150px] font-mono text-sm'
              />
            </div>
          </div>
        </div>
        <DialogFooter>
          <Button variant='outline' onClick={() => onOpenChange(false)}>
            {t('common:cancel')}
          </Button>
          <Button onClick={handleSave}>{t('common:confirm')}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
