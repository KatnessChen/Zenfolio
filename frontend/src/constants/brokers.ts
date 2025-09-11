// Broker constants for the application

export interface Broker {
  id: string
  name: string
}

export const BROKERS: Broker[] = [
  { id: 'fidelity', name: 'Fidelity' },
  { id: 'schwab', name: 'Charles Schwab' },
  { id: 'etrade', name: 'E*TRADE' },
  { id: 'td-ameritrade', name: 'TD Ameritrade' },
  { id: 'robinhood', name: 'Robinhood' },
  { id: 'interactive-brokers', name: 'Interactive Brokers' },
  { id: 'wealthsimple', name: 'Wealth Simple' },
  { id: 'masterlink', name: '元富' },
  { id: 'firstrade', name: 'Firstrade' },
  { id: 'others', name: 'Others' },
]
