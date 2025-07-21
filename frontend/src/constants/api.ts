// API endpoint constants to match backend
export const API_ENDPOINTS = {
  // Authentication endpoints
  AUTH: {
    LOGIN: '/login',
    SIGNUP: '/signup',
    LOGOUT: '/logout',
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

  // Portfolio endpoints
  PORTFOLIO: {
    HOLDINGS: '/portfolio/holdings',
  },

  // Other endpoints
  HELLO_WORLD: '/hello-world',
} as const

// Typed endpoint paths
export type ApiEndpoint =
  | (typeof API_ENDPOINTS.AUTH)[keyof typeof API_ENDPOINTS.AUTH]
  | (typeof API_ENDPOINTS.HEALTH)[keyof typeof API_ENDPOINTS.HEALTH]
  | (typeof API_ENDPOINTS.TRANSACTIONS)[keyof typeof API_ENDPOINTS.TRANSACTIONS]
  | (typeof API_ENDPOINTS.PORTFOLIO)[keyof typeof API_ENDPOINTS.PORTFOLIO]
  | typeof API_ENDPOINTS.HELLO_WORLD
