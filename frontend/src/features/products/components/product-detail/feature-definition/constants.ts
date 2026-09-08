// TSL模板 - 名称使用翻译键，实际显示名称在渲染时翻译
export const TSL_TEMPLATES = {
  empty: {
    nameKey: 'productDetail.featureDefinition.templates.empty',
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
} as const

// 常用工程单位
export const COMMON_UNITS = [
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
