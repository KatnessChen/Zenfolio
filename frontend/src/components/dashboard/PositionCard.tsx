import React, { useEffect, useState, useRef } from 'react'
import { Link } from 'react-router-dom'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { formatCurrency, formatPercent } from '@/utils'
import { ROUTES } from '@/constants'
import { fetchAllHoldings } from '@/services/portfolio.service'
import { useToast } from '@/hooks/useToast'

interface Position {
  symbol: string
  quantity: number
  last: number
  marketValue: number
  unitCost: number
  totalCost: number
  gainLossAmount: number
  gainLossPercent: number
  annualizedReturnRate: number
}

export const PositionCard: React.FC = () => {
  const { showToast } = useToast()
  const [positions, setPositions] = useState<Position[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const didFetch = useRef(false)

  const formatNumber = (num: number, decimals: number = 2) => {
    return new Intl.NumberFormat('en-US', {
      minimumFractionDigits: decimals,
      maximumFractionDigits: decimals,
    }).format(num)
  }

  const formatChange = (amount: number) => {
    const sign = amount >= 0 ? '+' : ''
    return `${sign}${amount.toFixed(2)}`
  }

  useEffect(() => {
    if (didFetch.current) return
    didFetch.current = true

    const fetchData = async () => {
      try {
        setLoading(true)
        setError(null)

        const response = await fetchAllHoldings()

        // Transform API data to match our component interface
        const transformedPositions: Position[] = response.holdings.map((holding) => ({
          symbol: holding.symbol,
          quantity: holding.total_quantity,
          last: holding.current_price,
          marketValue: holding.market_value,
          unitCost: holding.unit_cost,
          totalCost: holding.total_cost,
          gainLossAmount: holding.unrealized_gain_loss,
          gainLossPercent: holding.simple_return_rate,
          annualizedReturnRate: holding.annualized_return_rate,
        }))

        setPositions(transformedPositions)
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Failed to fetch holdings data'
        setError(errorMessage)
        showToast({
          title: 'Error',
          message: errorMessage,
          type: 'error',
        })
      } finally {
        setLoading(false)
      }
    }
    fetchData()
  }, [showToast])

  const getChangeColor = (value: number) => {
    if (value > 0) return 'text-profit'
    if (value < 0) return 'text-loss'
    return 'text-muted-foreground'
  }

  if (loading) {
    return (
      <Card className="bg-card border-border/50">
        <CardHeader>
          <CardTitle className="text-lg text-foreground">Positions</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-center h-32">
            <div className="text-muted-foreground">Loading positions...</div>
          </div>
        </CardContent>
      </Card>
    )
  }

  if (error) {
    return (
      <Card className="bg-card border-border/50">
        <CardHeader>
          <CardTitle className="text-lg text-foreground">Positions</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-center h-32">
            <div className="text-muted-foreground">Failed to load positions</div>
          </div>
        </CardContent>
      </Card>
    )
  }

  if (positions.length === 0) {
    return (
      <Card className="bg-card border-border/50">
        <CardHeader>
          <CardTitle className="text-lg text-foreground">Positions</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-center h-32">
            <div className="text-muted-foreground">No positions found</div>
          </div>
        </CardContent>
      </Card>
    )
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
                <TableHead className="text-right">Annualized Return</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {positions.map((position) => (
                <TableRow key={position.symbol}>
                  <TableCell className="">
                    <Link
                      to={`${ROUTES.PORTFOLIO_HOLDING}/${position.symbol}`}
                      className="text-primary font-medium underline"
                    >
                      {position.symbol}
                    </Link>
                  </TableCell>
                  <TableCell className="text-right">{formatNumber(position.quantity, 4)}</TableCell>
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
                  <TableCell
                    className={`text-right ${getChangeColor(position.annualizedReturnRate)}`}
                  >
                    {formatPercent(position.annualizedReturnRate)}
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
