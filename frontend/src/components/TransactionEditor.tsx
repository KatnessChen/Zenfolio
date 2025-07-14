import { useState, useMemo, useCallback, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Dropdown, DropdownItem } from '@/components/ui/dropdown'
import DropdownTrigger from '@/components/ui/dropdown-trigger'
import { Label } from '@/components/ui/label'
import { TransactionService } from '@/services/transaction.service'
import { fetchCurrencies, fetchBrokers, fetchSymbols } from '@/services/meta.service'
import { useToast } from '@/hooks/useToast'
import { ROUTES } from '@/constants/routes'
import type { TransactionData } from '@/types'
import type { Currency, Broker, Symbol } from '@/services/meta.service'

export interface TransactionEditorProps {
  initial?: TransactionData
}

const defaultState: TransactionData = {
  id: '',
  transaction_date: '',
  symbol: '',
  trade_type: 'Buy',
  price: 0,
  quantity: 0,
  amount: 0,
  broker: '',
  currency: 'USD',
  exchange: '',
  user_notes: '',
}

export function TransactionEditor({ initial }: TransactionEditorProps) {
  const navigate = useNavigate()
  const { showToast } = useToast()
  const [loading, setLoading] = useState(true)
  const [transaction, setTransaction] = useState<TransactionData>(initial || defaultState)
  const [currencies, setCurrencies] = useState<Currency[]>([])
  const [brokers, setBrokers] = useState<Broker[]>([])
  const [symbols, setSymbols] = useState<Symbol[]>([])
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchCurrencies().then(setCurrencies)
    fetchBrokers().then(setBrokers)
    fetchSymbols().then(setSymbols)

    setLoading(false)
  }, [])

  useEffect(() => {
    setTransaction(initial || defaultState)
  }, [initial])

  const handleInputChange = useCallback((field: keyof TransactionData, value: string) => {
    if (field === 'price' || field === 'quantity') {
      setTransaction((prev) => ({ ...prev, [field]: value === '' ? 0 : Number(value) }))
    } else {
      setTransaction((prev) => ({ ...prev, [field]: value }))
    }
  }, [])

  const calculatedAmount = useMemo(
    () => transaction.price * transaction.quantity,
    [transaction.price, transaction.quantity]
  )

  const formatAmount = useMemo(() => {
    if (calculatedAmount === 0) return '$0.00'
    const formattedValue = calculatedAmount.toFixed(2)
    switch (transaction.trade_type) {
      case 'Sell':
        return `-$${formattedValue}`
      case 'Buy':
      case 'Dividends':
      default:
        return `$${formattedValue}`
    }
  }, [calculatedAmount, transaction.trade_type])

  const onCancel = useCallback(() => {
    window.history.back()
  }, [])

  const isFormValid =
    transaction.transaction_date &&
    transaction.symbol &&
    transaction.trade_type &&
    transaction.price &&
    transaction.quantity &&
    transaction.currency

  const handleCreateTransaction = async (tx: TransactionData) => {
    setSubmitting(true)
    setError(null)
    try {
      await TransactionService.importTransactions([{ ...tx }])

      // Show success toast
      showToast({
        type: 'success',
        title: 'Transaction Created',
        message: 'Transaction has been created successfully',
        duration: 3000,
      })

      // Navigate back to transactions page
      navigate(ROUTES.TRANSACTIONS)
    } catch (e: unknown) {
      const errorMessage = e instanceof Error ? e.message : 'Failed to create transaction'
      setError(errorMessage)
    } finally {
      setSubmitting(false)
    }
  }

  const handleUpdateTransaction = async (tx: TransactionData) => {
    setSubmitting(true)
    setError(null)
    try {
      // Use transaction_id as the path parameter for PUT /transaction-history/:id
      const transactionId = tx.transaction_id || tx.id
      await TransactionService.updateTransaction(transactionId, tx)

      // Show success toast
      showToast({
        type: 'success',
        title: 'Transaction Updated',
        message: 'Transaction has been updated successfully',
        duration: 3000,
      })

      // Navigate back to transactions page
      navigate(ROUTES.TRANSACTIONS)
    } catch (e: unknown) {
      const errorMessage = e instanceof Error ? e.message : 'Failed to update transaction'
      setError(errorMessage)
    } finally {
      setSubmitting(false)
    }
  }

  const isEdit = Boolean(initial?.transaction_id || initial?.id)

  const handleSubmit = async () => {
    if (!isFormValid) return

    const transactionWithAmount = { ...transaction, amount: calculatedAmount }

    if (transactionWithAmount.transaction_id || transactionWithAmount.id) {
      await handleUpdateTransaction(transactionWithAmount)
    } else {
      await handleCreateTransaction(transactionWithAmount)
    }
  }

  return (
    <Card>
      <CardContent className="p-8">
        <div className="space-y-6">
          {/* Trade Date */}
          <div>
            <Label htmlFor="trade-date-input">Trade Date *</Label>
            <Input
              id="trade-date-input"
              type="date"
              value={transaction.transaction_date}
              onChange={(e) => handleInputChange('transaction_date', e.target.value)}
              className="w-full"
            />
          </div>
          {/* Symbol */}
          <div>
            <Label htmlFor="symbol-dropdown">Symbol *</Label>
            <Dropdown
              trigger={
                <DropdownTrigger className="w-full">
                  {transaction.symbol || 'Select or type symbol'}
                </DropdownTrigger>
              }
              className="w-full"
            >
              {symbols.map((symbol) => (
                <DropdownItem
                  key={symbol.symbol}
                  onClick={() => handleInputChange('symbol', symbol.symbol)}
                >
                  {symbol.symbol} - {symbol.name}
                </DropdownItem>
              ))}
              {symbols.length === 0 && (
                <DropdownItem onClick={() => {}}>
                  {loading ? 'Loading symbols...' : 'No symbols available'}
                </DropdownItem>
              )}
            </Dropdown>
            <Input
              id="symbol-dropdown"
              type="text"
              placeholder="Or type symbol manually (e.g., AAPL, TSLA)"
              value={transaction.symbol}
              onChange={(e) => handleInputChange('symbol', e.target.value.toUpperCase())}
              className="w-full mt-2"
            />
          </div>
          {/* Trade Type */}
          <div>
            <Label htmlFor="trade-type-dropdown">Trade Type *</Label>
            <Dropdown
              trigger={
                <DropdownTrigger className="w-full">
                  {transaction.trade_type || 'Select trade type'}
                </DropdownTrigger>
              }
              className="w-full"
            >
              <DropdownItem onClick={() => handleInputChange('trade_type', 'Buy')}>
                Buy
              </DropdownItem>
              <DropdownItem onClick={() => handleInputChange('trade_type', 'Sell')}>
                Sell
              </DropdownItem>
              <DropdownItem onClick={() => handleInputChange('trade_type', 'Dividends')}>
                Dividends
              </DropdownItem>
            </Dropdown>
          </div>
          {/* Price and Quantity Row */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <Label htmlFor="price-input">Price *</Label>
              <Input
                id="price-input"
                type="number"
                step="0.01"
                placeholder="0.00"
                value={transaction.price}
                onChange={(e) => handleInputChange('price', e.target.value)}
                className="w-full"
              />
            </div>
            <div>
              <Label htmlFor="quantity-input">Quantity *</Label>
              <Input
                id="quantity-input"
                type="number"
                step="1"
                placeholder="0"
                value={transaction.quantity}
                onChange={(e) => handleInputChange('quantity', e.target.value)}
                className="w-full"
              />
            </div>
          </div>
          {/* Amount */}
          <div>
            <Label>Total Amount</Label>
            <div className="flex h-10 w-full rounded-md border border-input bg-muted/30 px-3 py-2 text-sm text-foreground items-center justify-end">
              {formatAmount}
            </div>
            <p className="text-xs text-muted-foreground mt-1">
              Auto-calculated from price × quantity
            </p>
          </div>
          {/* Currency */}
          <div>
            <Label htmlFor="currency-dropdown">Currency *</Label>
            <Dropdown
              trigger={
                <DropdownTrigger className="w-full">
                  {transaction.currency || 'Select currency'}
                </DropdownTrigger>
              }
              className="w-full"
            >
              {currencies.map((currency) => (
                <DropdownItem
                  key={currency.code}
                  onClick={() => handleInputChange('currency', currency.code)}
                >
                  {currency.code} - {currency.name}
                </DropdownItem>
              ))}
              {currencies.length === 0 && (
                <DropdownItem onClick={() => {}}>
                  {loading ? 'Loading currencies...' : 'No currencies available'}
                </DropdownItem>
              )}
            </Dropdown>
          </div>
          {/* Broker */}
          <div>
            <Label htmlFor="broker-dropdown">Broker</Label>
            <Dropdown
              trigger={
                <DropdownTrigger className="w-full">
                  {transaction.broker || 'Select broker'}
                </DropdownTrigger>
              }
              className="w-full"
            >
              {brokers.map((broker) => (
                <DropdownItem
                  key={broker.id}
                  onClick={() => handleInputChange('broker', broker.name)}
                >
                  {broker.name}
                </DropdownItem>
              ))}
              {brokers.length === 0 && (
                <DropdownItem onClick={() => {}}>
                  {loading ? 'Loading brokers...' : 'No brokers available'}
                </DropdownItem>
              )}
            </Dropdown>
            <Input
              id="broker-dropdown"
              type="text"
              placeholder="Or type broker name manually"
              value={transaction.broker}
              onChange={(e) => handleInputChange('broker', e.target.value)}
              className="w-full mt-2"
            />
          </div>
          {/* Notes */}
          <div>
            <Label htmlFor="notes-textarea">Notes</Label>
            <textarea
              id="notes-textarea"
              placeholder="Optional notes about this transaction"
              value={transaction.user_notes}
              onChange={(e) => handleInputChange('user_notes', e.target.value.slice(0, 100))}
              className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm text-foreground ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
              rows={3}
              maxLength={100}
            />
            <p className="text-xs text-muted-foreground mt-1 text-right">
              {100 - transaction.user_notes.length} characters remaining
            </p>
          </div>
        </div>
        {error && <div className="text-red-500 mt-4">{error}</div>}
        <div className="flex justify-end gap-4 mt-8">
          <Button variant="secondary" onClick={onCancel} className="px-6">
            Cancel
          </Button>
          <Button
            variant="default"
            onClick={handleSubmit}
            disabled={!isFormValid || submitting}
            className="px-6"
          >
            {submitting ? 'Saving...' : isEdit ? 'Save Changes' : 'Add Transaction'}
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}
