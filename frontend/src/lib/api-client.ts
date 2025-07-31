import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { API_ENDPOINTS } from '@/constants/api'
import { logger } from '@/lib/logger'

// API configuration
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
const API_VERSION = '/api/v1'

class ApiClient {
  // Centralized API call logger
  private logApiCall(method: string, url: string, dataOrConfig?: unknown, config?: unknown) {
    if (method === 'GET' || method === 'DELETE') {
      logger.info(`[API CALL] ${method}`, url, dataOrConfig)
    } else {
      logger.info(`[API CALL] ${method}`, url, dataOrConfig, config)
    }
  }
  private instance: AxiosInstance

  constructor() {
    this.instance = axios.create({
      baseURL: `${API_BASE_URL}${API_VERSION}`,
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    this.setupInterceptors()
  }

  private setupInterceptors() {
    // Request interceptor to add auth token
    this.instance.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('auth_token')
        if (token) {
          config.headers.Authorization = `Bearer ${token}`
        }
        return config
      },
      (error) => {
        return Promise.reject(error)
      }
    )

    // Response interceptor to handle auth errors
    this.instance.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          // Token expired or invalid
          this.clearAuthData()
          window.location.href = API_ENDPOINTS.AUTH.LOGIN
        }
        return Promise.reject(error)
      }
    )
  }

  private clearAuthData() {
    localStorage.removeItem('auth_token')
    localStorage.removeItem('auth_user')
  }

  // Generic request method
  public async request<T>(config: AxiosRequestConfig): Promise<AxiosResponse<T>> {
    this.logApiCall(config.method?.toUpperCase() || 'REQUEST', config.url || '', config)
    return this.instance.request<T>(config)
  }

  // Convenience methods
  public async get<T>(url: string, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> {
    this.logApiCall('GET', url, config)
    return this.instance.get<T>(url, config)
  }

  public async post<T>(
    url: string,
    data?: unknown,
    config?: AxiosRequestConfig
  ): Promise<AxiosResponse<T>> {
    this.logApiCall('POST', url, data, config)
    return this.instance.post<T>(url, data, config)
  }

  public async put<T>(
    url: string,
    data?: unknown,
    config?: AxiosRequestConfig
  ): Promise<AxiosResponse<T>> {
    this.logApiCall('PUT', url, data, config)
    return this.instance.put<T>(url, data, config)
  }

  public async delete<T>(url: string, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> {
    this.logApiCall('DELETE', url, config)
    return this.instance.delete<T>(url, config)
  }
}

// Export singleton instance
export const apiClient = new ApiClient()
