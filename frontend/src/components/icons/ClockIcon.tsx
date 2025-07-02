import React from 'react'

interface ClockIconProps {
  className?: string
}

export const ClockIcon: React.FC<ClockIconProps> = ({ className = 'w-4 h-4' }) => {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
      xmlns="http://www.w3.org/2000/svg"
    >
      <circle cx="12" cy="12" r="10" />
      <path d="M12 6v6l4 2" />
    </svg>
  )
}
