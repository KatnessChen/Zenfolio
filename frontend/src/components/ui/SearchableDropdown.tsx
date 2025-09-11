import React, { useState, useRef, useEffect, useCallback } from 'react'
import { Input } from '@/components/ui/input'
import { DropdownItem } from '@/components/ui/dropdown'
import { cn } from '@/lib/utils'
import { fuzzySearch } from '@/services/listing.service'

export interface SearchableDropdownOption {
  id: string | number
  label: string
  value: string
  subtitle?: string
}

interface SearchableDropdownProps {
  value: string
  onChange: (value: string) => void
  options: SearchableDropdownOption[]
  placeholder?: string
  className?: string
  disabled?: boolean
  noResultsText?: string
  loadingText?: string
  isLoading?: boolean
  onSearch?: (query: string) => void
  searchDelay?: number
  allowCustomValue?: boolean
}

export function SearchableDropdown({
  value,
  onChange,
  options,
  placeholder = 'Search...',
  className,
  disabled = false,
  noResultsText = 'No results found',
  loadingText = 'Searching...',
  isLoading = false,
  onSearch,
  searchDelay = 50,
  allowCustomValue = true,
}: SearchableDropdownProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [searchQuery, setSearchQuery] = useState(value)
  const [filteredOptions, setFilteredOptions] = useState<SearchableDropdownOption[]>(options)
  const dropdownRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  // Handle click outside to close dropdown
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  // Filter options or call external search function
  useEffect(() => {
    const delayedSearch = setTimeout(() => {
      // Local filtering with scored ranking using fuzzySearch
      const filtered = fuzzySearch(searchQuery, options, (option) => option.value)
      setFilteredOptions(filtered)
    }, searchDelay)

    return () => clearTimeout(delayedSearch)
  }, [searchQuery, options, onSearch, searchDelay])

  // Update search query when value prop changes (only if not actively searching)
  useEffect(() => {
    if (!isOpen) {
      setSearchQuery(value)
    }
  }, [value, isOpen])

  const handleInputChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const newValue = e.target.value
      setSearchQuery(newValue)

      if (allowCustomValue && newValue !== value) {
        onChange(newValue)
      }

      if (!isOpen && newValue.length >= 0) {
        setIsOpen(true)
      }
    },
    [onChange, isOpen, allowCustomValue, value]
  )

  const handleOptionSelect = useCallback(
    (option: SearchableDropdownOption) => {
      onChange(option.value)
      setSearchQuery(option.value)
      setIsOpen(false)
      inputRef.current?.blur()
    },
    [onChange]
  )

  const handleInputFocus = useCallback(() => {
    setIsOpen(true)
  }, [])

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if (e.key === 'Escape') {
        setIsOpen(false)
        inputRef.current?.blur()
      } else if (e.key === 'Enter' && filteredOptions.length > 0) {
        e.preventDefault()
        handleOptionSelect(filteredOptions[0])
      }
    },
    [filteredOptions, handleOptionSelect]
  )

  return (
    <div ref={dropdownRef} className={cn('relative', className)}>
      <Input
        ref={inputRef}
        type="text"
        value={searchQuery}
        onChange={handleInputChange}
        onFocus={handleInputFocus}
        onKeyDown={handleKeyDown}
        placeholder={placeholder}
        disabled={disabled}
        className="w-full"
      />

      {isOpen && (
        <div className="absolute top-full left-0 right-0 z-50 mt-1 bg-card border border-border/50 rounded-md shadow-lg backdrop-blur-sm bg-card/95 animate-in fade-in-0 zoom-in-95 duration-200 max-h-60 overflow-y-auto">
          <div className="p-1">
            {isLoading ? (
              <div className="px-3 py-2 text-sm text-muted-foreground">{loadingText}</div>
            ) : filteredOptions.length > 0 ? (
              filteredOptions.map((option) => (
                <DropdownItem
                  key={option.id}
                  onClick={() => handleOptionSelect(option)}
                  className="cursor-pointer hover:bg-primary/10 hover:text-primary"
                >
                  <div className="flex items-center justify-between w-full">
                    <span className="font-medium">{option.label}</span>
                    {option.subtitle && (
                      <span className="text-muted-foreground text-xs truncate ml-2">
                        {option.subtitle}
                      </span>
                    )}
                  </div>
                </DropdownItem>
              ))
            ) : searchQuery.length > 0 ? (
              <div className="px-3 py-2 text-sm text-muted-foreground">
                {noResultsText} "{searchQuery}"
              </div>
            ) : (
              <div className="px-3 py-2 text-sm text-muted-foreground">Type to search...</div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}

export default SearchableDropdown
