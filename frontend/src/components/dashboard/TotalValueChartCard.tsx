import React, { useState } from 'react'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts'
// Custom tooltip for chart
import type { TooltipProps } from 'recharts'
type CustomTooltipProps = TooltipProps<number, string> & {
  payload?: Array<{ value: number }>
  label?: string | number
  active?: boolean
}
const CustomTooltip: React.FC<CustomTooltipProps> = (props) => {
  const { active, payload, label } = props
  if (active && payload && payload.length) {
    return (
      <div className="rounded border bg-card p-3 shadow">
        <div className="text-xs text-muted-foreground mb-1">{label}</div>
        <div className="font-semibold text-sm">
          Total Value:{' '}
          {new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD',
            minimumFractionDigits: 0,
            maximumFractionDigits: 0,
          }).format(payload[0].value)}
        </div>
      </div>
    )
  }
  return null
}
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Loader2, TrendingUp, TrendingDown } from 'lucide-react'
import { useTotalValueTrend } from '@/hooks/usePortfolio'
import type { TimeFrame } from '@/types/portfolio'
import { Selector } from '@/components/ui/selector'

const timeframeOptions: { value: TimeFrame; label: string }[] = [
  { value: '1W', label: '1W' },
  { value: '1M', label: '1M' },
  { value: '3M', label: '3M' },
  { value: '6M', label: '6M' },
  { value: 'YTD', label: 'YTD' },
  { value: '1Y', label: '1Y' },
  { value: '5Y', label: '5Y' },
  { value: 'ALL', label: 'ALL' },
]

export const TotalValueChartCard: React.FC = () => {
  const [selectedTimeframe, setSelectedTimeframe] = useState<TimeFrame>('1M')
  const { data, loading, error } = useTotalValueTrend(selectedTimeframe)

  // Transform data for the chart
  const chartData =
    data?.data_points.map((point) => ({
      timestamp: new Date(point.timestamp).toLocaleDateString(),
      value: point.market_value,
      change: point.day_change_percent,
      fullTimestamp: point.timestamp,
    })) || []

  // Format currency values
  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(value)
  }

  // Format percentage values
  const formatPercentage = (value: number) => {
    if (value === 0) return '∞'
    return `${value >= 0 ? '+' : ''}${value.toFixed(2)}%`
  }

  // Determine trend color based on overall performance
  const getTrendColor = () => {
    if (!data?.summary) return 'text-gray-500'
    return data.summary.change_percent >= 0 ? 'text-profit' : 'text-loss'
  }

  const getTrendIcon = () => {
    if (!data?.summary) return null
    return data.summary.change_percent >= 0 ? (
      <TrendingUp className="h-4 w-4" />
    ) : (
      <TrendingDown className="h-4 w-4" />
    )
  }

  // Get brand colors for chart
  const getChartStrokeColor = () => {
    if (!data?.summary) return 'hsl(var(--muted))'
    return data.summary.change_percent >= 0 ? 'hsl(var(--profit))' : 'hsl(var(--loss))'
  }

  return (
    <Card className="w-full">
      <CardHeader>
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <CardTitle className="text-xl font-semibold">Total Value</CardTitle>
          </div>
          <div className="flex flex-wrap gap-1">
            <Selector
              options={timeframeOptions}
              value={selectedTimeframe}
              onChange={(val) => setSelectedTimeframe(val as TimeFrame)}
              className="w-full sm:w-auto"
              type="horizontal"
            />
          </div>
        </div>
      </CardHeader>
      <CardContent>
        {loading && (
          <div className="flex items-center justify-center h-64">
            <Loader2 className="h-8 w-8 animate-spin text-gray-500" />
            <span className="ml-2 text-gray-500">Loading chart data...</span>
          </div>
        )}

        {error && (
          <div className="flex items-center justify-center h-64">
            <div className="text-center">
              <p className="text-warning font-medium">Failed to load chart data</p>
              <p className="text-sm text-gray-500 mt-1">{error}</p>
            </div>
          </div>
        )}

        {data && !loading && !error && (
          <>
            {/* Summary Statistics */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
              <div className="text-center">
                <p className="text-sm text-muted-foreground">Change</p>
                <p className={`text-lg font-semibold ${getTrendColor()}`}>
                  {formatCurrency(data.summary.change)}
                </p>
              </div>
              <div className="text-center">
                <p className="text-sm text-muted-foreground">Change %</p>
                <div
                  className={`flex items-center justify-center text-lg font-semibold ${getTrendColor()}`}
                >
                  {getTrendIcon()}
                  <span className="ml-1">
                    {data.summary.change_percent === 0 ? (
                      <>&#8734;</>
                    ) : (
                      formatPercentage(data.summary.change_percent)
                    )}
                  </span>
                </div>
              </div>
              <div className="text-center">
                <p className="text-sm text-muted-foreground">Max Value</p>
                <p className="text-lg font-semibold">{formatCurrency(data.summary.max_value)}</p>
              </div>
              <div className="text-center">
                <p className="text-sm text-muted-foreground">Volatility</p>
                <p className="text-lg font-semibold">{data.summary.volatility.toFixed(2)}%</p>
              </div>
            </div>

            {/* Chart */}
            <div className="h-64 w-full">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={chartData}>
                  <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--muted))" />
                  <XAxis
                    dataKey="timestamp"
                    tick={{ fontSize: 12 }}
                    stroke="hsl(var(--muted-foreground))"
                  />
                  <YAxis
                    tickFormatter={formatCurrency}
                    tick={{ fontSize: 12 }}
                    stroke="hsl(var(--muted-foreground))"
                  />
                  <Tooltip content={<CustomTooltip />} />
                  <Line
                    type="monotone"
                    dataKey="value"
                    stroke={getChartStrokeColor()}
                    strokeWidth={2}
                    dot={false}
                    activeDot={{ r: 4, stroke: 'hsl(var(--card))', strokeWidth: 2 }}
                  />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </>
        )}

        {data && data.data_points.length === 0 && !loading && !error && (
          <div className="flex items-center justify-center h-64">
            <div className="text-center">
              <p className="text-gray-600 font-medium">No data available</p>
              <p className="text-sm text-gray-500 mt-1">
                No transactions found for the selected timeframe
              </p>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
