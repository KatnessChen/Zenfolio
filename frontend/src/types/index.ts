import type { TRADE_TYPE } from '@/constants'

// Common API response structure
export interface ApiResponse<T> {
  success: boolean
  data?: T
  message?: string
  error?: string
}

// Transaction related types
export type TradeType = (typeof TRADE_TYPE)[keyof typeof TRADE_TYPE]

export interface Transaction {
  id: string
  ticker: string
  tickerLabel?: string
  tradeType: TradeType
  quantity: number
  price: number
  amount: number // Total transaction amount (price * quantity)
  tradeDate: string
  uploadDate: string
  broker?: string
  exchange?: string
  currency?: string
  userNotes?: string
  transactionHistory?: string // Brief summary/description of the transaction
}

// Authentication types
export interface User {
  email: string
  firstName: string
  lastName?: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  success: boolean
  message: string
  data: {
    token: string
    user: User
  }
}

export interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  error: string | null
}
