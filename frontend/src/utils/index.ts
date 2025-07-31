// Format percent values with sign
export function formatPercent(value: number): string {
  const sign = value >= 0 ? '+' : ''
  return `${sign}${value.toFixed(2)}%`
}
import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

// API base URL configuration
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

// Format currency values
export function formatCurrency(amount: number, currency = 'USD'): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
    minimumFractionDigits: 0,
    maximumFractionDigits: 4,
  }).format(amount)
}

// Validate file types for image uploads
export function isValidImageFile(file: File): boolean {
  const validTypes = ['image/png', 'image/jpeg', 'image/gif', 'image/webp']
  return validTypes.includes(file.type)
}

// Generate unique IDs
export function generateId(): string {
  return Math.random().toString(36).substring(2) + Date.now().toString(36)
}

// Get absolute value of a number
export function abs(value: number): number {
  return Math.abs(value)
}

export function toSnakeCase<T extends object>(obj: T): Record<string, unknown> {
  if (typeof obj !== 'object' || obj === null) return obj
  if (Array.isArray(obj)) return obj.map(toSnakeCase) as unknown as Record<string, unknown>
  return Object.fromEntries(
    Object.entries(obj).map(([key, value]) => [
      key.replace(/([A-Z])/g, '_$1').toLowerCase(),
      typeof value === 'object' && value !== null ? toSnakeCase(value) : value,
    ])
  )
}

export function snakeToCamel<T>(obj: T): T {
  if (obj === null) return obj
  if (Array.isArray(obj)) {
    return obj.map(snakeToCamel) as unknown as T
  } else if (typeof obj === 'object') {
    return Object.entries(obj as Record<string, unknown>).reduce(
      (acc, [key, value]) => {
        const camelKey = key.replace(/_([a-z])/g, (_, c) => c.toUpperCase())
        ;(acc as Record<string, unknown>)[camelKey] = snakeToCamel(value)
        return acc
      },
      {} as Record<string, unknown>
    ) as T
  }
  return obj
}
