// Meta data service for fetching currencies, brokers, and symbols
// TODO: Replace with real API calls when backend endpoints are available

export interface Currency {
  code: string
  name: string
}

export interface Broker {
  id: string
  name: string
}

export interface Symbol {
  symbol: string
  name: string
}

export async function fetchCurrencies(): Promise<Currency[]> {
  // TODO: Replace with real API call
  return [
    { code: 'USD', name: 'US Dollar' },
    { code: 'CAD', name: 'Canadian Dollar' },
  ]
}

export async function fetchBrokers(): Promise<Broker[]> {
  // TODO: Replace with real API call
  return [
    { id: 'fidelity', name: 'Fidelity' },
    { id: 'schwab', name: 'Charles Schwab' },
    { id: 'etrade', name: 'E*TRADE' },
    { id: 'td-ameritrade', name: 'TD Ameritrade' },
    { id: 'robinhood', name: 'Robinhood' },
    { id: 'interactive-brokers', name: 'Interactive Brokers' },
  ]
}

export async function fetchSymbols(): Promise<Symbol[]> {
  // TODO: Replace with real API call
  return [
    { symbol: 'AAPL', name: 'Apple Inc.' },
    { symbol: 'GOOGL', name: 'Alphabet Inc.' },
    { symbol: 'MSFT', name: 'Microsoft Corporation' },
    { symbol: 'TSLA', name: 'Tesla, Inc.' },
    { symbol: 'AMZN', name: 'Amazon.com, Inc.' },
    { symbol: 'NVDA', name: 'NVIDIA Corporation' },
  ]
}
