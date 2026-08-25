/**
 * API Clients
 *
 * This file exports configured API client instances for all services.
 * Each service has its own generated API client from OpenAPI specifications.
 * All clients share a common Axios instance with interceptors for auth and error handling.
 */
import axios, { type AxiosInstance } from 'axios'
// Import auth-service APIs when generated
// import { Configuration as AuthConfiguration, UsersApi } from './generated/auth-service';
// Import notification-service APIs when generated
// import { Configuration as NotificationConfiguration, NotificationsApi } from './generated/notification-service';

import { useAuthStore } from '@/stores/auth-store'
import { API_TIMEOUT } from './config'

// Import other service URLs when needed
// import { AUTH_SERVICE_BASE_URL, NOTIFICATION_SERVICE_BASE_URL } from './config';

// ============================================================================
// Custom Axios Instance with Interceptors
// ============================================================================

/**
 * Create a custom Axios instance with request/response interceptors
 * This instance is shared across all API clients for consistent behavior
 */
const createAxiosInstance = (): AxiosInstance => {
  const instance = axios.create({
    timeout: API_TIMEOUT,
    headers: {
      'Content-Type': 'application/json',
    },
    // Configure params serializer for standard REST API array format
    // Arrays will be serialized as: ?state=online&state=offline
    // instead of: ?state[]=online&state[]=offline
    paramsSerializer: {
      serialize: (params) => {
        const searchParams = new URLSearchParams()
        Object.entries(params).forEach(([key, value]) => {
          if (Array.isArray(value)) {
            // Add each array element as a separate parameter
            value.forEach((item) => searchParams.append(key, String(item)))
          } else if (value !== undefined && value !== null) {
            searchParams.append(key, String(value))
          }
        })
        return searchParams.toString()
      },
    },
  })

  // Request interceptor - add auth tokens, logging, etc.
  instance.interceptors.request.use(
    (config) => {
      // Add authorization token if available
      const token =
        useAuthStore.getState().auth.accessToken ||
        localStorage.getItem('authToken')
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }

      // Log request in development
      if (import.meta.env.DEV) {
        // eslint-disable-next-line no-console
        console.log(
          `[API Request] ${config.method?.toUpperCase()} ${config.url}`,
          config.data
        )
      }

      return config
    },
    (error) => {
      // eslint-disable-next-line no-console
      console.error('[API Request Error]', error)
      return Promise.reject(error)
    }
  )

  // Response interceptor - handle common response scenarios
  instance.interceptors.response.use(
    (response) => {
      // Log response in development
      if (import.meta.env.DEV) {
        // eslint-disable-next-line no-console
        console.log(
          `[API Response] ${response.config.method?.toUpperCase()} ${response.config.url}`,
          response.data
        )
      }
      return response
    },
    (error) => {
      // Handle common error scenarios
      if (error.response) {
        const { status } = error.response

        // Handle 401 Unauthorized - token expired or invalid
        if (status === 401) {
          // eslint-disable-next-line no-console
          console.error('[API Error] Unauthorized - clearing token')
          localStorage.removeItem('authToken')
          // You can redirect to login page here
          // window.location.href = '/login';
        }

        // Handle 403 Forbidden
        if (status === 403) {
          // eslint-disable-next-line no-console
          console.error('[API Error] Forbidden - insufficient permissions')
        }

        // Handle 500 Server Error
        if (status >= 500) {
          // eslint-disable-next-line no-console
          console.error('[API Error] Server error')
        }
      }

      // Log error in development
      if (import.meta.env.DEV) {
        // eslint-disable-next-line no-console
        console.error('[API Error]', error.response?.data || error.message)
      }

      // Let the error propagate to be handled by QueryClient error handlers
      return Promise.reject(error)
    }
  )

  return instance
}

/**
 * Shared Axios instance for all API clients
 */
const axiosInstance = createAxiosInstance()

/**
 * Export axios instance for custom API calls (e.g., multipart/form-data uploads)
 */
export { axiosInstance }

// ============================================================================
// Token Management
// ============================================================================

/**
 * Update auth token in localStorage
 * The axios interceptor will automatically pick it up for subsequent requests
 *
 * @param token - JWT token or access token
 */
export const setAuthToken = (token: string) => {
  localStorage.setItem('authToken', token)
}

/**
 * Clear auth token from localStorage
 */
export const clearAuthToken = () => {
  localStorage.removeItem('authToken')
}

/**
 * Get current auth token from localStorage
 */
export const getAuthToken = (): string | null => {
  return localStorage.getItem('authToken')
}
