import * as React from 'react'

import { cn } from '@/lib/utils'

function Input({ className, type, ...props }: React.ComponentProps<'input'>) {
  return (
    <input
      type={type}
      data-slot="input"
      className={cn(
        // Zenfolio Design: Dark Grey-Green background (#3A4B40) with Light Off-White text
        'bg-input text-foreground placeholder:text-muted-foreground',
        // Focus state: Brighter Sage Green border and ring
        'focus-visible:border-primary focus-visible:ring-primary/50 focus-visible:ring-2',
        // Base styling with subtle borders
        'border-border/40 flex h-10 w-full min-w-0 rounded-md border px-3 py-2 text-sm',
        // Transitions and states
        'transition-[color,border-color,box-shadow] outline-none',
        // File input styling
        'file:text-foreground file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-sm file:font-medium',
        // Selection styling with primary colors
        'selection:bg-primary selection:text-primary-foreground',
        // Disabled states
        'disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50',
        // Error states with Alert Red
        'aria-invalid:ring-destructive/30 aria-invalid:border-destructive',
        className
      )}
      {...props}
    />
  )
}

export { Input }
