import React from 'react'

interface SpinnerIconProps {
  size?: number
  className?: string
}

export const SpinnerIcon: React.FC<SpinnerIconProps> = ({ size = 24, className = '' }) => {
  return (
    <svg
      className={`${className} animate-spin`}
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      {/* Outer faded ring */}
      <circle className="opacity-20" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
      {/* Inner colored arc with rounded ends */}
      <path
        className="opacity-90"
        fill="none"
        stroke="currentColor"
        strokeWidth="4"
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M12 4
          a8 8 0 0 1 7.5 5.2"
      />
    </svg>
  )
}
