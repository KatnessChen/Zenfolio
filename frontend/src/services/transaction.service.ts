import { apiClient } from '@/lib/api-client'
import { API_ENDPOINTS } from '@/constants/api'
import type { TransactionData, ExtractResponse } from '@/types'
import type { GetTransactionHistoryResponse } from '@/store/transactionHistorySlice'

export interface ImportTransactionsRequest {
  transactions: TransactionData[]
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
    transactions: TransactionData[]
  ): Promise<ImportTransactionsResponse> {
    const payload: ImportTransactionsRequest = {
      transactions: transactions.map((transaction) => ({
        id: transaction.id,
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
      })),
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

  // Delete a transaction (for future use)
  static async deleteTransaction(id: string): Promise<void> {
    await apiClient.delete(`${API_ENDPOINTS.TRANSACTIONS.HISTORY}/${id}`)
  }

  // Update a transaction (for future use)
  static async updateTransaction(
    id: string,
    transaction: Partial<TransactionData>
  ): Promise<TransactionData> {
    const response = await apiClient.put<{ success: boolean; data: TransactionData }>(
      `${API_ENDPOINTS.TRANSACTIONS.HISTORY}/${id}`,
      transaction
    )
    return response.data.data
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
