import React from 'react'

interface RotateLeftIconProps {
  size?: number
  className?: string
}

export const RotateLeftIcon: React.FC<RotateLeftIconProps> = ({ size = 24, className }) => {
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
      <path d="M2.5 2v6h6" />
      <path d="M2.66 15.57a10 10 0 1 0 .57-8.38" />
    </svg>
  )
}
