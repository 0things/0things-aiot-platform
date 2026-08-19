import type { AxiosRequestConfig } from 'axios'
import { axiosInstance } from './clients'
import { DEVICE_SERVICE_BASE_URL } from './config'

export const orvalAxios = <T>(config: AxiosRequestConfig): Promise<T> =>
  axiosInstance({
    ...config,
    url: `${DEVICE_SERVICE_BASE_URL}/v1${config.url}`,
  }).then(({ data }) => data)

export type ErrorType<Error> = Error
export type BodyType<BodyData> = BodyData
