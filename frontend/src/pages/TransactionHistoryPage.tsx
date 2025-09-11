import { useEffect, useState, useCallback, useRef } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { useNavigate } from 'react-router-dom'
import { ROUTES, TRADE_TYPE, CURRENCY } from '@/constants'
import type { TransactionData } from '@/types'
import {
  Card,
  CardContent,
  Button,
  Input,
  Label,
  Checkbox,
  Dropdown,
  DropdownItem,
  DropdownTrigger,
  ConfirmationModal,
} from '@/components/ui'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination'
import { EditIcon, DeleteIcon } from '@/components/icons'
import EmptyTransactionsState from '@/components/EmptyTransactionsState'
import { TransactionCard } from '@/components/TransactionCard'
import Title from '@/components/ui/title'
import { TransactionService } from '@/services/transaction.service'
import {
  fetchTransactionHistoryStart,
  fetchTransactionHistorySuccess,
  fetchTransactionHistoryFailure,
} from '@/store/transactionHistorySlice'
import type { GetTransactionHistoryResponse } from '@/store/transactionHistorySlice'
import type { RootState } from '@/store'
import { useToast } from '@/hooks/useToast'

type SortField = keyof TransactionData
type SortDirection = 'asc' | 'desc' | null

export default function TransactionHistoryPage() {
  const navigate = useNavigate()
  const dispatch = useDispatch()
  const { showToast } = useToast()
  const { transactions } = useSelector((state: RootState) => state.transactionHistory)
  const loading = useSelector((state: RootState) => state.transactionHistory.loading)
  const error = useSelector((state: RootState) => state.transactionHistory.error)

  const hasFetchedTransactions = useRef(false)

  const fetchTransactions = useCallback(() => {
    dispatch(fetchTransactionHistoryStart())
    TransactionService.getTransactionHistory()
      .then((data: GetTransactionHistoryResponse) => {
        dispatch(fetchTransactionHistorySuccess(data))
      })
      .catch((err) => {
        dispatch(
          fetchTransactionHistoryFailure(err?.message || 'Failed to fetch transaction history')
        )
        showToast({
          type: 'error',
          title: 'Failed to fetch transactions',
          message: err?.message || 'An error occurred while fetching transaction history.',
        })
      })
  }, [dispatch, showToast])

  useEffect(() => {
    if (hasFetchedTransactions.current) return

    fetchTransactions()

    hasFetchedTransactions.current = true
  }, [fetchTransactions, hasFetchedTransactions])

  const [sortField, setSortField] = useState<SortField>('transaction_date')
  const [sortDirection, setSortDirection] = useState<SortDirection>('desc')

  // Filter states
  const [dateRange, setDateRange] = useState('Last 30 Days')
  const [symbolFilter, setSymbolFilter] = useState('')
  const [currencyFilter, setCurrencyFilter] = useState('All')
  const [brokerFilter, setBrokerFilter] = useState('All')
  const [tradeTypeFilter, setTradeTypeFilter] = useState('All')

  // Pagination states
  const [currentPage, setCurrentPage] = useState(1)
  const itemsPerPage = 10
  const totalPages = Math.ceil(transactions.length / itemsPerPage)

  // Delete confirmation modal states
  const [deleteModalOpen, setDeleteModalOpen] = useState(false)
  const [transactionToDelete, setTransactionToDelete] = useState<TransactionData | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  // Multi-select states
  const [selectedTransactions, setSelectedTransactions] = useState<string[]>([])
  const [bulkDeleteModalOpen, setBulkDeleteModalOpen] = useState(false)

  // Get current page transactions
  const startIndex = (currentPage - 1) * itemsPerPage
  const endIndex = startIndex + itemsPerPage
  const currentTransactions = transactions
    .slice(startIndex, endIndex)
    .filter((transaction) => !!transaction.transaction_id)

  const handleSort = (field: SortField) => {
    if (sortField === field) {
      setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc')
    } else {
      setSortField(field)
      setSortDirection('desc')
    }
  }

  const handleApplyFilters = () => {
    // TODO: Implement filter logic
    console.log('Applying filters:', {
      dateRange,
      symbolFilter,
      currencyFilter,
      brokerFilter,
      tradeTypeFilter,
    })
  }

  const handleClearFilters = () => {
    setDateRange('Last 30 Days')
    setSymbolFilter('')
    setCurrencyFilter('All')
    setBrokerFilter('All')
    setTradeTypeFilter('All')
  }

  const handleEditTransaction = (transaction: TransactionData) => {
    navigate(ROUTES.TRANSACTIONS_MANUAL_ADD, { state: { initial: transaction } })
  }

  const handleDeleteTransaction = (transaction: TransactionData) => {
    setTransactionToDelete(transaction)
    setDeleteModalOpen(true)
  }

  const handleDeleteTransactionById = (id: string | undefined) => {
    const transaction = transactions.find((t) => t.transaction_id === id)
    if (transaction) {
      handleDeleteTransaction(transaction)
    }
  }

  const confirmDeleteTransaction = async () => {
    if (!transactionToDelete) return

    setIsDeleting(true)
    try {
      if (!transactionToDelete.transaction_id) {
        throw new Error('Transaction ID is missing')
      }

      const result = await TransactionService.deleteTransaction(transactionToDelete.transaction_id)

      if (result.success) {
        showToast({
          type: 'success',
          title: 'Transaction Deleted',
          message: result.message || 'Transaction has been deleted successfully',
          duration: 3000,
        })

        // Refresh the transaction list
        fetchTransactions()
      } else {
        showToast({
          type: 'error',
          title: 'Delete Failed',
          message: result.message || 'Failed to delete transaction',
          duration: 5000,
        })
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to delete transaction'
      showToast({
        type: 'error',
        title: 'Delete Failed',
        message: errorMessage,
        duration: 5000,
      })
    } finally {
      setIsDeleting(false)
      setDeleteModalOpen(false)
      setTransactionToDelete(null)
    }
  }

  const cancelDeleteTransaction = () => {
    setDeleteModalOpen(false)
    setTransactionToDelete(null)
  }

  // Multi-select handlers
  const handleSelectAll = (checked: boolean) => {
    if (checked) {
      const allIds = currentTransactions
        .map((t) => t.transaction_id)
        .filter((transaction_id): transaction_id is string => transaction_id !== undefined)
      setSelectedTransactions(allIds)
    } else {
      setSelectedTransactions([])
    }
  }

  const handleSelectTransaction = (transaction_id: string | undefined, checked: boolean) => {
    if (!transaction_id) return

    if (checked) {
      setSelectedTransactions((prev) => [...prev, transaction_id])
    } else {
      setSelectedTransactions((prev) => prev.filter((selectedId) => selectedId !== transaction_id))
    }
  }

  const handleBulkDelete = () => {
    if (selectedTransactions.length === 0) return
    setBulkDeleteModalOpen(true)
  }

  const confirmBulkDelete = async () => {
    if (selectedTransactions.length === 0) return

    setIsDeleting(true)
    try {
      const result = await TransactionService.deleteTransactions(selectedTransactions)

      if (result.success) {
        showToast({
          type: 'success',
          title: 'Transactions Deleted',
          message: `${selectedTransactions.length} transaction(s) deleted successfully`,
          duration: 3000,
        })

        // Refresh the transaction list
        fetchTransactions()
        setSelectedTransactions([])
      } else {
        showToast({
          type: 'error',
          title: 'Delete Failed',
          message: result.message || 'Failed to delete transactions',
          duration: 5000,
        })
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to delete transactions'
      showToast({
        type: 'error',
        title: 'Delete Failed',
        message: errorMessage,
        duration: 5000,
      })
    } finally {
      setIsDeleting(false)
      setBulkDeleteModalOpen(false)
    }
  }

  const cancelBulkDelete = () => {
    setBulkDeleteModalOpen(false)
  }

  // Optionally show loading and error states in the UI
  if (loading) {
    return (
      <div className="container mx-auto py-6 px-4">
        <Title as="h1">Transaction History</Title>
        <div className="flex justify-center items-center h-64">
          <span className="text-muted-foreground">Loading transactions...</span>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="container mx-auto py-6 px-4">
        <Title as="h1">Transaction History</Title>
        <div className="flex justify-center items-center h-64">
          <span className="text-destructive">{error}</span>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto py-6 px-4 space-y-8">
      {/* Page Title */}
      <Title as="h1">Transaction History</Title>

      {transactions.length === 0 ? (
        <EmptyTransactionsState onUploadClick={() => navigate(ROUTES.TRANSACTIONS_UPLOAD)} />
      ) : (
        <>
          {/* Filters Section */}
          <section className="space-y-3">
            <div className="flex flex-wrap gap-3 mb-4">
              {/* Date Range Filter */}
              <div>
                <Label htmlFor="date-range" className="text-sm text-muted-foreground">
                  Date Range
                </Label>
                <Dropdown
                  trigger={
                    <DropdownTrigger className="w-full justify-between border border-input">
                      <span className="truncate">{dateRange}</span>
                    </DropdownTrigger>
                  }
                >
                  <DropdownItem onClick={() => setDateRange('Last 7 Days')}>
                    Last 7 Days
                  </DropdownItem>
                  <DropdownItem onClick={() => setDateRange('Last 30 Days')}>
                    Last 30 Days
                  </DropdownItem>
                  <DropdownItem onClick={() => setDateRange('Last 3 Months')}>
                    Last 3 Months
                  </DropdownItem>
                  <DropdownItem onClick={() => setDateRange('Last Year')}>Last Year</DropdownItem>
                </Dropdown>
              </div>

              {/* Symbol Filter */}
              <div>
                <Label htmlFor="symbol" className="text-sm text-muted-foreground">
                  Symbol
                </Label>
                <Input
                  id="symbol"
                  placeholder="e.g. AAPL"
                  className="w-full"
                  value={symbolFilter}
                  onChange={(e) => setSymbolFilter(e.target.value)}
                />
              </div>

              {/* Currency Filter */}
              <div>
                <Label htmlFor="currency" className="text-sm text-muted-foreground">
                  Currency
                </Label>
                <Dropdown
                  trigger={
                    <DropdownTrigger className="w-full justify-between border border-input">
                      <span className="truncate">{currencyFilter}</span>
                    </DropdownTrigger>
                  }
                >
                  <DropdownItem onClick={() => setCurrencyFilter('All')}>All</DropdownItem>
                  <DropdownItem onClick={() => setCurrencyFilter(CURRENCY.USD)}>
                    {CURRENCY.USD}
                  </DropdownItem>
                  <DropdownItem onClick={() => setCurrencyFilter(CURRENCY.TWD)}>
                    {CURRENCY.TWD}
                  </DropdownItem>
                  <DropdownItem onClick={() => setCurrencyFilter(CURRENCY.CAD)}>
                    {CURRENCY.CAD}
                  </DropdownItem>
                </Dropdown>
              </div>

              {/* Broker Filter */}
              <div>
                <Label htmlFor="broker" className="text-sm text-muted-foreground">
                  Broker
                </Label>
                <Dropdown
                  trigger={
                    <DropdownTrigger className="w-full justify-between border border-input">
                      <span className="truncate">{brokerFilter}</span>
                    </DropdownTrigger>
                  }
                >
                  <DropdownItem onClick={() => setBrokerFilter('All')}>All</DropdownItem>
                  <DropdownItem onClick={() => setBrokerFilter('Fidelity')}>Fidelity</DropdownItem>
                  <DropdownItem onClick={() => setBrokerFilter('Schwab')}>Schwab</DropdownItem>
                  <DropdownItem onClick={() => setBrokerFilter('E*TRADE')}>E*TRADE</DropdownItem>
                  <DropdownItem onClick={() => setBrokerFilter('TD Ameritrade')}>
                    TD Ameritrade
                  </DropdownItem>
                </Dropdown>
              </div>

              {/* Trade Type Filter */}
              <div>
                <Label htmlFor="trade-type" className="text-sm text-muted-foreground">
                  Trade Type
                </Label>
                <Dropdown
                  trigger={
                    <DropdownTrigger className="w-full justify-between border border-input">
                      <span className="truncate">{tradeTypeFilter}</span>
                    </DropdownTrigger>
                  }
                >
                  <DropdownItem onClick={() => setTradeTypeFilter('All')}>All</DropdownItem>
                  <DropdownItem onClick={() => setTradeTypeFilter(TRADE_TYPE.BUY)}>
                    {TRADE_TYPE.BUY}
                  </DropdownItem>
                  <DropdownItem onClick={() => setTradeTypeFilter(TRADE_TYPE.SELL)}>
                    {TRADE_TYPE.SELL}
                  </DropdownItem>
                  <DropdownItem onClick={() => setTradeTypeFilter(TRADE_TYPE.DIVIDEND)}>
                    {TRADE_TYPE.DIVIDEND}
                  </DropdownItem>
                </Dropdown>
              </div>
            </div>

            {/* Filter Action Buttons */}
            <div className="flex gap-3 items-center">
              <Button variant="default" onClick={handleApplyFilters}>
                Apply Filters
              </Button>
              <Button variant="secondary" onClick={handleClearFilters}>
                Clear Filters
              </Button>
              {selectedTransactions.length > 0 && (
                <Button variant="destructive" onClick={handleBulkDelete} className="ml-4">
                  Delete Selected ({selectedTransactions.length})
                </Button>
              )}
            </div>
          </section>

          {/* Transaction History Table Section */}
          <section className="space-y-4">
            {/* Desktop Table View - Hidden on mobile */}
            <div className="hidden md:block">
              <Card className="bg-card">
                <CardContent className="p-0">
                  <div className="rounded-md overflow-hidden">
                    <Table>
                      <TableHeader>
                        <TableRow>
                          <TableHead className="w-[50px]">
                            <Checkbox
                              checked={
                                selectedTransactions.length === currentTransactions.length &&
                                currentTransactions.length > 0
                              }
                              indeterminate={
                                selectedTransactions.length > 0 &&
                                selectedTransactions.length < currentTransactions.length
                              }
                              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                                handleSelectAll(e.target.checked)
                              }
                            />
                          </TableHead>
                          <TableHead
                            sortable
                            sortDirection={sortField === 'transaction_date' ? sortDirection : null}
                            onSort={() => handleSort('transaction_date')}
                            className="min-w-[120px]"
                          >
                            Trade Date
                          </TableHead>
                          <TableHead
                            sortable
                            sortDirection={sortField === 'symbol' ? sortDirection : null}
                            onSort={() => handleSort('symbol')}
                          >
                            Symbol
                          </TableHead>
                          <TableHead
                            sortable
                            sortDirection={sortField === 'trade_type' ? sortDirection : null}
                            onSort={() => handleSort('trade_type')}
                          >
                            Trade Type
                          </TableHead>
                          <TableHead
                            sortable
                            sortDirection={sortField === 'price' ? sortDirection : null}
                            onSort={() => handleSort('price')}
                          >
                            Price
                          </TableHead>
                          <TableHead
                            sortable
                            sortDirection={sortField === 'quantity' ? sortDirection : null}
                            onSort={() => handleSort('quantity')}
                          >
                            Quantity
                          </TableHead>
                          <TableHead
                            sortable
                            sortDirection={sortField === 'amount' ? sortDirection : null}
                            onSort={() => handleSort('amount')}
                            className="text-right"
                          >
                            Amount
                          </TableHead>
                          <TableHead
                            sortable
                            sortDirection={sortField === 'currency' ? sortDirection : null}
                            onSort={() => handleSort('currency')}
                          >
                            Currency
                          </TableHead>
                          <TableHead
                            sortable
                            sortDirection={sortField === 'broker' ? sortDirection : null}
                            onSort={() => handleSort('broker')}
                          >
                            Broker
                          </TableHead>
                          <TableHead
                            sortable
                            sortDirection={sortField === 'user_notes' ? sortDirection : null}
                            onSort={() => handleSort('user_notes')}
                          >
                            Notes
                          </TableHead>
                          <TableHead>Action</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {currentTransactions.map((transaction) => (
                          <TableRow key={transaction.transaction_id}>
                            <TableCell>
                              <Checkbox
                                checked={
                                  !!transaction.transaction_id &&
                                  selectedTransactions.includes(transaction.transaction_id)
                                }
                                disabled={!transaction.transaction_id}
                                onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                                  handleSelectTransaction(
                                    transaction.transaction_id,
                                    e.target.checked
                                  )
                                }
                              />
                            </TableCell>
                            <TableCell>{transaction.transaction_date}</TableCell>
                            <TableCell className="font-medium">{transaction.symbol}</TableCell>
                            <TableCell
                              className={`font-medium ${
                                transaction.trade_type === TRADE_TYPE.BUY
                                  ? 'text-primary'
                                  : transaction.trade_type === TRADE_TYPE.SELL
                                    ? 'text-chart-3'
                                    : 'text-chart-1'
                              }`}
                            >
                              {transaction.trade_type}
                            </TableCell>
                            <TableCell className="text-right text-muted-foreground">
                              ${transaction.price}
                            </TableCell>
                            <TableCell className="text-right text-muted-foreground">
                              {transaction.trade_type === TRADE_TYPE.SELL ? '-' : ''}
                              {transaction.quantity}
                            </TableCell>
                            <TableCell className="text-right">
                              {transaction.trade_type === TRADE_TYPE.BUY ? '-' : ''}$
                              {transaction.amount}
                            </TableCell>
                            <TableCell className="text-right">{transaction.currency}</TableCell>
                            <TableCell>{transaction.broker}</TableCell>
                            <TableCell className="max-w-32 truncate">
                              {transaction.user_notes}
                            </TableCell>
                            <TableCell>
                              <div className="flex items-center gap-2">
                                <Button
                                  variant="outline"
                                  onClick={() => handleEditTransaction(transaction)}
                                  size="sm"
                                  className="h-8 w-8 p-0"
                                  title="Edit transaction"
                                >
                                  <EditIcon size={16} />
                                </Button>

                                <Button
                                  variant="outline"
                                  size="sm"
                                  onClick={() =>
                                    handleDeleteTransactionById(transaction.transaction_id)
                                  }
                                  className="h-8 w-8 p-0"
                                  title="Delete transaction"
                                >
                                  <DeleteIcon size={16} />
                                </Button>
                              </div>
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* Mobile Cards View - Visible only on mobile */}
            <div className="md:hidden space-y-3">
              {currentTransactions.map((transaction) => (
                <TransactionCard
                  key={transaction.transaction_id || transaction.id}
                  transaction={transaction}
                  onEdit={handleEditTransaction}
                  onDelete={() => handleDeleteTransactionById(transaction.transaction_id)}
                  isSelected={
                    !!transaction.transaction_id &&
                    selectedTransactions.includes(transaction.transaction_id)
                  }
                  onSelect={handleSelectTransaction}
                />
              ))}
            </div>

            {/* Pagination */}
            <Pagination>
              <PaginationContent>
                <PaginationItem>
                  <PaginationPrevious
                    href="#"
                    onClick={(e) => {
                      e.preventDefault()
                      if (currentPage > 1) setCurrentPage(currentPage - 1)
                    }}
                    style={{
                      opacity: currentPage === 1 ? 0.5 : 1,
                      pointerEvents: currentPage === 1 ? 'none' : 'auto',
                    }}
                  />
                </PaginationItem>

                {/* Generate page numbers */}
                {Array.from({ length: Math.min(totalPages, 5) }, (_, i) => {
                  const pageNumber = i + 1
                  return (
                    <PaginationItem key={pageNumber}>
                      <PaginationLink
                        href="#"
                        isActive={pageNumber === currentPage}
                        onClick={(e) => {
                          e.preventDefault()
                          setCurrentPage(pageNumber)
                        }}
                      >
                        {pageNumber}
                      </PaginationLink>
                    </PaginationItem>
                  )
                })}

                {totalPages > 5 && (
                  <PaginationItem>
                    <PaginationEllipsis />
                  </PaginationItem>
                )}

                <PaginationItem>
                  <PaginationNext
                    href="#"
                    onClick={(e) => {
                      e.preventDefault()
                      if (currentPage < totalPages) setCurrentPage(currentPage + 1)
                    }}
                    style={{
                      opacity: currentPage === totalPages ? 0.5 : 1,
                      pointerEvents: currentPage === totalPages ? 'none' : 'auto',
                    }}
                  />
                </PaginationItem>
              </PaginationContent>
            </Pagination>
          </section>
        </>
      )}

      {/* Delete Confirmation Modal */}
      <ConfirmationModal
        isOpen={deleteModalOpen}
        title="Delete Transaction"
        message={
          transactionToDelete ? (
            <>
              Are you sure you want to delete the transaction for{' '}
              <strong>{transactionToDelete.symbol}</strong> (
              <strong>{transactionToDelete.trade_type}</strong>)?
            </>
          ) : (
            'Are you sure you want to delete this transaction?'
          )
        }
        confirmText="Delete"
        cancelText="Cancel"
        confirmVariant="destructive"
        onConfirm={confirmDeleteTransaction}
        onCancel={cancelDeleteTransaction}
        isLoading={isDeleting}
      />

      {/* Bulk Delete Confirmation Modal */}
      <ConfirmationModal
        isOpen={bulkDeleteModalOpen}
        title="Delete Multiple Transactions"
        message={
          <>
            Are you sure you want to delete <strong>{selectedTransactions.length}</strong> selected
            transaction(s)?
          </>
        }
        confirmText={`Delete ${selectedTransactions.length} Transaction${selectedTransactions.length > 1 ? 's' : ''}`}
        cancelText="Cancel"
        confirmVariant="destructive"
        onConfirm={confirmBulkDelete}
        onCancel={cancelBulkDelete}
        isLoading={isDeleting}
      />
    </div>
  )
}
