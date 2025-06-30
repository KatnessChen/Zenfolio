import { useState, useEffect, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import { ROUTES } from '@/constants'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { EditIcon, DeleteIcon } from '@/components/icons'
import { Title } from '@/components/ui/title'

interface BatchTransaction {
  id: string
  tradeDate: string
  symbol: string
  tradeType: 'Buy' | 'Sell' | 'Dividend' | ''
  price: string
  quantity: string
  amount: string
  broker: string
  currency: string
  notes: string
}

export default function BatchReviewPage() {
  const navigate = useNavigate()
  const [batch, setBatch] = useState<BatchTransaction[]>([])

  useEffect(() => {
    // Load batch from localStorage
    const storedBatch = localStorage.getItem('transaction-batch')
    if (storedBatch) {
      setBatch(JSON.parse(storedBatch))
    } else {
      // If no batch data, redirect back to manual transaction page
      navigate(ROUTES.TRANSACTIONS_MANUAL_ADD)
    }
  }, [navigate])

  const handleEdit = useCallback(
    (transactionId: string) => {
      // For now, just show an alert. In a full implementation, this would open an edit modal
      const transaction = batch.find((t) => t.id === transactionId)
      if (transaction) {
        alert(`Edit functionality would open for: ${transaction.symbol} ${transaction.tradeType}`)
      }
    },
    [batch]
  )

  const handleDelete = useCallback((transactionId: string) => {
    setBatch((prev) => {
      const updated = prev.filter((t) => t.id !== transactionId)
      localStorage.setItem('transaction-batch', JSON.stringify(updated))
      return updated
    })
  }, [])

  const handleAddMore = useCallback(() => {
    // Update localStorage with current batch and navigate back
    localStorage.setItem('transaction-batch', JSON.stringify(batch))
    navigate(ROUTES.TRANSACTIONS_MANUAL_ADD)
  }, [navigate, batch])

  const handleDiscardAll = useCallback(() => {
    setBatch([])
    localStorage.removeItem('transaction-batch')
    navigate(ROUTES.TRANSACTIONS) // or wherever the user should go after discarding
  }, [navigate])

  const handleConfirmAll = useCallback(() => {
    // TODO: Implement saving all transactions to backend
    console.log('Saving all transactions:', batch)

    // Clear batch and navigate to success or transactions page
    localStorage.removeItem('transaction-batch')
    navigate(ROUTES.TRANSACTIONS)
  }, [navigate, batch])

  const getTradeTypeColor = (tradeType: string) => {
    switch (tradeType) {
      case 'Buy':
        return 'text-green-600'
      case 'Sell':
        return 'text-red-600'
      case 'Dividend':
        return 'text-blue-600'
      default:
        return 'text-foreground'
    }
  }

  const getAmountDisplay = (amount: string, tradeType: string) => {
    const value = parseFloat(amount) || 0
    const formattedValue = value.toFixed(2)

    switch (tradeType) {
      case 'Sell':
        return `-$${formattedValue}`
      case 'Buy':
      case 'Dividend':
      default:
        return `$${formattedValue}`
    }
  }

  const totalValue = batch.reduce((sum, transaction) => {
    const amount = parseFloat(transaction.amount) || 0
    return transaction.tradeType === 'Sell' ? sum - amount : sum + amount
  }, 0)

  if (batch.length === 0) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <p className="text-lg text-muted-foreground">No transactions in batch.</p>
          <Button onClick={() => navigate(ROUTES.TRANSACTIONS_MANUAL_ADD)} className="mt-4">
            Add Transactions
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <main className="container mx-auto py-8 px-4">
        {/* Page Title */}
        <div className="text-center mb-8">
          <Title as="h1" className="mb-4">
            Batch Review
          </Title>
          <p className="text-lg text-muted-foreground">
            Please review the transactions you've added. You can make corrections or add more before
            confirming.
          </p>
        </div>

        <Card>
          <CardContent className="p-6">
            {/* Summary */}
            <div className="mb-6 p-4 bg-muted/50 rounded-lg">
              <div className="flex justify-between items-center">
                <span className="text-sm font-medium text-foreground">
                  Total transactions: {batch.length}
                </span>
                <span className="text-sm font-medium text-foreground">
                  Total value: {totalValue >= 0 ? '$' : '-$'}
                  {Math.abs(totalValue).toFixed(2)}
                </span>
              </div>
            </div>

            {/* Transactions Table */}
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-border">
                    <th className="text-left py-3 px-2 text-sm font-medium text-foreground">
                      Trade Date
                    </th>
                    <th className="text-left py-3 px-2 text-sm font-medium text-foreground">
                      Symbol
                    </th>
                    <th className="text-left py-3 px-2 text-sm font-medium text-foreground">
                      Trade Type
                    </th>
                    <th className="text-left py-3 px-2 text-sm font-medium text-foreground">
                      Price
                    </th>
                    <th className="text-left py-3 px-2 text-sm font-medium text-foreground">
                      Quantity
                    </th>
                    <th className="text-left py-3 px-2 text-sm font-medium text-foreground">
                      Amount
                    </th>
                    <th className="text-left py-3 px-2 text-sm font-medium text-foreground">
                      Currency
                    </th>
                    <th className="text-left py-3 px-2 text-sm font-medium text-foreground">
                      Broker
                    </th>
                    <th className="text-left py-3 px-2 text-sm font-medium text-foreground">
                      Notes
                    </th>
                    <th className="text-left py-3 px-2 text-sm font-medium text-foreground">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {batch.map((transaction) => (
                    <tr key={transaction.id} className="border-b border-border/50">
                      <td className="py-3 px-2 text-sm text-foreground">{transaction.tradeDate}</td>
                      <td className="py-3 px-2 text-sm font-medium text-foreground">
                        {transaction.symbol}
                      </td>
                      <td
                        className={`py-3 px-2 text-sm font-medium ${getTradeTypeColor(transaction.tradeType)}`}
                      >
                        {transaction.tradeType}
                      </td>
                      <td className="py-3 px-2 text-sm text-foreground">
                        ${parseFloat(transaction.price).toFixed(2)}
                      </td>
                      <td className="py-3 px-2 text-sm text-foreground">{transaction.quantity}</td>
                      <td className="py-3 px-2 text-sm font-medium text-foreground">
                        {getAmountDisplay(transaction.amount, transaction.tradeType)}
                      </td>
                      <td className="py-3 px-2 text-sm text-foreground">{transaction.currency}</td>
                      <td className="py-3 px-2 text-sm text-foreground">
                        {transaction.broker || '-'}
                      </td>
                      <td className="py-3 px-2 text-sm text-muted-foreground max-w-32 truncate">
                        {transaction.notes || '-'}
                      </td>
                      <td className="py-3 px-2">
                        <div className="flex gap-2">
                          <button
                            onClick={() => handleEdit(transaction.id)}
                            className="p-1 text-muted-foreground hover:text-foreground transition-colors"
                            title="Edit transaction"
                          >
                            <EditIcon className="h-4 w-4" />
                          </button>
                          <button
                            onClick={() => handleDelete(transaction.id)}
                            className="p-1 text-muted-foreground hover:text-red-600 transition-colors"
                            title="Delete transaction"
                          >
                            <DeleteIcon className="h-4 w-4" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Action Buttons */}
            <div className="flex flex-col sm:flex-row justify-center gap-4 mt-8">
              <Button variant="secondary" onClick={handleAddMore} className="px-6">
                Add More Transactions
              </Button>
              <Button variant="secondary" onClick={handleDiscardAll} className="px-6">
                Discard All
              </Button>
              <Button variant="default" onClick={handleConfirmAll} className="px-6">
                Confirm All & Save
              </Button>
            </div>
          </CardContent>
        </Card>
      </main>
    </div>
  )
}
