import { useEffect, useState, useCallback, useRef } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, CardContent } from '@/components/ui/card'
import { Breadcrumb } from '@/components/ui/breadcrumb'
import { ROUTES } from '@/constants'
import { formatCurrency, formatPercent } from '@/utils'
import { fetchSingleHoldingBasicInfo } from '@/services/portfolio.service'
import { useToast } from '@/hooks/useToast'

interface HoldingDetailData {
  symbol: string
  companyName: string
  currentPrice: number
  priceChangePercent: number
  totalQuantity: number
  totalCost: number
  marketValue: number
  unitCost: number
  simpleReturn: number
  annualizedReturn: number
  realizedGainLoss: number
  unrealizedGainLoss: number
}

export default function SingleHoldingDetailPage() {
  const { symbol } = useParams<{ symbol: string }>()
  const navigate = useNavigate()
  const { showToast } = useToast()
  const [holdingData, setHoldingData] = useState<HoldingDetailData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const didFetch = useRef(false)

  const fetchData = useCallback(async () => {
    if (!symbol) {
      setError('Symbol not provided')
      setLoading(false)
      return
    }

    try {
      setLoading(true)
      setError(null)

      const response = await fetchSingleHoldingBasicInfo(symbol.toUpperCase(), 'basic')

      // Transform API data to match our component interface
      setHoldingData({
        symbol: response.symbol,
        companyName: response.symbol,
        currentPrice: response.current_price,
        priceChangePercent: response.simple_return_rate,
        totalQuantity: response.total_quantity,
        totalCost: response.total_cost,
        marketValue: response.market_value,
        unitCost: response.unit_cost,
        simpleReturn: response.simple_return_rate,
        annualizedReturn: response.annualized_return_rate,
        realizedGainLoss: response.realized_gain_loss,
        unrealizedGainLoss: response.unrealized_gain_loss,
      })
    } catch (err) {
      console.error('Failed to fetch holding data:', err)

      let errorMessage = 'Failed to fetch holding data'

      if (err instanceof Error) {
        if (
          err.message.includes('no transactions found') ||
          err.message.includes('no current holdings')
        ) {
          errorMessage = `No holdings found for ${symbol.toUpperCase()}`
        } else if (err.message.includes('Unable to fetch current price data')) {
          errorMessage = 'Unable to fetch current price data. Please try again later.'
        } else if (err.message.includes('User not authenticated')) {
          errorMessage = 'Please log in to view your holdings'
          // Redirect to login page after a short delay
          setTimeout(() => {
            navigate(ROUTES.LOGIN)
          }, 2000)
        } else {
          errorMessage = err.message
        }
      }

      setError(errorMessage)
      showToast({
        type: 'error',
        title: 'Failed to Load Holding',
        message: errorMessage,
        duration: 5000,
      })
    } finally {
      setLoading(false)
    }
  }, [symbol, navigate, showToast])

  useEffect(() => {
    if (didFetch.current) return
    didFetch.current = true
    fetchData()
  }, [fetchData])

  const breadcrumbItems = [
    { label: 'Dashboard', href: ROUTES.DASHBOARD },
    { label: 'Position', href: ROUTES.DASHBOARD },
    { label: symbol?.toUpperCase() || '' },
  ]

  const formatValue = (value: number) => {
    const sign = value >= 0 ? '+' : ''
    return `${sign}${formatCurrency(Math.abs(value))}`
  }

  const getValueColor = (value: number) => {
    return value >= 0 ? 'text-profit' : 'text-loss'
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-background">
        <div className="container mx-auto py-6 px-4">
          <div className="flex justify-center items-center h-64">
            <span className="text-muted-foreground">Loading holding details...</span>
          </div>
        </div>
      </div>
    )
  }

  if (error || !holdingData) {
    return (
      <div className="min-h-screen bg-background">
        <div className="container mx-auto py-6 px-4">
          {/* Breadcrumb Navigation */}
          <Breadcrumb items={breadcrumbItems} />

          <div className="flex flex-col items-center justify-center h-64 space-y-4">
            <span className="text-destructive text-lg font-medium">
              {error || 'Holding not found'}
            </span>
            <div className="flex gap-3">
              <button
                onClick={() => fetchData()}
                className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition-colors"
              >
                Try Again
              </button>
              <button
                onClick={() => navigate(ROUTES.DASHBOARD)}
                className="px-4 py-2 bg-secondary text-secondary-foreground rounded-md hover:bg-secondary/90 transition-colors"
              >
                Back to Dashboard
              </button>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto py-6 px-4 space-y-6">
        {/* Breadcrumb Navigation */}
        <Breadcrumb items={breadcrumbItems} />

        {/* Stock Title Section */}
        <div className="flex flex-col md:flex-row md:items-center md:gap-6 space-y-2 md:space-y-0">
          <h1 className="text-3xl font-bold text-foreground flex-shrink-0">
            {holdingData.companyName}
          </h1>
          <div
            className={`flex items-center gap-3 ${getValueColor(holdingData.priceChangePercent)}`}
          >
            <span className="text-2xl font-bold">{formatCurrency(holdingData.currentPrice)}</span>
            <span className={'text-lg font-medium'}>
              ({formatPercent(holdingData.priceChangePercent)})
            </span>
          </div>
        </div>

        {/* Chart Placeholder - will be implemented later */}
        <Card>
          <CardContent>
            <div className="h-64 bg-muted/20 rounded-lg flex items-center justify-center">
              <span className="text-muted-foreground">
                Interactive Stock Chart (To be implemented)
              </span>
            </div>
          </CardContent>
        </Card>

        {/* Holdings Information Section */}
        <Card>
          <CardContent>
            <h2 className="text-xl font-semibold text-foreground mb-6">Holding Summary</h2>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-16">
              {/* Left Column */}
              <div className="space-y-4">
                <div className="flex justify-between items-center">
                  <span className="text-muted-foreground">Total Quantity</span>
                  <span className="font-medium text-foreground">
                    {holdingData.totalQuantity.toLocaleString()} shares
                  </span>
                </div>

                <div className="flex justify-between items-center">
                  <span className="text-muted-foreground">Total Cost</span>
                  <span className="font-medium text-foreground">
                    {formatCurrency(holdingData.totalCost)}
                  </span>
                </div>

                <div className="flex justify-between items-center">
                  <span className="text-muted-foreground">Simple Return</span>
                  <span className={`font-medium ${getValueColor(holdingData.simpleReturn)}`}>
                    {formatPercent(holdingData.simpleReturn)}
                  </span>
                </div>

                <div className="flex justify-between items-center">
                  <span className="text-muted-foreground">Realized Gain/Loss</span>
                  <span className={`font-medium ${getValueColor(holdingData.realizedGainLoss)}`}>
                    {formatValue(holdingData.realizedGainLoss)}
                  </span>
                </div>
              </div>

              {/* Right Column */}
              <div className="space-y-4">
                <div className="flex justify-between items-center">
                  <span className="text-muted-foreground">Market Value</span>
                  <span className="font-medium text-foreground">
                    {formatCurrency(holdingData.marketValue)}
                  </span>
                </div>

                <div className="flex justify-between items-center">
                  <span className="text-muted-foreground">Unit Cost</span>
                  <span className="font-medium text-foreground">
                    {formatCurrency(holdingData.unitCost)}
                  </span>
                </div>

                <div className="flex justify-between items-center">
                  <span className="text-muted-foreground">Annualized Return</span>
                  <span className={`font-medium ${getValueColor(holdingData.annualizedReturn)}`}>
                    {formatPercent(holdingData.annualizedReturn)}
                  </span>
                </div>

                <div className="flex justify-between items-center">
                  <span className="text-muted-foreground">Unrealized Gain/Loss</span>
                  <span className={`font-medium ${getValueColor(holdingData.unrealizedGainLoss)}`}>
                    {formatValue(holdingData.unrealizedGainLoss)}
                  </span>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
