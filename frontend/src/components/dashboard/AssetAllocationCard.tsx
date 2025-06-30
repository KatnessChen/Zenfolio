import React from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { PieChart, Pie, Cell, ResponsiveContainer } from 'recharts'
import { formatCurrency } from '@/utils'

interface AllocationData {
  name: string
  value: number
  percentage: number
  amount: number
  color: string
  label: string
}

export const AssetAllocationCard: React.FC = () => {
  // Mock data - in real app this would come from API
  const allocationData: AllocationData[] = [
    {
      name: 'Equities',
      value: 45,
      percentage: 45,
      amount: 39352.61,
      color: 'hsl(var(--chart-1))',
      label: `Equities ${formatCurrency(39352.61)} (45%)`,
    },
    {
      name: 'Cryptocurrency',
      value: 25,
      percentage: 25,
      amount: 21862.56,
      color: 'hsl(var(--chart-2))',
      label: `Cryptocurrency ${formatCurrency(21862.56)} (25%)`,
    },
    {
      name: 'ETFs',
      value: 20,
      percentage: 20,
      amount: 17490.05,
      color: 'hsl(var(--chart-3))',
      label: `ETFs ${formatCurrency(17490.05)} (20%)`,
    },
    {
      name: 'Bonds',
      value: 7,
      percentage: 7,
      amount: 6121.52,
      color: 'hsl(var(--chart-4))',
      label: `Bonds ${formatCurrency(6121.52)} (7%)`,
    },
    {
      name: 'Cash',
      value: 3,
      percentage: 3,
      amount: 2623.51,
      color: 'hsl(var(--chart-5))',
      label: `Cash ${formatCurrency(2623.51)} (3%)`,
    },
  ]

  return (
    <Card className="bg-card border-border/50">
      <CardHeader>
        <CardTitle className="text-lg text-foreground">Asset Allocation</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Recharts Donut Chart with Built-in Labels */}
        <div className="w-full h-60 focus:outline-none" style={{ fontSize: '0.8em' }}>
          <ResponsiveContainer width="100%" height="100%" style={{ outline: 'none' }}>
            <PieChart tabIndex={-1} style={{ outline: 'none' }}>
              <Pie
                data={allocationData}
                cx="50%"
                cy="50%"
                labelLine={true}
                label={(entry) => entry.label}
                outerRadius={80}
                innerRadius={48}
                dataKey="value"
                tabIndex={-1}
                style={{ outline: 'none' }}
              >
                {allocationData.map((entry, index) => (
                  <Cell
                    key={`cell-${index}`}
                    fill={entry.color}
                    tabIndex={-1}
                    style={{ outline: 'none' }}
                  />
                ))}
              </Pie>
            </PieChart>
          </ResponsiveContainer>
        </div>
      </CardContent>
    </Card>
  )
}
