// Common API response structure
export interface ApiResponse<T> {
  success: boolean
  data?: T
  message?: string
  error?: string
}

// Shared user type (used in multiple features)
export interface User {
  email: string
  firstName: string
  lastName?: string
}

// Transaction extraction types to match backend API
export type TradeType = 'Buy' | 'Sell' | 'Dividends'

export interface TransactionData {
  symbol: string
  type: TradeType
  quantity: number
  price: number
  amount: number
  currency: string
  broker: string
  account: string
  transaction_date: string
  user_notes: string
  exchange: string
}

export interface ExtractResponseData {
  transactions: TransactionData[]
  transaction_count: number
  image_name: string
}

export interface ExtractResponse {
  data?: ExtractResponseData
  success: boolean
  message: string
}

// Processing states for frontend
export type FileProcessingStatus = 'pending' | 'processing' | 'completed' | 'error'

export interface FileProcessingState {
  file: File
  status: FileProcessingStatus
  result?: ExtractResponseData
  error?: string
  progress?: number
}

// Legacy interface for backward compatibility during migration
export interface ProcessedTransaction {
  id: string
  ticker: string
  tradeType: 'Buy' | 'Sell' | 'Dividend'
  quantity: number
  price: number
  amount: number
  tradeDate: string
  broker: string
  currency: string
  userNotes: string
  confidence: number
}
