import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { ROUTES } from '@/constants'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Dropdown, DropdownItem } from '@/components/ui/dropdown'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { DeleteIcon, PlusIcon, ClockIcon, SpinnerIcon, CheckIcon, XIcon } from '@/components/icons'
import { Title } from '@/components/ui/title'
import type { RootState } from '@/store'
import type { FileProcessingState, TradeType } from '@/types'

interface ProcessedTransaction {
  id: string
  symbol: string
  type: TradeType
  quantity: number
  price: number
  amount: number
  transaction_date: string
  broker: string
  currency: string
  user_notes: string
  account: string
  exchange: string
}

export default function ExtractedDataReviewPage() {
  const navigate = useNavigate()

  // Get data from Redux store
  const { files, fileStates, extractResults } = useSelector(
    (state: RootState) => state.fileProcessing
  )

  const [currentFileIndex, setCurrentFileIndex] = useState(-1)
  const [allTransactions, setAllTransactions] = useState<ProcessedTransaction[]>([])
  const [selectedRows, setSelectedRows] = useState<Set<string>>(new Set())
  const [showDeleteConfirmation, setShowDeleteConfirmation] = useState(false)

  useEffect(() => {
    // Redirect if no files in store
    if (files.length === 0) {
      navigate(ROUTES.TRANSACTIONS_UPLOAD)
      return
    }

    // Find the first available (completed) file index
    const firstAvailableIndex = fileStates.findIndex((fs) => fs.status === 'completed')
    setCurrentFileIndex(firstAvailableIndex >= 0 ? firstAvailableIndex : 0)
  }, [navigate, files, fileStates])

  // Update transactions when extract results change
  useEffect(() => {
    const convertedTransactions: ProcessedTransaction[] = []

    extractResults.forEach((result, fileIndex) => {
      result.transactions.forEach((transaction, transactionIndex) => {
        convertedTransactions.push({
          id: `${fileIndex}-${transactionIndex}`,
          symbol: transaction.symbol,
          type: transaction.type,
          quantity: transaction.quantity,
          price: transaction.price,
          amount: transaction.amount,
          transaction_date: transaction.transaction_date,
          broker: transaction.broker,
          currency: transaction.currency,
          user_notes: transaction.user_notes,
          account: transaction.account,
          exchange: transaction.exchange,
        })
      })
    })

    setAllTransactions(convertedTransactions)
  }, [extractResults])

  // Helper functions
  const getFileStatus = (fileIndex: number): FileProcessingState['status'] => {
    if (!fileStates[fileIndex]) return 'completed'
    return fileStates[fileIndex].status
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'pending':
        return <ClockIcon className="w-4 h-4" />
      case 'processing':
        return <SpinnerIcon className="w-4 h-4" />
      case 'completed':
        return <CheckIcon className="w-4 h-4" />
      case 'error':
        return <XIcon className="w-4 h-4" />
      default:
        return <CheckIcon className="w-4 h-4" />
    }
  }

  const getTransactionCount = (fileIndex: number) => {
    if (!extractResults[fileIndex]) return 0
    return extractResults[fileIndex].transaction_count
  }

  const areAllFilesProcessed = () => {
    return fileStates.every((fs) => fs.status === 'completed' || fs.status === 'error')
  }

  // Get transactions for current file
  const getCurrentFileTransactions = (): ProcessedTransaction[] => {
    if (!extractResults[currentFileIndex]) return []

    const fileResult = extractResults[currentFileIndex]
    return fileResult.transactions.map((transaction, index) => ({
      id: `${currentFileIndex}-${index}`,
      symbol: transaction.symbol,
      type: transaction.type,
      quantity: transaction.quantity,
      price: transaction.price,
      amount: transaction.amount,
      transaction_date: transaction.transaction_date,
      broker: transaction.broker,
      currency: transaction.currency,
      user_notes: transaction.user_notes,
      account: transaction.account,
      exchange: transaction.exchange,
    }))
  }

  const currentTransactions = getCurrentFileTransactions()

  const handleCellEdit = (
    id: string,
    field: keyof ProcessedTransaction,
    value: string | number
  ) => {
    setAllTransactions((prev) => prev.map((t) => (t.id === id ? { ...t, [field]: value } : t)))
  }

  const handleAddRow = () => {
    const newTransaction: ProcessedTransaction = {
      id: Date.now().toString(),
      symbol: '',
      type: 'Buy',
      quantity: 0,
      price: 0,
      amount: 0,
      transaction_date: new Date().toISOString().split('T')[0],
      broker: '',
      currency: 'USD',
      user_notes: '',
      account: '',
      exchange: '',
    }
    setAllTransactions((prev) => [...prev, newTransaction])
  }

  const handleDeleteRows = () => {
    setShowDeleteConfirmation(true)
  }

  const confirmDeleteRows = () => {
    setAllTransactions((prev) => prev.filter((t) => !selectedRows.has(t.id)))
    setSelectedRows(new Set())
    setShowDeleteConfirmation(false)
  }

  const cancelDeleteRows = () => {
    setShowDeleteConfirmation(false)
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
    console.log('Importing transactions:', allTransactions)
    navigate(ROUTES.TRANSACTIONS, {
      state: { message: 'Transactions imported successfully!' },
    })
  }

  const handleCancel = () => {
    navigate(ROUTES.TRANSACTIONS)
  }

  return (
    <div className="container mx-auto p-6 space-y-6">
      {/* Page Title */}
      <div className="text-center mb-8">
        <Title as="h1" className="mb-4">
          Review & Confirm Transactions
        </Title>
        <p className="text-muted-foreground">
          {areAllFilesProcessed()
            ? 'All files processed. Please review the extracted data below and make any necessary corrections.'
            : 'Some data is ready, other files are still processing. You can start reviewing now.'}
        </p>
      </div>

      {/* File tabs */}
      <div className="flex flex-wrap gap-2 mb-6">
        {files.map((file, index) => {
          const status = getFileStatus(index)
          const count = getTransactionCount(index)
          const isActive = index === currentFileIndex
          const isClickable = status === 'completed' || 'error'

          return (
            <Button
              key={index}
              onClick={() => isClickable && setCurrentFileIndex(index)}
              disabled={!isClickable}
              variant={isActive ? 'default' : 'outline'}
              size="sm"
            >
              <div className="flex items-center gap-2">
                <span>{getStatusIcon(status)}</span>
                <span className="truncate max-w-[150px]">{file.name}</span>
                {status === 'completed' && count > 0 && (
                  <span className="text-xs">({count} records)</span>
                )}
                {status === 'processing' && <span className="text-xs">Processing...</span>}
                {status === 'error' && <span className="text-xs text-red-500">Error</span>}
              </div>
            </Button>
          )
        })}
      </div>

      {/* Content area */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Screenshot display */}
        <Card>
          <CardContent>
            <h3 className="text-lg font-semibold mb-4">Original Screenshot</h3>
            {files[currentFileIndex] && (
              <div className="bg-muted rounded-lg p-4 text-center">
                <img
                  src={URL.createObjectURL(files[currentFileIndex])}
                  alt={files[currentFileIndex].name}
                  className="max-w-full h-auto rounded-lg"
                />
              </div>
            )}
          </CardContent>
        </Card>

        {/* Transaction data */}
        <Card>
          <CardContent>
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-lg font-semibold">Extracted Data</h3>
              <div className="flex gap-2">
                <Button onClick={handleAddRow} variant="outline" size="sm">
                  <PlusIcon className="w-4 h-4 mr-2" />
                  Add Row
                </Button>
                {selectedRows.size > 0 && (
                  <Button
                    onClick={handleDeleteRows}
                    variant="outline"
                    size="sm"
                    className="text-red-600 hover:text-red-700"
                  >
                    <DeleteIcon className="w-4 h-4 mr-2" />
                    Delete Selected ({selectedRows.size})
                  </Button>
                )}
              </div>
            </div>

            {currentTransactions.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                {getFileStatus(currentFileIndex) === 'processing'
                  ? 'Processing file...'
                  : getFileStatus(currentFileIndex) === 'error'
                    ? 'Error processing this file'
                    : 'No transactions found in this file'}
              </div>
            ) : (
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead className="w-[50px]">Select</TableHead>
                      <TableHead>Date</TableHead>
                      <TableHead>Symbol</TableHead>
                      <TableHead>Type</TableHead>
                      <TableHead>Quantity</TableHead>
                      <TableHead>Price</TableHead>
                      <TableHead>Amount</TableHead>
                      <TableHead>Exchange</TableHead>
                      <TableHead>Currency</TableHead>
                      <TableHead>Notes</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {currentTransactions.map((transaction) => (
                      <TableRow key={transaction.id}>
                        <TableCell>
                          <input
                            type="checkbox"
                            checked={selectedRows.has(transaction.id)}
                            onChange={() => handleRowSelect(transaction.id)}
                            className="rounded"
                          />
                        </TableCell>
                        <TableCell>
                          <Input
                            type="date"
                            value={transaction.transaction_date}
                            onChange={(e) =>
                              handleCellEdit(transaction.id, 'transaction_date', e.target.value)
                            }
                            className="w-[150px]"
                          />
                        </TableCell>
                        <TableCell>
                          <Input
                            value={transaction.symbol}
                            onChange={(e) =>
                              handleCellEdit(transaction.id, 'symbol', e.target.value)
                            }
                            className="w-[100px]"
                          />
                        </TableCell>
                        <TableCell>
                          <Dropdown
                            trigger={
                              <Button variant="outline" size="sm">
                                {transaction.type}
                              </Button>
                            }
                          >
                            <DropdownItem
                              onClick={() => handleCellEdit(transaction.id, 'type', 'Buy')}
                            >
                              Buy
                            </DropdownItem>
                            <DropdownItem
                              onClick={() => handleCellEdit(transaction.id, 'type', 'Sell')}
                            >
                              Sell
                            </DropdownItem>
                            <DropdownItem
                              onClick={() => handleCellEdit(transaction.id, 'type', 'Dividends')}
                            >
                              Dividends
                            </DropdownItem>
                          </Dropdown>
                        </TableCell>
                        <TableCell>
                          <Input
                            type="number"
                            value={transaction.quantity}
                            onChange={(e) =>
                              handleCellEdit(
                                transaction.id,
                                'quantity',
                                parseFloat(e.target.value) || 0
                              )
                            }
                            className="w-[100px]"
                          />
                        </TableCell>
                        <TableCell>
                          <Input
                            type="number"
                            step="0.01"
                            value={transaction.price}
                            onChange={(e) =>
                              handleCellEdit(
                                transaction.id,
                                'price',
                                parseFloat(e.target.value) || 0
                              )
                            }
                            className="w-[100px]"
                          />
                        </TableCell>
                        <TableCell>
                          <span
                            className={`font-medium ${
                              transaction.type === 'Buy'
                                ? 'text-red-600'
                                : transaction.type === 'Sell'
                                  ? 'text-green-600'
                                  : 'text-blue-600'
                            }`}
                          >
                            ${transaction.amount.toFixed(2)}
                          </span>
                        </TableCell>
                        <TableCell>
                          <Input
                            value={transaction.exchange}
                            onChange={(e) =>
                              handleCellEdit(transaction.id, 'exchange', e.target.value)
                            }
                            className="w-[100px]"
                          />
                        </TableCell>
                        <TableCell>
                          <Input
                            value={transaction.currency}
                            onChange={(e) =>
                              handleCellEdit(transaction.id, 'currency', e.target.value)
                            }
                            className="w-[80px]"
                          />
                        </TableCell>
                        <TableCell>
                          <Input
                            value={transaction.user_notes}
                            onChange={(e) =>
                              handleCellEdit(transaction.id, 'user_notes', e.target.value)
                            }
                            className="w-[150px]"
                            placeholder="Add notes..."
                          />
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Action buttons */}
      <div className="flex justify-between pt-6">
        <Button onClick={handleCancel} variant="outline">
          Cancel
        </Button>
        <Button
          onClick={handleConfirmImport}
          disabled={!areAllFilesProcessed()}
          className="bg-green-600 hover:bg-green-700 text-white"
        >
          {areAllFilesProcessed()
            ? 'Confirm & Import'
            : 'Import (waiting for all files to complete)'}
        </Button>
      </div>

      {/* Delete confirmation modal */}
      {showDeleteConfirmation && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <Card className="w-[400px]">
            <CardContent className="p-6">
              <h3 className="text-lg font-semibold mb-4">Confirm Deletion</h3>
              <p className="text-muted-foreground mb-6">
                Are you sure you want to delete the selected rows? This action cannot be undone.
              </p>
              <div className="flex justify-end gap-3">
                <Button onClick={cancelDeleteRows} variant="outline">
                  Cancel
                </Button>
                <Button
                  onClick={confirmDeleteRows}
                  className="bg-red-600 hover:bg-red-700 text-white"
                >
                  Delete
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  )
}
