import React from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { formatCurrency } from '@/utils'

interface Position {
  symbol: string
  quantity: number
  last: number
  changeAmount: number
  changePercent: number
  marketValue: number
  unitCost: number
  totalCost: number
  gainLossAmount: number
  gainLossPercent: number
}

export const PositionCard: React.FC = () => {
  const formatNumber = (num: number, decimals: number = 2) => {
    return new Intl.NumberFormat('en-US', {
      minimumFractionDigits: decimals,
      maximumFractionDigits: decimals,
    }).format(num)
  }

  const formatPercent = (percent: number) => {
    const sign = percent >= 0 ? '+' : ''
    return `${sign}${percent.toFixed(2)}%`
  }

  const formatChange = (amount: number) => {
    const sign = amount >= 0 ? '+' : ''
    return `${sign}${amount.toFixed(2)}`
  }

  // Mock data based on the screenshot
  const positions: Position[] = [
    {
      symbol: 'GOOGL',
      quantity: 32.29027,
      last: 174.15,
      changeAmount: 0.61,
      changePercent: 0.35,
      marketValue: 5623.35,
      unitCost: 170.40923,
      totalCost: 5502.56,
      gainLossAmount: 120.79,
      gainLossPercent: 2.2,
    },
    {
      symbol: 'NVDA',
      quantity: 61.37719,
      last: 157.33,
      changeAmount: 2.31,
      changePercent: 1.49,
      marketValue: 9656.47,
      unitCost: 118.03799,
      totalCost: 7244.84,
      gainLossAmount: 2411.63,
      gainLossPercent: 33.29,
    },
    {
      symbol: 'ON',
      quantity: 9.36329,
      last: 53.27,
      changeAmount: -0.38,
      changePercent: -0.71,
      marketValue: 498.78,
      unitCost: 53.40003,
      totalCost: 500.0,
      gainLossAmount: -1.22,
      gainLossPercent: -0.24,
    },
    {
      symbol: 'PL',
      quantity: 83,
      last: 6.52,
      changeAmount: 0.46,
      changePercent: 7.59,
      marketValue: 541.16,
      unitCost: 6.03916,
      totalCost: 501.25,
      gainLossAmount: 39.91,
      gainLossPercent: 7.96,
    },
    {
      symbol: 'QQQ',
      quantity: 22.94481,
      last: 549.38,
      changeAmount: 3.16,
      changePercent: 0.58,
      marketValue: 12605.42,
      unitCost: 506.55508,
      totalCost: 11622.81,
      gainLossAmount: 982.61,
      gainLossPercent: 8.45,
    },
    {
      symbol: 'SGOV',
      quantity: 42.84301,
      last: 100.68,
      changeAmount: 0.03,
      changePercent: 0.03,
      marketValue: 4313.43,
      unitCost: 100.64792,
      totalCost: 4312.06,
      gainLossAmount: 1.37,
      gainLossPercent: 0.03,
    },
  ]

  const getChangeColor = (value: number) => {
    if (value > 0) return 'text-profit'
    if (value < 0) return 'text-loss'
    return 'text-muted-foreground'
  }

  return (
    <Card className="bg-card border-border/50">
      <CardHeader>
        <CardTitle className="text-lg text-foreground">Positions</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Symbol</TableHead>
                <TableHead className="text-right">Quantity</TableHead>
                <TableHead className="text-right">Last</TableHead>
                <TableHead className="text-right">Market Value</TableHead>
                <TableHead className="text-right">Gain/Loss($)</TableHead>
                <TableHead className="text-right">Gain/Loss(%)</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {positions.map((position) => (
                <TableRow key={position.symbol}>
                  <TableCell className="">
                    <span className="text-primary font-medium">{position.symbol}</span>
                  </TableCell>
                  <TableCell className="text-right">{formatNumber(position.quantity, 5)}</TableCell>
                  <TableCell className="text-right">{formatCurrency(position.last)}</TableCell>
                  <TableCell className="text-right">
                    {formatCurrency(position.marketValue)}
                  </TableCell>
                  <TableCell className={`text-right ${getChangeColor(position.gainLossAmount)}`}>
                    {formatChange(position.gainLossAmount)}
                  </TableCell>
                  <TableCell className={`text-right ${getChangeColor(position.gainLossPercent)}`}>
                    {formatPercent(position.gainLossPercent)}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  )
}
