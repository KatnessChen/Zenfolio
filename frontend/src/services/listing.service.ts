// Meta data service for fetching currencies, brokers, and symbols
// TODO: Replace with real API calls when backend endpoints are available

import { SYMBOLS, type Symbol } from '@/constants/symbols'
import { CURRENCIES, type Currency } from '@/constants/currencies'
import { BROKERS, type Broker } from '@/constants/brokers'

export async function fetchCurrencies(): Promise<Currency[]> {
  // TODO: Replace with real API call
  return CURRENCIES
}

export async function fetchBrokers(): Promise<Broker[]> {
  // TODO: Replace with real API call
  return BROKERS
}

export async function fetchSymbols(): Promise<Symbol[]> {
  // TODO: Replace with real API call
  return SYMBOLS
}

/**
 * Generic fuzzy search function for any options with symbol-like values
 * @param query - Search query string
 * @param options - Array of options to search through
 * @param getSymbol - Function to extract symbol value from option
 * @returns Filtered and sorted array of options
 */
export function fuzzySearch<T>(query: string, options: T[], getValue: (option: T) => string): T[] {
  if (!query.trim()) return options

  const normalizedQuery = query.toLowerCase().trim()

  // Score each option based on match quality (symbol field only)
  const scoredResults = options.map((option) => {
    const valueLower = getValue(option).toLowerCase()
    let score = 0

    // Exact symbol match (highest priority)
    if (valueLower === normalizedQuery) {
      score = 1000
    }
    // Symbol starts with query
    else if (valueLower.startsWith(normalizedQuery)) {
      score = 900
    }
    // Symbol contains query
    else if (valueLower.includes(normalizedQuery)) {
      score = 800
    }
    // Fuzzy match - check if all characters of query exist in order in symbol
    else {
      let queryIndex = 0
      for (let i = 0; i < valueLower.length && queryIndex < normalizedQuery.length; i++) {
        if (valueLower[i] === normalizedQuery[queryIndex]) {
          queryIndex++
        }
      }
      if (queryIndex === normalizedQuery.length) {
        score = 500
      }
    }

    return { option, score }
  })

  // Filter out non-matches and sort by score
  return scoredResults
    .filter((result) => result.score > 0)
    .sort((a, b) => b.score - a.score)
    .map((result) => result.option)
    .slice(0, 10) // Limit to top 10 results
}
