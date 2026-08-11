export interface DeviceSDK {
  id: string
  name: string
  docUrl: string
  packageUrl: string
}

export const deviceSDKs: DeviceSDK[] = [
  {
    id: 'c',
    name: 'C',
    docUrl: 'https://example.com/docs/c',
    packageUrl: 'https://example.com/download/c',
  },
  {
    id: 'c-extend',
    name: 'C (Extend版)',
    docUrl: 'https://example.com/docs/c-extend',
    packageUrl: 'https://example.com/download/c-extend',
  },
  {
    id: 'android',
    name: 'Android',
    docUrl: 'https://example.com/docs/android',
    packageUrl: 'https://example.com/download/android',
  },
  {
    id: 'java',
    name: 'Java',
    docUrl: 'https://example.com/docs/java',
    packageUrl: 'https://example.com/download/java',
  },
  {
    id: 'python',
    name: 'Python',
    docUrl: 'https://example.com/docs/python',
    packageUrl: 'https://example.com/download/python',
  },
  {
    id: 'nodejs',
    name: 'Node.js',
    docUrl: 'https://example.com/docs/nodejs',
    packageUrl: 'https://example.com/download/nodejs',
  },
  {
    id: 'ios',
    name: 'iOS',
    docUrl: 'https://example.com/docs/ios',
    packageUrl: 'https://example.com/download/ios',
  },
]
