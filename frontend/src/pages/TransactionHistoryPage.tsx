import type { Transaction } from '@/types'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { formatCurrency, formatDate } from '@/utils'

// Mock data for demonstration
const mockTransactions: Transaction[] = [
  {
    id: '1',
    ticker: 'AAPL',
    tickerLabel: 'Apple Inc.',
    tradeType: 'Buy',
    quantity: 100,
    price: 150.25,
    tradeDate: '2024-01-15',
    exchange: 'NASDAQ',
    currency: 'USD',
  },
  {
    id: '2',
    ticker: 'GOOGL',
    tickerLabel: 'Alphabet Inc.',
    tradeType: 'Sell',
    quantity: 50,
    price: 2750.8,
    tradeDate: '2024-01-14',
    exchange: 'NASDAQ',
    currency: 'USD',
  },
]

export default function TransactionHistoryPage() {
  return (
    <div className="container mx-auto py-10">
      <div className="space-y-8">
        <div className="text-center space-y-4">
          <h1 className="text-3xl font-bold tracking-tight">Transaction History</h1>
          <p className="text-lg text-muted-foreground">
            View and manage your extracted trading transactions
          </p>
        </div>

        <div className="space-y-4">
          {mockTransactions.map((transaction) => (
            <Card key={transaction.id}>
              <CardHeader>
                <div className="flex justify-between items-start">
                  <div>
                    <CardTitle className="text-lg">
                      {transaction.ticker}
                      {transaction.tickerLabel && (
                        <span className="text-sm font-normal text-muted-foreground ml-2">
                          ({transaction.tickerLabel})
                        </span>
                      )}
                    </CardTitle>
                    <CardDescription>
                      {transaction.tradeType} • {formatDate(transaction.tradeDate)}
                    </CardDescription>
                  </div>
                  <div className="text-right">
                    <div className="text-lg font-semibold">
                      {formatCurrency(transaction.price * transaction.quantity)}
                    </div>
                    <div className="text-sm text-muted-foreground">
                      {transaction.quantity} shares @ {formatCurrency(transaction.price)}
                    </div>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="flex justify-between text-sm text-muted-foreground">
                  <span>Exchange: {transaction.exchange}</span>
                  <span>Currency: {transaction.currency}</span>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </div>
  )
}
