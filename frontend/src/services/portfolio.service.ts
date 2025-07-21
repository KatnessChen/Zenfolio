// Portfolio API service for frontend integration
// This service will integrate with the backend portfolio endpoint

import { API_ENDPOINTS } from '@/constants/api'
import { apiClient } from '@/lib/api-client'

export interface SingleHoldingBasicInfo {
  symbol: string
  total_quantity: number
  total_cost: number
  unit_cost: number
  current_price: number
  market_value: number
  simple_return_rate: number
  annualized_return_rate: number
  realized_gain_loss: number
  unrealized_gain_loss: number
  timestamp: string
}

export interface SingleHoldingBasicInfoResponse {
  data: SingleHoldingBasicInfo
  timestamp: string
}

export interface ApiResponse<T> {
  success: boolean
  message: string
  data?: T
}

/**
 * Fetch basic single holding information for a specific symbol
 * @param symbol - Stock symbol (e.g., 'AAPL', 'GOOGL')
 * @param analysisType - Type of analysis (default: 'basic')
 * @returns Promise with stock basic info
 */
export async function fetchSingleHoldingBasicInfo(
  symbol: string,
  analysisType: 'basic' = 'basic'
): Promise<SingleHoldingBasicInfoResponse> {
  const url = `${API_ENDPOINTS.PORTFOLIO.HOLDINGS}/${symbol.toUpperCase()}?analysis_type=${analysisType}`

  try {
    const response = await apiClient.get<ApiResponse<SingleHoldingBasicInfoResponse>>(url)
    const result = response.data

    if (!result.success) {
      throw new Error(result.message || 'Failed to fetch stock basic info')
    }

    return result.data!
  } catch (error) {
    console.error('Error fetching stock basic info:', error)
    throw error
  }
}
