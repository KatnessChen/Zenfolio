import { apiClient } from '@/lib/api-client'
import type { LoginRequest, LoginResponse, User } from '@/types'

export class AuthService {
  // Login user
  static async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await apiClient.post<LoginResponse>('/login', credentials)
    return response.data
  }

  // Logout user
  static async logout(): Promise<void> {
    try {
      await apiClient.post('/logout')
    } catch (error) {
      // Even if logout fails on server, we should clear local data
      console.error('Logout request failed:', error)
    } finally {
      // Always clear local storage
      this.clearAuthData()
    }
  }

  // Get current user info
  static async getCurrentUser(): Promise<User> {
    const response = await apiClient.get<{ success: boolean; data: User }>('/me')
    return response.data.data
  }

  // Token management
  static saveToken(token: string): void {
    localStorage.setItem('auth_token', token)
  }

  static getToken(): string | null {
    return localStorage.getItem('auth_token')
  }

  static clearToken(): void {
    localStorage.removeItem('auth_token')
  }

  // User data management
  static saveUser(user: User): void {
    localStorage.setItem('auth_user', JSON.stringify(user))
  }

  static getStoredUser(): User | null {
    const userData = localStorage.getItem('auth_user')
    return userData ? JSON.parse(userData) : null
  }

  static clearUser(): void {
    localStorage.removeItem('auth_user')
  }

  // Clear all auth data
  static clearAuthData(): void {
    this.clearToken()
    this.clearUser()
  }

  // Check if user is authenticated
  static isAuthenticated(): boolean {
    return !!this.getToken()
  }
}
