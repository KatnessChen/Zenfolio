import React from 'react'

interface SelectorOption {
  value: string
  label: string
}

interface SelectorProps {
  options: SelectorOption[]
  value: string
  onChange: (value: string) => void
  className?: string
  type?: 'horizontal' | 'vertical'
}

export const Selector: React.FC<SelectorProps> = ({
  options,
  value,
  onChange,
  className = '',
  type = 'horizontal',
}) => {
  return (
    <div
      className={`${
        type === 'vertical' ? 'flex flex-col items-stretch gap-1' : 'flex justify-center gap-1'
      } ${className}`}
    >
      {options.map((option) => (
        <button
          key={option.value}
          onClick={() => onChange(option.value)}
          className={`px-3 py-1.5 text-sm rounded-md transition-colors text-left ${
            value === option.value
              ? 'bg-muted text-foreground font-bold'
              : 'text-muted-foreground hover:text-foreground hover:bg-muted/50'
          }`}
        >
          {option.label}
        </button>
      ))}
    </div>
  )
}
