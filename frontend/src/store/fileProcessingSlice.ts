import { createSlice, type PayloadAction } from '@reduxjs/toolkit'
import type { FileProcessingState, ExtractResponseData } from '@/types'
import { FILE_STATUS_PENDING, FILE_STATUS_PROCESSING } from '@/constants/fileProcessingStatus'

interface FileProcessingSliceState {
  files: File[]
  fileStates: FileProcessingState[]
  extractResults: ExtractResponseData[]
  isProcessing: boolean
}

const initialState: FileProcessingSliceState = {
  files: [],
  fileStates: [],
  extractResults: [],
  isProcessing: false,
}

const fileProcessingSlice = createSlice({
  name: 'fileProcessing',
  initialState,
  reducers: {
    initializeFiles: (
      state,
      action: PayloadAction<{
        files: File[]
        fileStates?: FileProcessingState[]
        extractResults?: ExtractResponseData[]
      }>
    ) => {
      state.files = action.payload.files
      state.fileStates =
        action.payload.fileStates ||
        action.payload.files.map((file) => ({
          file,
          status: FILE_STATUS_PENDING,
          progress: 0,
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
        progress?: number
      }>
    ) => {
      const { fileIndex, status, result, error, progress } = action.payload

      if (state.fileStates[fileIndex]) {
        state.fileStates[fileIndex] = {
          ...state.fileStates[fileIndex],
          status,
          result,
          error,
          progress:
            progress ??
            (status === 'completed' || status === 'error'
              ? 100
              : state.fileStates[fileIndex].progress),
        }
      }

      // Update extract results if successful
      if (result && status === 'completed') {
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
      state.isProcessing = false
    },
  },
})

export const { initializeFiles, updateFileStatus, clearFiles } = fileProcessingSlice.actions
export default fileProcessingSlice.reducer
