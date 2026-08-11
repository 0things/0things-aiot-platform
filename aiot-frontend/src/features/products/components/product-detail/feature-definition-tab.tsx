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
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { productTSLServiceApi } from '@/api/clients'
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { Textarea } from '@/components/ui/textarea'

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
  | 'int'
  | 'float'
  | 'double'
  | 'bool'
  | 'string'
  | 'enum'
  | 'struct'
  | 'array'
  specs?: DataTypeSpecs
}

interface Property {
  identifier: string
  name: string
  accessMode: 'r' | 'rw'
  required: boolean
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
  schema: z.string(),
  profile: z.object({
    productKey: z.string().min(1),
  }),
  properties: z.array(propertySchema),
  events: z.array(eventSchema),
  services: z.array(serviceSchema),
})

interface FeatureDefinitionTabProps {
  productKey: string
}

export function FeatureDefinitionTab({
  productKey,
}: FeatureDefinitionTabProps) {
  const { t } = useTranslation('deviceManagement')
  const queryClient = useQueryClient()

  // Query for fetching TSL data
  const {
    data: tslData,
    isLoading: isLoadingTSL,
    error: tslError,
  } = useQuery({
    queryKey: ['product-tsl', productKey],
    queryFn: async () => {
      try {
        const response =
          await productTSLServiceApi.productTSLServiceGetProductTSL({
            productKey,
          })
        // Parse the TSL string from API into TSLModel
        if (response.data.productTsl?.tsl) {
          return JSON.parse(response.data.productTsl.tsl) as TSLModel
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
      const response =
        await productTSLServiceApi.productTSLServiceCreateProductTSL({
          productKey,
          productTslV1CreateProductTSLRequest: { tsl: tslJsonString },
        })
      return response.data
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
      const response =
        await productTSLServiceApi.productTSLServiceDeleteProductTSL({
          productKey,
        })
      return response.data
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
    TSL_TEMPLATES.empty.value as TSLModel
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

  // Sync fetched TSL data to editor state
  useEffect(() => {
    if (tslData) {
      setTslModel(tslData)
      syncToJson(tslData)
      setError('')
      setSuccess('')
      setCurrentTemplate('empty') // Reset template selection
    } else if (!isLoadingTSL && !tslError) {
      // No TSL data exists, use empty template
      const emptyModel = TSL_TEMPLATES.empty.value as TSLModel
      setTslModel(emptyModel)
      syncToJson(emptyModel)
      setCurrentTemplate('empty')
    }
  }, [tslData, isLoadingTSL, tslError])

  // 从 JSON 更新模型
  const syncFromJson = (text: string) => {
    try {
      const parsed = JSON.parse(text)
      setTslModel(parsed)
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
    const formatted = JSON.stringify(model, null, 2)
    setTslText(formatted)
  }

  // 当编辑模式切换时同步数据
  useEffect(() => {
    if (editMode === 'json') {
      syncToJson(tslModel)
    } else {
      syncFromJson(tslText)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [editMode])

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
      const newModel = template.value as TSLModel
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

      const formatted = JSON.stringify(parsed, null, 2)
      setTslText(formatted)
      setTslModel(parsed)
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
      const modelToSave = editMode === 'json' ? JSON.parse(tslText) : tslModel
      saveTSLMutation.mutate(modelToSave)
    }
  }

  // 清空编辑器
  const clearTSL = () => {
    const emptyModel = TSL_TEMPLATES.empty.value as TSLModel
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
        accessMode: 'r',
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
    const newModel = { ...tslModel }
    if (editingPropertyIndex >= 0) {
      newModel.properties[editingPropertyIndex] = property
    } else {
      newModel.properties.push(property)
    }
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
    const newModel = { ...tslModel }
    newModel.properties.splice(index, 1)
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
    const newModel = { ...tslModel }
    if (editingServiceIndex >= 0) {
      newModel.services[editingServiceIndex] = service
    } else {
      newModel.services.push(service)
    }
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
    const newModel = { ...tslModel }
    newModel.services.splice(index, 1)
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
    const newModel = { ...tslModel }
    if (editingEventIndex >= 0) {
      newModel.events[editingEventIndex] = event
    } else {
      newModel.events.push(event)
    }
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
    const newModel = { ...tslModel }
    newModel.events.splice(index, 1)
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
            {t('productDetail.featureDefinition.buttons.clear')}
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
        onValueChange={(value) => setEditMode(value as 'visual' | 'json')}
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
                {t('productDetail.featureDefinition.status.loading')}
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
                <Badge variant='secondary'>{tslModel.properties.length}</Badge>
              </div>
              <Button size='sm' onClick={() => openPropertyDialog()}>
                <Plus className='mr-1 h-3 w-3' />
                {t('productDetail.featureDefinition.buttons.addProperty')}
              </Button>
            </div>
            <div className='p-4'>
              {tslModel.properties.length === 0 ? (
                <div className='py-8 text-center text-sm text-muted-foreground'>
                  {t('productDetail.featureDefinition.empty.properties')}
                </div>
              ) : (
                <div className='space-y-2'>
                  {tslModel.properties.map((property, index) => (
                    <div
                      key={index}
                      className='flex items-center justify-between rounded-md border p-3 hover:bg-muted/50'
                    >
                      <div className='flex-1'>
                        <div className='flex items-center gap-2'>
                          <span className='font-mono text-sm font-medium'>
                            {property.identifier}
                          </span>
                          <Badge variant='outline' className='text-xs'>
                            {t(
                              `featureDefinition.accessMode.${property.accessMode}`
                            )}
                          </Badge>
                          <Badge variant='outline' className='text-xs'>
                            {t(
                              `featureDefinition.dataTypes.${property.dataType.type}`
                            )}
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
                <Badge variant='secondary'>{tslModel.services.length}</Badge>
              </div>
              <Button size='sm' onClick={() => openServiceDialog()}>
                <Plus className='mr-1 h-3 w-3' />
                {t('productDetail.featureDefinition.buttons.addService')}
              </Button>
            </div>
            <div className='p-4'>
              {tslModel.services.length === 0 ? (
                <div className='py-8 text-center text-sm text-muted-foreground'>
                  {t('productDetail.featureDefinition.empty.services')}
                </div>
              ) : (
                <div className='space-y-2'>
                  {tslModel.services.map((service, index) => (
                    <div
                      key={index}
                      className='flex items-center justify-between rounded-md border p-3 hover:bg-muted/50'
                    >
                      <div className='flex-1'>
                        <div className='flex items-center gap-2'>
                          <span className='font-mono text-sm font-medium'>
                            {service.identifier}
                          </span>
                          <Badge variant='outline' className='text-xs'>
                            {t(
                              `featureDefinition.serviceDialog.${service.callType}`
                            )}
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
                <Badge variant='secondary'>{tslModel.events.length}</Badge>
              </div>
              <Button size='sm' onClick={() => openEventDialog()}>
                <Plus className='mr-1 h-3 w-3' />
                {t('productDetail.featureDefinition.buttons.addEvent')}
              </Button>
            </div>
            <div className='p-4'>
              {tslModel.events.length === 0 ? (
                <div className='py-8 text-center text-sm text-muted-foreground'>
                  {t('productDetail.featureDefinition.empty.events')}
                </div>
              ) : (
                <div className='space-y-2'>
                  {tslModel.events.map((event, index) => (
                    <div
                      key={index}
                      className='flex items-center justify-between rounded-md border p-3 hover:bg-muted/50'
                    >
                      <div className='flex-1'>
                        <div className='flex items-center gap-2'>
                          <span className='font-mono text-sm font-medium'>
                            {event.identifier}
                          </span>
                          <Badge variant='outline' className='text-xs'>
                            {t(`featureDefinition.eventTypes.${event.type}`)}
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
            {t('productDetail.featureDefinition.buttons.delete')}
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
          {t('productDetail.featureDefinition.buttons.save')}
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
              {t('productDetail.featureDefinition.buttons.cancel')}
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
              {t('productDetail.featureDefinition.buttons.confirmDelete')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

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
  const { t } = useTranslation('deviceManagement')
  const [formData, setFormData] = useState<Property | null>(null)

  // 只在对话框打开时初始化表单数据
  useEffect(() => {
    if (open && property) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setFormData(JSON.parse(JSON.stringify(property)))
    }
  }, [open, property])

  const handleSave = () => {
    if (formData && formData.identifier && formData.name) {
      onSave(formData)
    }
  }

  if (!formData) return null

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='max-h-[90vh] max-w-2xl overflow-hidden'>
        <DialogHeader>
          <DialogTitle>
            {t('productDetail.featureDefinition.propertyDialog.title', {
              isEdit: formData.identifier
                ? t('productDetail.featureDefinition.propertyDialog.edit')
                : t('productDetail.featureDefinition.propertyDialog.add'),
            })}
          </DialogTitle>
          <DialogDescription>
            {t('productDetail.featureDefinition.propertyDialog.description')}
          </DialogDescription>
        </DialogHeader>
        <div className='max-h-[calc(90vh-200px)] overflow-y-auto'>
          <div className='space-y-4'>
            <div className='grid grid-cols-2 gap-4'>
              <div className='space-y-2'>
                <Label>
                  {t(
                    'productDetail.featureDefinition.propertyDialog.identifier'
                  )}
                </Label>
                <Input
                  value={formData.identifier}
                  onChange={(e) =>
                    setFormData({ ...formData, identifier: e.target.value })
                  }
                  placeholder={t(
                    'productDetail.featureDefinition.placeholders.propertyId'
                  )}
                />
              </div>
              <div className='space-y-2'>
                <Label>
                  {t('productDetail.featureDefinition.propertyDialog.name')}
                </Label>
                <Input
                  value={formData.name}
                  onChange={(e) =>
                    setFormData({ ...formData, name: e.target.value })
                  }
                  placeholder={t(
                    'productDetail.featureDefinition.placeholders.propertyName'
                  )}
                />
              </div>
            </div>
            <div className='grid grid-cols-2 gap-4'>
              <div className='space-y-2'>
                <Label>
                  {t(
                    'productDetail.featureDefinition.propertyDialog.accessMode'
                  )}
                </Label>
                <Select
                  value={formData.accessMode}
                  onValueChange={(value: 'r' | 'rw') =>
                    setFormData({ ...formData, accessMode: value })
                  }
                >
                  <SelectTrigger className='w-full'>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value='r'>
                      {t('productDetail.featureDefinition.accessMode.r')}
                    </SelectItem>
                    <SelectItem value='rw'>
                      {t('productDetail.featureDefinition.accessMode.rw')}
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className='space-y-2'>
                <Label>
                  {t('productDetail.featureDefinition.propertyDialog.dataType')}
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
                  <SelectTrigger className='w-full'>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value='int'>
                      {t('productDetail.featureDefinition.dataTypes.int')}
                    </SelectItem>
                    <SelectItem value='float'>
                      {t('productDetail.featureDefinition.dataTypes.float')}
                    </SelectItem>
                    <SelectItem value='double'>
                      {t('productDetail.featureDefinition.dataTypes.double')}
                    </SelectItem>
                    <SelectItem value='bool'>
                      {t('productDetail.featureDefinition.dataTypes.bool')}
                    </SelectItem>
                    <SelectItem value='string'>
                      {t('productDetail.featureDefinition.dataTypes.string')}
                    </SelectItem>
                    <SelectItem value='enum'>
                      {t('productDetail.featureDefinition.dataTypes.enum')}
                    </SelectItem>
                    <SelectItem value='struct'>
                      {t('productDetail.featureDefinition.dataTypes.struct')}
                    </SelectItem>
                    <SelectItem value='array'>
                      {t('productDetail.featureDefinition.dataTypes.array')}
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            {/* 数值类型的规格 */}
            {['int', 'float', 'double'].includes(formData.dataType.type) && (
              <div className='space-y-2'>
                <Label>
                  {t(
                    'productDetail.featureDefinition.propertyDialog.dataSpecs'
                  )}
                </Label>
                <div className='grid grid-cols-4 gap-2'>
                  <Input
                    placeholder={t(
                      'productDetail.featureDefinition.propertyDialog.min'
                    )}
                    value={formData.dataType.specs?.min || ''}
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
                  <Input
                    placeholder={t(
                      'productDetail.featureDefinition.propertyDialog.max'
                    )}
                    value={formData.dataType.specs?.max || ''}
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
                  <Input
                    placeholder={t(
                      'productDetail.featureDefinition.propertyDialog.step'
                    )}
                    value={formData.dataType.specs?.step || ''}
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
                  <Input
                    placeholder={t(
                      'productDetail.featureDefinition.propertyDialog.unit'
                    )}
                    value={formData.dataType.specs?.unit || ''}
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
                  />
                </div>
              </div>
            )}

            {/* 枚举类型的规格 */}
            {formData.dataType.type === 'enum' && (
              <div className='space-y-2'>
                <Label>
                  {t(
                    'productDetail.featureDefinition.propertyDialog.enumValues'
                  )}
                </Label>
                <Textarea
                  placeholder='{"0": "Off", "1": "On"}'
                  value={JSON.stringify(formData.dataType.specs || {}, null, 2)}
                  onChange={(e) => {
                    try {
                      const specs = JSON.parse(e.target.value)
                      setFormData({
                        ...formData,
                        dataType: { ...formData.dataType, specs },
                      })
                    } catch {
                      // 忽略 JSON 解析错误，用户可能还在输入
                    }
                  }}
                  className='font-mono text-sm'
                />
              </div>
            )}

            {/* 结构体类型的规格 */}
            {formData.dataType.type === 'struct' && (
              <div className='space-y-2'>
                <Label>
                  {t(
                    'productDetail.featureDefinition.propertyDialog.structDef'
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
                      // 忽略 JSON 解析错误，用户可能还在输入
                    }
                  }}
                  className='min-h-[120px] font-mono text-sm'
                />
                <p className='text-xs text-muted-foreground'>
                  {t(
                    'productDetail.featureDefinition.propertyDialog.structDefDesc'
                  )}
                </p>
              </div>
            )}

            {/* 数组类型的规格 */}
            {formData.dataType.type === 'array' && (
              <div className='space-y-2'>
                <Label>
                  {t(
                    'productDetail.featureDefinition.propertyDialog.arrayElemType'
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
                      // 忽略 JSON 解析错误，用户可能还在输入
                    }
                  }}
                  className='min-h-[100px] font-mono text-sm'
                />
                <p className='text-xs text-muted-foreground'>
                  {t(
                    'productDetail.featureDefinition.propertyDialog.arrayElemDesc'
                  )}
                </p>
              </div>
            )}
          </div>
        </div>
        <DialogFooter>
          <Button variant='outline' onClick={() => onOpenChange(false)}>
            {t('productDetail.featureDefinition.buttons.cancel')}
          </Button>
          <Button onClick={handleSave}>
            {t('productDetail.featureDefinition.buttons.saveChange')}
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
  const { t } = useTranslation('deviceManagement')
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
            {t('productDetail.featureDefinition.buttons.cancel')}
          </Button>
          <Button onClick={handleSave}>
            {t('productDetail.featureDefinition.buttons.saveChange')}
          </Button>
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
  const { t } = useTranslation('deviceManagement')
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
            {t('productDetail.featureDefinition.buttons.cancel')}
          </Button>
          <Button onClick={handleSave}>
            {t('productDetail.featureDefinition.buttons.saveChange')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
