import { configureStore } from '@reduxjs/toolkit'
import authReducer from './authSlice'
import fileProcessingReducer from './fileProcessingSlice'
import transactionHistoryReducer from './transactionHistorySlice'

export const store = configureStore({
  reducer: {
    auth: authReducer,
    fileProcessing: fileProcessingReducer,
    transactionHistory: transactionHistoryReducer,
  },
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
