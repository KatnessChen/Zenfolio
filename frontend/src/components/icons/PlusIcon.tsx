import React from 'react'

interface PlusIconProps extends React.SVGAttributes<SVGElement> {
  size?: number
}

const PlusIcon: React.FC<PlusIconProps> = ({ size = 24, ...props }) => {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      {...props}
    >
      <path d="M12 5v14" />
      <path d="M5 12h14" />
    </svg>
  )
}

export { PlusIcon }
