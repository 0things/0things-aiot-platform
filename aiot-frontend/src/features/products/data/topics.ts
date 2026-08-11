// Topic 操作权限类型
export type TopicPermission = 'publish' | 'subscribe'

// Topic 分类
export type TopicCategory =
  | 'otaUpgrade'
  | 'deviceTag'
  | 'timeSync'
  | 'deviceShadow'
  | 'configUpdate'
  | 'broadcast'
  | 'propertyPost'
  | 'propertySet'
  | 'eventPost'
  | 'serviceCall'

// 基础通信 Topic 配置
export interface BasicTopic {
  category: TopicCategory
  topic: string
  permission: TopicPermission
  description: string
}

// 物模型通信 Topic 配置
export interface ThingModelTopic {
  category: TopicCategory
  topic: string
  permission: TopicPermission
  description: string
}

// 自定义 Topic 配置
export interface CustomTopic {
  id: string
  topic: string
  permission: TopicPermission
  allowBroadcast: boolean
  compressionEnabled: boolean
  description: string
}

// 基础通信 Topic 列表
export const basicTopics: BasicTopic[] = [
  // OTA 升级
  {
    category: 'otaUpgrade',
    topic: '/ota/device/inform/ht0yxt4KxDo/${deviceName}',
    permission: 'publish',
    description: 'otaInform',
  },
  {
    category: 'otaUpgrade',
    topic: '/ota/device/upgrade/ht0yxt4KxDo/${deviceName}',
    permission: 'subscribe',
    description: 'otaUpgrade',
  },
  {
    category: 'otaUpgrade',
    topic: '/ota/device/progress/ht0yxt4KxDo/${deviceName}',
    permission: 'publish',
    description: 'otaProgress',
  },
  // 设备标签
  {
    category: 'deviceTag',
    topic: '/sys/ht0yxt4KxDo/${deviceName}/thing/ota/firmware/get',
    permission: 'publish',
    description: 'firmwareGet',
  },
  {
    category: 'deviceTag',
    topic: '/sys/ht0yxt4KxDo/${deviceName}/thing/deviceinfo/update',
    permission: 'publish',
    description: 'deviceInfoUpdate',
  },
  {
    category: 'deviceTag',
    topic: '/sys/ht0yxt4KxDo/${deviceName}/thing/deviceinfo/update_reply',
    permission: 'subscribe',
    description: 'deviceInfoUpdateReply',
  },
  {
    category: 'deviceTag',
    topic: '/sys/ht0yxt4KxDo/${deviceName}/thing/deviceinfo/delete',
    permission: 'subscribe',
    description: 'deviceInfoDelete',
  },
  {
    category: 'deviceTag',
    topic: '/sys/ht0yxt4KxDo/${deviceName}/thing/deviceinfo/delete_reply',
    permission: 'publish',
    description: 'deviceInfoDeleteReply',
  },
  // 时转同步
  {
    category: 'timeSync',
    topic: '/ext/ntp/ht0yxt4KxDo/${deviceName}/request',
    permission: 'publish',
    description: 'ntpRequest',
  },
  {
    category: 'timeSync',
    topic: '/ext/ntp/ht0yxt4KxDo/${deviceName}/response',
    permission: 'subscribe',
    description: 'ntpResponse',
  },
  // 设备影子
  {
    category: 'deviceShadow',
    topic: '/shadow/update/ht0yxt4KxDo/${deviceName}',
    permission: 'publish',
    description: 'shadowUpdate',
  },
  {
    category: 'deviceShadow',
    topic: '/shadow/get/ht0yxt4KxDo/${deviceName}',
    permission: 'subscribe',
    description: 'shadowGet',
  },
  // 配置更新
  {
    category: 'configUpdate',
    topic: '/sys/ht0yxt4KxDo/${deviceName}/thing/config/push',
    permission: 'subscribe',
    description: 'configPush',
  },
  {
    category: 'configUpdate',
    topic: '/sys/ht0yxt4KxDo/${deviceName}/thing/config/get',
    permission: 'publish',
    description: 'configGet',
  },
  {
    category: 'configUpdate',
    topic: '/sys/ht0yxt4KxDo/${deviceName}/thing/config/get_reply',
    permission: 'subscribe',
    description: 'configGetReply',
  },
  // 广播
  {
    category: 'broadcast',
    topic: '/broadcast/ht0yxt4KxDo/${identifier}',
    permission: 'subscribe',
    description: 'broadcast',
  },
]

// 物模型通信 Topic 列表
export const thingModelTopics: ThingModelTopic[] = [
  // 属性上报
  {
    category: 'propertyPost',
    topic: '/sys/ht0yxt4KxDo/${deviceName}/thing/event/property/post',
    permission: 'publish',
    description: 'propertyPost',
  },
  {
    category: 'propertyPost',
    topic: '/sys/ht0yxt4KxDo/${deviceName}/thing/event/property/post_reply',
    permission: 'subscribe',
    description: 'propertyPostReply',
  },
  // 属性设置
  {
    category: 'propertySet',
    topic: '/sys/ht0yxt4KxDo/${deviceName}/thing/service/property/set',
    permission: 'subscribe',
    description: 'propertySet',
  },
  // 事件上报
  {
    category: 'eventPost',
    topic:
      '/sys/ht0yxt4KxDo/${deviceName}/thing/event/${tsl.event.identifier}/post',
    permission: 'publish',
    description: 'eventPost',
  },
  {
    category: 'eventPost',
    topic:
      '/sys/ht0yxt4KxDo/${deviceName}/thing/event/${tsl.event.identifier}/post_reply',
    permission: 'subscribe',
    description: 'eventPostReply',
  },
  // 服务调用
  {
    category: 'serviceCall',
    topic:
      '/sys/ht0yxt4KxDo/${deviceName}/thing/service/${tsl.service.identifier}',
    permission: 'subscribe',
    description: 'serviceCall',
  },
  {
    category: 'serviceCall',
    topic:
      '/sys/ht0yxt4KxDo/${deviceName}/thing/service/${tsl.service.identifier}_reply',
    permission: 'publish',
    description: 'serviceCallReply',
  },
]

/**
 * 替换 Topic 中的 productKey 占位符
 */
export function replaceProductKey(topic: string, productKey: string): string {
  return topic.replace(/ht0yxt4KxDo/g, productKey)
}
