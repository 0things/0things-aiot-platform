/**
 * API Client Configuration
 *
 * This file contains the base configuration for the generated API clients.
 * Update these values according to your backend API settings.
 */

/**
 * Base URLs for different services
 * Can be overridden via environment variables
 */
export const DEVICE_SERVICE_BASE_URL =
  import.meta.env.VITE_DEVICE_SERVICE_URL || 'http://localhost:8000'

export const AUTH_SERVICE_BASE_URL =
  import.meta.env.VITE_AUTH_SERVICE_URL || 'http://localhost:8003'

export const NOTIFICATION_SERVICE_BASE_URL =
  import.meta.env.VITE_NOTIFICATION_SERVICE_URL || 'http://localhost:8004'

/**
 * Default timeout for API requests (in milliseconds)
 */
export const API_TIMEOUT = 30000

/**
 * Common headers for all API requests
 */
export const DEFAULT_HEADERS = {
  'Content-Type': 'application/json',
}
