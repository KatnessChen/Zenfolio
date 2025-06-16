// Common API response structure
export interface ApiResponse<T> {
  success: boolean
  data?: T
  message?: string
  error?: string
}

// Transaction related types
export interface Transaction {
  id: string
  ticker: string
  tickerLabel?: string
  tradeType: 'Buy' | 'Sell' | 'Dividends'
  quantity: number
  price: number
  tradeDate: string
  exchange?: string
  currency?: string
}

// User authentication types
export interface User {
  id: string
  email: string
  name?: string
}

export interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
}

// API endpoints
export interface ApiEndpoints {
  AUTH: '/api/v1/login'
  HEALTH: '/api/v1/health'
  EXTRACT_TRANSACTIONS: '/api/v1/extract-transactions'
}
