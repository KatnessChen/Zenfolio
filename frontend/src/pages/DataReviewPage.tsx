import { useState, useEffect } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { ROUTES } from '@/constants/routes'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Dropdown, DropdownItem } from '@/components/ui/dropdown'
import DropdownTrigger from '@/components/ui/dropdown-trigger'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { DeleteIcon, PlusIcon } from '@/components/icons'
import { Title } from '@/components/ui/title'

interface ProcessedTransaction {
  id: string
  ticker: string
  tradeType: 'Buy' | 'Sell' | 'Dividend'
  quantity: number
  price: number
  amount: number
  tradeDate: string
  broker: string
  currency: string
  userNotes: string
  confidence: number
}

interface LocationState {
  files: string[]
  processedData: ProcessedTransaction[]
}

export default function DataReviewPage() {
  const navigate = useNavigate()
  const location = useLocation()
  const state = location.state as LocationState

  const [currentFileIndex, setCurrentFileIndex] = useState(0)
  const [transactions, setTransactions] = useState<ProcessedTransaction[]>(
    state?.processedData || []
  )
  const [selectedRows, setSelectedRows] = useState<Set<string>>(new Set())

  useEffect(() => {
    if (!state?.files || !state?.processedData) {
      navigate(ROUTES.TRANSACTIONS_UPLOAD)
    }
  }, [state, navigate])

  const handleCellEdit = (
    id: string,
    field: keyof ProcessedTransaction,
    value: string | number
  ) => {
    setTransactions((prev) => prev.map((t) => (t.id === id ? { ...t, [field]: value } : t)))
  }

  const handleAddRow = () => {
    const newTransaction: ProcessedTransaction = {
      id: Date.now().toString(),
      ticker: '',
      tradeType: 'Buy',
      quantity: 0,
      price: 0,
      amount: 0,
      tradeDate: new Date().toISOString().split('T')[0],
      broker: '',
      currency: 'USD',
      userNotes: '',
      confidence: 1.0,
    }
    setTransactions((prev) => [...prev, newTransaction])
  }

  const handleDeleteRows = () => {
    setTransactions((prev) => prev.filter((t) => !selectedRows.has(t.id)))
    setSelectedRows(new Set())
  }

  const handleRowSelect = (id: string) => {
    const newSelected = new Set(selectedRows)
    if (newSelected.has(id)) {
      newSelected.delete(id)
    } else {
      newSelected.add(id)
    }
    setSelectedRows(newSelected)
  }

  const handleConfirmImport = () => {
    // TODO: Send data to backend
    console.log('Importing transactions:', transactions)
    navigate(ROUTES.TRANSACTIONS, {
      state: { message: 'Transactions imported successfully!' },
    })
  }

  const handleCancel = () => {
    navigate(ROUTES.TRANSACTIONS)
  }

  const files = state?.files || []
  const hasMultipleFiles = files.length > 1

  return (
    <div className="min-h-screen bg-background">
      {/* Main Content */}
      <main className="container mx-auto py-6 px-4">
        {/* Page Title */}
        <div className="mb-6">
          <Title as="h1" className="mb-2">
            Data Review
          </Title>
          <p className="text-muted-foreground">
            Please review the extracted data below. Make any necessary corrections before
            confirming.
          </p>
        </div>

        {/* File Navigation */}
        {hasMultipleFiles && (
          <Card className="mb-6">
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <Button
                  variant="secondary"
                  onClick={() => setCurrentFileIndex(Math.max(0, currentFileIndex - 1))}
                  disabled={currentFileIndex === 0}
                >
                  ← Prev File
                </Button>

                <div className="flex items-center space-x-2">
                  <span className="text-sm text-muted-foreground">
                    File {currentFileIndex + 1} of {files.length} - Reviewing
                  </span>
                  {/* TODO: Add warning/error indicators */}
                </div>

                <Button
                  variant="secondary"
                  onClick={() =>
                    setCurrentFileIndex(Math.min(files.length - 1, currentFileIndex + 1))
                  }
                  disabled={currentFileIndex === files.length - 1}
                >
                  Next File →
                </Button>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Main Review Layout */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
          {/* Left Half - Screenshot Display */}
          <Card className="h-fit">
            <CardContent className="p-4">
              <h2 className="text-lg font-medium text-foreground mb-4">Original Screenshot</h2>
              <div className="bg-muted rounded-lg p-8 text-center">
                <p className="text-muted-foreground">
                  {files[currentFileIndex] || 'No file selected'}
                </p>
                <p className="text-sm text-muted-foreground mt-2">
                  Screenshot preview would be displayed here
                </p>
              </div>
            </CardContent>
          </Card>

          {/* Right Half - Extracted Data */}
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-lg font-medium text-foreground">Extracted Data</h2>
                <div className="flex gap-2">
                  <Button variant="secondary" size="sm" onClick={handleAddRow}>
                    <PlusIcon size={16} />
                    <span>Add Row</span>
                  </Button>
                  {selectedRows.size > 0 && (
                    <Button
                      variant="secondary"
                      size="sm"
                      onClick={handleDeleteRows}
                      className="text-destructive hover:bg-destructive/10"
                    >
                      <DeleteIcon size={16} />
                      <span>Delete Selected</span>
                    </Button>
                  )}
                </div>
              </div>

              {/* Data Table */}
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead className="w-12">
                        <input
                          type="checkbox"
                          checked={
                            selectedRows.size === transactions.length && transactions.length > 0
                          }
                          onChange={(e) => {
                            if (e.target.checked) {
                              setSelectedRows(new Set(transactions.map((t) => t.id)))
                            } else {
                              setSelectedRows(new Set())
                            }
                          }}
                        />
                      </TableHead>
                      <TableHead>Trade Date</TableHead>
                      <TableHead>Symbol</TableHead>
                      <TableHead>Trade Type</TableHead>
                      <TableHead>Price</TableHead>
                      <TableHead>Quantity</TableHead>
                      <TableHead>Amount</TableHead>
                      <TableHead>Broker</TableHead>
                      <TableHead>Notes</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {transactions.map((transaction) => (
                      <TableRow key={transaction.id}>
                        <TableCell>
                          <input
                            type="checkbox"
                            checked={selectedRows.has(transaction.id)}
                            onChange={() => handleRowSelect(transaction.id)}
                          />
                        </TableCell>
                        <TableCell>
                          <Input
                            type="date"
                            value={transaction.tradeDate}
                            onChange={(e) =>
                              handleCellEdit(transaction.id, 'tradeDate', e.target.value)
                            }
                            className="w-full"
                          />
                        </TableCell>
                        <TableCell>
                          <Input
                            value={transaction.ticker}
                            onChange={(e) =>
                              handleCellEdit(transaction.id, 'ticker', e.target.value)
                            }
                            className="w-full font-mono"
                          />
                        </TableCell>
                        <TableCell>
                          <Dropdown
                            trigger={
                              <DropdownTrigger className="w-full">
                                {transaction.tradeType}
                              </DropdownTrigger>
                            }
                          >
                            <DropdownItem
                              onClick={() => handleCellEdit(transaction.id, 'tradeType', 'Buy')}
                            >
                              Buy
                            </DropdownItem>
                            <DropdownItem
                              onClick={() => handleCellEdit(transaction.id, 'tradeType', 'Sell')}
                            >
                              Sell
                            </DropdownItem>
                            <DropdownItem
                              onClick={() =>
                                handleCellEdit(transaction.id, 'tradeType', 'Dividend')
                              }
                            >
                              Dividend
                            </DropdownItem>
                          </Dropdown>
                        </TableCell>
                        <TableCell>
                          <Input
                            type="number"
                            step="0.01"
                            value={transaction.price}
                            onChange={(e) =>
                              handleCellEdit(transaction.id, 'price', parseFloat(e.target.value))
                            }
                            className="w-full"
                          />
                        </TableCell>
                        <TableCell>
                          <Input
                            type="number"
                            value={transaction.quantity}
                            onChange={(e) =>
                              handleCellEdit(transaction.id, 'quantity', parseInt(e.target.value))
                            }
                            className="w-full"
                          />
                        </TableCell>
                        <TableCell>
                          <Input
                            type="number"
                            step="0.01"
                            value={transaction.amount}
                            onChange={(e) =>
                              handleCellEdit(transaction.id, 'amount', parseFloat(e.target.value))
                            }
                            className={`w-full ${
                              transaction.tradeType === 'Buy'
                                ? 'text-[#9DC0B2]'
                                : transaction.tradeType === 'Sell'
                                  ? 'text-[#E6B9B3]'
                                  : 'text-foreground'
                            }`}
                          />
                        </TableCell>
                        <TableCell>
                          <Input
                            value={transaction.broker}
                            onChange={(e) =>
                              handleCellEdit(transaction.id, 'broker', e.target.value)
                            }
                            className="w-full"
                          />
                        </TableCell>
                        <TableCell>
                          <Input
                            value={transaction.userNotes}
                            onChange={(e) =>
                              handleCellEdit(transaction.id, 'userNotes', e.target.value)
                            }
                            className="w-full"
                            placeholder="Add notes..."
                          />
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Action Buttons */}
        <div className="flex items-center justify-end space-x-4">
          <Button variant="secondary" onClick={handleCancel}>
            Cancel
          </Button>
          <Button variant="default" onClick={handleConfirmImport}>
            Confirm & Import
          </Button>
        </div>
      </main>
    </div>
  )
}
