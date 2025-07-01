export const ROUTES = {
  HOME: '/',
  DASHBOARD: '/dashboard',
  SIGN_UP: '/sign-up',
  LOGIN: '/login',
  TRANSACTIONS: '/transactions',
  TRANSACTIONS_UPLOAD: '/transactions/upload',
  TRANSACTIONS_UPLOAD_PROCESSING: '/transactions/upload-processing',
  TRANSACTIONS_UPLOAD_REVIEW: '/transactions/upload-review',
  TRANSACTIONS_MANUAL_ADD: '/transactions/manual-add',
  TRANSACTIONS_MANUAL_REVIEW: '/transactions/manual-review',
  SETTINGS: '/settings',
  UI_DEMO: '/ui-demo',
  CONTACT: '/contact',
  LOGOUT: '/logout',
  ME: '/me',
} as const

// Type for route values
export type RouteValue = (typeof ROUTES)[keyof typeof ROUTES]
