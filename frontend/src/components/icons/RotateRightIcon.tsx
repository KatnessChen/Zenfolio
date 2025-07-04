import React from 'react'

interface RotateRightIconProps {
  size?: number
  className?: string
}

export const RotateRightIcon: React.FC<RotateRightIconProps> = ({ size = 24, className }) => {
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
      <path d="M21.5 2v6h-6" />
      <path d="M21.34 15.57a10 10 0 1 1-.57-8.38" />
    </svg>
  )
}
