import { apiClient } from '@/lib/api-client'
import { API_ENDPOINTS } from '@/constants/api'
import type { TransactionData, ExtractResponse } from '@/types'
import type { GetTransactionHistoryResponse } from '@/store/transactionHistorySlice'

export type TransactionDataRequest = Omit<TransactionData, 'id'>
export interface ImportTransactionsRequest {
  transactions: TransactionDataRequest[]
}

export interface ImportTransactionsResponse {
  success: boolean
  message: string
  data: {
    transactions: TransactionData[]
    count: number
  }
}

export class TransactionService {
  // Import transactions from file processing
  static async importTransactions(
    transactions: TransactionDataRequest[]
  ): Promise<ImportTransactionsResponse> {
    // Filter out fields that shouldn't be sent to the backend for create operations
    const filteredTransactions = transactions.map((transaction) => ({
      symbol: transaction.symbol,
      trade_type: transaction.trade_type,
      quantity: transaction.quantity,
      price: transaction.price,
      amount: transaction.amount,
      currency: transaction.currency,
      broker: transaction.broker,
      transaction_date: transaction.transaction_date,
      user_notes: transaction.user_notes,
      exchange: transaction.exchange,
    }))

    const payload: ImportTransactionsRequest = {
      transactions: filteredTransactions,
    }

    const response = await apiClient.post<ImportTransactionsResponse>(
      API_ENDPOINTS.TRANSACTIONS.HISTORY,
      payload
    )
    return response.data
  }

  // Get transaction history (for future use)
  static async getTransactionHistory(): Promise<GetTransactionHistoryResponse> {
    const response = await apiClient.get<{
      success: boolean
      message: string
      data: GetTransactionHistoryResponse
    }>(API_ENDPOINTS.TRANSACTIONS.HISTORY)
    return response.data.data
  }

  // Delete a single transaction
  static async deleteTransaction(id: string): Promise<{ success: boolean; message: string }> {
    const response = await apiClient.delete<{
      success: boolean
      message: string
      data?: { deleted_ids: string[] }
    }>(`${API_ENDPOINTS.TRANSACTIONS.HISTORY}/${id}`)
    return {
      success: response.data.success,
      message: response.data.message,
    }
  }

  // Delete multiple transactions in batch
  static async deleteTransactions(
    ids: string[]
  ): Promise<{ success: boolean; message: string; deletedIds: string[] }> {
    const response = await apiClient.delete<{
      success: boolean
      message: string
      data?: { deleted_ids: string[] }
    }>(`${API_ENDPOINTS.TRANSACTIONS.HISTORY}`, {
      data: { ids },
    })
    return {
      success: response.data.success,
      message: response.data.message,
      deletedIds: response.data.data?.deleted_ids || [],
    }
  }

  // Update a transaction (for future use)
  static async updateTransaction(
    transaction_id: string,
    transaction: TransactionData
  ): Promise<TransactionData> {
    // Remove unnecessary property to backend
    delete transaction.transaction_id

    const response = await apiClient.put<{
      success: boolean
      message: string
      data: { transaction: TransactionData }
    }>(`${API_ENDPOINTS.TRANSACTIONS.HISTORY}/${transaction_id}`, transaction)
    return response.data.data.transaction
  }

  /**
   * Extract transactions from a single image file
   * @param file - The image file to process
   * @returns Promise with extraction results
   */
  static async extractTransactions(file: File): Promise<ExtractResponse> {
    const formData = new FormData()
    formData.append('file', file)

    try {
      const response = await apiClient.post<ExtractResponse>(
        API_ENDPOINTS.TRANSACTIONS.EXTRACT,
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
          timeout: 60000, // 60 seconds for AI processing
        }
      )
      return response.data
    } catch (error: unknown) {
      console.error('Transaction extraction failed:', error)
      const errorMessage =
        error instanceof Error
          ? error.message
          : (error as { response?: { data?: { message?: string } } })?.response?.data?.message ||
            'Failed to extract transactions'
      return {
        success: false,
        message: errorMessage,
      }
    }
  }

  /**
   * Process multiple files in parallel
   * @param files - Array of files to process
   * @param onProgress - Callback for progress updates
   * @returns Promise with array of results
   */
  static async extractTransactionsParallel(
    files: File[],
    onProgress?: (fileIndex: number, result: ExtractResponse, error?: string) => void
  ): Promise<ExtractResponse[]> {
    const promises = files.map(async (file, index) => {
      try {
        const result = await this.extractTransactions(file)
        onProgress?.(index, result)
        return result
      } catch (error: unknown) {
        const errorResult: ExtractResponse = {
          success: false,
          message: `Failed to process ${file.name}: ${error instanceof Error ? error.message : 'Unknown error'}`,
        }
        onProgress?.(index, errorResult, error instanceof Error ? error.message : 'Unknown error')
        return errorResult
      }
    })
    return Promise.all(promises)
  }
}
