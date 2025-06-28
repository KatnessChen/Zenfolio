import React, { useState } from 'react'
import { Card, CardContent } from './ui/card'
import { Input } from './ui/input'
import { EditIcon } from './icons/EditIcon'
import { DeleteIcon } from './icons/DeleteIcon'
import type { Transaction } from '../types'
import { formatDate } from '../utils'

interface TransactionCardProps {
  transaction: Transaction
  onEdit: (transaction: Transaction) => void
  onDelete: (id: string) => void
  onUpdateNotes?: (id: string, notes: string) => void
}

export function TransactionCard({
  transaction,
  onEdit,
  onDelete,
  onUpdateNotes,
}: TransactionCardProps) {
  const [isExpanded, setIsExpanded] = useState(false)
  const [isEditingNotes, setIsEditingNotes] = useState(false)
  const [notes, setNotes] = useState(transaction.userNotes || '')

  const handleCardClick = (e: React.MouseEvent) => {
    // Don't toggle if clicking on action buttons or input field
    if (
      (e.target as HTMLElement).closest('button') ||
      (e.target as HTMLElement).closest('input') ||
      (e.target as HTMLElement).closest('[data-action]')
    ) {
      return
    }
    setIsExpanded(!isExpanded)
  }

  const handleNotesSubmit = () => {
    if (onUpdateNotes) {
      onUpdateNotes(transaction.id, notes)
    }
    setIsEditingNotes(false)
  }

  const handleNotesKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleNotesSubmit()
    } else if (e.key === 'Escape') {
      setNotes(transaction.userNotes || '')
      setIsEditingNotes(false)
    }
  }

  // Trade Type colors based on wireframe
  const getTradeTypeColor = (tradeType: string) => {
    switch (tradeType.toLowerCase()) {
      case 'buy':
        return '#9DC0B2' // Brighter Sage Green
      case 'sell':
        return '#E6B9B3' // Brighter Soft Salmon Pink
      case 'dividend':
        return '#A0B0A6' // Medium Grey-Green
      default:
        return '#F5F5F5' // Pure Light Gray fallback
    }
  }

  // Amount formatting with sign
  const formatAmount = (amount: number, tradeType: string) => {
    const sign = tradeType.toLowerCase() === 'sell' ? '-' : ''
    return `${sign}${Math.abs(amount).toFixed(2)}`
  }

  return (
    <Card
      className="bg-[#3A4B40] border-none rounded-lg cursor-pointer transition-all duration-200 ease-in-out hover:bg-[#414D47] hover:shadow-lg hover:scale-[1.01] active:scale-[0.99] active:transition-transform active:duration-100"
      onClick={handleCardClick}
    >
      <CardContent className="p-4">
        {/* Row 1: Symbol + Trade Type */}
        <div className="flex justify-between items-center mb-2 transition-all duration-200 hover:translate-x-1">
          <span className="text-[#F5F5F5] font-semibold text-lg">{transaction.ticker}</span>
          <span
            className="font-medium text-sm transition-all duration-200"
            style={{ color: getTradeTypeColor(transaction.tradeType) }}
          >
            {transaction.tradeType}
          </span>
        </div>

        {/* Row 2: Trade Date + Amount */}
        <div className="flex justify-between items-center mb-2 transition-all duration-200 hover:translate-x-1">
          <span className="text-[#F5F5F5] text-sm">{formatDate(transaction.tradeDate)}</span>
          <span className="text-[#F5F5F5] text-sm font-medium text-right">
            {formatAmount(transaction.amount, transaction.tradeType)}
          </span>
        </div>

        {/* Row 3: Broker + Upload Date */}
        <div className="flex justify-between items-center mb-2">
          <span className="text-[#A0B0A6] text-xs">{transaction.broker}</span>
          <div className="flex items-center gap-2">
            <span className="text-[#A0B0A6] text-xs">{formatDate(transaction.uploadDate)}</span>
            {/* Expand/Collapse indicator */}
            <svg
              className={`w-4 h-4 text-[#A0B0A6] transition-transform duration-300 ease-in-out ${
                isExpanded ? 'rotate-180' : ''
              }`}
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 9l-7 7-7-7"
              />
            </svg>
          </div>
        </div>

        {/* Action Icons - Always visible */}
        <div className="flex justify-end gap-2 mt-3" data-action>
          <button
            onClick={(e) => {
              e.stopPropagation()
              onEdit(transaction)
            }}
            className="p-1 text-[#F5F5F5] hover:text-[#9DC0B2] transition-all duration-200 hover:scale-110"
            title="Edit transaction"
          >
            <EditIcon className="w-4 h-4" />
          </button>
          <button
            onClick={(e) => {
              e.stopPropagation()
              onDelete(transaction.id)
            }}
            className="p-1 text-[#F5F5F5] hover:text-[#E6B9B3] transition-all duration-200 hover:scale-110"
            title="Delete transaction"
          >
            <DeleteIcon className="w-4 h-4" />
          </button>
        </div>

        {/* Expanded Section with smooth animation */}
        <div
          className={`overflow-hidden transition-all duration-300 ease-in-out ${
            isExpanded ? 'max-h-96 opacity-100' : 'max-h-0 opacity-0'
          }`}
        >
          <div className="mt-4 pt-4 border-t border-[#5A6B61] space-y-3">
            {/* Price */}
            <div
              className={`flex justify-between transform transition-all duration-300 hover:translate-x-1 ${
                isExpanded ? 'translate-y-0 opacity-100 delay-75' : 'translate-y-2 opacity-0'
              }`}
            >
              <span className="text-[#A0B0A6] text-sm">Price:</span>
              <span className="text-[#F5F5F5] text-sm">${transaction.price.toFixed(2)}</span>
            </div>

            {/* Quantity */}
            <div
              className={`flex justify-between transform transition-all duration-300 hover:translate-x-1 ${
                isExpanded ? 'translate-y-0 opacity-100 delay-100' : 'translate-y-2 opacity-0'
              }`}
            >
              <span className="text-[#A0B0A6] text-sm">Quantity:</span>
              <span className="text-[#F5F5F5] text-sm">{transaction.quantity}</span>
            </div>

            {/* Currency */}
            <div
              className={`flex justify-between transform transition-all duration-300 hover:translate-x-1 ${
                isExpanded ? 'translate-y-0 opacity-100 delay-150' : 'translate-y-2 opacity-0'
              }`}
            >
              <span className="text-[#A0B0A6] text-sm">Currency:</span>
              <span className="text-[#F5F5F5] text-sm">{transaction.currency}</span>
            </div>

            {/* User Notes - Editable */}
            <div
              className={`space-y-2 transform transition-all duration-300 hover:translate-x-1 ${
                isExpanded ? 'translate-y-0 opacity-100 delay-200' : 'translate-y-2 opacity-0'
              }`}
            >
              <span className="text-[#A0B0A6] text-sm block">User Notes:</span>
              {isEditingNotes ? (
                <Input
                  value={notes}
                  onChange={(e) => setNotes(e.target.value)}
                  onBlur={handleNotesSubmit}
                  onKeyDown={handleNotesKeyDown}
                  placeholder="Add notes..."
                  className="bg-[#2A3530] border-[#5A6B61] text-[#F5F5F5] text-sm transition-all duration-200 focus:border-[#9DC0B2] focus:ring-1 focus:ring-[#9DC0B2]/20"
                  autoFocus
                  data-action
                />
              ) : (
                <div
                  onClick={(e) => {
                    e.stopPropagation()
                    setIsEditingNotes(true)
                  }}
                  className="text-[#F5F5F5] text-sm min-h-[32px] p-2 rounded border border-transparent hover:border-[#5A6B61] cursor-text transition-all duration-200 hover:bg-[#2A3530]/50"
                  data-action
                >
                  {transaction.userNotes || 'Click to add notes...'}
                </div>
              )}
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
