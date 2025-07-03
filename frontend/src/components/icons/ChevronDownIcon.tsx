import React from 'react'

interface ChevronDownIconProps {
  size?: number
  className?: string
}

export const ChevronDownIcon: React.FC<ChevronDownIconProps> = ({ size = 24, className }) => {
  return (
    <svg
      className={className}
      width={size}
      height={size}
      viewBox="0 0 16 16"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        d="M4 6L8 10L12 6"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}
