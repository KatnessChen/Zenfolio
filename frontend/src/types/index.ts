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

export type TransactionData = {
  id: string
  symbol: string
  trade_type: TradeType
  quantity: number
  price: number
  amount: number
  currency: string
  broker: string
  transaction_date: string
  user_notes: string
  exchange: string
  transaction_id?: string // from backend
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

export interface SerializableFile {
  name: string
  size: number
  type: string
  lastModified: number
  dataUrl: string // Base64 encoded file data
}

export interface FileProcessingState {
  file: File
  status: FileProcessingStatus
  result?: ExtractResponseData
  error?: string
}
