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

// User authentication types
export interface User {
  id: string
  email: string
  name?: string
}
