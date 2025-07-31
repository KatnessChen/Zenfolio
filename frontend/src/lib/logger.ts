// src/lib/logger.ts
// Production-grade logger utility for browser and Node (Vite/React)

const isDev = import.meta.env.DEV
const logLevel = import.meta.env.VITE_LOG_LEVEL || 'info'

function shouldLog(level: 'debug' | 'info' | 'warn' | 'error') {
  if (isDev) return true
  if (level === 'error' || level === 'warn') return true
  if (level === 'info' && logLevel === 'info') return true
  if (level === 'debug' && logLevel === 'debug') return true
  return false
}

export const logger = {
  debug: (...args: unknown[]) => {
    if (shouldLog('debug')) console.debug('[DEBUG]', ...args)
  },
  info: (...args: unknown[]) => {
    if (shouldLog('info')) console.info('[INFO]', ...args)
  },
  warn: (...args: unknown[]) => {
    if (shouldLog('warn')) console.warn('[WARN]', ...args)
  },
  error: (...args: unknown[]) => {
    if (shouldLog('error')) console.error('[ERROR]', ...args)
    // Optionally send to Sentry or remote error tracker here
  },
}
