export interface DataTypeSpecs {
  min?: string | number
  max?: string | number
  unit?: string
  step?: string | number
  [key: string]: string | number | undefined
}

export interface DataType {
  type:
    'int' | 'float' | 'double' | 'bool' | 'string' | 'enum' | 'struct' | 'array'
  specs?: DataTypeSpecs
}

export interface Property {
  identifier: string
  name: string
  accessMode: 'r' | 'rw'
  required: boolean
  desc?: string
  description?: string
  dataType: DataType
}

export interface ServiceParam {
  identifier: string
  name: string
  dataType: DataType
}

export interface Service {
  identifier: string
  name: string
  required: boolean
  callType: 'async' | 'sync'
  inputData: ServiceParam[]
  outputData: ServiceParam[]
  desc?: string
  description?: string
}

export interface EventParam {
  identifier: string
  name: string
  dataType: DataType
}

export interface Event {
  identifier: string
  name: string
  type: 'info' | 'alert' | 'error'
  required: boolean
  outputData: EventParam[]
  desc?: string
  description?: string
}

export interface TSLModel {
  schema: string
  version: string
  profile: {
    productKey: string
  }
  properties: Property[]
  events: Event[]
  services: Service[]
}

export interface FeatureDefinitionTabProps {
  productKey: string
}
