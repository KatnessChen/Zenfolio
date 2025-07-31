import React from 'react'
import { formatCurrency } from '@/utils'
import { usePortfolioSummary } from '@/hooks/usePortfolio'

interface PortfolioData {
  totalValue: number
  dailyChange: number
  dailyChangePercent: number
  totalGainLoss: number
  totalGainLossPercent: number
}

export const PortfolioSummaryCard: React.FC = () => {
  const { data: portfolioSummary, loading, error, refetch } = usePortfolioSummary()

  // Transform API data to component data format
  const portfolioData: PortfolioData = portfolioSummary
    ? {
        totalValue: portfolioSummary.market_value,
        dailyChange: portfolioSummary.daily_change,
        dailyChangePercent: portfolioSummary.daily_change_percentage,
        totalGainLoss: portfolioSummary.total_return,
        totalGainLossPercent: portfolioSummary.total_return_percentage,
      }
    : {
        totalValue: 0,
        dailyChange: 0,
        dailyChangePercent: 0,
        totalGainLoss: 0,
        totalGainLossPercent: 0,
      }

  return (
    <div className="w-full space-y-6">
      {/* Loading State */}
      {loading && (
        <div className="flex items-center justify-center h-32">
          <div className="text-muted-foreground">Loading portfolio summary...</div>
        </div>
      )}

      {/* Error State */}
      {error && (
        <div className="flex items-center justify-center h-32">
          <div className="text-destructive">
            Failed to load portfolio summary: {error}
            <button
              onClick={refetch}
              className="ml-2 text-sm underline text-primary hover:text-primary/80"
            >
              Retry
            </button>
          </div>
        </div>
      )}

      {/* Success State */}
      {!loading && !error && portfolioSummary && (
        <>
          {/* Check if user has transactions */}
          {!portfolioSummary.has_transactions ? (
            /* No Transactions State */
            <div className="flex items-center justify-center h-32">
              <div className="text-center space-y-2">
                <div className="text-lg font-medium text-muted-foreground">No transactions yet</div>
                <div className="text-sm text-muted-foreground">
                  Start by adding your first transaction to see your portfolio summary
                </div>
              </div>
            </div>
          ) : (
            /* Portfolio Data Display */
            <>
              {/* Portfolio Value and Gain/Loss (side by side) */}
              <div className="flex items-center gap-6">
                {/* Total Portfolio Value */}
                <div className="flex items-center gap-2">
                  <div className="text-3xl lg:text-4xl font-bold text-foreground">
                    {formatCurrency(portfolioData.totalValue)}
                  </div>
                  <span className="text-lg text-muted-foreground">USD</span>
                </div>
                {/* Total Gain/Loss - All Time */}
                <div
                  className={`text-lg font-medium ${portfolioData.totalGainLoss >= 0 ? 'text-profit' : 'text-loss'}`}
                >
                  {portfolioData.totalGainLoss >= 0 ? '+' : ''}
                  {formatCurrency(portfolioData.totalGainLoss)} (
                  {portfolioData.totalGainLoss >= 0 ? '+' : ''}
                  {portfolioData.totalGainLossPercent.toFixed(2)}%) all time
                </div>
              </div>
              {/* Annualized Return Rate */}
              {portfolioSummary.annualized_return_rate !== undefined && (
                <div className="flex items-center gap-2 mt-2">
                  <span className="text-muted-foreground text-base">Annualized Return:</span>
                  <span
                    className={`text-lg font-semibold ${portfolioSummary.annualized_return_rate >= 0 ? 'text-profit' : 'text-loss'}`}
                  >
                    {portfolioSummary.annualized_return_rate >= 0 ? '+' : ''}
                    {portfolioSummary.annualized_return_rate.toFixed(2)}%
                  </span>
                </div>
              )}
            </>
          )}
        </>
      )}
    </div>
  )
}
