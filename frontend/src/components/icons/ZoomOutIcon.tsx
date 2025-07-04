import React from 'react'

interface ZoomOutIconProps {
  size?: number
  className?: string
}

export const ZoomOutIcon: React.FC<ZoomOutIconProps> = ({ size = 24, className }) => {
  return (
    <svg
      className={className}
      width={size}
      height={size}
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      viewBox="0 0 24 24"
      xmlns="http://www.w3.org/2000/svg"
    >
      <circle cx="11" cy="11" r="8" />
      <path d="M21 21l-4.35-4.35" />
      <line x1="8" y1="11" x2="14" y2="11" />
    </svg>
  )
}
