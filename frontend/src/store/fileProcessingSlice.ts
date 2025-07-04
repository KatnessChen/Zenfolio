import { createSlice, type PayloadAction } from '@reduxjs/toolkit'
import type {
  FileProcessingState,
  ExtractResponseData,
  TransactionData,
  SerializableFile,
} from '@/types'
import {
  FILE_STATUS_COMPLETED,
  FILE_STATUS_PENDING,
  FILE_STATUS_PROCESSING,
} from '@/constants/fileProcessingStatus'

interface FileProcessingSliceState {
  files: SerializableFile[]
  fileStates: FileProcessingState[]
  extractResults: ExtractResponseData[]
  validationErrors: Record<string, string>
  isProcessing: boolean
}

const initialState: FileProcessingSliceState = {
  files: [],
  fileStates: [],
  extractResults: [],
  validationErrors: {},
  isProcessing: false,
}

const fileProcessingSlice = createSlice({
  name: 'fileProcessing',
  initialState,
  reducers: {
    initializeFiles: (
      state,
      action: PayloadAction<{
        files: SerializableFile[]
        fileStates?: FileProcessingState[]
        extractResults?: ExtractResponseData[]
      }>
    ) => {
      state.files = action.payload.files
      state.fileStates =
        action.payload.fileStates ||
        action.payload.files.map((file) => ({
          file: {
            name: file.name,
            size: file.size,
            type: file.type,
            lastModified: file.lastModified,
          } as File,
          status: FILE_STATUS_PENDING,
        }))
      state.extractResults = action.payload.extractResults || []
      state.isProcessing = state.fileStates.some(
        (fs) => fs.status === FILE_STATUS_PENDING || fs.status === FILE_STATUS_PROCESSING
      )
    },

    updateFileStatus: (
      state,
      action: PayloadAction<{
        fileIndex: number
        status: FileProcessingState['status']
        result?: ExtractResponseData
        error?: string
      }>
    ) => {
      const { fileIndex, status, result, error } = action.payload

      if (state.fileStates[fileIndex]) {
        state.fileStates[fileIndex] = {
          ...state.fileStates[fileIndex],
          status,
          result,
          error,
        }
      }

      // Update extract results if successful
      if (result && status === FILE_STATUS_COMPLETED) {
        state.extractResults[fileIndex] = result
      }

      // Update processing status
      state.isProcessing = state.fileStates.some(
        (fs) => fs.status === FILE_STATUS_PENDING || fs.status === FILE_STATUS_PROCESSING
      )
    },

    clearFiles: (state) => {
      state.files = []
      state.fileStates = []
      state.extractResults = []
      state.validationErrors = {}
      state.isProcessing = false
    },

    updateExtractedTransactions: (
      state,
      action: PayloadAction<{
        fileIndex: number
        transactions: TransactionData[]
      }>
    ) => {
      const { fileIndex, transactions } = action.payload

      if (state.extractResults[fileIndex]) {
        state.extractResults[fileIndex] = {
          ...state.extractResults[fileIndex],
          transactions,
          transaction_count: transactions.length,
        }
      }
    },

    setValidationError: (
      state,
      action: PayloadAction<{
        errorKey: string
        message: string
      }>
    ) => {
      const { errorKey, message } = action.payload
      state.validationErrors[errorKey] = message
    },

    clearValidationError: (
      state,
      action: PayloadAction<{
        errorKey: string
      }>
    ) => {
      const { errorKey } = action.payload
      delete state.validationErrors[errorKey]
    },

    clearAllValidationErrors: (state) => {
      state.validationErrors = {}
    },

    clearCurrentFile: (
      state,
      action: PayloadAction<{
        fileIndex: number
      }>
    ) => {
      const { fileIndex } = action.payload

      // Remove the file from files array
      state.files.splice(fileIndex, 1)

      // Remove the file state
      state.fileStates.splice(fileIndex, 1)

      // Remove the extract results
      state.extractResults.splice(fileIndex, 1)

      // Clear validation errors for this file
      const filePrefix = `file-${fileIndex}-`
      Object.keys(state.validationErrors).forEach((key) => {
        if (key.startsWith(filePrefix)) {
          delete state.validationErrors[key]
        }
      })

      // Update validation error keys for remaining files (shift indexes down)
      const newValidationErrors: Record<string, string> = {}
      Object.entries(state.validationErrors).forEach(([key, value]) => {
        const match = key.match(/^file-(\d+)-(.+)$/)
        if (match) {
          const keyFileIndex = parseInt(match[1])
          const restOfKey = match[2]

          if (keyFileIndex > fileIndex) {
            // Shift index down by 1
            const newKey = `file-${keyFileIndex - 1}-${restOfKey}`
            newValidationErrors[newKey] = value
          } else if (keyFileIndex < fileIndex) {
            // Keep the same key
            newValidationErrors[key] = value
          }
          // keyFileIndex === fileIndex entries are already deleted above
        } else {
          // Non-file-specific validation errors
          newValidationErrors[key] = value
        }
      })
      state.validationErrors = newValidationErrors

      // Update processing status
      state.isProcessing = state.fileStates.some(
        (fs) => fs.status === FILE_STATUS_PENDING || fs.status === FILE_STATUS_PROCESSING
      )
    },
  },
})

export const {
  initializeFiles,
  updateFileStatus,
  clearFiles,
  updateExtractedTransactions,
  setValidationError,
  clearValidationError,
  clearAllValidationErrors,
  clearCurrentFile,
} = fileProcessingSlice.actions
export default fileProcessingSlice.reducer
