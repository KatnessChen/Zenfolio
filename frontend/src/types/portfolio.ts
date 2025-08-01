// Portfolio-related TypeScript interfaces

export type TimeFrame = '1D' | '1W' | '1M' | '3M' | '6M' | 'YTD' | '1Y' | '5Y' | 'ALL'

export type Granularity = 'hourly' | 'daily' | 'weekly' | 'monthly'

export interface TotalValueDataPoint {
  timestamp: string
  market_value: number
  day_change: number
  day_change_percent: number
}

export interface TotalValueTrendSummary {
  change: number
  change_percent: number
  volatility: number
  max_value: number
  min_value: number
}

export interface TotalValueTrendData {
  timeframe: TimeFrame
  granularity: Granularity
  period: {
    start_date: string
    end_date: string
  }
  data_points: TotalValueDataPoint[]
  summary: TotalValueTrendSummary
}

export interface TotalValueTrendResponse {
  success: boolean
  data: TotalValueTrendData
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
