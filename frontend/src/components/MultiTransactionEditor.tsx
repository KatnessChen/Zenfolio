import { useState, useCallback, useEffect } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import {
  Card,
  CardContent,
  Button,
  Input,
  Dropdown,
  DropdownItem,
  DropdownTrigger,
  SearchableDropdown,
} from '@/components/ui'
import { PlusIcon, DeleteIcon } from '@/components/icons'
import { TransactionService } from '@/services/transaction.service'
import { fetchCurrencies, fetchBrokers, fetchSymbols } from '@/services/listing.service'
import { useToast } from '@/hooks/useToast'
import { ROUTES, TRADE_TYPE } from '@/constants'
import type { Currency, Broker, Symbol } from '@/constants'
import type { TransactionData, TradeType } from '@/types'

const createEmptyRow = (id: string): TransactionData => ({
  id,
  symbol: '',
  trade_type: TRADE_TYPE.BUY as TradeType,
  transaction_date: new Date().toISOString().split('T')[0],
  quantity: 0,
  price: 0,
  amount: 0,
  broker: '',
  currency: 'USD',
  exchange: '',
  user_notes: '',
})

export function MultiTransactionEditor() {
  const navigate = useNavigate()
  const location = useLocation()
  const { showToast } = useToast()

  // Check if we're editing an existing transaction
  const initialTransaction = location.state?.initial as TransactionData | undefined

  const [rows, setRows] = useState<TransactionData[]>(() => {
    if (initialTransaction) {
      // Pre-fill with existing transaction data
      return [
        {
          id: initialTransaction.transaction_id || '',
          transaction_id: initialTransaction.transaction_id,
          symbol: initialTransaction.symbol || '',
          transaction_date:
            initialTransaction.transaction_date || new Date().toISOString().split('T')[0],
          trade_type: initialTransaction.trade_type || TRADE_TYPE.BUY,
          price: initialTransaction.price || 0,
          quantity: initialTransaction.quantity || 0,
          amount: initialTransaction.amount || 0,
          broker: initialTransaction.broker || '',
          currency: initialTransaction.currency || 'USD',
          exchange: initialTransaction.exchange || '',
          user_notes: initialTransaction.user_notes || '',
        },
      ]
    }
    return [createEmptyRow('1')]
  })

  const [currencies, setCurrencies] = useState<Currency[]>([])
  const [brokers, setBrokers] = useState<Broker[]>([])
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [allSymbols, setAllSymbols] = useState<Symbol[]>([])

  useEffect(() => {
    const loadData = async () => {
      const [currenciesData, brokersData, symbolsData] = await Promise.all([
        fetchCurrencies(),
        fetchBrokers(),
        fetchSymbols(),
      ])
      setCurrencies(currenciesData)
      setBrokers(brokersData)
      setAllSymbols(symbolsData)
    }
    loadData()
  }, [])

  const handleRowChange = useCallback(
    (rowId: string, field: keyof TransactionData, value: string | number) => {
      setRows((prev) => prev.map((row) => (row.id === rowId ? { ...row, [field]: value } : row)))
    },
    []
  )

  const handleNumericInputChange = useCallback(
    (rowId: string, field: 'price' | 'quantity', value: string) => {
      const numericValue = value === '' ? 0 : Number(value)
      handleRowChange(rowId, field, numericValue)
    },
    [handleRowChange]
  )

  const handleSymbolChange = useCallback(
    (rowId: string, value: string) => {
      handleRowChange(rowId, 'symbol', value.toUpperCase())
    },
    [handleRowChange]
  )

  const addRow = useCallback(() => {
    const newId = (rows.length + 1).toString()
    setRows((prev) => [...prev, createEmptyRow(newId)])
  }, [rows.length])

  const removeRow = useCallback(
    (rowId: string) => {
      if (rows.length > 1) {
        setRows((prev) => prev.filter((row) => row.id !== rowId))
      }
    },
    [rows.length]
  )

  const onCancel = useCallback(() => {
    navigate(ROUTES.TRANSACTIONS)
  }, [navigate])

  const validateRows = () => {
    return rows.every(
      (row) =>
        row.transaction_date &&
        row.symbol &&
        row.trade_type &&
        row.price > 0 &&
        row.quantity > 0 &&
        row.currency
    )
  }

  const handleSubmit = async () => {
    if (!validateRows()) {
      setError('Please fill in all required fields for all rows')
      return
    }

    setSubmitting(true)
    setError(null)

    try {
      if (initialTransaction) {
        // We're in edit mode - update existing transaction
        const updatedTransaction: TransactionData = {
          ...rows[0],
          id: rows[0].id || '',
          transaction_id: rows[0].transaction_id,
          symbol: rows[0].symbol.toUpperCase(),
          amount: rows[0].price * rows[0].quantity,
          exchange: rows[0].exchange || '', // Preserve existing exchange if available
        }

        // Use PUT to update the transaction
        await TransactionService.updateTransaction(
          updatedTransaction.transaction_id || '',
          updatedTransaction
        )

        showToast({
          type: 'success',
          title: 'Transaction Updated',
          message: 'Transaction updated successfully',
          duration: 3000,
        })
      } else {
        // Create new transactions mode
        const transactions: TransactionData[] = rows.map((row) => ({
          id: '',
          transaction_date: row.transaction_date,
          symbol: row.symbol.toUpperCase(),
          trade_type: row.trade_type,
          price: row.price,
          quantity: row.quantity,
          amount: row.price * row.quantity,
          broker: row.broker,
          currency: row.currency,
          user_notes: row.user_notes,
          exchange: '', // Will be populated by backend
        }))

        await TransactionService.importTransactions(transactions)

        showToast({
          type: 'success',
          title: 'Transactions Created',
          message: `${transactions.length} transaction(s) created successfully`,
          duration: 3000,
        })
      }

      navigate(ROUTES.TRANSACTIONS)
    } catch (e: unknown) {
      const errorMessage = e instanceof Error ? e.message : 'Failed to process transaction'
      setError(errorMessage)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <Card>
      <CardContent className="p-6">
        <div className="space-y-4">
          <div className="grid grid-cols-12 gap-2 text-sm font-medium text-muted-foreground border-b pb-2">
            <div className="col-span-2">Date *</div>
            <div className="col-span-2">Symbol *</div>
            <div className="col-span-1">Type *</div>
            <div className="col-span-1">Price *</div>
            <div className="col-span-1">Quantity *</div>
            <div className="col-span-1">Amount</div>
            <div className="col-span-1">Currency *</div>
            <div className="col-span-2">Broker</div>
            {!initialTransaction && <div className="col-span-1">Actions</div>}
          </div>

          {/* Transaction Rows */}
          {rows.map((row, index) => (
            <div key={row.id} className="grid grid-cols-12 gap-2 items-center">
              {/* Date */}
              <div className="col-span-2">
                <Input
                  type="date"
                  value={row.transaction_date}
                  onChange={(e) => handleRowChange(row.id, 'transaction_date', e.target.value)}
                  className="w-full"
                />
              </div>

              {/* Symbol */}
              <div className="col-span-2">
                <SearchableDropdown
                  value={row.symbol}
                  onChange={(value: string) => handleSymbolChange(row.id, value)}
                  options={allSymbols.map((symbol) => ({
                    id: symbol.symbol,
                    label: symbol.symbol,
                    value: symbol.symbol,
                    subtitle: symbol.name,
                  }))}
                  placeholder="Search symbols..."
                  className="w-full"
                  noResultsText="No symbols found for"
                  allowCustomValue={true}
                />
              </div>

              {/* Trade Type */}
              <div className="col-span-1">
                <Dropdown
                  trigger={
                    <DropdownTrigger className="w-full text-sm">{row.trade_type}</DropdownTrigger>
                  }
                >
                  <DropdownItem
                    onClick={() => handleRowChange(row.id, 'trade_type', TRADE_TYPE.BUY)}
                  >
                    {TRADE_TYPE.BUY}
                  </DropdownItem>
                  <DropdownItem
                    onClick={() => handleRowChange(row.id, 'trade_type', TRADE_TYPE.SELL)}
                  >
                    {TRADE_TYPE.SELL}
                  </DropdownItem>
                  <DropdownItem
                    onClick={() => handleRowChange(row.id, 'trade_type', TRADE_TYPE.DIVIDEND)}
                  >
                    {TRADE_TYPE.DIVIDEND}
                  </DropdownItem>
                </Dropdown>
              </div>

              {/* Price */}
              <div className="col-span-1">
                <Input
                  type="number"
                  step="0.01"
                  placeholder="0.00"
                  value={row.price || ''}
                  onChange={(e) => handleNumericInputChange(row.id, 'price', e.target.value)}
                  className="w-full"
                />
              </div>

              {/* Quantity */}
              <div className="col-span-1">
                <Input
                  type="number"
                  step="1"
                  placeholder="0"
                  value={row.quantity || ''}
                  onChange={(e) => handleNumericInputChange(row.id, 'quantity', e.target.value)}
                  className="w-full"
                />
              </div>

              {/* Amount (calculated) */}
              <div className="col-span-1">
                <div className="flex items-center text-sm text-muted-foreground">
                  ${(row.price * row.quantity).toFixed(2)}
                </div>
              </div>

              {/* Currency */}
              <div className="col-span-1">
                <Dropdown
                  trigger={
                    <DropdownTrigger className="w-full text-sm">{row.currency}</DropdownTrigger>
                  }
                >
                  {currencies.map((currency) => (
                    <DropdownItem
                      key={currency.code}
                      onClick={() => handleRowChange(row.id, 'currency', currency.code)}
                    >
                      {currency.code}
                    </DropdownItem>
                  ))}
                </Dropdown>
              </div>

              {/* Broker */}
              <div className="col-span-2">
                <SearchableDropdown
                  value={row.broker}
                  onChange={(value: string) => handleRowChange(row.id, 'broker', value)}
                  options={brokers.map((broker) => ({
                    id: broker.id,
                    label: broker.name,
                    value: broker.name,
                  }))}
                  placeholder="Search brokers..."
                  className="w-full"
                  noResultsText="No brokers found for"
                  allowCustomValue={true}
                />
              </div>

              {/* Actions - only show in create mode */}
              {!initialTransaction && (
                <div className="col-span-1">
                  <div className="flex gap-1">
                    {index === rows.length - 1 && (
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={addRow}
                        className="h-8 w-8 p-0"
                        title="Add row"
                      >
                        <PlusIcon className="h-4 w-4" />
                      </Button>
                    )}
                    {rows.length > 1 && (
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={() => removeRow(row.id)}
                        className="h-8 w-8 p-0"
                        title="Remove row"
                      >
                        <DeleteIcon className="h-4 w-4" />
                      </Button>
                    )}
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>

        {error && <div className="text-red-500 mt-4 text-sm">{error}</div>}

        <div className="flex justify-end mt-6 pt-4 border-t gap-4">
          <Button variant="secondary" onClick={onCancel} className="px-6">
            Cancel
          </Button>
          <Button
            variant="default"
            onClick={handleSubmit}
            disabled={!validateRows() || submitting}
            className="px-6"
          >
            {submitting
              ? 'Saving...'
              : initialTransaction
                ? 'Update Transaction'
                : `Add ${rows.length} Transaction${rows.length !== 1 ? 's' : ''}`}
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}
