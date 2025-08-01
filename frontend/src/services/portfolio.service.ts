// Portfolio API service for frontend integration
// This service will integrate with the backend portfolio endpoint

import { API_ENDPOINTS } from '@/constants/api'
import { apiClient } from '@/lib/api-client'
import type { TimeFrame, Granularity, TotalValueTrendResponse } from '@/types/portfolio'

export interface SingleHoldingBasicInfoResponse {
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

export interface AllHoldingsResponse {
  holdings: SingleHoldingBasicInfoResponse[]
  timestamp: string
}

export interface PortfolioSummaryResponse {
  timestamp: string
  currency: string
  market_value: number
  total_cost: number
  total_return: number
  total_return_percentage: number
  holdings_count: number
  has_transactions: boolean
  annualized_return_rate: number
  last_updated: string
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
  analysisType: 'basic'
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

/**
 * Fetch all portfolio holdings for the current user
 * @returns Promise with all holdings data
 */
export async function fetchAllHoldings(): Promise<AllHoldingsResponse> {
  const url = API_ENDPOINTS.PORTFOLIO.HOLDINGS

  try {
    const response = await apiClient.get<AllHoldingsResponse>(url)
    return response.data
  } catch (error) {
    console.error('Error fetching all holdings:', error)
    throw error
  }
}

/**
 * Fetch portfolio summary for the current user
 * @returns Promise with portfolio summary data
 */
export async function fetchPortfolioSummary(): Promise<PortfolioSummaryResponse> {
  const url = API_ENDPOINTS.PORTFOLIO.SUMMARY

  try {
    const response = await apiClient.get<PortfolioSummaryResponse>(url)
    return response.data
  } catch (error) {
    console.error('Error fetching portfolio summary:', error)
    throw error
  }
}

/**
 * Fetch historical market value data for portfolio chart
 * @param timeframe - Time period for historical data (1D, 1W, 1M, etc.)
 * @param granularity - Optional data point frequency (hourly, daily, weekly, monthly)
 * @returns Promise with historical market value data
 */
export async function fetchHistoricalMarketValue(
  timeframe: TimeFrame,
  granularity?: Granularity
): Promise<TotalValueTrendResponse> {
  const url = API_ENDPOINTS.PORTFOLIO.HISTORICAL_CHART
  const params = new URLSearchParams({ timeframe })

  if (granularity) {
    params.append('granularity', granularity)
  }

  try {
    const response = await apiClient.get<TotalValueTrendResponse>(`${url}?${params}`)

    if (!response.data.success) {
      throw new Error('Failed to fetch historical market value data')
    }

    return response.data
  } catch (error) {
    console.error('Error fetching historical market value:', error)
    throw error
  }
}
