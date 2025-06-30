import React from 'react'

interface TitleProps {
  children: React.ReactNode
  as?: 'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6'
  className?: string
}

export const Title: React.FC<TitleProps> = ({ children, as = 'h1', className = '' }) => {
  const Tag = as
  return (
    <Tag
      className={`font-semibold text-foreground ${
        as === 'h1' ? 'text-3xl' : as === 'h2' ? 'text-2xl' : as === 'h3' ? 'text-xl' : ''
      } ${className}`}
    >
      {children}
    </Tag>
  )
}

export default Title
