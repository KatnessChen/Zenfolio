import React, { useState } from 'react'
import { XAxis, YAxis, ResponsiveContainer, Area, AreaChart, Tooltip, Dot } from 'recharts'
import { Selector } from '@/components/ui/selector'
import { formatCurrency } from '@/utils'

interface PortfolioData {
  totalValue: number
  dailyChange: number
  dailyChangePercent: number
  totalGainLoss: number
  totalGainLossPercent: number
}

// Type for Recharts tooltip props
interface CustomTooltipProps {
  active?: boolean
  payload?: Array<{ payload: { date: string; value: number } }>
}

// Type for Recharts dot props
interface CustomDotProps {
  cx?: number
  cy?: number
  payload?: { value: number; date: string }
}

export const PortfolioSummaryCard: React.FC = () => {
  const [selectedPeriod, setSelectedPeriod] = useState('ALL')

  // Mock data - in real app this would come from API
  const portfolioData: PortfolioData = {
    totalValue: 80085.72,
    dailyChange: 1205.5,
    dailyChangePercent: 1.4,
    totalGainLoss: 4331.67,
    totalGainLossPercent: 5.72,
  }

  // Time period options for the selector
  const timePeriodOptions = [
    { value: '1D', label: '1D' },
    { value: '1W', label: '1W' },
    { value: '1M', label: '1M' },
    { value: '3M', label: '3M' },
    { value: '6M', label: '6M' },
    { value: 'YTD', label: 'YTD' },
    { value: '1Y', label: '1Y' },
    { value: 'ALL', label: 'ALL' },
  ]

  // Mock chart data for portfolio performance with realistic fluctuation
  const chartData = [
    { time: '1', date: '2024-12-01', value: 75000 },
    { time: '2', date: '2024-12-02', value: 74850 },
    { time: '3', date: '2024-12-03', value: 76420 },
    { time: '4', date: '2024-12-04', value: 75980 },
    { time: '5', date: '2024-12-05', value: 76850 },
    { time: '6', date: '2024-12-06', value: 74320 },
    { time: '7', date: '2024-12-07', value: 77650 },
    { time: '8', date: '2024-12-08', value: 78120 },
    { time: '9', date: '2024-12-09', value: 76890 },
    { time: '10', date: '2024-12-10', value: 78450 },
    { time: '11', date: '2024-12-11', value: 77980 },
    { time: '12', date: '2024-12-12', value: 79320 },
    { time: '13', date: '2024-12-13', value: 78750 },
    { time: '14', date: '2024-12-14', value: 76540 },
    { time: '15', date: '2024-12-15', value: 78890 },
    { time: '16', date: '2024-12-16', value: 79680 },
    { time: '17', date: '2024-12-17', value: 78230 },
    { time: '18', date: '2024-12-18', value: 79150 },
    { time: '19', date: '2024-12-19', value: 77820 },
    { time: '20', date: '2024-12-20', value: 77200.72 },
  ]

  // Find highest and lowest points
  const highestPoint = chartData.reduce((prev, current) =>
    prev.value > current.value ? prev : current
  )
  const lowestPoint = chartData.reduce((prev, current) =>
    prev.value < current.value ? prev : current
  )

  // Determine if portfolio is gaining or losing (compare first and last values)
  const isGaining = chartData[chartData.length - 1].value > chartData[0].value
  const chartColor = isGaining ? 'rgb(34, 197, 94)' : 'rgb(255, 59, 48)' // profit green : loss red

  // Custom tooltip component
  const CustomTooltip = ({ active, payload }: CustomTooltipProps) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload
      return (
        <div className="bg-card border border-border rounded-lg p-3 shadow-lg">
          <p className="text-foreground font-medium">
            {new Date(data.date).toLocaleDateString('en-US', {
              month: 'short',
              day: 'numeric',
              year: 'numeric',
            })}
          </p>
          <p className="text-foreground">Value: {formatCurrency(data.value)}</p>
        </div>
      )
    }
    return null
  }

  // Custom dot component for highest/lowest points
  const CustomDot = (props: CustomDotProps) => {
    const { cx, cy, payload } = props
    
    // Early return if required props are undefined
    if (cx === undefined || cy === undefined || !payload) {
      return null
    }
    
    const isHighest = payload.value === highestPoint.value
    const isLowest = payload.value === lowestPoint.value

    if (isHighest || isLowest) {
      return (
        <g>
          <Dot
            cx={cx}
            cy={cy}
            r={4}
            fill={isHighest ? 'rgb(34, 197, 94)' : 'rgb(255, 59, 48)'}
            stroke="white"
            strokeWidth={2}
          />
          <text
            x={cx}
            y={cy - 10}
            textAnchor="middle"
            fill="currentColor"
            fontSize="12"
            className="text-foreground font-medium"
          >
            {formatCurrency(payload.value)}
          </text>
          <text
            x={cx}
            y={cy - 30}
            textAnchor="middle"
            fill="currentColor"
            fontSize="11"
            className="text-muted-foreground"
          >
            {new Date(payload.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
          </text>
        </g>
      )
    }
    return null
  }

  return (
    <div className="w-full space-y-6">
      {/* Portfolio Value and Gain/Loss */}
      <div className="space-y-2">
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

      {/* Recharts Area Chart */}
      <div className="h-72 bg-muted/20 rounded-lg relative overflow-hidden">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={chartData} margin={{ top: 35 }}>
            <defs>
              <linearGradient id="portfolioGradient" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor={chartColor} stopOpacity={0.3} />
                <stop offset="100%" stopColor={chartColor} stopOpacity={0.05} />
              </linearGradient>
            </defs>
            <XAxis dataKey="time" hide />
            <YAxis hide domain={['dataMin - 1000', 'dataMax + 1000']} />
            <Tooltip content={<CustomTooltip />} />
            <Area
              type="monotone"
              dataKey="value"
              stroke={chartColor}
              strokeWidth={2}
              fill="url(#portfolioGradient)"
              fillOpacity={1}
              dot={<CustomDot />}
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>

      {/* Time Period Selector */}
      <Selector options={timePeriodOptions} value={selectedPeriod} onChange={setSelectedPeriod} />
    </div>
  )
}
