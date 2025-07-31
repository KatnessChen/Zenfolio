import { useState, useEffect, useCallback, useRef } from 'react'
import {
  fetchPortfolioSummary,
  fetchHistoricalMarketValue,
  type PortfolioSummaryResponse,
} from '@/services/portfolio.service'
import type { TimeFrame, Granularity, TotalValueTrendData } from '@/types/portfolio'

interface UsePortfolioSummaryReturn {
  data: PortfolioSummaryResponse | null
  loading: boolean
  error: string | null
  refetch: () => Promise<void>
}

interface UseHistoricalPortfolioReturn {
  data: TotalValueTrendData | null
  loading: boolean
  error: string | null
  refetch: () => Promise<void>
}

export const usePortfolioSummary = (): UsePortfolioSummaryReturn => {
  const [data, setData] = useState<PortfolioSummaryResponse | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const didFetch = useRef(false)

  const loadData = useCallback(async () => {
    setLoading(true)
    setError(null)

    try {
      const result = await fetchPortfolioSummary()
      setData(result)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch portfolio summary')
    } finally {
      setLoading(false)
      didFetch.current = true
    }
  }, [])

  useEffect(() => {
    if (didFetch.current) return
    loadData()
  }, [loadData])

  return {
    data,
    loading,
    error,
    refetch: loadData,
  }
}

export const useTotalValueTrend = (
  timeframe: TimeFrame,
  granularity?: Granularity
): UseHistoricalPortfolioReturn => {
  const [data, setData] = useState<TotalValueTrendData | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetchData = useCallback(async () => {
    setLoading(true)
    setError(null)

    try {
      const result = await fetchHistoricalMarketValue(timeframe, granularity)
      setData(result.data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch historical data')
    } finally {
      setLoading(false)
    }
  }, [timeframe, granularity])

  useEffect(() => {
    const loadData = async () => {
      setLoading(true)
      setError(null)

      try {
        const result = await fetchHistoricalMarketValue(timeframe, granularity)
        setData(result.data)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch historical data')
      } finally {
        setLoading(false)
      }
    }

    loadData()
  }, [timeframe, granularity]) // Only depend on the actual parameters

  return {
    data,
    loading,
    error,
    refetch: fetchData,
  }
}
