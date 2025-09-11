// Currency constants for the application

export interface Currency {
  code: string
  name: string
}

export const CURRENCIES: Currency[] = [
  { code: 'USD', name: 'US Dollar' },
  { code: 'CAD', name: 'Canadian Dollar' },
  { code: 'TWD', name: 'New Taiwan Dollar' },
]
