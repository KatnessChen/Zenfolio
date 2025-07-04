import { createSlice, type PayloadAction } from '@reduxjs/toolkit'
import type { TransactionData } from '../types'

interface Pagination {
  page: number
  page_size: number
  total_records: number
  total_pages: number
  has_previous: boolean
  has_next: boolean
}

interface FiltersApplied {
  symbol?: string[]
  broker?: string[]
  exchange?: string[]
}

interface TransactionHistoryState {
  transactions: TransactionData[]
  pagination: Pagination
  filters_applied: FiltersApplied
  loading: boolean
  error: string | null
}

export interface GetTransactionHistoryResponse {
  transactions: TransactionData[]
  pagination: Pagination
  filters_applied: FiltersApplied
}

const initialState: TransactionHistoryState = {
  transactions: [],
  pagination: {
    page: 0,
    page_size: 0,
    total_records: 0,
    total_pages: 0,
    has_previous: false,
    has_next: false,
  },
  filters_applied: {
    symbol: undefined,
    broker: undefined,
    exchange: undefined,
  },
  loading: false,
  error: null,
}

const transactionHistorySlice = createSlice({
  name: 'transactionHistory',
  initialState,
  reducers: {
    fetchTransactionHistoryStart(state) {
      state.loading = true
      state.error = null
    },
    fetchTransactionHistorySuccess(state, action: PayloadAction<GetTransactionHistoryResponse>) {
      state.transactions = action.payload.transactions
      state.pagination = action.payload.pagination
      state.filters_applied = action.payload.filters_applied
      state.loading = false
      state.error = null
    },
    fetchTransactionHistoryFailure(state, action: PayloadAction<string>) {
      state.loading = false
      state.error = action.payload
    },
    clearTransactionHistory(state) {
      state.transactions = []
      state.loading = false
      state.error = null
    },
  },
})

export const {
  fetchTransactionHistoryStart,
  fetchTransactionHistorySuccess,
  fetchTransactionHistoryFailure,
  clearTransactionHistory,
} = transactionHistorySlice.actions

export default transactionHistorySlice.reducer
