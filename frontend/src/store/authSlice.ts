import { createSlice, createAsyncThunk } from '@reduxjs/toolkit'
import type { PayloadAction } from '@reduxjs/toolkit'
import { AuthService } from '@/services/auth.service'
import type { AuthState, LoginRequest } from '@/types/auth'
import type { User } from '@/types'

// Initial state
const initialState: AuthState = {
  user: AuthService.getStoredUser(),
  token: AuthService.getToken(),
  isAuthenticated: AuthService.isAuthenticated(),
  isLoading: false,
  error: null,
}

// Async thunks
export const loginUser = createAsyncThunk(
  'auth/login',
  async (credentials: LoginRequest, { rejectWithValue }) => {
    try {
      const response = await AuthService.login(credentials)

      if (response.success) {
        // Save token and user data
        AuthService.saveToken(response.data.token)
        AuthService.saveUser(response.data.user)

        return {
          user: response.data.user,
          token: response.data.token,
        }
      } else {
        return rejectWithValue(response.message || 'Login failed')
      }
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'An unexpected error occurred'
      return rejectWithValue(errorMessage)
    }
  }
)

export const logoutUser = createAsyncThunk('auth/logout', async (_, { rejectWithValue }) => {
  try {
    await AuthService.logout()
  } catch (error: unknown) {
    const errorMessage = error instanceof Error ? error.message : 'Logout failed'
    return rejectWithValue(errorMessage)
  }
})

export const getCurrentUser = createAsyncThunk(
  'auth/getCurrentUser',
  async (_, { rejectWithValue }) => {
    try {
      const user = await AuthService.getCurrentUser()
      AuthService.saveUser(user)
      return user
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to get user'
      return rejectWithValue(errorMessage)
    }
  }
)

// Auth slice
const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null
    },
    clearAuth: (state) => {
      state.user = null
      state.token = null
      state.isAuthenticated = false
      state.error = null
      AuthService.clearAuthData()
    },
  },
  extraReducers: (builder) => {
    // Login
    builder
      .addCase(loginUser.pending, (state) => {
        state.isLoading = true
        state.error = null
      })
      .addCase(
        loginUser.fulfilled,
        (state, action: PayloadAction<{ user: User; token: string }>) => {
          state.isLoading = false
          state.user = action.payload.user
          state.token = action.payload.token
          state.isAuthenticated = true
          state.error = null
        }
      )
      .addCase(loginUser.rejected, (state, action) => {
        state.isLoading = false
        state.error = action.payload as string
        state.isAuthenticated = false
      })

    // Logout
    builder
      .addCase(logoutUser.pending, (state) => {
        state.isLoading = true
      })
      .addCase(logoutUser.fulfilled, (state) => {
        state.user = null
        state.token = null
        state.isAuthenticated = false
        state.isLoading = false
        state.error = null
      })
      .addCase(logoutUser.rejected, (state) => {
        state.user = null
        state.token = null
        state.isAuthenticated = false
        state.isLoading = false
        // Don't set error for logout failures
      })

    // Get current user
    builder
      .addCase(getCurrentUser.pending, (state) => {
        state.isLoading = true
      })
      .addCase(getCurrentUser.fulfilled, (state, action: PayloadAction<User>) => {
        state.isLoading = false
        state.user = action.payload
        state.error = null
      })
      .addCase(getCurrentUser.rejected, (state, action) => {
        state.isLoading = false
        state.error = action.payload as string
        // Clear auth if getting user fails (token might be invalid)
        state.user = null
        state.token = null
        state.isAuthenticated = false
        AuthService.clearAuthData()
      })
  },
})

export const { clearError, clearAuth } = authSlice.actions
export default authSlice.reducer
