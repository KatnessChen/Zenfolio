// Common API response structure
export interface ApiResponse<T> {
  success: boolean
  data?: T
  message?: string
  error?: string
}

// Shared user type (used in multiple features)
export interface User {
  email: string
  firstName: string
  lastName?: string
}
