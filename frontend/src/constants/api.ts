// API endpoint constants to match backend
export const API_ENDPOINTS = {
  // Authentication endpoints
  AUTH: {
    LOGIN: '/login',
    SIGNUP: '/signup',
    LOGOUT: '/logout',
    REFRESH_TOKEN: '/refresh-token',
    ME: '/me',
  },

  // Health check endpoints
  HEALTH: {
    MAIN: '/health',
    DATABASE: '/health/database',
  },

  // Transaction endpoints
  TRANSACTIONS: {
    EXTRACT: '/extract-transactions',
    HISTORY: '/transaction-history',
  },

  // Other endpoints
  HELLO_WORLD: '/hello-world',
} as const

// Helper to get full API path
export const getApiPath = (endpoint: string): string => {
  return `/api/v1${endpoint}`
}

// Typed endpoint paths
export type ApiEndpoint =
  | (typeof API_ENDPOINTS.AUTH)[keyof typeof API_ENDPOINTS.AUTH]
  | (typeof API_ENDPOINTS.HEALTH)[keyof typeof API_ENDPOINTS.HEALTH]
  | (typeof API_ENDPOINTS.TRANSACTIONS)[keyof typeof API_ENDPOINTS.TRANSACTIONS]
  | typeof API_ENDPOINTS.HELLO_WORLD
