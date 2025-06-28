import React, { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { formatCurrency } from '@/utils'

type PerformerView = 'Top 5' | 'Worst 5'

interface StockPerformance {
  symbol: string
  name: string
  currentPrice: number
  dailyChange: number
  dailyChangePercent: number
}

export const TopPerformersCard: React.FC = () => {
  const [selectedView, setSelectedView] = useState<PerformerView>('Top 5')

  // Mock data - in real app this would come from API
  const performanceData: Record<PerformerView, StockPerformance[]> = {
    'Top 5': [
      {
        symbol: 'NVDA',
        name: 'NVIDIA Corp',
        currentPrice: 885.5,
        dailyChange: 15.25,
        dailyChangePercent: 1.75,
      },
      {
        symbol: 'AAPL',
        name: 'Apple Inc',
        currentPrice: 150.25,
        dailyChange: 2.5,
        dailyChangePercent: 1.69,
      },
      {
        symbol: 'MSFT',
        name: 'Microsoft Corp',
        currentPrice: 342.8,
        dailyChange: 5.1,
        dailyChangePercent: 1.51,
      },
      {
        symbol: 'GOOGL',
        name: 'Alphabet Inc',
        currentPrice: 2750.8,
        dailyChange: 35.2,
        dailyChangePercent: 1.3,
      },
      {
        symbol: 'AMZN',
        name: 'Amazon.com Inc',
        currentPrice: 3125.45,
        dailyChange: 28.75,
        dailyChangePercent: 0.93,
      },
    ],
    'Worst 5': [
      {
        symbol: 'TSLA',
        name: 'Tesla Inc',
        currentPrice: 250.75,
        dailyChange: -8.5,
        dailyChangePercent: -3.28,
      },
      {
        symbol: 'META',
        name: 'Meta Platforms',
        currentPrice: 485.2,
        dailyChange: -12.3,
        dailyChangePercent: -2.47,
      },
      {
        symbol: 'NFLX',
        name: 'Netflix Inc',
        currentPrice: 425.6,
        dailyChange: -9.8,
        dailyChangePercent: -2.25,
      },
      {
        symbol: 'AMD',
        name: 'AMD Inc',
        currentPrice: 142.35,
        dailyChange: -2.85,
        dailyChangePercent: -1.96,
      },
      {
        symbol: 'UBER',
        name: 'Uber Technologies',
        currentPrice: 68.9,
        dailyChange: -1.2,
        dailyChangePercent: -1.71,
      },
    ],
  }

  const currentData = performanceData[selectedView]

  const formatPercentage = (percent: number) => {
    return `${percent >= 0 ? '+' : ''}${percent.toFixed(2)}%`
  }

  return (
    <Card className="bg-card border-border/50">
      <CardHeader>
        <CardTitle className="text-lg text-foreground">Top & Worst Performers</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Toggle Buttons */}
        <div className="flex justify-center gap-1">
          <button
            onClick={() => setSelectedView('Top 5')}
            className={`px-3 py-1.5 text-sm rounded-md transition-colors ${
              selectedView === 'Top 5'
                ? 'bg-muted text-foreground'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted/50'
            }`}
          >
            Top 5
          </button>
          <button
            onClick={() => setSelectedView('Worst 5')}
            className={`px-3 py-1.5 text-sm rounded-md transition-colors ${
              selectedView === 'Worst 5'
                ? 'bg-muted text-foreground'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted/50'
            }`}
          >
            Worst 5
          </button>
        </div>

        {/* Performance List */}
        <div className="space-y-3">
          {currentData.map((stock, index) => (
            <div
              key={stock.symbol}
              className="flex items-center justify-between p-3 bg-muted/20 rounded-md"
            >
              <div className="flex items-center space-x-3">
                {/* Rank */}
                <div className="flex items-center justify-center w-6 h-6 rounded-full bg-primary/20 text-primary text-xs font-medium">
                  {index + 1}
                </div>

                {/* Stock Info */}
                <div>
                  <div className="font-medium text-foreground">{stock.symbol}</div>
                  <div className="text-xs text-muted-foreground truncate max-w-[120px]">
                    {stock.name}
                  </div>
                </div>
              </div>

              {/* Price and Change */}
              <div className="text-right">
                <div className="text-sm font-medium text-foreground">
                  {formatCurrency(stock.currentPrice)}
                </div>
                <div
                  className={`text-sm font-medium ${
                    stock.dailyChange >= 0 ? 'text-profit' : 'text-destructive'
                  }`}
                >
                  {formatCurrency(Math.abs(stock.dailyChange))} (
                  {formatPercentage(stock.dailyChangePercent)})
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* Summary Info */}
        <div className="text-xs text-muted-foreground text-center pt-2 border-t border-border">
          {selectedView === 'Top 5'
            ? 'Showing best performing stocks today'
            : 'Showing worst performing stocks today'}
        </div>
      </CardContent>
    </Card>
  )
}
